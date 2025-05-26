package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/handlers"
	"callable-api/internal/models"
)

func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.GET("/health", handlers.HealthCheck)

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the response body
	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "API is running", response.Message)
}

func TestGetData(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.GET("/api/v1/data", handlers.GetData)

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "/api/v1/data", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the response body
	var response models.ListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Data retrieved successfully", response.Message)

	// Verify there is data in the response
	assert.NotNil(t, response.Data)
}

func TestGetDataById(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.GET("/api/v1/data/:id", handlers.GetDataById)

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "/api/v1/data/123", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the response body
	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Verify the ID was returned correctly
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "123", data["id"])
}

func TestPostData(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.POST("/api/v1/data", handlers.PostData)

	// Prepare test data
	input := models.InputData{
		Name:        "Test Item",
		Value:       "ABC123",
		Description: "Test Description",
		Email:       "test@example.com",
	}

	jsonData, err := json.Marshal(input)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify the response body
	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Data saved successfully", response.Message)

	// Verify data was returned correctly
	assert.NotNil(t, response.Data)
}

func TestPostDataInvalid(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.POST("/api/v1/data", handlers.PostData)

	// Prepare invalid data
	invalidInput := `{"name":"", "value":""}`

	// Create a test request
	req, err := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBufferString(invalidInput))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the error was returned correctly
	assert.Equal(t, http.StatusBadRequest, w.Code)
}