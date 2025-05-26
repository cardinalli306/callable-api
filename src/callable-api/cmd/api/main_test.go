package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/handlers"
	"callable-api/internal/middleware"
	"callable-api/internal/models"
	"callable-api/pkg/config"
)

// Helper function to set up the test router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/data", handlers.GetData)
		v1.GET("/data/:id", handlers.GetDataById)

		auth := v1.Group("/")
		
		// Load configuration for token
		cfg := config.Load()
		auth.Use(middleware.TokenAuthMiddleware(cfg.DemoApiToken))
		{
			auth.POST("/data", handlers.PostData)
		}
	}

	return router
}

func TestIntegrationHealthCheck(t *testing.T) {
	router := setupTestRouter()

	// Test health check endpoint
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestIntegrationGetData(t *testing.T) {
	router := setupTestRouter()

	// Test GET /api/v1/data endpoint
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/data", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegrationGetDataById(t *testing.T) {
	router := setupTestRouter()

	// Test GET /api/v1/data/:id endpoint
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/data/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "123", data["id"])
}

func TestIntegrationPostDataWithAuth(t *testing.T) {
	router := setupTestRouter()

	// Prepare data for POST
	input := models.InputData{
		Name:        "Integration Test Item",
		Value:       "INT123",
		Description: "Integration Test Description",
	}

	jsonData, _ := json.Marshal(input)

	// Test POST with token
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Use the correct token from configuration
	cfg := config.Load()
	req.Header.Set("Authorization", "Bearer "+cfg.DemoApiToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegrationPostDataWithoutAuth(t *testing.T) {
	router := setupTestRouter()

	// Prepare data for POST
	input := models.InputData{
		Name:        "Integration Test Item",
		Value:       "INT123",
		Description: "Integration Test Description",
	}

	jsonData, _ := json.Marshal(input)

	// Test POST without token
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}