// internal/handlers/handlers.go
package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"callable-api/internal/models"
	"callable-api/internal/service"
	"callable-api/pkg/errors"
	)

// ItemHandler gerencia as requisições HTTP relacionadas a itens
type ItemHandler struct {
	itemService *service.ItemService
}

// NewItemHandler cria uma nova instância de ItemHandler
func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// GetData retorna uma lista paginada de itens
// (Mantendo a assinatura original para compatibilidade com swagger)
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
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"message": "Callable API is up and running",
	})
}