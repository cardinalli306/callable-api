package background

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// JobFunc é o tipo de função que será executada em background
type JobFunc func(ctx context.Context, updateStatus func(progress int, estimatedCompletion *time.Time, result any)) error

// ScheduleJob agenda uma nova tarefa para execução em background
func ScheduleJob(manager *JobManager, job JobFunc) string {
	jobID := uuid.New().String()
	manager.StartJob(jobID, 30*time.Minute, job) // 30 minutos como timeout padrão
	return jobID
}