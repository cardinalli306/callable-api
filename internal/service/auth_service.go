package service

import (
	"callable-api/internal/models"
	"callable-api/internal/repository"
	"callable-api/pkg/auth"
	"callable-api/pkg/config"
	"callable-api/pkg/errors"
	"callable-api/pkg/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService gerencia autenticação e usuários
type AuthService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

// NewAuthService cria uma nova instância do AuthService
func NewAuthService(repo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

// Register registra um novo usuário
func (s *AuthService) Register(input *models.RegisterUserInput) (*models.UserResponse, error) {
	// Validação adicional pode ser feita aqui
	validationErr := errors.NewValidationError("Dados de entrada inválidos")
	validInputs := true

	if len(input.Password) < 6 {
		validationErr.AddFieldError("password", "Senha deve ter pelo menos 6 caracteres")
		validInputs = false
	}

	if !validInputs {
		return nil, validationErr
	}

	// Verificar se o email já está em uso
	_, err := s.repo.FindByEmail(input.Email)
	if err == nil {
		return nil, errors.NewConflictError("Email já está em uso", nil)
	}

	// Só continuar se o erro for "não encontrado"
	if _, ok := err.(*errors.AppError); !ok || err.(*errors.AppError).Type != "NOT_FOUND" {
		return nil, errors.NewInternalServerError("Erro ao verificar disponibilidade do email", err)
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalServerError("Erro ao processar senha", err)
	}

	// Criar usuário
	user := &models.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: string(hashedPassword),
		Role:     "user", // Papel padrão
	}

	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, errors.NewInternalServerError("Erro ao criar usuário", err)
	}

	logger.Info("Usuário registrado com sucesso", map[string]interface{}{
		"userId": createdUser.ID,
		"email":  createdUser.Email,
	})

	return &models.UserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Name:      createdUser.Name,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

// Login autentica um usuário e retorna tokens JWT
func (s *AuthService) Login(input *models.LoginInput) (*models.TokenPair, *models.UserResponse, error) {
	// Autenticar usuário
	user, err := s.repo.Authenticate(input.Email, input.Password)
	if err != nil {
		return nil, nil, err // O repositório já retorna o erro adequado
	}

	// Gerar tokens
	tokenPair, err := auth.GenerateTokenPair(user, s.cfg)
	if err != nil {
		return nil, nil, errors.NewInternalServerError("Erro ao gerar tokens", err)
	}

	logger.Info("Login de usuário bem-sucedido", map[string]interface{}{
		"userId": user.ID,
		"email":  user.Email,
	})

	return tokenPair, &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

// RefreshToken atualiza os tokens JWT usando um token de atualização
func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenPair, error) {
	// Validar o token de atualização
	claims, err := auth.ValidateToken(refreshToken, true, s.cfg)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Token de atualização inválido", err)
	}

	// Buscar o usuário
	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Usuário não encontrado", err)
	}

	// Gerar novos tokens
	tokenPair, err := auth.GenerateTokenPair(user, s.cfg)
	if err != nil {
		return nil, errors.NewInternalServerError("Erro ao gerar tokens", err)
	}

	logger.Info("Tokens atualizados com sucesso", map[string]interface{}{
		"userId": user.ID,
		"email":  user.Email,
	})

	return tokenPair, nil
}

// GetUserProfile retorna o perfil do usuário
func (s *AuthService) GetUserProfile(userID string) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err // O repositório já retorna o erro adequado
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateUserProfile atualiza o perfil do usuário
func (s *AuthService) UpdateUserProfile(userID string, name string) (*models.UserResponse, error) {
	// Buscar usuário atual
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	user.Name = name
	user.UpdatedAt = time.Now()

	// Salvar usuário
	updatedUser, err := s.repo.Update(user)
	if err != nil {
		return nil, errors.NewInternalServerError("Erro ao atualizar perfil", err)
	}

	return &models.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Name:      updatedUser.Name,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.CreatedAt,
	}, nil
}