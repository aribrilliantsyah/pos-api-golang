package schemas

import "time"

type ProductHistoryData struct {
	ID             int64     `json:"id"`
	TrxRef         string    `json:"trx_ref"`
	ProductID      int64     `json:"product_id"`
	QuantityChange int32     `json:"quantity_change"`
	Type           string    `json:"type"` // "in" or "out"
	Reason         string    `json:"reason"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateProductHistory struct {
	ProductID      int64  `json:"product_id" binding:"required"`
	QuantityChange int32  `json:"quantity_change" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=in out"`
	Reason         string `json:"reason" binding:"required"`
}
