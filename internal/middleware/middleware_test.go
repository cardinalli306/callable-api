package middleware_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"callable-api/internal/middleware"
	"callable-api/internal/models"
	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/logger"
)

// Mock para SecretManager
type MockSecretManager struct {
	mock.Mock
}

func (m *MockSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	args := m.Called(ctx, secretName)
	return args.String(0), args.Error(1)
}

func (m *MockSecretManager) GetSecretWithCache(ctx context.Context, secretName string, ttl time.Duration) (string, error) {
	args := m.Called(ctx, secretName, ttl)
	return args.String(0), args.Error(1)
}

// Mock para Logger
type MockLogger struct {
	mock.Mock
}

// Debug registra mensagens de nível debug. Implementação simulada para testes.
func (m *MockLogger) Debug(msg string, data ...map[string]interface{}) {
	// Implementação vazia pois não precisamos processar logs durante os testes
}

// Info registra mensagens informativas. Implementação simulada para testes.
func (m *MockLogger) Info(msg string, data ...map[string]interface{}) {
	// Implementação vazia pois não precisamos processar logs durante os testes
}

// Warn registra avisos. Implementação simulada para testes.
func (m *MockLogger) Warn(msg string, data ...map[string]interface{}) {
	// Implementação vazia pois não precisamos processar logs durante os testes
}

// Error registra erros. Implementação simulada para testes.
func (m *MockLogger) Error(msg string, err error, data ...map[string]interface{}) {
	// Implementação vazia pois não precisamos processar logs durante os testes
}

// Fatal registra erros fatais. Implementação simulada para testes.
func (m *MockLogger) Fatal(msg string, err error, data ...map[string]interface{}) {
	// Implementação vazia pois não precisamos processar logs durante os testes
}

// Close finaliza o logger e libera recursos. Retorna o valor simulado configurado em testes.
func (m *MockLogger) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Constante para a chave secreta usada nos testes
const testSecretKey = "test-secret-key"

// testJWTAuthMiddleware é um adaptador para teste que permite usar
// um SecretProvider com o JWTAuthMiddleware
func testJWTAuthMiddleware(_ *auth.SecretProvider) gin.HandlerFunc {
    // Cria uma configuração mock para o teste
    cfg := &config.Config{
        JWTSecret: testSecretKey,
        JWTConfig: config.JWTConfig{
            SecretKey: testSecretKey,
        },
    }

    // Retorna o middleware real com a config
    return middleware.JWTAuthMiddleware(cfg)
}

func TestJWTAuthMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	// Criar as dependências necessárias
	mockSecretMgr := new(MockSecretManager)
	mockLogger := new(MockLogger)

	// Mock o método Close para o logger
	mockLogger.On("Close").Return(nil)

	cfg := &config.Config{
		JWTConfig: config.JWTConfig{
			SecretKey: "test-secret-key",
		},
		JWTSecret: "test-secret-key", // Adicionado para compatibilidade
	}

	// Configurar o mock para retornar a chave secreta
	testSecret := "test-secret-key"
	mockSecretMgr.On("GetSecret", mock.Anything, mock.Anything).Return(testSecret, nil)
	mockSecretMgr.On("GetSecretWithCache", mock.Anything, mock.Anything, mock.Anything).Return(testSecret, nil)

	// Criar o provedor JWT
	var loggerInterface logger.Logger = mockLogger // Verificação de tipo
	jwtProvider := auth.NewSecretProvider(cfg, mockSecretMgr, loggerInterface)

	t.Run("Sem token de autorização retorna erro", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		// Usar o adaptador para o middleware
		testJWTAuthMiddleware(jwtProvider)(c)

		// Verifica resposta
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var apiErr models.APIError
		err := json.Unmarshal(w.Body.Bytes(), &apiErr)
		assert.NoError(t, err)
		assert.Equal(t, "error", apiErr.Status)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
		assert.Contains(t, apiErr.Message, "Token de autenticação não fornecido")
	})

	t.Run("Token inválido retorna erro", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set("Authorization", "Bearer token-invalido")

		// Usar o adaptador para o middleware
		testJWTAuthMiddleware(jwtProvider)(c)

		// Verifica resposta
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var apiErr models.APIError
		err := json.Unmarshal(w.Body.Bytes(), &apiErr)
		assert.NoError(t, err)
		assert.Equal(t, "error", apiErr.Status)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
		assert.Contains(t, apiErr.Message, "inválido")
	})

	t.Run("Token válido permite acesso", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Cria um router e adiciona o middleware usando o adaptador
		router := gin.New()
		router.Use(testJWTAuthMiddleware(jwtProvider))

		// Adiciona uma rota que só será alcançada se o middleware permitir
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Acesso autorizado",
			})
		})

		// Cria um token válido usando a mesma chave secreta
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "user-123",      // Alterado para corresponder à estrutura Claims
			"email":   "test@example.com",
			"name":    "Test User",
			"role":    "user",
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		tokenString, err := token.SignedString([]byte(testSecret))
		assert.NoError(t, err)

		// Faz uma requisição com o token
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		router.ServeHTTP(w, req)

		// Verifica se o acesso foi permitido
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Acesso autorizado", response.Message)
	})

	t.Run("Token expirado retorna erro", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Cria um token expirado
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "user-123",      // Alterado para corresponder à estrutura Claims
			"email":   "test@example.com",
			"name":    "Test User",
			"role":    "user",
			"exp":     time.Now().Add(-time.Hour).Unix(), // Expirado há 1 hora
		})
		tokenString, err := token.SignedString([]byte(testSecret))
		assert.NoError(t, err)

		// Configura requisição com token expirado
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)

		// Usar o adaptador para o middleware
		testJWTAuthMiddleware(jwtProvider)(c)

		// Verifica resposta
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var apiErr models.APIError
		err = json.Unmarshal(w.Body.Bytes(), &apiErr)
		assert.NoError(t, err)
		assert.Equal(t, "error", apiErr.Status)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
		assert.Contains(t, apiErr.Message, "inválido")
	})
}

func TestValidationErrorMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("Erro de validação é formatado corretamente", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(middleware.ValidationErrorMiddleware())

		// Adiciona rota de teste que gera um erro de validação
		router.POST("/test", func(c *gin.Context) {
			// Adiciona um erro diretamente no contexto Gin
			c.Error(fmt.Errorf("erro de validação"))
			c.AbortWithStatusJSON(http.StatusBadRequest, models.APIError{
				Status:      "error",
				Code:        http.StatusBadRequest,
				Message:     "Dados inválidos",
				FieldErrors: map[string]string{
					"name":  "O nome é obrigatório",
					"email": "Email em formato inválido",
				},
			})
		})

		// Faz uma requisição
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var apiErr models.APIError
		err := json.Unmarshal(w.Body.Bytes(), &apiErr)
		assert.NoError(t, err)
		assert.Equal(t, "error", apiErr.Status)
		assert.Equal(t, http.StatusBadRequest, apiErr.Code)
		assert.Equal(t, "Dados inválidos", apiErr.Message)
		assert.Equal(t, "O nome é obrigatório", apiErr.FieldErrors["name"])
		assert.Equal(t, "Email em formato inválido", apiErr.FieldErrors["email"])
	})

	t.Run("Middleware passa adiante requisições sem erro", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(middleware.ValidationErrorMiddleware())

		// Adiciona rota de teste sem erro
		router.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Operação bem-sucedida",
				Data:    map[string]string{"result": "ok"},
			})
		})

		// Faz uma requisição
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/success", nil)
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Operação bem-sucedida", response.Message)
	})
}

func TestLoggerMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("Logger middleware processa requisições corretamente", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(middleware.LoggerMiddleware())

		// Adiciona rota de teste
		router.GET("/log-test", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Requisição registrada",
			})
		})

		// Faz uma requisição
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/log-test", nil)
		router.ServeHTTP(w, req)

		// Verifica resposta - como o logger apenas registra, só verificamos se
		// a requisição foi processada normalmente
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Requisição registrada", response.Message)
	})
}

func TestCORSMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("CORS middleware adiciona headers corretos", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(middleware.CORSMiddleware())

		// Adiciona rota de teste
		router.GET("/cors-test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Faz uma requisição com Origin
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/cors-test", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		// Verifica se os headers CORS foram adicionados
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})

	t.Run("CORS middleware responde a requisição OPTIONS", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(middleware.CORSMiddleware())

		// Faz uma requisição OPTIONS
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/any-path", nil)
		router.ServeHTTP(w, req)

		// Verifica se a resposta OPTIONS está correta
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "OPTIONS")
	})
}