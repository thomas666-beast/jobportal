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
	"job-portal/internal/service"
)

type RejectJobRequest struct {
	Reason string `json:"reason"`
}

type AdminJobsHandler struct {
	Queries       *dbgen.Queries
	Notifications *service.NotificationService
}

func NewAdminJobsHandler(q *dbgen.Queries, n *service.NotificationService) *AdminJobsHandler {
	return &AdminJobsHandler{Queries: q, Notifications: n}
}

// List handles GET /api/admin/jobs?status=pending_review&limit=50&offset=0
func (h *AdminJobsHandler) List(c echo.Context) error {
	status := c.QueryParam("status")
	if status == "" {
		status = "pending_review"
	}

	validStatuses := map[string]bool{
		"draft":          true,
		"pending_review": true,
		"published":      true,
		"closed":         true,
		"archived":       true,
	}
	if !validStatuses[status] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid status filter",
		})
	}

	limit := 50
	offset := 0
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	jobs, err := h.Queries.ListJobsByStatus(c.Request().Context(), dbgen.ListJobsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list jobs",
		})
	}

	total, err := h.Queries.CountJobsByStatus(c.Request().Context(), status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to count jobs",
		})
	}

	out := make([]JobResponse, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, jobToResponse(j))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"jobs":  out,
		"total": total,
	})
}

// Approve handles POST /api/admin/jobs/:id/approve
func (h *AdminJobsHandler) Approve(c echo.Context) error {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid job id",
		})
	}

	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
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

	if job.Status != "pending_review" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "only jobs in pending_review can be approved",
		})
	}

	updated, err := h.Queries.UpdateJobStatus(c.Request().Context(), dbgen.UpdateJobStatusParams{
		ID:              jobID,
		Status:          "published",
		RejectionReason: nil,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to approve job",
		})
	}

	// Look up employer via company.
	if company, cErr := h.Queries.GetCompanyByID(c.Request().Context(), job.CompanyID); cErr == nil {
		h.Notifications.NotifyJobApproved(
			c.Request().Context(),
			company.OwnerUserID,
			job.Title,
			jobID.String(),
		)
	}

	return c.JSON(http.StatusOK, jobToResponse(updated))
}

// Reject handles POST /api/admin/jobs/:id/reject
func (h *AdminJobsHandler) Reject(c echo.Context) error {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid job id",
		})
	}

	req := new(RejectJobRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	req.Reason = strings.TrimSpace(req.Reason)
	if req.Reason == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "reason is required",
		})
	}
	if len(req.Reason) > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "reason must be 1000 characters or fewer",
		})
	}

	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
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

	if job.Status != "pending_review" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "only jobs in pending_review can be rejected",
		})
	}

	reason := req.Reason
	updated, err := h.Queries.UpdateJobStatus(c.Request().Context(), dbgen.UpdateJobStatusParams{
		ID:              jobID,
		Status:          "draft",
		RejectionReason: &reason,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to reject job",
		})
	}

	return c.JSON(http.StatusOK, jobToResponse(updated))
}

type AdminJobDetailResponse struct {
	JobResponse
	Company struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
		Location    *string `json:"location"`
		LogoURL     *string `json:"logo_url"`
	} `json:"company"`
	Employer struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
	} `json:"employer"`
}

// Get handles GET /api/admin/jobs/:id
func (h *AdminJobsHandler) Get(c echo.Context) error {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid job id",
		})
	}

	row, err := h.Queries.GetJobWithCompanyByID(c.Request().Context(), jobID)
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

	resp := AdminJobDetailResponse{
		JobResponse: JobResponse{
			ID:              row.ID.String(),
			CompanyID:       row.CompanyID.String(),
			Title:           row.Title,
			Description:     row.Description,
			Requirements:    row.Requirements,
			Location:        row.Location,
			RemoteType:      row.RemoteType,
			SalaryMin:       row.SalaryMin,
			SalaryMax:       row.SalaryMax,
			SalaryCurrency:  row.SalaryCurrency,
			EmploymentType:  row.EmploymentType,
			Category:        row.Category,
			Tags:            row.Tags,
			Status:          row.Status,
			RejectionReason: row.RejectionReason,
			CreatedAt:       row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	if row.Tags == nil {
		resp.Tags = []string{}
	}
	if row.PublishedAt.Valid {
		s := row.PublishedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.PublishedAt = &s
	}

	resp.Company.ID = row.CompanyID.String()
	resp.Company.Name = row.CompanyName
	resp.Company.Description = row.CompanyDescription
	resp.Company.Website = row.CompanyWebsite
	resp.Company.Location = row.CompanyLocation
	resp.Company.LogoURL = row.CompanyLogoUrl

	resp.Employer.Username = row.EmployerUsername
	resp.Employer.Email = row.EmployerEmail

	return c.JSON(http.StatusOK, resp)
}
