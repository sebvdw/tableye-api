package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type GameSummaryRouteController struct {
	gameSummaryController controllers.GameSummaryController
}

func NewRouteGameSummaryController(gameSummaryController controllers.GameSummaryController) GameSummaryRouteController {
	return GameSummaryRouteController{gameSummaryController}
}

func (gsc *GameSummaryRouteController) GameSummaryRoute(rg *gin.RouterGroup) {
	router := rg.Group("game-summaries")

	router.POST("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), gsc.gameSummaryController.CreateGameSummary)
	router.GET("/", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gsc.gameSummaryController.FindGameSummaries)
	router.GET("/:gameSummaryId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gsc.gameSummaryController.FindGameSummaryById)
	router.PUT("/:gameSummaryId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin", "dealer"), gsc.gameSummaryController.UpdateGameSummary)
	router.DELETE("/:gameSummaryId", middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"), gsc.gameSummaryController.DeleteGameSummary)
}
