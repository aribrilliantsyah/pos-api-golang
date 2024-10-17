// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: category.sql

package db

import (
	"context"
	"database/sql"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (name, created_by, created_at)
VALUES ($1, $2, CURRENT_TIMESTAMP)
RETURNING id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
`

type CreateCategoryParams struct {
	Name      string        `json:"name"`
	CreatedBy sql.NullInt64 `json:"created_by"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.queryRow(ctx, q.createCategoryStmt, createCategory, arg.Name, arg.CreatedBy)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteCategoryByID = `-- name: DeleteCategoryByID :one
DELETE FROM categories
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteCategoryByID(ctx context.Context, id int64) (int64, error) {
	row := q.queryRow(ctx, q.deleteCategoryByIDStmt, deleteCategoryByID, id)
	err := row.Scan(&id)
	return id, err
}

const getAllCategories = `-- name: GetAllCategories :many

SELECT id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
FROM categories
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetAllCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// #CATEGORY
func (q *Queries) GetAllCategories(ctx context.Context, arg GetAllCategoriesParams) ([]Category, error) {
	rows, err := q.query(ctx, q.getAllCategoriesStmt, getAllCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllDeletedCategories = `-- name: GetAllDeletedCategories :many
SELECT id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
FROM categories
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetAllDeletedCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllDeletedCategories(ctx context.Context, arg GetAllDeletedCategoriesParams) ([]Category, error) {
	rows, err := q.query(ctx, q.getAllDeletedCategoriesStmt, getAllDeletedCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategoryByID = `-- name: GetCategoryByID :one
SELECT id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
FROM categories
WHERE id = $1
`

func (q *Queries) GetCategoryByID(ctx context.Context, id int64) (Category, error) {
	row := q.queryRow(ctx, q.getCategoryByIDStmt, getCategoryByID, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const softDeleteCategoryByID = `-- name: SoftDeleteCategoryByID :one
UPDATE categories
SET deleted_by = $2, deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
`

type SoftDeleteCategoryByIDParams struct {
	ID        int64         `json:"id"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

func (q *Queries) SoftDeleteCategoryByID(ctx context.Context, arg SoftDeleteCategoryByIDParams) (Category, error) {
	row := q.queryRow(ctx, q.softDeleteCategoryByIDStmt, softDeleteCategoryByID, arg.ID, arg.DeletedBy)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories
SET name = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at
`

type UpdateCategoryParams struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.queryRow(ctx, q.updateCategoryStmt, updateCategory, arg.ID, arg.Name, arg.UpdatedBy)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
