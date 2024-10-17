-- #CUSTOMER

-- name: GetAllCustomers :many
SELECT *
FROM customers
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllDeletedCustomers :many
SELECT *
FROM customers
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: CreateCustomer :one
INSERT INTO customers (member_code, name, phone, email, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET member_code = $2, name = $3, phone = $4, email = $5, updated_by = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetCustomerByID :one
SELECT *
FROM customers
WHERE id = $1
LIMIT 1;

-- name: GetCustomerByEmail :one
SELECT *
FROM customers
WHERE email = $1;

-- name: GetCustomerByEmailExceptID :one
SELECT *
FROM customers
WHERE email = $1 AND id != $2
LIMIT 1; 

-- name: GetCustomerByPhone :one
SELECT *
FROM customers
WHERE phone = $1;

-- name: GetCustomerByPhoneExceptID :one
SELECT *
FROM customers
WHERE phone = $1 AND id != $2
LIMIT 1;

-- name: SoftDeleteCustomerByID :one
UPDATE customers
SET deleted_by = $2, deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCustomerByID :one
DELETE FROM customers
WHERE id = $1
RETURNING id;
