-- +goose Up
CREATE TABLE seeker_profiles (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    full_name   TEXT NOT NULL,
    phone       TEXT,
    location    TEXT,
    headline    TEXT,
    bio         TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE seeker_profiles;
