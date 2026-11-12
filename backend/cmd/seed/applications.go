package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var applicationStatuses = []string{
	"submitted", "submitted", "submitted",
	"viewed", "viewed",
	"shortlisted",
	"rejected", "rejected",
	"hired",
	"withdrawn",
}

var coverLetterPool = coverLetters

func (s *seeder) seedApplications(ctx context.Context, total int) error {
	if len(s.jobIDs) == 0 || len(s.seekerIDs) == 0 {
		return nil
	}

	// UNIQUE(job_id, seeker_id) — deduplicate while generating.
	seen := make(map[string]struct{}, total)
	cols := []string{"id", "job_id", "seeker_id", "cover_letter", "status", "created_at", "updated_at"}
	rows := make([][]any, 0, total)

	attempts := 0
	maxAttempts := total * 3 // avoid infinite loop if the space is exhausted
	for len(rows) < total && attempts < maxAttempts {
		attempts++
		jobID := s.jobIDs[s.r.Intn(len(s.jobIDs))]
		seekerID := s.seekerIDs[s.r.Intn(len(s.seekerIDs))]

		key := jobID.String() + ":" + seekerID.String()
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		id := uuid.New()
		status := pick(s.r, applicationStatuses)
		createdAt := s.randomTimestamp(180)
		updatedAt := createdAt.Add(randomDuration(s.r, 0, 72*3600))

		var cover any
		if s.r.Float64() < 0.5 {
			cover = pick(s.r, coverLetterPool)
		}

		rows = append(rows, []any{
			id, jobID, seekerID, cover, status, createdAt, updatedAt,
		})
	}

	if len(rows) == 0 {
		return nil
	}

	return s.batchCopy(ctx, "applications", cols, rows)
}

func (s *seeder) seedSavedJobs(ctx context.Context, total int) error {
	if len(s.jobIDs) == 0 || len(s.seekerIDs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, total)
	cols := []string{"seeker_id", "job_id", "created_at"}
	rows := make([][]any, 0, total)

	attempts := 0
	maxAttempts := total * 3
	for len(rows) < total && attempts < maxAttempts {
		attempts++
		jobID := s.jobIDs[s.r.Intn(len(s.jobIDs))]
		seekerID := s.seekerIDs[s.r.Intn(len(s.seekerIDs))]

		key := jobID.String() + ":" + seekerID.String()
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		rows = append(rows, []any{
			seekerID, jobID, s.randomTimestamp(90),
		})
	}

	if len(rows) == 0 {
		return nil
	}

	return s.batchCopy(ctx, "saved_jobs", cols, rows)
}

func (s *seeder) seedSavedSearches(ctx context.Context, total int) error {
	if len(s.seekerIDs) == 0 {
		return nil
	}

	cols := []string{"id", "seeker_id", "name", "query", "location", "remote_type", "created_at"}
	rows := make([][]any, 0, total)

	// Track per-seeker names to avoid UNIQUE violations.
	used := make(map[string]struct{}, total)

	attempts := 0
	maxAttempts := total * 3
	for len(rows) < total && attempts < maxAttempts {
		attempts++
		seekerID := s.seekerIDs[s.r.Intn(len(s.seekerIDs))]
		name := pick(s.r, searchNames)

		key := seekerID.String() + "|" + name
		if _, exists := used[key]; exists {
			continue
		}
		used[key] = struct{}{}

		var query, location, remoteType any
		if s.r.Float64() < 0.6 {
			query = skillTags[s.r.Intn(len(skillTags))]
		}
		if s.r.Float64() < 0.4 {
			location = cities[s.r.Intn(len(cities))]
		}
		if s.r.Float64() < 0.5 {
			remoteType = []string{"remote", "hybrid", "on_site"}[s.r.Intn(3)]
		}

		// Ensure at least one filter is set — required by the handler.
		if query == nil && location == nil && remoteType == nil {
			query = skillTags[s.r.Intn(len(skillTags))]
		}

		rows = append(rows, []any{
			uuid.New(), seekerID, name, query, location, remoteType,
			s.randomTimestamp(120),
		})
	}

	if len(rows) == 0 {
		return nil
	}

	return s.batchCopy(ctx, "saved_searches", cols, rows)
}

// --- helpers ---

// randomTimestamp returns a time in the past `days` days, at a random hour.
func (s *seeder) randomTimestamp(days int) time.Time {
	now := time.Now().UTC()
	secs := int64(days * 24 * 3600)
	offset := time.Duration(s.r.Int63n(secs)) * time.Second
	return now.Add(-offset)
}

func randomDuration(r *rand.Rand, minSec, maxSec int) time.Duration {
	if maxSec <= minSec {
		return time.Duration(minSec) * time.Second
	}
	return time.Duration(minSec+r.Intn(maxSec-minSec)) * time.Second
}
