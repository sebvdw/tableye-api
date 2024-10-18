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

func TestSignUpUser(t *testing.T) {
	router := GetTestRouter()

	t.Run("Successful SignUp", func(t *testing.T) {
		// Generate random email and password
		randomEmail := fmt.Sprintf("test%d@example.com", time.Now().UnixNano())
		randomPassword := fmt.Sprintf("password%d", time.Now().UnixNano())

		payload := models.SignUpInput{
			Name:            "Test User",
			Email:           randomEmail,
			Password:        randomPassword,
			PasswordConfirm: randomPassword,
		}

		jsonPayload, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
	})

	t.Run("Password Mismatch", func(t *testing.T) {
		// Generate random email
		randomEmail := fmt.Sprintf("test%d@example.com", time.Now().UnixNano())

		payload := models.SignUpInput{
			Name:            "Test User",
			Email:           randomEmail,
			Password:        "password123",
			PasswordConfirm: "password456",
		}

		jsonPayload, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "Passwords do not match", response["message"])
	})
}

func TestSignInUser(t *testing.T) {
	router := GetTestRouter()

	// Helper function to create a user if not exists
	createUserIfNotExists := func(email, password string) {
		// Check if user exists
		w := httptest.NewRecorder()
		signInPayload := models.SignInInput{
			Email:    email,
			Password: password,
		}
		jsonSignInPayload, _ := json.Marshal(signInPayload)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonSignInPayload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code == http.StatusBadRequest {
			// User doesn't exist, create one
			signUpPayload := models.SignUpInput{
				Name:            "Test User",
				Email:           email,
				Password:        password,
				PasswordConfirm: password,
			}
			jsonSignUpPayload, _ := json.Marshal(signUpPayload)
			w = httptest.NewRecorder()
			req, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonSignUpPayload))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusCreated, w.Code)
		}
	}

	t.Run("Successful SignIn", func(t *testing.T) {
		email := "signin_test@example.com"
		password := "password123"

		// Ensure user exists
		createUserIfNotExists(email, password)

		// Attempt to sign in
		signInPayload := models.SignInInput{
			Email:    email,
			Password: password,
		}
		jsonSignInPayload, _ := json.Marshal(signInPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonSignInPayload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.SignInResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotNil(t, response.User)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		payload := models.SignInInput{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		}
		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "Invalid email or Password", response["message"])
	})
}
