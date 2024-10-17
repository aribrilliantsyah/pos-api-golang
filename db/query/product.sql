-- #PRODUCT

-- name: GetAllProducts :many
SELECT * 
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllDeletedProducts :many
SELECT * 
FROM products
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products (name, price, stock, category_id, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, stock = $4, category_id = $5, updated_by = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: SoftDeleteProductByID :one
UPDATE products
SET deleted_by = $2, deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteProductByID :one
DELETE FROM products
WHERE id = $1
RETURNING id;
