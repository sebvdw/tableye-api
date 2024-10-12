package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(DB *gorm.DB) TransactionController {
	return TransactionController{DB}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction with the given input data
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Create transaction request"
// @Success 201 {object} models.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /transactions [post]
func (tc *TransactionController) CreateTransaction(ctx *gin.Context) {
	var payload *models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	gameSummaryID, err := uuid.Parse(payload.GameSummaryID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid game summary ID"})
		return
	}

	playerID, err := uuid.Parse(payload.PlayerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid player ID"})
		return
	}

	now := time.Now()
	newTransaction := models.Transaction{
		GameSummaryID: gameSummaryID,
		PlayerID:      playerID,
		Amount:        payload.Amount,
		Type:          payload.Type,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := tc.DB.Create(&newTransaction)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newTransaction})
}

// GetTransactionsByGameSummary godoc
// @Summary Get transactions for a game summary
// @Description Get all transactions associated with a specific game summary
// @Tags transactions
// @Accept json
// @Produce json
// @Param gameSummaryId path string true "Game Summary ID"
// @Success 200 {array} models.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /transactions/game-summary/{gameSummaryId} [get]
func (tc *TransactionController) GetTransactionsByGameSummary(ctx *gin.Context) {
	gameSummaryId := ctx.Param("gameSummaryId")

	var transactions []models.Transaction
	result := tc.DB.Where("game_summary_id = ?", gameSummaryId).Find(&transactions)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No transactions found for this game summary"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": transactions})
}

// Additional methods like UpdateTransaction, DeleteTransaction, etc. can be added here as needed.
