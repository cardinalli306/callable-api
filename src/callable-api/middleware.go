package main

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

// Middleware para logging de requisições
func RequestLogger() gin.HandlerFunc {
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
		
		// Log formatado
		fmt.Printf("[API] %v | %3d | %13v | %15s | %s | %s\n",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			requestPath,
		)
	}
}

// Middleware para verificação de token (simulado)
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		
		// Simples verificação de token (em produção seria mais complexo)
		if token == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "API token required",
			})
			c.Abort()
			return
		}
		
		// Se passar pela verificação, continua
		c.Next()
	}
}

// Middleware para tratamento de erros de validação
func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Verifica se há erros após o processamento
		if len(c.Errors) > 0 {
			c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Validation error: " + c.Errors.String(),
			})
			c.Abort()
			return
		}
	}
}