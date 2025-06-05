package handlers

import (
	"callable-api/internal/background"
	"callable-api/internal/models"
	"callable-api/pkg/errors"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ItemServiceInterface define os métodos que o handler espera do serviço de itens
type ItemServiceInterface interface {
	GetItems(page, limit int) ([]models.Item, int, error)
	GetItemByID(id string) (*models.Item, error)
	CreateItem(ctx context.Context, input *models.InputData) (*models.Item, error)
}

// ItemHandler gerencia as requisições HTTP relacionadas a itens
type ItemHandler struct {
	itemService    ItemServiceInterface
	jobManager     *background.JobManager
	handlerTimeout time.Duration
}

// NewItemHandler cria uma nova instância de ItemHandler
func NewItemHandler(itemService ItemServiceInterface, jobManager *background.JobManager, handlerTimeout time.Duration) *ItemHandler {
	if handlerTimeout == 0 {
		handlerTimeout = 30 * time.Second // Valor padrão se não for especificado
	}
	return &ItemHandler{
		itemService:    itemService,
		jobManager:     jobManager,
		handlerTimeout: handlerTimeout,
	}
}

// GetData retorna uma lista paginada de itens
// @Summary Listar dados
// @Description Retorna uma lista paginada de itens
// @Tags data
// @Produce json
// @Param page query int false "Número da página (default: 1)"
// @Param limit query int false "Itens por página (default: 10, max: 100)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data [get]
func (h *ItemHandler) GetData(c *gin.Context) {
	// Parse query parameters for pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := h.itemService.GetItems(page, limit)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data: map[string]interface{}{
			"items": items,
			"meta": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

// GetDataById retorna um item específico pelo ID
// @Summary Obter dados por ID
// @Description Retorna um item específico pelo ID
// @Tags data
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data/{id} [get]
func (h *ItemHandler) GetDataById(c *gin.Context) {
	id := c.Param("id")

	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		errors.HandleErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    item,
	})
}

// PostData cria um novo item
// @Summary Criar novo item
// @Description Cria um novo item de dados
// @Tags data
// @Accept json
// @Produce json
// @Security Bearer
// @Param item body models.InputData true "Dados do item"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 408 {object} models.APIError "Request Timeout"
// @Failure 500 {object} models.APIError
// @Router /api/v1/data [post]
func (h *ItemHandler) PostData(c *gin.Context) {
    // Gerar ID de requisição para rastreamento
    reqID := uuid.New().String()
    logger := log.With().Str("request_id", reqID).Str("handler", "PostData").Logger()
    logger.Info().Msg("Iniciando processamento da requisição POST /api/v1/data")

	var input models.InputData

	// Validação de entrada
	logger.Debug().Msg("Iniciando validação de dados de entrada")
	startTime := time.Now()
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error().Err(err).Msg("Erro no bind JSON")
		apiError := errors.NewBadRequestError("Invalid input data", err).ToAPIError()
		c.AbortWithStatusJSON(http.StatusBadRequest, apiError)
		return
	}
	logger.Debug().
		Dur("duration_ms", time.Since(startTime)).
		Interface("input", input).
		Msg("Validação concluída com sucesso")

	// Estratégia de processamento assíncrono imediato com resposta rápida
	// Isso evita timeouts enquanto mantém a semântica de criação síncrona
	jobID := uuid.New().String()
	logger = logger.With().Str("job_id", jobID).Logger()

	// Inicia um job para processamento em background imediatamente
    jobHandle := background.ScheduleJob(h.jobManager, func(ctx context.Context, updateStatus func(progress int, estimatedCompletion *time.Time, result any)) error {
        jobLogger := log.With().
            Str("request_id", reqID).
            Str("job_id", jobID).
            Str("handler", "PostData(async)").
            Logger()

		// Inicialização - 10%
		jobLogger.Info().Msg("Iniciando processamento em background")
		updateStatus(10, nil, nil)

		// Preparação dos dados - 25%
		jobLogger.Debug().Msg("Preparando dados para processamento")
		time.Sleep(200 * time.Millisecond)
		updateStatus(25, nil, nil)

		// Validação - 40%
		jobLogger.Debug().Msg("Validando dados de entrada")
		time.Sleep(200 * time.Millisecond)
		updateStatus(40, nil, nil)

		// Preparação do contexto - 50%
		bgCtx := context.Background()
		jobLogger.Debug().Msg("Preparando contexto para criação do item")
		updateStatus(50, nil, nil)

		// Início da criação do item - 60%
		startTime := time.Now()
		jobLogger.Debug().Msg("Chamando serviço para criar item")
		updateStatus(60, nil, nil)
		
		// Processamento em background - 75%
		time.Sleep(300 * time.Millisecond)
		updateStatus(75, nil, nil)
		
		// Chamada ao serviço de criação
		item, err := h.itemService.CreateItem(bgCtx, &input)
		if err != nil {
			jobLogger.Error().
				Err(err).
				Dur("duration_ms", time.Since(startTime)).
				Msg("Erro ao criar item no serviço")
			return err
		}

		// Finalização - 90%
		jobLogger.Debug().Msg("Finalizando processamento do item")
		updateStatus(90, nil, nil)
		time.Sleep(100 * time.Millisecond)

		jobLogger.Info().
			Dur("duration_ms", time.Since(startTime)).
			Interface("item_id", item.ID).
			Msg("Item criado com sucesso")
		
		// Conclusão - 100%
		updateStatus(100, nil, item)
		return nil
	})

	// Responde imediatamente ao cliente com status 202 (Accepted)
	logger.Info().
        Str("job_id", jobHandle).
        Msg("Solicitação aceita para processamento assíncrono")
	
	 c.JSON(http.StatusAccepted, models.Response{
        Status:  "accepted",
        Message: "Sua solicitação foi aceita e está sendo processada",
        Data: map[string]interface{}{
            "job_id":     jobHandle,
            "status_url": "/api/v1/jobs/" + jobHandle,
        },
    })
}

