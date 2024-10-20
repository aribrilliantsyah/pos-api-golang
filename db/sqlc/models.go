// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Category struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
	DeletedAt sql.NullTime  `json:"deleted_at"`
}

type Customer struct {
	ID         int64          `json:"id"`
	MemberCode string         `json:"member_code"`
	Name       string         `json:"name"`
	Phone      sql.NullString `json:"phone"`
	Email      sql.NullString `json:"email"`
	CreatedBy  sql.NullInt64  `json:"created_by"`
	UpdatedBy  sql.NullInt64  `json:"updated_by"`
	DeletedBy  sql.NullInt64  `json:"deleted_by"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
}

type Order struct {
	ID            int64         `json:"id"`
	TrxNumber     string        `json:"trx_number"`
	CashierID     sql.NullInt64 `json:"cashier_id"`
	CustomerID    sql.NullInt64 `json:"customer_id"`
	TotalAmount   string        `json:"total_amount"`
	PaymentMethod string        `json:"payment_method"`
	Status        string        `json:"status"`
	OrderDate     sql.NullTime  `json:"order_date"`
	UpdatedBy     sql.NullInt64 `json:"updated_by"`
	UpdatedAt     sql.NullTime  `json:"updated_at"`
}

type OrderItem struct {
	ID         int64          `json:"id"`
	OrderID    sql.NullInt64  `json:"order_id"`
	ProductID  sql.NullInt64  `json:"product_id"`
	OldProduct sql.NullString `json:"old_product"`
	Quantity   int32          `json:"quantity"`
	UnitPrice  string         `json:"unit_price"`
	CreatedBy  sql.NullInt64  `json:"created_by"`
	CreatedAt  sql.NullTime   `json:"created_at"`
}

type Product struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	Price      string        `json:"price"`
	Stock      int32         `json:"stock"`
	CategoryID sql.NullInt64 `json:"category_id"`
	CreatedBy  sql.NullInt64 `json:"created_by"`
	UpdatedBy  sql.NullInt64 `json:"updated_by"`
	DeletedBy  sql.NullInt64 `json:"deleted_by"`
	CreatedAt  sql.NullTime  `json:"created_at"`
	UpdatedAt  sql.NullTime  `json:"updated_at"`
	DeletedAt  sql.NullTime  `json:"deleted_at"`
}

type ProductHistory struct {
	ID             int64          `json:"id"`
	TrxRef         string         `json:"trx_ref"`
	ProductID      sql.NullInt64  `json:"product_id"`
	QuantityChange int32          `json:"quantity_change"`
	Type           string         `json:"type"`
	Reason         sql.NullString `json:"reason"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	CreatedAt      sql.NullTime   `json:"created_at"`
}

type Refund struct {
	ID        int64         `json:"id"`
	OrderID   sql.NullInt64 `json:"order_id"`
	Reason    string        `json:"reason"`
	RefundAt  sql.NullTime  `json:"refund_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
}

type User struct {
	ID           int64          `json:"id"`
	Username     string         `json:"username"`
	PasswordHash string         `json:"password_hash"`
	Role         string         `json:"role"`
	FullName     string         `json:"full_name"`
	CurrentToken sql.NullString `json:"current_token"`
	CreatedBy    sql.NullInt64  `json:"created_by"`
	UpdatedBy    sql.NullInt64  `json:"updated_by"`
	DeletedBy    sql.NullInt64  `json:"deleted_by"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}
