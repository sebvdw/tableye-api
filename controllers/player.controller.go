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

type PlayerController struct {
	DB *gorm.DB
}

func NewPlayerController(DB *gorm.DB) PlayerController {
	return PlayerController{DB}
}

// CreatePlayer godoc
//	@Summary		Create a new player
//	@Description	Create a new player with the input payload
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			player	body		models.CreatePlayerRequest	true	"Create player request"
//	@Success		201		{object}	models.Player
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]interface{}
//	@Failure		502		{object}	map[string]interface{}
//	@Router			/players [post]
func (pc *PlayerController) CreatePlayer(ctx *gin.Context) {
	var payload *models.CreatePlayerRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newPlayer := models.Player{
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
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Player with that nickname already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPlayer})
}

// UpdatePlayer godoc
//	@Summary		Update a player
//	@Description	Update a player's information by ID
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string						true	"Player ID"
//	@Param			player		body		models.UpdatePlayerRequest	true	"Update player request"
//	@Success		200			{object}	models.Player
//	@Failure		400			{object}	map[string]interface{}
//	@Failure		404			{object}	map[string]interface{}
//	@Router			/players/{playerId} [put]
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

// FindPlayerById godoc
//	@Summary		Get a player by ID
//	@Description	Get details of a player by ID
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string	true	"Player ID"
//	@Success		200			{object}	models.Player
//	@Failure		404			{object}	map[string]interface{}
//	@Router			/players/{playerId} [get]
func (pc *PlayerController) FindPlayerById(ctx *gin.Context) {
	playerId := ctx.Param("playerId")

	var player models.Player
	result := pc.DB.Preload("PlayedGames").First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": player})
}

// FindPlayers godoc
//	@Summary		List players
//	@Description	Get a list of players
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"				default(1)
//	@Param			limit	query		int	false	"Number of items per page"	default(10)
//	@Success		200		{object}	map[string]interface{}
//	@Failure		502		{object}	map[string]interface{}
//	@Router			/players [get]
func (pc *PlayerController) FindPlayers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var players []models.Player
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&players)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(players), "data": players})
}

// DeletePlayer godoc
//	@Summary		Delete a player
//	@Description	Delete a player by ID
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path	string	true	"Player ID"
//	@Success		204			"No Content"
//	@Failure		404			{object}	map[string]interface{}
//	@Router			/players/{playerId} [delete]
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

// FindPlayerStats godoc
//	@Summary		Get player statistics
//	@Description	Get statistics of a player by ID
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string	true	"Player ID"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		404			{object}	map[string]interface{}
//	@Router			/players/{playerId}/stats [get]
func (pc *PlayerController) FindPlayerStats(ctx *gin.Context) {
	playerId := ctx.Param("playerId")

	var player models.Player
	result := pc.DB.First(&player, "id = ?", playerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No player with that ID exists"})
		return
	}

	// Calculate win rate
	var winCount int64
	pc.DB.Model(&models.GameSummary{}).
		Joins("JOIN game_players ON game_summaries.id = game_players.game_summary_id").
		Where("game_players.player_id = ? AND game_summaries.status = ?", player.ID, "Completed").
		Count(&winCount)

	winRate := float64(0)
	if player.GamesPlayed > 0 {
		winRate = float64(winCount) / float64(player.GamesPlayed) * 100
	}

	stats := gin.H{
		"total_winnings": player.TotalWinnings,
		"games_played":   player.GamesPlayed,
		"win_rate":       winRate,
		"rank":           player.Rank,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": stats})
}
