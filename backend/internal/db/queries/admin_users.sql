-- name: ListUsers :many
SELECT
    u.id,
    u.username,
    u.email,
    u.role,
    u.onboarding_completed,
    u.created_at,
    sp.full_name   AS seeker_full_name,
    sp.headline    AS seeker_headline,
    c.name         AS company_name
FROM users u
LEFT JOIN seeker_profiles sp ON sp.user_id = u.id
LEFT JOIN companies c        ON c.owner_user_id = u.id
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsersFiltered :one
SELECT COUNT(*)::bigint FROM users u
WHERE (sqlc.narg('role')::text IS NULL OR u.role = sqlc.narg('role')::text);

-- name: ListUsersByRole :many
SELECT
    u.id,
    u.username,
    u.email,
    u.role,
    u.onboarding_completed,
    u.created_at,
    sp.full_name   AS seeker_full_name,
    sp.headline    AS seeker_headline,
    c.name         AS company_name
FROM users u
LEFT JOIN seeker_profiles sp ON sp.user_id = u.id
LEFT JOIN companies c        ON c.owner_user_id = u.id
WHERE u.role = $1
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserDetailByID :one
SELECT
    u.id,
    u.username,
    u.email,
    u.role,
    u.onboarding_completed,
    u.created_at,
    u.updated_at,
    sp.full_name      AS seeker_full_name,
    sp.phone          AS seeker_phone,
    sp.location       AS seeker_location,
    sp.headline       AS seeker_headline,
    sp.bio            AS seeker_bio,
    sp.created_at     AS seeker_created_at,
    c.id              AS company_id,
    c.name            AS company_name,
    c.description     AS company_description,
    c.website         AS company_website,
    c.location        AS company_location,
    c.logo_url        AS company_logo_url,
    c.created_at      AS company_created_at
FROM users u
LEFT JOIN seeker_profiles sp ON sp.user_id = u.id
LEFT JOIN companies c        ON c.owner_user_id = u.id
WHERE u.id = $1;
