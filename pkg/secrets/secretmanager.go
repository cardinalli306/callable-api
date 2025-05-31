package secrets

import (
	"context"
	"fmt"
	"sync"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
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
}

type cachedSecret struct {
	value      string
	expiration time.Time
}

// NewGCPSecretManager cria uma nova instância do gerenciador de segredos GCP
func NewGCPSecretManager(projectID string) SecretManager {
	return &GCPSecretManager{
		projectID: projectID,
		cache:     make(map[string]cachedSecret),
	}
}

// GetSecret busca um segredo do Secret Manager
func (m *GCPSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	// Criar cliente
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("falha ao criar cliente do secretmanager: %v", err)
	}
	defer client.Close()

	// Construir o nome do recurso
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", m.projectID, secretName)

	// Acessar o segredo
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("falha ao acessar versão do segredo: %v", err)
	}

	return string(result.Payload.Data), nil
}

// GetSecretWithCache busca um segredo com cache para reduzir chamadas à API
func (m *GCPSecretManager) GetSecretWithCache(ctx context.Context, secretName string, cacheDuration time.Duration) (string, error) {
	now := time.Now()

	// Check cache (thread-safe)
	m.mutex.RLock()
	cached, exists := m.cache[secretName]
	m.mutex.RUnlock()

	if exists && now.Before(cached.expiration) {
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

	return value, nil
}
