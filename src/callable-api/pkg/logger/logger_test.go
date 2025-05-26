// pkg/logger/logger_test.go
package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// Estrutura para analisar o JSON do log
type logJSON struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
}

// Função helper para capturar saída de stdout/stderr
func captureOutput(f func()) (stdout, stderr string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	os.Stdout = wOut
	os.Stderr = wErr
	
	f()
	
	wOut.Close()
	wErr.Close()
	
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)
	
	return bufOut.String(), bufErr.String()
}

func TestSetLevel(t *testing.T) {
	// Test each level setting
	testCases := []struct {
		levelStr  string
		levelEnum Level
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"invalid", INFO}, // default to INFO
	}

	for _, tc := range testCases {
		t.Run("SetLevel_"+tc.levelStr, func(t *testing.T) {
			SetLevel(tc.levelStr)
			if currentLevel != tc.levelEnum {
				t.Errorf("SetLevel(%s) should set currentLevel to %v, got %v", tc.levelStr, tc.levelEnum, currentLevel)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	// Test string representation of each level
	testCases := []struct {
		level  Level
		strVal string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, tc := range testCases {
		if tc.level.String() != tc.strVal {
			t.Errorf("Level %v should be %s, got %s", tc.level, tc.strVal, tc.level.String())
		}
	}
}

func TestLog(t *testing.T) {
	// Reset log level to capture everything
	SetLevel("debug")

	// Test all log levels
	testCases := []struct {
		level     Level
		levelName string
		toStderr  bool
	}{
		{DEBUG, "DEBUG", false},
		{INFO, "INFO", false},
		{WARN, "WARN", false},
		{ERROR, "ERROR", true},
	}

	for _, tc := range testCases {
		t.Run("Log_"+tc.levelName, func(t *testing.T) {
			fields := map[string]interface{}{"test": "value", "number": 42}
			message := "Test message for " + tc.levelName
			
			stdout, stderr := captureOutput(func() {
				Log(tc.level, message, fields)
			})
			
			var output string
			if tc.toStderr {
				output = stderr
			} else {
				output = stdout
			}
			
			if output == "" {
				t.Fatalf("Expected log output, got nothing for level %s", tc.levelName)
			}
			
			var logData logJSON
			err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logData)
			if err != nil {
				t.Fatalf("Failed to parse JSON log: %v\nOutput: %s", err, output)
			}
			
			if logData.Level != tc.levelName {
				t.Errorf("Expected level %s, got %s", tc.levelName, logData.Level)
			}
			
			if logData.Message != message {
				t.Errorf("Expected message %q, got %q", message, logData.Message)
			}
			
			if val, ok := logData.Fields["test"]; !ok || val != "value" {
				t.Errorf("Expected field test=value, got %v", logData.Fields["test"])
			}
			
			if val, ok := logData.Fields["number"]; !ok || val != float64(42) {
				t.Errorf("Expected field number=42, got %v", logData.Fields["number"])
			}
		})
	}
}

func TestLogLevelFiltering(t *testing.T) {
	// Set log level to WARN
	SetLevel("warn")
	
	// INFO logs should not appear
	stdout, _ := captureOutput(func() {
		Info("This should not appear", nil)
	})
	if stdout != "" {
		t.Errorf("INFO log appeared when level is WARN: %s", stdout)
	}
	
	// DEBUG logs should not appear
	stdout, _ = captureOutput(func() {
		Debug("This should not appear", nil)
	})
	if stdout != "" {
		t.Errorf("DEBUG log appeared when level is WARN: %s", stdout)
	}
	
	// WARN logs should appear
	stdout, _ = captureOutput(func() {
		Warn("This should appear", nil)
	})
	if stdout == "" {
		t.Errorf("WARN log did not appear when level is WARN")
	}
	
	// ERROR logs should appear
	_, stderr := captureOutput(func() {
		Error("This should appear", nil)
	})
	if stderr == "" {
		t.Errorf("ERROR log did not appear when level is WARN")
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Reset log level to capture everything
	SetLevel("debug")
	
	// Test each convenience function
	testCases := []struct {
		name     string
		logFunc  func(string, map[string]interface{})
		level    string
		toStderr bool
	}{
		{"Debug", Debug, "DEBUG", false},
		{"Info", Info, "INFO", false},
		{"Warn", Warn, "WARN", false},
		{"Error", Error, "ERROR", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name+"Function", func(t *testing.T) {
			message := "Test " + tc.name + " function"
			fields := map[string]interface{}{"func": tc.name}
			
			stdout, stderr := captureOutput(func() {
				tc.logFunc(message, fields)
			})
			
			var output string
			if tc.toStderr {
				output = stderr
			} else {
				output = stdout
			}
			
			if output == "" {
				t.Fatalf("Expected output from %s, got nothing", tc.name)
			}
			
			var logData logJSON
			err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logData)
			if err != nil {
				t.Fatalf("Failed to parse JSON log: %v\nOutput: %s", err, output)
			}
			
			if logData.Level != tc.level {
				t.Errorf("Expected level %s, got %s", tc.level, logData.Level)
			}
			
			if logData.Message != message {
				t.Errorf("Expected message %q, got %q", message, logData.Message)
			}
			
			if val, ok := logData.Fields["func"]; !ok || val != tc.name {
				t.Errorf("Expected field func=%s, got %v", tc.name, logData.Fields["func"])
			}
		})
	}
}

func TestNilFields(t *testing.T) {
	SetLevel("debug")
	
	stdout, _ := captureOutput(func() {
		Info("Message with nil fields", nil)
	})
	
	if stdout == "" {
		t.Fatal("Expected log output with nil fields, got nothing")
	}
	
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &logData)
	if err != nil {
		t.Fatalf("Failed to parse JSON log: %v\nOutput: %s", err, stdout)
	}
	
	// Check if fields key exists but is null or empty
	fields, exists := logData["fields"]
	if exists {
		// If fields exists, it should be null or empty
		fieldsMap, isMap := fields.(map[string]interface{})
		if isMap && len(fieldsMap) > 0 {
			t.Errorf("Expected empty or nil fields, got %v", fields)
		}
	}
}