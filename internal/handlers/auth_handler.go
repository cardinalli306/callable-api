package handlers

import (
	"callable-api/internal/models"
	"callable-api/internal/service"
	"callable-api/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler processa requisições relacionadas a autenticação
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register registra um novo usuário
// @Summary Registrar um novo usuário
// @Description Cria uma nova conta de usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterUserInput true "Dados de registro"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.APIError
// @Failure 409 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		validationErr := errors.NewValidationError("Dados de registro inválidos")
		validationErr.AddFieldError("request", "Formato de dados inválido")
		errors.HandleErrors(c, validationErr)
		return
	}

	user, err := h.service.Register(&input)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login autentica um usuário
// @Summary Login de usuário
// @Description Autentica o usuário e retorna os tokens JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginInput true "Credenciais de login"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		validationErr := errors.NewValidationError("Dados de login inválidos")
		validationErr.AddFieldError("request", "Formato de dados inválido")
		errors.HandleErrors(c, validationErr)
		return
	}

	tokens, user, err := h.service.Login(&input)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"user":   user,
	})
}

// RefreshToken renova os tokens JWT
// @Summary Atualizar tokens
// @Description Renova os tokens JWT usando um token de atualização
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Token de atualização"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		validationErr := errors.NewValidationError("Dados inválidos")
		validationErr.AddFieldError("refresh_token", "Token de atualização é obrigatório")
		errors.HandleErrors(c, validationErr)
		return
	}

	tokens, err := h.service.RefreshToken(request.RefreshToken)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Profile retorna o perfil do usuário autenticado
// @Summary Perfil do usuário
// @Description Retorna os dados do perfil do usuário autenticado
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		err := errors.NewUnauthorizedError("ID de usuário inválido", nil)
		errors.HandleErrors(c, err)
		return
	}

	profile, err := h.service.GetUserProfile(userIDStr)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile atualiza o perfil do usuário
// @Summary Atualizar perfil
// @Description Atualiza os dados do perfil do usuário autenticado
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]string true "Dados para atualização do perfil"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		err := errors.NewUnauthorizedError("ID de usuário inválido", nil)
		errors.HandleErrors(c, err)
		return
	}

	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		validationErr := errors.NewValidationError("Dados inválidos")
		validationErr.AddFieldError("name", "Nome é obrigatório")
		errors.HandleErrors(c, validationErr)
		return
	}

	profile, err := h.service.UpdateUserProfile(userIDStr, request.Name)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}