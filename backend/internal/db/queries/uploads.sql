-- name: SetSeekerResume :one
UPDATE seeker_profiles
SET resume_url         = $2,
    resume_filename    = $3,
    resume_uploaded_at = now(),
    updated_at         = now()
WHERE user_id = $1
RETURNING user_id, full_name, phone, location, headline, bio, resume_url, resume_filename, resume_uploaded_at, created_at, updated_at;

-- name: ClearSeekerResume :one
UPDATE seeker_profiles
SET resume_url         = NULL,
    resume_filename    = NULL,
    resume_uploaded_at = NULL,
    updated_at         = now()
WHERE user_id = $1
RETURNING user_id, full_name, phone, location, headline, bio, resume_url, resume_filename, resume_uploaded_at, created_at, updated_at;

-- name: SetCompanyLogo :one
UPDATE companies
SET logo_url      = $2,
    logo_filename = $3,
    updated_at    = now()
WHERE owner_user_id = $1
RETURNING id, owner_user_id, name, description, website, location, logo_url, logo_filename, created_at, updated_at;

-- name: ClearCompanyLogo :one
UPDATE companies
SET logo_url      = NULL,
    logo_filename = NULL,
    updated_at    = now()
WHERE owner_user_id = $1
RETURNING id, owner_user_id, name, description, website, location, logo_url, logo_filename, created_at, updated_at;
