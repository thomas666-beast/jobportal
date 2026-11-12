package main

import (
	"context"

	"github.com/google/uuid"
)

func (s *seeder) seedCompanies(ctx context.Context, total int) error {
	if total > len(s.employerIDs) {
		total = len(s.employerIDs)
	}

	// Pick a random subset of employers to have companies.
	perm := s.r.Perm(len(s.employerIDs))
	chosen := make([]uuid.UUID, 0, total)
	for i := 0; i < total; i++ {
		chosen = append(chosen, s.employerIDs[perm[i]])
	}

	s.companyIDs = make([]uuid.UUID, 0, total)

	cols := []string{
		"id", "owner_user_id", "name", "description", "website",
		"location", "logo_url", "logo_filename", "created_at", "updated_at",
	}
	rows := make([][]any, 0, total)

	for _, ownerID := range chosen {
		id := uuid.New()
		s.companyIDs = append(s.companyIDs, id)

		name := pick(s.r, companyPrefixes) + " " + pick(s.r, companySuffixes)
		desc := pick(s.r, companyDescriptions)
		location := pick(s.r, cities)
		createdAt := s.randomTimestamp(365)

		// Mixed completeness — 70% have a logo, 90% have a website.
		var logoURL, logoFilename, website any
		if s.r.Float64() < 0.7 {
			idx := s.r.Intn(len(s.logoURLs))
			logoURL = s.logoURLs[idx]
			logoFilename = "logo-" + itoa(idx+1) + ".png"
		}
		if s.r.Float64() < 0.9 {
			slug := toSlug(name)
			website = "https://" + slug + ".example.com"
		}

		rows = append(rows, []any{
			id,
			ownerID,
			name,
			desc,
			website,
			location,
			logoURL,
			logoFilename,
			createdAt,
			createdAt,
		})
	}

	return s.batchCopy(ctx, "companies", cols, rows)
}

func toSlug(s string) string {
	out := make([]byte, 0, len(s))
	prevDash := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
			prevDash = false
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			out = append(out, c)
			prevDash = false
		default:
			if !prevDash && len(out) > 0 {
				out = append(out, '-')
				prevDash = true
			}
		}
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}
	return string(out)
}
