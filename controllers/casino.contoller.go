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

type CasinoController struct {
	DB *gorm.DB
}

func NewCasinoController(DB *gorm.DB) CasinoController {
	return CasinoController{DB}
}

// CreateCasino godoc
// @Summary Create a new casino
// @Description Create a new casino with the input payload
// @Tags casinos
// @Accept json
// @Produce json
// @Param casino body models.CreateCasinoRequest true "Create casino request"
// @Success 201 {object} models.Casino
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /casinos [post]
func (cc *CasinoController) CreateCasino(ctx *gin.Context) {
	var payload *models.CreateCasinoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	ownerID, err := uuid.Parse(payload.OwnerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid owner ID"})
		return
	}

	newCasino := models.Casino{
		Name:          payload.Name,
		Location:      payload.Location,
		LicenseNumber: payload.LicenseNumber,
		Description:   payload.Description,
		OpeningHours:  payload.OpeningHours,
		Website:       payload.Website,
		PhoneNumber:   payload.PhoneNumber,
		MaxCapacity:   payload.MaxCapacity,
		Status:        payload.Status,
		OwnerID:       ownerID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := cc.DB.Create(&newCasino)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Casino with that name or license number already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newCasino})
}

// UpdateCasino godoc
// @Summary Update a casino
// @Description Update a casino with the input payload
// @Tags casinos
// @Accept json
// @Produce json
// @Param casinoId path string true "Casino ID"
// @Param casino body models.UpdateCasinoRequest true "Update casino request"
// @Success 200 {object} models.Casino
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /casinos/{casinoId} [put]
func (cc *CasinoController) UpdateCasino(ctx *gin.Context) {
	casinoId := ctx.Param("casinoId")
	var payload *models.UpdateCasinoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var casino models.Casino
	result := cc.DB.First(&casino, "id = ?", casinoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No casino with that ID exists"})
		return
	}

	now := time.Now()
	casinoToUpdate := models.Casino{
		Name:          payload.Name,
		Location:      payload.Location,
		LicenseNumber: payload.LicenseNumber,
		Description:   payload.Description,
		OpeningHours:  payload.OpeningHours,
		Website:       payload.Website,
		PhoneNumber:   payload.PhoneNumber,
		MaxCapacity:   payload.MaxCapacity,
		Status:        payload.Status,
		Rating:        payload.Rating,
		UpdatedAt:     now,
	}

	cc.DB.Model(&casino).Updates(casinoToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": casino})
}

// FindCasinoById godoc
// @Summary Get a casino by ID
// @Description Get a single casino by its ID
// @Tags casinos
// @Produce json
// @Param casinoId path string true "Casino ID"
// @Success 200 {object} models.Casino
// @Failure 404 {object} map[string]interface{}
// @Router /casinos/{casinoId} [get]
func (cc *CasinoController) FindCasinoById(ctx *gin.Context) {
	casinoId := ctx.Param("casinoId")

	var casino models.Casino
	result := cc.DB.Preload("Owner").Preload("Dealers").Preload("Games").First(&casino, "id = ?", casinoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No casino with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": casino})
}

// FindCasinos godoc
// @Summary List casinos
// @Description Get a list of casinos
// @Tags casinos
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /casinos [get]
func (cc *CasinoController) FindCasinos(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var casinos []models.Casino
	results := cc.DB.Preload("Owner").Limit(intLimit).Offset(offset).Find(&casinos)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	// Create a slice to hold the modified casino data
	casinosResponse := make([]gin.H, len(casinos))

	for i, casino := range casinos {
		casinoResponse := gin.H{
			"id":             casino.ID,
			"name":           casino.Name,
			"location":       casino.Location,
			"license_number": casino.LicenseNumber,
			"description":    casino.Description,
			"opening_hours":  casino.OpeningHours,
			"website":        casino.Website,
			"phone_number":   casino.PhoneNumber,
			"max_capacity":   casino.MaxCapacity,
			"status":         casino.Status,
			"rating":         casino.Rating,
			"created_at":     casino.CreatedAt,
			"updated_at":     casino.UpdatedAt,
			"owner": gin.H{
				"id":    casino.Owner.ID,
				"name":  casino.Owner.Name,
				"email": casino.Owner.Email,
			},
		}
		casinosResponse[i] = casinoResponse
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(casinosResponse), "data": casinosResponse})
}

// DeleteCasino godoc
// @Summary Delete a casino
// @Description Delete a casino by its ID
// @Tags casinos
// @Produce json
// @Param casinoId path string true "Casino ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /casinos/{casinoId} [delete]
func (cc *CasinoController) DeleteCasino(ctx *gin.Context) {
	casinoId := ctx.Param("casinoId")

	var casino models.Casino
	result := cc.DB.First(&casino, "id = ?", casinoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No casino with that ID exists"})
		return
	}

	cc.DB.Delete(&casino)
	ctx.JSON(http.StatusNoContent, nil)
}
