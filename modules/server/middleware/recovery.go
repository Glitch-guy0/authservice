package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware provides custom recovery with structured logging
type RecoveryMiddleware struct {
	logger logger.Logger
	config RecoveryConfig
}

// RecoveryConfig holds configuration for recovery middleware
type RecoveryConfig struct {
	EnableStackTrace bool     `json:"enableStackTrace"`
	SkipPaths        []string `json:"skipPaths"`
	StackSize        int      `json:"stackSize"` // Maximum stack trace lines to log
}

// DefaultRecoveryConfig returns default recovery configuration
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		EnableStackTrace: true,
		SkipPaths:        []string{"/health", "/metrics", "/ping"},
		StackSize:        50, // Log up to 50 stack frames
	}
}

// NewRecoveryMiddleware creates a new recovery middleware instance
func NewRecoveryMiddleware(logger logger.Logger, config RecoveryConfig) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
		config: config,
	}
}

// NewRecoveryMiddlewareWithDefaults creates a recovery middleware with default config
func NewRecoveryMiddlewareWithDefaults(logger logger.Logger) *RecoveryMiddleware {
	return NewRecoveryMiddleware(logger, DefaultRecoveryConfig())
}

// Middleware returns the Gin recovery middleware function
func (rm *RecoveryMiddleware) Middleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Get request ID if available
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "unknown"
		}

		// Log the panic with structured information
		rm.logPanic(c, recovered, requestID)

		// Send appropriate response
		rm.sendErrorResponse(c, recovered, requestID)
	})
}

// logPanic logs detailed information about the panic
func (rm *RecoveryMiddleware) logPanic(c *gin.Context, recovered interface{}, requestID interface{}) {
	// Basic panic information
	rm.logger.Error("Panic recovered",
		"requestID", requestID,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"query", c.Request.URL.RawQuery,
		"clientIP", c.ClientIP(),
		"userAgent", c.Request.UserAgent(),
		"panic", fmt.Sprintf("%v", recovered),
		"contentType", c.GetHeader("Content-Type"),
	)

	// Log stack trace if enabled
	if rm.config.EnableStackTrace {
		stackTrace := string(debug.Stack())

		// Limit stack trace size
		if rm.config.StackSize > 0 {
			lines := strings.Split(stackTrace, "\n")
			if len(lines) > rm.config.StackSize {
				stackTrace = strings.Join(lines[:rm.config.StackSize], "\n")
			}
		}

		rm.logger.Error("Stack trace",
			"requestID", requestID,
			"stack", stackTrace,
		)
	}

	// Log request headers (sanitized)
	rm.logRequestHeaders(c, requestID)

	// Log request body if available and safe
	rm.logRequestBody(c, requestID)
}

// logRequestHeaders logs request headers (sanitized)
func (rm *RecoveryMiddleware) logRequestHeaders(c *gin.Context, requestID interface{}) {
	headers := make(map[string]string)

	for key, values := range c.Request.Header {
		// Sanitize sensitive headers
		if isSensitiveHeader(key) {
			headers[key] = "[REDACTED]"
		} else {
			headers[key] = strings.Join(values, ", ")
		}
	}

	rm.logger.Debug("Request headers during panic",
		"requestID", requestID,
		"headers", headers,
	)
}

// logRequestBody logs request body if safe to do so
func (rm *RecoveryMiddleware) logRequestBody(c *gin.Context, requestID interface{}) {
	// Only log body for safe content types and reasonable sizes
	contentType := c.GetHeader("Content-Type")
	contentLength := c.Request.ContentLength

	if shouldLogBodyForPanic(contentType, contentLength) {
		if c.Request.Body != nil {
			// Note: Body might have been read already, so we can't always log it
			rm.logger.Debug("Request body during panic",
				"requestID", requestID,
				"contentType", contentType,
				"contentLength", contentLength,
				"note", "Body may have been consumed by other middleware",
			)
		}
	}
}

