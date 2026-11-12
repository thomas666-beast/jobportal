-- +goose Up
ALTER TABLE users
    ADD COLUMN role TEXT NOT NULL DEFAULT 'job_seeker'
        CHECK (role IN ('job_seeker', 'employer', 'admin')),
    ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX idx_users_role ON users (role);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;

ALTER TABLE users
    DROP COLUMN onboarding_completed,
    DROP COLUMN role;
