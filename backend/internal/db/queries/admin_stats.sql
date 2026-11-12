-- name: CountUsersByRole :many
SELECT role, COUNT(*)::bigint AS count
FROM users
GROUP BY role;

-- name: CountJobsByStatusAdmin :many
SELECT status, COUNT(*)::bigint AS count
FROM jobs
GROUP BY status;

-- name: CountApplicationsByStatus :many
SELECT status, COUNT(*)::bigint AS count
FROM applications
GROUP BY status;

-- name: CountAllUsers :one
SELECT COUNT(*)::bigint FROM users;

-- name: CountAllCompanies :one
SELECT COUNT(*)::bigint FROM companies;

-- name: CountAllJobs :one
SELECT COUNT(*)::bigint FROM jobs;

-- name: CountAllApplications :one
SELECT COUNT(*)::bigint FROM applications;

-- name: ListRecentUsers :many
SELECT id, username, email, role, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1;

-- name: ListRecentApplications :many
SELECT
    a.id,
    a.status,
    a.created_at,
    u.username   AS seeker_username,
    j.title      AS job_title,
    c.name       AS company_name
FROM applications a
JOIN users u      ON u.id = a.seeker_id
JOIN jobs j       ON j.id = a.job_id
JOIN companies c  ON c.id = j.company_id
ORDER BY a.created_at DESC
LIMIT $1;
