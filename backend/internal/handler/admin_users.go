package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
)

type AdminUserListItem struct {
	ID                  string  `json:"id"`
	Username            string  `json:"username"`
	Email               string  `json:"email"`
	Role                string  `json:"role"`
	OnboardingCompleted bool    `json:"onboarding_completed"`
	CreatedAt           string  `json:"created_at"`
	SeekerFullName      *string `json:"seeker_full_name,omitempty"`
	SeekerHeadline      *string `json:"seeker_headline,omitempty"`
	CompanyName         *string `json:"company_name,omitempty"`
}

type AdminUserDetail struct {
	ID                  string `json:"id"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	Role                string `json:"role"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`

	Seeker *struct {
		FullName  *string `json:"full_name"`
		Phone     *string `json:"phone"`
		Location  *string `json:"location"`
		Headline  *string `json:"headline"`
		Bio       *string `json:"bio"`
		CreatedAt *string `json:"created_at"`
	} `json:"seeker,omitempty"`

	Company *struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
		Location    *string `json:"location"`
		LogoURL     *string `json:"logo_url"`
		CreatedAt   string  `json:"created_at"`
	} `json:"company,omitempty"`
}

type AdminUsersHandler struct {
	Queries *dbgen.Queries
}

func NewAdminUsersHandler(q *dbgen.Queries) *AdminUsersHandler {
	return &AdminUsersHandler{Queries: q}
}

func (h *AdminUsersHandler) List(c echo.Context) error {
	role := c.QueryParam("role")
	validRoles := map[string]bool{"": true, "job_seeker": true, "employer": true, "admin": true}
	if !validRoles[role] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role filter"})
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

	ctx := c.Request().Context()

	if role == "" {
		rows, err := h.Queries.ListUsers(ctx, dbgen.ListUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
		}
		total, err := h.Queries.CountUsersFiltered(ctx, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count users"})
		}
		out := make([]AdminUserListItem, 0, len(rows))
		for _, r := range rows {
			out = append(out, toAdminUserListItem(
				r.ID, r.Username, r.Email, r.Role, r.OnboardingCompleted,
				r.CreatedAt.Time, r.SeekerFullName, r.SeekerHeadline, r.CompanyName,
			))
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"users": out, "total": total})
	}

	rows, err := h.Queries.ListUsersByRole(ctx, dbgen.ListUsersByRoleParams{
		Role:   role,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
	}
	rolePtr := role
	total, err := h.Queries.CountUsersFiltered(ctx, &rolePtr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count users"})
	}
	out := make([]AdminUserListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, toAdminUserListItem(
			r.ID, r.Username, r.Email, r.Role, r.OnboardingCompleted,
			r.CreatedAt.Time, r.SeekerFullName, r.SeekerHeadline, r.CompanyName,
		))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"users": out, "total": total})
}

func (h *AdminUsersHandler) Detail(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	r, err := h.Queries.GetUserDetailByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch user"})
	}

	resp := AdminUserDetail{
		ID:                  r.ID.String(),
		Username:            r.Username,
		Email:               r.Email,
		Role:                r.Role,
		OnboardingCompleted: r.OnboardingCompleted,
		CreatedAt:           r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:           r.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if r.SeekerFullName != nil || r.SeekerHeadline != nil || r.SeekerBio != nil {
		resp.Seeker = &struct {
			FullName  *string `json:"full_name"`
			Phone     *string `json:"phone"`
			Location  *string `json:"location"`
			Headline  *string `json:"headline"`
			Bio       *string `json:"bio"`
			CreatedAt *string `json:"created_at"`
		}{
			FullName: r.SeekerFullName,
			Phone:    r.SeekerPhone,
			Location: r.SeekerLocation,
			Headline: r.SeekerHeadline,
			Bio:      r.SeekerBio,
		}
		if r.SeekerCreatedAt.Valid {
			s := r.SeekerCreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			resp.Seeker.CreatedAt = &s
		}
	}

	if r.CompanyID != nil {
		c := &struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Description *string `json:"description"`
			Website     *string `json:"website"`
			Location    *string `json:"location"`
			LogoURL     *string `json:"logo_url"`
			CreatedAt   string  `json:"created_at"`
		}{
			ID:          r.CompanyID.String(),
			Name:        derefStr(r.CompanyName),
			Description: r.CompanyDescription,
			Website:     r.CompanyWebsite,
			Location:    r.CompanyLocation,
			LogoURL:     r.CompanyLogoUrl,
		}
		if r.CompanyCreatedAt.Valid {
			c.CreatedAt = r.CompanyCreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		resp.Company = c
	}

	return c.JSON(http.StatusOK, resp)
}

func toAdminUserListItem(
	id uuid.UUID, username, email, role string, onboarding bool,
	createdAt interface{ Format(string) string },
	seekerFullName, seekerHeadline, companyName *string,
) AdminUserListItem {
	return AdminUserListItem{
		ID:                  id.String(),
		Username:            username,
		Email:               email,
		Role:                role,
		OnboardingCompleted: onboarding,
		CreatedAt:           createdAt.Format("2006-01-02T15:04:05Z07:00"),
		SeekerFullName:      seekerFullName,
		SeekerHeadline:      seekerHeadline,
		CompanyName:         companyName,
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
