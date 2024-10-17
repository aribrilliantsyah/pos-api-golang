package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	authController := *controllers.NewAuthController(db, ctx)
	router := rg.Group("auth")
	router.POST("/login", authController.Login)
	// router.POST("/register", authController.Register)
}
