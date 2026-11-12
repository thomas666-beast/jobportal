package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *seeder) seedUsers(ctx context.Context, total int) error {
	seekerCount := total * 85 / 100
	employerCount := total - seekerCount

	// Pre-generate IDs so we can reference them later without RETURNING.
	s.seekerIDs = make([]uuid.UUID, seekerCount)
	s.employerIDs = make([]uuid.UUID, employerCount)

	// --- Users ---
	userCols := []string{"id", "username", "email", "password_hash", "role", "onboarding_completed", "created_at"}
	userRows := make([][]any, 0, total)

	// Weights for role distribution.
	for i := 0; i < seekerCount; i++ {
		id := uuid.New()
		s.seekerIDs[i] = id

		createdAt := s.randomTimestamp(365)

		userRows = append(userRows, []any{
			id,
			randomUsername(s.r, i+1),
			fmt.Sprintf("seeker%d%s", i+1, seedDomain),
			s.passHash,
			"job_seeker",
			true,
			createdAt,
		})
	}

	for i := 0; i < employerCount; i++ {
		id := uuid.New()
		s.employerIDs[i] = id

		createdAt := s.randomTimestamp(365)

		userRows = append(userRows, []any{
			id,
			randomUsername(s.r, seekerCount+i+1),
			fmt.Sprintf("employer%d%s", i+1, seedDomain),
			s.passHash,
			"employer",
			true,
			createdAt,
		})
	}

	if err := s.batchCopy(ctx, "users", userCols, userRows); err != nil {
		return fmt.Errorf("copy users: %w", err)
	}

	// --- Seeker profiles (every seeker has a completed profile) ---
	profileCols := []string{
		"user_id", "full_name", "phone", "location", "headline", "bio",
		"resume_url", "resume_filename", "resume_uploaded_at",
		"created_at", "updated_at",
	}
	profileRows := make([][]any, 0, seekerCount)

	for _, id := range s.seekerIDs {
		fullName := randomFullName(s.r)
		location := pick(s.r, cities)
		headline := pick(s.r, seekerHeadlines)
		bio := pick(s.r, seekerBios)
		createdAt := s.randomTimestamp(365)

		profileRows = append(profileRows, []any{
			id,
			fullName,
			"+358 40 " + fmt.Sprintf("%07d", s.r.Intn(10000000)),
			location,
			headline,
			bio,
			s.resumeURL,
			"resume.pdf",
			createdAt,
			createdAt,
			createdAt,
		})
	}

	if err := s.batchCopy(ctx, "seeker_profiles", profileCols, profileRows); err != nil {
		return fmt.Errorf("copy seeker_profiles: %w", err)
	}

	return nil
}

// batchCopy runs CopyFrom in chunks to bound memory and give progress feedback.
func (s *seeder) batchCopy(ctx context.Context, table string, cols []string, rows [][]any) error {
	total := len(rows)
	for start := 0; start < total; start += batchSize {
		end := start + batchSize
		if end > total {
			end = total
		}
		chunk := rows[start:end]
		_, err := s.copyFrom(ctx, table, cols, chunk)
		if err != nil {
			return fmt.Errorf("%s batch [%d:%d]: %w", table, start, end, err)
		}
		if end%10000 == 0 || end == total {
			fmt.Printf("    %s: %d/%d\n", table, end, total)
		}
	}
	return nil
}

var _ = pgx.Identifier{}
