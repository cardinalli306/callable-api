// internal/repository/item_repository.go
package repository

import (
	"callable-api/internal/models"
	"callable-api/pkg/errors"
	"sync"
	"fmt"
)

// ItemRepository define a interface para acessar dados de items
type ItemRepository interface {
	// FindAll retorna todos os itens com paginação
	FindAll(page, limit int) ([]models.Item, int, error)
	
	// FindByID retorna um item pelo seu ID
	FindByID(id string) (*models.Item, error)
	
	// Create cria um novo item
	Create(input *models.InputData) (*models.Item, error)
}

// InMemoryItemRepository implementa ItemRepository com armazenamento em memória
// para simplificar demonstrações e testes
type InMemoryItemRepository struct {
	items      map[string]models.Item
	mutex      sync.RWMutex
	nextID     int
}

// NewInMemoryItemRepository cria uma nova instância de InMemoryItemRepository
func NewInMemoryItemRepository() *InMemoryItemRepository {
	repo := &InMemoryItemRepository{
		items:  make(map[string]models.Item),
		nextID: 1,
	}
	
	// Pré-popular com alguns dados de exemplo
	repo.seedData()
	
	return repo
}

// seedData popula o repositório com dados iniciais de exemplo
func (r *InMemoryItemRepository) seedData() {
	// Adicionar alguns itens de exemplo
	for i := 1; i <= 10; i++ {
		id := r.generateID()
		r.items[id] = models.Item{
			ID:          id,
			Name:        "Item " + id,
			Value:       "Value-" + id,
			Description: "Description for item " + id,
			Email:       "user" + id + "@example.com",
			CreatedAt:   "2023-06-01T09:30:00Z",
		}
	}
}

// generateID gera um novo ID único para itens
func (r *InMemoryItemRepository) generateID() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := r.nextID
	r.nextID++
	return fmt.Sprint(id)
}

// FindAll implementa ItemRepository.FindAll
func (r *InMemoryItemRepository) FindAll(page, limit int) ([]models.Item, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	// Calcular o índice inicial e final para paginação
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit
	
	// Coletar todos os itens em um slice
	allItems := make([]models.Item, 0, len(r.items))
	for _, item := range r.items {
		allItems = append(allItems, item)
	}
	
	// Total de itens
	totalItems := len(allItems)
	
	// Verificar se não há itens ou se está além dos limites
	if totalItems == 0 || startIdx >= totalItems {
		return []models.Item{}, totalItems, nil
	}
	
	// Ajustar o índice final se necessário
	if endIdx > totalItems {
		endIdx = totalItems
	}
	
	// Retornar o subconjunto de itens para a página solicitada
	return allItems[startIdx:endIdx], totalItems, nil
}

// FindByID implementa ItemRepository.FindByID
func (r *InMemoryItemRepository) FindByID(id string) (*models.Item, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	item, exists := r.items[id]
	if !exists {
		return nil, errors.NewNotFoundError("Item não encontrado", nil)
	}
	
	return &item, nil
}

// Create implementa ItemRepository.Create
func (r *InMemoryItemRepository) Create(input *models.InputData) (*models.Item, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	id := r.generateID()
	newItem := models.Item{
		ID:          id,
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		Email:       input.Email,
		CreatedAt:   "2023-07-01T10:00:00Z", // Normalmente você usaria time.Now()
	}
	
	r.items[id] = newItem
	
	return &newItem, nil
}