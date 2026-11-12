-- name: ListAdminApplications :many
SELECT
    a.id,
    a.status,
    a.created_at,
    a.updated_at,
    a.job_id,
    a.seeker_id,
    u.username   AS seeker_username,
    u.email      AS seeker_email,
    sp.full_name AS seeker_full_name,
    sp.headline  AS seeker_headline,
    j.title      AS job_title,
    c.id         AS company_id,
    c.name       AS company_name
FROM applications a
JOIN users u                 ON u.id = a.seeker_id
JOIN jobs j                  ON j.id = a.job_id
JOIN companies c             ON c.id = j.company_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE (
    sqlc.narg('status')::text IS NULL
    OR a.status = sqlc.narg('status')::text
)
  AND (
    sqlc.narg('job_id')::uuid IS NULL
    OR a.job_id = sqlc.narg('job_id')::uuid
)
  AND (
    sqlc.narg('seeker_q')::text IS NULL
    OR u.username ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
    OR u.email ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
    OR sp.full_name ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
)
ORDER BY a.created_at DESC
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');

-- name: CountAdminApplications :one
SELECT COUNT(*)::bigint
FROM applications a
JOIN users u                 ON u.id = a.seeker_id
JOIN jobs j                  ON j.id = a.job_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE (
    sqlc.narg('status')::text IS NULL
    OR a.status = sqlc.narg('status')::text
)
  AND (
    sqlc.narg('job_id')::uuid IS NULL
    OR a.job_id = sqlc.narg('job_id')::uuid
)
  AND (
    sqlc.narg('seeker_q')::text IS NULL
    OR u.username ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
    OR u.email ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
    OR sp.full_name ILIKE '%' || sqlc.narg('seeker_q')::text || '%'
);

-- name: GetAdminApplicationByID :one
SELECT
    a.id,
    a.status,
    a.cover_letter,
    a.created_at,
    a.updated_at,
    a.job_id,
    a.seeker_id,
    u.username   AS seeker_username,
    u.email      AS seeker_email,
    sp.full_name AS seeker_full_name,
    sp.phone     AS seeker_phone,
    sp.location  AS seeker_location,
    sp.headline  AS seeker_headline,
    sp.bio       AS seeker_bio,
    sp.resume_url      AS seeker_resume_url,
    sp.resume_filename AS seeker_resume_filename,
    j.title      AS job_title,
    j.status     AS job_status,
    c.id         AS company_id,
    c.name       AS company_name
FROM applications a
JOIN users u                 ON u.id = a.seeker_id
JOIN jobs j                  ON j.id = a.job_id
JOIN companies c             ON c.id = j.company_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE a.id = $1;
