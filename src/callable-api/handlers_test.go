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

func TestHealthCheck(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste
	r := gin.Default()
	r.GET("/health", healthCheck)

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica o código de status
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica o corpo da resposta
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "API is running", response.Message)
}

func TestGetData(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste
	r := gin.Default()
	r.GET("/api/v1/data", getData)

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodGet, "/api/v1/data", nil)
	assert.NoError(t, err)

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica o código de status
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica o corpo da resposta
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Data retrieved successfully", response.Message)

	// Verifica se há dados na resposta
	assert.NotNil(t, response.Data)
}

func TestGetDataById(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste
	r := gin.Default()
	r.GET("/api/v1/data/:id", getDataById)

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodGet, "/api/v1/data/123", nil)
	assert.NoError(t, err)

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica o código de status
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica o corpo da resposta
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Verifica se o ID foi retornado corretamente
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "123", data["id"])
}

func TestPostData(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste
	r := gin.Default()
	r.POST("/api/v1/data", postData)

	// Prepara os dados de teste
	input := InputData{
		Name:        "Test Item",
		Value:       "ABC123",
		Description: "Test Description",
		Email:       "test@example.com",
	}

	jsonData, err := json.Marshal(input)
	assert.NoError(t, err)

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica o código de status
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verifica o corpo da resposta
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Data saved successfully", response.Message)

	// Verifica se os dados foram retornados corretamente
	assert.NotNil(t, response.Data)
}

func TestPostDataInvalid(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste
	r := gin.Default()
	r.POST("/api/v1/data", postData)

	// Prepara dados inválidos
	invalidInput := `{"name":"", "value":""}`

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodPost, "/api/v1/data", bytes.NewBufferString(invalidInput))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica se o erro foi retornado corretamente
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
