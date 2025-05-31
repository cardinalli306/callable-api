package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"callable-api/internal/models"
)

// AppError é uma estrutura que representa um erro da aplicação
type AppError struct {
	StatusCode int
	Type       string
	Message    string
	Details    string
	Err        error
	Stack      string
}

// Error implementa a interface error
func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ToAPIError converte um AppError para models.APIError para resposta HTTP
func (e AppError) ToAPIError() models.APIError {
	return models.APIError{
		Code:    e.StatusCode,
		Status:  "error",
		Message: e.Message,
		Details: e.Details,
	}
}

// ValidationFieldError representa um erro de validação para um campo específico
type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError representa um erro com múltiplos campos de validação
type ValidationError struct {
	AppError
	FieldErrors []ValidationFieldError `json:"field_errors,omitempty"`
}

// NewValidationError cria um novo erro de validação
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			StatusCode: http.StatusBadRequest,
			Type:       "VALIDATION_ERROR",
			Message:    message,
			Stack:      captureStack(),
		},
		FieldErrors: make([]ValidationFieldError, 0),
	}
}

// AddFieldError adiciona um erro de campo à lista de erros de validação
func (e *ValidationError) AddFieldError(field, message string) *ValidationError {
	e.FieldErrors = append(e.FieldErrors, ValidationFieldError{
		Field:   field,
		Message: message,
	})
	return e
}

// ToAPIError sobrescreve o método para incluir erros de campo
func (e ValidationError) ToAPIError() models.APIError {
	apiErr := e.AppError.ToAPIError()
	
	// Converter erros de campo para o formato esperado pela API
	fieldErrors := make(map[string]string, len(e.FieldErrors))
	for _, fieldErr := range e.FieldErrors {
		fieldErrors[fieldErr.Field] = fieldErr.Message
	}
	
	return apiErr.WithFieldErrors(fieldErrors)
}

// captureStack captura a pilha de chamadas para ajudar na depuração
func captureStack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "callable-api") {
			if more {
				continue
			}
			break
		}
		fmt.Fprintf(&builder, "%s:%d\n", frame.Function, frame.Line)
		if !more {
			break
		}
	}
	return builder.String()
}

// New cria um novo erro da aplicação
func New(statusCode int, errType string, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Type:       errType,
		Message:    message,
		Err:        err,
		Stack:      captureStack(),
	}
}

// WithDetails adiciona detalhes ao erro
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Funções helpers para criar erros específicos
func NewBadRequestError(message string, err error) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}

func NewUnauthorizedError(message string, err error) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message, err)
}

func NewForbiddenError(message string, err error) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, err)
}

func NewConflictError(message string, err error) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, err)
}

func NewTooManyRequestsError(message string, err error) *AppError {
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message, err)
}

func NewInternalServerError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER", message, err)
}

func NewServiceUnavailableError(message string, err error) *AppError {
	return New(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, err)
}

func NewPaymentRequiredError(message string, err error) *AppError {
	return New(http.StatusPaymentRequired, "PAYMENT_REQUIRED", message, err)
}

func NewMethodNotAllowedError(message string, err error) *AppError {
	return New(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", message, err)
}