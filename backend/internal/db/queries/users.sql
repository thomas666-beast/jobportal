-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email, role, onboarding_completed, created_at;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, role, onboarding_completed, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, username, email, role, onboarding_completed, created_at
FROM users
WHERE id = $1;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2,
    onboarding_completed = true,
    updated_at = now()
WHERE id = $1
RETURNING id, username, email, role, onboarding_completed, created_at;

-- name: GetUserWithPasswordByID :one
SELECT id, username, email, password_hash, role, onboarding_completed, created_at
FROM users
WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2,
    updated_at    = now()
WHERE id = $1
RETURNING id, username, email, role, onboarding_completed, created_at;

-- name: SetUserPasswordByID :exec
UPDATE users
SET password_hash = $2,
    updated_at    = now()
WHERE id = $1;

-- name: GetUserEmailVerifiedByID :one
SELECT email_verified_at FROM users WHERE id = $1;
