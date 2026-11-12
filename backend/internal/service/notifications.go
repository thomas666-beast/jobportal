package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"job-portal/internal/db/dbgen"
)

// NotificationService creates notifications for users. Failures are logged
// but never returned to the caller — a notification failing must not break
// the primary action (e.g. an application being submitted).
type NotificationService struct {
	queries *dbgen.Queries
}

func NewNotificationService(q *dbgen.Queries) *NotificationService {
	return &NotificationService{queries: q}
}

type NotifyParams struct {
	UserID uuid.UUID
	Type   string
	Title  string
	Body   string
	Link   string
}

func (s *NotificationService) Notify(ctx context.Context, p NotifyParams) {
	var body, link *string
	if p.Body != "" {
		body = &p.Body
	}
	if p.Link != "" {
		link = &p.Link
	}

	_, err := s.queries.CreateNotification(ctx, dbgen.CreateNotificationParams{
		UserID: p.UserID,
		Type:   p.Type,
		Title:  p.Title,
		Body:   body,
		Link:   link,
	})
	if err != nil {
		slog.Warn("notification failed", "user_id", p.UserID.String(), "type", p.Type, "error", err)
	}
}

// NotifyApplicationSubmitted notifies the employer that a new application arrived.
func (s *NotificationService) NotifyApplicationSubmitted(
	ctx context.Context,
	employerID uuid.UUID,
	seekerName, jobTitle, applicationID string,
) {
	s.Notify(ctx, NotifyParams{
		UserID: employerID,
		Type:   "application_submitted",
		Title:  "New application for " + jobTitle,
		Body:   seekerName + " applied to your job.",
		Link:   "/employer/applications/" + applicationID,
	})
}

// NotifyApplicationStatusChanged notifies the seeker that their status changed.
func (s *NotificationService) NotifyApplicationStatusChanged(
	ctx context.Context,
	seekerID uuid.UUID,
	jobTitle, newStatus, applicationID string,
) {
	s.Notify(ctx, NotifyParams{
		UserID: seekerID,
		Type:   "application_status_changed",
		Title:  "Your application was " + humanStatus(newStatus),
		Body:   "Update on your application to " + jobTitle + ".",
		Link:   "/applications",
	})
}

// NotifyJobApproved notifies the employer that their job went live.
func (s *NotificationService) NotifyJobApproved(
	ctx context.Context,
	employerID uuid.UUID,
	jobTitle, jobID string,
) {
	s.Notify(ctx, NotifyParams{
		UserID: employerID,
		Type:   "job_approved",
		Title:  jobTitle + " is now published",
		Body:   "Your job is live and accepting applications.",
		Link:   "/jobs/" + jobID + "/edit",
	})
}

// NotifyJobRejected notifies the employer that their job was rejected.
func (s *NotificationService) NotifyJobRejected(
	ctx context.Context,
	employerID uuid.UUID,
	jobTitle, reason, jobID string,
) {
	s.Notify(ctx, NotifyParams{
		UserID: employerID,
		Type:   "job_rejected",
		Title:  jobTitle + " needs changes",
		Body:   reason,
		Link:   "/jobs/" + jobID + "/edit",
	})
}

func humanStatus(s string) string {
	switch s {
	case "viewed":
		return "viewed"
	case "shortlisted":
		return "shortlisted"
	case "rejected":
		return "rejected"
	case "hired":
		return "marked as hired"
	case "withdrawn":
		return "withdrawn"
	}
	return s
}
