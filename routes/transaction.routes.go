package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type TransactionRouteController struct {
	transactionController controllers.TransactionController
}

func NewRouteTransactionController(transactionController controllers.TransactionController) TransactionRouteController {
	return TransactionRouteController{transactionController}
}

func (tc *TransactionRouteController) TransactionRoute(rg *gin.RouterGroup) {
	router := rg.Group("transactions")

	router.POST("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), tc.transactionController.CreateTransaction)
	router.GET("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), tc.transactionController.FindTransactions)
	router.GET("/:transactionId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), tc.transactionController.FindTransactionById)
	router.PUT("/:transactionId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), tc.transactionController.UpdateTransaction)
	router.DELETE("/:transactionId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), tc.transactionController.DeleteTransaction)
}
