package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suidevv/tableye-api/models"
)

func TestGameController(t *testing.T) {
	router := GetTestRouter()

	getAccessToken := func() string {
		signInPayload := models.SignInInput{
			Email:    "user13@example.com",
			Password: "password13",
		}
		jsonSignInPayload, _ := json.Marshal(signInPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonSignInPayload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var signInResponse models.SignInResponse
		json.Unmarshal(w.Body.Bytes(), &signInResponse)
		return signInResponse.AccessToken
	}

	accessToken := getAccessToken()

	t.Run("CreateGame", func(t *testing.T) {
		payload := models.CreateGameRequest{
			Name:       fmt.Sprintf("Test Game %d", time.Now().UnixNano()),
			Type:       "Poker",
			MaxPlayers: 8,
			MinPlayers: 2,
			MinBet:     10,
			MaxBet:     1000,
		}

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/games/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetGames", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/games/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
