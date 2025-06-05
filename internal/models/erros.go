package models

import "net/http"

// APIError defines a standardized API error
// APIError representa um erro da API que será retornado ao cliente
type APIError struct {
    Status      string                 `json:"status"`
    Message     string                 `json:"message"`
    Details     string                 `json:"details,omitempty"`
    Code        int                    `json:"code"`                 // Manter como int para compatibilidade
    CodeString  string                 `json:"code_string,omitempty"` // Adicionar campo string opcional
    FieldErrors map[string]string      `json:"field_errors,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Adicionar campo para armazenar erros de validação
func (e *APIError) WithFieldErrors(fieldErrors map[string]string) *APIError {
    e.FieldErrors = fieldErrors
    return e
}

// Adicionar método para configurar metadados
func (e *APIError) WithMetadata(metadata map[string]interface{}) *APIError {
    e.Metadata = metadata
    return e
}

// Adicionar método para configurar código como string
func (e *APIError) WithCodeString(codeString string) *APIError {
    e.CodeString = codeString
    return e
}

// Common predefined errors
var (
	ErrInvalidInput = APIError{
		Code:    http.StatusBadRequest,
		Status:  "error",
		Message: "Invalid input data",
	}

	ErrResourceNotFound = APIError{ 
		Code:    http.StatusNotFound,
		Status:  "error",
		Message: "Resource not found",
	}

	ErrUnauthorized = APIError{
		Code:    http.StatusUnauthorized,
		Status:  "error",
		Message: "Authentication required",
	}

	ErrInternalServer = APIError{
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Message: "Internal server error",
	}
)