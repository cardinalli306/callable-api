// errors.go
package models

import "net/http"

// APIError defines a standardized API error
type APIError struct {
    Code    int    `json:"-"`              // HTTP code (not exposed in response)
    Status  string `json:"status"`         // Always "error"
    Message string `json:"message"`        // User-friendly message
    Details string `json:"details,omitempty"` // Technical details (optional)
}

// WithDetails adds details to the error
func (e APIError) WithDetails(details string) APIError {
    e.Details = details
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