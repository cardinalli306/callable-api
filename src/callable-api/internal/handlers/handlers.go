package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"callable-api/internal/models"
	"callable-api/pkg/logger"
)

// HealthCheck godoc
// @Summary Check API status
// @Description Returns a 200 status if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} models.Response "API is running"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "API is running",
	})
}

// GetData godoc
// @Summary Get data list
// @Description Returns a paginated list of available items
// @Tags items
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) maximum(100)
// @Success 200 {object} models.ListResponse{data=[]models.Item} "Data retrieved successfully"
// @Failure 400 {object} models.Response "Invalid request"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/v1/data [get]
func GetData(c *gin.Context) {
	// Extract pagination parameters
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

	// Simulating response data with pagination
	data := []models.Item{
		{
			ID: "1", 
			Name: "Item 1", 
			Value: "ABC123", 
			Description: "Description for Item 1",
			Email: "user1@example.com",
			CreatedAt: "2023-05-22T14:56:32Z",
		},
		{
			ID: "2", 
			Name: "Item 2", 
			Value: "XYZ456", 
			Description: "Description for Item 2",
			Email: "user2@example.com",
			CreatedAt: "2023-05-23T10:15:45Z",
		},
	}
	
	totalItems := 42 // Simulating total count from database

	logger.Info("Data retrieved successfully", map[string]interface{}{
		"page": page,
		"limit": limit,
		"total": totalItems,
	})

	c.JSON(http.StatusOK, models.ListResponse{
		Status:    "success",
		Message:   "Data retrieved successfully",
		Data:      data,
		Page:      page,
		PageSize:  limit,
		TotalRows: totalItems,
	})
}

// PostData godoc
// @Summary Create new item
// @Description Add a new item based on provided data
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body models.InputData true "Item data"
// @Success 201 {object} models.Response{data=models.Item} "Item created"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/v1/data [post]
func PostData(c *gin.Context) {
	var input models.InputData
	
	// Request body validation
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Warn("Invalid input data", map[string]interface{}{
			"error": err.Error(),
		})
		
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid input: " + err.Error(),
		})
		return
	}
	
	// Create a simulated response item
	newItem := models.Item{
		ID:          "new-generated-id",
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		Email:       input.Email,
		CreatedAt:   input.CreatedAt,
	}
	
	logger.Info("New item created", map[string]interface{}{
		"name": input.Name,
		"id":   newItem.ID,
	})
	
	// Process received data and return the new item
	c.JSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Data saved successfully",
		Data:    newItem,
	})
}

// GetDataById godoc
// @Summary Get item by ID
// @Description Returns a specific item based on provided ID
// @Tags items
// @Produce json
// @Security Bearer
// @Param id path string true "Item ID" format(uuid)
// @Success 200 {object} models.Response{data=models.Item} "Item found"
// @Failure 400 {object} models.Response "Invalid ID format"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 404 {object} models.Response "Item not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/v1/data/{id} [get]
func GetDataById(c *gin.Context) {
	// Get ID from URL
	id := c.Param("id")
	
	logger.Info("Fetching item by ID", map[string]interface{}{
		"id": id,
	})
	
	// Simulate item found
	item := models.Item{
		ID:          id,
		Name:        "Item " + id,
		Value:       "Value-" + id,
		Description: "Description for item " + id,
		Email:       "user" + id + "@example.com",
		CreatedAt:   "2023-06-01T09:30:00Z",
	}
	
	// Here you would normally fetch data from a database
	// In a real implementation, you would check if the item exists
	// and return a 404 if not found
	
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    item,
	})
}