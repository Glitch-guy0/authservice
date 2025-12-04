package logger

import (
	"log"
	"os"
)

// Logger interface defines the logging methods
type Logger interface {
	Create() *log.Logger
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Critical(msg string, args ...interface{})
}

// StandardLogger implements the Logger interface
type StandardLogger struct {
	logger *log.Logger
}

// New creates a new logger instance
func New() *StandardLogger {
	return &StandardLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Create returns the underlying logger instance
func (l *StandardLogger) Create() *log.Logger {
	return l.logger
}

// Info logs an informational message
func (l *StandardLogger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Warn logs a warning message
func (l *StandardLogger) Warn(msg string, args ...interface{}) {
	l.logger.Printf("[WARN] "+msg, args...)
}

// Error logs an error message
func (l *StandardLogger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// Critical logs a critical message
func (l *StandardLogger) Critical(msg string, args ...interface{}) {
	l.logger.Printf("[CRITICAL] "+msg, args...)
}
