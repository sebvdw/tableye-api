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

// CreateGameSummary godoc
// @Summary Create a new game summary
// @Description Create a new game summary with the given input data
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummary body models.CreateGameSummaryRequest true "Create game summary request"
// @Success 201 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries [post]
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Add players to the game summary
	for _, playerID := range playerIDs {
		gsc.DB.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", newGameSummary.ID, playerID)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newGameSummary})
}

// UpdateGameSummary godoc
// @Summary Update a game summary
// @Description Update a game summary with the given input data
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Param gameSummary body models.UpdateGameSummaryRequest true "Update game summary request"
// @Success 200 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [put]
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

	gsc.DB.Model(&gameSummary).Updates(gameSummaryToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gameSummary})
}

// FindGameSummaryById godoc
// @Summary Get a game summary by ID
// @Description Get details of a game summary by its ID, including transactions
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Success 200 {object} models.GameSummaryResponse
// @Failure 404 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [get]
func (gsc *GameSummaryController) FindGameSummaryById(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var gameSummary models.GameSummary
	result := gsc.DB.Preload("Game").
		Preload("Casino", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Owner")
		}).
		Preload("Players").
		Preload("Dealer").
		Preload("Transactions").
		First(&gameSummary, "id = ?", gameSummaryId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gameSummary})
}

// FindGameSummaries godoc
// @Summary List game summaries
// @Description Get a list of game summaries with pagination, including transactions
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries [get]
func (gsc *GameSummaryController) FindGameSummaries(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var gameSummaries []models.GameSummary
	results := gsc.DB.Preload("Game").
		Preload("Casino", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Owner")
		}).
		Preload("Dealer").
		Preload("Transactions").
		Limit(intLimit).
		Offset(offset).
		Find(&gameSummaries)

	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(gameSummaries), "data": gameSummaries})
}

// DeleteGameSummary godoc
// @Summary Delete a game summary
// @Description Delete a game summary by its ID
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [delete]
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
