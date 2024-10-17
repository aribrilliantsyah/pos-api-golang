package schemas

import "time"

// CreateCategory digunakan untuk payload pembuatan kategori baru
type CreateCategory struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCategory digunakan untuk payload pembaruan kategori
type UpdateCategory struct {
	Name string `json:"name" binding:"required"`
}

// CategoryData digunakan untuk menampilkan data kategori di response
type CategoryData struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedBy int64     `json:"modified_by,omitempty"`
	UpdatedAt time.Time `json:"modified_at,omitempty"`
	DeletedBy int64     `json:"deleted_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}
