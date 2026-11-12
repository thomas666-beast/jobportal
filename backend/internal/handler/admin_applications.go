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

type AdminApplicationListItem struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	JobID       string `json:"job_id"`
	SeekerID    string `json:"seeker_id"`
	JobTitle    string `json:"job_title"`
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
	Seeker      struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		FullName *string `json:"full_name"`
		Headline *string `json:"headline"`
	} `json:"seeker"`
}

type AdminApplicationDetail struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	CoverLetter *string `json:"cover_letter"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	JobID       string  `json:"job_id"`
	JobTitle    string  `json:"job_title"`
	JobStatus   string  `json:"job_status"`
	CompanyID   string  `json:"company_id"`
	CompanyName string  `json:"company_name"`
	SeekerID    string  `json:"seeker_id"`
	Seeker      struct {
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

type AdminApplicationsHandler struct {
	Queries *dbgen.Queries
}

func NewAdminApplicationsHandler(q *dbgen.Queries) *AdminApplicationsHandler {
	return &AdminApplicationsHandler{Queries: q}
}

var validApplicationStatuses = map[string]bool{
	"submitted":   true,
	"viewed":      true,
	"shortlisted": true,
	"rejected":    true,
	"hired":       true,
	"withdrawn":   true,
}

func (h *AdminApplicationsHandler) List(c echo.Context) error {
	status := strings.TrimSpace(c.QueryParam("status"))
	if status != "" && !validApplicationStatuses[status] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid status filter"})
	}

	seekerQ := strings.TrimSpace(c.QueryParam("seeker_q"))

	var jobIDPtr *uuid.UUID
	if v := strings.TrimSpace(c.QueryParam("job_id")); v != "" {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job_id filter"})
		}
		jobIDPtr = &parsed
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

	var statusPtr, seekerQPtr *string
	if status != "" {
		statusPtr = &status
	}
	if seekerQ != "" {
		seekerQPtr = &seekerQ
	}

	ctx := c.Request().Context()

	rows, err := h.Queries.ListAdminApplications(ctx, dbgen.ListAdminApplicationsParams{
		Status:  statusPtr,
		JobID:   jobIDPtr,
		SeekerQ: seekerQPtr,
		Lim:     int32(limit),
		Off:     int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list applications"})
	}

	total, err := h.Queries.CountAdminApplications(ctx, dbgen.CountAdminApplicationsParams{
		Status:  statusPtr,
		JobID:   jobIDPtr,
		SeekerQ: seekerQPtr,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count applications"})
	}

	out := make([]AdminApplicationListItem, 0, len(rows))
	for _, r := range rows {
		item := AdminApplicationListItem{
			ID:          r.ID.String(),
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   r.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			JobID:       r.JobID.String(),
			SeekerID:    r.SeekerID.String(),
			JobTitle:    r.JobTitle,
			CompanyID:   r.CompanyID.String(),
			CompanyName: r.CompanyName,
		}
		item.Seeker.Username = r.SeekerUsername
		item.Seeker.Email = r.SeekerEmail
		item.Seeker.FullName = r.SeekerFullName
		item.Seeker.Headline = r.SeekerHeadline
		out = append(out, item)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"applications": out,
		"total":        total,
	})
}

func (h *AdminApplicationsHandler) Detail(c echo.Context) error {
	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid application id"})
	}

	r, err := h.Queries.GetAdminApplicationByID(c.Request().Context(), appID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch application"})
	}

	resp := AdminApplicationDetail{
		ID:          r.ID.String(),
		Status:      r.Status,
		CoverLetter: r.CoverLetter,
		CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   r.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		JobID:       r.JobID.String(),
		JobTitle:    r.JobTitle,
		JobStatus:   r.JobStatus,
		CompanyID:   r.CompanyID.String(),
		CompanyName: r.CompanyName,
		SeekerID:    r.SeekerID.String(),
	}
	resp.Seeker.Username = r.SeekerUsername
	resp.Seeker.Email = r.SeekerEmail
	resp.Seeker.FullName = r.SeekerFullName
	resp.Seeker.Phone = r.SeekerPhone
	resp.Seeker.Location = r.SeekerLocation
	resp.Seeker.Headline = r.SeekerHeadline
	resp.Seeker.Bio = r.SeekerBio
	resp.Seeker.ResumeURL = r.SeekerResumeUrl
	resp.Seeker.ResumeFilename = r.SeekerResumeFilename

	return c.JSON(http.StatusOK, resp)
}
