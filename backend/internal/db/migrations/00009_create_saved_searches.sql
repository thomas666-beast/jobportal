-- +goose Up
CREATE TABLE saved_searches (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seeker_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    query       TEXT,
    location    TEXT,
    remote_type TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (seeker_id, name)
);

CREATE INDEX idx_saved_searches_seeker_id ON saved_searches (seeker_id);

-- +goose Down
DROP TABLE saved_searches;
