package logger

import (
	"context"

	"github.com/google/uuid"
)

// ContextKey is used for context values
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// CorrelationIDKey is the context key for correlation ID
	CorrelationIDKey ContextKey = "correlation_id"
)

// ContextLogger extends Logger with context support
type ContextLogger interface {
	Logger
	WithContext(ctx context.Context) Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger
	WithCorrelationID(correlationID string) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// ContextualLogger implements ContextLogger
type ContextualLogger struct {
	*StandardLogger
	requestID     string
	userID        string
	correlationID string
	fields        map[string]interface{}
}

// NewContextual creates a new contextual logger
func NewContextual() *ContextualLogger {
	return &ContextualLogger{
		StandardLogger: New(),
		fields:         make(map[string]interface{}),
	}
}

// WithContext returns a logger with context values
func (l *ContextualLogger) WithContext(ctx context.Context) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Extract context values
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		newLogger.requestID = requestID
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		newLogger.userID = userID
	}
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		newLogger.correlationID = correlationID
	}

	return newLogger
}

// WithRequestID adds request ID to the logger
func (l *ContextualLogger) WithRequestID(requestID string) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		requestID:      requestID,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithUserID adds user ID to the logger
func (l *ContextualLogger) WithUserID(userID string) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		userID:         userID,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithCorrelationID adds correlation ID to the logger
func (l *ContextualLogger) WithCorrelationID(correlationID string) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		correlationID:  correlationID,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithField adds a field to the logger
func (l *ContextualLogger) WithField(key string, value interface{}) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		requestID:      l.requestID,
		userID:         l.userID,
		correlationID:  l.correlationID,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value

	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *ContextualLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := &ContextualLogger{
		StandardLogger: l.StandardLogger,
		requestID:      l.requestID,
		userID:         l.userID,
		correlationID:  l.correlationID,
		fields:         make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// WithRequestContext creates a context with request ID
func WithRequestContext(ctx context.Context) context.Context {
	requestID := GenerateRequestID()
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return correlationID
	}
	return ""
}
