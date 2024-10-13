package controllers

import (
	"errors"
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
// @Description Create a new game summary with the input payload
// @Tags game-summaries
// @Accept json
// @Produce json
// @Param payload body models.CreateGameSummaryRequest true "Create game summary payload"
// @Success 201 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries [post]
func (gsc *GameSummaryController) CreateGameSummary(ctx *gin.Context) {
	var payload models.CreateGameSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newGameSummary, err := gsc.createGameSummaryFromPayload(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := gsc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newGameSummary).Error; err != nil {
			return err
		}
		return gsc.addPlayersToGameSummary(tx, newGameSummary.ID, payload.PlayerIDs)
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create game summary"})
		return
	}

	response, err := gsc.getGameSummaryResponse(newGameSummary.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch complete game summary"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": response})
}

// UpdateGameSummary godoc
// @Summary Update a game summary
// @Description Update an existing game summary
// @Tags game-summaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Param payload body models.UpdateGameSummaryRequest true "Update game summary payload"
// @Success 200 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [put]
func (gsc *GameSummaryController) UpdateGameSummary(ctx *gin.Context) {
	gameSummaryId, err := uuid.Parse(ctx.Param("gameSummaryId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid game summary ID"})
		return
	}

	var payload models.UpdateGameSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := gsc.DB.Model(&models.GameSummary{}).Where("id = ?", gameSummaryId).Updates(payload).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update game summary"})
		}
		return
	}

	response, err := gsc.getGameSummaryResponse(gameSummaryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch updated game summary"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

// FindGameSummaryById godoc
// @Summary Get a game summary by ID
// @Description Retrieve a game summary by its ID
// @Tags game-summaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Success 200 {object} models.GameSummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [get]
func (gsc *GameSummaryController) FindGameSummaryById(ctx *gin.Context) {
	gameSummaryId, err := uuid.Parse(ctx.Param("gameSummaryId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid game summary ID"})
		return
	}

	response, err := gsc.getGameSummaryResponse(gameSummaryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch game summary"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

// FindGameSummaries godoc
// @Summary List game summaries
// @Description Retrieve a list of game summaries with pagination
// @Tags game-summaries
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} models.GameSummaryResponse
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries [get]
func (gsc *GameSummaryController) FindGameSummaries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var gameSummaries []models.GameSummary
	if err := gsc.DB.Limit(limit).Offset(offset).Find(&gameSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch game summaries"})
		return
	}

	responses, err := gsc.getMultipleGameSummaryResponses(gameSummaries)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process game summaries"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(responses), "data": responses})
}

