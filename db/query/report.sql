-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetAllOrders :many
SELECT 
    o.*,
    c.name as customer_name,
    u.username as cashier_name
FROM orders o
LEFT JOIN customers c ON o.customer_id = c.id
LEFT JOIN users u ON o.cashier_id = u.id
ORDER BY o.order_date DESC
LIMIT $1 OFFSET $2;

-- name: GetFastMovingProducts :many

SELECT 
    p.id,
    p.name,
    p.category_id,
    c.name as category_name,
    COALESCE(SUM(oi.quantity), 0) as total_quantity,
    COALESCE(COUNT(DISTINCT o.id), 0) as total_orders
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id
LEFT JOIN orders o ON oi.order_id = o.id AND
    (o.order_date >= $1 AND  o.order_date < $2) AND
    o.status = 'order'
LEFT JOIN categories c ON p.category_id = c.id
GROUP BY p.id, p.name, p.category_id, c.name
ORDER BY total_quantity DESC
LIMIT 10;

-- name: GetSlowMovingProducts :many
SELECT 
    p.id,
    p.name,
    p.category_id,
    c.name as category_name,
    COALESCE(SUM(oi.quantity), 0) as total_quantity,
    COALESCE(COUNT(DISTINCT o.id), 0) as total_orders
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id
LEFT JOIN orders o ON oi.order_id = o.id AND
    (o.order_date >= $1 AND  o.order_date < $2) AND
    o.status = 'order'
LEFT JOIN categories c ON p.category_id = c.id
GROUP BY p.id, p.name, p.category_id, c.name
ORDER BY total_quantity ASC
LIMIT 10;

-- name: GetTopCashiers :many
SELECT 
    u.id as cashier_id,
    u.username,
    u.full_name,
    COUNT(DISTINCT o.id) as total_orders,
    SUM(o.total_amount) as total_amount,
    AVG(o.total_amount) as average_order_amount
FROM users u
JOIN orders o ON u.id = o.cashier_id
WHERE 
    (o.order_date >= $1 AND  o.order_date < $2) AND
    o.status = 'order'
GROUP BY u.id, u.username, u.full_name
ORDER BY total_amount DESC
LIMIT 10;

-- name: GetTopCustomers :many
SELECT 
    c.id as customer_id,
    c.member_code,
    c.name,
    c.phone,
    c.email,
    COUNT(DISTINCT o.id) as total_orders,
    SUM(o.total_amount)::DECIMAL as total_spent,
    AVG(o.total_amount) as average_order_amount
FROM customers c
JOIN orders o ON c.id = o.customer_id
WHERE 
    (o.order_date >= $1 AND  o.order_date < $2) AND
    o.status = 'order'
GROUP BY c.id, c.member_code, c.name, c.phone, c.email
ORDER BY total_spent DESC
LIMIT 10;