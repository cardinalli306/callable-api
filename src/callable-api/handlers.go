package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// @Summary Verifica o status da API
// @Description Retorna um status 200 se a API estiver rodando
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "API is running",
	})
}

// @Summary Obtém lista de dados
// @Description Retorna uma lista de itens disponíveis
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/data [get]
func getData(c *gin.Context) {
	// Simulando dados de resposta
	data := []map[string]interface{}{
		{"id": 1, "name": "Item 1", "description": "Description for Item 1"},
		{"id": 2, "name": "Item 2", "description": "Description for Item 2"},
	}

	log.Printf("Dados recuperados com sucesso: %d itens", len(data))

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    data,
	})
}

// @Summary Cria um novo item
// @Description Adiciona um novo item com base nos dados fornecidos
// @Accept json
// @Produce json
// @Param data body InputData true "Dados do item"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/data [post]
func postData(c *gin.Context) {
	var input InputData
	
	// Validação do corpo da requisição
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid input: " + err.Error(),
		})
		return
	}
	
	log.Printf("Novo item criado: %s", input.Name)
	
	// Processa os dados recebidos (aqui apenas retornamos como exemplo)
	c.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "Data saved successfully",
		Data:    input,
	})
}

// @Summary Obtém item por ID
// @Description Retorna um item específico com base no ID fornecido
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} Response
// @Router /api/v1/data/{id} [get]
func getDataById(c *gin.Context) {
	// Obtém o ID da URL
	id := c.Param("id")
	
	log.Printf("Buscando item por ID: %s", id)
	
	// Aqui você normalmente buscaria dados em um banco de dados
	// Neste exemplo, apenas retornamos o ID recebido
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    map[string]interface{}{"id": id},
	})
}