package auth

import (
	"context"
	"time"

	"callable-api/pkg/config"
	"callable-api/pkg/logger"
	"callable-api/pkg/secrets"
)

const (
	// Nomes dos segredos no Secret Manager
	JWTSecretName        = "jwt-secret"
	JWTRefreshSecretName = "jwt-refresh-secret"

	// Duração do cache de segredos (para evitar muitas chamadas à API)
	secretCacheDuration = 1 * time.Hour
)

// SecretProvider gerencia as chaves para JWT
type SecretProvider struct {
	config    *config.Config
	secretMgr secrets.SecretManager
	logger    logger.Logger
}

// NewSecretProvider cria um novo provedor de segredos
func NewSecretProvider(cfg *config.Config, secretMgr secrets.SecretManager, log logger.Logger) *SecretProvider {
	return &SecretProvider{
		config:    cfg,
		secretMgr: secretMgr,
		logger:    log,
	}
}

// GetJWTSecret obtém a chave secreta para tokens JWT (do Secret Manager ou config)
func (p *SecretProvider) GetJWTSecret(ctx context.Context) (string, error) {
	// Se não estamos usando Secret Manager, use o valor da config
	if !p.config.UseSecretManager || p.config.GCPProjectID == "" || p.secretMgr == nil {
		p.logger.Debug("Usando chave JWT da configuração local")
		return p.config.JWTSecret, nil
	}

	// Buscar do Secret Manager com cache
	secret, err := p.secretMgr.GetSecretWithCache(ctx, JWTSecretName, secretCacheDuration)
	if err != nil {
		p.logger.Error("Falha ao buscar JWT secret do Secret Manager", err)
		// Fallback para o valor da config em caso de falha
		return p.config.JWTSecret, nil
	}

	p.logger.Debug("JWT secret obtido do Secret Manager")
	return secret, nil
}

// GetJWTRefreshSecret obtém a chave de refresh para tokens JWT
func (p *SecretProvider) GetJWTRefreshSecret(ctx context.Context) (string, error) {
	// Se não estamos usando Secret Manager, use o valor da config
	if !p.config.UseSecretManager || p.config.GCPProjectID == "" || p.secretMgr == nil {
		p.logger.Debug("Usando chave JWT refresh da configuração local")
		return p.config.JWTRefreshSecret, nil
	}

	// Buscar do Secret Manager com cache
	secret, err := p.secretMgr.GetSecretWithCache(ctx, JWTRefreshSecretName, secretCacheDuration)
	if err != nil {
		p.logger.Error("Falha ao buscar JWT refresh secret do Secret Manager", err)
		// Fallback para o valor da config em caso de falha
		return p.config.JWTRefreshSecret, nil
	}

	p.logger.Debug("JWT refresh secret obtido do Secret Manager")
	return secret, nil
}
