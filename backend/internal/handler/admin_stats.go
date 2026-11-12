package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
)

type AdminStatsHandler struct {
	Queries *dbgen.Queries
}

func NewAdminStatsHandler(q *dbgen.Queries) *AdminStatsHandler {
	return &AdminStatsHandler{Queries: q}
}

type recentUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type recentApplication struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	SeekerUsername string `json:"seeker_username"`
	JobTitle       string `json:"job_title"`
	CompanyName    string `json:"company_name"`
}

// Stats handles GET /api/admin/stats
func (h *AdminStatsHandler) Stats(c echo.Context) error {
	ctx := c.Request().Context()

	totalUsers, err := h.Queries.CountAllUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	totalCompanies, err := h.Queries.CountAllCompanies(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	totalJobs, err := h.Queries.CountAllJobs(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	totalApplications, err := h.Queries.CountAllApplications(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	usersByRole, err := h.Queries.CountUsersByRole(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	jobsByStatus, err := h.Queries.CountJobsByStatusAdmin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	appsByStatus, err := h.Queries.CountApplicationsByStatus(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	recentUsers, err := h.Queries.ListRecentUsers(ctx, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	recentApps, err := h.Queries.ListRecentApplications(ctx, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
	}

	roleMap := make(map[string]int64)
	for _, r := range usersByRole {
		roleMap[r.Role] = r.Count
	}
	jobMap := make(map[string]int64)
	for _, r := range jobsByStatus {
		jobMap[r.Status] = r.Count
	}
	appMap := make(map[string]int64)
	for _, r := range appsByStatus {
		appMap[r.Status] = r.Count
	}

	ru := make([]recentUser, 0, len(recentUsers))
	for _, r := range recentUsers {
		ru = append(ru, recentUser{
			ID:        r.ID.String(),
			Username:  r.Username,
			Email:     r.Email,
			Role:      r.Role,
			CreatedAt: r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	ra := make([]recentApplication, 0, len(recentApps))
	for _, r := range recentApps {
		ra = append(ra, recentApplication{
			ID:             r.ID.String(),
			Status:         r.Status,
			CreatedAt:      r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			SeekerUsername: r.SeekerUsername,
			JobTitle:       r.JobTitle,
			CompanyName:    r.CompanyName,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totals": map[string]int64{
			"users":        totalUsers,
			"companies":    totalCompanies,
			"jobs":         totalJobs,
			"applications": totalApplications,
		},
		"users_by_role":          roleMap,
		"jobs_by_status":         jobMap,
		"applications_by_status": appMap,
		"recent_users":           ru,
		"recent_applications":    ra,
	})
}
