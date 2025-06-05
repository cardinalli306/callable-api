package background

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// JobStatus representa o status atual de um job em background
type JobStatus struct {
	ID                  string    `json:"id"`
	State               string    `json:"state"` // "pending", "processing", "completed", "failed"
	Progress            int       `json:"progress"`
	StartTime           time.Time `json:"start_time"`
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
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				m.jobsLock.Lock()
				status.State = "failed"
				status.Error = "Panic in background job: " + stringify(r)
				m.jobsLock.Unlock()
			}
			cancel()
		}()
		
		status.State = "processing"
		
		updateStatus := func(progress int, estimatedCompletion *time.Time, result any) {
			m.jobsLock.Lock()
			status.Progress = progress
			if estimatedCompletion != nil {
				status.EstimatedCompletion = *estimatedCompletion
			}
			if result != nil {
				status.Result = result
			}
			m.jobsLock.Unlock()
		}
		
		err := job(ctx, updateStatus)
		
		m.jobsLock.Lock()
		if err != nil {
			status.State = "failed"
			status.Error = err.Error()
		} else {
			status.State = "completed"
			status.Progress = 100
		}
		m.jobsLock.Unlock()
	}()
}

// GetJobStatus retorna o status atual de um job
func (m *JobManager) GetJobStatus(jobID string) *JobStatus {
	m.jobsLock.RLock()
	defer m.jobsLock.RUnlock()
	
	if status, ok := m.jobs[jobID]; ok {
		// Retorna uma cópia para evitar condições de corrida
		statusCopy := *status
		return &statusCopy
	}
	return nil
}

// CleanupCompletedJobs remove jobs concluídos há mais de certo tempo
func (m *JobManager) CleanupCompletedJobs(olderThan time.Duration) {
	threshold := time.Now().Add(-olderThan)
	
	m.jobsLock.Lock()
	defer m.jobsLock.Unlock()
	
	for id, status := range m.jobs {
		if (status.State == "completed" || status.State == "failed") && 
		   status.StartTime.Before(threshold) {
			delete(m.jobs, id)
		}
	}
}

// stringify converte qualquer valor para string
func stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v)
}