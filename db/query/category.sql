-- #CATEGORY

-- name: GetAllCategories :many
SELECT *
FROM categories
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetAllDeletedCategories :many
SELECT *
FROM categories
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateCategory :one
INSERT INTO categories (name, created_by, created_at)
VALUES ($1, $2, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetCategoryByID :one
SELECT *
FROM categories
WHERE id = $1;

-- name: SoftDeleteCategoryByID :one
UPDATE categories
SET deleted_by = $2, deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCategoryByID :one
DELETE FROM categories
WHERE id = $1
RETURNING id;
