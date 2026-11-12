-- name: SaveJob :exec
INSERT INTO saved_jobs (seeker_id, job_id)
VALUES ($1, $2)
ON CONFLICT (seeker_id, job_id) DO NOTHING;

-- name: UnsaveJob :exec
DELETE FROM saved_jobs
WHERE seeker_id = $1 AND job_id = $2;

-- name: IsJobSaved :one
SELECT EXISTS (
    SELECT 1 FROM saved_jobs
    WHERE seeker_id = $1 AND job_id = $2
) AS saved;

-- name: ListSavedJobsBySeeker :many
SELECT
    s.job_id,
    s.created_at AS saved_at,
    j.title,
    j.description,
    j.location,
    j.remote_type,
    j.salary_min,
    j.salary_max,
    j.salary_currency,
    j.employment_type,
    j.category,
    j.tags,
    j.published_at,
    j.status AS job_status,
    c.id       AS company_id,
    c.name     AS company_name,
    c.location AS company_location,
    c.logo_url AS company_logo_url
FROM saved_jobs s
JOIN jobs j      ON j.id = s.job_id
JOIN companies c ON c.id = j.company_id
WHERE s.seeker_id = $1
  AND j.status = 'published'
ORDER BY s.created_at DESC;
