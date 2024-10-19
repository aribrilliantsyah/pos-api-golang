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

type CustomerController struct {
	db  *db.Queries
	ctx context.Context
}

func NewCustomerController(db *db.Queries, ctx context.Context) *CustomerController {
	return &CustomerController{db, ctx}
}

// CreateCustomer godoc
// @Security BearerAuth
// @Summary Create a new customer
// @Description Create a new customer with the given payload
// @Tags customers
// @Accept json
// @Produce json
// @Param payload body schemas.CreateCustomer true "Customer Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers [post]
func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var payload *schemas.CreateCustomer

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

	CheckEmail, _ := c.db.GetCustomerByEmail(ctx, sql.NullString{String: payload.Email, Valid: payload.Email != ""})
	if common.ConvertNullString(CheckEmail.Email) != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "email already exists",
		})
		return
	}

	CheckPhone, _ := c.db.GetCustomerByPhone(ctx, sql.NullString{String: payload.Phone, Valid: true})
	if common.ConvertNullString(CheckPhone.Phone) != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "email already exists",
		})
		return
	}

	MemberCode, err := c.db.GenerateMemberCode(ctx)
	if err != nil {
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
		Name:       payload.Name,
		Phone:      sql.NullString{String: payload.Phone, Valid: payload.Phone != ""},
		Email:      sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		CreatedBy:  sql.NullInt64{Int64: UserID, Valid: true},
	}

	customer, err := c.db.CreateCustomer(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CustomerData{
		ID:         customer.ID,
		MemberCode: customer.MemberCode,
		Name:       customer.Name,
		Phone:      common.ConvertNullString(customer.Phone),
		Email:      common.ConvertNullString(customer.Email),
		CreatedBy:  common.ConvertNullInt64(customer.CreatedBy),
		CreatedAt:  common.ConvertNullTime(customer.CreatedAt),
		UpdatedBy:  common.ConvertNullInt64(customer.UpdatedBy),
		UpdatedAt:  common.ConvertNullTime(customer.UpdatedAt),
		DeletedAt:  common.ConvertNullTime(customer.DeletedAt),
		DeletedBy:  common.ConvertNullInt64(customer.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "created successfully",
		"data":    data,
	})
}

// UpdateCustomer godoc
// @Security BearerAuth
// @Summary Update an existing customer
// @Description Update a customer with the given ID and payload
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param payload body schemas.UpdateCustomer true "Customer Update Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers/{id} [put]
func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	var payload *schemas.UpdateCustomer
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid customer id",
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

	if _, err := c.db.GetCustomerByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve customer with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	argsCheckEmail := &db.GetCustomerByEmailExceptIDParams{
		Email: sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		ID:    id,
	}
	CheckEmail, _ := c.db.GetCustomerByEmailExceptID(ctx, *argsCheckEmail)
	if common.ConvertNullString(CheckEmail.Email) != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "email already exists",
		})
		return
	}

	argsCheckPhone := &db.GetCustomerByPhoneExceptIDParams{
		Phone: sql.NullString{String: payload.Phone, Valid: payload.Phone != ""},
		ID:    id,
	}
	CheckPhone, _ := c.db.GetCustomerByPhoneExceptID(ctx, *argsCheckPhone)
	if common.ConvertNullString(CheckPhone.Phone) != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "phone already exists",
		})
		return
	}

	args := &db.UpdateCustomerParams{
		ID:        id,
		Name:      payload.Name,
		Phone:     sql.NullString{String: payload.Phone, Valid: payload.Phone != ""},
		Email:     sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}

	customer, err := c.db.UpdateCustomer(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve customer with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CustomerData{
		ID:         customer.ID,
		MemberCode: customer.MemberCode,
		Name:       customer.Name,
		Phone:      common.ConvertNullString(customer.Phone),
		Email:      common.ConvertNullString(customer.Email),
		CreatedBy:  common.ConvertNullInt64(customer.CreatedBy),
		CreatedAt:  common.ConvertNullTime(customer.CreatedAt),
		UpdatedBy:  common.ConvertNullInt64(customer.UpdatedBy),
		UpdatedAt:  common.ConvertNullTime(customer.UpdatedAt),
		DeletedAt:  common.ConvertNullTime(customer.DeletedAt),
		DeletedBy:  common.ConvertNullInt64(customer.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "updated successfully",
		"data":    data,
	})
}

