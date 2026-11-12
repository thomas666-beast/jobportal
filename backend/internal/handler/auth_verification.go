package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/email"
	"job-portal/internal/service"
)

type AuthVerificationHandler struct {
	Queries *dbgen.Queries
	Tokens  *service.AuthTokenService
	Email   email.Sender
	BaseURL string
	AppName string
}

func NewAuthVerificationHandler(q *dbgen.Queries, tokens *service.AuthTokenService, e email.Sender, baseURL, appName string) *AuthVerificationHandler {
	return &AuthVerificationHandler{Queries: q, Tokens: tokens, Email: e, BaseURL: baseURL, AppName: appName}
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// VerifyEmail handles POST /api/auth/verify-email
func (h *AuthVerificationHandler) VerifyEmail(c echo.Context) error {
	req := new(VerifyEmailRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	req.Token = strings.TrimSpace(req.Token)
	if req.Token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "token is required"})
	}

	userID, err := h.Tokens.Consume(c.Request().Context(), req.Token, "email_verification")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenExpired):
			return c.JSON(http.StatusGone, map[string]string{"error": "verification link expired"})
		case errors.Is(err, service.ErrTokenUsed):
			return c.JSON(http.StatusGone, map[string]string{"error": "verification link already used"})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid verification link"})
		}
	}

	if _, err := h.Queries.SetUserEmailVerified(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to verify email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "email verified"})
}

// ResendVerification handles POST /api/auth/resend-verification (protected, rate-limited)
func (h *AuthVerificationHandler) ResendVerification(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	// If already verified, no-op.
	verified, err := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID)
	if err == nil && verified.Valid {
		return c.JSON(http.StatusOK, map[string]string{"message": "already verified"})
	}

	user, err := h.Queries.GetUserWithPasswordByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load user"})
	}

	token, err := h.Tokens.Create(c.Request().Context(), userID, "email_verification", 24*time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create token"})
	}

	msg := email.RenderVerificationEmail(h.AppName, h.BaseURL, token)
	if err := h.Email.Send(user.Email, msg.Subject, msg.HTML, msg.Text); err != nil {
		slog.Warn("resend verification email failed", "to", user.Email, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "verification email sent"})
}
