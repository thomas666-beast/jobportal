package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type MeHandler struct {
	Queries *dbgen.Queries
}

func NewMeHandler(q *dbgen.Queries) *MeHandler {
	return &MeHandler{Queries: q}
}

// ChangePassword handles POST /api/me/password
func (h *MeHandler) ChangePassword(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id",
		})
	}

	req := new(ChangePasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "current_password and new_password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "new_password must be at least 8 characters",
		})
	}

	if req.CurrentPassword == req.NewPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "new_password must be different from current_password",
		})
	}

	user, err := h.Queries.GetUserWithPasswordByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to load user",
		})
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.CurrentPassword),
	); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "current password is incorrect",
		})
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to hash password",
		})
	}

	if _, err := h.Queries.UpdateUserPassword(c.Request().Context(), dbgen.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: string(newHash),
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update password",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password updated",
	})
}
