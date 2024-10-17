package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"pos-api/app/schemas"
	db "pos-api/db/sqlc"
	"pos-api/util/common"
	"pos-api/util/jwt"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	db  *db.Queries
	ctx context.Context
}

func NewCategoryController(db *db.Queries, ctx context.Context) *CategoryController {
	return &CategoryController{db, ctx}
}

// CreateCategory godoc
// @Security BearerAuth
// @Summary Create a new category
// @Description Create a new category with the given payload
// @Tags categories
// @Accept json
// @Produce json
// @Param payload body schemas.CreateCategory true "Category Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var payload *schemas.CreateCategory

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	userInfo, err := jwt.GetUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	UserID := userInfo.UserID

	args := &db.CreateCategoryParams{
		Name:      payload.Name,
		CreatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	category, err := c.db.CreateCategory(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CategoryData{
		ID:        category.ID,
		Name:      category.Name,
		CreatedBy: common.ConvertNullInt64(category.CreatedBy),
		CreatedAt: common.ConvertNullTime(category.CreatedAt),
		UpdatedBy: common.ConvertNullInt64(category.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(category.UpdatedAt),
		DeletedAt: common.ConvertNullTime(category.DeletedAt),
		DeletedBy: common.ConvertNullInt64(category.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "created successfully",
		"data":    data,
	})
}

// UpdateCategory godoc
// @Security BearerAuth
// @Summary Update an existing category
// @Description Update a category with the given ID and payload
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param payload body schemas.UpdateCategory true "Category Update Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	var payload *schemas.UpdateCategory
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid category id",
		})
		return
	}

	userInfo, err := jwt.GetUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	UserID := userInfo.UserID

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if _, err := c.db.GetCategoryByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve category with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	args := &db.UpdateCategoryParams{
		ID:        id,
		Name:      payload.Name,
		UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	category, err := c.db.UpdateCategory(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve category with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CategoryData{
		ID:        category.ID,
		Name:      category.Name,
		CreatedBy: common.ConvertNullInt64(category.CreatedBy),
		CreatedAt: common.ConvertNullTime(category.CreatedAt),
		UpdatedBy: common.ConvertNullInt64(category.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(category.UpdatedAt),
		DeletedAt: common.ConvertNullTime(category.DeletedAt),
		DeletedBy: common.ConvertNullInt64(category.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "updated successfully",
		"data":    data,
	})
}

// GetCategoryById godoc
// @Security BearerAuth
// @Summary Get a category by ID
// @Description Retrieve a category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/{id} [get]
func (c *CategoryController) GetCategoryById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid category id",
		})
		return
	}

	category, err := c.db.GetCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve category with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CategoryData{
		ID:        category.ID,
		Name:      category.Name,
		CreatedBy: common.ConvertNullInt64(category.CreatedBy),
		CreatedAt: common.ConvertNullTime(category.CreatedAt),
		UpdatedBy: common.ConvertNullInt64(category.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(category.UpdatedAt),
		DeletedAt: common.ConvertNullTime(category.DeletedAt),
		DeletedBy: common.ConvertNullInt64(category.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllCategories godoc
// @Security BearerAuth
// @Summary Get all categories
// @Description Retrieve all categories with pagination
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories [get]
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllCategoriesParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	categories, err := c.db.GetAllCategories(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if categories == nil {
		categories = []db.Category{}
	}

	data := make([]schemas.CategoryData, len(categories))
	for i, category := range categories {
		data[i] = schemas.CategoryData{
			ID:        category.ID,
			Name:      category.Name,
			CreatedBy: common.ConvertNullInt64(category.CreatedBy),
			CreatedAt: common.ConvertNullTime(category.CreatedAt),
			UpdatedBy: common.ConvertNullInt64(category.UpdatedBy),
			UpdatedAt: common.ConvertNullTime(category.UpdatedAt),
			DeletedAt: common.ConvertNullTime(category.DeletedAt),
			DeletedBy: common.ConvertNullInt64(category.DeletedBy),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllDeletedCategories godoc
// @Security BearerAuth
// @Summary Get all deleted categories
// @Description Retrieve all categories with pagination
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/deleted [get]
func (c *CategoryController) GetAllDeletedCategories(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllDeletedCategoriesParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	categories, err := c.db.GetAllDeletedCategories(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if categories == nil {
		categories = []db.Category{}
	}

	data := make([]schemas.CategoryData, len(categories))
	for i, category := range categories {
		data[i] = schemas.CategoryData{
			ID:        category.ID,
			Name:      category.Name,
			CreatedBy: common.ConvertNullInt64(category.CreatedBy),
			CreatedAt: common.ConvertNullTime(category.CreatedAt),
			UpdatedBy: common.ConvertNullInt64(category.UpdatedBy),
			UpdatedAt: common.ConvertNullTime(category.UpdatedAt),
			DeletedAt: common.ConvertNullTime(category.DeletedAt),
			DeletedBy: common.ConvertNullInt64(category.DeletedBy),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// DeleteCategoryById godoc
// @Security BearerAuth
// @Summary Delete a category by ID
// @Description Delete a category with the given ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/{id} [delete]
func (c *CategoryController) DeleteCategoryById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid category id",
		})
		return
	}

	if _, err := c.db.GetCategoryByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve category with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	_, err = c.db.DeleteCategoryByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "deleted successfully",
	})
}

// SoftDeleteCategoryById godoc
// @Security BearerAuth
// @Summary Soft delete a category by ID
// @Description Soft delete a category with the given ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/{id}/soft [delete]
func (c *CategoryController) SoftDeleteCategoryById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid category id",
		})
		return
	}

	if _, err := c.db.GetCategoryByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve category with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	args := &db.SoftDeleteCategoryByIDParams{
		ID:        id,
		DeletedBy: sql.NullInt64{Int64: id, Valid: true},
	}
	_, err = c.db.SoftDeleteCategoryByID(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "soft deleted successfully",
	})
}
