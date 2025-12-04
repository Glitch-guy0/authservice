package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for CORS middleware
type CORSConfig struct {
	// Allowed origins (wildcards supported)
	AllowedOrigins []string `json:"allowedOrigins"`

	// Allowed methods
	AllowedMethods []string `json:"allowedMethods"`

	// Allowed headers
	AllowedHeaders []string `json:"allowedHeaders"`

	// Exposed headers (headers that can be exposed to the client)
	ExposedHeaders []string `json:"exposedHeaders"`

	// Allow credentials (cookies, authorization headers, etc.)
	AllowCredentials bool `json:"allowCredentials"`

	// Max age for preflight requests
	MaxAge time.Duration `json:"maxAge"`

	// Debug mode for logging CORS decisions
	Debug bool `json:"debug"`
}

// DefaultCORSConfig returns a secure default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",  // Common development port
			"http://localhost:8080",  // Common development port
			"https://localhost:3000", // HTTPS development
			"https://localhost:8080", // HTTPS development
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
			"X-API-Key",
			"X-Client-Version",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
			"X-Total-Count",
			"X-Page-Count",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 12 hours for preflight caching
		Debug:            false,
	}
}

// ProductionCORSConfig returns a production-ready CORS configuration
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: false,         // More secure for production
		MaxAge:           1 * time.Hour, // Shorter cache for production
		Debug:            false,
	}
}

// DevelopmentCORSConfig returns a permissive development CORS configuration
func DevelopmentCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"*", // Allow all origins in development
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*", // Allow all headers in development
		},
		ExposedHeaders: []string{
			"*", // Expose all headers in development
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour, // Long cache for development convenience
		Debug:            true,           // Enable debug logging in development
	}
}

// CORSMiddleware provides CORS functionality
type CORSMiddleware struct {
	config CORSConfig
	logger logger.Logger
}

// NewCORSMiddleware creates a new CORS middleware instance
func NewCORSMiddleware(logger logger.Logger, config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
		logger: logger,
	}
}

// NewCORSMiddlewareWithDefaults creates a CORS middleware with default config
func NewCORSMiddlewareWithDefaults(logger logger.Logger) *CORSMiddleware {
	return NewCORSMiddleware(logger, DefaultCORSConfig())
}

