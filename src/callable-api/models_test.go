package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseStructure(t *testing.T) {
	// Cria uma resposta de exemplo
	resp := Response{
		Status:  "success",
		Message: "Test message",
		Data:    map[string]interface{}{"key": "value"},
	}

	// Serializa para JSON
	jsonData, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Deserializa de volta para validar estrutura
	var decodedResp Response
	err = json.Unmarshal(jsonData, &decodedResp)
	assert.NoError(t, err)

	// Verifica se os campos foram preservados
	assert.Equal(t, "success", decodedResp.Status)
	assert.Equal(t, "Test message", decodedResp.Message)

	// Verifica se os dados foram preservados
	data, ok := decodedResp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestInputDataValidation(t *testing.T) {
	// Teste 1: JSON válido
	validJSON := `{
		"name": "Valid Name",
		"value": "Valid Value",
		"description": "Valid Description",
		"email": "valid@example.com",
		"created_at": "2023-05-22T14:56:32Z"
	}`

	var validInput InputData
	err := json.Unmarshal([]byte(validJSON), &validInput)
	assert.NoError(t, err)
	assert.Equal(t, "Valid Name", validInput.Name)
	assert.Equal(t, "Valid Value", validInput.Value)

	// Teste 2: JSON com email inválido
	invalidJSON := `{
		"name": "Valid Name",
		"value": "Valid Value",
		"email": "invalid-email"
	}`

	var invalidInput InputData
	err = json.Unmarshal([]byte(invalidJSON), &invalidInput)
	assert.NoError(t, err) // O unmarshalling não falha com email inválido
	// A validação aconteceria pelo binding do Gin
}
