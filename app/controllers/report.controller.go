package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"pos-api/app/schemas"
	db "pos-api/db/sqlc"
	"pos-api/util/common"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	db  *db.Queries
	ctx context.Context
}

func NewReportController(db *db.Queries, ctx context.Context) *ReportController {
	return &ReportController{db, ctx}
}

// GetDetailOrderByID godoc
// @Security BearerAuth
// @Summary Get detailed order information by ID
// @Description Get order details including customer and items
// @Tags reports
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} schemas.Response
// @Failure 400,404,502 {object} schemas.Response
// @Router /api/v1/reports/orders/{id} [get]
func (c *ReportController) GetDetailOrderByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid order id",
		})
		return
	}

	order, err := c.db.GetOrderByID(ctx, id)
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

	// Get order items
	items, err := c.db.GetOrderItemsByOrderID(ctx, sql.NullInt64{Int64: id, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	var customer *schemas.CustomerResponse
	if order.CustomerID.Valid {
		customerData, err := c.db.GetCustomerByID(ctx, order.CustomerID.Int64)
		if err == nil {
			customer = &schemas.CustomerResponse{
				ID:         customerData.ID,
				MemberCode: customerData.MemberCode,
				Name:       customerData.Name,
				Phone:      common.ConvertNullString(customerData.Phone),
				Email:      common.ConvertNullString(customerData.Email),
			}
		}
	}

	orderItems := make([]schemas.OrderItemDetail, len(items))
	for i, item := range items {
		UnitPrice, err := strconv.ParseFloat(item.UnitPrice, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid price for product with id " + strconv.FormatInt(common.ConvertNullInt64(item.ProductID), 10),
				"error":   err.Error(),
			})
			return
		}

		orderItems[i] = schemas.OrderItemDetail{
			ID:         item.ID,
			ProductID:  common.ConvertNullInt64(item.ProductID),
			OldProduct: common.ConvertNullString(item.OldProduct),
			Quantity:   item.Quantity,
			UnitPrice:  UnitPrice,
		}
	}

	TotalAmount, err := strconv.ParseFloat(order.TotalAmount, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid price for order with id " + strconv.FormatInt(order.ID, 10),
			"error":   err.Error(),
		})
		return
	}

	response := schemas.OrderDetailResponse{
		ID:            order.ID,
		TrxNumber:     order.TrxNumber,
		CashierID:     common.ConvertNullInt64(order.CashierID),
		CustomerID:    order.CustomerID,
		TotalAmount:   TotalAmount,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
		OrderDate:     common.ConvertNullTime(order.OrderDate),
		UpdatedBy:     order.UpdatedBy,
		UpdatedAt:     order.UpdatedAt,
		Customer:      customer,
		OrderItems:    orderItems,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    response,
	})
}

// GetAllOrders godoc
// @Security BearerAuth
// @Summary Get all orders with optional filters
// @Description Get orders list with pagination and filters
// @Tags reports
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param customer_id query int false "Customer ID filter"
// @Param cashier_id query int false "Cashier ID filter"
// @Param status query string false "Order status filter"
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/reports/orders [get]
func (c *ReportController) GetAllOrders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	// customerID, _ := strconv.ParseInt(ctx.Query("customer_id"), 10, 64)
	// cashierID, _ := strconv.ParseInt(ctx.Query("cashier_id"), 10, 64)
	// status := ctx.Query("status")

	offset := (page - 1) * limit

	args := &db.GetAllOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	orders, err := c.db.GetAllOrders(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    orders,
	})
}

// FastMoving godoc
// @Security BearerAuth
// @Summary Get fast moving products
// @Description Get list of fast moving products for specific month and year
// @Tags reports
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} schemas.Response
// @Failure 400,502 {object} schemas.Response
// @Router /api/v1/reports/fast-moving [get]
func (c *ReportController) FastMoving(ctx *gin.Context) {
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil || month < 1 || month > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid month",
		})
		return
	}

	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid year",
		})
		return
	}

	StartDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	EndDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	args := &db.GetFastMovingProductsParams{
		OrderDate:   sql.NullTime{Time: StartDate, Valid: true},
		OrderDate_2: sql.NullTime{Time: EndDate, Valid: true},
	}

	products, err := c.db.GetFastMovingProducts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    products,
	})
}

// SlowMoving godoc
// @Security BearerAuth
// @Summary Get slow moving products
// @Description Get list of slow moving products for specific month and year
// @Tags reports
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} schemas.Response
// @Failure 400,502 {object} schemas.Response
// @Router /api/v1/reports/slow-moving [get]
func (c *ReportController) SlowMoving(ctx *gin.Context) {
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil || month < 1 || month > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid month",
		})
		return
	}

	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid year",
		})
		return
	}

	StartDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	EndDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	args := &db.GetSlowMovingProductsParams{
		OrderDate:   sql.NullTime{Time: StartDate, Valid: true},
		OrderDate_2: sql.NullTime{Time: EndDate, Valid: true},
	}

	products, err := c.db.GetSlowMovingProducts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    products,
	})
}

// TopCashier godoc
// @Security BearerAuth
// @Summary Get top performing cashiers
// @Description Get list of top performing cashiers for specific month and year
// @Tags reports
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} schemas.Response
// @Failure 400,502 {object} schemas.Response
// @Router /api/v1/reports/top-cashiers [get]
func (c *ReportController) TopCashier(ctx *gin.Context) {
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil || month < 1 || month > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid month",
		})
		return
	}

	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid year",
		})
		return
	}

	StartDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	EndDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	args := &db.GetTopCashiersParams{
		OrderDate:   sql.NullTime{Time: StartDate, Valid: true},
		OrderDate_2: sql.NullTime{Time: EndDate, Valid: true},
	}

	cashiers, err := c.db.GetTopCashiers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    cashiers,
	})
}

// TopCustomer godoc
// @Security BearerAuth
// @Summary Get top customers
// @Description Get list of top customers for specific month and year
// @Tags reports
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} schemas.Response
// @Failure 400,502 {object} schemas.Response
// @Router /api/v1/reports/top-customers [get]
func (c *ReportController) TopCustomer(ctx *gin.Context) {
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil || month < 1 || month > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid month",
		})
		return
	}

	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid year",
		})
		return
	}

	StartDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	EndDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	args := &db.GetTopCustomersParams{
		OrderDate:   sql.NullTime{Time: StartDate, Valid: true},
		OrderDate_2: sql.NullTime{Time: EndDate, Valid: true},
	}

	customers, err := c.db.GetTopCustomers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    customers,
	})
}
