package main

// Response representa o formato padrão de resposta da API
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InputData representa os dados de entrada da API com validação aprimorada
type InputData struct {
	Name        string `json:"name" binding:"required,min=3,max=50" example:"Nome do Item"`
	Value       string `json:"value" binding:"required,min=1" example:"123ABC"`
	Description string `json:"description" binding:"omitempty,max=200" example:"Descrição detalhada do item"`
	Email       string `json:"email" binding:"omitempty,email" example:"usuario@exemplo.com"`
	CreatedAt   string `json:"created_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00" example:"2023-05-22T14:56:32Z"`
}