// GetCustomerById godoc
// @Security BearerAuth
// @Summary Get a customer by ID
// @Description Retrieve a customer by its ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers/{id} [get]
func (c *CustomerController) GetCustomerById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid customer id",
		})
		return
	}

	customer, err := c.db.GetCustomerByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve customer with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.CustomerData{
		ID:         customer.ID,
		MemberCode: customer.MemberCode,
		Name:       customer.Name,
		Phone:      common.ConvertNullString(customer.Phone),
		Email:      common.ConvertNullString(customer.Email),
		CreatedBy:  common.ConvertNullInt64(customer.CreatedBy),
		CreatedAt:  common.ConvertNullTime(customer.CreatedAt),
		UpdatedBy:  common.ConvertNullInt64(customer.UpdatedBy),
		UpdatedAt:  common.ConvertNullTime(customer.UpdatedAt),
		DeletedAt:  common.ConvertNullTime(customer.DeletedAt),
		DeletedBy:  common.ConvertNullInt64(customer.DeletedBy),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "customer retrieved successfully",
		"data":    data,
	})
}

// GetAllCustomers godoc
// @Security BearerAuth
// @Summary Get all customers
// @Description Retrieve all customers with pagination
// @Tags customers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers [get]
func (c *CustomerController) GetAllCustomers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllCustomersParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	customers, err := c.db.GetAllCustomers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if customers == nil {
		customers = []db.Customer{}
	}

	data := make([]schemas.CustomerData, len(customers))
	for i, customer := range customers {
		data[i] = schemas.CustomerData{
			ID:         customer.ID,
			MemberCode: customer.MemberCode,
			Name:       customer.Name,
			Phone:      common.ConvertNullString(customer.Phone),
			Email:      common.ConvertNullString(customer.Email),
			CreatedBy:  common.ConvertNullInt64(customer.CreatedBy),
			CreatedAt:  common.ConvertNullTime(customer.CreatedAt),
			UpdatedBy:  common.ConvertNullInt64(customer.UpdatedBy),
			UpdatedAt:  common.ConvertNullTime(customer.UpdatedAt),
			DeletedAt:  common.ConvertNullTime(customer.DeletedAt),
			DeletedBy:  common.ConvertNullInt64(customer.DeletedBy),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllCustomers godoc
// @Security BearerAuth
// @Summary Get all deleted customers
// @Description Retrieve all customers with pagination
// @Tags customers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers/deleted [get]
func (c *CustomerController) GetAllDeletedCustomers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllDeletedCustomersParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	customers, err := c.db.GetAllDeletedCustomers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if customers == nil {
		customers = []db.Customer{}
	}

	data := make([]schemas.CustomerData, len(customers))
	for i, customer := range customers {
		data[i] = schemas.CustomerData{
			ID:         customer.ID,
			MemberCode: customer.MemberCode,
			Name:       customer.Name,
			Phone:      common.ConvertNullString(customer.Phone),
			Email:      common.ConvertNullString(customer.Email),
			CreatedBy:  common.ConvertNullInt64(customer.CreatedBy),
			CreatedAt:  common.ConvertNullTime(customer.CreatedAt),
			UpdatedBy:  common.ConvertNullInt64(customer.UpdatedBy),
			UpdatedAt:  common.ConvertNullTime(customer.UpdatedAt),
			DeletedAt:  common.ConvertNullTime(customer.DeletedAt),
			DeletedBy:  common.ConvertNullInt64(customer.DeletedBy),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// DeleteCustomerById godoc
// @Security BearerAuth
// @Summary Delete a customer by ID
// @Description Delete a customer with the given ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers/{id} [delete]
func (c *CustomerController) DeleteCustomerById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid customer id",
		})
		return
	}

	if _, err := c.db.GetCustomerByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve customer with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	_, err = c.db.DeleteCustomerByID(ctx, id)
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

// SoftDeleteCustomerById godoc
// @Security BearerAuth
// @Summary Soft delete a customer by ID
// @Description Soft delete a customer with the given ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/customers/{id}/soft [delete]
func (c *CustomerController) SoftDeleteCustomerById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid customer id",
		})
		return
	}

	if _, err := c.db.GetCustomerByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve customer with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
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

	args := &db.SoftDeleteCustomerByIDParams{
		ID:        id,
		DeletedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}
	_, err = c.db.SoftDeleteCustomerByID(ctx, *args)
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
