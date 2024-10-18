package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupProductHistoryRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	productHistoryController := *controllers.NewProductHistoryController(db, ctx)
	router := rg.Group("product-history")
	router.GET("/", productHistoryController.GetAllProductHistory)
	router.POST("/", productHistoryController.CreateProductHistory)
}
