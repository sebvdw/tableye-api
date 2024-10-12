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
// @Success 201 {object} models.Transaction
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /transactions [post]
func (tc *TransactionController) CreateTransaction(ctx *gin.Context) {
	var payload *models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	playerID, err := uuid.Parse(payload.PlayerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid player ID"})
		return
	}

	gameSummaryID, err := uuid.Parse(payload.GameSummaryID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid player ID"})
		return
	}

	if payload.Outcome == "win" && payload.Amount < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Win amount cannot be negative"})
		return
	}
	if payload.Outcome == "loss" && payload.Amount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Loss amount should be negative"})
		return
	}

	now := time.Now()
	newTransaction := models.Transaction{
		GameSummaryID: gameSummaryID,
		PlayerID:      playerID,
		Amount:        payload.Amount,
		Type:          payload.Type,
		Outcome:       payload.Outcome,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	result := tc.DB.Create(&newTransaction)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	var player models.Player
	if err := tc.DB.First(&player, "id = ?", playerID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve player data"})
		return
	}

	response := struct {
		Transaction models.Transaction `json:"transaction"`
		Player      models.Player      `json:"player"`
	}{
		Transaction: newTransaction,
		Player:      player,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": response})
}

// FindTransactions godoc
// @Summary List transactions
// @Description Get a list of transactions with pagination
// @Tags transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /transactions [get]
func (tc *TransactionController) FindTransactions(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var transactions []models.Transaction
	results := tc.DB.Preload("Player").Limit(intLimit).Offset(offset).Find(&transactions)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(transactions), "data": transactions})
}

// FindTransactionById godoc
// @Summary Get a transaction by ID
// @Description Get details of a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param transactionId path string true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 404 {object} map[string]interface{}
// @Router /transactions/{transactionId} [get]
func (tc *TransactionController) FindTransactionById(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	var transaction models.Transaction
	result := tc.DB.Preload("Player").First(&transaction, "id = ?", transactionId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No transaction with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": transaction})
}

func (tc *TransactionController) UpdateTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	var payload models.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var transaction models.Transaction
	result := tc.DB.First(&transaction, "id = ?", transactionId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No transaction with that ID exists"})
		return
	}

	updatedTransaction := models.Transaction{
		Amount: payload.Amount,
		Type:   payload.Type,
	}

	tc.DB.Model(&transaction).Updates(updatedTransaction)

	// Fetch the updated transaction with player data
	tc.DB.Preload("Player").First(&transaction, "id = ?", transactionId)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": transaction})
}

func (tc *TransactionController) DeleteTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	result := tc.DB.Delete(&models.Transaction{}, "id = ?", transactionId)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No transaction with that ID exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
