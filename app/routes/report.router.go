package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	reportController := *controllers.NewReportController(db, ctx)

	router := rg.Group("reports")
	router.GET("/orders/:id", reportController.GetDetailOrderByID)
	router.GET("/orders", reportController.GetAllOrders)
	router.GET("/fast-moving", reportController.FastMoving)
	router.GET("/slow-moving", reportController.SlowMoving)
	router.GET("/top-cashiers", reportController.TopCashier)
	router.GET("/top-customers", reportController.TopCustomer)
}
