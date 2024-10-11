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

	router.POST("/", middleware.DeserializeUser(), gc.gameController.CreateGame)
	router.GET("/", gc.gameController.FindGames)
	router.GET("/:gameId", gc.gameController.FindGameById)
	router.PUT("/:gameId", middleware.DeserializeUser(), gc.gameController.UpdateGame)
	router.DELETE("/:gameId", middleware.DeserializeUser(), gc.gameController.DeleteGame)
}
