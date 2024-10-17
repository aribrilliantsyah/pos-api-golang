package schemas

import "time"

// CreateCustomer digunakan untuk payload pembuatan customer baru
type CreateCustomer struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// UpdateCustomer digunakan untuk payload pembaruan customer
type UpdateCustomer struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// CustomerData digunakan untuk menampilkan data customer di response
type CustomerData struct {
	ID         int64     `json:"id"`
	MemberCode string    `json:"member_code"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone,omitempty"`
	Email      string    `json:"email,omitempty"`
	CreatedBy  int64     `json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedBy  int64     `json:"modified_by,omitempty"`
	UpdatedAt  time.Time `json:"modified_at,omitempty"`
	DeletedBy  int64     `json:"deleted_by,omitempty"`
	DeletedAt  time.Time `json:"deleted_at,omitempty"`
}
