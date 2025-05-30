// auth_service_test.go
package service

import (
	"callable-api/internal/models"
	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock do repositório de usuário
type MockUserRepository struct {
	mock.Mock
}

// List implements repository.UserRepository.
func (m *MockUserRepository) List(page int, limit int) ([]models.User, int, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Authenticate(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Configurações para testes
func getTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:                "test-secret",
		JWTExpirationMinutes:     15,
		JWTRefreshExpirationDays: 7,
	}
}

// Helper para criar usuário de teste
func createTestUser() *models.User {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	now := time.Now()

	return &models.User{
		ID:        "user123",
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  string(hashedPwd),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Testes para Register
func TestRegister_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Configurar comportamento do mock
	mockRepo.On("FindByEmail", "new@example.com").Return(nil,
		errors.NewNotFoundError("Usuário não encontrado", nil))

	// Mock da criação do usuário
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(
		&models.User{
			ID:        "new123",
			Email:     "new@example.com",
			Name:      "New User",
			Password:  "hashedpassword", // Não importa o valor exato aqui
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Dados de entrada para o registro
	input := &models.RegisterUserInput{
		Email:    "new@example.com",
		Name:     "New User",
		Password: "password123",
	}

	// Chamar método
	userResponse, err := authService.Register(input)

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, userResponse)
	assert.Equal(t, "new123", userResponse.ID)
	assert.Equal(t, "new@example.com", userResponse.Email)
	assert.Equal(t, "New User", userResponse.Name)
	assert.Equal(t, "user", userResponse.Role)

	// Verificar que os métodos mock foram chamados conforme esperado
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Email já existe no sistema
	existingUser := createTestUser()
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Dados de entrada para o registro
	input := &models.RegisterUserInput{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	// Chamar método
	userResponse, err := authService.Register(input)

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, userResponse)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "CONFLICT", appErr.Type)

	// Verificar que os métodos mock foram chamados conforme esperado
	mockRepo.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
    // Configurar mock
    mockRepo := new(MockUserRepository)
    
    // Criar serviço com mock
    authService := NewAuthService(mockRepo, getTestConfig())
    
    // Dados de entrada inválidos (senha muito curta)
    input := &models.RegisterUserInput{
        Email:    "test@example.com",
        Name:     "Test User",
        Password: "12345", // Menos de 6 caracteres
    }
    
    // Chamar método
    userResponse, err := authService.Register(input)
    
    // Verificações básicas
    assert.Error(t, err)
    assert.Nil(t, userResponse)
    
    // Verificar o tipo específico do erro
    validationErr, ok := err.(*errors.ValidationError)
    assert.True(t, ok, "O erro deveria ser do tipo *errors.ValidationError")
    
    // Se quiser verificar a mensagem específica
    if ok {
        assert.Equal(t, "Dados de entrada inválidos", validationErr.Error())
    }
    
    // Mock não deve ser chamado para verificar email quando há erro de validação
    mockRepo.AssertNotCalled(t, "FindByEmail")
}
func TestRegister_RepositoryError(t *testing.T) {
    // Configurar mock
    mockRepo := new(MockUserRepository)
    
    // Configurar comportamento do mock para retornar erro no FindByEmail
    mockRepo.On("FindByEmail", "error@example.com").Return(nil, 
        errors.NewInternalServerError("Erro de banco de dados", nil))
    
    // Criar serviço com mock
    authService := NewAuthService(mockRepo, getTestConfig())
    
    // Dados de entrada
    input := &models.RegisterUserInput{
        Email:    "error@example.com",
        Name:     "Error User",
        Password: "password123",
    }
    
    // Chamar método
    userResponse, err := authService.Register(input)
    
    // Verificações
    assert.Error(t, err)
    assert.Nil(t, userResponse)
    
    // Verificar tipo de erro
    appErr, ok := err.(*errors.AppError)
    assert.True(t, ok)
    assert.Equal(t, "INTERNAL_SERVER", appErr.Type) // Corrigido para o valor real
    
    mockRepo.AssertExpectations(t)
}

// Testes para Login
func TestLogin_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Usuário existente
	user := createTestUser()

	// Mock de autenticação bem-sucedida
	mockRepo.On("Authenticate", "test@example.com", "password123").Return(user, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Dados de entrada para login
	input := &models.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Chamar método
	tokenPair, userResponse, err := authService.Login(input)

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotNil(t, userResponse)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, user.ID, userResponse.ID)
	assert.Equal(t, user.Email, userResponse.Email)

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Mock de autenticação falha
	mockRepo.On("Authenticate", "test@example.com", "wrongpassword").Return(nil,
		errors.NewUnauthorizedError("Credenciais inválidas", nil))

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Dados de entrada para login
	input := &models.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Chamar método
	tokenPair, userResponse, err := authService.Login(input)

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Nil(t, userResponse)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "UNAUTHORIZED", appErr.Type)

	mockRepo.AssertExpectations(t)
}

