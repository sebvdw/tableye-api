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

func TestPlayerController(t *testing.T) {
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

	t.Run("CreatePlayer", func(t *testing.T) {
		payload := models.CreatePlayerRequest{
			Nickname: fmt.Sprintf("TestPlayer%d", time.Now().UnixNano()),
		}

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/players/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetPlayers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/players/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
