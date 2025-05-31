package errors

import (
	"callable-api/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware é um middleware que captura e trata erros de forma centralizada
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prossegue com as outras funções
		c.Next()

		// Verifica se há erros
		if len(c.Errors) > 0 {
			// Pega o último erro
			err := c.Errors.Last()

			// Verifica se é um ValidationError
			if validationErr, ok := err.Err.(*ValidationError); ok {
				// Caso especial para erros de validação
				apiError := validationErr.ToAPIError()

				// Registra o erro no log
				logger.Error("Validation error", map[string]interface{}{
					"error":  validationErr.Error(),
					"type":   validationErr.Type,
					"fields": validationErr.FieldErrors,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				})

				c.JSON(validationErr.StatusCode, apiError)
				c.Abort()
				return
			}

			// Verifica se é um AppError
			if appError, ok := err.Err.(*AppError); ok {
				// Cria resposta API padronizada
				apiError := appError.ToAPIError()

				// Registra o erro no log
				logger.Error("Request error", map[string]interface{}{
					"error":   appError.Error(),
					"type":    appError.Type,
					"status":  appError.StatusCode,
					"stack":   appError.Stack,
					"details": appError.Details,
					"path":    c.Request.URL.Path,
					"method":  c.Request.Method,
				})

				// Responde com erro adequado
				c.JSON(appError.StatusCode, apiError)
				c.Abort()
				return
			}

			// Erro genérico, não é um AppError nem ValidationError
			appError := NewInternalServerError("Ocorreu um erro inesperado", err.Err)

			// Registra o erro no log
			logger.Error("Unexpected error", map[string]interface{}{
				"error":  err.Err.Error(),
				"stack":  appError.Stack,
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})

			// Responde com erro adequado
			c.JSON(http.StatusInternalServerError, appError.ToAPIError())
			c.Abort()
		}
	}
}

// HandleErrors é um helper para manipular erros em handlers
func HandleErrors(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Se já é um AppError ou ValidationError, usa diretamente
	if _, ok := err.(*AppError); ok {
		c.Error(err)
		return
	}

	if _, ok := err.(*ValidationError); ok {
		c.Error(err)
		return
	}

	// Caso contrário, cria um erro interno
	c.Error(NewInternalServerError("Erro interno ao processar requisição", err))
}
