package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type DealerRouteController struct {
	dealerController controllers.DealerController
}

func NewRouteDealerController(dealerController controllers.DealerController) DealerRouteController {
	return DealerRouteController{dealerController}
}

func (dc *DealerRouteController) DealerRoute(rg *gin.RouterGroup) {
	router := rg.Group("dealers")

	router.POST("/", middleware.DeserializeUser(), dc.dealerController.CreateDealer)
	router.GET("/", dc.dealerController.FindDealers)
	router.GET("/:dealerId", dc.dealerController.FindDealerById)
	router.PUT("/:dealerId", middleware.DeserializeUser(), dc.dealerController.UpdateDealer)
	router.DELETE("/:dealerId", middleware.DeserializeUser(), dc.dealerController.DeleteDealer)
}
