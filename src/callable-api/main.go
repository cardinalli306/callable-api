package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "callable-api/docs" // Para geração da documentação Swagger
)

// @title Callable API
// @version 1.0
// @description Uma API simples construída em Go usando Gin framework.
// @host localhost:8080
// @BasePath /

func main() {
	// Inicializa o router do Gin
	router := gin.Default()
	
	// Adiciona o middleware de logging
	router.Use(RequestLogger())

	// Rota de verificação de saúde
	router.GET("/health", healthCheck)

	// Grupo de rotas para API v1
	v1 := router.Group("/api/v1")
	{
		// Rotas públicas
		v1.GET("/data", getData)
		v1.GET("/data/:id", getDataById)
		
		// Rotas que requerem autenticação
		auth := v1.Group("/")
		auth.Use(TokenAuthMiddleware())
		{
			auth.POST("/data", postData)
		}
	}

	// Rota para acessar a documentação Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configuração do servidor com timeout
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Inicia o servidor em uma goroutine separada
	go func() {
		log.Printf("Servidor iniciado na porta 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando o servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}