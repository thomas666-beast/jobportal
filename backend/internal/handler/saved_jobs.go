package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type SavedJobItem struct {
	JobID          string   `json:"job_id"`
	SavedAt        string   `json:"saved_at"`
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

type SavedJobsHandler struct {
	Queries *dbgen.Queries
}

func NewSavedJobsHandler(q *dbgen.Queries) *SavedJobsHandler {
	return &SavedJobsHandler{Queries: q}
}

// Save handles POST /api/seeker/saved-jobs/:job_id
func (h *SavedJobsHandler) Save(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	jobID, err := uuid.Parse(c.Param("job_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	// Confirm the job exists and is published.
	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch job"})
	}
	if job.Status != "published" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "only published jobs can be saved",
		})
	}

	if err := h.Queries.SaveJob(c.Request().Context(), dbgen.SaveJobParams{
		SeekerID: seekerID,
		JobID:    jobID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save job"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"saved":  true,
		"job_id": jobID.String(),
	})
}

// Unsave handles DELETE /api/seeker/saved-jobs/:job_id
func (h *SavedJobsHandler) Unsave(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	jobID, err := uuid.Parse(c.Param("job_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	if err := h.Queries.UnsaveJob(c.Request().Context(), dbgen.UnsaveJobParams{
		SeekerID: seekerID,
		JobID:    jobID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to unsave job"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"saved":  false,
		"job_id": jobID.String(),
	})
}

// Check handles GET /api/seeker/saved-jobs/:job_id — returns whether the job is saved.
func (h *SavedJobsHandler) Check(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	jobID, err := uuid.Parse(c.Param("job_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	saved, err := h.Queries.IsJobSaved(c.Request().Context(), dbgen.IsJobSavedParams{
		SeekerID: seekerID,
		JobID:    jobID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"saved":  saved,
		"job_id": jobID.String(),
	})
}

// List handles GET /api/seeker/saved-jobs
func (h *SavedJobsHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	rows, err := h.Queries.ListSavedJobsBySeeker(c.Request().Context(), seekerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list saved jobs"})
	}

	out := make([]SavedJobItem, 0, len(rows))
	for _, r := range rows {
		item := SavedJobItem{
			JobID:          r.JobID.String(),
			SavedAt:        r.SavedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			Title:          r.Title,
			Description:    r.Description,
			Location:       r.Location,
			RemoteType:     r.RemoteType,
			SalaryMin:      r.SalaryMin,
			SalaryMax:      r.SalaryMax,
			SalaryCurrency: r.SalaryCurrency,
			EmploymentType: r.EmploymentType,
			Category:       r.Category,
			Tags:           r.Tags,
		}
		if item.Tags == nil {
			item.Tags = []string{}
		}
		if r.PublishedAt.Valid {
			s := r.PublishedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			item.PublishedAt = &s
		}
		item.Company.ID = r.CompanyID.String()
		item.Company.Name = r.CompanyName
		item.Company.Location = r.CompanyLocation
		item.Company.LogoURL = r.CompanyLogoUrl
		out = append(out, item)
	}

	return c.JSON(http.StatusOK, out)
}
