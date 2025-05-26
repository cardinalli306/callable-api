package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "callable-api/docs" // Isso importa o pacote apenas pelos seus efeitos colaterais // For Swagger documentation generation
	"callable-api/internal/handlers"
	"callable-api/internal/middleware"
	"callable-api/pkg/config"
	"callable-api/pkg/logger"
)

// @title Callable API
// @version 1.0
// @description Uma API robusta construída em Go usando o framework Gin, oferecendo endpoints para gerenciamento de dados com validação completa.
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

func main() {
	// Load configuration
	cfg := config.Load()

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

	// Initialize Gin router
	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	// Health check route
	router.GET("/health", handlers.HealthCheck)

	// API v1 route group
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/data", handlers.GetData)
		v1.GET("/data/:id", handlers.GetDataById)

		// Routes requiring authentication
		auth := v1.Group("/")
		auth.Use(middleware.TokenAuthMiddleware(cfg.DemoApiToken))
		{
			auth.POST("/data", handlers.PostData)
		}
	}

	// Route to access Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Server configuration with timeout
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSecs) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSecs) * time.Second,
	}

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

	logger.Info("Server exited gracefully", nil)
}