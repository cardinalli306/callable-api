// pkg/config/config_test.go
package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Salvar o estado atual das variáveis de ambiente relevantes
	originalVars := saveEnvVars([]string{
		"SERVER_PORT", "SERVER_HOST", 
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"JWT_SECRET", "JWT_EXPIRATION", "JWT_REFRESH_SECRET", "JWT_REFRESH_EXPIRATION",
		"GCP_PROJECT_ID", "GCP_STORAGE_BUCKET", "USE_SECRET_MANAGER", "USE_CLOUD_LOGGING", "LOGGING_NAME",
		"LOG_LEVEL", "PORT", "READ_TIMEOUT_SECS", "WRITE_TIMEOUT_SECS", "GRACEFUL_TIMEOUT_SECS",
		"JWT_ISSUER", "JWT_EXPIRATION_MINUTES", "JWT_REFRESH_EXPIRATION_DAYS",
	})
	
	// Limpar todas as variáveis antes de cada teste
	clearEnvVars(originalVars)

	// Restaurar variáveis de ambiente no final
	defer restoreEnvVars(originalVars)

	// Test case 1: Valores padrão
	t.Run("Default Values", func(t *testing.T) {
		// Limpar variáveis para garantir uso de valores padrão
		clearEnvVars(originalVars)
		
		cfg := Load()
		
		// Verificar valores padrão do servidor
		if cfg.ServerPort != "8080" {
			t.Errorf("Default ServerPort should be 8080, got %s", cfg.ServerPort)
		}
		if cfg.ServerHost != "" {
			t.Errorf("Default ServerHost should be empty, got %s", cfg.ServerHost)
		}
		
		// Verificar valores padrão do banco de dados
		if cfg.DBHost != "localhost" {
			t.Errorf("Default DBHost should be localhost, got %s", cfg.DBHost)
		}
		if cfg.DBPort != "5432" {
			t.Errorf("Default DBPort should be 5432, got %s", cfg.DBPort)
		}
		
		// Verificar valores padrão do JWT
		if cfg.JWTSecret != "default-secret-key" {
			t.Errorf("Default JWTSecret incorrect, got %s", cfg.JWTSecret)
		}
		if cfg.JWTExpiration != 3600*time.Second {
			t.Errorf("Default JWTExpiration should be 3600s, got %v", cfg.JWTExpiration)
		}
		
		// Verificar configurações de logs e timeout
		if cfg.LogLevel != "info" {
			t.Errorf("Default LogLevel should be info, got %s", cfg.LogLevel)
		}
		if cfg.ReadTimeoutSecs != 10 {
			t.Errorf("Default ReadTimeoutSecs should be 10, got %d", cfg.ReadTimeoutSecs)
		}
		if cfg.WriteTimeoutSecs != 10 {
			t.Errorf("Default WriteTimeoutSecs should be 10, got %d", cfg.WriteTimeoutSecs)
		}
		if cfg.GracefulTimeoutSecs != 15 {
			t.Errorf("Default GracefulTimeoutSecs should be 15, got %d", cfg.GracefulTimeoutSecs)
		}
		
		// Verificar valores de tempo JWT
		if cfg.JWTExpirationMinutes != 720 {
			t.Errorf("Default JWTExpirationMinutes should be 720, got %d", cfg.JWTExpirationMinutes)
		}
		if cfg.JWTRefreshExpirationDays != 7 {
			t.Errorf("Default JWTRefreshExpirationDays should be 7, got %d", cfg.JWTRefreshExpirationDays)
		}
		if cfg.JWTIssuer != "callable-api" {
			t.Errorf("Default JWTIssuer should be callable-api, got %s", cfg.JWTIssuer)
		}
	})

	// Test case 2: Sobrescrita com variáveis de ambiente
	t.Run("Environment Variables Override", func(t *testing.T) {
		// Limpar e definir as variáveis para o teste
		clearEnvVars(originalVars)
		
		os.Setenv("SERVER_PORT", "9000")
		os.Setenv("SERVER_HOST", "api.example.com")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("JWT_SECRET", "custom-secret")
		os.Setenv("JWT_EXPIRATION", "7200")
		os.Setenv("LOG_LEVEL", "error")
		os.Setenv("READ_TIMEOUT_SECS", "30")
		os.Setenv("JWT_EXPIRATION_MINUTES", "60")
		os.Setenv("JWT_ISSUER", "custom-issuer")

		cfg := Load()
		
		// Verificar se as variáveis foram aplicadas corretamente
		if cfg.ServerPort != "9000" {
			t.Errorf("ServerPort should be 9000, got %s", cfg.ServerPort)
		}
		if cfg.ServerHost != "api.example.com" {
			t.Errorf("ServerHost should be api.example.com, got %s", cfg.ServerHost)
		}
		if cfg.DBHost != "db.example.com" {
			t.Errorf("DBHost should be db.example.com, got %s", cfg.DBHost)
		}
		if cfg.DBPort != "5433" {
			t.Errorf("DBPort should be 5433, got %s", cfg.DBPort)
		}
		if cfg.JWTSecret != "custom-secret" {
			t.Errorf("JWTSecret should be custom-secret, got %s", cfg.JWTSecret)
		}
		if cfg.JWTExpiration != 7200*time.Second {
			t.Errorf("JWTExpiration should be 7200s, got %v", cfg.JWTExpiration)
		}
		if cfg.LogLevel != "error" {
			t.Errorf("LogLevel should be error, got %s", cfg.LogLevel)
		}
		if cfg.ReadTimeoutSecs != 30 {
			t.Errorf("ReadTimeoutSecs should be 30, got %d", cfg.ReadTimeoutSecs)
		}
		if cfg.JWTExpirationMinutes != 60 {
			t.Errorf("JWTExpirationMinutes should be 60, got %d", cfg.JWTExpirationMinutes)
		}
		if cfg.JWTIssuer != "custom-issuer" {
			t.Errorf("JWTIssuer should be custom-issuer, got %s", cfg.JWTIssuer)
		}
	})

	// Test case 3: Configurações parciais
	t.Run("Partial Environment Variables", func(t *testing.T) {
		// Limpar e definir apenas algumas variáveis
		clearEnvVars(originalVars)
		
		os.Setenv("SERVER_PORT", "5000")
		os.Setenv("JWT_SECRET", "partial-test-secret")

		cfg := Load()
		
		// Verificar se as variáveis definidas foram aplicadas
		if cfg.ServerPort != "5000" {
			t.Errorf("ServerPort should be 5000, got %s", cfg.ServerPort)
		}
		if cfg.JWTSecret != "partial-test-secret" {
			t.Errorf("JWTSecret should be partial-test-secret, got %s", cfg.JWTSecret)
		}
		
		// Verificar se os outros valores permaneceram como padrão
		if cfg.DBHost != "localhost" {
			t.Errorf("DBHost should fall back to default localhost, got %s", cfg.DBHost)
		}
		if cfg.LogLevel != "info" {
			t.Errorf("LogLevel should fall back to default info, got %s", cfg.LogLevel)
		}
	})
}

// Função auxiliar para salvar o estado atual das variáveis de ambiente
func saveEnvVars(vars []string) map[string]string {
	env := make(map[string]string)
	for _, key := range vars {
		env[key] = os.Getenv(key)
	}
	return env
}

// Função auxiliar para limpar as variáveis de ambiente
func clearEnvVars(env map[string]string) {
	for key := range env {
		os.Unsetenv(key)
	}
}

// Função auxiliar para restaurar as variáveis de ambiente
func restoreEnvVars(env map[string]string) {
	for key, value := range env {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}