-- +goose Up
CREATE TABLE saved_jobs (
    seeker_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id     UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (seeker_id, job_id)
);

CREATE INDEX idx_saved_jobs_seeker_id ON saved_jobs (seeker_id);
CREATE INDEX idx_saved_jobs_job_id ON saved_jobs (job_id);

-- +goose Down
DROP TABLE saved_jobs;
