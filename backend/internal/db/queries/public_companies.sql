-- name: ListCompanies :many
SELECT
    c.id,
    c.name,
    c.description,
    c.website,
    c.location,
    c.logo_url,
    (SELECT COUNT(*)::bigint
       FROM jobs j
      WHERE j.company_id = c.id
        AND j.status = 'published') AS published_job_count
FROM companies c
WHERE (
    sqlc.narg('q')::text IS NULL
    OR c.name ILIKE '%' || sqlc.narg('q')::text || '%'
)
ORDER BY c.created_at DESC
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');

-- name: CountCompanies :one
SELECT COUNT(*)::bigint FROM companies c
WHERE (
    sqlc.narg('q')::text IS NULL
    OR c.name ILIKE '%' || sqlc.narg('q')::text || '%'
);

-- name: GetCompanyPublicByID :one
SELECT
    c.id,
    c.name,
    c.description,
    c.website,
    c.location,
    c.logo_url,
    c.created_at,
    (SELECT COUNT(*)::bigint
       FROM jobs j
      WHERE j.company_id = c.id
        AND j.status = 'published') AS published_job_count
FROM companies c
WHERE c.id = $1;

-- name: ListPublishedJobsByCompany :many
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
    j.published_at
FROM jobs j
WHERE j.company_id = $1
  AND j.status = 'published'
ORDER BY j.published_at DESC NULLS LAST, j.created_at DESC;
