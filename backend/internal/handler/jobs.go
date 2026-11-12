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

type JobResponse struct {
	ID              string   `json:"id"`
	CompanyID       string   `json:"company_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Requirements    *string  `json:"requirements"`
	Location        *string  `json:"location"`
	RemoteType      string   `json:"remote_type"`
	SalaryMin       *int32   `json:"salary_min"`
	SalaryMax       *int32   `json:"salary_max"`
	SalaryCurrency  *string  `json:"salary_currency"`
	EmploymentType  string   `json:"employment_type"`
	Category        *string  `json:"category"`
	Tags            []string `json:"tags"`
	Status          string   `json:"status"`
	RejectionReason *string  `json:"rejection_reason"`
	PublishedAt     *string  `json:"published_at"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type JobRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Requirements   string   `json:"requirements"`
	Location       string   `json:"location"`
	RemoteType     string   `json:"remote_type"`
	SalaryMin      *int32   `json:"salary_min"`
	SalaryMax      *int32   `json:"salary_max"`
	SalaryCurrency string   `json:"salary_currency"`
	EmploymentType string   `json:"employment_type"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
}

type JobsHandler struct {
	Queries *dbgen.Queries
}

func NewJobsHandler(q *dbgen.Queries) *JobsHandler {
	return &JobsHandler{Queries: q}
}

var validRemoteTypes = map[string]bool{"remote": true, "hybrid": true, "on_site": true}
var validEmploymentTypes = map[string]bool{"full_time": true, "part_time": true, "contract": true, "internship": true}

// validateJobRequest returns an error message or "" if valid.
func validateJobRequest(req *JobRequest) string {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return "title is required"
	}
	if len(req.Title) > 200 {
		return "title must be 200 characters or fewer"
	}

	req.Description = strings.TrimSpace(req.Description)
	if req.Description == "" {
		return "description is required"
	}
	if len(req.Description) > 10000 {
		return "description must be 10000 characters or fewer"
	}

	if len(req.Requirements) > 10000 {
		return "requirements must be 10000 characters or fewer"
	}

	if len(req.Location) > 120 {
		return "location must be 120 characters or fewer"
	}

	if !validRemoteTypes[req.RemoteType] {
		return "remote_type must be one of: remote, hybrid, on_site"
	}

	if !validEmploymentTypes[req.EmploymentType] {
		return "employment_type must be one of: full_time, part_time, contract, internship"
	}

	if req.SalaryMin != nil && req.SalaryMax != nil && *req.SalaryMin > *req.SalaryMax {
		return "salary_min cannot be greater than salary_max"
	}

	if req.SalaryMin != nil && *req.SalaryMin < 0 {
		return "salary_min cannot be negative"
	}
	if req.SalaryMax != nil && *req.SalaryMax < 0 {
		return "salary_max cannot be negative"
	}

	if len(req.SalaryCurrency) > 3 {
		return "salary_currency must be 3 characters or fewer"
	}

	if len(req.Category) > 100 {
		return "category must be 100 characters or fewer"
	}

	if len(req.Tags) > 10 {
		return "at most 10 tags allowed"
	}
	for _, t := range req.Tags {
		if len(strings.TrimSpace(t)) == 0 {
			return "tags cannot contain empty values"
		}
		if len(t) > 30 {
			return "each tag must be 30 characters or fewer"
		}
	}

	return ""
}

func jobToResponse(j dbgen.Job) JobResponse {
	resp := JobResponse{
		ID:              j.ID.String(),
		CompanyID:       j.CompanyID.String(),
		Title:           j.Title,
		Description:     j.Description,
		Requirements:    j.Requirements,
		Location:        j.Location,
		RemoteType:      j.RemoteType,
		SalaryMin:       j.SalaryMin,
		SalaryMax:       j.SalaryMax,
		SalaryCurrency:  j.SalaryCurrency,
		EmploymentType:  j.EmploymentType,
		Category:        j.Category,
		Tags:            j.Tags,
		Status:          j.Status,
		RejectionReason: j.RejectionReason,
		CreatedAt:       j.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       j.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	if j.Tags == nil {
		resp.Tags = []string{}
	}
	if j.PublishedAt.Valid {
		s := j.PublishedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.PublishedAt = &s
	}
	return resp
}

// getOwnedJob loads a job and verifies the caller's company owns it.
// Returns nil and writes the error response if not found or not owned.
func (h *JobsHandler) getOwnedJob(c echo.Context, jobID uuid.UUID, companyID uuid.UUID) (*dbgen.Job, error) {
	job, err := h.Queries.GetJobByID(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
			return nil, err
		}
		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch job"})
		return nil, err
	}
	if job.CompanyID != companyID {
		_ = c.JSON(http.StatusForbidden, map[string]string{"error": "not your job"})
		return nil, errors.New("not owned")
	}
	return &job, nil
}

func (h *JobsHandler) Create(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "create your company profile first",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load company"})
	}

	req := new(JobRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if msg := validateJobRequest(req); msg != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	job, err := h.Queries.CreateJob(c.Request().Context(), dbgen.CreateJobParams{
		CompanyID:       company.ID,
		CreatedByUserID: &userID,
		Title:           req.Title,
		Description:     req.Description,
		Requirements:    toNullable(req.Requirements),
		Location:        toNullable(req.Location),
		RemoteType:      req.RemoteType,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		SalaryCurrency:  toNullable(req.SalaryCurrency),
		EmploymentType:  req.EmploymentType,
		Category:        toNullable(req.Category),
		Tags:            req.Tags,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create job"})
	}

	return c.JSON(http.StatusCreated, jobToResponse(job))
}

func (h *JobsHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusOK, []JobResponse{})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load company"})
	}

	jobs, err := h.Queries.ListJobsByCompany(c.Request().Context(), company.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list jobs"})
	}

	out := make([]JobResponse, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, jobToResponse(j))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *JobsHandler) Get(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company profile first"})
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	job, err := h.getOwnedJob(c, jobID, company.ID)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, jobToResponse(*job))
}

func (h *JobsHandler) Update(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company profile first"})
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	job, err := h.getOwnedJob(c, jobID, company.ID)
	if err != nil {
		return nil
	}

	if job.Status != "draft" && job.Status != "pending_review" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "job can only be edited while in draft or pending_review",
		})
	}

	req := new(JobRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if msg := validateJobRequest(req); msg != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	updated, err := h.Queries.UpdateJob(c.Request().Context(), dbgen.UpdateJobParams{
		ID:             jobID,
		Title:          req.Title,
		Description:    req.Description,
		Requirements:   toNullable(req.Requirements),
		Location:       toNullable(req.Location),
		RemoteType:     req.RemoteType,
		SalaryMin:      req.SalaryMin,
		SalaryMax:      req.SalaryMax,
		SalaryCurrency: toNullable(req.SalaryCurrency),
		EmploymentType: req.EmploymentType,
		Category:       toNullable(req.Category),
		Tags:           req.Tags,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update job"})
	}

	return c.JSON(http.StatusOK, jobToResponse(updated))
}

func (h *JobsHandler) Submit(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	verified, vErr := h.Queries.GetUserEmailVerifiedByID(c.Request().Context(), userID)
	if vErr != nil || !verified.Valid {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "please verify your email before submitting jobs",
		})
	}

	return h.transition(c, "draft", "pending_review")
}

func (h *JobsHandler) Close(c echo.Context) error {
	return h.transition(c, "published", "closed")
}

func (h *JobsHandler) Archive(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company profile first"})
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	job, err := h.getOwnedJob(c, jobID, company.ID)
	if err != nil {
		return nil
	}

	if job.Status == "archived" {
		return c.JSON(http.StatusConflict, map[string]string{"error": "job is already archived"})
	}

	updated, err := h.Queries.UpdateJobStatus(c.Request().Context(), dbgen.UpdateJobStatusParams{
		ID:              jobID,
		Status:          "archived",
		RejectionReason: nil,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to archive job"})
	}
	return c.JSON(http.StatusOK, jobToResponse(updated))
}

// transition handles simple state transitions requiring ownership.
func (h *JobsHandler) transition(c echo.Context, fromStatus, toStatus string) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, _ := uuid.Parse(claims.UserID)

	company, err := h.Queries.GetCompanyByOwner(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "create your company profile first"})
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job id"})
	}

	job, err := h.getOwnedJob(c, jobID, company.ID)
	if err != nil {
		return nil
	}

	if job.Status != fromStatus {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "invalid state transition from " + job.Status + " to " + toStatus,
		})
	}

	updated, err := h.Queries.UpdateJobStatus(c.Request().Context(), dbgen.UpdateJobStatusParams{
		ID:              jobID,
		Status:          toStatus,
		RejectionReason: nil,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update status"})
	}
	return c.JSON(http.StatusOK, jobToResponse(updated))
}
