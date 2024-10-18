package routes

import (
	"context"
	"database/sql"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(db *db.Queries, ctx context.Context, sqlDB *sql.DB, rg *gin.RouterGroup) {
	transactionHistoryController := *controllers.NewTransactionController(db, sqlDB, ctx)
	router := rg.Group("transaction")
	router.POST("/order", transactionHistoryController.CreateOrder)
	router.POST("/refund", transactionHistoryController.CreateRefund)
}
