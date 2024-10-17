package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	categoryController := *controllers.NewCategoryController(db, ctx)
	router := rg.Group("categories")
	router.POST("/", categoryController.CreateCategory)
	router.GET("/", categoryController.GetAllCategories)
	router.GET("/deleted", categoryController.GetAllDeletedCategories)
	router.PUT("/:id", categoryController.UpdateCategory)
	router.GET("/:id", categoryController.GetCategoryById)
	router.DELETE("/:id", categoryController.DeleteCategoryById)
	router.DELETE("/:id/soft", categoryController.SoftDeleteCategoryById)
}
