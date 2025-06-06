package background

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// JobStatus representa o status atual de um job em background
type JobStatus struct {
	ID                  string    `json:"id"`
	State               string    `json:"state"` // "pending", "processing", "completed", "failed"
	Progress            int       `json:"progress"`
	StartTime           time.Time `json:"start_time"`
	CompletionTime      time.Time `json:"completion_time,omitempty"`
	EstimatedCompletion time.Time `json:"estimated_completion,omitempty"`
	Error               string    `json:"error,omitempty"`
	Result              any       `json:"result,omitempty"`
}

// JobManager gerencia todas as goroutines em background
type JobManager struct {
	jobs     map[string]*JobStatus
	jobsLock sync.RWMutex
}

// NewJobManager cria uma nova instância do gerenciador de jobs
func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*JobStatus),
	}
}

// StartJob inicia uma nova tarefa em background
func (m *JobManager) StartJob(jobID string, maxDuration time.Duration, job func(ctx context.Context, updateStatus func(progress int, estimatedCompletion *time.Time, result any)) error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	
	// Inicializa o status do job
	status := &JobStatus{
		ID:        jobID,
		State:     "pending",
		Progress:  0,
		StartTime: time.Now(),
	}
	
	m.jobsLock.Lock()
	m.jobs[jobID] = status
	m.jobsLock.Unlock()
	
	// Registra início do job
	log.Info().
		Str("job_id", jobID).
		Msg("Job iniciado")
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				m.jobsLock.Lock()
				status.State = "failed"
				status.Error = "Panic in background job: " + stringify(r)
				status.CompletionTime = time.Now()
				m.jobsLock.Unlock()
				
				log.Error().
					Str("job_id", jobID).
					Interface("panic", r).
					Msg("Job falhou com panic")
			}
			cancel()
		}()
		
		m.jobsLock.Lock()
		status.State = "processing"
		m.jobsLock.Unlock()
		
		log.Debug().
			Str("job_id", jobID).
			Msg("Job em processamento")
		
		// Implementação melhorada da função updateStatus
		updateStatus := func(progress int, estimatedCompletion *time.Time, result any) {
			m.jobsLock.Lock()
			defer m.jobsLock.Unlock()
			
			prevProgress := status.Progress
			status.Progress = progress
			
			if estimatedCompletion != nil {
				status.EstimatedCompletion = *estimatedCompletion
			}
			
			if result != nil {
				status.Result = result
				// Quando recebemos um resultado, consideramos o job completado automaticamente
				if status.State == "processing" && progress >= 100 {
					status.State = "completed"
					status.CompletionTime = time.Now()
				}
			}
			
			// Registrar mudanças significativas no progresso
			if progress != prevProgress {
				log.Debug().
					Str("job_id", jobID).
					Int("progress", progress).
					Msg("Progresso atualizado")
			}
		}
		
		// Monitoramento de timeout separado
		done := make(chan struct{})
		
		go func() {
			err := job(ctx, updateStatus)
			
			m.jobsLock.Lock()
			if err != nil {
				status.State = "failed"
				status.Error = err.Error()
				status.CompletionTime = time.Now()
				
				log.Error().
					Str("job_id", jobID).
					Err(err).
					Msg("Job falhou com erro")
			} else if status.State != "completed" {
				// Certifique-se de que seja marcado como concluído mesmo se
				// a função updateStatus não foi chamada com progress=100
				status.State = "completed"
				status.Progress = 100
				status.CompletionTime = time.Now()
				
				log.Info().
					Str("job_id", jobID).
					Msg("Job concluído com sucesso")
			}
			m.jobsLock.Unlock()
			
			close(done)
		}()
		
		// Aguardar conclusão ou timeout
		select {
		case <-done:
			// Job concluído normalmente
		case <-ctx.Done():
			// Timeout ocorreu
			m.jobsLock.Lock()
			if status.State == "processing" {
				status.State = "failed"
				status.Error = "Job timeout: excedeu o tempo máximo permitido"
				status.CompletionTime = time.Now()
				
				log.Error().
					Str("job_id", jobID).
					Dur("max_duration", maxDuration).
					Msg("Job cancelado por timeout")
			}
			m.jobsLock.Unlock()
		}
	}()
}

// GetJobStatus retorna o status atual de um job
func (m *JobManager) GetJobStatus(jobID string) (*JobStatus, error) {
	m.jobsLock.RLock()
	defer m.jobsLock.RUnlock()
	
	if status, ok := m.jobs[jobID]; ok {
		// Registrar acesso ao status
		log.Debug().
			Str("job_id", jobID).
			Str("state", status.State).
			Int("progress", status.Progress).
			Msg("Obtendo status do job")
			
		// Retorna uma cópia para evitar condições de corrida
		statusCopy := *status
		return &statusCopy, nil
	}
	
	log.Debug().
		Str("job_id", jobID).
		Msg("Job não encontrado ao consultar status")
	
	return nil, fmt.Errorf("job não encontrado: %s", jobID)
}

// CleanupCompletedJobs remove jobs concluídos há mais de certo tempo
func (m *JobManager) CleanupCompletedJobs(olderThan time.Duration) {
	threshold := time.Now().Add(-olderThan)
	count := 0
	
	m.jobsLock.Lock()
	defer m.jobsLock.Unlock()
	
	for id, status := range m.jobs {
		if (status.State == "completed" || status.State == "failed") && 
		   status.StartTime.Before(threshold) {
			delete(m.jobs, id)
			count++
		}
	}
	
	if count > 0 {
		log.Info().
			Int("count", count).
			Dur("older_than", olderThan).
			Msg("Jobs antigos removidos")
	}
}

// stringify converte qualquer valor para string
func stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v)
}