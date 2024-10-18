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

type ProductController struct {
	db  *db.Queries
	ctx context.Context
}

func NewProductController(db *db.Queries, ctx context.Context) *ProductController {
	return &ProductController{db, ctx}
}

// CreateProduct godoc
// @Security BearerAuth
// @Summary Create a new product
// @Description Create a new product with the given payload
// @Tags products
// @Accept json
// @Produce json
// @Param payload body schemas.CreateProduct true "Product Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products [post]
func (p *ProductController) CreateProduct(ctx *gin.Context) {
	var payload *schemas.CreateProduct

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

	if _, err := p.db.GetCategoryByID(ctx, payload.CategoryID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "category id not found",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if payload.Price < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid price (need positive number/decimal)",
		})
		return
	}

	args := &db.CreateProductParams{
		Name:       payload.Name,
		Price:      strconv.FormatFloat(payload.Price, 'f', 2, 64),
		CategoryID: sql.NullInt64{Int64: payload.CategoryID, Valid: true},
		CreatedBy:  sql.NullInt64{Int64: UserID, Valid: true},
	}

	product, err := p.db.CreateProduct(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "invalid price format (need float number/decimal)",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	data := schemas.ProductData{
		ID:         product.ID,
		Name:       product.Name,
		Price:      price,
		Stock:      product.Stock,
		CategoryID: common.ConvertNullInt64(product.CategoryID),
		CreatedBy:  common.ConvertNullInt64(product.CreatedBy),
		CreatedAt:  common.ConvertNullTime(product.CreatedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "created successfully",
		"data":    data,
	})
}

// UpdateProduct godoc
// @Security BearerAuth
// @Summary Update an existing product
// @Description Update a product with the given ID and payload
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param payload body schemas.UpdateProduct true "Product Update Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products/{id} [put]
func (p *ProductController) UpdateProduct(ctx *gin.Context) {
	var payload *schemas.UpdateProduct
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid product id",
		})
		return
	}

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

	if _, err := p.db.GetCategoryByID(ctx, payload.CategoryID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "category id not found",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if _, err := p.db.GetProductByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve product with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if payload.Price < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid price (need positive number/decimal)",
		})
		return
	}

	args := &db.UpdateProductParams{
		ID:         id,
		Name:       payload.Name,
		Price:      strconv.FormatFloat(payload.Price, 'f', 2, 64),
		CategoryID: sql.NullInt64{Int64: payload.CategoryID, Valid: true},
		UpdatedBy:  sql.NullInt64{Int64: UserID, Valid: true},
	}

	product, err := p.db.UpdateProduct(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "invalid price format (need float number/decimal)",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	data := schemas.ProductData{
		ID:         product.ID,
		Name:       product.Name,
		Price:      price,
		Stock:      product.Stock,
		CategoryID: common.ConvertNullInt64(product.CategoryID),
		UpdatedBy:  common.ConvertNullInt64(product.UpdatedBy),
		UpdatedAt:  common.ConvertNullTime(product.UpdatedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "updated successfully",
		"data":    data,
	})
}

// GetProductById godoc
// @Security BearerAuth
// @Summary Get a product by ID
// @Description Retrieve a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products/{id} [get]
func (p *ProductController) GetProductById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid product id",
		})
		return
	}

	product, err := p.db.GetProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve product with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "invalid price format (need float number/decimal)",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	data := schemas.ProductData{
		ID:         product.ID,
		Name:       product.Name,
		Price:      price,
		Stock:      product.Stock,
		CategoryID: common.ConvertNullInt64(product.CategoryID),
		CreatedBy:  common.ConvertNullInt64(product.CreatedBy),
		CreatedAt:  common.ConvertNullTime(product.CreatedAt),
		UpdatedBy:  common.ConvertNullInt64(product.UpdatedBy),
		UpdatedAt:  common.ConvertNullTime(product.UpdatedAt),
		DeletedBy:  common.ConvertNullInt64(product.DeletedBy),
		DeletedAt:  common.ConvertNullTime(product.DeletedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllProducts godoc
// @Security BearerAuth
// @Summary Get all products
// @Description Retrieve all products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products [get]
func (p *ProductController) GetAllProducts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllProductsParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	products, err := p.db.GetAllProducts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if products == nil {
		products = []db.Product{}
	}

	data := make([]schemas.ProductData, len(products))
	for i, product := range products {
		price, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			continue
		}

		data[i] = schemas.ProductData{
			ID:         product.ID,
			Name:       product.Name,
			Price:      price,
			Stock:      product.Stock,
			CategoryID: common.ConvertNullInt64(product.CategoryID),
			CreatedBy:  common.ConvertNullInt64(product.CreatedBy),
			CreatedAt:  common.ConvertNullTime(product.CreatedAt),
			UpdatedBy:  common.ConvertNullInt64(product.UpdatedBy),
			UpdatedAt:  common.ConvertNullTime(product.UpdatedAt),
			DeletedBy:  common.ConvertNullInt64(product.DeletedBy),
			DeletedAt:  common.ConvertNullTime(product.DeletedAt),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllDeletedProducts godoc
// @Security BearerAuth
// @Summary Get all deleted products
// @Description Retrieve all deleted products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products/deleted [get]
func (p *ProductController) GetAllDeletedProducts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllDeletedProductsParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	products, err := p.db.GetAllDeletedProducts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if products == nil {
		products = []db.Product{}
	}

	data := make([]schemas.ProductData, len(products))
	for i, product := range products {
		price, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			continue
		}

		data[i] = schemas.ProductData{
			ID:         product.ID,
			Name:       product.Name,
			Price:      price,
			Stock:      product.Stock,
			CategoryID: common.ConvertNullInt64(product.CategoryID),
			CreatedBy:  common.ConvertNullInt64(product.CreatedBy),
			CreatedAt:  common.ConvertNullTime(product.CreatedAt),
			UpdatedBy:  common.ConvertNullInt64(product.UpdatedBy),
			UpdatedAt:  common.ConvertNullTime(product.UpdatedAt),
			DeletedBy:  common.ConvertNullInt64(product.DeletedBy),
			DeletedAt:  common.ConvertNullTime(product.DeletedAt),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// DeleteProductById godoc
// @Security BearerAuth
// @Summary Delete a product by ID
// @Description Delete a product with the given ID
// @Tags categories
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/categories/{id} [delete]
func (c *ProductController) DeleteProductById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid product id",
		})
		return
	}

	if _, err := c.db.GetProductByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve product with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	_, err = c.db.DeleteProductByID(ctx, id)
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

// SoftDeleteProductById godoc
// @Security BearerAuth
// @Summary Soft delete a product by ID
// @Description Soft delete a product with the given ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/products/{id}/soft [delete]
func (p *ProductController) SoftDeleteProductById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid category id",
		})
		return
	}

	if _, err := p.db.GetProductByID(ctx, id); err != nil {
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

	args := &db.SoftDeleteProductByIDParams{
		ID:        id,
		DeletedBy: sql.NullInt64{Int64: id, Valid: true},
	}
	_, err = p.db.SoftDeleteProductByID(ctx, *args)
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
