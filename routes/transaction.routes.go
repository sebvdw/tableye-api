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

	router.POST("/", middleware.DeserializeUser(), tc.transactionController.CreateTransaction)
	router.GET("/", middleware.DeserializeUser(), tc.transactionController.FindTransactions)
	router.GET("/:transactionId", middleware.DeserializeUser(), tc.transactionController.FindTransactionById)
	router.PUT("/:transactionId", middleware.DeserializeUser(), tc.transactionController.UpdateTransaction)
	router.DELETE("/:transactionId", middleware.DeserializeUser(), tc.transactionController.DeleteTransaction)
	router.GET("/game-summary/:gameSummaryId", middleware.DeserializeUser(), tc.transactionController.GetTransactionsByGameSummary)
}
