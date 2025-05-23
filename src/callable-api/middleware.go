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

// Middleware para verificação de token JWT
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// Imprime o header para depuração
		fmt.Printf("Auth header: '%s'\n", authHeader)
		
		// Verifica se o header existe
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, Response{
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
		
		// Apenas para testes, aceitamos qualquer token não vazio
		if token == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid or empty token",
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