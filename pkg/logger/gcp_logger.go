package logger

import (
	"context"
	
	"log"
	"os"
	"time"
)

// Níveis de log
const (
	LogLevelDEBUG = "DEBUG"
	LogLevelINFO  = "INFO"
	LogLevelWARN  = "WARN"
	LogLevelERROR = "ERROR"
	LogLevelFATAL = "FATAL"
)

// Logger interface para abstração do logger
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, err error, fields ...map[string]interface{})
	Fatal(msg string, err error, fields ...map[string]interface{})
	Close() error
}

// GCPLogger implementa Logger para Cloud Logging
type GCPLogger struct {
	stdLog *log.Logger
	// Campos simulados - não são usados realmente
	mockProjectID string
	mockLogName   string
}

// NewGCPLogger cria uma nova instância de logger simulado
func NewGCPLogger(ctx context.Context, projectID, logName string, useGCP bool) (Logger, error) {
	stdLog := log.New(os.Stdout, "", log.LstdFlags)

	return &GCPLogger{
		stdLog:        stdLog,
		mockProjectID: projectID,
		mockLogName:   logName,
	}, nil
}

// createEntry cria uma entrada de log com campos adicionais
func (l *GCPLogger) createEntry(msg string, fields ...map[string]interface{}) map[string]interface{} {
	entry := map[string]interface{}{
		"message":   msg,
		"timestamp": time.Now().Format(time.RFC3339),
		"project":   l.mockProjectID, // adicionado para simular integração GCP
		"log_name":  l.mockLogName,   // adicionado para simular integração GCP
	}

	if len(fields) > 0 {
		for k, v := range fields[0] {
			entry[k] = v
		}
	}

	return entry
}

// logToStdout envia log para stdout
func (l *GCPLogger) logToStdout(level, msg string, err error, fields ...map[string]interface{}) {
	entry := l.createEntry(msg, fields...)
	if err != nil {
		entry["error"] = err.Error()
	}

	l.stdLog.Printf("[%s] %+v", level, entry)
}

// Debug registra mensagem de nível debug
func (l *GCPLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelDEBUG, msg, nil, fields...)
}

// Info registra mensagem de nível info
func (l *GCPLogger) Info(msg string, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelINFO, msg, nil, fields...)
}

// Warn registra mensagem de nível warning
func (l *GCPLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelWARN, msg, nil, fields...)
}

// Error registra mensagem de nível erro
func (l *GCPLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelERROR, msg, err, fields...)
}

// Fatal registra mensagem de nível fatal
func (l *GCPLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelFATAL, msg, err, fields...)
	os.Exit(1)
}

// Close simula o fechamento do cliente de logging
func (l *GCPLogger) Close() error {
	l.stdLog.Printf("[INFO] GCP Logger mock fechado com sucesso")
	return nil
}