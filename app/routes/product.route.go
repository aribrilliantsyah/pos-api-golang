package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	productController := *controllers.NewProductController(db, ctx)
	router := rg.Group("products")
	router.POST("/", productController.CreateProduct)
	router.GET("/", productController.GetAllProducts)
	router.GET("/deleted", productController.GetAllDeletedProducts)
	router.PUT("/:id", productController.UpdateProduct)
	router.GET("/:id", productController.GetProductById)
	router.DELETE("/:id", productController.DeleteProductById)
	router.DELETE("/:id/soft", productController.SoftDeleteProductById)
}
