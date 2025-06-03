package secrets

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SecretManager interface para acesso a segredos
type SecretManager interface {
	GetSecret(ctx context.Context, secretName string) (string, error)
	GetSecretWithCache(ctx context.Context, secretName string, cacheDuration time.Duration) (string, error)
}

// GCPSecretManager implementa SecretManager para GCP
type GCPSecretManager struct {
	projectID string
	cache     map[string]cachedSecret
	mutex     sync.RWMutex
	// Mapa simulado de segredos para testes
	mockSecrets map[string]string
}

type cachedSecret struct {
	value      string
	expiration time.Time
}

// NewGCPSecretManager cria uma nova instância do gerenciador de segredos GCP simulado
func NewGCPSecretManager(projectID string) SecretManager {
	// Criar alguns segredos simulados para testes
	mockSecrets := map[string]string{
		"api-key":        "mock-api-key-12345",
		"database-pass":  "mock-db-password",
		"jwt-secret":     "mock-jwt-secret-token",
		"storage-key":    "mock-storage-access-key",
		"test-secret":    "mock-test-secret-value",
		"webhook-token":  "mock-webhook-auth-token",
	}

	return &GCPSecretManager{
		projectID:   projectID,
		cache:       make(map[string]cachedSecret),
		mockSecrets: mockSecrets,
	}
}

// GetSecret busca um segredo do Secret Manager simulado
func (m *GCPSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	// Verificar se o segredo existe no mapa de simulação
	if val, exists := m.mockSecrets[secretName]; exists {
		fmt.Printf("[MOCK] Acessando segredo simulado: %s\n", secretName)
		return val, nil
	}
	
	// Se o segredo não existe no mapa de simulação, retornamos um valor padrão com o nome do segredo
	mockValue := fmt.Sprintf("mock-value-for-%s", secretName)
	fmt.Printf("[MOCK] Criando segredo simulado on-demand: %s\n", secretName)
	return mockValue, nil
}

// GetSecretWithCache busca um segredo com cache simulado
func (m *GCPSecretManager) GetSecretWithCache(ctx context.Context, secretName string, cacheDuration time.Duration) (string, error) {
	now := time.Now()

	// Check cache (thread-safe)
	m.mutex.RLock()
	cached, exists := m.cache[secretName]
	m.mutex.RUnlock()

	if exists && now.Before(cached.expiration) {
		fmt.Printf("[MOCK] Usando segredo em cache: %s\n", secretName)
		return cached.value, nil
	}

	// Buscar um novo valor
	value, err := m.GetSecret(ctx, secretName)
	if err != nil {
		return "", err
	}

	// Atualizar cache (thread-safe)
	m.mutex.Lock()
	m.cache[secretName] = cachedSecret{
		value:      value,
		expiration: now.Add(cacheDuration),
	}
	m.mutex.Unlock()
	
	fmt.Printf("[MOCK] Segredo atualizado no cache: %s\n", secretName)
	return value, nil
}