// PostDataAsync cria um novo item de forma assíncrona
// @Summary Criar novo item de forma assíncrona
// @Description Inicia o processamento assíncrono para criar um novo item
// @Tags data
// @Accept json
// @Produce json
// @Security Bearer
// @Param item body models.InputData true "Dados do item"
// @Success 202 {object} models.Response
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/data/async [post]
func (h *ItemHandler) PostDataAsync(c *gin.Context) {
    reqID := uuid.New().String()
    logger := log.With().Str("request_id", reqID).Str("handler", "PostDataAsync").Logger()
    logger.Info().Msg("Recebendo requisição assíncrona")

    var input models.InputData
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.Error().Err(err).Msg("Erro ao validar dados de entrada")
        apiError := errors.NewBadRequestError("Invalid input data", err).ToAPIError()
        c.AbortWithStatusJSON(http.StatusBadRequest, apiError)
        return
    }

    // Criar o jobID antes de passar para o ScheduleJob
    jobID := uuid.New().String()
    
    // Agende o job diretamente para processamento assíncrono
    background.ScheduleJob(h.jobManager, func(ctx context.Context, updateStatus func(progress int, estimatedCompletion *time.Time, result any)) error {
        jobLogger := log.With().
            Str("request_id", reqID).
            Str("job_id", jobID).
            Str("handler", "PostDataAsync").
            Logger()
			
		// Inicialização - 10%
		jobLogger.Info().Msg("Processando requisição assíncrona")
		updateStatus(10, nil, nil)
		
		// Preparação dos dados - 25%
		jobLogger.Debug().Msg("Preparando dados para processamento")
		time.Sleep(200 * time.Millisecond) 
		updateStatus(25, nil, nil)
		
		// Simulando algum processamento inicial - 40%
		jobLogger.Debug().Msg("Realizando processamento inicial")
		time.Sleep(200 * time.Millisecond)
		updateStatus(40, nil, nil)
		
		// Adicionando request_id ao contexto - 50%
		ctx = context.WithValue(ctx, "request_id", reqID)
		jobLogger.Debug().Msg("Contexto preparado para processamento")
		updateStatus(50, nil, nil)
		
		// Preparação para chamada ao serviço - 60%
		startTime := time.Now()
		jobLogger.Debug().Msg("Preparando chamada ao serviço")
		updateStatus(60, nil, nil)
		
		// Processamento principal - 75%
		jobLogger.Debug().Msg("Chamando serviço para criar item")
		updateStatus(75, nil, nil)
		
		// Chamada ao serviço
		result, err := h.itemService.CreateItem(ctx, &input)
		if err != nil {
			jobLogger.Error().
				Err(err).
				Dur("duration_ms", time.Since(startTime)).
				Msg("Erro no processamento assíncrono")
			return err
		}
		
		// Finalização - 90%
		jobLogger.Debug().Msg("Finalizando processamento do item")
		updateStatus(90, nil, nil)
		time.Sleep(100 * time.Millisecond)
		
		jobLogger.Info().
			Dur("duration_ms", time.Since(startTime)).
			Interface("item_id", result.ID).
			Msg("Processamento assíncrono concluído com sucesso")
			
		// Conclusão - 100%
		updateStatus(100, nil, result)
		return nil
	})

	c.JSON(http.StatusAccepted, models.Response{
		Status:  "success",
		Message: "Request accepted for asynchronous processing",
		Data: map[string]interface{}{
			"job_id":     jobID,
			"status_url": "/api/v1/jobs/" + jobID,
		},
	})
}

// JobStatus retorna o status atual de um job
// @Summary Obter status de um job
// @Description Retorna o status atual de um job em execução
// @Tags jobs
// @Produce json
// @Param id path string true "ID do job"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.APIError
// @Router /api/v1/jobs/{id} [get]
func (h *ItemHandler) JobStatus(c *gin.Context) {
	jobID := c.Param("id")
	logger := log.With().Str("job_id", jobID).Str("handler", "JobStatus").Logger()
	logger.Debug().Msg("Verificando status de job")
	
	status := h.jobManager.GetJobStatus(jobID)
	if status == nil {
		logger.Warn().Msg("Job não encontrado")
		c.JSON(http.StatusNotFound, models.APIError{
   		Status:  "error",
    	Message: "Job not found",
    	Code:    http.StatusNotFound,
		})
		return
	}
	
	logger.Debug().Interface("status", status).Msg("Status de job recuperado")
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Job status retrieved successfully",
		Data:    status,
	})
}

// HealthCheck responde com informações de status da API
// @Summary Verificar status da API
// @Description Retorna o status atual da API
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"message": "Callable API is up and running",
	})
}