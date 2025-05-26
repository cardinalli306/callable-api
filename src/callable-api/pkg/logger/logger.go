// pkg/logger/logger.go
package logger

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

// Level define os níveis de log
type Level int

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
)

// String converte nível para texto
func (l Level) String() string {
    return [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[l]
}

// Configuração global
var currentLevel = INFO

// SetLevel define o nível mínimo de log
func SetLevel(level string) {
    switch level {
    case "debug":
        currentLevel = DEBUG
    case "info":
        currentLevel = INFO
    case "warn":
        currentLevel = WARN
    case "error":
        currentLevel = ERROR
    default:
        currentLevel = INFO
    }
}

// Estrutura do log
type logEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Log registra uma mensagem no nível especificado
func Log(level Level, message string, fields map[string]interface{}) {
    // Ignora logs abaixo do nível configurado
    if level < currentLevel {
        return
    }

    entry := logEntry{
        Timestamp: time.Now().Format(time.RFC3339),
        Level:     level.String(),
        Message:   message,
        Fields:    fields,
    }

    // Converte para JSON
    jsonData, err := json.Marshal(entry)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error logging: %v\n", err)
        return
    }

    // Logs de erro vão para stderr, outros para stdout
    if level == ERROR {
        fmt.Fprintln(os.Stderr, string(jsonData))
    } else {
        fmt.Println(string(jsonData))
    }
}

// Funções de conveniência
func Debug(message string, fields map[string]interface{}) {
    Log(DEBUG, message, fields)
}

func Info(message string, fields map[string]interface{}) {
    Log(INFO, message, fields)
}

func Warn(message string, fields map[string]interface{}) {
    Log(WARN, message, fields)
}

func Error(message string, fields map[string]interface{}) {
    Log(ERROR, message, fields)
}