// Testes para GetUserProfile
func TestGetUserProfile_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Usuário existente
	user := createTestUser()

	// Mock de busca por ID
	mockRepo.On("FindByID", "user123").Return(user, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Chamar método
	userResponse, err := authService.GetUserProfile("user123")

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, userResponse)
	assert.Equal(t, user.ID, userResponse.ID)
	assert.Equal(t, user.Email, userResponse.Email)
	assert.Equal(t, user.Name, userResponse.Name)
	assert.Equal(t, user.Role, userResponse.Role)

	mockRepo.AssertExpectations(t)
}

func TestGetUserProfile_UserNotFound(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Mock de usuário não encontrado
	mockRepo.On("FindByID", "nonexistent").Return(nil,
		errors.NewNotFoundError("Usuário não encontrado", nil))

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Chamar método
	userResponse, err := authService.GetUserProfile("nonexistent")

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, userResponse)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", appErr.Type)

	mockRepo.AssertExpectations(t)
}

// Testes para UpdateUserProfile
func TestUpdateUserProfile_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Usuário existente
	user := createTestUser()

	// Mock de busca por ID
	mockRepo.On("FindByID", "user123").Return(user, nil)

	// Cópia do usuário com nome atualizado
	updatedUser := *user
	updatedUser.Name = "Updated Name"

	// Mock de atualização
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(&updatedUser, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Chamar método
	userResponse, err := authService.UpdateUserProfile("user123", "Updated Name")

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, userResponse)
	assert.Equal(t, "Updated Name", userResponse.Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserProfile_UserNotFound(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Mock de usuário não encontrado
	mockRepo.On("FindByID", "nonexistent").Return(nil,
		errors.NewNotFoundError("Usuário não encontrado", nil))

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Chamar método
	userResponse, err := authService.UpdateUserProfile("nonexistent", "New Name")

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, userResponse)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", appErr.Type)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserProfile_UpdateError(t *testing.T) {
    // Configurar mock
    mockRepo := new(MockUserRepository)
    
    // Usuário existente
    user := createTestUser()
    
    // Mock de busca por ID
    mockRepo.On("FindByID", "user123").Return(user, nil)
    
    // Mock de erro na atualização
    mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil, 
        errors.NewInternalServerError("Erro ao atualizar", nil))
    
    // Criar serviço com mock
    authService := NewAuthService(mockRepo, getTestConfig())
    
    // Chamar método
    userResponse, err := authService.UpdateUserProfile("user123", "Updated Name")
    
    // Verificações
    assert.Error(t, err)
    assert.Nil(t, userResponse)
    
    // Verificar tipo de erro
    appErr, ok := err.(*errors.AppError)
    assert.True(t, ok)
    assert.Equal(t, "INTERNAL_SERVER", appErr.Type) // Corrigido para o valor real
    
    mockRepo.AssertExpectations(t)
}

// Testes para RefreshToken
func TestRefreshToken_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)
	cfg := getTestConfig()

	// Usuário existente
	user := createTestUser()

	// Gerar um token de teste
	tokenPair, _ := auth.GenerateTokenPair(user, cfg)

	// Mock de busca por ID (será chamado após validação do token)
	mockRepo.On("FindByID", "user123").Return(user, nil)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, cfg)

	// Chamar método
	newTokenPair, err := authService.RefreshToken(tokenPair.RefreshToken)

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, newTokenPair)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.NotEmpty(t, newTokenPair.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockUserRepository)

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, getTestConfig())

	// Token inválido
	invalidToken := "invalid.token.string"

	// Chamar método
	newTokenPair, err := authService.RefreshToken(invalidToken)

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, newTokenPair)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "UNAUTHORIZED", appErr.Type)

	// Não deve chamar o repositório
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	// Este teste é um pouco mais complexo pois precisa de um token válido
	// mas para um usuário que não existe mais

	// Configurar mock
	mockRepo := new(MockUserRepository)
	cfg := getTestConfig()

	// Usuário que será "removido" depois
	tempUser := &models.User{
		ID:        "deleted123",
		Email:     "deleted@example.com",
		Name:      "Deleted User",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	// Gerar token para o usuário temporário
	tokenPair, _ := auth.GenerateTokenPair(tempUser, cfg)

	// Mock - usuário não existe mais quando tentamos buscá-lo
	mockRepo.On("FindByID", "deleted123").Return(nil,
		errors.NewNotFoundError("Usuário não encontrado", nil))

	// Criar serviço com mock
	authService := NewAuthService(mockRepo, cfg)

	// Chamar método com o refresh token
	newTokenPair, err := authService.RefreshToken(tokenPair.RefreshToken)

	// Verificações
	assert.Error(t, err)
	assert.Nil(t, newTokenPair)

	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "UNAUTHORIZED", appErr.Type)

	mockRepo.AssertExpectations(t)
}
