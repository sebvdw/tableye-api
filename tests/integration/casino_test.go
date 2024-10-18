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

func TestCasinoController(t *testing.T) {
	router := GetTestRouter()

	// Helper function to create a user and get access token
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

	t.Run("CreateCasino", func(t *testing.T) {
		payload := models.CreateCasinoRequest{
			Name:          fmt.Sprintf("Test Casino %d", time.Now().UnixNano()),
			Location:      "Test Location",
			LicenseNumber: fmt.Sprintf("LN%d", time.Now().UnixNano()),
			MaxCapacity:   1000,
			Status:        "Active",
		}

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
	})

	t.Run("CreateCasinoWithDuplicateName", func(t *testing.T) {
		casinoName := fmt.Sprintf("Duplicate Casino %d", time.Now().UnixNano())
		payload := models.CreateCasinoRequest{
			Name:          casinoName,
			Location:      "Test Location",
			LicenseNumber: fmt.Sprintf("LN%d", time.Now().UnixNano()),
			MaxCapacity:   1000,
			Status:        "Active",
		}

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Try to create another casino with the same name
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "fail", response["status"])
		assert.Contains(t, response["message"], "already exists")
	})

	t.Run("GetCasinos", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/casinos/", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
	})

	t.Run("GetCasinoById", func(t *testing.T) {
		// First, create a casino
		createPayload := models.CreateCasinoRequest{
			Name:          fmt.Sprintf("Test Casino for GetById %d", time.Now().UnixNano()),
			Location:      "Test Location",
			LicenseNumber: fmt.Sprintf("LN%d", time.Now().UnixNano()),
			MaxCapacity:   1000,
			Status:        "Active",
		}
		jsonCreatePayload, _ := json.Marshal(createPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonCreatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create casino")

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		assert.NoError(t, err, "Failed to unmarshal create response")

		t.Logf("Create casino response: %+v", createResponse)

		casinoData, ok := createResponse["data"].(map[string]interface{})
		assert.True(t, ok, "Failed to extract casino data from create response")

		casinoId, ok := casinoData["id"].(string)
		assert.True(t, ok, "Failed to extract casino ID from create response")

		t.Logf("Created casino ID: %s", casinoId)

		// Now, get the casino by ID
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/casinos/"+casinoId, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		t.Logf("Get casino response code: %d", w.Code)
		t.Logf("Get casino response body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code, "Failed to get casino by ID")

		var getResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		assert.NoError(t, err, "Failed to unmarshal get response")

		assert.Equal(t, "success", getResponse["status"], "Get casino response status is not success")
		assert.NotNil(t, getResponse["data"], "Get casino response data is nil")

		getCasinoData, ok := getResponse["data"].(map[string]interface{})
		assert.True(t, ok, "Failed to extract casino data from get response")

		assert.Equal(t, casinoId, getCasinoData["id"], "Retrieved casino ID does not match created casino ID")
		assert.Equal(t, createPayload.Name, getCasinoData["name"], "Retrieved casino name does not match created casino name")
	})

	t.Run("UpdateCasino", func(t *testing.T) {
		// First, create a casino
		createPayload := models.CreateCasinoRequest{
			Name:          fmt.Sprintf("Test Casino for Update %d", time.Now().UnixNano()),
			Location:      "Test Location",
			LicenseNumber: fmt.Sprintf("LN%d", time.Now().UnixNano()),
			MaxCapacity:   1000,
			Status:        "Active",
		}
		jsonCreatePayload, _ := json.Marshal(createPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonCreatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create casino")

		var createResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResponse)
		casinoData := createResponse["data"].(map[string]interface{})
		casinoId := casinoData["id"].(string)

		// Now, update the casino
		updatePayload := models.UpdateCasinoRequest{
			Name:     "Updated Casino Name",
			Location: "Updated Location",
			Status:   "Inactive",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/api/casinos/"+casinoId, bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Failed to update casino")

		var updateResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &updateResponse)
		assert.Equal(t, "success", updateResponse["status"])
		assert.NotNil(t, updateResponse["data"])

		updatedCasinoData := updateResponse["data"].(map[string]interface{})
		assert.Equal(t, "Updated Casino Name", updatedCasinoData["name"])
		assert.Equal(t, "Updated Location", updatedCasinoData["location"])
		assert.Equal(t, "Inactive", updatedCasinoData["status"])
	})

	t.Run("DeleteCasino", func(t *testing.T) {
		// First, create a casino
		createPayload := models.CreateCasinoRequest{
			Name:          fmt.Sprintf("Test Casino for Delete %d", time.Now().UnixNano()),
			Location:      "Test Location",
			LicenseNumber: fmt.Sprintf("LN%d", time.Now().UnixNano()),
			MaxCapacity:   1000,
			Status:        "Active",
		}
		jsonCreatePayload, _ := json.Marshal(createPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/casinos/", bytes.NewBuffer(jsonCreatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create casino")

		var createResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResponse)
		casinoData := createResponse["data"].(map[string]interface{})
		casinoId := casinoData["id"].(string)

		// Now, delete the casino
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/casinos/"+casinoId, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code, "Failed to delete casino")

		// Try to get the deleted casino
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/casinos/"+casinoId, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Deleted casino should not be found")
	})
}
