package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"job-portal/internal/db/dbgen"
	"job-portal/internal/service"
)

type NotificationResponse struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Body      *string `json:"body"`
	Link      *string `json:"link"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

type NotificationsHandler struct {
	Queries *dbgen.Queries
}

func NewNotificationsHandler(q *dbgen.Queries) *NotificationsHandler {
	return &NotificationsHandler{Queries: q}
}

func notificationToResponse(n dbgen.Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID.String(),
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		Link:      n.Link,
		CreatedAt: n.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	if n.ReadAt.Valid {
		s := n.ReadAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.ReadAt = &s
	}
	return resp
}

// List handles GET /api/notifications?limit=20&offset=0
func (h *NotificationsHandler) List(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	limit := 20
	offset := 0
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	rows, err := h.Queries.ListNotifications(c.Request().Context(), dbgen.ListNotificationsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list notifications"})
	}

	out := make([]NotificationResponse, 0, len(rows))
	for _, n := range rows {
		out = append(out, notificationToResponse(n))
	}

	unread, err := h.Queries.CountUnreadNotifications(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": out,
		"unread":        unread,
	})
}

// Recent handles GET /api/notifications/recent — last 5 for the dropdown.
func (h *NotificationsHandler) Recent(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	rows, err := h.Queries.GetRecentNotifications(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list notifications"})
	}

	out := make([]NotificationResponse, 0, len(rows))
	for _, n := range rows {
		out = append(out, notificationToResponse(n))
	}

	unread, err := h.Queries.CountUnreadNotifications(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": out,
		"unread":        unread,
	})
}

// Count handles GET /api/notifications/unread-count — cheap endpoint for polling.
func (h *NotificationsHandler) Count(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	unread, err := h.Queries.CountUnreadNotifications(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count"})
	}

	return c.JSON(http.StatusOK, map[string]int64{"unread": unread})
}

// MarkRead handles POST /api/notifications/:id/read
func (h *NotificationsHandler) MarkRead(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid notification id"})
	}

	n, err := h.Queries.MarkNotificationRead(c.Request().Context(), dbgen.MarkNotificationReadParams{
		ID:     notifID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Already read or doesn't belong to user — treat as success.
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to mark read"})
	}

	return c.JSON(http.StatusOK, notificationToResponse(n))
}

// MarkAllRead handles POST /api/notifications/read-all
func (h *NotificationsHandler) MarkAllRead(c echo.Context) error {
	claims, ok := c.Get("user").(*service.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no user in context"})
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id"})
	}

	if err := h.Queries.MarkAllNotificationsRead(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to mark all read"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
