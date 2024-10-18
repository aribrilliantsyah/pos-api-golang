package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pos-api/app/schemas"
	db "pos-api/db/sqlc"
	"pos-api/util/common"
	"pos-api/util/jwt"

	"github.com/gin-gonic/gin"
)

type ProductHistoryController struct {
	db  *db.Queries
	ctx context.Context
}

func NewProductHistoryController(db *db.Queries, ctx context.Context) *ProductHistoryController {
	return &ProductHistoryController{db, ctx}
}

// CreateProductHistory godoc
// @Security BearerAuth
// @Summary Create product stock history
// @Description Create a new product stock history and update product stock
// @Tags product-history
// @Accept json
// @Produce json
// @Param payload body schemas.CreateProductHistory true "Product History Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/product-history [post]
func (p *ProductHistoryController) CreateProductHistory(ctx *gin.Context) {
	var payload schemas.CreateProductHistory

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Get product first to check stock
	product, err := p.db.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "product not found",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if common.ConvertNullInt64(product.DeletedBy) != 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "product deleted (soft)",
		})
		return
	}

	// Validate stock for "out" type
	if payload.Type == "out" && product.Stock < payload.QuantityChange {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "insufficient stock",
		})
		return
	}

	userInfo, err := jwt.GetUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "unauthorized user",
			"error":   err.Error(),
		})
		return
	}
	UserID := userInfo.UserID

	// Generate transaction reference
	trxRef := fmt.Sprintf("TRX-%s-%d", time.Now().Format("20060102150405"), UserID)

	// Create history record
	historyArgs := db.CreateProductHistoryParams{
		TrxRef:         trxRef,
		ProductID:      sql.NullInt64{Int64: payload.ProductID, Valid: true},
		QuantityChange: payload.QuantityChange,
		Type:           payload.Type,
		Reason:         sql.NullString{String: payload.Reason, Valid: true},
		CreatedBy:      sql.NullInt64{Int64: UserID, Valid: true},
	}

	history, err := p.db.CreateProductHistory(ctx, historyArgs)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": "failed to create product history",
			"error":   err.Error(),
		})
		return
	}

	// Update product stock
	newStock := product.Stock
	if payload.Type == "in" {
		newStock += payload.QuantityChange
	} else {
		newStock -= payload.QuantityChange
	}

	updateArgs := db.UpdateProductStockParams{
		ID:        payload.ProductID,
		Stock:     newStock,
		UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	if _, err = p.db.UpdateProductStock(ctx, updateArgs); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": "failed to update product stock",
			"error":   err.Error(),
		})
		return
	}

	data := schemas.ProductHistoryData{
		ID:             history.ID,
		TrxRef:         history.TrxRef,
		ProductID:      common.ConvertNullInt64(history.ProductID),
		QuantityChange: history.QuantityChange,
		Type:           history.Type,
		Reason:         common.ConvertNullString(history.Reason),
		CreatedBy:      common.ConvertNullInt64(history.CreatedBy),
		CreatedAt:      common.ConvertNullTime(history.CreatedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "stock updated successfully",
		"data":    data,
	})
}

// GetAllProductHistory godoc
// @Security BearerAuth
// @Summary Get all product history
// @Description Retrieve all product history with pagination
// @Tags product-history
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/product-history [get]
func (p *ProductHistoryController) GetAllProductHistory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	args := db.GetAllProductHistoryParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	histories, err := p.db.GetAllProductHistory(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": "failed to retrieve product history",
			"error":   err.Error(),
		})
		return
	}

	data := make([]schemas.ProductHistoryData, len(histories))
	for i, history := range histories {
		data[i] = schemas.ProductHistoryData{
			ID:             history.ID,
			TrxRef:         history.TrxRef,
			ProductID:      common.ConvertNullInt64(history.ProductID),
			QuantityChange: history.QuantityChange,
			Type:           history.Type,
			Reason:         common.ConvertNullString(history.Reason),
			CreatedBy:      common.ConvertNullInt64(history.CreatedBy),
			CreatedAt:      common.ConvertNullTime(history.CreatedAt),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}
