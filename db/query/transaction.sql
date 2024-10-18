-- name: CreateOrder :one
INSERT INTO orders (
    trx_number,
    cashier_id,
    customer_id,
    total_amount,
    payment_method,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    old_product,
    quantity,
    unit_price,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: CreateRefund :one
INSERT INTO refunds (
    order_id,
    reason,
    created_by
) VALUES (
    $1, $2, $3
) RETURNING *;

-- GetOrderByTrxNumber
-- name: GetOrderByTrxNumber :one
SELECT * FROM orders 
WHERE trx_number = $1 
LIMIT 1;

-- GetOrderItemsByOrderID
-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items 
WHERE order_id = $1;

-- UpdateOrderStatus
-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $1, 
    updated_by = $2, 
    updated_at = CURRENT_TIMESTAMP 
WHERE id = $3 
RETURNING *;
