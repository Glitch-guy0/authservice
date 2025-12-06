package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/Glitch-guy0/authService/src/modules/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware creates a request logging middleware
type LoggerMiddleware struct {
	logger logger.Logger
}

// NewLoggerMiddleware creates a new logger middleware instance
func NewLoggerMiddleware(logger logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

// Middleware returns the Gin middleware function
func (lm *LoggerMiddleware) Middleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Use structured logging with the application logger
			lm.logger.Info("HTTP Request",
				"method", param.Method,
				"path", param.Path,
				"status", param.StatusCode,
				"latency", param.Latency,
				"clientIP", param.ClientIP,
				"userAgent", param.Request.UserAgent(),
				"requestID", param.Keys["request_id"],
				"errorMessage", param.ErrorMessage,
			)
			return "" // Return empty string since we're using our own logger
		},
		Output: io.Discard, // Discard gin's default output
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists (from upstream)
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Set request ID in context and header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// DetailedLoggerMiddleware provides detailed request/response logging
type DetailedLoggerMiddleware struct {
	logger logger.Logger
}

// NewDetailedLoggerMiddleware creates a new detailed logger middleware
func NewDetailedLoggerMiddleware(logger logger.Logger) *DetailedLoggerMiddleware {
	return &DetailedLoggerMiddleware{
		logger: logger,
	}
}

// Middleware returns the detailed logging middleware
func (dlm *DetailedLoggerMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID
		requestID, _ := c.Get("request_id")
		if requestID == nil {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}

		// Record start time
		start := time.Now()

		// Read request body for logging (only for specific content types)
		var requestBody []byte
		if c.Request.Body != nil && shouldLogRequestBody(c) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Log request
		dlm.logger.Info("Request started",
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"clientIP", c.ClientIP(),
			"userAgent", c.Request.UserAgent(),
			"contentType", c.GetHeader("Content-Type"),
			"contentLength", c.Request.ContentLength,
		)

		// Log request body if present and appropriate
		if len(requestBody) > 0 {
			dlm.logger.Debug("Request body",
				"requestID", requestID,
				"body", string(requestBody),
			)
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		dlm.logger.Info("Request completed",
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
			"responseSize", c.Writer.Size(),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				dlm.logger.Error("Request error",
					"requestID", requestID,
					"error", err.Error(),
					"type", err.Type,
				)
			}
		}
	}
}

// shouldLogRequestBody determines if we should log the request body
func shouldLogRequestBody(c *gin.Context) bool {
	contentType := c.GetHeader("Content-Type")

	// Don't log binary data, large files, or sensitive data
	return contentType == "application/json" ||
		contentType == "application/xml" ||
		contentType == "text/plain" ||
		contentType == "application/x-www-form-urlencoded" &&
			c.Request.ContentLength < 1024*1024 // Less than 1MB
}

// SecurityLoggerMiddleware logs security-related events
type SecurityLoggerMiddleware struct {
	logger logger.Logger
}

// NewSecurityLoggerMiddleware creates a new security logger middleware
func NewSecurityLoggerMiddleware(logger logger.Logger) *SecurityLoggerMiddleware {
	return &SecurityLoggerMiddleware{
		logger: logger,
	}
}

// Middleware returns the security logging middleware
func (slm *SecurityLoggerMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("request_id")
		if requestID == nil {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}

		// Log suspicious activities
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Check for suspicious patterns
		if isSuspiciousUserAgent(userAgent) {
			slm.logger.Warn("Suspicious user agent detected",
				"requestID", requestID,
				"clientIP", clientIP,
				"userAgent", userAgent,
			)
		}

		// Log authentication attempts
		if c.Request.URL.Path == "/api/v1/auth/login" || c.Request.URL.Path == "/api/v1/auth/register" {
			slm.logger.Info("Authentication attempt",
				"requestID", requestID,
				"clientIP", clientIP,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)
		}

		c.Next()
	}
}

// isSuspiciousUserAgent checks for suspicious user agent patterns
func isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"sqlmap",
		"nikto",
		"nmap",
		"masscan",
		"zap",
		"burp",
		"scanner",
		"bot",
		"crawler",
	}

	for _, pattern := range suspiciousPatterns {
		if containsIgnoreCase(userAgent, pattern) {
			return true
		}
	}

	return false
}

// containsIgnoreCase performs case-insensitive substring check
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsIgnoreCaseMiddle(s, substr))
}

func containsIgnoreCaseMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// LoggerConfig holds configuration for logger middleware
type LoggerConfig struct {
	SkipPaths         []string `json:"skipPaths"`
	LogRequestBody    bool     `json:"logRequestBody"`
	LogResponseBody   bool     `json:"logResponseBody"`
	MaxBodySize       int64    `json:"maxBodySize"` // Maximum body size to log in bytes
	EnableSecurityLog bool     `json:"enableSecurityLog"`
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipPaths:         []string{"/health", "/metrics", "/ping"},
		LogRequestBody:    true,
		LogResponseBody:   false,       // Usually disabled for security
		MaxBodySize:       1024 * 1024, // 1MB
		EnableSecurityLog: true,
	}
}

// CreateLoggerMiddlewareChain creates a chain of logger middlewares
func CreateLoggerMiddlewareChain(logger logger.Logger, config LoggerConfig) []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc

	// Add request ID middleware first
	middlewares = append(middlewares, RequestIDMiddleware())

	// Add basic logger middleware
	basicLogger := NewLoggerMiddleware(logger)
	middlewares = append(middlewares, basicLogger.Middleware())

	// Add detailed logger if enabled
	if config.LogRequestBody || config.LogResponseBody {
		detailedLogger := NewDetailedLoggerMiddleware(logger)
		middlewares = append(middlewares, detailedLogger.Middleware())
	}

	// Add security logger if enabled
	if config.EnableSecurityLog {
		securityLogger := NewSecurityLoggerMiddleware(logger)
		middlewares = append(middlewares, securityLogger.Middleware())
	}

	return middlewares
}
