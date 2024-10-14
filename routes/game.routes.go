package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type GameRouteController struct {
	gameController controllers.GameController
}

func NewRouteGameController(gameController controllers.GameController) GameRouteController {
	return GameRouteController{gameController}
}

func (gc *GameRouteController) GameRoute(rg *gin.RouterGroup) {
	router := rg.Group("games")

	router.POST("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gc.gameController.CreateGame)
	router.GET("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), gc.gameController.FindGames)
	router.GET("/:gameId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), gc.gameController.FindGameById)
	router.PUT("/:gameId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gc.gameController.UpdateGame)
	router.DELETE("/:gameId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gc.gameController.DeleteGame)
}