// Middleware returns the Gin CORS middleware function
func (cm *CORSMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		requestID, _ := c.Get("request_id")

		// Log CORS request if debug mode is enabled
		if cm.config.Debug {
			cm.logger.Debug("CORS request",
				"requestID", requestID,
				"origin", origin,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
			)
		}

		// Check if origin is allowed
		if cm.isOriginAllowed(origin) {
			cm.setCORSHeaders(c, origin)

			if cm.config.Debug {
				cm.logger.Debug("CORS allowed",
					"requestID", requestID,
					"origin", origin,
				)
			}
		} else {
			if cm.config.Debug {
				cm.logger.Warn("CORS denied - origin not allowed",
					"requestID", requestID,
					"origin", origin,
					"allowedOrigins", cm.config.AllowedOrigins,
				)
			}
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			cm.handlePreflight(c, origin, requestID)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func (cm *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		// Same-origin requests don't need CORS
		return true
	}

	for _, allowedOrigin := range cm.config.AllowedOrigins {
		if cm.matchOrigin(origin, allowedOrigin) {
			return true
		}
	}

	return false
}

// matchOrigin matches an origin against an allowed origin pattern
func (cm *CORSMiddleware) matchOrigin(origin, pattern string) bool {
	// Exact match
	if origin == pattern {
		return true
	}

	// Wildcard matching
	if pattern == "*" {
		return true
	}

	// Subdomain wildcard matching (e.g., *.example.com)
	if strings.HasPrefix(pattern, "*.") {
		domain := pattern[2:] // Remove "*."
		if strings.HasSuffix(origin, domain) {
			originParts := strings.Split(origin, ".")
			patternParts := strings.Split(domain, ".")

			// Check if the origin has at least the same number of parts as the pattern
			if len(originParts) >= len(patternParts) {
				// Compare the domain parts
				for i, part := range patternParts {
					if originParts[len(originParts)-len(patternParts)+i] != part {
						return false
					}
				}
				return true
			}
		}
	}

	return false
}

// setCORSHeaders sets the appropriate CORS headers
func (cm *CORSMiddleware) setCORSHeaders(c *gin.Context, origin string) {
	// Set Access-Control-Allow-Origin
	if cm.containsWildcard(cm.config.AllowedOrigins) {
		c.Header("Access-Control-Allow-Origin", "*")
	} else {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	// Set Access-Control-Allow-Methods
	if len(cm.config.AllowedMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", strings.Join(cm.config.AllowedMethods, ", "))
	}

	// Set Access-Control-Allow-Headers
	if len(cm.config.AllowedHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", strings.Join(cm.config.AllowedHeaders, ", "))
	}

	// Set Access-Control-Expose-Headers
	if len(cm.config.ExposedHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", strings.Join(cm.config.ExposedHeaders, ", "))
	}

	// Set Access-Control-Allow-Credentials
	if cm.config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// Set Access-Control-Max-Age
	if cm.config.MaxAge > 0 {
		c.Header("Access-Control-Max-Age", cm.config.MaxAge.String())
	}

	// Set Vary header for proper caching
	c.Header("Vary", "Origin")
}

// handlePreflight handles OPTIONS preflight requests
func (cm *CORSMiddleware) handlePreflight(c *gin.Context, origin string, requestID interface{}) {
	if cm.isOriginAllowed(origin) {
		cm.setCORSHeaders(c, origin)

		if cm.config.Debug {
			cm.logger.Debug("CORS preflight request allowed",
				"requestID", requestID,
				"origin", origin,
				"method", c.Request.Method,
			)
		}

		c.Status(http.StatusNoContent)
	} else {
		if cm.config.Debug {
			cm.logger.Warn("CORS preflight request denied",
				"requestID", requestID,
				"origin", origin,
				"method", c.Request.Method,
			)
		}

		c.Status(http.StatusForbidden)
	}
}

// containsWildcard checks if any origin pattern contains a wildcard
func (cm *CORSMiddleware) containsWildcard(origins []string) bool {
	for _, origin := range origins {
		if origin == "*" || strings.HasPrefix(origin, "*.") {
			return true
		}
	}
	return false
}

// CORSSecurityMiddleware adds additional security features for CORS
type CORSSecurityMiddleware struct {
	*CORSMiddleware
	maxRequestsPerMinute int
	requestCounts        map[string]int
	lastReset            time.Time
}

// NewCORSSecurityMiddleware creates a CORS middleware with security features
func NewCORSSecurityMiddleware(logger logger.Logger, config CORSConfig, maxRequestsPerMinute int) *CORSSecurityMiddleware {
	base := NewCORSMiddleware(logger, config)

	return &CORSSecurityMiddleware{
		CORSMiddleware:       base,
		maxRequestsPerMinute: maxRequestsPerMinute,
		requestCounts:        make(map[string]int),
		lastReset:            time.Now(),
	}
}

// Middleware returns the security-enhanced CORS middleware
func (csm *CORSSecurityMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rate limiting for preflight requests
		if c.Request.Method == http.MethodOptions {
			origin := c.Request.Header.Get("Origin")

			if !csm.checkRateLimit(origin) {
				csm.logger.Warn("CORS preflight rate limit exceeded",
					"origin", origin,
					"clientIP", c.ClientIP(),
				)

				c.Status(http.StatusTooManyRequests)
				return
			}
		}

		// Call base CORS middleware
		csm.CORSMiddleware.Middleware()(c)
	}
}

// checkRateLimit checks if the origin has exceeded the rate limit
func (csm *CORSSecurityMiddleware) checkRateLimit(origin string) bool {
	now := time.Now()

	// Reset counters if more than a minute has passed
	if now.Sub(csm.lastReset) > time.Minute {
		csm.requestCounts = make(map[string]int)
		csm.lastReset = now
	}

	// Increment counter for this origin
	csm.requestCounts[origin]++

	// Check if limit exceeded
	return csm.requestCounts[origin] <= csm.maxRequestsPerMinute
}

// CreateCORSMiddlewareChain creates a CORS middleware chain based on environment
func CreateCORSMiddlewareChain(logger logger.Logger, isProduction bool, allowedOrigins []string) gin.HandlerFunc {
	var config CORSConfig

	if isProduction {
		if len(allowedOrigins) == 0 {
			// Default production origins if none provided
			allowedOrigins = []string{
				"https://yourdomain.com",
			}
		}
		config = ProductionCORSConfig(allowedOrigins)
	} else {
		config = DevelopmentCORSConfig()
	}

	// Add security features in production
	if isProduction {
		security := NewCORSSecurityMiddleware(logger, config, 60) // 60 preflight requests per minute
		return security.Middleware()
	}

	cors := NewCORSMiddleware(logger, config)
	return cors.Middleware()
}

// CORSMiddlewareConfigBuilder helps build CORS configurations
type CORSMiddlewareConfigBuilder struct {
	config CORSConfig
}

// NewCORSMiddlewareConfigBuilder creates a new CORS config builder
func NewCORSMiddlewareConfigBuilder() *CORSMiddlewareConfigBuilder {
	return &CORSMiddlewareConfigBuilder{
		config: DefaultCORSConfig(),
	}
}

// WithOrigins sets allowed origins
func (b *CORSMiddlewareConfigBuilder) WithOrigins(origins ...string) *CORSMiddlewareConfigBuilder {
	b.config.AllowedOrigins = origins
	return b
}

// WithMethods sets allowed methods
func (b *CORSMiddlewareConfigBuilder) WithMethods(methods ...string) *CORSMiddlewareConfigBuilder {
	b.config.AllowedMethods = methods
	return b
}

// WithHeaders sets allowed headers
func (b *CORSMiddlewareConfigBuilder) WithHeaders(headers ...string) *CORSMiddlewareConfigBuilder {
	b.config.AllowedHeaders = headers
	return b
}

// WithExposedHeaders sets exposed headers
func (b *CORSMiddlewareConfigBuilder) WithExposedHeaders(headers ...string) *CORSMiddlewareConfigBuilder {
	b.config.ExposedHeaders = headers
	return b
}

// WithCredentials sets whether credentials are allowed
func (b *CORSMiddlewareConfigBuilder) WithCredentials(allow bool) *CORSMiddlewareConfigBuilder {
	b.config.AllowCredentials = allow
	return b
}

// WithMaxAge sets the max age for preflight requests
func (b *CORSMiddlewareConfigBuilder) WithMaxAge(maxAge time.Duration) *CORSMiddlewareConfigBuilder {
	b.config.MaxAge = maxAge
	return b
}

// WithDebug enables or disables debug mode
func (b *CORSMiddlewareConfigBuilder) WithDebug(debug bool) *CORSMiddlewareConfigBuilder {
	b.config.Debug = debug
	return b
}

// Build creates the final CORS configuration
func (b *CORSMiddlewareConfigBuilder) Build() CORSConfig {
	return b.config
}

// BuildMiddleware creates a CORS middleware with the built configuration
func (b *CORSMiddlewareConfigBuilder) BuildMiddleware(logger logger.Logger) gin.HandlerFunc {
	config := b.Build()
	cors := NewCORSMiddleware(logger, config)
	return cors.Middleware()
}
