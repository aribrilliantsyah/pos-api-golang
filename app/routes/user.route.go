package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	userController := *controllers.NewUserController(db, ctx)
	router := rg.Group("users")
	router.POST("/", userController.CreateUser)
	router.GET("/", userController.GetAllUsers)
	router.GET("/deleted", userController.GetAllDeletedUsers)
	router.PUT("/:id", userController.UpdateUser)
	router.GET("/:id", userController.GetUserById)
	router.DELETE("/:id", userController.DeleteUserById)
	router.DELETE("/:id/soft", userController.SoftDeleteUserById)
}
