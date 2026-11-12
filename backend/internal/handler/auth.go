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

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token               string `json:"token"`
	UserID              string `json:"user_id"`
	Username            string `json:"username"`
	Role                string `json:"role"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
}

type OnboardingRequest struct {
	Role string `json:"role"`
}

type AuthHandler struct {
	Queries *dbgen.Queries
	JWT     *service.JWTService
	Tokens  *service.AuthTokenService
	Email   email.Sender
	BaseURL string
	AppName string
}

func NewAuthHandler(q *dbgen.Queries, jwtSvc *service.JWTService, tokens *service.AuthTokenService, e email.Sender, baseURL, appName string) *AuthHandler {
	return &AuthHandler{
		Queries: q,
		JWT:     jwtSvc,
		Tokens:  tokens,
		Email:   e,
		BaseURL: baseURL,
		AppName: appName,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "username, email and password are required",
		})
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password must be at least 8 characters",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to hash password",
		})
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
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "email already registered",
				})
			case "users_username_key":
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "username already taken",
				})
			default:
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "user already exists",
				})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}

	// Generate verification token and send email. Failures here must not
	// fail registration — the user can resend from the banner.
	if token, err := h.Tokens.Create(c.Request().Context(), user.ID, "email_verification", 24*time.Hour); err == nil {
		msg := email.RenderVerificationEmail(h.AppName, h.BaseURL, token)
		if sendErr := h.Email.Send(user.Email, msg.Subject, msg.HTML, msg.Text); sendErr != nil {
			slog.Warn("verification email failed", "to", user.Email, "error", sendErr)
		}
	} else {
		slog.Warn("verification token create failed", "user_id", user.ID.String(), "error", err)
	}

	return c.JSON(http.StatusCreated, RegisterResponse{
		Message: "registration successful",
		UserID:  user.ID.String(),
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email and password are required",
		})
	}

	user, err := h.Queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to look up user",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid email or password",
		})
	}

	token, err := h.JWT.Generate(user.ID, user.Username, user.Role, user.OnboardingCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to issue token",
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token:               token,
		UserID:              user.ID.String(),
		Username:            user.Username,
		Role:                user.Role,
		OnboardingCompleted: user.OnboardingCompleted,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	userID, _ := uuid.Parse(claims.UserID)
	var emailVerified bool
	if row, err := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID); err == nil {
		emailVerified = row.Valid
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":              claims.UserID,
		"username":             claims.Username,
		"role":                 claims.Role,
		"onboarding_completed": claims.OnboardingCompleted,
		"email_verified":       emailVerified,
	})
}

func (h *AuthHandler) OnboardingRole(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	if claims.OnboardingCompleted {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "onboarding already completed",
		})
	}

	req := new(OnboardingRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Role != "job_seeker" && req.Role != "employer" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "role must be 'job_seeker' or 'employer'",
		})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id in token",
		})
	}

	updated, err := h.Queries.UpdateUserRole(c.Request().Context(), dbgen.UpdateUserRoleParams{
		ID:   userID,
		Role: req.Role,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update role",
		})
	}

	newToken, err := h.JWT.Generate(updated.ID, updated.Username, updated.Role, updated.OnboardingCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to issue token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":                newToken,
		"user_id":              updated.ID.String(),
		"username":             updated.Username,
		"role":                 updated.Role,
		"onboarding_completed": updated.OnboardingCompleted,
	})
}
