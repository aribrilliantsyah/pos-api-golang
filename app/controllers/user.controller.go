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

type UserController struct {
	db  *db.Queries
	ctx context.Context
}

func NewUserController(db *db.Queries, ctx context.Context) *UserController {
	return &UserController{db, ctx}
}

// CreateUser godoc
// @Security BearerAuth
// @Summary Create a new user
// @Description Create a new user with the given payload
// @Tags users
// @Accept json
// @Produce json
// @Param payload body schemas.CreateUser true "User Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users [post]
func (u *UserController) CreateUser(ctx *gin.Context) {
	var payload *schemas.CreateUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if payload.Role != "admin" && payload.Role != "cashier" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid role (admin or cashier)",
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

	hashedPassword, err := jwt.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to hash password",
		})
		return
	}

	CheckUsername, _ := u.db.GetUserByUsername(ctx, payload.Username)
	if CheckUsername.Username != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "username already exists",
		})
		return
	}

	args := &db.CreateUserParams{
		Username:     payload.Username,
		PasswordHash: string(hashedPassword),
		Role:         payload.Role,
		FullName:     payload.FullName,
		CreatedBy:    sql.NullInt64{Int64: UserID, Valid: true},
	}

	user, err := u.db.CreateUser(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.UserData{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedBy: common.ConvertNullInt64(user.CreatedBy),
		CreatedAt: common.ConvertNullTime(user.CreatedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "created successfully",
		"data":    data,
	})
}

// UpdateUser godoc
// @Security BearerAuth
// @Summary Update an existing user
// @Description Update a user with the given ID and payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param payload body schemas.UpdateUser true "User Update Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users/{id} [put]
func (u *UserController) UpdateUser(ctx *gin.Context) {
	var payload *schemas.UpdateUser
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
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

	if payload.Role != "admin" && payload.Role != "cashier" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid role (admin or cashier)",
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

	if _, err := u.db.GetUserByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve user with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	argsCheckUsername := &db.GetUserByUsernameExceptIDParams{
		Username: payload.Username,
		ID:       id,
	}
	CheckUsername, _ := u.db.GetUserByUsernameExceptID(ctx, *argsCheckUsername)
	if CheckUsername.Username != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "username already exists",
		})
		return
	}

	var errExc error
	var user db.User

	if payload.Password != "" {
		hashedPassword, err := jwt.HashPassword(payload.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "failed to hash password",
			})
			return
		}

		args := &db.UpdateUserWithPasswordParams{
			ID:           id,
			Username:     payload.Username,
			Role:         payload.Role,
			FullName:     payload.FullName,
			UpdatedBy:    sql.NullInt64{Int64: UserID, Valid: true},
			PasswordHash: string(hashedPassword),
		}

		user, errExc = u.db.UpdateUserWithPassword(ctx, *args)
	} else {
		args := &db.UpdateUserParams{
			ID:        id,
			Username:  payload.Username,
			Role:      payload.Role,
			FullName:  payload.FullName,
			UpdatedBy: sql.NullInt64{Int64: UserID, Valid: true},
		}

		user, errExc = u.db.UpdateUser(ctx, *args)
	}

	if errExc != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": errExc.Error(),
		})
		return
	}

	data := schemas.UserData{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		UpdatedBy: common.ConvertNullInt64(user.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(user.UpdatedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "updated successfully",
		"data":    data,
	})
}

// GetUserById godoc
// @Security BearerAuth
// @Summary Get a user by ID
// @Description Retrieve a user by its ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users/{id} [get]
func (u *UserController) GetUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	user, err := u.db.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve user with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	data := schemas.UserData{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedBy: common.ConvertNullInt64(user.CreatedBy),
		CreatedAt: common.ConvertNullTime(user.CreatedAt),
		UpdatedBy: common.ConvertNullInt64(user.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(user.UpdatedAt),
		DeletedBy: common.ConvertNullInt64(user.DeletedBy),
		DeletedAt: common.ConvertNullTime(user.DeletedAt),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllUsers godoc
// @Security BearerAuth
// @Summary Get all users
// @Description Retrieve all users with pagination
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllUsersParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	users, err := c.db.GetAllUsers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if users == nil {
		users = []db.User{}
	}

	data := make([]schemas.UserData, len(users))
	for i, user := range users {
		data[i] = schemas.UserData{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			CreatedBy: common.ConvertNullInt64(user.CreatedBy),
			CreatedAt: common.ConvertNullTime(user.CreatedAt),
			UpdatedBy: common.ConvertNullInt64(user.UpdatedBy),
			UpdatedAt: common.ConvertNullTime(user.UpdatedAt),
			DeletedBy: common.ConvertNullInt64(user.DeletedBy),
			DeletedAt: common.ConvertNullTime(user.DeletedAt),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// GetAllDeletedUsers godoc
// @Security BearerAuth
// @Summary Get all deleted users
// @Description Retrieve all users with pagination
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users/deleted [get]
func (c *UserController) GetAllDeletedUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.GetAllDeletedUsersParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	users, err := c.db.GetAllDeletedUsers(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if users == nil {
		users = []db.User{}
	}

	data := make([]schemas.UserData, len(users))
	for i, user := range users {
		data[i] = schemas.UserData{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			CreatedBy: common.ConvertNullInt64(user.CreatedBy),
			CreatedAt: common.ConvertNullTime(user.CreatedAt),
			UpdatedBy: common.ConvertNullInt64(user.UpdatedBy),
			UpdatedAt: common.ConvertNullTime(user.UpdatedAt),
			DeletedBy: common.ConvertNullInt64(user.DeletedBy),
			DeletedAt: common.ConvertNullTime(user.DeletedAt),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "retrieved successfully",
		"data":    data,
	})
}

// DeleteUserById godoc
// @Security BearerAuth
// @Summary Delete a user by ID
// @Description Delete a user with the given ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	if _, err := c.db.GetUserByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve user with this id",
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	_, err = c.db.DeleteUserByID(ctx, id)
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

// SoftDeleteUserById godoc
// @Security BearerAuth
// @Summary Soft delete a user by ID
// @Description Soft delete a user with the given ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 404 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/users/{id}/soft [delete]
func (c *UserController) SoftDeleteUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	if _, err := c.db.GetUserByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "failed to retrieve user with this id",
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

	args := &db.SoftDeleteUserByIDParams{
		ID:        id,
		DeletedBy: sql.NullInt64{Int64: UserID, Valid: true},
	}
	_, err = c.db.SoftDeleteUserByID(ctx, *args)
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
