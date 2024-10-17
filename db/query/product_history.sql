-- #PRODUCT_HISTORY

-- name: GetAllProductHistory :many
SELECT *
FROM product_history
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProductHistory :one
INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
RETURNING *;
