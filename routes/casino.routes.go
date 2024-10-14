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

	router.POST("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), cc.casinoController.CreateCasino)
	router.GET("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), cc.casinoController.FindCasinos)
	router.GET("/:casinoId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), cc.casinoController.FindCasinoById)
	router.PUT("/:casinoId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), cc.casinoController.UpdateCasino)
	router.DELETE("/:casinoId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), cc.casinoController.DeleteCasino)
}
