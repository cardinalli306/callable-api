package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/models"
	
)



func TestTokenAuthMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("Sem token de autorização retorna erro", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(TokenAuthMiddleware("test-token"))

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Acesso autorizado",
			})
		})

		// Faz uma requisição sem token
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Authorization token required", response.Message)
	})

	t.Run("Token inválido retorna erro", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(TokenAuthMiddleware("test-token"))

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Acesso autorizado",
			})
		})

		// Faz uma requisição com token inválido
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Invalid or empty token", response.Message)
	})

	t.Run("Token válido permite acesso", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(TokenAuthMiddleware("test-token"))

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Acesso autorizado",
			})
		})

		// Faz uma requisição com token válido
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Acesso autorizado", response.Message)
	})

	t.Run("Token direto sem Bearer também funciona", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(TokenAuthMiddleware("test-token"))

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Acesso autorizado",
			})
		})

		// Faz uma requisição com token válido sem prefixo Bearer
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "test-token")
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Acesso autorizado", response.Message)
	})
}

func TestValidationErrorMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("Middleware passa adiante requisições sem erro", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(ValidationErrorMiddleware())

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

	t.Run("Erros de validação são tratados", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(ValidationErrorMiddleware())

		// Adiciona rota que produz erro
		router.GET("/validation-error", func(c *gin.Context) {
			c.Error(gin.Error{
				Err:  errors.New("campo obrigatório"),
				Type: gin.ErrorTypePrivate,
				Meta: nil,
			})
			c.Next() // Deixe o middleware tratar o erro
		})

		// Faz uma requisição
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/validation-error", nil)
		router.ServeHTTP(w, req)

		// Verifica resposta
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Validation error")
	})
}

func TestLoggerMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("Logger middleware processa requisições corretamente", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(LoggerMiddleware())

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

		// Verifica resposta - apenas verifica se a requisição foi processada
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCORSMiddleware(t *testing.T) {
	// Configuração para testes
	gin.SetMode(gin.TestMode)

	t.Run("CORS middleware adiciona headers corretos", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(CORSMiddleware())

		// Adiciona rota de teste
		router.GET("/cors-test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Faz uma requisição com Origin
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/cors-test", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		// Verifica headers CORS
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})

	t.Run("CORS middleware responde a requisição OPTIONS", func(t *testing.T) {
		// Configura router com middleware
		router := gin.New()
		router.Use(CORSMiddleware())

		// Faz uma requisição OPTIONS
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/any-path", nil)
		router.ServeHTTP(w, req)

		// Verifica resposta OPTIONS
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRequestLogger(t *testing.T) {
	// Teste simples para verificar se RequestLogger é apenas um alias
	// para LoggerMiddleware
	gin.SetMode(gin.TestMode)
	
	// Configura dois routers, um com cada middleware
	routerA := gin.New()
	routerA.Use(LoggerMiddleware())
	routerA.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	
	routerB := gin.New()
	routerB.Use(RequestLogger())
	routerB.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	
	// Faz requisições para ambos
	wA := httptest.NewRecorder()
	reqA := httptest.NewRequest(http.MethodGet, "/test", nil)
	routerA.ServeHTTP(wA, reqA)
	
	wB := httptest.NewRecorder()
	reqB := httptest.NewRequest(http.MethodGet, "/test", nil)
	routerB.ServeHTTP(wB, reqB)
	
	// Verifica se ambos retornam o mesmo resultado
	assert.Equal(t, wA.Code, wB.Code)
	assert.Equal(t, wA.Body.String(), wB.Body.String())
}