package schemas

import (
	"database/sql"
	"time"
)

type OrderDetailResponse struct {
	ID            int64             `json:"id"`
	TrxNumber     string            `json:"trx_number"`
	CashierID     int64             `json:"cashier_id"`
	CustomerID    sql.NullInt64     `json:"customer_id"`
	TotalAmount   float64           `json:"total_amount"`
	PaymentMethod string            `json:"payment_method"`
	Status        string            `json:"status"`
	OrderDate     time.Time         `json:"order_date"`
	UpdatedBy     sql.NullInt64     `json:"updated_by"`
	UpdatedAt     sql.NullTime      `json:"updated_at"`
	Customer      *CustomerResponse `json:"customer,omitempty"`
	OrderItems    []OrderItemDetail `json:"order_items"`
}

type CustomerResponse struct {
	ID         int64  `json:"id"`
	MemberCode string `json:"member_code"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type OrderItemDetail struct {
	ID         int64   `json:"id"`
	ProductID  int64   `json:"product_id"`
	OldProduct string  `json:"old_product"`
	Quantity   int32   `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

type OrderListParams struct {
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	CustomerID sql.NullInt64  `json:"customer_id"`
	CashierID  sql.NullInt64  `json:"cashier_id"`
	Status     sql.NullString `json:"status"`
}
