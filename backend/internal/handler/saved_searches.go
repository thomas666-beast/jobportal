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

type SavedSearchResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Query      *string `json:"query"`
	Location   *string `json:"location"`
	RemoteType *string `json:"remote_type"`
	CreatedAt  string  `json:"created_at"`
	MatchCount int64   `json:"match_count"`
}

type CreateSavedSearchRequest struct {
	Name       string `json:"name"`
	Query      string `json:"query"`
	Location   string `json:"location"`
	RemoteType string `json:"remote_type"`
}

type SavedSearchesHandler struct {
	Queries *dbgen.Queries
}

func NewSavedSearchesHandler(q *dbgen.Queries) *SavedSearchesHandler {
	return &SavedSearchesHandler{Queries: q}
}

var validSavedSearchRemoteTypes = map[string]bool{"": true, "remote": true, "hybrid": true, "on_site": true}

// toPtr converts "" -> nil, otherwise &trimmed.
func toPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func (h *SavedSearchesHandler) Create(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	req := new(CreateSavedSearchRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}
	if len(req.Name) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name must be 100 characters or fewer"})
	}

	if !validSavedSearchRemoteTypes[req.RemoteType] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid remote_type"})
	}

	query := strings.TrimSpace(req.Query)
	location := strings.TrimSpace(req.Location)

	if query == "" && location == "" && req.RemoteType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "at least one filter (query, location, or remote_type) is required",
		})
	}

	created, err := h.Queries.CreateSavedSearch(c.Request().Context(), dbgen.CreateSavedSearchParams{
		SeekerID:   seekerID,
		Name:       req.Name,
		Query:      toPtr(query),
		Location:   toPtr(location),
		RemoteType: toPtr(req.RemoteType),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "a saved search with this name already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save search"})
	}

	matchCount := h.countMatches(c, created.Query, created.Location, created.RemoteType)

	return c.JSON(http.StatusCreated, SavedSearchResponse{
		ID:         created.ID.String(),
		Name:       created.Name,
		Query:      created.Query,
		Location:   created.Location,
		RemoteType: created.RemoteType,
		CreatedAt:  created.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		MatchCount: matchCount,
	})
}

func (h *SavedSearchesHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	rows, err := h.Queries.ListSavedSearchesBySeeker(c.Request().Context(), seekerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list saved searches"})
	}

	out := make([]SavedSearchResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, SavedSearchResponse{
			ID:         r.ID.String(),
			Name:       r.Name,
			Query:      r.Query,
			Location:   r.Location,
			RemoteType: r.RemoteType,
			CreatedAt:  r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			MatchCount: h.countMatches(c, r.Query, r.Location, r.RemoteType),
		})
	}

	return c.JSON(http.StatusOK, out)
}

func (h *SavedSearchesHandler) Delete(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	seekerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	searchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid search id"})
	}

	// Ownership check.
	existing, err := h.Queries.GetSavedSearchByID(c.Request().Context(), searchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "saved search not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch saved search"})
	}
	if existing.SeekerID != seekerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "not your saved search"})
	}

	if err := h.Queries.DeleteSavedSearch(c.Request().Context(), dbgen.DeleteSavedSearchParams{
		ID:       searchID,
		SeekerID: seekerID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete"})
	}

	return c.NoContent(http.StatusNoContent)
}

// countMatches returns how many published jobs match a saved search's filters.
// Reuses CountPublishedJobs from jobs.sql.
func (h *SavedSearchesHandler) countMatches(c echo.Context, q, location, remoteType *string) int64 {
	count, err := h.Queries.CountPublishedJobs(c.Request().Context(), dbgen.CountPublishedJobsParams{
		Q:          q,
		Location:   location,
		RemoteType: remoteType,
	})
	if err != nil {
		return 0
	}
	return count
}
