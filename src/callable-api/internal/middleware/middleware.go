// internal/middleware/middleware.go
package middleware

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"callable-api/internal/models"
	"callable-api/pkg/logger"
)

// RequestLogger middleware for logging request information
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request start time
		startTime := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate processing time
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		
		// Get request details
		requestPath := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		
		// Log with structured logger
		logger.Info("Request processed", map[string]interface{}{
			"timestamp":  endTime.Format("2006/01/02 - 15:04:05"),
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
			"method":     method,
			"path":       requestPath,
		})
	}
}

// TokenAuthMiddleware for JWT token verification
func TokenAuthMiddleware(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// Log for debugging
		logger.Debug("Authorization header received", map[string]interface{}{
			"header": authHeader,
		})
		
		// Check if header exists
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  "error",
				Message: "Authorization token required",
			})
			c.Abort()
			return
		}
		
		// Extract token
		var token string
		
		// Check if in format "Bearer {token}"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Try to use entire value as token
			token = authHeader
		}
		
		// For demo purposes, we check against the configured API token
		// In a real app, this would validate a JWT token
		if token == "" || (apiToken != "" && token != apiToken) {
			logger.Warn("Authentication failed", map[string]interface{}{
				"reason": "Invalid or empty token",
			})
			
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  "error",
				Message: "Invalid or empty token",
			})
			c.Abort()
			return
		}
		
		// If verification passes, continue
		c.Next()
	}
}

// ValidationErrorMiddleware for handling validation errors
func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Check for errors after processing
		if len(c.Errors) > 0 {
			logger.Warn("Validation errors", map[string]interface{}{
				"errors": c.Errors.String(),
			})
			
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "error",
				Message: "Validation error: " + c.Errors.String(),
			})
			c.Abort()
			return
		}
	}
}