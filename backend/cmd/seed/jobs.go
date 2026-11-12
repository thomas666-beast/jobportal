package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *seeder) seedJobs(ctx context.Context, total int) error {
	if len(s.companyIDs) == 0 {
		return fmt.Errorf("no companies to attach jobs to")
	}

	cols := []string{
		"id", "company_id", "created_by_user_id",
		"title", "description", "requirements",
		"location", "remote_type",
		"salary_min", "salary_max", "salary_currency",
		"employment_type", "category", "tags",
		"status", "rejection_reason", "published_at",
		"created_at", "updated_at",
	}

	rows := make([][]any, 0, total)
	s.jobIDs = make([]uuid.UUID, 0, total)

	// Status distribution: 60% published, 16% draft, 12% pending, 8% closed, 4% archived
	type dist struct {
		status string
		weight int
	}
	statuses := []dist{
		{"published", 60},
		{"draft", 16},
		{"pending_review", 12},
		{"closed", 8},
		{"archived", 4},
	}
	totalWeight := 0
	for _, d := range statuses {
		totalWeight += d.weight
	}

	for i := 0; i < total; i++ {
		id := uuid.New()
		companyID := s.companyIDs[s.r.Intn(len(s.companyIDs))]

		title := pick(s.r, jobTitles)
		description := pick(s.r, jobDescriptions)
		requirements := pick(s.r, jobRequirements)
		location := pick(s.r, cities)
		remoteType := pick(s.r, remoteTypes)
		employmentType := pick(s.r, employmentTypes)
		category := pick(s.r, jobCategories)
		tags := pickTags(s.r, 2, 6)
		currency := pick(s.r, currencies)

		var salaryMin, salaryMax any
		if s.r.Float64() < 0.75 {
			base := 30000 + s.r.Intn(90000)
			salaryMin = base
			salaryMax = base + 5000 + s.r.Intn(40000)
		}

		// Pick status by weight.
		r := s.r.Intn(totalWeight)
		acc := 0
		status := "published"
		for _, d := range statuses {
			acc += d.weight
			if r < acc {
				status = d.status
				break
			}
		}

		createdAt := s.randomTimestamp(365)

		var publishedAt, rejectionReason any
		if status == "published" || status == "closed" || status == "archived" {
			publishedAt = createdAt.Add(randomDuration(s.r, 0, 24*3600))
			s.jobIDs = append(s.jobIDs, id)
		}
		if status == "draft" && s.r.Float64() < 0.15 {
			rejectionReason = "Please provide more detail in the description."
		}

		rows = append(rows, []any{
			id, companyID, nil,
			title, description, requirements,
			location, remoteType,
			salaryMin, salaryMax, currency,
			employmentType, category, tags,
			status, rejectionReason, publishedAt,
			createdAt, createdAt,
		})
	}

	if err := s.batchCopy(ctx, "jobs", cols, rows); err != nil {
		return err
	}

	if len(s.jobIDs) == 0 {
		return fmt.Errorf("no published jobs — applications will be empty")
	}
	return nil
}
