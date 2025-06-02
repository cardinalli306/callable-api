package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"

	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/logger"
	"callable-api/pkg/secrets"
)

// GCPDemoHandler gerencia as rotas de demonstração da integração GCP
type GCPDemoHandler struct {
	config      *config.Config
	logger      logger.Logger
	secretMgr   secrets.SecretManager
	storage     *storage.Client
	jwtProvider *auth.SecretProvider
}

// NewGCPDemoHandler cria um novo handler de demonstração
func NewGCPDemoHandler(
	cfg *config.Config,
	log logger.Logger,
	secretMgr secrets.SecretManager,
	storage *storage.Client,
) *GCPDemoHandler {
	var jwtProvider *auth.SecretProvider
	// Cria o jwtProvider apenas se todos os componentes necessários estiverem disponíveis
	if cfg != nil && secretMgr != nil && log != nil {
		jwtProvider = auth.NewSecretProvider(cfg, secretMgr, log)
	}
	
	return &GCPDemoHandler{
		config:      cfg,
		logger:      log,
		secretMgr:   secretMgr,
		storage:     storage,
		jwtProvider: jwtProvider,
	}
}

// TestIntegration testa todas as integrações
// @Summary Teste de integração GCP
// @Description Testa a integração com vários serviços Google Cloud Platform (Logging, Secret Manager e Storage)
// @Tags gcp
// @Produce json
// @Success 200 {object} map[string]interface{} "Resultado dos testes de integração"
// @Failure 503 {object} map[string]interface{} "Erro de serviços não disponíveis"
// @Router /api/test/gcp [get]
func (h *GCPDemoHandler) TestIntegration(c *gin.Context) {
	// Verificar se os serviços GCP necessários estão disponíveis
	if h.logger == nil || h.secretMgr == nil || h.storage == nil || h.jwtProvider == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "GCP integration not configured",
			"details": gin.H{
				"logger_available":       h.logger != nil,
				"secret_mgr_available":   h.secretMgr != nil,
				"storage_available":      h.storage != nil,
				"jwt_provider_available": h.jwtProvider != nil,
			},
		})
		return
	}

	// Se chegou aqui, os serviços estão disponíveis
	ctx := c.Request.Context()
	response := gin.H{
		"status":    "success",
		"tests":     gin.H{},
		"timestamp": time.Now().Format(time.RFC3339),
	}
	tests := response["tests"].(gin.H)

	// Teste de logging
	// Corrigido para usar a versão variádica conforme definido na interface Logger
	h.logger.Info("Teste de integração GCP iniciado", map[string]interface{}{
		"handler": "GCPDemoHandler",
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})
	tests["logging"] = gin.H{
		"status":  "success",
		"message": "Log enviado com sucesso",
	}

	// Teste de Secret Manager
	secretTest, err := h.testSecretManager(ctx)
	tests["secret_manager"] = secretTest
	if err != nil {
		response["status"] = "partial_success"
	}

	// Teste de Storage
	storageTest, err := h.testStorage(ctx)
	tests["storage"] = storageTest
	if err != nil {
		response["status"] = "partial_success"
	}

	// Retorna a resposta
	c.JSON(http.StatusOK, response)
}

// testSecretManager testa o Secret Manager
func (h *GCPDemoHandler) testSecretManager(ctx context.Context) (gin.H, error) {
	result := gin.H{
		"status": "success",
	}

	// Verifica se o JWT Provider está disponível
	if h.jwtProvider == nil {
		result["status"] = "error"
		result["message"] = "JWT Provider não configurado"
		return result, fmt.Errorf("jwt provider não configurado")
	}

	// Tenta buscar o segredo JWT
	_, err := h.jwtProvider.GetJWTSecret(ctx)
	if err != nil {
		result["status"] = "error"
		result["message"] = fmt.Sprintf("Erro ao buscar segredo JWT: %v", err)
		return result, err
	}

	result["message"] = "Secret Manager testado com sucesso"
	return result, nil
}

// testStorage testa o Cloud Storage
func (h *GCPDemoHandler) testStorage(ctx context.Context) (gin.H, error) {
	result := gin.H{
		"status": "success",
	}

	// Verifica se o Storage e o bucket estão configurados
	if h.storage == nil || h.config.GCPStorageBucket == "" {
		result["status"] = "error"
		result["message"] = "Cloud Storage não configurado"
		return result, fmt.Errorf("cloud storage não configurado")
	}

	// Lista os objetos do bucket para testar o acesso
	bucket := h.storage.Bucket(h.config.GCPStorageBucket)
	it := bucket.Objects(ctx, nil)
	
	// Tenta buscar o primeiro objeto apenas para testar
	_, err := it.Next()
	if err != nil && err != storage.ErrObjectNotExist {
		result["status"] = "error"
		result["message"] = fmt.Sprintf("Erro ao listar objetos do bucket: %v", err)
		return result, err
	}

	result["message"] = "Cloud Storage testado com sucesso"
	result["bucket"] = h.config.GCPStorageBucket
	return result, nil
}

// RegisterRoutes registra as rotas do handler
// @Summary Registrar rotas GCP
// @Description Registra as rotas de teste de integração com GCP no router
// @Tags internal
func (h *GCPDemoHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/test/gcp", h.TestIntegration)
}