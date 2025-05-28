package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/logger"
	"callable-api/pkg/secrets"
	"callable-api/pkg/storage"
)

// GCPDemoHandler demonstra a integração com GCP
type GCPDemoHandler struct {
	config      *config.Config
	logger      logger.Logger
	secretMgr   secrets.SecretManager
	storage     *storage.CloudStorage
	jwtProvider *auth.SecretProvider
}

// NewGCPDemoHandler cria um novo handler de demonstração
func NewGCPDemoHandler(
	cfg *config.Config,
	log logger.Logger,
	secretMgr secrets.SecretManager,
	storage *storage.CloudStorage,
) *GCPDemoHandler {
	return &GCPDemoHandler{
		config:      cfg,
		logger:      log,
		secretMgr:   secretMgr,
		storage:     storage,
		jwtProvider: auth.NewSecretProvider(cfg, secretMgr, log),
	}
}

// TestIntegration testa todas as integrações
func (h *GCPDemoHandler) TestIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := map[string]interface{}{
		"status":    "success",
		"tests":     make(map[string]interface{}),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// 1. Testar Logging
	h.logger.Info("Teste de integração GCP iniciado", map[string]interface{}{
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"method":      r.Method,
	})
	response["tests"].(map[string]interface{})["logging"] = map[string]interface{}{
		"status":              "success",
		"message":             "Logs enviados com sucesso",
		"cloud_logging_enabled": h.config.UseCloudLogging,
	}

	// 2. Testar Secret Manager (com fallback)
	jwtSecret, err := h.jwtProvider.GetJWTSecret(ctx)
	if err != nil {
		h.logger.Error("Falha no teste de Secret Manager", err)
		response["tests"].(map[string]interface{})["secret_manager"] = map[string]interface{}{
			"status":         "error",
			"message":        "Falha ao acessar segredos",
			"using_fallback": true,
			"error":          err.Error(),
		}
	} else {
		secretLen := len(jwtSecret)
		secretPreview := ""
		if secretLen > 0 {
			previewLen := min(3, secretLen)
			secretPreview = jwtSecret[:previewLen] + "..."
		}

		response["tests"].(map[string]interface{})["secret_manager"] = map[string]interface{}{
			"status":               "success",
			"secret_length":        secretLen,
			"preview":              secretPreview,
			"using_secret_manager": h.config.UseSecretManager,
		}
	}

	// 3. Testar Cloud Storage
	if h.storage != nil && h.config.GCPStorageBucket != "" {
		testData := []byte("Teste de integração com Cloud Storage - " + time.Now().Format(time.RFC3339))
		objectName := fmt.Sprintf("demo/test-%s.txt", time.Now().Format("20060102-150405"))

		err := h.storage.UploadFile(ctx, objectName, bytes.NewReader(testData))
		if err != nil {
			h.logger.Error("Erro no upload para Cloud Storage", err)
			response["tests"].(map[string]interface{})["storage"] = map[string]interface{}{
				"status":  "error",
				"message": "Falha no upload",
				"error":   err.Error(),
			}
		} else {
			signedURL, urlErr := h.storage.GetSignedURL(ctx, objectName, 15*time.Minute)
			signedURLStatus := "success"
			signedURLMsg := ""

			if urlErr != nil {
				signedURLStatus = "error"
				signedURLMsg = urlErr.Error()
				h.logger.Error("Erro ao gerar URL assinada", urlErr)
			}

			response["tests"].(map[string]interface{})["storage"] = map[string]interface{}{
				"status":      "success",
				"bucket":      h.config.GCPStorageBucket,
				"object_name": objectName,
				"signed_url": map[string]interface{}{
					"status":     signedURLStatus,
					"url":        signedURL,
					"error":      signedURLMsg,
					"expiration": "15 minutos",
				},
				"data_size": len(testData),
			}
		}
	} else {
		response["tests"].(map[string]interface{})["storage"] = map[string]interface{}{
			"status":            "skipped",
			"message":           "Cloud Storage não configurado",
			"bucket_configured": h.config.GCPStorageBucket != "",
		}
	}

	// Registrar o resultado completo
	h.logger.Info("Teste de integração GCP concluído", map[string]interface{}{
		"success": true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// min retorna o menor de dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}