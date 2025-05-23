package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste com o middleware
	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Cria uma requisição de teste
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)

	// Registra a resposta
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verifica se a resposta foi bem-sucedida
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTokenAuthMiddleware(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste com o middleware
	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "access granted"})
	})

	// Caso 1: Requisição sem token
	req1, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Verifica se o acesso foi negado
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// Caso 2: Requisição com token
	req2, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("Authorization", "test-token")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Verifica se o acesso foi permitido
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestValidationErrorMiddleware(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Cria um router de teste com o middleware
	r := gin.New()
	r.Use(ValidationErrorMiddleware())

	r.POST("/validate", func(c *gin.Context) {
		var input InputData
		if err := c.ShouldBindJSON(&input); err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "valid"})
	})

	// Simula requisição com erro de validação
	req, _ := http.NewRequest(http.MethodPost, "/validate", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Deve retornar erro de validação
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
