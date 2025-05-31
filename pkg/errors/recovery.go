package errors

import (
	"callable-api/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware é um middleware personalizado que recupera de panics
// e converte-os em respostas de erro estruturadas
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Captura stack trace
				stack := string(debug.Stack())
				
				// Log detalhado do panic
				logger.Error("Recovered from panic", map[string]interface{}{
					"error":      fmt.Sprintf("%v", r),
					"stacktrace": stack,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				})
				
				// Criar um erro AppError para padronização
				var errMsg string
				if err, ok := r.(error); ok {
					errMsg = err.Error()
				} else {
					errMsg = fmt.Sprintf("%v", r)
				}
				
				appErr := NewInternalServerError("O servidor encontrou um erro inesperado", nil).
					WithDetails(errMsg)
				
				// Responda com o erro estruturado
				apiErr := appErr.ToAPIError()
				c.JSON(http.StatusInternalServerError, apiErr)
				
				// Aborta o processamento
				c.Abort()
			}
		}()
		
		// Continua com a execução normal
		c.Next()
	}
}