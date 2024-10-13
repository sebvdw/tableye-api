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
// @Description Create a new game summary with the given input data, including game, casino, dealer, and players information
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummary body models.CreateGameSummaryRequest true "Create game summary request"
// @Success 201 {object} models.GameSummaryResponse "Successfully created game summary"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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

	// Start a database transaction
	tx := gsc.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to start transaction"})
		return
	}

	// Create the game summary
	if err := tx.Create(&newGameSummary).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Add players to the game summary
	for _, playerID := range playerIDs {
		if err := tx.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", newGameSummary.ID, playerID).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add player to game summary"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to commit transaction"})
		return
	}

	// Fetch the complete game summary with all related data
	var completeSummary models.GameSummary
	if err := gsc.DB.Preload("Game").
		Preload("Casino").
		Preload("Dealer").
		Preload("Dealer.User").
		Preload("Players").
		First(&completeSummary, newGameSummary.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch complete game summary"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": completeSummary})
}

// UpdateGameSummary godoc
// @Summary Update a game summary
// @Description Update a game summary with the given input data and return the updated summary with all related information
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

	// Start a transaction
	tx := gsc.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to start transaction"})
		return
	}

	// Update the game summary
	if err := tx.Model(&gameSummary).Updates(gameSummaryToUpdate).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update game summary"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to commit transaction"})
		return
	}

	// Fetch the updated game summary with all related data
	var updatedSummary models.GameSummary
	if err := gsc.DB.Preload("Game").
		Preload("Casino").
		Preload("Dealer").
		Preload("Dealer.User").
		Preload("Players").
		Preload("Transactions").
		Preload("Transactions.Player").
		First(&updatedSummary, gameSummary.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch updated game summary"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedSummary})
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
		Preload("Casino").
		Preload("Players").
		Preload("Dealer").
		Preload("Dealer.User").
		Preload("Transactions").
		Preload("Transactions.Player").
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
		Preload("Casino").
		Preload("Players").
		Preload("Dealer").
		Preload("Dealer.User").
		Preload("Transactions").
		Preload("Transactions.Player").
		Limit(intLimit).Offset(offset).
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

// AddTransactionToGameSummary godoc
// @Summary Add a transaction to a game summary
// @Description Add a new transaction to an existing game summary
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Param transaction body models.CreateTransactionRequest true "Transaction details"
// @Success 200 {object} models.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId}/transactions [post]
func (gsc *GameSummaryController) AddTransactionToGameSummary(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var payload models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var gameSummary models.GameSummary
	if err := gsc.DB.First(&gameSummary, "id = ?", gameSummaryId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Game summary not found"})
		return
	}

	playerID, err := uuid.Parse(payload.PlayerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid player ID"})
		return
	}

	transaction := models.Transaction{
		GameSummaryID: gameSummary.ID,
		PlayerID:      playerID,
		Amount:        payload.Amount,
		Outcome:       payload.Outcome,
		//Type:          payload.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := gsc.DB.Create(&transaction).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create transaction"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": transaction})
}

// AddPlayersToGameSummary godoc
// @Summary Add players to a game summary
// @Description Add an array of players to an existing game summary
// @Tags gameSummaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Param players body []string true "Array of player IDs"
// @Success 200 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId}/players [post]
func (gsc *GameSummaryController) AddPlayersToGameSummary(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var playerIDs []string
	if err := ctx.ShouldBindJSON(&playerIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var gameSummary models.GameSummary
	if err := gsc.DB.First(&gameSummary, "id = ?", gameSummaryId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Game summary not found"})
		return
	}

	var players []models.Player
	for _, playerID := range playerIDs {
		var player models.Player
		if err := gsc.DB.First(&player, "id = ?", playerID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Player not found: " + playerID})
			return
		}
		players = append(players, player)
	}

	if err := gsc.DB.Model(&gameSummary).Association("Players").Append(players); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add players"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Players added successfully"})
}
