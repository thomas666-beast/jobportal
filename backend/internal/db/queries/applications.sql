-- name: CreateApplication :one
INSERT INTO applications (job_id, seeker_id, cover_letter)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetApplicationByID :one
SELECT * FROM applications
WHERE id = $1;

-- name: GetApplicationByJobAndSeeker :one
SELECT * FROM applications
WHERE job_id = $1 AND seeker_id = $2;

-- name: ListApplicationsBySeeker :many
SELECT
    a.id,
    a.job_id,
    a.seeker_id,
    a.cover_letter,
    a.status,
    a.created_at,
    a.updated_at,
    j.title         AS job_title,
    j.location      AS job_location,
    j.remote_type   AS job_remote_type,
    j.status        AS job_status,
    c.id            AS company_id,
    c.name          AS company_name,
    c.logo_url      AS company_logo_url
FROM applications a
JOIN jobs j       ON j.id = a.job_id
JOIN companies c  ON c.id = j.company_id
WHERE a.seeker_id = $1
ORDER BY a.created_at DESC;

-- name: ListApplicationsByCompany :many
SELECT
    a.id,
    a.job_id,
    a.seeker_id,
    a.cover_letter,
    a.status,
    a.created_at,
    a.updated_at,
    j.title      AS job_title,
    u.username   AS seeker_username,
    u.email      AS seeker_email,
    sp.full_name AS seeker_full_name,
    sp.headline  AS seeker_headline,
    sp.location  AS seeker_location
FROM applications a
JOIN jobs j             ON j.id = a.job_id
JOIN users u            ON u.id = a.seeker_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE j.company_id = $1
ORDER BY a.created_at DESC;

-- name: ListApplicationsByJob :many
SELECT
    a.id,
    a.job_id,
    a.seeker_id,
    a.cover_letter,
    a.status,
    a.created_at,
    a.updated_at,
    u.username   AS seeker_username,
    u.email      AS seeker_email,
    sp.full_name AS seeker_full_name,
    sp.headline  AS seeker_headline,
    sp.location  AS seeker_location
FROM applications a
JOIN users u                 ON u.id = a.seeker_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE a.job_id = $1
ORDER BY a.created_at DESC;

-- name: GetApplicationDetailForEmployer :one
SELECT
    a.id,
    a.job_id,
    a.seeker_id,
    a.cover_letter,
    a.status,
    a.created_at,
    a.updated_at,
    j.title      AS job_title,
    j.company_id AS company_id,
    u.username   AS seeker_username,
    u.email      AS seeker_email,
    sp.full_name AS seeker_full_name,
    sp.phone     AS seeker_phone,
    sp.location  AS seeker_location,
    sp.headline  AS seeker_headline,
    sp.bio       AS seeker_bio,
    sp.resume_url      AS seeker_resume_url,
    sp.resume_filename AS seeker_resume_filename
FROM applications a
JOIN jobs j                  ON j.id = a.job_id
JOIN users u                 ON u.id = a.seeker_id
LEFT JOIN seeker_profiles sp ON sp.user_id = a.seeker_id
WHERE a.id = $1;

-- name: UpdateApplicationStatus :one
UPDATE applications
SET status     = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
