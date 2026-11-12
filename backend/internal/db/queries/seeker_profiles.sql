-- name: GetSeekerProfile :one
SELECT user_id, full_name, phone, location, headline, bio, resume_url, resume_filename, resume_uploaded_at, created_at, updated_at
FROM seeker_profiles
WHERE user_id = $1;

-- name: UpsertSeekerProfile :one
INSERT INTO seeker_profiles (user_id, full_name, phone, location, headline, bio)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id) DO UPDATE
SET full_name  = EXCLUDED.full_name,
    phone      = EXCLUDED.phone,
    location   = EXCLUDED.location,
    headline   = EXCLUDED.headline,
    bio        = EXCLUDED.bio,
    updated_at = now()
RETURNING user_id, full_name, phone, location, headline, bio, resume_url, resume_filename, resume_uploaded_at, created_at, updated_at;
