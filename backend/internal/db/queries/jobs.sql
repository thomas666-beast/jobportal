-- name: CreateJob :one
INSERT INTO jobs (
    company_id, created_by_user_id, title, description, requirements,
    location, remote_type, salary_min, salary_max, salary_currency,
    employment_type, category, tags
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetJobByID :one
SELECT * FROM jobs
WHERE id = $1;

-- name: ListJobsByCompany :many
SELECT * FROM jobs
WHERE company_id = $1
ORDER BY created_at DESC;

-- name: UpdateJob :one
UPDATE jobs
SET title           = $2,
    description     = $3,
    requirements    = $4,
    location        = $5,
    remote_type     = $6,
    salary_min      = $7,
    salary_max      = $8,
    salary_currency = $9,
    employment_type = $10,
    category        = $11,
    tags            = $12,
    updated_at      = now()
WHERE id = $1
RETURNING *;

-- name: UpdateJobStatus :one
UPDATE jobs
SET status           = $2,
    rejection_reason = $3,
    published_at     = CASE
                         WHEN $2 = 'published' AND published_at IS NULL THEN now()
                         ELSE published_at
                       END,
    updated_at       = now()
WHERE id = $1
RETURNING *;

-- name: ListJobsByStatus :many
SELECT * FROM jobs
WHERE status = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountJobsByStatus :one
SELECT COUNT(*) FROM jobs
WHERE status = $1;

-- name: ListPublishedJobs :many
SELECT
    j.id,
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
    c.id           AS company_id,
    c.name         AS company_name,
    c.location     AS company_location,
    c.logo_url     AS company_logo_url
FROM jobs j
JOIN companies c ON c.id = j.company_id
WHERE j.status = 'published'
  AND (
    sqlc.narg('q')::text IS NULL
    OR j.search_vector @@ to_tsquery('simple',
        regexp_replace(sqlc.narg('q')::text, '\s+', ':* & ', 'g') || ':*'
    )
  )
  AND (
    sqlc.narg('location')::text IS NULL
    OR j.location ILIKE '%' || sqlc.narg('location')::text || '%'
  )
  AND (
    sqlc.narg('remote_type')::text IS NULL
    OR j.remote_type = sqlc.narg('remote_type')::text
  )
ORDER BY
    CASE WHEN sqlc.narg('q')::text IS NULL THEN NULL
         ELSE ts_rank(j.search_vector, to_tsquery('simple',
             regexp_replace(sqlc.narg('q')::text, '\s+', ':* & ', 'g') || ':*'
         ))
    END DESC NULLS LAST,
    j.published_at DESC NULLS LAST,
    j.created_at DESC
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');

-- name: CountPublishedJobs :one
SELECT COUNT(*) FROM jobs j
WHERE j.status = 'published'
  AND (
    sqlc.narg('q')::text IS NULL
    OR j.search_vector @@ to_tsquery('simple',
        regexp_replace(sqlc.narg('q')::text, '\s+', ':* & ', 'g') || ':*'
    )
  )
  AND (
    sqlc.narg('location')::text IS NULL
    OR j.location ILIKE '%' || sqlc.narg('location')::text || '%'
  )
  AND (
    sqlc.narg('remote_type')::text IS NULL
    OR j.remote_type = sqlc.narg('remote_type')::text
  );

-- name: GetPublishedJobByID :one
SELECT
    j.id,
    j.title,
    j.description,
    j.requirements,
    j.location,
    j.remote_type,
    j.salary_min,
    j.salary_max,
    j.salary_currency,
    j.employment_type,
    j.category,
    j.tags,
    j.published_at,
    j.created_at,
    c.id           AS company_id,
    c.name         AS company_name,
    c.description  AS company_description,
    c.website      AS company_website,
    c.location     AS company_location,
    c.logo_url     AS company_logo_url
FROM jobs j
JOIN companies c ON c.id = j.company_id
WHERE j.id = $1
  AND j.status = 'published';

-- name: GetJobWithCompanyByID :one
SELECT
    j.id,
    j.company_id,
    j.title,
    j.description,
    j.requirements,
    j.location,
    j.remote_type,
    j.salary_min,
    j.salary_max,
    j.salary_currency,
    j.employment_type,
    j.category,
    j.tags,
    j.status,
    j.rejection_reason,
    j.published_at,
    j.created_at,
    j.updated_at,
    c.name            AS company_name,
    c.description     AS company_description,
    c.website         AS company_website,
    c.location        AS company_location,
    c.logo_url        AS company_logo_url,
    u.username        AS employer_username,
    u.email           AS employer_email
FROM jobs j
JOIN companies c ON c.id = j.company_id
LEFT JOIN users u ON u.id = c.owner_user_id
WHERE j.id = $1;
