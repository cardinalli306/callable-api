package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "callable-api/docs" // Para geração de documentação Swagger
	"callable-api/internal/handlers"
	"callable-api/internal/middleware"
	"callable-api/internal/repository"
	"callable-api/internal/service"
	"callable-api/pkg/config"
	"callable-api/pkg/errors"
	"callable-api/pkg/logger"

	// Importações para GCP
	gcplogger "callable-api/pkg/logger" // Renomeando para evitar conflito
	"callable-api/pkg/secrets"
	mystorage "callable-api/pkg/storage"
)

// @title Callable API
// @version 1.0
// @description Uma API robusta construída em Go usando o framework Gin, oferecendo endpoints para gerenciamento de dados com validação completa e autenticação JWT.
// @contact.name Desenvolvedor
// @contact.email dev@exemplo.com
// @contact.url https://exemplo.com
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Insert your JWT token in the format: Bearer {token}

// SetupEnv configures the environment based on config
func SetupEnv(cfg *config.Config) {
	// Configure logger
	logger.SetLevel(cfg.LogLevel)
	logger.Info("Starting API", map[string]interface{}{
		"port": cfg.Port,
	})

	// Set Gin mode based on log level
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// Modificação apenas para a parte relevante do arquivo main.go
// Particularmente a função SetupGCPServices

// SetupGCPServices configura e inicializa os serviços simulados do GCP
func SetupGCPServices(cfg *config.Config) (gcplogger.Logger, secrets.SecretManager, *mystorage.CloudStorage) {
	ctx := context.Background()

	// Inicializar o logger com suporte a GCP simulado
	log, err := gcplogger.NewGCPLogger(ctx, cfg.GCPProjectID, cfg.LoggingName, true)
	if err != nil {
		logger.Error("Erro ao inicializar mock do logger GCP", map[string]interface{}{
			"error": err.Error(),
		})
		// Continuar com o logger padrão em caso de erro
	} else {
		logger.Info("Mock do GCP Logger inicializado com sucesso", map[string]interface{}{
			"useCloudLogging": true,
			"mockEnabled": true,
		})
	}

	// Inicializar Secret Manager simulado
	secretManager := secrets.NewGCPSecretManager(cfg.GCPProjectID)
	logger.Info("Mock do Secret Manager inicializado", map[string]interface{}{
		"project_id": cfg.GCPProjectID,
		"mockEnabled": true,
	})

	// Inicializar Cloud Storage simulado
	cloudStorage := mystorage.NewCloudStorage(cfg.GCPStorageBucket)
	logger.Info("Mock do Cloud Storage inicializado", map[string]interface{}{
		"bucket": cfg.GCPStorageBucket,
		"mockEnabled": true,
	})

	return log, secretManager, cloudStorage
}

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config, gcpLog gcplogger.Logger, secretMgr secrets.SecretManager, cloudStorage *mystorage.CloudStorage, storageClient *storage.Client) *gin.Engine {
	// Initialize Gin router
	router := gin.New()

	// Adicionar middlewares
	router.Use(errors.RecoveryMiddleware()) // Primeiro o recovery
	router.Use(errors.ErrorMiddleware())    // Depois o tratamento de erros
	router.Use(middleware.RequestLogger())  // Por último o logger

	// Criar as instâncias dos repositórios
	itemRepo := repository.NewInMemoryItemRepository()
	userRepo := repository.NewInMemoryUserRepository()

	// Criar as instâncias dos serviços
	itemService := service.NewItemService(itemRepo)
	authService := service.NewAuthService(userRepo, cfg)

	// Criar as instâncias dos handlers
	itemHandler := handlers.NewItemHandler(itemService)
	authHandler := handlers.NewAuthHandler(authService)

	// Criar handler de demonstração do GCP (se configurado)
	var gcpDemoHandler *handlers.GCPDemoHandler
	if cfg.GCPProjectID != "" {
		gcpDemoHandler = handlers.NewGCPDemoHandler(cfg, gcpLog, secretMgr, storageClient)
		logger.Info("GCP Demo Handler inicializado", nil)
	} else {
		logger.Info("GCP Demo Handler não inicializado - GCP não configurado", nil)
	}

	// Health check route
	// @Summary Verificar status da API
	// @Description Retorna o status atual da API
	// @Tags system
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	router.GET("/health", handlers.HealthCheck)

	// Rota para testar integração GCP (endpoint original da documentação Swagger)
	// @Summary Teste de integração GCP
	// @Description Testa a integração com serviços Google Cloud Platform
	// @Tags gcp
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 503 {object} map[string]interface{}
	// @Router /api/test/gcp [get]
	router.GET("/api/test/gcp", func(c *gin.Context) {
		logger.Info("Endpoint /api/test/gcp chamado", nil)
		if gcpDemoHandler != nil {
			gcpDemoHandler.TestIntegration(c)
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "GCP integration not configured",
			})
		}
	})

	// Rota para testar integração GCP (endpoint alternativo)
	// @Summary Teste de integração GCP (Alternativo)
	// @Description Testa a integração com serviços Google Cloud Platform através do endpoint alternativo
	// @Tags gcp
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 503 {object} map[string]interface{}
	// @Router /api/test-gcp-integration [get]
	router.GET("/api/test-gcp-integration", func(c *gin.Context) {
		logger.Info("Endpoint /api/test-gcp-integration chamado", nil)
		if gcpDemoHandler != nil {
			gcpDemoHandler.TestIntegration(c)
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "GCP integration not configured",
			})
		}
	})

	// API v1 route group
	v1 := router.Group("/api/v1")
	{
		// Rotas públicas
		v1.GET("/data", itemHandler.GetData)
		v1.GET("/data/:id", itemHandler.GetDataById)

		// Rotas de autenticação
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Rotas autenticadas
			protected := auth.Group("/")
			protected.Use(middleware.JWTAuthMiddleware(cfg))
			{
				protected.GET("/profile", authHandler.Profile)
				protected.PUT("/profile", authHandler.UpdateProfile)
			}
		}

		// Rotas que exigem autenticação
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(cfg))
		{
			// Rotas básicas autenticadas
			protected.POST("/data", itemHandler.PostData)

			// Rotas que exigem papel de admin
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				// Aqui você pode adicionar rotas administrativas
				// Exemplo: admin.GET("/users", userHandler.ListUsers)
			}
		}
	}

	// Route to access Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// SetupServer configures and returns the HTTP server
func SetupServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSecs) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSecs) * time.Second,
	}
}

// StartServer starts the HTTP server and sets up graceful shutdown
func StartServer(server *http.Server, cfg *config.Config, gcpLog gcplogger.Logger, storageClient *storage.Client) {
	// Start server in a separate goroutine
	go func() {
		logger.Info("Server started", map[string]interface{}{
			"port": cfg.Port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeoutSecs)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Fechar o logger do GCP se estiver configurado
	if gcpLog != nil {
		if err := gcpLog.Close(); err != nil {
			logger.Error("Error closing GCP logger", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Fechar o cliente de armazenamento GCP se estiver configurado
	if storageClient != nil {
		if err := storageClient.Close(); err != nil {
			logger.Error("Error closing GCP storage client", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	logger.Info("Server exited gracefully", nil)
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup environment
	SetupEnv(cfg)

	// Setup GCP Services
	gcpLog, secretMgr, cloudStorage := SetupGCPServices(cfg)

	// Obter cliente de armazenamento para o manipulador GCP
	var storageClient *storage.Client
	if cloudStorage != nil {
		var err error
		storageClient, err = cloudStorage.GetClient(context.Background())
		if err != nil {
			logger.Error("Erro ao criar cliente Cloud Storage", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Setup router with GCP services
	router := SetupRouter(cfg, gcpLog, secretMgr, cloudStorage, storageClient)

	// Setup server
	server := SetupServer(cfg, router)

	// Start server with graceful shutdown
	StartServer(server, cfg, gcpLog, storageClient)
}