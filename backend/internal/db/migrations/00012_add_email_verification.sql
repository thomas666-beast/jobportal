-- +goose Up
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ;

-- Mark all existing users as verified so they don't suddenly need to re-verify.
UPDATE users SET email_verified_at = now() WHERE email_verified_at IS NULL;

CREATE TABLE auth_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    purpose    TEXT NOT NULL CHECK (purpose IN ('email_verification', 'password_reset')),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_auth_tokens_user_purpose ON auth_tokens (user_id, purpose);
CREATE INDEX idx_auth_tokens_expires_at ON auth_tokens (expires_at);

-- +goose Down
DROP TABLE auth_tokens;
ALTER TABLE users DROP COLUMN email_verified_at;
