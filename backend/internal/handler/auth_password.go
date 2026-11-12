package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/email"
	"job-portal/internal/service"
)

type AuthPasswordHandler struct {
	Queries *dbgen.Queries
	Tokens  *service.AuthTokenService
	Email   email.Sender
	BaseURL string
	AppName string
}

func NewAuthPasswordHandler(q *dbgen.Queries, tokens *service.AuthTokenService, e email.Sender, baseURL, appName string) *AuthPasswordHandler {
	return &AuthPasswordHandler{Queries: q, Tokens: tokens, Email: e, BaseURL: baseURL, AppName: appName}
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ForgotPassword handles POST /api/auth/forgot-password
// Always returns 200 to prevent account enumeration.
func (h *AuthPasswordHandler) ForgotPassword(c echo.Context) error {
	req := new(ForgotPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required"})
	}

	// Always return 200, even if user not found.
	defer func() {
		_ = c.JSON(http.StatusOK, map[string]string{
			"message": "if an account exists, a reset email has been sent",
		})
	}()

	user, err := h.Queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		slog.Warn("forgot-password query failed", "error", err)
		return nil
	}

	token, err := h.Tokens.Create(c.Request().Context(), user.ID, "password_reset", time.Hour)
	if err != nil {
		slog.Warn("forgot-password token create failed", "error", err)
		return nil
	}

	msg := email.RenderPasswordResetEmail(h.AppName, h.BaseURL, token)
	if err := h.Email.Send(user.Email, msg.Subject, msg.HTML, msg.Text); err != nil {
		slog.Warn("forgot-password email failed", "to", user.Email, "error", err)
	}

	return nil
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ResetPassword handles POST /api/auth/reset-password
func (h *AuthPasswordHandler) ResetPassword(c echo.Context) error {
	req := new(ResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.Token = strings.TrimSpace(req.Token)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if req.Token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "token is required"})
	}
	if len(req.NewPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "new_password must be at least 8 characters"})
	}

	userID, err := h.Tokens.Consume(c.Request().Context(), req.Token, "password_reset")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenExpired):
			return c.JSON(http.StatusGone, map[string]string{"error": "reset link expired"})
		case errors.Is(err, service.ErrTokenUsed):
			return c.JSON(http.StatusGone, map[string]string{"error": "reset link already used"})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid reset link"})
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	if err := h.Queries.SetUserPasswordByID(c.Request().Context(), dbgen.SetUserPasswordByIDParams{
		ID:           userID,
		PasswordHash: string(hash),
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "password updated"})
}
