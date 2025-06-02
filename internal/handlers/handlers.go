package handlers

import (
	"callable-api/internal/models"
	"callable-api/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemServiceInterface define os métodos que o handler espera do serviço de itens
type ItemServiceInterface interface {
	GetItems(page, limit int) ([]models.Item, int, error)
	GetItemByID(id string) (*models.Item, error)
	CreateItem(input *models.InputData) (*models.Item, error)
}

// ItemHandler gerencia as requisições HTTP relacionadas a itens
type ItemHandler struct {
	itemService ItemServiceInterface
}

// NewItemHandler cria uma nova instância de ItemHandler
func NewItemHandler(itemService ItemServiceInterface) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// GetData retorna uma lista paginada de itens
// @Summary Listar dados
// @Description Retorna uma lista paginada de itens
// @Tags data
// @Produce json
// @Param page query int false "Número da página (default: 1)"
// @Param limit query int false "Itens por página (default: 10, max: 100)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data [get]
func (h *ItemHandler) GetData(c *gin.Context) {
	// Parse query parameters for pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := h.itemService.GetItems(page, limit)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data: map[string]interface{}{
			"items": items,
			"meta": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

// GetDataById retorna um item específico pelo ID
// @Summary Obter dados por ID
// @Description Retorna um item específico pelo ID
// @Tags data
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data/{id} [get]
func (h *ItemHandler) GetDataById(c *gin.Context) {
	id := c.Param("id")

	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    item,
	})
}

// PostData cria um novo item
// @Summary Criar novo item
// @Description Cria um novo item de dados
// @Tags data
// @Accept json
// @Produce json
// @Security Bearer
// @Param item body models.InputData true "Dados do item"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data [post]
func (h *ItemHandler) PostData(c *gin.Context) {
	var input models.InputData

	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleErrors(c, errors.NewBadRequestError("Invalid input data", err))
		return
	}

	item, err := h.itemService.CreateItem(&input)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Data created successfully",
		Data:    item,
	})
}

// HealthCheck responde com informações de status da API
// @Summary Verificar status da API
// @Description Retorna o status atual da API
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"message": "Callable API is up and running",
	})
}
