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

func TestGameSummaryController(t *testing.T) {
	router := GetTestRouter()

	var accessToken string
	var dealerID string
	var casinoID string

	// Helper function to sign in and get necessary IDs
	signInAndGetIDs := func() {
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
		accessToken = signInResponse.AccessToken

		if signInResponse.Dealer != nil {
			dealerID = signInResponse.Dealer.ID.String()
		} else {
			t.Fatal("Logged in user is not a dealer")
		}

		if signInResponse.Casino != nil {
			casinoID = signInResponse.Casino.ID.String()
		} else {
			t.Fatal("Logged in user is not associated with a casino")
		}
	}

	signInAndGetIDs()

	// Helper function to create a game and return its ID
	createGame := func() string {
		gamePayload := models.CreateGameRequest{
			Name:       fmt.Sprintf("Test Game %d", time.Now().UnixNano()),
			Type:       "Poker",
			MaxPlayers: 8,
			MinPlayers: 2,
			MinBet:     10,
			MaxBet:     1000,
		}
		jsonGamePayload, _ := json.Marshal(gamePayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/games/", bytes.NewBuffer(jsonGamePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		var createGameResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createGameResponse)
		gameData := createGameResponse["data"].(map[string]interface{})
		return gameData["id"].(string)
	}

	// Helper function to get player IDs
	getPlayerIDs := func(count int) []string {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/players/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		players := response["data"].([]interface{})
		playerIDs := make([]string, 0, count)
		for i := 0; i < count && i < len(players); i++ {
			player := players[i].(map[string]interface{})
			playerIDs = append(playerIDs, player["id"].(string))
		}

		return playerIDs
	}

	t.Run("CreateGameSummary", func(t *testing.T) {
		gameID := createGame()
		playerIDs := getPlayerIDs(2) // Get 2 player IDs

		payload := models.CreateGameSummaryRequest{
			GameID:    gameID,
			CasinoID:  casinoID,
			StartTime: time.Now(),
			DealerID:  dealerID,
			PlayerIDs: playerIDs,
		}

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/game-summaries/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
	})

	t.Run("GetGameSummaries", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/game-summaries/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
	})

	t.Run("GetGameSummaryById", func(t *testing.T) {
		gameID := createGame()
		playerIDs := getPlayerIDs(2) // Get 2 player IDs

		createPayload := models.CreateGameSummaryRequest{
			GameID:    gameID,
			CasinoID:  casinoID,
			StartTime: time.Now(),
			DealerID:  dealerID,
			PlayerIDs: playerIDs,
		}

		jsonCreatePayload, _ := json.Marshal(createPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/game-summaries/", bytes.NewBuffer(jsonCreatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create game summary")

		var createResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResponse)
		gameSummaryData := createResponse["data"].(map[string]interface{})
		gameSummaryID := gameSummaryData["id"].(string)

		// Now, get the game summary by ID
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/game-summaries/"+gameSummaryID, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Failed to get game summary by ID")

		var getResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &getResponse)
		assert.Equal(t, "success", getResponse["status"])
		assert.NotNil(t, getResponse["data"])

		retrievedGameSummaryData := getResponse["data"].(map[string]interface{})
		assert.Equal(t, gameSummaryID, retrievedGameSummaryData["id"], "Retrieved game summary ID does not match created game summary ID")
	})
}
