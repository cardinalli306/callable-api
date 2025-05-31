package middleware

import (
	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/errors"
	"callable-api/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware verifica a validade do token JWT
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o token Authorization do header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			err := errors.NewUnauthorizedError("Token de autenticação não fornecido", nil)
			errors.HandleErrors(c, err)
			c.Abort()
			return
		}

		// O header deve ter o formato "Bearer {token}"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			err := errors.NewUnauthorizedError("Formato de token inválido", nil)
			errors.HandleErrors(c, err)
			c.Abort()
			return
		}

		tokenString := headerParts[1]

		// Validar o token
		claims, err := auth.ValidateToken(tokenString, false, cfg)
		if err != nil {
			logger.Error("Falha na validação do token", map[string]interface{}{
				"error": err.Error(),
			})
			err := errors.NewUnauthorizedError("Token inválido ou expirado", nil)
			errors.HandleErrors(c, err)
			c.Abort()
			return
		}

		// Armazenar os claims no contexto para uso posterior
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userName", claims.Name)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireRole verifica se o usuário tem um papel específico
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			err := errors.NewForbiddenError("Acesso negado", nil)
			errors.HandleErrors(c, err)
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			logger.Warn("Tentativa de acesso não autorizado", map[string]interface{}{
				"requiredRoles": roles,
				"userRole":      userRole,
				"path":          c.Request.URL.Path,
				"method":        c.Request.Method,
			})
			err := errors.NewForbiddenError("Você não tem permissão para acessar este recurso", nil)
			errors.HandleErrors(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}