// controllers/auth.controller.go

package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/models"
	"github.com/suidevv/tableye-api/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.SignUpInput true "User registration details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /auth/register [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "dealer", // Default role is now "dealer"
		Verified:  true,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      newUser.Role,
		Provider:  newUser.Provider,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

// SignInUser godoc
// @Summary Login a user
// @Description Authenticate a user and return access token, user info, associated dealer and casino details
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.SignInInput true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Success 200 {object} models.SignInResponse
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Router /auth/login [post]
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	// Fetch associated dealer
	var dealer models.Dealer
	dealerResult := ac.DB.First(&dealer, "user_id = ?", user.ID)

	// Fetch associated casino
	var casino models.Casino
	var casinoID *string
	if dealerResult.Error == nil {
		casinoResult := ac.DB.Table("casino_dealers").
			Select("casino_id").
			Where("dealer_id = ?", dealer.ID).
			Limit(1).
			Scan(&casinoID)

		if casinoResult.Error == nil && casinoID != nil {
			ac.DB.First(&casino, "id = ?", *casinoID)
		}
	}

	response := models.SignInResponse{
		AccessToken: access_token,
		User: models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}

	if dealerResult.Error == nil {
		response.Dealer = &models.DealerResponse{
			ID:         dealer.ID,
			DealerCode: dealer.DealerCode,
			Status:     dealer.Status,
		}
	}

	if casinoID != nil {
		response.Casino = &models.CasinoResponse{
			ID:   casino.ID,
			Name: casino.Name,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// RefreshAccessToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags authentication
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no longer exists"})
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(config.AccessTokenExpiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access_token, err := token.SignedString([]byte(config.AccessTokenPrivateKey))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

// LogoutUser godoc
// @Summary Logout user
// @Description Clear authentication cookies
// @Tags authentication
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	config, _ := initializers.LoadConfig(".")

	ctx.SetCookie("access_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
