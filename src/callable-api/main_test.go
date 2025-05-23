package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Função helper para configurar o router de teste
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/data", getData)
		v1.GET("/data/:id", getDataById)

		auth := v1.Group("/")
		auth.Use(TokenAuthMiddleware())
		{
			auth.POST("/data", postData)
		}
	}

	return router
}

func TestIntegrationHealthCheck(t *testing.T) {
	router := setupTestRouter()

	// Testa endpoint de health check
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestIntegrationGetData(t *testing.T) {
	router := setupTestRouter()

	// Testa endpoint GET /api/v1/data
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/data", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegrationGetDataById(t *testing.T) {
	router := setupTestRouter()

	// Testa endpoint GET /api/v1/data/:id
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/data/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "123", data["id"])
}

func TestIntegrationPostDataWithAuth(t *testing.T) {
	router := setupTestRouter()

	// Prepara dados para POST
	input := InputData{
		Name:        "Integration Test Item",
		Value:       "INT123",
		Description: "Integration Test Description",
	}

	jsonData, _ := json.Marshal(input)

	// Testa POST com token
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegrationPostDataWithoutAuth(t *testing.T) {
	router := setupTestRouter()

	// Prepara dados para POST
	input := InputData{
		Name:        "Integration Test Item",
		Value:       "INT123",
		Description: "Integration Test Description",
	}

	jsonData, _ := json.Marshal(input)

	// Testa POST sem token
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
