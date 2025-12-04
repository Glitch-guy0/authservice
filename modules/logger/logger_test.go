package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLoggerMethods(t *testing.T) {
	// Capture output
	var buf bytes.Buffer

	// Create logger with custom output
	logger := &StandardLogger{
		logger: log.New(&buf, "", log.LstdFlags),
	}

	tests := []struct {
		name     string
		method   func(string, ...interface{})
		message  string
		expected string
	}{
		{
			name:     "Info",
			method:   logger.Info,
			message:  "test info message",
			expected: "[INFO] test info message",
		},
		{
			name:     "Warn",
			method:   logger.Warn,
			message:  "test warning message",
			expected: "[WARN] test warning message",
		},
		{
			name:     "Error",
			method:   logger.Error,
			message:  "test error message",
			expected: "[ERROR] test error message",
		},
		{
			name:     "Critical",
			method:   logger.Critical,
			message:  "test critical message",
			expected: "[CRITICAL] test critical message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.method(tt.message)
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	logger := New()
	createdLogger := logger.Create()
	if createdLogger == nil {
		t.Error("Create() should return a non-nil logger")
	}
}

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Error("New() should return a non-nil logger")
	}
	if logger.logger == nil {
		t.Error("New() should initialize the internal logger")
	}
}
