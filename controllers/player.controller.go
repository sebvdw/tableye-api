package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type PlayerController struct {
	DB *gorm.DB
}

func NewPlayerController(DB *gorm.DB) PlayerController {
	return PlayerController{DB}
}

func (pc *PlayerController) CreatePlayer(ctx *gin.Context) {
	var payload *models.CreatePlayerRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID"})
		return
	}

	now := time.Now()
	newPlayer := models.Player{
		UserID:        userID,
		Nickname:      payload.Nickname,
		TotalWinnings: 0,
		GamesPlayed:   0,
		Rank:          "Beginner",
		Status:        "Active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := pc.DB.Create(&newPlayer)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Player with that user ID already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPlayer})
}

func (pc *PlayerController) UpdatePlayer(ctx *gin.Context) {
	playerId := ctx.Param("playerId")
	var payload *models.UpdatePlayerRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var player models.Player
	result := pc.DB.First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	now := time.Now()
	playerToUpdate := models.Player{
		Nickname:      payload.Nickname,
		TotalWinnings: payload.TotalWinnings,
		GamesPlayed:   payload.GamesPlayed,
		Rank:          payload.Rank,
		Status:        payload.Status,
		UpdatedAt:     now,
	}

	pc.DB.Model(&player).Updates(playerToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": player})
}

func (pc *PlayerController) FindPlayerById(ctx *gin.Context) {
	playerId := ctx.Param("playerId")

	var player models.Player
	result := pc.DB.Preload("User").Preload("PlayedGames").Preload("WonGames").First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": player})
}

func (pc *PlayerController) FindPlayers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var players []models.Player
	results := pc.DB.Preload("User").Limit(intLimit).Offset(offset).Find(&players)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(players), "data": players})
}

func (pc *PlayerController) DeletePlayer(ctx *gin.Context) {
	playerId := ctx.Param("playerId")

	var player models.Player
	result := pc.DB.First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	pc.DB.Delete(&player)
	ctx.JSON(http.StatusNoContent, nil)
}

func (pc *PlayerController) FindPlayerStats(ctx *gin.Context) {
	playerId := ctx.Param("playerId")

	var player models.Player
	result := pc.DB.Preload("WonGames").First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	// Calculate win rate
	winRate := float64(0)
	if player.GamesPlayed > 0 {
		winRate = float64(len(player.WonGames)) / float64(player.GamesPlayed) * 100
	}

	stats := gin.H{
		"total_winnings": player.TotalWinnings,
		"games_played":   player.GamesPlayed,
		"win_rate":       winRate,
		"rank":           player.Rank,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": stats})
}
