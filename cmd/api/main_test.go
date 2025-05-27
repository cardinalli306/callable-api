package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/models"
	"callable-api/pkg/config"
)

func TestSetupRouter(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Load config
	cfg := config.Load()

	// Test the router setup function
	router := SetupRouter(cfg)
	assert.NotNil(t, router)

	// Test health endpoint
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestSetupEnv(t *testing.T) {
	// Test debug mode
	debugCfg := &config.Config{
		LogLevel: "debug",
	}
	SetupEnv(debugCfg)
	assert.Equal(t, gin.DebugMode, gin.Mode())

	// Test release mode
	releaseCfg := &config.Config{
		LogLevel: "info",
	}
	SetupEnv(releaseCfg)
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}

func TestSetupServer(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		ReadTimeoutSecs:  10,
		WriteTimeoutSecs: 10,
	}
	router := gin.New()
	server := SetupServer(cfg, router)

	assert.Equal(t, ":8080", server.Addr)
	assert.Equal(t, 10*time.Second, server.ReadTimeout)
	assert.Equal(t, 10*time.Second, server.WriteTimeout)
}

func TestIntegrationHealthCheck(t *testing.T) {
	// Use the actual router setup from main.go
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	router := SetupRouter(cfg)

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
	// Use the actual router setup from main.go
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	router := SetupRouter(cfg)

	// Test GET /api/v1/data endpoint
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/data", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegrationGetDataById(t *testing.T) {
	// Use the actual router setup from main.go
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	router := SetupRouter(cfg)

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
	// Use the actual router setup from main.go
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	router := SetupRouter(cfg)

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
	req.Header.Set("Authorization", "Bearer "+cfg.DemoApiToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegrationPostDataWithoutAuth(t *testing.T) {
	// Use the actual router setup from main.go
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	router := SetupRouter(cfg)

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