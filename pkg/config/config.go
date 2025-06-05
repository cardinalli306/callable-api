package config

import (
	"os"
	"strconv"
	"time"
)

// Config representa as configurações da aplicação
type Config struct {
	// Servidor
	ServerPort string
	ServerHost string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret            string
	JWTExpiration        time.Duration
	JWTRefreshSecret     string
	JWTRefreshExpiration time.Duration

	// GCP - Novas configurações
	GCPProjectID     string
	GCPStorageBucket string
	UseSecretManager bool
	UseCloudLogging  bool
	LoggingName      string

	// Novas configurações
	LogLevel          string
	Port              string
	ReadTimeoutSecs   int
	WriteTimeoutSecs  int
	GracefulTimeoutSecs int
	JWTIssuer string
	JWTExpirationMinutes int
	JWTRefreshExpirationDays int // Adicionado para o tempo de expiração do refresh token
}

// Load carrega as configurações do ambiente
func Load() *Config {
	cfg := &Config{}

	// Servidor
	cfg.ServerPort = getEnv("SERVER_PORT", "8080")
	cfg.ServerHost = getEnv("SERVER_HOST", "")

	// Database
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "postgres")
	cfg.DBName = getEnv("DB_NAME", "postgres")
	cfg.DBSSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT
	cfg.JWTSecret = getEnv("JWT_SECRET", "default-secret-key")
	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "3600"))
	cfg.JWTExpiration = time.Duration(jwtExp) * time.Second

	cfg.JWTRefreshSecret = getEnv("JWT_REFRESH_SECRET", "default-refresh-secret-key")
	jwtRefreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION", "604800"))
	cfg.JWTRefreshExpiration = time.Duration(jwtRefreshExp) * time.Second

	// GCP configurações
	cfg.GCPProjectID = getEnv("GCP_PROJECT_ID", "")
	cfg.GCPStorageBucket = getEnv("GCP_STORAGE_BUCKET", "")
	cfg.UseSecretManager = getEnv("USE_SECRET_MANAGER", "false") == "true"
	cfg.UseCloudLogging = getEnv("USE_CLOUD_LOGGING", "false") == "true"
	cfg.LoggingName = getEnv("LOGGING_NAME", "api-service")

	// Novas configurações
	cfg.LogLevel = getEnv("LOG_LEVEL", "debug") // Alterado de "info" para "debug" para mais detalhes nos logs

	cfg.Port = getEnv("PORT", "8080")

	// Aumentando os timeouts
	readTimeout, _ := strconv.Atoi(getEnv("READ_TIMEOUT_SECS", "60"))  // 60 segundos
	cfg.ReadTimeoutSecs = readTimeout

	writeTimeout, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT_SECS", "60")) // 60 segundos
	cfg.WriteTimeoutSecs = writeTimeout

	gracefulTimeout, _ := strconv.Atoi(getEnv("GRACEFUL_TIMEOUT_SECS", "30")) // Aumentado de 15 para 30 segundos
	cfg.GracefulTimeoutSecs = gracefulTimeout

	jwtExpirationMinutes, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_MINUTES", "720"))
	cfg.JWTExpirationMinutes = jwtExpirationMinutes

	jwtRefreshExpirationDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_DAYS", "7"))
	cfg.JWTRefreshExpirationDays = jwtRefreshExpirationDays

	cfg.JWTIssuer = getEnv("JWT_ISSUER", "callable-api")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}