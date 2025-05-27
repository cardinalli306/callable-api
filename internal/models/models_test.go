package models_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

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

func TestResponseWithEmptyData(t *testing.T) {
	// Teste com data vazio
	resp := models.Response{
		Status:  "error",
		Message: "Not found",
		// Data não especificado
	}

	jsonData, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Verificar se o campo 'data' não está presente quando vazio
	assert.NotContains(t, string(jsonData), "data")
	
	// Deserializar de volta
	var decodedResp models.Response
	err = json.Unmarshal(jsonData, &decodedResp)
	assert.NoError(t, err)
	
	assert.Equal(t, "error", decodedResp.Status)
	assert.Equal(t, "Not found", decodedResp.Message)
	assert.Nil(t, decodedResp.Data)
}

func TestInputDataValidations(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		// Este seria um exemplo de dado válido
		input := models.InputData{
			Name:        "Valid Name",
			Value:       "Valid Value",
			Description: "This is a valid description",
			Email:       "valid@example.com",
			CreatedAt:   "2023-05-22T14:56:32Z",
		}

		// Serialize para JSON para verificar se os campos são formatados corretamente
		jsonData, err := json.Marshal(input)
		assert.NoError(t, err)

		// Deserialize para verificar se os campos são preservados
		var decoded models.InputData
		err = json.Unmarshal(jsonData, &decoded)
		assert.NoError(t, err)

		assert.Equal(t, "Valid Name", decoded.Name)
		assert.Equal(t, "Valid Value", decoded.Value)
		assert.Equal(t, "This is a valid description", decoded.Description)
		assert.Equal(t, "valid@example.com", decoded.Email)
		assert.Equal(t, "2023-05-22T14:56:32Z", decoded.CreatedAt)
	})

	// Validamos a interpretação do formato de data manualmente
	t.Run("Date Format Validation", func(t *testing.T) {
		validDate := "2023-05-22T14:56:32Z"
		_, err := time.Parse("2006-01-02T15:04:05Z07:00", validDate)
		assert.NoError(t, err, "Data em formato válido deve ser interpretada corretamente")

		invalidDate := "2023-13-42T99:99:99Z"
		_, err = time.Parse("2006-01-02T15:04:05Z07:00", invalidDate)
		assert.Error(t, err, "Data em formato inválido deve gerar erro")
	})
}

// Testes para os novos métodos adicionados
func TestResponseMethods(t *testing.T) {
	t.Run("IsSuccess", func(t *testing.T) {
		resp := models.Response{Status: "success"}
		assert.True(t, resp.IsSuccess())
		assert.False(t, resp.IsError())
	})

	t.Run("IsError", func(t *testing.T) {
		resp := models.Response{Status: "error"}
		assert.True(t, resp.IsError())
		assert.False(t, resp.IsSuccess())
	})

	t.Run("Other Status", func(t *testing.T) {
		resp := models.Response{Status: "pending"}
		assert.False(t, resp.IsSuccess())
		assert.False(t, resp.IsError())
	})
}

func TestListResponseMethods(t *testing.T) {
	t.Run("GetTotalPages", func(t *testing.T) {
		// Caso com divisão exata
		resp1 := models.ListResponse{TotalRows: 20, PageSize: 10}
		assert.Equal(t, 2, resp1.GetTotalPages())

		// Caso com resto
		resp2 := models.ListResponse{TotalRows: 25, PageSize: 10}
		assert.Equal(t, 3, resp2.GetTotalPages())

		// Caso com pageSize zero ou negativo
		resp3 := models.ListResponse{TotalRows: 25, PageSize: 0}
		assert.Equal(t, 0, resp3.GetTotalPages())

		resp4 := models.ListResponse{TotalRows: 25, PageSize: -5}
		assert.Equal(t, 0, resp4.GetTotalPages())
	})

	t.Run("HasNextPage", func(t *testing.T) {
		// Não tem próxima página (última página)
		resp1 := models.ListResponse{Page: 3, PageSize: 10, TotalRows: 30}
		assert.False(t, resp1.HasNextPage())

		// Tem próxima página
		resp2 := models.ListResponse{Page: 2, PageSize: 10, TotalRows: 30}
		assert.True(t, resp2.HasNextPage())

		// Não tem próxima página (PageSize inválido)
		resp3 := models.ListResponse{Page: 1, PageSize: 0, TotalRows: 30}
		assert.False(t, resp3.HasNextPage())
	})

	t.Run("HasPreviousPage", func(t *testing.T) {
		// Não tem página anterior (primeira página)
		resp1 := models.ListResponse{Page: 1, PageSize: 10, TotalRows: 30}
		assert.False(t, resp1.HasPreviousPage())

		// Tem página anterior
		resp2 := models.ListResponse{Page: 2, PageSize: 10, TotalRows: 30}
		assert.True(t, resp2.HasPreviousPage())
	})
}

func TestItemMethods(t *testing.T) {
	t.Run("HasDescription", func(t *testing.T) {
		// Com descrição
		item1 := models.Item{Description: "Some description"}
		assert.True(t, item1.HasDescription())

		// Sem descrição
		item2 := models.Item{Description: ""}
		assert.False(t, item2.HasDescription())
	})

	t.Run("HasEmail", func(t *testing.T) {
		// Com email
		item1 := models.Item{Email: "test@example.com"}
		assert.True(t, item1.HasEmail())

		// Sem email
		item2 := models.Item{Email: ""}
		assert.False(t, item2.HasEmail())
	})

	t.Run("GetCreatedAtTime", func(t *testing.T) {
		// Data válida
		item1 := models.Item{CreatedAt: "2023-05-22T14:56:32Z"}
		time1, err := item1.GetCreatedAtTime()
		assert.NoError(t, err)
		assert.Equal(t, 2023, time1.Year())
		assert.Equal(t, time.May, time1.Month())
		assert.Equal(t, 22, time1.Day())

		// Data inválida
		item2 := models.Item{CreatedAt: "invalid-date"}
		_, err = item2.GetCreatedAtTime()
		assert.Error(t, err)
	})
}

func TestInputDataValidate(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		input := models.InputData{
			Name:        "Valid Name",
			Value:       "Valid Value",
			Description: "Valid Description",
			Email:       "valid@example.com",
			CreatedAt:   "2023-05-22T14:56:32Z",
		}
		err := input.Validate()
		assert.NoError(t, err)
	})

	t.Run("Name Too Short", func(t *testing.T) {
		input := models.InputData{
			Name:  "AB", // menos de 3 caracteres
			Value: "Valid Value",
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("Name Too Long", func(t *testing.T) {
		longName := strings.Repeat("X", 51) // 51 caracteres
		input := models.InputData{
			Name:  longName,
			Value: "Valid Value",
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("Missing Value", func(t *testing.T) {
		input := models.InputData{
			Name:  "Valid Name",
			Value: "", // valor em branco (obrigatório)
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value")
	})

	t.Run("Description Too Long", func(t *testing.T) {
		longDesc := strings.Repeat("X", 201) // 201 caracteres
		input := models.InputData{
			Name:        "Valid Name",
			Value:       "Valid Value",
			Description: longDesc,
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description")
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		input := models.InputData{
			Name:  "Valid Name",
			Value: "Valid Value",
			Email: "invalid-email", // sem @ e sem ponto
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		input := models.InputData{
			Name:      "Valid Name",
			Value:     "Valid Value",
			CreatedAt: "2023-13-42T99:99:99Z", // data inválida
		}
		err := input.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "date")
	})
}