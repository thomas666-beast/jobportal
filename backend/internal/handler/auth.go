package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/email"
	"job-portal/internal/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID              string `json:"user_id"`
	Username            string `json:"username"`
	Role                string `json:"role"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
	EmailVerified       bool   `json:"email_verified"`
}

type OnboardingRequest struct {
	Role string `json:"role"`
}

type AuthHandler struct {
	Queries       *dbgen.Queries
	JWT           *service.JWTService
	RefreshTokens *service.RefreshTokenService
	Tokens        *service.AuthTokenService
	Email         email.Sender
	BaseURL       string
	AppName       string
}

func NewAuthHandler(
	q *dbgen.Queries,
	jwtSvc *service.JWTService,
	refreshSvc *service.RefreshTokenService,
	tokens *service.AuthTokenService,
	e email.Sender,
	baseURL, appName string,
) *AuthHandler {
	return &AuthHandler{
		Queries:       q,
		JWT:           jwtSvc,
		RefreshTokens: refreshSvc,
		Tokens:        tokens,
		Email:         e,
		BaseURL:       baseURL,
		AppName:       appName,
	}
}

// issueSession creates a fresh access + refresh + CSRF set of cookies.
func (h *AuthHandler) issueSession(c echo.Context, userID uuid.UUID, username, role string, onboardingCompleted bool) error {
	accessToken, err := h.JWT.Generate(userID, username, role, onboardingCompleted)
	if err != nil {
		return err
	}

	refreshToken, err := h.RefreshTokens.Create(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	csrfToken, err := service.GenerateCSRFToken()
	if err != nil {
		return err
	}

	// Access token TTL — read from same env as JWT service uses.
	accessTTL := 15 * time.Minute
	service.SetAccessCookie(c, accessToken, accessTTL)
	service.SetRefreshCookie(c, refreshToken, h.RefreshTokens.TTL())
	service.SetCSRFCookie(c, csrfToken, h.RefreshTokens.TTL())

	return nil
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username, email and password are required"})
	}
	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "password must be at least 8 characters"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	user, err := h.Queries.CreateUser(c.Request().Context(), dbgen.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "job_seeker",
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return c.JSON(http.StatusConflict, map[string]string{"error": "email already registered"})
			case "users_username_key":
				return c.JSON(http.StatusConflict, map[string]string{"error": "username already taken"})
			default:
				return c.JSON(http.StatusConflict, map[string]string{"error": "user already exists"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	// Verification email (best-effort).
	if token, err := h.Tokens.Create(c.Request().Context(), user.ID, "email_verification", 24*time.Hour); err == nil {
		msg := email.RenderVerificationEmail(h.AppName, h.BaseURL, token)
		if sendErr := h.Email.Send(user.Email, msg.Subject, msg.HTML, msg.Text); sendErr != nil {
			slog.Warn("verification email failed", "to", user.Email, "error", sendErr)
		}
	} else {
		slog.Warn("verification token create failed", "user_id", user.ID.String(), "error", err)
	}

	// Auto-login on register — issue session cookies.
	if err := h.issueSession(c, user.ID, user.Username, user.Role, user.OnboardingCompleted); err != nil {
		slog.Warn("issue session on register failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "session setup failed"})
	}

	return c.JSON(http.StatusCreated, UserResponse{
		UserID:              user.ID.String(),
		Username:            user.Username,
		Role:                user.Role,
		OnboardingCompleted: user.OnboardingCompleted,
		EmailVerified:       false,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password are required"})
	}

	user, err := h.Queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to look up user"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
	}

	if err := h.issueSession(c, user.ID, user.Username, user.Role, user.OnboardingCompleted); err != nil {
		slog.Warn("issue session on login failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "session setup failed"})
	}

	verified, _ := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), user.ID)

	return c.JSON(http.StatusOK, UserResponse{
		UserID:              user.ID.String(),
		Username:            user.Username,
		Role:                user.Role,
		OnboardingCompleted: user.OnboardingCompleted,
		EmailVerified:       verified.Valid,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Best-effort revoke of the refresh token in the DB.
	if cookie, err := c.Cookie(service.RefreshCookieName); err == nil && cookie.Value != "" {
		h.RefreshTokens.Revoke(c.Request().Context(), cookie.Value)
	}
	service.ClearAuthCookies(c)
	return c.JSON(http.StatusOK, map[string]string{"status": "logged out"})
}

// Refresh issues a new access token using a valid refresh token.
// Called automatically by the frontend when the access token expires.
func (h *AuthHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie(service.RefreshCookieName)
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no refresh token"})
	}

	userID, err := h.RefreshTokens.Validate(c.Request().Context(), cookie.Value)
	if err != nil {
		service.ClearAuthCookies(c)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}

	// Look up the user to build new claims.
	user, err := h.Queries.GetUserWithPasswordByID(c.Request().Context(), userID)
	if err != nil {
		service.ClearAuthCookies(c)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	accessToken, err := h.JWT.Generate(user.ID, user.Username, user.Role, user.OnboardingCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
	}

	// Rotate the CSRF token on every refresh (cheap, extra defense).
	csrfToken, err := service.GenerateCSRFToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "csrf generation failed"})
	}

	service.SetAccessCookie(c, accessToken, 15*time.Minute)
	service.SetCSRFCookie(c, csrfToken, h.RefreshTokens.TTL())

	verified, _ := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID)

	return c.JSON(http.StatusOK, UserResponse{
		UserID:              user.ID.String(),
		Username:            user.Username,
		Role:                user.Role,
		OnboardingCompleted: user.OnboardingCompleted,
		EmailVerified:       verified.Valid,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}

	userID, _ := uuid.Parse(claims.UserID)
	var emailVerified bool
	if row, err := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID); err == nil {
		emailVerified = row.Valid
	}

	return c.JSON(http.StatusOK, UserResponse{
		UserID:              claims.UserID,
		Username:            claims.Username,
		Role:                claims.Role,
		OnboardingCompleted: claims.OnboardingCompleted,
		EmailVerified:       emailVerified,
	})
}

func (h *AuthHandler) OnboardingRole(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	if claims.OnboardingCompleted {
		return c.JSON(http.StatusConflict, map[string]string{"error": "onboarding already completed"})
	}

	req := new(OnboardingRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if req.Role != "job_seeker" && req.Role != "employer" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "role must be 'job_seeker' or 'employer'"})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	updated, err := h.Queries.UpdateUserRole(c.Request().Context(), dbgen.UpdateUserRoleParams{
		ID:   userID,
		Role: req.Role,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update role"})
	}

	// Re-issue access token with the new role/onboarding state. Refresh token
	// stays the same — the session continues across the role change.
	accessToken, err := h.JWT.Generate(updated.ID, updated.Username, updated.Role, updated.OnboardingCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
	}
	service.SetAccessCookie(c, accessToken, 15*time.Minute)

	verified, _ := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID)

	return c.JSON(http.StatusOK, UserResponse{
		UserID:              updated.ID.String(),
		Username:            updated.Username,
		Role:                updated.Role,
		OnboardingCompleted: updated.OnboardingCompleted,
		EmailVerified:       verified.Valid,
	})
}
