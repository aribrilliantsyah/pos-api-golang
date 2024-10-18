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

type TransactionController struct {
	db    *db.Queries
	sqlDB *sql.DB
	ctx   context.Context
}

func NewTransactionController(db *db.Queries, sqlDB *sql.DB, ctx context.Context) *TransactionController {
	return &TransactionController{db, sqlDB, ctx}
}

// CreateOrder godoc
// @Security BearerAuth
// @Summary Create product stock history
// @Description Create a new order
// @Tags transaction
// @Accept json
// @Produce json
// @Param payload body schemas.CreateOrder true "Order Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/transaction/order [post]
func (p *TransactionController) CreateOrder(ctx *gin.Context) {
	var payload schemas.CreateOrder

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request data",
			"error":   err.Error(),
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
	tx, err := p.sqlDB.Begin()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	defer tx.Rollback()
	qtx := p.db.WithTx(tx)

	//Customer
	var CustomerID int64
	if payload.Type == "guest" {

	} else if payload.Type == "member" {
		CustomerID = payload.CustomerID
	} else if payload.Type == "new" {
		MemberCode, err := qtx.GenerateMemberCode(ctx)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "failed",
					"message": "failed to retrieve member code",
				})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		args := &db.CreateCustomerParams{
			MemberCode: common.ConvertNullString(MemberCode),
			Name:       payload.Customer.Name,
			Email:      sql.NullString{String: payload.Customer.Email, Valid: true},
			Phone:      sql.NullString{String: payload.Customer.Phone, Valid: true},
			CreatedBy:  sql.NullInt64{Int64: 1, Valid: true},
		}

		Customer, err := qtx.CreateCustomer(ctx, *args)
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
		CustomerID = Customer.ID
	}

	//Order
	TrxNumber := fmt.Sprintf("TRX-%d-%s-%d", CustomerID, time.Now().Format("20060102150405"), UserID)
	TotalAmount := 0.0

	Items := make([]db.CreateOrderItemParams, 0)

	for _, item := range payload.Items {
		Product, err := qtx.GetProductByID(ctx, item.ProductID)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "failed",
					"message": "failed to retrieve product with id " + strconv.FormatInt(item.ProductID, 10),
				})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		if Product.Stock < item.Quantity {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "stock not enough product with id " + strconv.FormatInt(item.ProductID, 10),
			})
			return
		}

		Items = append(Items, db.CreateOrderItemParams{
			ProductID: sql.NullInt64{Int64: Product.ID, Valid: true},
			Quantity:  item.Quantity,
			UnitPrice: Product.Price,
			CreatedBy: sql.NullInt64{Int64: 1, Valid: true},
		})

		Price, err := strconv.ParseFloat(Product.Price, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid price for product with id " + strconv.FormatInt(item.ProductID, 10),
				"error":   err.Error(),
			})
			return
		}
		TotalAmount += Price * float64(item.Quantity)
	}

	args := &db.CreateOrderParams{
		TrxNumber:     TrxNumber,
		CashierID:     sql.NullInt64{Int64: UserID, Valid: true},
		CustomerID:    sql.NullInt64{Int64: int64(CustomerID), Valid: CustomerID != 0},
		TotalAmount:   strconv.FormatFloat(TotalAmount, 'f', 2, 64),
		PaymentMethod: payload.PaymentMethod,
		Status:        "order",
	}

	Order, err := qtx.CreateOrder(ctx, *args)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	for _, item := range Items {

		Product, err := qtx.GetProductByID(ctx, item.ProductID.Int64)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "failed",
					"message": "failed to retrieve product with id " + strconv.FormatInt(item.ProductID.Int64, 10),
				})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		//order item
		args := &db.CreateOrderItemParams{
			OrderID:    sql.NullInt64{Int64: Order.ID, Valid: true},
			ProductID:  item.ProductID,
			OldProduct: sql.NullString{String: Product.Name, Valid: true},
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			CreatedBy:  sql.NullInt64{Int64: 1, Valid: true},
		}

		if _, err := qtx.CreateOrderItem(ctx, *args); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		// Update product stock
		newStock := Product.Stock
		if payload.Type == "in" {
			newStock += item.Quantity
		} else {
			newStock -= item.Quantity
		}

		updateArgs := db.UpdateProductStockParams{
			ID:        Product.ID,
			Stock:     newStock,
			UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
		}

		if _, err = p.db.UpdateProductStock(ctx, updateArgs); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": "failed to update product stock with id " + strconv.FormatInt(item.ProductID.Int64, 10),
				"error":   err.Error(),
			})
			return
		}

		//product history
		historyArgs := &db.CreateProductHistoryParams{
			TrxRef:         Order.TrxNumber,
			ProductID:      sql.NullInt64{Int64: item.ProductID.Int64, Valid: true},
			QuantityChange: item.Quantity,
			Type:           "out",
			Reason:         sql.NullString{String: "Order", Valid: true},
			CreatedBy:      sql.NullInt64{Int64: UserID, Valid: true},
		}

		_, err = qtx.CreateProductHistory(ctx, *historyArgs)
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	}

	tx.Commit()

	data := schemas.OrderData{
		ID:            Order.ID,
		TrxNumber:     Order.TrxNumber,
		CashierID:     common.ConvertNullInt64(Order.CashierID),
		CustomerID:    common.ConvertNullInt64(Order.CustomerID),
		TotalAmount:   Order.TotalAmount,
		PaymentMethod: Order.PaymentMethod,
		Status:        Order.Status,
		OrderDate:     common.ConvertNullTime(Order.OrderDate),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "order created successfully",
		"data":    data,
	})
}

