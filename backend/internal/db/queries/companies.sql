-- name: GetCompanyByOwner :one
SELECT id, owner_user_id, name, description, website, location, logo_url, logo_filename, created_at, updated_at
FROM companies
WHERE owner_user_id = $1;

-- name: GetCompanyByID :one
SELECT id, owner_user_id, name, description, website, location, logo_url, logo_filename, created_at, updated_at
FROM companies
WHERE id = $1;

-- name: UpsertCompany :one
INSERT INTO companies (owner_user_id, name, description, website, location, logo_url)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (owner_user_id) DO UPDATE
SET name        = EXCLUDED.name,
    description = EXCLUDED.description,
    website     = EXCLUDED.website,
    location    = EXCLUDED.location,
    logo_url    = EXCLUDED.logo_url,
    updated_at  = now()
RETURNING id, owner_user_id, name, description, website, location, logo_url, logo_filename, created_at, updated_at;
