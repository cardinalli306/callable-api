package service

import (
	"callable-api/internal/models"
	"callable-api/pkg/errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do repositório de itens
type MockItemRepository struct {
	mock.Mock
}

// Implementação dos métodos da interface repository.ItemRepository para o mock
func (m *MockItemRepository) FindAll(page, limit int) ([]models.Item, int, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Item), args.Int(1), args.Error(2)
}

func (m *MockItemRepository) FindByID(id string) (*models.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

func (m *MockItemRepository) Create(input *models.InputData) (*models.Item, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// Helper para criar um item de teste
func createTestItem() *models.Item {
	return &models.Item{
		ID:    "item123",
		Name:  "Test Item",
		Email: "test@example.com",
		Value: "100.00",
	}
}

// Helper para criar uma lista de itens de teste
func createTestItems(count int) []models.Item {
	items := make([]models.Item, count)
	for i := 0; i < count; i++ {
		items[i] = models.Item{
			ID:    "item" + strconv.Itoa(i),
			Name:  "Item " + strconv.Itoa(i),
			Email: "email" + strconv.Itoa(i) + "@example.com",
			Value: strconv.Itoa((i + 1) * 100) + ".00",
		}
	}
	return items
}

// Testes para GetItems
func TestGetItems_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Criar dados de teste
	testItems := createTestItems(3)
	totalItems := 10
	
	// Configurar comportamento do mock
	mockRepo.On("FindAll", 1, 10).Return(testItems, totalItems, nil)
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	items, total, err := itemService.GetItems(1, 10)
	
	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, testItems, items)
	assert.Equal(t, totalItems, total)
	
	mockRepo.AssertExpectations(t)
}

func TestGetItems_Error(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Configurar comportamento do mock para retornar erro
	mockRepo.On("FindAll", 1, 10).Return(nil, 0, errors.NewInternalServerError("erro de banco de dados", nil))
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	items, total, err := itemService.GetItems(1, 10)
	
	// Verificações
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, 0, total)
	
	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_SERVER", appErr.Type)
	
	mockRepo.AssertExpectations(t)
}

// Testes para GetItemByID
func TestGetItemByID_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Criar dados de teste
	testItem := createTestItem()
	
	// Configurar comportamento do mock
	mockRepo.On("FindByID", "item123").Return(testItem, nil)
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	item, err := itemService.GetItemByID("item123")
	
	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, testItem, item)
	
	mockRepo.AssertExpectations(t)
}

func TestGetItemByID_EmptyID(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método com ID vazio
	item, err := itemService.GetItemByID("")
	
	// Verificações
	assert.Error(t, err)
	assert.Nil(t, item)
	
	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "BAD_REQUEST", appErr.Type)
	
	// O repositório não deve ser chamado quando o ID é vazio
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestGetItemByID_NotFound(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Configurar mock para retornar "não encontrado"
	mockRepo.On("FindByID", "nonexistent").Return(nil, errors.NewNotFoundError("Item não encontrado", nil))
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	item, err := itemService.GetItemByID("nonexistent")
	
	// Verificações
	assert.Error(t, err)
	assert.Nil(t, item)
	
	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", appErr.Type)
	
	mockRepo.AssertExpectations(t)
}

func TestGetItemByID_RepositoryError(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Configurar mock para retornar erro de repositório
	mockRepo.On("FindByID", "error").Return(nil, errors.NewInternalServerError("Erro de banco de dados", nil))
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	item, err := itemService.GetItemByID("error")
	
	// Verificações
	assert.Error(t, err)
	assert.Nil(t, item)
	
	// O erro deve ser passado adiante sem modificação
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_SERVER", appErr.Type)
	
	mockRepo.AssertExpectations(t)
}

// Teste da função validateEmail
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"test@example", false},
		{"test.example.com", false},
		{"@example.com", true},  // Esta é uma validação simples que só verifica @ e .
		{"test@.com", true},     // Idem
		{"", false},
	}
	
	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			result := validateEmail(test.email)
			assert.Equal(t, test.valid, result, "Email: %s", test.email)
		})
	}
}

// Testes para CreateItem
func TestCreateItem_Success(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Dados de entrada
	input := &models.InputData{
		Name:  "New Item",
		Email: "new@example.com",
		Value: "150.00",
	}
	
	// Item que seria criado
	createdItem := &models.Item{
		ID:    "new123",
		Name:  "New Item",
		Email: "new@example.com",
		Value: "150.00",
	}
	
	// Configurar comportamento do mock
	mockRepo.On("Create", input).Return(createdItem, nil)
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	item, err := itemService.CreateItem(input)
	
	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, createdItem, item)
	
	mockRepo.AssertExpectations(t)
}

func TestCreateItem_ValidationError(t *testing.T) {
	// Vários casos de teste para validação
	tests := []struct {
		name          string
		input         *models.InputData
		expectedField string
	}{
		{
			name: "Empty Name",
			input: &models.InputData{
				Name:  "",
				Email: "test@example.com",
				Value: "100.00",
			},
			expectedField: "name",
		},
		{
			name: "Short Name",
			input: &models.InputData{
				Name:  "AB",
				Email: "test@example.com",
				Value: "100.00",
			},
			expectedField: "name",
		},
		{
			name: "Empty Email",
			input: &models.InputData{
				Name:  "Test Item",
				Email: "",
				Value: "100.00",
			},
			expectedField: "email",
		},
		{
			name: "Invalid Email",
			input: &models.InputData{
				Name:  "Test Item",
				Email: "invalid-email",
				Value: "100.00",
			},
			expectedField: "email",
		},
		{
			name: "Empty Value",
			input: &models.InputData{
				Name:  "Test Item",
				Email: "test@example.com",
				Value: "",
			},
			expectedField: "value",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockRepo := new(MockItemRepository)
			
			// Criar serviço com mock
			itemService := NewItemService(mockRepo)
			
			// Chamar método
			item, err := itemService.CreateItem(tt.input)
			
			// Verificações
			assert.Error(t, err)
			assert.Nil(t, item)
			
			// Verificar tipo de erro
			validationErr, ok := err.(*errors.ValidationError)
			assert.True(t, ok, "Deve ser um erro de validação")
			
			// Verificar se o campo específico foi mencionado no erro
			hasFieldError := false
			for _, fieldErr := range validationErr.FieldErrors {
				if fieldErr.Field == tt.expectedField {
					hasFieldError = true
					break
				}
			}
			assert.True(t, hasFieldError, "Deveria conter erro para o campo %s", tt.expectedField)
			
			// O repositório não deve ser chamado quando há erros de validação
			mockRepo.AssertNotCalled(t, "Create")
		})
	}
}

func TestCreateItem_RepositoryError(t *testing.T) {
	// Configurar mock
	mockRepo := new(MockItemRepository)
	
	// Dados de entrada válidos
	input := &models.InputData{
		Name:  "New Item",
		Email: "new@example.com",
		Value: "150.00",
	}
	
	// Configurar mock para retornar erro
	mockRepo.On("Create", input).Return(nil, errors.NewInternalServerError("erro de banco de dados", nil))
	
	// Criar serviço com mock
	itemService := NewItemService(mockRepo)
	
	// Chamar método
	item, err := itemService.CreateItem(input)
	
	// Verificações
	assert.Error(t, err)
	assert.Nil(t, item)
	
	// Verificar tipo de erro
	appErr, ok := err.(*errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_SERVER", appErr.Type)
	
	mockRepo.AssertExpectations(t)
}