// CreateRefund godoc
// @Security BearerAuth
// @Summary Create product stock history
// @Description Create a new order
// @Tags transaction
// @Accept json
// @Produce json
// @Param payload body schemas.CreateRefund true "Refund Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/transaction/refund [post]
func (p *TransactionController) CreateRefund(ctx *gin.Context) {
	var payload schemas.CreateRefund

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request data",
			"error":   err.Error(),
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
	tx, err := p.sqlDB.Begin()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	defer tx.Rollback()
	qtx := p.db.WithTx(tx)

	// Get Order by TrxNumber
	order, err := qtx.GetOrderByTrxNumber(ctx, payload.TrxNumber)
	Order := order

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "order not found",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Check if order is already refunded
	if order.Status == "refunded" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "order already refunded",
		})
		return
	}

	// Get Order Items
	orderItems, err := qtx.GetOrderItemsByOrderID(ctx, sql.NullInt64{Int64: order.ID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Process each item for refund
	for _, item := range orderItems {
		// Get current product
		product, err := qtx.GetProductByID(ctx, item.ProductID.Int64)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "failed",
					"message": "product not found with id " + strconv.FormatInt(item.ProductID.Int64, 10),
				})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		// Update product stock (add back the quantity)
		newStock := product.Stock + item.Quantity
		updateArgs := db.UpdateProductStockParams{
			ID:        product.ID,
			Stock:     newStock,
			UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
		}

		if _, err = qtx.UpdateProductStock(ctx, updateArgs); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": "failed to update product stock",
				"error":   err.Error(),
			})
			return
		}

		// Create product history for refund
		historyArgs := &db.CreateProductHistoryParams{
			TrxRef:         order.TrxNumber + "-REFUND",
			ProductID:      sql.NullInt64{Int64: item.ProductID.Int64, Valid: true},
			QuantityChange: item.Quantity,
			Type:           "in",
			Reason:         sql.NullString{String: "Refund", Valid: true},
			CreatedBy:      sql.NullInt64{Int64: UserID, Valid: true},
		}

		_, err = qtx.CreateProductHistory(ctx, *historyArgs)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	}

	//create refund
	refundArgs := &db.CreateRefundParams{
		OrderID:   sql.NullInt64{Int64: order.ID, Valid: true},
		Reason:    payload.Reason,
		CreatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	_, err = qtx.CreateRefund(ctx, *refundArgs)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": "failed to create refund",
			"error":   err.Error(),
		})
		return
	}

	// Update order status to refunded
	updateOrderArgs := &db.UpdateOrderStatusParams{
		ID:        order.ID,
		Status:    "refunded",
		UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	_, err = qtx.UpdateOrderStatus(ctx, *updateOrderArgs)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": "failed to update order status",
			"error":   err.Error(),
		})
		return
	}

	tx.Commit()

	data := schemas.OrderData{
		ID:            Order.ID,
		TrxNumber:     Order.TrxNumber,
		CashierID:     common.ConvertNullInt64(Order.CashierID),
		CustomerID:    common.ConvertNullInt64(Order.CustomerID),
		TotalAmount:   Order.TotalAmount,
		PaymentMethod: Order.PaymentMethod,
		Status:        Order.Status,
		OrderDate:     common.ConvertNullTime(Order.OrderDate),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "order refunded successfully",
		"data":    data,
	})
}
