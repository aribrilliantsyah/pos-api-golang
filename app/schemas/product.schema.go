package schemas

import "time"

type ProductData struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int32     `json:"stock"`
	CategoryID int64     `json:"category_id"`
	CreatedBy  int64     `json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedBy  int64     `json:"updated_by,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	DeletedBy  int64     `json:"deleted_by,omitempty"`
	DeletedAt  time.Time `json:"deleted_at,omitempty"`
}

type CreateProduct struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID int64   `json:"category_id" binding:"required"`
}

type UpdateProduct struct {
	Name       string  `json:"name,omitempty"`
	Price      float64 `json:"price,omitempty"`
	CategoryID int64   `json:"category_id,omitempty"`
}
