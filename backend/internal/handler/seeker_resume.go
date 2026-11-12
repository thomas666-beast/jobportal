package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
	"job-portal/internal/storage"
)

const maxResumeBytes = 5 * 1024 * 1024 // 5 MB

type SeekerResumeHandler struct {
	Queries *dbgen.Queries
	Storage storage.FileStorage
}

func NewSeekerResumeHandler(q *dbgen.Queries, s storage.FileStorage) *SeekerResumeHandler {
	return &SeekerResumeHandler{Queries: q, Storage: s}
}

func (h *SeekerResumeHandler) Upload(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	// Ensure profile exists.
	current, err := h.Queries.GetSeekerProfile(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "create your profile first",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load profile"})
	}

	upload, err := parseUpload(c, maxResumeBytes, resumeAllowedMIME)
	if err != nil {
		return mapUploadError(c, err)
	}
	defer upload.Reader.(interface{ Close() error }).Close()

	name, err := newStorageName(upload.SniffedMIME)
	if err != nil {
		return mapUploadError(c, err)
	}

	// Delete old file if it exists.
	if current.ResumeUrl != nil && *current.ResumeUrl != "" {
		_ = h.Storage.Delete(*current.ResumeUrl)
	}

	url, err := h.Storage.Save("resumes", name, upload.Reader)
	if err != nil {
		return mapUploadError(c, err)
	}

	profile, err := h.Queries.SetSeekerResume(c.Request().Context(), dbgen.SetSeekerResumeParams{
		UserID:         userID,
		ResumeUrl:      &url,
		ResumeFilename: &upload.Original,
	})
	if err != nil {
		// Clean up on DB failure.
		_ = h.Storage.Delete(url)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save resume"})
	}

	return c.JSON(http.StatusOK, profileToResponse(profile))
}

func (h *SeekerResumeHandler) Delete(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	current, err := h.Queries.GetSeekerProfile(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load profile"})
	}

	if current.ResumeUrl != nil && *current.ResumeUrl != "" {
		_ = h.Storage.Delete(*current.ResumeUrl)
	}

	profile, err := h.Queries.ClearSeekerResume(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to clear resume"})
	}

	return c.JSON(http.StatusOK, profileToResponse(profile))
}
