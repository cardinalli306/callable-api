package middleware

import (
	"net/http"
	
	"time"

	"github.com/gin-gonic/gin"

	"callable-api/internal/models"
	"callable-api/pkg/logger"
)

// LoggerMiddleware para registrar informações da requisição
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tempo de início da requisição
		startTime := time.Now()
		
		// Processa a requisição
		c.Next()
		
		// Calcula o tempo de processamento
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		
		// Obtém detalhes da requisição
		requestPath := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		
		// Registra com logger estruturado
		logger.Info("Requisição processada", map[string]interface{}{
			"timestamp":  endTime.Format("2006/01/02 - 15:04:05"),
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
			"method":     method,
			"path":       requestPath,
		})
	}
}

// TokenAuthMiddleware para verificação de token simples (compatibilidade)
func TokenAuthMiddleware(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// Verifica se o header existe
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  "error",
				Message: "Authorization token required",
			})
			c.Abort()
			return
		}
		
		// Extrai o token
		var token string
		
		// Verifica se está no formato "Bearer {token}"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Tenta usar o valor inteiro como token
			token = authHeader
		}
		
		// Para fins de demonstração, verificamos contra o token API configurado
		if token == "" || (apiToken != "" && token != apiToken) {
			logger.Warn("Falha de autenticação", map[string]interface{}{
				"reason": "Token inválido ou vazio",
			})
			
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  "error",
				Message: "Invalid or empty token",
			})
			c.Abort()
			return
		}
		
		// Se a verificação passar, continua
		c.Next()
	}
}

// ValidationErrorMiddleware para tratamento de erros de validação
func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Verifica erros após processamento
		if len(c.Errors) > 0 {
			logger.Warn("Erros de validação", map[string]interface{}{
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

// CORSMiddleware configura as políticas CORS
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        c.Next()
    }
}

// RequestLogger mantido para compatibilidade
func RequestLogger() gin.HandlerFunc {
	return LoggerMiddleware()
}