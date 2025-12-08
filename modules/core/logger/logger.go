package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger interface defines the logging methods
type Logger interface {
	Create() *logrus.Logger
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// StandardLogger implements the Logger interface
type StandardLogger struct {
	logger *logrus.Logger
}

// New creates a new logger instance
func New() *StandardLogger {
	config := DefaultConfig()
	return &StandardLogger{
		logger: ConfigureLogger(config),
	}
}

// NewWithConfig creates a new logger instance with custom configuration
func NewWithConfig(config *LogConfig) *StandardLogger {
	return &StandardLogger{
		logger: ConfigureLogger(config),
	}
}

// Create returns the underlying logger instance
func (l *StandardLogger) Create() *logrus.Logger {
	return l.logger
}

// Info logs an informational message
func (l *StandardLogger) Info(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Warn logs a warning message
func (l *StandardLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

// Error logs an error message
func (l *StandardLogger) Error(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

// Debug logs a debug message
func (l *StandardLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

// Fatal logs a fatal message and exits
func (l *StandardLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

// Panic logs a panic message and panics
func (l *StandardLogger) Panic(msg string, args ...interface{}) {
	l.logger.Panicf(msg, args...)
}

// WithField adds a field to the logger
func (l *StandardLogger) WithField(key string, value interface{}) Logger {
	return &EntryLogger{
		entry: l.logger.WithField(key, value),
	}
}

// WithFields adds multiple fields to the logger
func (l *StandardLogger) WithFields(fields map[string]interface{}) Logger {
	return &EntryLogger{
		entry: l.logger.WithFields(fields),
	}
}

// EntryLogger implements Logger interface using logrus.Entry
type EntryLogger struct {
	entry *logrus.Entry
}

// Create returns the underlying logger entry
func (l *EntryLogger) Create() *logrus.Logger {
	return l.entry.Logger
}

// Info logs an informational message
func (l *EntryLogger) Info(msg string, args ...interface{}) {
	l.entry.Infof(msg, args...)
}

// Warn logs a warning message
func (l *EntryLogger) Warn(msg string, args ...interface{}) {
	l.entry.Warnf(msg, args...)
}

// Error logs an error message
func (l *EntryLogger) Error(msg string, args ...interface{}) {
	l.entry.Errorf(msg, args...)
}

// Debug logs a debug message
func (l *EntryLogger) Debug(msg string, args ...interface{}) {
	l.entry.Debugf(msg, args...)
}

// Fatal logs a fatal message and exits
func (l *EntryLogger) Fatal(msg string, args ...interface{}) {
	l.entry.Fatalf(msg, args...)
}

// Panic logs a panic message and panics
func (l *EntryLogger) Panic(msg string, args ...interface{}) {
	l.entry.Panicf(msg, args...)
}

// WithField adds a field to the logger
func (l *EntryLogger) WithField(key string, value interface{}) Logger {
	return &EntryLogger{
		entry: l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the logger
func (l *EntryLogger) WithFields(fields map[string]interface{}) Logger {
	return &EntryLogger{
		entry: l.entry.WithFields(fields),
	}
}
