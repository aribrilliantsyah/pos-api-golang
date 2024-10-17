package schemas

import "time"

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

// UserData digunakan untuk menampilkan data pengguna di response
type UserData struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	FullName     string    `json:"full_name"`
	CurrentToken string    `json:"current_token,omitempty"`
	CreatedBy    int64     `json:"created_by,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedBy    int64     `json:"updated_by,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	DeletedBy    int64     `json:"deleted_by,omitempty"`
	DeletedAt    time.Time `json:"deleted_at,omitempty"`
}

// CreateUser digunakan untuk payload pembuatan pengguna baru
type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUser digunakan untuk payload pembaruan pengguna
type UpdateUser struct {
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
}
