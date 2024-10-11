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

type DealerController struct {
	DB *gorm.DB
}

func NewDealerController(DB *gorm.DB) DealerController {
	return DealerController{DB}
}

// CreateDealer godoc
// @Summary Create a new dealer
// @Description Create a new dealer with the input payload
// @Tags dealers
// @Accept json
// @Produce json
// @Param dealer body models.CreateDealerRequest true "Create dealer request"
// @Success 201 {object} models.Dealer
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /dealers [post]
func (dc *DealerController) CreateDealer(ctx *gin.Context) {
	var payload *models.CreateDealerRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID"})
		return
	}

	newDealer := models.Dealer{
		UserID:     userID,
		DealerCode: payload.DealerCode,
		Status:     payload.Status,
		GamesDealt: 0,
		Rating:     0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result := dc.DB.Create(&newDealer)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Dealer with that dealer code already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newDealer})
}

// UpdateDealer godoc
// @Summary Update a dealer
// @Description Update a dealer with the input payload
// @Tags dealers
// @Accept json
// @Produce json
// @Param dealerId path string true "Dealer ID"
// @Param dealer body models.UpdateDealerRequest true "Update dealer request"
// @Success 200 {object} models.Dealer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /dealers/{dealerId} [put]
func (dc *DealerController) UpdateDealer(ctx *gin.Context) {
	dealerId := ctx.Param("dealerId")
	var payload *models.UpdateDealerRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var dealer models.Dealer
	result := dc.DB.First(&dealer, "id = ?", dealerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No dealer with that ID exists"})
		return
	}

	now := time.Now()
	dealerToUpdate := models.Dealer{
		Status:       payload.Status,
		GamesDealt:   payload.GamesDealt,
		Rating:       payload.Rating,
		LastActiveAt: payload.LastActiveAt,
		UpdatedAt:    now,
	}

	dc.DB.Model(&dealer).Updates(dealerToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dealer})
}

// FindDealerById godoc
// @Summary Get a dealer by ID
// @Description Get a single dealer by its ID
// @Tags dealers
// @Produce json
// @Param dealerId path string true "Dealer ID"
// @Success 200 {object} models.Dealer
// @Failure 404 {object} map[string]interface{}
// @Router /dealers/{dealerId} [get]
func (dc *DealerController) FindDealerById(ctx *gin.Context) {
	dealerId := ctx.Param("dealerId")

	var dealer models.Dealer
	result := dc.DB.Preload("User").Preload("Casinos").Preload("GameSummaries").First(&dealer, "id = ?", dealerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No dealer with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dealer})
}

// FindDealers godoc
// @Summary List dealers
// @Description Get a list of dealers
// @Tags dealers
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /dealers [get]
func (dc *DealerController) FindDealers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var dealers []models.Dealer
	results := dc.DB.Preload("User").Limit(intLimit).Offset(offset).Find(&dealers)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(dealers), "data": dealers})
}

// DeleteDealer godoc
// @Summary Delete a dealer
// @Description Delete a dealer by its ID
// @Tags dealers
// @Produce json
// @Param dealerId path string true "Dealer ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /dealers/{dealerId} [delete]
func (dc *DealerController) DeleteDealer(ctx *gin.Context) {
	dealerId := ctx.Param("dealerId")

	var dealer models.Dealer
	result := dc.DB.First(&dealer, "id = ?", dealerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No dealer with that ID exists"})
		return
	}

	dc.DB.Delete(&dealer)
	ctx.JSON(http.StatusNoContent, nil)
}
