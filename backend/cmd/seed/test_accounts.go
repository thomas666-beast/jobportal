package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// createTestAccounts inserts admin1, seeker1, employer1 with linked data.
// These are stable across reseeds — if they already exist, they're left alone.
func (s *seeder) createTestAccounts(ctx context.Context) error {
	// admin1@seed.local
	if err := s.ensureTestUser(ctx, "admin1", "admin"); err != nil {
		return err
	}

	// seeker1@seed.local — with profile, resume, some applications, saves.
	if err := s.ensureTestUser(ctx, "seeker1", "job_seeker"); err != nil {
		return err
	}

	// employer1@seed.local — with company, published jobs, some applications.
	if err := s.ensureTestUser(ctx, "employer1", "employer"); err != nil {
		return err
	}

	return nil
}

func (s *seeder) ensureTestUser(ctx context.Context, username, role string) error {
	email := username + seedDomain

	var existing uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&existing)
	if err == nil {
		// Already exists — nothing to do.
		return nil
	}

	userID := uuid.New()
	now := time.Now().UTC()

	_, err = s.pool.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, role, onboarding_completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, $6, $6)
	`, userID, username, email, s.passHash, role, now)
	if err != nil {
		return fmt.Errorf("insert %s: %w", email, err)
	}

	switch role {
	case "job_seeker":
		if err := s.linkTestSeeker(ctx, userID, now); err != nil {
			return err
		}
	case "employer":
		if err := s.linkTestEmployer(ctx, userID, now); err != nil {
			return err
		}
	}

	return nil
}

func (s *seeder) linkTestSeeker(ctx context.Context, userID uuid.UUID, now time.Time) error {
	// Profile with resume.
	_, err := s.pool.Exec(ctx, `
		INSERT INTO seeker_profiles (
			user_id, full_name, phone, location, headline, bio,
			resume_url, resume_filename, resume_uploaded_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
	`, userID, "Alex Tester", "+358 40 1234567", "Helsinki",
		"Senior Go developer with a passion for clean systems",
		"I'm the seeker1 test account. Everything here is seeded.",
		s.resumeURL, "resume.pdf", now, now)
	if err != nil {
		return fmt.Errorf("seeker1 profile: %w", err)
	}

	// Apply to up to 3 published jobs.
	rows, err := s.pool.Query(ctx, `
		SELECT id FROM jobs WHERE status = 'published' ORDER BY random() LIMIT 3
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var jobID uuid.UUID
		if err := rows.Scan(&jobID); err != nil {
			return err
		}
		_, err := s.pool.Exec(ctx, `
			INSERT INTO applications (id, job_id, seeker_id, cover_letter, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, 'submitted', $5, $5)
			ON CONFLICT (job_id, seeker_id) DO NOTHING
		`, uuid.New(), jobID, userID, "Hi, I'm very interested in this role.", now)
		if err != nil {
			return fmt.Errorf("seeker1 application: %w", err)
		}
	}

	// Save 2 jobs.
	_, err = s.pool.Exec(ctx, `
		INSERT INTO saved_jobs (seeker_id, job_id, created_at)
		SELECT $1, id, $2 FROM jobs WHERE status = 'published' ORDER BY random() LIMIT 2
		ON CONFLICT DO NOTHING
	`, userID, now)
	if err != nil {
		return fmt.Errorf("seeker1 saved jobs: %w", err)
	}

	// Save one search.
	_, err = s.pool.Exec(ctx, `
		INSERT INTO saved_searches (id, seeker_id, name, query, remote_type, created_at)
		VALUES ($1, $2, 'Remote Go jobs', 'golang', 'remote', $3)
	`, uuid.New(), userID, now)
	if err != nil {
		return fmt.Errorf("seeker1 saved search: %w", err)
	}

	return nil
}

func (s *seeder) linkTestEmployer(ctx context.Context, userID uuid.UUID, now time.Time) error {
	companyID := uuid.New()

	var logoURL any
	if len(s.logoURLs) > 0 {
		logoURL = s.logoURLs[0]
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO companies (
			id, owner_user_id, name, description, website, location,
			logo_url, logo_filename, created_at, updated_at
		) VALUES ($1, $2, 'Test Employer Co', $3, 'https://test-employer.example.com',
		          'Helsinki', $4, 'logo-1.png', $5, $5)
	`, companyID, userID,
		"A company created for the employer1 test account. Feel free to explore.",
		logoURL, now)
	if err != nil {
		return fmt.Errorf("employer1 company: %w", err)
	}

	// Publish 3 jobs.
	titles := []string{
		"Senior Go Engineer (test)",
		"Backend Developer (test)",
		"Platform Engineer (test)",
	}
	for i, title := range titles {
		_, err := s.pool.Exec(ctx, `
			INSERT INTO jobs (
				id, company_id, created_by_user_id,
				title, description, requirements,
				location, remote_type, salary_min, salary_max, salary_currency,
				employment_type, category, tags,
				status, published_at, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, 'Helsinki', 'hybrid', 60000, 90000, 'EUR',
			          'full_time', 'Engineering', ARRAY['golang','postgres'],
			          'published', $7, $7, $7)
		`, uuid.New(), companyID, userID,
			title,
			"This is a seeded test job created by employer1. It exercises every screen end-to-end.",
			"- Go experience\n- Postgres\n- Distributed systems",
			now)
		if err != nil {
			return fmt.Errorf("employer1 job %d: %w", i, err)
		}
	}

	return nil
}
