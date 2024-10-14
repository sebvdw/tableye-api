package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type PlayerRouteController struct {
	playerController controllers.PlayerController
}

func NewRoutePlayerController(playerController controllers.PlayerController) PlayerRouteController {
	return PlayerRouteController{playerController}
}

func (pc *PlayerRouteController) PlayerRoute(rg *gin.RouterGroup) {
	router := rg.Group("players")

	router.POST("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), pc.playerController.CreatePlayer)
	router.GET("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), pc.playerController.FindPlayers)
	router.GET("/:playerId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), pc.playerController.FindPlayerById)
	router.PUT("/:playerId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), pc.playerController.UpdatePlayer)
	router.DELETE("/:playerId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), pc.playerController.DeletePlayer)
	router.GET("/:playerId/stats", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), pc.playerController.FindPlayerStats)
}
