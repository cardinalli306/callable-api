// pkg/config/config_test.go
package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// Salvar o estado atual das variáveis de ambiente
	originalPort := os.Getenv("API_PORT")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalOrigins := os.Getenv("ALLOWED_ORIGINS")
	originalToken := os.Getenv("DEMO_API_TOKEN")

	// Limpar as variáveis para o teste
	os.Unsetenv("API_PORT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ALLOWED_ORIGINS")
	os.Unsetenv("DEMO_API_TOKEN")

	// Test case 1: Default values
	t.Run("Default Values", func(t *testing.T) {
		cfg := Load()
		if cfg.Port != "8080" {
			t.Errorf("Default port should be 8080, got %s", cfg.Port)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("Default log level should be debug, got %s", cfg.LogLevel)
		}
		if !reflect.DeepEqual(cfg.AllowedOrigins, []string{"localhost:*", "127.0.0.1:*"}) {
			t.Errorf("Default allowed origins should be [localhost:*, 127.0.0.1:*], got %v", cfg.AllowedOrigins)
		}
		if cfg.DemoApiToken != "api-token-123" {
			t.Errorf("Default API token should be api-token-123, got %s", cfg.DemoApiToken)
		}
		if cfg.ReadTimeoutSecs != 15 {
			t.Errorf("Default read timeout should be 15, got %d", cfg.ReadTimeoutSecs)
		}
		if cfg.WriteTimeoutSecs != 15 {
			t.Errorf("Default write timeout should be 15, got %d", cfg.WriteTimeoutSecs)
		}
		if cfg.GracefulTimeoutSecs != 5 {
			t.Errorf("Default graceful timeout should be 5, got %d", cfg.GracefulTimeoutSecs)
		}
	})

	// Test case 2: Environment variables override
	t.Run("Environment Variables Override", func(t *testing.T) {
		os.Setenv("API_PORT", "9000")
		os.Setenv("LOG_LEVEL", "error")
		os.Setenv("ALLOWED_ORIGINS", "example.com,api.example.com")
		os.Setenv("DEMO_API_TOKEN", "custom-token")

		cfg := Load()
		if cfg.Port != "9000" {
			t.Errorf("Port should be 9000, got %s", cfg.Port)
		}
		if cfg.LogLevel != "error" {
			t.Errorf("Log level should be error, got %s", cfg.LogLevel)
		}
		expectedOrigins := []string{"example.com", "api.example.com"}
		if !reflect.DeepEqual(cfg.AllowedOrigins, expectedOrigins) {
			t.Errorf("Allowed origins should be %v, got %v", expectedOrigins, cfg.AllowedOrigins)
		}
		if cfg.DemoApiToken != "custom-token" {
			t.Errorf("API token should be custom-token, got %s", cfg.DemoApiToken)
		}
	})

	// Test case 3: Partial environment variables
	t.Run("Partial Environment Variables", func(t *testing.T) {
		os.Unsetenv("API_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Setenv("ALLOWED_ORIGINS", "test.com")
		os.Unsetenv("DEMO_API_TOKEN")

		cfg := Load()
		if cfg.Port != "8080" {
			t.Errorf("Port should fall back to default 8080, got %s", cfg.Port)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("Log level should fall back to default debug, got %s", cfg.LogLevel)
		}
		expectedOrigins := []string{"test.com"}
		if !reflect.DeepEqual(cfg.AllowedOrigins, expectedOrigins) {
			t.Errorf("Allowed origins should be %v, got %v", expectedOrigins, cfg.AllowedOrigins)
		}
		if cfg.DemoApiToken != "api-token-123" {
			t.Errorf("API token should fall back to default api-token-123, got %s", cfg.DemoApiToken)
		}
	})

	// Restore original environment variables
	if originalPort != "" {
		os.Setenv("API_PORT", originalPort)
	} else {
		os.Unsetenv("API_PORT")
	}
	if originalLogLevel != "" {
		os.Setenv("LOG_LEVEL", originalLogLevel)
	} else {
		os.Unsetenv("LOG_LEVEL")
	}
	if originalOrigins != "" {
		os.Setenv("ALLOWED_ORIGINS", originalOrigins)
	} else {
		os.Unsetenv("ALLOWED_ORIGINS")
	}
	if originalToken != "" {
		os.Setenv("DEMO_API_TOKEN", originalToken)
	} else {
		os.Unsetenv("DEMO_API_TOKEN")
	}
}