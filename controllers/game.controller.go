package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type GameController struct {
	DB *gorm.DB
}

func NewGameController(DB *gorm.DB) GameController {
	return GameController{DB}
}

func (gc *GameController) CreateGame(ctx *gin.Context) {
	var payload *models.CreateGameRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newGame := models.Game{
		Name:        payload.Name,
		Type:        payload.Type,
		Description: payload.Description,
		MaxPlayers:  payload.MaxPlayers,
		MinPlayers:  payload.MinPlayers,
		MinBet:      payload.MinBet,
		MaxBet:      payload.MaxBet,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := gc.DB.Create(&newGame)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Game with that name already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newGame})
}

func (gc *GameController) UpdateGame(ctx *gin.Context) {
	gameId := ctx.Param("gameId")
	var payload *models.UpdateGameRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var game models.Game
	result := gc.DB.First(&game, "id = ?", gameId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game with that ID exists"})
		return
	}

	now := time.Now()
	gameToUpdate := models.Game{
		Name:        payload.Name,
		Type:        payload.Type,
		Description: payload.Description,
		MaxPlayers:  payload.MaxPlayers,
		MinPlayers:  payload.MinPlayers,
		MinBet:      payload.MinBet,
		MaxBet:      payload.MaxBet,
		UpdatedAt:   now,
	}

	gc.DB.Model(&game).Updates(gameToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": game})
}

func (gc *GameController) FindGameById(ctx *gin.Context) {
	gameId := ctx.Param("gameId")

	var game models.Game
	result := gc.DB.Preload("Casinos").Preload("GameSummaries").First(&game, "id = ?", gameId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": game})
}

func (gc *GameController) FindGames(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var games []models.Game
	results := gc.DB.Limit(intLimit).Offset(offset).Find(&games)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(games), "data": games})
}

func (gc *GameController) DeleteGame(ctx *gin.Context) {
	gameId := ctx.Param("gameId")

	var game models.Game
	result := gc.DB.First(&game, "id = ?", gameId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game with that ID exists"})
		return
	}

	gc.DB.Delete(&game)
	ctx.JSON(http.StatusNoContent, nil)
}
