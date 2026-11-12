package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
)

type CompanyListItem struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	Website           *string `json:"website"`
	Location          *string `json:"location"`
	LogoURL           *string `json:"logo_url"`
	PublishedJobCount int64   `json:"published_job_count"`
}

type CompanyPublicDetail struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       *string            `json:"description"`
	Website           *string            `json:"website"`
	Location          *string            `json:"location"`
	LogoURL           *string            `json:"logo_url"`
	CreatedAt         string             `json:"created_at"`
	PublishedJobCount int64              `json:"published_job_count"`
	Jobs              []CompanyPublicJob `json:"jobs"`
}

type CompanyPublicJob struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Location       *string  `json:"location"`
	RemoteType     string   `json:"remote_type"`
	SalaryMin      *int32   `json:"salary_min"`
	SalaryMax      *int32   `json:"salary_max"`
	SalaryCurrency *string  `json:"salary_currency"`
	EmploymentType string   `json:"employment_type"`
	Category       *string  `json:"category"`
	Tags           []string `json:"tags"`
	PublishedAt    *string  `json:"published_at"`
}

type PublicCompaniesHandler struct {
	Queries *dbgen.Queries
}

func NewPublicCompaniesHandler(q *dbgen.Queries) *PublicCompaniesHandler {
	return &PublicCompaniesHandler{Queries: q}
}

func (h *PublicCompaniesHandler) List(c echo.Context) error {
	q := strings.TrimSpace(c.QueryParam("q"))

	limit := 20
	offset := 0
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	var qPtr *string
	if q != "" {
		qPtr = &q
	}

	rows, err := h.Queries.ListCompanies(c.Request().Context(), dbgen.ListCompaniesParams{
		Q:   qPtr,
		Lim: int32(limit),
		Off: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list companies",
		})
	}

	total, err := h.Queries.CountCompanies(c.Request().Context(), qPtr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to count companies",
		})
	}

	out := make([]CompanyListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, CompanyListItem{
			ID:                r.ID.String(),
			Name:              r.Name,
			Description:       r.Description,
			Website:           r.Website,
			Location:          r.Location,
			LogoURL:           r.LogoUrl,
			PublishedJobCount: r.PublishedJobCount,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"companies": out,
		"total":     total,
	})
}

func (h *PublicCompaniesHandler) Detail(c echo.Context) error {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid company id"})
	}

	company, err := h.Queries.GetCompanyPublicByID(c.Request().Context(), companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "company not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch company"})
	}

	jobs, err := h.Queries.ListPublishedJobsByCompany(c.Request().Context(), companyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch jobs"})
	}

	resp := CompanyPublicDetail{
		ID:                company.ID.String(),
		Name:              company.Name,
		Description:       company.Description,
		Website:           company.Website,
		Location:          company.Location,
		LogoURL:           company.LogoUrl,
		CreatedAt:         company.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		PublishedJobCount: company.PublishedJobCount,
		Jobs:              make([]CompanyPublicJob, 0, len(jobs)),
	}

	for _, j := range jobs {
		item := CompanyPublicJob{
			ID:             j.ID.String(),
			Title:          j.Title,
			Description:    j.Description,
			Location:       j.Location,
			RemoteType:     j.RemoteType,
			SalaryMin:      j.SalaryMin,
			SalaryMax:      j.SalaryMax,
			SalaryCurrency: j.SalaryCurrency,
			EmploymentType: j.EmploymentType,
			Category:       j.Category,
			Tags:           j.Tags,
		}
		if item.Tags == nil {
			item.Tags = []string{}
		}
		if j.PublishedAt.Valid {
			s := j.PublishedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			item.PublishedAt = &s
		}
		resp.Jobs = append(resp.Jobs, item)
	}

	return c.JSON(http.StatusOK, resp)
}
