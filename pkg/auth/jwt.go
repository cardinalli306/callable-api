package auth

import (
	"callable-api/internal/models"
	"callable-api/pkg/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims representa o payload do JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokenPair gera um par de tokens JWT (access e refresh)
func GenerateTokenPair(user *models.User, cfg *config.Config) (*models.TokenPair, error) {
	// Configurações para o token de acesso
	accessTokenExpiry := time.Now().Add(time.Minute * time.Duration(cfg.JWTExpirationMinutes))
	accessClaims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWTIssuer,
		},
	}

	// Criar token de acesso
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Configurações para o token de atualização
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * time.Duration(cfg.JWTRefreshExpirationDays))
	refreshClaims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWTIssuer,
		},
	}

	// Criar token de atualização
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWTRefreshSecret))
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateToken valida um token JWT
func ValidateToken(tokenString string, isRefresh bool, cfg *config.Config) (*Claims, error) {
	var secret string
	if isRefresh {
		secret = cfg.JWTRefreshSecret
	} else {
		secret = cfg.JWTSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Garantir que o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}