package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type CreateApplicationRequest struct {
	JobID       string `json:"job_id"`
	CoverLetter string `json:"cover_letter"`
}

type SeekerApplicationResponse struct {
	ID          string  `json:"id"`
	JobID       string  `json:"job_id"`
	CoverLetter *string `json:"cover_letter"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Job         struct {
		Title      string  `json:"title"`
		Location   *string `json:"location"`
		RemoteType string  `json:"remote_type"`
		Status     string  `json:"status"`
	} `json:"job"`
	Company struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		LogoURL *string `json:"logo_url"`
	} `json:"company"`
}

type SeekerApplicationsHandler struct {
	Queries       *dbgen.Queries
	Notifications *service.NotificationService
}

func NewSeekerApplicationsHandler(q *dbgen.Queries, n *service.NotificationService) *SeekerApplicationsHandler {
	return &SeekerApplicationsHandler{Queries: q, Notifications: n}
}

func (h *SeekerApplicationsHandler) Create(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)

	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	// Unverified users cannot apply.
	verified, vErr := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), seekerID)
	if vErr != nil || !verified.Valid {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "please verify your email before applying",
		})
	}

	req := new(CreateApplicationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	jobID, err := uuid.Parse(strings.TrimSpace(req.JobID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job_id"})
	}

	cover := strings.TrimSpace(req.CoverLetter)
	if len(cover) > 5000 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "cover_letter must be 5000 characters or fewer",
		})
	}

	// Seeker must have a profile.
	if _, err := h.Queries.GetSeekerProfile(c.Request().Context(), seekerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "complete your profile before applying",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check profile"})
	}

	// Job must exist and be published.
	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch job"})
	}
	if job.Status != "published" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "job is not accepting applications",
		})
	}

	app, err := h.Queries.CreateApplication(c.Request().Context(), dbgen.CreateApplicationParams{
		JobID:       jobID,
		SeekerID:    seekerID,
		CoverLetter: toNullable(cover),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "you have already applied to this job",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create application"})
	}

	// Notify the employer who owns the job.
	company, cErr := h.Queries.GetCompanyByID(c.Request().Context(), job.CompanyID)
	if cErr == nil {
		seekerName := claims.Username
		if profile, pErr := h.Queries.GetSeekerProfile(c.Request().Context(), seekerID); pErr == nil && profile.FullName != "" {
			seekerName = profile.FullName
		}
		h.Notifications.NotifyApplicationSubmitted(
			c.Request().Context(),
			company.OwnerUserID,
			seekerName,
			job.Title,
			app.ID.String(),
		)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":           app.ID.String(),
		"job_id":       app.JobID.String(),
		"cover_letter": app.CoverLetter,
		"status":       app.Status,
		"created_at":   app.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":   app.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *SeekerApplicationsHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	rows, err := h.Queries.ListApplicationsBySeeker(c.Request().Context(), seekerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list applications"})
	}

	out := make([]SeekerApplicationResponse, 0, len(rows))
	for _, r := range rows {
		item := SeekerApplicationResponse{
			ID:          r.ID.String(),
			JobID:       r.JobID.String(),
			CoverLetter: r.CoverLetter,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   r.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
		item.Job.Title = r.JobTitle
		item.Job.Location = r.JobLocation
		item.Job.RemoteType = r.JobRemoteType
		item.Job.Status = r.JobStatus
		item.Company.ID = r.CompanyID.String()
		item.Company.Name = r.CompanyName
		item.Company.LogoURL = r.CompanyLogoUrl
		out = append(out, item)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *SeekerApplicationsHandler) Withdraw(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid application id"})
	}

	app, err := h.Queries.GetApplicationByID(c.Request().Context(), appID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch application"})
	}

	if app.SeekerID != seekerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "not your application"})
	}

	if app.Status != "submitted" && app.Status != "viewed" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "only submitted or viewed applications can be withdrawn",
		})
	}

	updated, err := h.Queries.UpdateApplicationStatus(c.Request().Context(), dbgen.UpdateApplicationStatusParams{
		ID:     appID,
		Status: "withdrawn",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to withdraw"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id":     updated.ID.String(),
		"status": updated.Status,
	})
}

// HasApplied handles GET /api/seeker/applications/status?job_id=...
// Returns { "applied": true|false, "status": "..." } so the job detail page
// can show "Applied" without a second query.
func (h *SeekerApplicationsHandler) HasApplied(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	jobID, err := uuid.Parse(c.QueryParam("job_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job_id"})
	}

	app, err := h.Queries.GetApplicationByJobAndSeeker(c.Request().Context(), dbgen.GetApplicationByJobAndSeekerParams{
		JobID:    jobID,
		SeekerID: seekerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"applied": false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"applied": true,
		"status":  app.Status,
	})
}
