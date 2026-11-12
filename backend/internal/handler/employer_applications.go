package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type EmployerApplicationListItem struct {
	ID          string  `json:"id"`
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	CoverLetter *string `json:"cover_letter"`
	CreatedAt   string  `json:"created_at"`
	Job         struct {
		Title string `json:"title"`
	} `json:"job"`
	Seeker struct {
		ID       string  `json:"id"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		FullName *string `json:"full_name"`
		Headline *string `json:"headline"`
		Location *string `json:"location"`
	} `json:"seeker"`
}

type EmployerApplicationDetail struct {
	ID          string  `json:"id"`
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	CoverLetter *string `json:"cover_letter"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	JobTitle    string  `json:"job_title"`
	Seeker      struct {
		ID             string  `json:"id"`
		Username       string  `json:"username"`
		Email          string  `json:"email"`
		FullName       *string `json:"full_name"`
		Phone          *string `json:"phone"`
		Location       *string `json:"location"`
		Headline       *string `json:"headline"`
		Bio            *string `json:"bio"`
		ResumeURL      *string `json:"resume_url"`
		ResumeFilename *string `json:"resume_filename"`
	} `json:"seeker"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status"`
}

type EmployerApplicationsHandler struct {
	Queries       *dbgen.Queries
	Notifications *service.NotificationService
}

func NewEmployerApplicationsHandler(q *dbgen.Queries, n *service.NotificationService) *EmployerApplicationsHandler {
	return &EmployerApplicationsHandler{Queries: q, Notifications: n}
}

// loadOwnedApplication checks that the employer's company owns the job this
// application belongs to.
func (h *EmployerApplicationsHandler) loadOwnedApplication(c echo.Context, appID, companyID uuid.UUID) (*dbgen.GetApplicationDetailForEmployerRow, error) {
	row, err := h.Queries.GetApplicationDetailForEmployer(c.Request().Context(), appID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
			return nil, err
		}
		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch application"})
		return nil, err
	}
	if row.CompanyID != companyID {
		_ = c.JSON(http.StatusForbidden, map[string]string{"error": "not your application"})
		return nil, errors.New("not owned")
	}
	return &row, nil
}

func (h *EmployerApplicationsHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company first"})
	}

	rows, err := h.Queries.ListApplicationsByCompany(c.Request().Context(), company.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list applications"})
	}

	out := make([]EmployerApplicationListItem, 0, len(rows))
	for _, r := range rows {
		item := EmployerApplicationListItem{
			ID:          r.ID.String(),
			JobID:       r.JobID.String(),
			Status:      r.Status,
			CoverLetter: r.CoverLetter,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
		item.Job.Title = r.JobTitle
		item.Seeker.ID = r.SeekerID.String()
		item.Seeker.Username = r.SeekerUsername
		item.Seeker.Email = r.SeekerEmail
		item.Seeker.FullName = r.SeekerFullName
		item.Seeker.Headline = r.SeekerHeadline
		item.Seeker.Location = r.SeekerLocation
		out = append(out, item)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *EmployerApplicationsHandler) ListByJob(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company first"})
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
	}
	if job.CompanyID != company.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "not your job"})
	}

	rows, err := h.Queries.ListApplicationsByJob(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list applications"})
	}

	out := make([]EmployerApplicationListItem, 0, len(rows))
	for _, r := range rows {
		item := EmployerApplicationListItem{
			ID:          r.ID.String(),
			JobID:       r.JobID.String(),
			Status:      r.Status,
			CoverLetter: r.CoverLetter,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
		item.Job.Title = job.Title
		item.Seeker.ID = r.SeekerID.String()
		item.Seeker.Username = r.SeekerUsername
		item.Seeker.Email = r.SeekerEmail
		item.Seeker.FullName = r.SeekerFullName
		item.Seeker.Headline = r.SeekerHeadline
		item.Seeker.Location = r.SeekerLocation
		out = append(out, item)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *EmployerApplicationsHandler) Detail(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company first"})
	}

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid application id"})
	}

	row, err := h.loadOwnedApplication(c, appID, company.ID)
	if err != nil {
		return nil
	}

	resp := EmployerApplicationDetail{
		ID:          row.ID.String(),
		JobID:       row.JobID.String(),
		Status:      row.Status,
		CoverLetter: row.CoverLetter,
		CreatedAt:   row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		JobTitle:    row.JobTitle,
	}
	resp.Seeker.ID = row.SeekerID.String()
	resp.Seeker.Username = row.SeekerUsername
	resp.Seeker.Email = row.SeekerEmail
	resp.Seeker.FullName = row.SeekerFullName
	resp.Seeker.Phone = row.SeekerPhone
	resp.Seeker.Location = row.SeekerLocation
	resp.Seeker.Headline = row.SeekerHeadline
	resp.Seeker.Bio = row.SeekerBio
	resp.Seeker.ResumeURL = row.SeekerResumeUrl
	resp.Seeker.ResumeFilename = row.SeekerResumeFilename

	return c.JSON(http.StatusOK, resp)
}

var validEmployerTransitions = map[string]bool{
	"viewed":      true,
	"shortlisted": true,
	"rejected":    true,
	"hired":       true,
}

func (h *EmployerApplicationsHandler) UpdateStatus(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company first"})
	}

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid application id"})
	}

	if _, err := h.loadOwnedApplication(c, appID, company.ID); err != nil {
		return nil
	}

	req := new(UpdateApplicationStatusRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	req.Status = strings.TrimSpace(req.Status)
	if !validEmployerTransitions[req.Status] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "status must be one of: viewed, shortlisted, rejected, hired",
		})
	}

	updated, err := h.Queries.UpdateApplicationStatus(c.Request().Context(), dbgen.UpdateApplicationStatusParams{
		ID:     appID,
		Status: req.Status,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update"})
	}

	// Notify the seeker whose application was updated.
	detail, dErr := h.Queries.GetApplicationDetailForEmployer(c.Request().Context(), appID)
	if dErr == nil {
		h.Notifications.NotifyApplicationStatusChanged(
			c.Request().Context(),
			detail.SeekerID,
			detail.JobTitle,
			req.Status,
			appID.String(),
		)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id":     updated.ID.String(),
		"status": updated.Status,
	})
}
