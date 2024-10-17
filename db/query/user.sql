-- name: GetUserByUsername :one
SELECT id, username, password_hash, role
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByUsernameExceptID :one
SELECT id, username, password_hash, role
FROM users
WHERE username = $1 AND id != $2
LIMIT 1;

-- name: CheckToken :one
SELECT *
FROM users
WHERE username = $1 AND current_token = $2
LIMIT 1;

-- name: SetCurrentToken :one
UPDATE users
SET current_token = $2
WHERE username = $1
RETURNING *;

-- #USER

-- name: GetAllUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllDeletedUsers :many
SELECT *
FROM users
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (username, password_hash, role, full_name, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, full_name = $3, role = $4, updated_by = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateUserWithPassword :one
UPDATE users
SET username = $2, full_name = $3, role = $4, password_hash = $5, updated_by = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: SoftDeleteUserByID :one
UPDATE users
SET deleted_by = $2, deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUserByID :one
DELETE FROM users
WHERE id = $1
RETURNING id;
