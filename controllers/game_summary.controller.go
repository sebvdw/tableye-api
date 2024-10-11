package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type GameSummaryController struct {
	DB *gorm.DB
}

func NewGameSummaryController(DB *gorm.DB) GameSummaryController {
	return GameSummaryController{DB}
}

func (gsc *GameSummaryController) CreateGameSummary(ctx *gin.Context) {
	var payload *models.CreateGameSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	gameID, err := uuid.Parse(payload.GameID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid game ID"})
		return
	}

	casinoID, err := uuid.Parse(payload.CasinoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid casino ID"})
		return
	}

	dealerID, err := uuid.Parse(payload.DealerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid dealer ID"})
		return
	}

	playerIDs := make([]uuid.UUID, len(payload.PlayerIDs))
	for i, playerIDStr := range payload.PlayerIDs {
		playerID, err := uuid.Parse(playerIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid player ID"})
			return
		}
		playerIDs[i] = playerID
	}

	now := time.Now()
	newGameSummary := models.GameSummary{
		GameID:       gameID,
		CasinoID:     casinoID,
		StartTime:    payload.StartTime,
		DealerID:     dealerID,
		Status:       "In Progress",
		RoundsPlayed: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := gsc.DB.Create(&newGameSummary)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Add players to the game summary
	for _, playerID := range playerIDs {
		gsc.DB.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", newGameSummary.ID, playerID)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newGameSummary})
}

func (gsc *GameSummaryController) UpdateGameSummary(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")
	var payload *models.UpdateGameSummaryRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var gameSummary models.GameSummary
	result := gsc.DB.First(&gameSummary, "id = ?", gameSummaryId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		return
	}

	now := time.Now()
	gameSummaryToUpdate := models.GameSummary{
		EndTime:      payload.EndTime,
		TotalPot:     payload.TotalPot,
		Status:       payload.Status,
		RoundsPlayed: payload.RoundsPlayed,
		HighestBet:   payload.HighestBet,
		UpdatedAt:    now,
	}

	if payload.WinnerID != "" {
		winnerID, err := uuid.Parse(payload.WinnerID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid winner ID"})
			return
		}
		gameSummaryToUpdate.WinnerID = winnerID
	}

	gsc.DB.Model(&gameSummary).Updates(gameSummaryToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gameSummary})
}

func (gsc *GameSummaryController) FindGameSummaryById(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var gameSummary models.GameSummary
	result := gsc.DB.Preload("Game").Preload("Casino").Preload("Players").Preload("Winner").Preload("Dealer").First(&gameSummary, "id = ?", gameSummaryId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gameSummary})
}

func (gsc *GameSummaryController) FindGameSummaries(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var gameSummaries []models.GameSummary
	results := gsc.DB.Preload("Game").Preload("Casino").Limit(intLimit).Offset(offset).Find(&gameSummaries)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(gameSummaries), "data": gameSummaries})
}

func (gsc *GameSummaryController) DeleteGameSummary(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var gameSummary models.GameSummary
	result := gsc.DB.First(&gameSummary, "id = ?", gameSummaryId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		return
	}

	gsc.DB.Delete(&gameSummary)
	ctx.JSON(http.StatusNoContent, nil)
}
