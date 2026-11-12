-- name: CreateSavedSearch :one
INSERT INTO saved_searches (seeker_id, name, query, location, remote_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListSavedSearchesBySeeker :many
SELECT * FROM saved_searches
WHERE seeker_id = $1
ORDER BY created_at DESC;

-- name: GetSavedSearchByID :one
SELECT * FROM saved_searches
WHERE id = $1;

-- name: DeleteSavedSearch :exec
DELETE FROM saved_searches
WHERE id = $1 AND seeker_id = $2;
