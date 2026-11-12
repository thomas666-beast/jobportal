-- +goose Up
CREATE TABLE applications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id       UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    seeker_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cover_letter TEXT,
    status       TEXT NOT NULL DEFAULT 'submitted'
        CHECK (status IN ('submitted', 'viewed', 'shortlisted', 'rejected', 'hired', 'withdrawn')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (job_id, seeker_id)
);

CREATE INDEX idx_applications_job_id ON applications (job_id);
CREATE INDEX idx_applications_seeker_id ON applications (seeker_id);
CREATE INDEX idx_applications_status ON applications (status);

-- +goose Down
DROP TABLE applications;
