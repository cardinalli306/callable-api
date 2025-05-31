package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/logging"
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
	client *logging.Client
	logger *logging.Logger
	useGCP bool
	stdLog *log.Logger
}

// NewGCPLogger cria uma nova instância de logger
func NewGCPLogger(ctx context.Context, projectID, logName string, useGCP bool) (Logger, error) {
	stdLog := log.New(os.Stdout, "", log.LstdFlags)

	if !useGCP || projectID == "" {
		return &GCPLogger{
			useGCP: false,
			stdLog: stdLog,
		}, nil
	}

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente de logging: %v", err)
	}

	return &GCPLogger{
		client: client,
		logger: client.Logger(logName),
		useGCP: true,
		stdLog: stdLog,
	}, nil
}

// createEntry cria uma entrada de log com campos adicionais
func (l *GCPLogger) createEntry(msg string, fields ...map[string]interface{}) map[string]interface{} {
	entry := map[string]interface{}{
		"message":   msg,
		"timestamp": time.Now().Format(time.RFC3339),
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

	if l.useGCP {
		l.logger.Log(logging.Entry{
			Severity: logging.Debug,
			Payload:  l.createEntry(msg, fields...),
		})
	}
}

// Info registra mensagem de nível info
func (l *GCPLogger) Info(msg string, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelINFO, msg, nil, fields...)

	if l.useGCP {
		l.logger.Log(logging.Entry{
			Severity: logging.Info,
			Payload:  l.createEntry(msg, fields...),
		})
	}
}

// Warn registra mensagem de nível warning
func (l *GCPLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelWARN, msg, nil, fields...)

	if l.useGCP {
		l.logger.Log(logging.Entry{
			Severity: logging.Warning,
			Payload:  l.createEntry(msg, fields...),
		})
	}
}

// Error registra mensagem de nível erro
func (l *GCPLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelERROR, msg, err, fields...)

	if l.useGCP && err != nil {
		payload := l.createEntry(msg, fields...)
		payload["error"] = err.Error()

		l.logger.Log(logging.Entry{
			Severity: logging.Error,
			Payload:  payload,
		})
	}
}

// Fatal registra mensagem de nível fatal
func (l *GCPLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.logToStdout(LogLevelFATAL, msg, err, fields...)

	if l.useGCP && err != nil {
		payload := l.createEntry(msg, fields...)
		payload["error"] = err.Error()

		l.logger.Log(logging.Entry{
			Severity: logging.Critical,
			Payload:  payload,
		})
	}

	os.Exit(1)
}

// Close fecha o cliente de logging
func (l *GCPLogger) Close() error {
	if l.useGCP {
		return l.client.Close()
	}
	return nil
}
