package models_test

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"callable-api/internal/models"
)

func TestResponseStructure(t *testing.T) {
	// Create an example response
	resp := models.Response{
		Status:  "success",
		Message: "Test message",
		Data:    map[string]interface{}{"key": "value"},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Deserialize back to validate structure
	var decodedResp models.Response
	err = json.Unmarshal(jsonData, &decodedResp)
	assert.NoError(t, err)

	// Check if fields were preserved
	assert.Equal(t, "success", decodedResp.Status)
	assert.Equal(t, "Test message", decodedResp.Message)

	// Check if data was preserved
	data, ok := decodedResp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestListResponseStructure(t *testing.T) {
	// Create an example list response
	resp := models.ListResponse{
		Status:    "success",
		Message:   "Data retrieved successfully",
		Data:      []map[string]interface{}{{"id": "1", "name": "Test Item"}},
		Page:      2,
		PageSize:  10,
		TotalRows: 42,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Deserialize back to validate structure
	var decodedResp models.ListResponse
	err = json.Unmarshal(jsonData, &decodedResp)
	assert.NoError(t, err)

	// Check if fields were preserved
	assert.Equal(t, "success", decodedResp.Status)
	assert.Equal(t, "Data retrieved successfully", decodedResp.Message)
	assert.Equal(t, 2, decodedResp.Page)
	assert.Equal(t, 10, decodedResp.PageSize)
	assert.Equal(t, 42, decodedResp.TotalRows)
}

func TestItemStructure(t *testing.T) {
	// Create an example item
	item := models.Item{
		ID:          "5f8d0e6e-6c0a-4f0a-8e0a-6c0a4f0a8e0a",
		Name:        "Test Item",
		Value:       "ABC123",
		Description: "Test Description",
		Email:       "user@example.com",
		CreatedAt:   "2023-05-22T14:56:32Z",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(item)
	assert.NoError(t, err)

	// Deserialize back to validate structure
	var decodedItem models.Item
	err = json.Unmarshal(jsonData, &decodedItem)
	assert.NoError(t, err)

	// Check if fields were preserved
	assert.Equal(t, "5f8d0e6e-6c0a-4f0a-8e0a-6c0a4f0a8e0a", decodedItem.ID)
	assert.Equal(t, "Test Item", decodedItem.Name)
	assert.Equal(t, "ABC123", decodedItem.Value)
	assert.Equal(t, "Test Description", decodedItem.Description)
	assert.Equal(t, "user@example.com", decodedItem.Email)
	assert.Equal(t, "2023-05-22T14:56:32Z", decodedItem.CreatedAt)
}

func TestInputDataValidation(t *testing.T) {
	// Set up validator
	validate := validator.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
	}

	// Test 1: Valid JSON
	validInput := models.InputData{
		Name:        "Valid Name",
		Value:       "Valid Value",
		Description: "Valid Description",
		Email:       "valid@example.com",
		CreatedAt:   "2023-05-22T14:56:32Z",
	}

	// Validate valid input
	err := validate.Struct(validInput)
	assert.NoError(t, err)

	// Test 2: Invalid name (too short)
	invalidInput1 := models.InputData{
		Name:  "ab", // less than required 3 chars
		Value: "Valid Value",
	}

	// Validate invalid name
	err = validate.Struct(invalidInput1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name")

	// Test 3: Invalid email format
	invalidInput2 := models.InputData{
		Name:  "Valid Name",
		Value: "Valid Value",
		Email: "invalid-email",
	}

	// Validate invalid email
	err = validate.Struct(invalidInput2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email")

	// Test 4: Missing required field
	invalidInput3 := models.InputData{
		Name: "Valid Name",
		// Value is missing but required
	}

	// Validate missing required field
	err = validate.Struct(invalidInput3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Value")
}