// DeleteGameSummary godoc
// @Summary Delete a game summary
// @Description Delete a game summary by its ID
// @Tags game-summaries
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /game-summaries/{gameSummaryId} [delete]
func (gsc *GameSummaryController) DeleteGameSummary(ctx *gin.Context) {
	gameSummaryId, err := uuid.Parse(ctx.Param("gameSummaryId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid game summary ID"})
		return
	}

	result := gsc.DB.Delete(&models.GameSummary{}, gameSummaryId)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete game summary"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No game summary with that ID exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (gsc *GameSummaryController) createGameSummaryFromPayload(payload models.CreateGameSummaryRequest) (models.GameSummary, error) {
	gameID, err := uuid.Parse(payload.GameID)
	if err != nil {
		return models.GameSummary{}, errors.New("invalid game ID")
	}

	casinoID, err := uuid.Parse(payload.CasinoID)
	if err != nil {
		return models.GameSummary{}, errors.New("invalid casino ID")
	}

	dealerID, err := uuid.Parse(payload.DealerID)
	if err != nil {
		return models.GameSummary{}, errors.New("invalid dealer ID")
	}

	return models.GameSummary{
		ID:           uuid.New(),
		GameID:       gameID,
		CasinoID:     casinoID,
		StartTime:    payload.StartTime,
		DealerID:     dealerID,
		Status:       "In Progress",
		RoundsPlayed: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (gsc *GameSummaryController) addPlayersToGameSummary(tx *gorm.DB, gameSummaryID uuid.UUID, playerIDs []string) error {
	for _, playerIDStr := range playerIDs {
		playerID, err := uuid.Parse(playerIDStr)
		if err != nil {
			return errors.New("invalid player ID")
		}
		if err := tx.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", gameSummaryID, playerID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (gsc *GameSummaryController) getGameSummaryResponse(id uuid.UUID) (models.GameSummaryResponse, error) {
	var gameSummary models.GameSummary
	if err := gsc.DB.Preload("Dealer").Preload("Players").Preload("Transactions.Player").First(&gameSummary, id).Error; err != nil {
		return models.GameSummaryResponse{}, err
	}

	var game models.Game
	if err := gsc.DB.First(&game, gameSummary.GameID).Error; err != nil {
		return models.GameSummaryResponse{}, err
	}

	var casino models.Casino
	if err := gsc.DB.First(&casino, gameSummary.CasinoID).Error; err != nil {
		return models.GameSummaryResponse{}, err
	}

	return convertToGameSummaryResponse(gameSummary, game, casino), nil
}

func (gsc *GameSummaryController) getMultipleGameSummaryResponses(gameSummaries []models.GameSummary) ([]models.GameSummaryResponse, error) {
	var responses []models.GameSummaryResponse
	for _, summary := range gameSummaries {
		response, err := gsc.getGameSummaryResponse(summary.ID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func convertToGameSummaryResponse(gameSummary models.GameSummary, game models.Game, casino models.Casino) models.GameSummaryResponse {
	return models.GameSummaryResponse{
		ID:           gameSummary.ID,
		Game:         models.GameResponse{ID: game.ID, Name: game.Name},
		Casino:       models.CasinoResponse{ID: casino.ID, Name: casino.Name},
		StartTime:    gameSummary.StartTime,
		EndTime:      gameSummary.EndTime,
		Players:      convertToPlayerResponses(gameSummary.Players),
		Dealer:       convertToDealerResponse(gameSummary.Dealer),
		TotalPot:     gameSummary.TotalPot,
		Status:       gameSummary.Status,
		RoundsPlayed: gameSummary.RoundsPlayed,
		HighestBet:   gameSummary.HighestBet,
		Transactions: convertToTransactionResponses(gameSummary.Transactions),
		CreatedAt:    gameSummary.CreatedAt,
		UpdatedAt:    gameSummary.UpdatedAt,
	}
}

func convertToPlayerResponses(players []models.Player) []models.PlayerResponse {
	responses := make([]models.PlayerResponse, len(players))
	for i, player := range players {
		responses[i] = models.PlayerResponse{
			ID:            player.ID,
			Nickname:      player.Nickname,
			TotalWinnings: player.TotalWinnings,
			Rank:          player.Rank,
			Status:        player.Status,
			CreatedAt:     player.CreatedAt,
			UpdatedAt:     player.UpdatedAt,
		}
	}
	return responses
}

func convertToDealerResponse(dealer models.Dealer) models.GameSummaryDealerResponse {
	return models.GameSummaryDealerResponse{
		ID:           dealer.ID,
		DealerCode:   dealer.DealerCode,
		Status:       dealer.Status,
		GamesDealt:   dealer.GamesDealt,
		Rating:       dealer.Rating,
		LastActiveAt: dealer.LastActiveAt,
		CreatedAt:    dealer.CreatedAt,
		UpdatedAt:    dealer.UpdatedAt,
	}
}

func convertToTransactionResponses(transactions []models.Transaction) []models.TransactionResponse {
	responses := make([]models.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = models.TransactionResponse{
			ID:        transaction.ID,
			Player:    convertToPlayerResponse(transaction.Player),
			Amount:    transaction.Amount,
			Outcome:   transaction.Outcome,
			CreatedAt: transaction.CreatedAt,
			UpdatedAt: transaction.UpdatedAt,
		}
	}
	return responses
}

func convertToPlayerResponse(player models.Player) models.PlayerResponse {
	return models.PlayerResponse{
		ID:            player.ID,
		Nickname:      player.Nickname,
		TotalWinnings: player.TotalWinnings,
		Rank:          player.Rank,
		Status:        player.Status,
		CreatedAt:     player.CreatedAt,
		UpdatedAt:     player.UpdatedAt,
	}
}
