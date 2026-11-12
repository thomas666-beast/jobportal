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

type SeekerProfileResponse struct {
	UserID           string  `json:"user_id"`
	FullName         string  `json:"full_name"`
	Phone            *string `json:"phone"`
	Location         *string `json:"location"`
	Headline         *string `json:"headline"`
	Bio              *string `json:"bio"`
	ResumeURL        *string `json:"resume_url"`
	ResumeFilename   *string `json:"resume_filename"`
	ResumeUploadedAt *string `json:"resume_uploaded_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type SeekerProfileRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Headline string `json:"headline"`
	Bio      string `json:"bio"`
}

type SeekerProfileHandler struct {
	Queries *dbgen.Queries
}

func NewSeekerProfileHandler(q *dbgen.Queries) *SeekerProfileHandler {
	return &SeekerProfileHandler{Queries: q}
}

// toNullable converts "" -> nil, otherwise &trimmed.
func toNullable(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func (h *SeekerProfileHandler) Get(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id in token",
		})
	}

	profile, err := h.Queries.GetSeekerProfile(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "profile not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch profile",
		})
	}

	return c.JSON(http.StatusOK, profileToResponse(profile))
}

func (h *SeekerProfileHandler) Upsert(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "no user in context",
		})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "invalid user id in token",
		})
	}

	req := new(SeekerProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "full_name is required",
		})
	}
	if len(req.FullName) > 120 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "full_name must be 120 characters or fewer",
		})
	}

	phone := strings.TrimSpace(req.Phone)
	if len(phone) > 32 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "phone must be 32 characters or fewer",
		})
	}

	headline := strings.TrimSpace(req.Headline)
	if len(headline) > 120 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "headline must be 120 characters or fewer",
		})
	}

	bio := strings.TrimSpace(req.Bio)
	if len(bio) > 2000 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bio must be 2000 characters or fewer",
		})
	}

	location := strings.TrimSpace(req.Location)
	if len(location) > 120 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "location must be 120 characters or fewer",
		})
	}

	profile, err := h.Queries.UpsertSeekerProfile(c.Request().Context(), dbgen.UpsertSeekerProfileParams{
		UserID:   userID,
		FullName: req.FullName,
		Phone:    toNullable(phone),
		Location: toNullable(location),
		Headline: toNullable(headline),
		Bio:      toNullable(bio),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to save profile",
		})
	}

	return c.JSON(http.StatusOK, profileToResponse(profile))
}

func profileRowToResponse(p dbgen.SeekerProfile) SeekerProfileResponse {
	return profileToResponse(p)
}

func profileToResponse(input interface{}) SeekerProfileResponse {
	b, err := json.Marshal(input)
	if err != nil {
		return SeekerProfileResponse{}
	}
	var p dbgen.SeekerProfile
	if err := json.Unmarshal(b, &p); err != nil {
		return SeekerProfileResponse{}
	}

	resp := SeekerProfileResponse{
		UserID:         p.UserID.String(),
		FullName:       p.FullName,
		Phone:          p.Phone,
		Location:       p.Location,
		Headline:       p.Headline,
		Bio:            p.Bio,
		ResumeURL:      p.ResumeUrl,
		ResumeFilename: p.ResumeFilename,
		CreatedAt:      p.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      p.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	if p.ResumeUploadedAt.Valid {
		s := p.ResumeUploadedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.ResumeUploadedAt = &s
	}
	return resp
}

func getSeekerProfileRowToModel(r dbgen.GetSeekerProfileRow) dbgen.SeekerProfile {
	return dbgen.SeekerProfile{
		UserID:           r.UserID,
		FullName:         r.FullName,
		Phone:            r.Phone,
		Location:         r.Location,
		Headline:         r.Headline,
		Bio:              r.Bio,
		ResumeUrl:        r.ResumeUrl,
		ResumeFilename:   r.ResumeFilename,
		ResumeUploadedAt: r.ResumeUploadedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func upsertSeekerProfileRowToModel(r dbgen.UpsertSeekerProfileRow) dbgen.SeekerProfile {
	return dbgen.SeekerProfile{
		UserID:           r.UserID,
		FullName:         r.FullName,
		Phone:            r.Phone,
		Location:         r.Location,
		Headline:         r.Headline,
		Bio:              r.Bio,
		ResumeUrl:        r.ResumeUrl,
		ResumeFilename:   r.ResumeFilename,
		ResumeUploadedAt: r.ResumeUploadedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func setSeekerResumeRowToModel(r dbgen.SetSeekerResumeRow) dbgen.SeekerProfile {
	return dbgen.SeekerProfile{
		UserID:           r.UserID,
		FullName:         r.FullName,
		Phone:            r.Phone,
		Location:         r.Location,
		Headline:         r.Headline,
		Bio:              r.Bio,
		ResumeUrl:        r.ResumeUrl,
		ResumeFilename:   r.ResumeFilename,
		ResumeUploadedAt: r.ResumeUploadedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func clearSeekerResumeRowToModel(r dbgen.ClearSeekerResumeRow) dbgen.SeekerProfile {
	return dbgen.SeekerProfile{
		UserID:           r.UserID,
		FullName:         r.FullName,
		Phone:            r.Phone,
		Location:         r.Location,
		Headline:         r.Headline,
		Bio:              r.Bio,
		ResumeUrl:        r.ResumeUrl,
		ResumeFilename:   r.ResumeFilename,
		ResumeUploadedAt: r.ResumeUploadedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}
