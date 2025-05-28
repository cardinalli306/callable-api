package service

import (
	"callable-api/internal/models"
	"callable-api/internal/repository"
	"callable-api/pkg/errors"
	"callable-api/pkg/logger"
	"strings"
)

// ItemService gerencia a lógica de negócios relacionada a itens
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService cria uma nova instância do ItemService
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

// GetItems retorna uma lista paginada de itens
func (s *ItemService) GetItems(page, limit int) ([]models.Item, int, error) {
	logger.Info("Buscando lista de itens", map[string]interface{}{
		"page":  page,
		"limit": limit,
	})
	
	items, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, errors.NewInternalServerError("Falha ao buscar itens", err)
	}
	
	return items, total, nil
}

// GetItemByID retorna um item específico pelo ID
func (s *ItemService) GetItemByID(id string) (*models.Item, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("ID não fornecido", nil)
	}
	
	logger.Info("Buscando item por ID", map[string]interface{}{
		"id": id,
	})
	
	item, err := s.repo.FindByID(id)
	if err != nil {
		// O repositório já retorna um erro NotFound se não encontrar
		return nil, err
	}
	
	return item, nil
}

// validateEmail realiza uma validação simples de email
func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// CreateItem cria um novo item
func (s *ItemService) CreateItem(input *models.InputData) (*models.Item, error) {
	// Validar input usando o novo sistema de erros de validação
	validationErr := errors.NewValidationError("Dados de entrada inválidos")
	validInputs := true
	
	if input.Name == "" {
		validationErr.AddFieldError("name", "Nome é obrigatório")
		validInputs = false
	} else if len(input.Name) < 3 {
		validationErr.AddFieldError("name", "Nome deve ter pelo menos 3 caracteres")
		validInputs = false
	}
	
	if input.Email == "" {
		validationErr.AddFieldError("email", "Email é obrigatório")
		validInputs = false
	} else if !validateEmail(input.Email) {
		validationErr.AddFieldError("email", "Email inválido")
		validInputs = false
	}
	
	if input.Value == "" {
		validationErr.AddFieldError("value", "Valor é obrigatório")
		validInputs = false
	}
	
	if !validInputs {
		return nil, validationErr
	}
	
	logger.Info("Criando novo item", map[string]interface{}{
		"name": input.Name,
		"email": input.Email,
	})
	
	item, err := s.repo.Create(input)
	if err != nil {
		return nil, errors.NewInternalServerError("Falha ao criar item", err)
	}
	
	return item, nil
}