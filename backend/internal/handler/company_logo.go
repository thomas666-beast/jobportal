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

const maxLogoBytes = 2 * 1024 * 1024 // 2 MB

type CompanyLogoHandler struct {
	Queries *dbgen.Queries
	Storage storage.FileStorage
}

func NewCompanyLogoHandler(q *dbgen.Queries, s storage.FileStorage) *CompanyLogoHandler {
	return &CompanyLogoHandler{Queries: q, Storage: s}
}

func (h *CompanyLogoHandler) Upload(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	current, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "create your company first",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load company"})
	}

	upload, err := parseUpload(c, maxLogoBytes, logoAllowedMIME)
	if err != nil {
		return mapUploadError(c, err)
	}
	defer upload.Close()

	name, err := newStorageName(upload.SniffedMIME)
	if err != nil {
		return mapUploadError(c, err)
	}

	// Delete old file if present.
	if current.LogoUrl != nil && *current.LogoUrl != "" {
		_ = h.Storage.Delete(*current.LogoUrl)
	}

	url, err := h.Storage.Save("logos", name, upload.Reader)
	if err != nil {
		return mapUploadError(c, err)
	}

	company, err := h.Queries.SetCompanyLogo(c.Request().Context(), dbgen.SetCompanyLogoParams{
		OwnerUserID:  userID,
		LogoUrl:      &url,
		LogoFilename: &upload.Original,
	})
	if err != nil {
		_ = h.Storage.Delete(url)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save logo"})
	}

	return c.JSON(http.StatusOK, companyToResponse(company))
}

func (h *CompanyLogoHandler) Delete(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	current, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "company not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load company"})
	}

	if current.LogoUrl != nil && *current.LogoUrl != "" {
		_ = h.Storage.Delete(*current.LogoUrl)
	}

	company, err := h.Queries.ClearCompanyLogo(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to clear logo"})
	}

	return c.JSON(http.StatusOK, companyToResponse(company))
}
