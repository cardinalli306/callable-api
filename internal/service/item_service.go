package service

import (
	"callable-api/internal/models"
	"callable-api/internal/repository"
	"callable-api/pkg/errors"
	"callable-api/pkg/logger"
	"context"
	"strings"
	"time"
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
// Modificado para aceitar um contexto para controle de cancelamento/timeout
func (s *ItemService) CreateItem(ctx context.Context, input *models.InputData) (*models.Item, error) {
	// Verificar se o contexto já foi cancelado
	if ctx.Err() != nil {
		return nil, errors.NewInternalServerError("Operação cancelada", ctx.Err())
	}
	
	// Log com requestID se disponível no contexto
	var logData map[string]interface{}
	if requestID, ok := ctx.Value("request_id").(string); ok {
		logData = map[string]interface{}{
			"request_id": requestID,
			"name":       input.Name,
			"email":      input.Email,
		}
	} else {
		logData = map[string]interface{}{
			"name":  input.Name,
			"email": input.Email,
		}
	}
	
	logger.Info("Validando dados para criação de item", logData)
	
	// Validar input usando o sistema de erros de validação
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
	
	// Verificar novamente o contexto após validação
	if ctx.Err() != nil {
		return nil, errors.NewInternalServerError("Operação cancelada após validação", ctx.Err())
	}
	
	logger.Info("Criando novo item", logData)
	
	// Simulando uma operação de longa duração (remover em produção)
	// Isso é apenas para testar o comportamento do timeout
	if input.Value == "demorado" {
		logger.Info("Simulando operação de longa duração", logData)
		
		// Loop para simular processamento e verificar contexto periodicamente
		for i := 0; i < 20; i++ {
			select {
			case <-ctx.Done():
				logger.Warn("Contexto cancelado durante processamento", logData)
				return nil, errors.NewInternalServerError("Operação cancelada durante processamento", ctx.Err())
			case <-time.After(500 * time.Millisecond):
				// Continua processamento
			}
		}
	}
	
	// Pode ser necessário modificar o repositório para aceitar contexto também
	// Por enquanto, estamos apenas passando o input
	item, err := s.repo.Create(input)
	if err != nil {
		logger.Error("Falha ao criar item no repositório", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.NewInternalServerError("Falha ao criar item", err)
	}
	
	logger.Info("Item criado com sucesso", map[string]interface{}{
		"item_id": item.ID,
	})
	
	return item, nil
}

// Você pode adicionar métodos adicionais conforme necessário