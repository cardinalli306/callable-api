// pkg/config/config.go
package config

import (
    "os"
    "strings"
)

// Config contém todas as configurações da aplicação
type Config struct {
    // Servidor
    Port               string
    ReadTimeoutSecs    int
    WriteTimeoutSecs   int
    GracefulTimeoutSecs int
    
    // Logging
    LogLevel           string // "debug", "info", "warn", "error"
    
    // CORS
    AllowedOrigins     []string
    
    // Autenticação para demo
    DemoApiToken       string
}

// Load carrega configurações de variáveis de ambiente ou usa valores padrão
func Load() *Config {
    // Valores default para desenvolvimento
    cfg := &Config{
        Port:               "8080",
        ReadTimeoutSecs:    15,
        WriteTimeoutSecs:   15,
        GracefulTimeoutSecs: 5,
        LogLevel:           "debug",
        AllowedOrigins:     []string{"localhost:*", "127.0.0.1:*"},
        DemoApiToken:       "api-token-123", // Token de demonstração
    }
    
    // Sobrescrever com variáveis de ambiente se existirem
    if port := os.Getenv("API_PORT"); port != "" {
        cfg.Port = port
    }
    
    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }
    
    if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
        cfg.AllowedOrigins = strings.Split(origins, ",")
    }
    
    if token := os.Getenv("DEMO_API_TOKEN"); token != "" {
        cfg.DemoApiToken = token
    }
    
    return cfg
}