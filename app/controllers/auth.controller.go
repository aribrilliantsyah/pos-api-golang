package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"pos-api/app/schemas"
	db "pos-api/db/sqlc"
	"pos-api/util/common"
	"pos-api/util/jwt"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	db  *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db, ctx}
}

// Login godoc
// @Summary Login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body schemas.Login true "Login Data"
// @Success 200 {object} schemas.Response
// @Failure 400 {object} schemas.Response
// @Failure 502 {object} schemas.Response
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var payload *schemas.Login

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	user, err := c.db.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "invalid credentials (1)",
		})
		return
	}

	if !jwt.CheckPasswordHash(payload.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "invalid credentials (2)",
		})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"failed":  "failed",
			"message": "failed to generate token",
		})
		return
	}

	Username := user.Username
	args := &db.SetCurrentTokenParams{
		Username: Username,
		CurrentToken: sql.NullString{
			String: string(token),
			Valid:  true,
		},
	}

	UserDetail, _ := c.db.GetUserByID(ctx, user.ID)
	if common.ConvertNullInt64(UserDetail.DeletedBy) != 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "user deleted (soft)",
		})
		return
	}

	_, err = c.db.SetCurrentToken(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to update token",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "successfully logged in",
		"data": gin.H{
			"token": token,
		},
	})
}

// ?Register godoc
// ?@Summary Register a new user
// ?@Description Register a new user with username and password
// ?@Tags auth
// ?@Accept json
// ?@Produce json
// ?@Param payload body schemas.Register true "Register Data"
// ?@Success 201 {object} schemas.Response
// ?@Failure 400 {object} schemas.Response
// ?@Failure 502 {object} schemas.Response
// ?@Router /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var payload *schemas.Register

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid input data",
		})
		return
	}

	hashedPassword, err := jwt.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to hash password",
		})
		return
	}

	CheckUsername, err := c.db.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

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
	}

	user, err := c.db.CreateUser(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to create user",
			"err":     err.Error(),
		})
		return
	}

	data := schemas.UserData{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		FullName:  user.FullName,
		CreatedBy: common.ConvertNullInt64(user.CreatedBy),
		CreatedAt: common.ConvertNullTime(user.CreatedAt),
		UpdatedBy: common.ConvertNullInt64(user.UpdatedBy),
		UpdatedAt: common.ConvertNullTime(user.UpdatedAt),
		DeletedBy: common.ConvertNullInt64(user.DeletedBy),
		DeletedAt: common.ConvertNullTime(user.DeletedAt),
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   data,
	})
}