// sendErrorResponse sends an appropriate error response to the client
func (rm *RecoveryMiddleware) sendErrorResponse(c *gin.Context, recovered interface{}, requestID interface{}) {
	// Don't send panic details in production mode
	if gin.Mode() == gin.ReleaseMode {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":      "INTERNAL_SERVER_ERROR",
				"message":   "An internal error occurred",
				"requestID": requestID,
				"timestamp": "now", // Would use actual timestamp
			},
		})
		return
	}

	// In debug mode, provide more information
	errorResponse := gin.H{
		"error": gin.H{
			"code":      "INTERNAL_SERVER_ERROR",
			"message":   "An internal error occurred",
			"requestID": requestID,
			"timestamp": "now", // Would use actual timestamp
		},
	}

	// Add panic details in debug mode
	if gin.Mode() == gin.DebugMode {
		errorResponse["debug"] = gin.H{
			"panic": fmt.Sprintf("%v", recovered),
		}

		// Add stack trace if enabled and in debug mode
		if rm.config.EnableStackTrace {
			stackTrace := string(debug.Stack())
			if rm.config.StackSize > 0 {
				lines := strings.Split(stackTrace, "\n")
				if len(lines) > rm.config.StackSize {
					stackTrace = strings.Join(lines[:rm.config.StackSize], "\n")
				}
			}
			errorResponse["debug"].(gin.H)["stack"] = stackTrace
		}
	}

	c.JSON(http.StatusInternalServerError, errorResponse)
}

// isSensitiveHeader checks if a header contains sensitive information
func isSensitiveHeader(key string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"set-cookie",
		"x-api-key",
		"x-auth-token",
		"x-session-token",
		"x-csrf-token",
		"x-forwarded-for",
		"x-real-ip",
	}

	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveHeaders {
		if keyLower == strings.ToLower(sensitive) {
			return true
		}
	}

	return false
}

// shouldLogBodyForPanic determines if it's safe to log request body during panic
func shouldLogBodyForPanic(contentType string, contentLength int64) bool {
	// Don't log if content is too large
	if contentLength > 1024*1024 { // 1MB
		return false
	}

	// Only log safe content types
	safeTypes := []string{
		"application/json",
		"application/xml",
		"text/plain",
		"application/x-www-form-urlencoded",
	}

	for _, safeType := range safeTypes {
		if strings.Contains(contentType, safeType) {
			return true
		}
	}

	return false
}

// shouldSkipPath checks if the path should be skipped by recovery middleware
func (rm *RecoveryMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range rm.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// EnhancedRecoveryMiddleware provides additional recovery features
type EnhancedRecoveryMiddleware struct {
	*RecoveryMiddleware
	metrics PanicMetrics
}

// PanicMetrics tracks panic statistics
type PanicMetrics struct {
	PanicCount    int64 `json:"panicCount"`
	LastPanicTime int64 `json:"lastPanicTime"` // Unix timestamp
}

// NewEnhancedRecoveryMiddleware creates an enhanced recovery middleware
func NewEnhancedRecoveryMiddleware(logger logger.Logger, config RecoveryConfig) *EnhancedRecoveryMiddleware {
	base := NewRecoveryMiddleware(logger, config)

	return &EnhancedRecoveryMiddleware{
		RecoveryMiddleware: base,
		metrics:            PanicMetrics{},
	}
}

// Middleware returns the enhanced recovery middleware function
func (erm *EnhancedRecoveryMiddleware) Middleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Update metrics
		erm.metrics.PanicCount++
		// erm.metrics.LastPanicTime = time.Now().Unix() // Would import time

		// Get request ID
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "unknown"
		}

		// Log enhanced panic information
		erm.logger.Error("Enhanced panic recovery",
			"requestID", requestID,
			"panicCount", erm.metrics.PanicCount,
			"panic", fmt.Sprintf("%v", recovered),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"clientIP", c.ClientIP(),
		)

		// Call base recovery
		erm.logPanic(c, recovered, requestID)
		erm.sendErrorResponse(c, recovered, requestID)
	})
}

// GetMetrics returns current panic metrics
func (erm *EnhancedRecoveryMiddleware) GetMetrics() PanicMetrics {
	return erm.metrics
}

// ResetMetrics resets panic metrics
func (erm *EnhancedRecoveryMiddleware) ResetMetrics() {
	erm.metrics = PanicMetrics{}
}

// CreateRecoveryMiddlewareChain creates a recovery middleware with appropriate configuration
func CreateRecoveryMiddlewareChain(logger logger.Logger, isProduction bool) gin.HandlerFunc {
	config := DefaultRecoveryConfig()

	if isProduction {
		config.EnableStackTrace = false // Disable stack traces in production
	}

	recovery := NewRecoveryMiddleware(logger, config)
	return recovery.Middleware()
}
