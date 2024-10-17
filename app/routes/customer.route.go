package routes

import (
	"context"
	"pos-api/app/controllers"
	db "pos-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(db *db.Queries, ctx context.Context, rg *gin.RouterGroup) {
	customerController := *controllers.NewCustomerController(db, ctx)
	router := rg.Group("customers")
	router.POST("/", customerController.CreateCustomer)
	router.GET("/", customerController.GetAllCustomers)
	router.GET("/deleted", customerController.GetAllDeletedCustomers)
	router.PUT("/:id", customerController.UpdateCustomer)
	router.GET("/:id", customerController.GetCustomerById)
	router.DELETE("/:id", customerController.DeleteCustomerById)
	router.DELETE("/:id/soft", customerController.SoftDeleteCustomerById)
}
