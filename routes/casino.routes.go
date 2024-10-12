package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type CasinoRouteController struct {
	casinoController controllers.CasinoController
}

func NewRouteCasinoController(casinoController controllers.CasinoController) CasinoRouteController {
	return CasinoRouteController{casinoController}
}

func (cc *CasinoRouteController) CasinoRoute(rg *gin.RouterGroup) {
	router := rg.Group("casinos")

	router.POST("/", middleware.DeserializeUser(), cc.casinoController.CreateCasino)
	router.GET("/", middleware.DeserializeUser(), cc.casinoController.FindCasinos)
	router.GET("/:casinoId", middleware.DeserializeUser(), cc.casinoController.FindCasinoById)
	router.PUT("/:casinoId", middleware.DeserializeUser(), cc.casinoController.UpdateCasino)
	router.DELETE("/:casinoId", middleware.DeserializeUser(), cc.casinoController.DeleteCasino)
}
