// controllers/admin.controller.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(DB *gorm.DB) AdminController {
	return AdminController{DB}
}

// AssignAdminRole godoc
// @Summary Assign admin role to a user
// @Description Assigns the admin role to a specified user
// @Tags admin
// @Accept json
// @Produce json
// @Param request body AdminRoleAssignRequest true "User ID to assign admin role"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/assign-admin [post]
func (ac *AdminController) AssignAdminRole(ctx *gin.Context) {
	var payload AdminRoleAssignRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID"})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	user.Role = "admin"
	ac.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User role updated to admin"})
}

// AdminRoleAssignRequest represents the request body for assigning admin role
type AdminRoleAssignRequest struct {
	UserID string `json:"userId" binding:"required"`
}
