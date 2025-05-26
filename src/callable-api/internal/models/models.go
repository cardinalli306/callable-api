package models

// Response represents the standard API response format
type Response struct {
    Status  string      `json:"status" example:"success"`
    Message string      `json:"message" example:"Operation completed successfully"`
    Data    interface{} `json:"data,omitempty"`
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

// Item represents a complete data item returned by the API
type Item struct {
    ID          string `json:"id" example:"5f8d0e6e-6c0a-4f0a-8e0a-6c0a4f0a8e0a"`
    Name        string `json:"name" example:"Item Name"`
    Value       string `json:"value" example:"ABC123"`
    Description string `json:"description,omitempty" example:"Detailed item description"`
    Email       string `json:"email,omitempty" example:"user@example.com"`
    CreatedAt   string `json:"created_at" example:"2023-05-22T14:56:32Z"`
}

// InputData represents API input data with enhanced validation
type InputData struct {
    Name        string `json:"name" binding:"required,min=3,max=50" example:"Item Name"`
    Value       string `json:"value" binding:"required,min=1" example:"123ABC"`
    Description string `json:"description" binding:"omitempty,max=200" example:"Detailed item description"`
    Email       string `json:"email" binding:"omitempty,email" example:"user@example.com"`
    CreatedAt   string `json:"created_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00" example:"2023-05-22T14:56:32Z"`
}