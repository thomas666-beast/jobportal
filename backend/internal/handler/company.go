package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type CompanyResponse struct {
	ID          string  `json:"id"`
	OwnerUserID string  `json:"owner_user_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
	Location    *string `json:"location"`
	LogoURL     *string `json:"logo_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Location    string `json:"location"`
	LogoURL     string `json:"logo_url"`
}

type CompanyHandler struct {
	Queries *dbgen.Queries
}

func NewCompanyHandler(q *dbgen.Queries) *CompanyHandler {
	return &CompanyHandler{Queries: q}
}

func (h *CompanyHandler) Get(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id in token",
		})
	}

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "company not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch company",
		})
	}

	return c.JSON(http.StatusOK, companyToResponse(company))
}

func (h *CompanyHandler) Upsert(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id in token",
		})
	}

	req := new(CompanyRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
	}
	if len(req.Name) > 200 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "name must be 200 characters or fewer",
		})
	}

	description := strings.TrimSpace(req.Description)
	if len(description) > 5000 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "description must be 5000 characters or fewer",
		})
	}

	website := strings.TrimSpace(req.Website)
	if len(website) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "website must be 255 characters or fewer",
		})
	}

	location := strings.TrimSpace(req.Location)
	if len(location) > 120 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "location must be 120 characters or fewer",
		})
	}

	logoURL := strings.TrimSpace(req.LogoURL)
	if len(logoURL) > 500 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "logo_url must be 500 characters or fewer",
		})
	}

	company, err := h.Queries.UpsertCompany(c.Request().Context(), dbgen.UpsertCompanyParams{
		OwnerUserID: ownerID,
		Name:        req.Name,
		Description: toNullable(description),
		Website:     toNullable(website),
		Location:    toNullable(location),
		LogoUrl:     toNullable(logoURL),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to save company",
		})
	}

	return c.JSON(http.StatusOK, companyToResponse(company))
}

func companyToResponse(input interface{}) CompanyResponse {
	b, err := json.Marshal(input)
	if err != nil {
		return CompanyResponse{}
	}
	var c dbgen.Company
	if err := json.Unmarshal(b, &c); err != nil {
		return CompanyResponse{}
	}

	return CompanyResponse{
		ID:          c.ID.String(),
		OwnerUserID: c.OwnerUserID.String(),
		Name:        c.Name,
		Description: c.Description,
		Website:     c.Website,
		Location:    c.Location,
		LogoURL:     c.LogoUrl,
		CreatedAt:   c.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   c.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}
