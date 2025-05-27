package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/middleware"
	"callable-api/internal/models"
	"callable-api/pkg/config"
)

func TestRequestLogger(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify the response was successful
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTokenAuthMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Get configuration
	cfg := config.Load()
	testToken := "test-token"
	if cfg.DemoApiToken != "" {
		testToken = cfg.DemoApiToken
	}

	// Create a test router with the middleware
	r := gin.New()
	r.Use(middleware.TokenAuthMiddleware(testToken))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "access granted"})
	})

	// Case 1: Request without token
	req1, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Verify access was denied
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// Case 2: Request with token
	req2, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("Authorization", testToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Verify access was granted
	assert.Equal(t, http.StatusOK, w2.Code)
	
	// Case 3: Request with Bearer token format
	req3, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req3.Header.Set("Authorization", "Bearer "+testToken)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	// Verify access was granted with Bearer format
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestValidationErrorMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	r := gin.New()
	r.Use(middleware.ValidationErrorMiddleware())

	// Create a test route that validates JSON input
	r.POST("/validate", func(c *gin.Context) {
		var input models.InputData
		if err := c.ShouldBindJSON(&input); err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "valid"})
	})

	// Simulate request with validation error (empty JSON)
	req, _ := http.NewRequest(http.MethodPost, "/validate", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return validation error
	assert.Equal(t, http.StatusBadRequest, w.Code)
}