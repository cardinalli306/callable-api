package models

import (
	"fmt"
	"strings"
	"time"
)

// Response represents the standard API response format
type Response struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// IsSuccess returns true if the response status is "success"
func (r *Response) IsSuccess() bool {
	return r.Status == "success"
}

// IsError returns true if the response status is "error"
func (r *Response) IsError() bool {
	return r.Status == "error"
}

// ListResponse is the model for paginated list responses
type ListResponse struct {
	Status    string      `json:"status" example:"success"`
	Message   string      `json:"message" example:"Data retrieved successfully"`
	Data      interface{} `json:"data"`
	Page      int         `json:"page" example:"1"`
	PageSize  int         `json:"page_size" example:"10"`
	TotalRows int         `json:"total_rows" example:"42"`
}

// GetTotalPages calculates the total number of pages based on total rows and page size
func (lr *ListResponse) GetTotalPages() int {
	if lr.PageSize <= 0 {
		return 0
	}
	totalPages := lr.TotalRows / lr.PageSize
	if lr.TotalRows%lr.PageSize > 0 {
		totalPages++
	}
	return totalPages
}

// HasNextPage returns true if there are more pages after the current one
func (lr *ListResponse) HasNextPage() bool {
	return lr.Page < lr.GetTotalPages()
}

// HasPreviousPage returns true if there are pages before the current one
func (lr *ListResponse) HasPreviousPage() bool {
	return lr.Page > 1
}

// Item represents a complete data item returned by the API
type Item struct {
	ID          string `json:"id" example:"5f8d0e6e-6c0a-4f0a-8e0a-6c0a4f0a8e0a"`
	Name        string `json:"name" example:"Item Name"`
	Value       string `json:"value" example:"ABC123"`
	Description string `json:"description,omitempty" example:"Detailed item description"`
	Email       string `json:"email,omitempty" example:"user@example.com"`
	CreatedAt   string `json:"created_at" example:"2023-05-22T14:56:32Z"`
}

// HasDescription returns true if the item has a non-empty description
func (i *Item) HasDescription() bool {
	return i.Description != ""
}

// HasEmail returns true if the item has a non-empty email
func (i *Item) HasEmail() bool {
	return i.Email != ""
}

// GetCreatedAtTime attempts to parse the CreatedAt field as a time.Time
func (i *Item) GetCreatedAtTime() (time.Time, error) {
	return time.Parse(time.RFC3339, i.CreatedAt)
}

// InputData represents API input data with enhanced validation
type InputData struct {
	Name        string `json:"name" binding:"required,min=3,max=50" example:"Item Name"`
	Value       string `json:"value" binding:"required,min=1" example:"123ABC"`
	Description string `json:"description" binding:"omitempty,max=200" example:"Detailed item description"`
	Email       string `json:"email" binding:"omitempty,email" example:"user@example.com"`
	CreatedAt   string `json:"created_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00" example:"2023-05-22T14:56:32Z"`
}

// Validate performs basic validation on the input data
func (i *InputData) Validate() error {
	if len(i.Name) < 3 || len(i.Name) > 50 {
		return fmt.Errorf("name must be between 3 and 50 characters")
	}
	if i.Value == "" {
		return fmt.Errorf("value is required")
	}
	if len(i.Description) > 200 {
		return fmt.Errorf("description must not exceed 200 characters")
	}
	if i.Email != "" {
		// Basic email validation
		if !strings.Contains(i.Email, "@") || !strings.Contains(i.Email, ".") {
			return fmt.Errorf("invalid email format")
		}
	}
	if i.CreatedAt != "" {
		_, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			return fmt.Errorf("invalid date format (should be RFC3339): %v", err)
		}
	}
	return nil
}