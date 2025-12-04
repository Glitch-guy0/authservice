package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerMethods(t *testing.T) {
	// Capture output
	var buf bytes.Buffer

	// Create logger with custom output
	logger := &StandardLogger{
		logger: logrus.New(),
	}
	logger.logger.SetOutput(&buf)
	logger.logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableQuote:     true,
	})

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
			expected: "level=info msg=test info message",
		},
		{
			name:     "Warn",
			method:   logger.Warn,
			message:  "test warning message",
			expected: "level=warning msg=test warning message",
		},
		{
			name:     "Error",
			method:   logger.Error,
			message:  "test error message",
			expected: "level=error msg=test error message",
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
		return
	}
	if logger.logger == nil {
		t.Error("New() should initialize the internal logger")
	}
}
