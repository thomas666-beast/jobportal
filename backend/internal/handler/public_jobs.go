package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"regexp"
)

type PublicJobListItem struct {
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
	Company        struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Location *string `json:"location"`
		LogoURL  *string `json:"logo_url"`
	} `json:"company"`
}

type PublicJobDetail struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Requirements   *string  `json:"requirements"`
	Location       *string  `json:"location"`
	RemoteType     string   `json:"remote_type"`
	SalaryMin      *int32   `json:"salary_min"`
	SalaryMax      *int32   `json:"salary_max"`
	SalaryCurrency *string  `json:"salary_currency"`
	EmploymentType string   `json:"employment_type"`
	Category       *string  `json:"category"`
	Tags           []string `json:"tags"`
	PublishedAt    *string  `json:"published_at"`
	CreatedAt      string   `json:"created_at"`
	Company        struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
		Location    *string `json:"location"`
		LogoURL     *string `json:"logo_url"`
	} `json:"company"`
}

type PublicJobsHandler struct {
	Queries *dbgen.Queries
}

func NewPublicJobsHandler(q *dbgen.Queries) *PublicJobsHandler {
	return &PublicJobsHandler{Queries: q}
}

// tsQuerySanitizer strips characters that would break to_tsquery syntax
// and collapses whitespace.
var tsQuerySanitizer = regexp.MustCompile(`[^\p{L}\p{N}\s]+`)

func sanitizeSearchQuery(q string) string {
	q = tsQuerySanitizer.ReplaceAllString(q, " ")
	q = strings.Join(strings.Fields(q), " ")
	if len(q) > 100 {
		q = q[:100]
	}
	return q
}

const descSnippetLen = 200

func snippet(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return strings.TrimSpace(s[:n]) + "…"
}

// List handles GET /api/jobs?q=&location=&remote_type=&limit=&offset=
func (h *PublicJobsHandler) List(c echo.Context) error {
	q := sanitizeSearchQuery(strings.TrimSpace(c.QueryParam("q")))
	location := strings.TrimSpace(c.QueryParam("location"))
	remoteType := strings.TrimSpace(c.QueryParam("remote_type"))

	validRemote := map[string]bool{"remote": true, "hybrid": true, "on_site": true}
	if remoteType != "" && !validRemote[remoteType] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid remote_type",
		})
	}

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

	var qPtr, locPtr, rtPtr *string
	if q != "" {
		qPtr = &q
	}
	if location != "" {
		locPtr = &location
	}
	if remoteType != "" {
		rtPtr = &remoteType
	}

	rows, err := h.Queries.ListPublishedJobs(c.Request().Context(), dbgen.ListPublishedJobsParams{
		Q:          qPtr,
		Location:   locPtr,
		RemoteType: rtPtr,
		Lim:        int32(limit),
		Off:        int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list jobs",
		})
	}

	total, err := h.Queries.CountPublishedJobs(c.Request().Context(), dbgen.CountPublishedJobsParams{
		Q:          qPtr,
		Location:   locPtr,
		RemoteType: rtPtr,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to count jobs",
		})
	}

	out := make([]PublicJobListItem, 0, len(rows))
	for _, r := range rows {
		item := PublicJobListItem{
			ID:             r.ID.String(),
			Title:          r.Title,
			Description:    snippet(r.Description, descSnippetLen),
			Location:       r.Location,
			RemoteType:     r.RemoteType,
			SalaryMin:      r.SalaryMin,
			SalaryMax:      r.SalaryMax,
			SalaryCurrency: r.SalaryCurrency,
			EmploymentType: r.EmploymentType,
			Category:       r.Category,
			Tags:           r.Tags,
			PublishedAt:    tsPtr(r.PublishedAt),
		}
		if item.Tags == nil {
			item.Tags = []string{}
		}

		item.Company.ID = r.CompanyID.String()
		item.Company.Name = r.CompanyName
		item.Company.Location = r.CompanyLocation
		item.Company.LogoURL = r.CompanyLogoUrl
		out = append(out, item)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"jobs":  out,
		"total": total,
	})
}

// Detail handles GET /api/jobs/:id
func (h *PublicJobsHandler) Detail(c echo.Context) error {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid job id",
		})
	}

	r, err := h.Queries.GetPublishedJobByID(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "job not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch job",
		})
	}

	resp := PublicJobDetail{
		ID:             r.ID.String(),
		Title:          r.Title,
		Description:    r.Description,
		Requirements:   r.Requirements,
		Location:       r.Location,
		RemoteType:     r.RemoteType,
		SalaryMin:      r.SalaryMin,
		SalaryMax:      r.SalaryMax,
		SalaryCurrency: r.SalaryCurrency,
		EmploymentType: r.EmploymentType,
		Category:       r.Category,
		Tags:           r.Tags,
		PublishedAt:    tsPtr(r.PublishedAt),
		CreatedAt:      r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if resp.Tags == nil {
		resp.Tags = []string{}
	}

	resp.Company.ID = r.CompanyID.String()
	resp.Company.Name = r.CompanyName
	resp.Company.Description = r.CompanyDescription
	resp.Company.Website = r.CompanyWebsite
	resp.Company.Location = r.CompanyLocation
	resp.Company.LogoURL = r.CompanyLogoUrl

	return c.JSON(http.StatusOK, resp)
}

func tsPtr(t pgtype.Timestamptz) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format("2006-01-02T15:04:05Z07:00")
	return &s
}
