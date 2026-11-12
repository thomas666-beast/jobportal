-- +goose Up
ALTER TABLE seeker_profiles
    ADD COLUMN resume_url         TEXT,
    ADD COLUMN resume_filename    TEXT,
    ADD COLUMN resume_uploaded_at TIMESTAMPTZ;

ALTER TABLE companies
    ADD COLUMN logo_filename TEXT;

-- +goose Down
ALTER TABLE companies
    DROP COLUMN logo_filename;

ALTER TABLE seeker_profiles
    DROP COLUMN resume_uploaded_at,
    DROP COLUMN resume_filename,
    DROP COLUMN resume_url;
