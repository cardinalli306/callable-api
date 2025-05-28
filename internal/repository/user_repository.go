package repository

import (
	"callable-api/internal/models"
	"callable-api/pkg/errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	userNotFoundMessage = "Usuário não encontrado" // Definição da constante
)

// UserRepository define as operações do repositório de usuários
type UserRepository interface {
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	List(page, limit int) ([]models.User, int, error)
	Delete(id string) error
	Authenticate(email, password string) (*models.User, error)
}

// InMemoryUserRepository implementa um repositório de usuários em memória
type InMemoryUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository cria um novo repositório de usuários em memória
func NewInMemoryUserRepository() *InMemoryUserRepository {
	// Criar com alguns usuários de exemplo
	repo := &InMemoryUserRepository{
		users: make(map[string]*models.User),
	}

	// Criar um usuário administrativo padrão
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminID := uuid.New().String()
	repo.users[adminID] = &models.User{
		ID:        adminID,
		Email:     "admin@example.com",
		Name:      "Admin User",
		Password:  string(adminPassword),
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Criar um usuário normal de exemplo
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	userID := uuid.New().String()
	repo.users[userID] = &models.User{
		ID:        userID,
		Email:     "user@example.com",
		Name:      "Regular User",
		Password:  string(userPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return repo
}

// FindByID busca um usuário pelo ID
func (r *InMemoryUserRepository) FindByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, errors.NewNotFoundError(userNotFoundMessage, nil) // Usando a constante
}

// FindByEmail busca um usuário pelo email
func (r *InMemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.NewNotFoundError(userNotFoundMessage, nil) // Usando a constante
}

// Create cria um novo usuário
func (r *InMemoryUserRepository) Create(user *models.User) (*models.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Verificar se o email já está em uso
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return nil, errors.NewConflictError("Email já está em uso", nil)
		}
	}

	// Gerar um novo ID se não for fornecido
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Definir timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Adicionar o usuário ao repositório
	r.users[user.ID] = user
	return user, nil
}

// Update atualiza um usuário existente
func (r *InMemoryUserRepository) Update(user *models.User) (*models.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return nil, errors.NewNotFoundError(userNotFoundMessage, nil) // Usando a constante
	}

	// Verificar se o email está sendo alterado e se o novo email já está em uso
	oldUser := r.users[user.ID]
	if user.Email != oldUser.Email {
		for _, existingUser := range r.users {
			if existingUser.Email == user.Email && existingUser.ID != user.ID {
				return nil, errors.NewConflictError("Email já está em uso", nil)
			}
		}
	}

	// Atualizar timestamp
	user.UpdatedAt = time.Now()
	user.CreatedAt = oldUser.CreatedAt // Preservar data de criação

	// Atualizar o usuário no repositório
	r.users[user.ID] = user
	return user, nil
}

// List retorna uma lista paginada de usuários
func (r *InMemoryUserRepository) List(page, limit int) ([]models.User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	total := len(r.users)

	// Converter o mapa para uma slice
	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}

	// Aplicar paginação
	end := offset + limit
	if end > total {
		end = total
	}

	if offset >= total {
		return []models.User{}, total, nil
	}

	return users[offset:end], total, nil
}

// Delete remove um usuário pelo ID
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.NewNotFoundError(userNotFoundMessage, nil) // Usando a constante
	}

	delete(r.users, id)
	return nil
}

// Authenticate verifica as credenciais do usuário e retorna o usuário se válido
func (r *InMemoryUserRepository) Authenticate(email, password string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var user *models.User
	for _, u := range r.users {
		if u.Email == email {
			user = u
			break
		}
	}

	if user == nil {
		return nil, errors.NewUnauthorizedError("Credenciais inválidas", nil)
	}

	// Verificar a senha
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.NewUnauthorizedError("Credenciais inválidas", nil)
	}

	return user, nil
}
