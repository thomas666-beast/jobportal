-- name: CreateAuthToken :one
INSERT INTO auth_tokens (user_id, token_hash, purpose, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAuthTokenByHash :one
SELECT * FROM auth_tokens
WHERE token_hash = $1 AND purpose = $2;

-- name: MarkAuthTokenUsed :exec
UPDATE auth_tokens
SET used_at = now()
WHERE id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM auth_tokens WHERE expires_at < now() OR used_at IS NOT NULL;

-- name: SetUserEmailVerified :one
UPDATE users
SET email_verified_at = now(),
    updated_at = now()
WHERE id = $1
RETURNING id, username, email, role, onboarding_completed, email_verified_at, created_at;
