-- +goose Up
CREATE TABLE jobs (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id         UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title              TEXT NOT NULL,
    description        TEXT NOT NULL,
    requirements       TEXT,
    location           TEXT,
    remote_type        TEXT NOT NULL DEFAULT 'on_site'
        CHECK (remote_type IN ('remote', 'hybrid', 'on_site')),
    salary_min         INT,
    salary_max         INT,
    salary_currency    TEXT,
    employment_type    TEXT NOT NULL DEFAULT 'full_time'
        CHECK (employment_type IN ('full_time', 'part_time', 'contract', 'internship')),
    category           TEXT,
    tags               TEXT[] NOT NULL DEFAULT '{}',
    status             TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'pending_review', 'published', 'closed', 'archived')),
    rejection_reason   TEXT,
    published_at       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (salary_min IS NULL OR salary_max IS NULL OR salary_min <= salary_max)
);

CREATE INDEX idx_jobs_company_id ON jobs (company_id);
CREATE INDEX idx_jobs_status ON jobs (status);
CREATE INDEX idx_jobs_published_at ON jobs (published_at DESC);

-- +goose Down
DROP TABLE jobs;
