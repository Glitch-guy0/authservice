package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
	appCtx     *core.AppContext
	logger     logger.Logger
	config     ServerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
	Mode         string        `json:"mode"` // "debug", "release", "test"
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:         "0.0.0.0",
		Port:         8080,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Mode:         "debug",
	}
}

// NewServer creates a new HTTP server instance
func NewServer(appCtx *core.AppContext, config ServerConfig) *Server {
	// Set Gin mode
	gin.SetMode(config.Mode)

	// Create Gin engine
	engine := gin.New()

	// Create server instance
	server := &Server{
		engine: engine,
		appCtx: appCtx,
		logger: appCtx.GetLogger(),
		config: config,
	}

	// Setup HTTP server
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      engine,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// Register graceful shutdown handler
	server.registerShutdownHandler()

	return server
}

// NewServerWithDefaults creates a server with default configuration
func NewServerWithDefaults(appCtx *core.AppContext) *Server {
	return NewServer(appCtx, DefaultServerConfig())
}

// NewServerFromConfig creates a server with configuration from app context
func NewServerFromConfig(appCtx *core.AppContext) *Server {
	config := appCtx.GetConfig()

	// Extract server configuration from the app context
	serverConfig := ServerConfig{
		Host:         "0.0.0.0", // default host
		Port:         8080,      // default port
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Mode:         "debug",
	}

	// Try to get server config from the configuration
	if serverConfigMap, ok := config["server"].(map[string]interface{}); ok {
		if host, ok := serverConfigMap["host"].(string); ok {
			serverConfig.Host = host
		}
		if port, ok := serverConfigMap["port"].(int); ok {
			serverConfig.Port = port
		}
		if mode, ok := serverConfigMap["debug"].(bool); ok {
			if mode {
				serverConfig.Mode = "debug"
			} else {
				serverConfig.Mode = "release"
			}
		}

		// Handle timeouts if present
		if timeout, ok := serverConfigMap["timeout"].(map[string]interface{}); ok {
			if read, ok := timeout["read"].(int); ok {
				serverConfig.ReadTimeout = time.Duration(read) * time.Second
			}
			if write, ok := timeout["write"].(int); ok {
				serverConfig.WriteTimeout = time.Duration(write) * time.Second
			}
			if idle, ok := timeout["idle"].(int); ok {
				serverConfig.IdleTimeout = time.Duration(idle) * time.Second
			}
		}
	}

	return NewServer(appCtx, serverConfig)
}

// GetEngine returns the Gin engine for route registration
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// GetAppContext returns the application context
func (s *Server) GetAppContext() *core.AppContext {
	return s.appCtx
}

// GetLogger returns the logger instance
func (s *Server) GetLogger() logger.Logger {
	return s.logger
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() ServerConfig {
	return s.config
}

// GetAddress returns the server address
func (s *Server) GetAddress() string {
	return s.httpServer.Addr
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		"address", s.GetAddress(),
		"host", s.config.Host,
		"port", s.config.Port,
		"mode", s.config.Mode,
		"readTimeout", s.config.ReadTimeout,
		"writeTimeout", s.config.WriteTimeout,
		"idleTimeout", s.config.IdleTimeout,
	)

	// Update health status
	s.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusHealthy,
		Message:   "HTTP server started",
		Timestamp: time.Now(),
	})

	// Start the server
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("Failed to start HTTP server", "error", err)

		// Update health status
		s.appCtx.UpdateHealthStatus("server", core.HealthStatus{
			Status:    core.StatusUnhealthy,
			Message:   fmt.Sprintf("Failed to start: %v", err),
			Timestamp: time.Now(),
		})

		return err
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")

	// Update health status
	s.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusDegraded,
		Message:   "HTTP server shutting down",
		Timestamp: time.Now(),
	})

	// Shutdown the HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to shutdown HTTP server gracefully", "error", err)

		// Update health status
		s.appCtx.UpdateHealthStatus("server", core.HealthStatus{
			Status:    core.StatusUnhealthy,
			Message:   fmt.Sprintf("Shutdown failed: %v", err),
			Timestamp: time.Now(),
		})

		return err
	}

	s.logger.Info("HTTP server shutdown successfully")

	// Update health status
	s.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusHealthy,
		Message:   "HTTP server shutdown successfully",
		Timestamp: time.Now(),
	})

	return nil
}

// registerShutdownHandler registers the server shutdown handler with AppContext
func (s *Server) registerShutdownHandler() {
	s.appCtx.RegisterShutdownHandler("gin-server", func(ctx context.Context) error {
		return s.Shutdown(ctx)
	}, 30*time.Second, 1) // High priority, 30s timeout
}

// IsRunning checks if the server is currently running
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}

// SetupMiddleware sets up the default middleware
func (s *Server) SetupMiddleware() {
	// Recovery middleware
	s.engine.Use(gin.Recovery())

	// Request logging middleware (will be implemented in T023)
	// s.engine.Use(s.requestLoggerMiddleware())

	// CORS middleware (will be implemented in T026)
	// s.engine.Use(s.corsMiddleware())

	s.logger.Info("Default middleware setup completed")
}

// SetupRoutes sets up the default routes
func (s *Server) SetupRoutes() {
	// Health check endpoint (will be implemented in T028-T031)
	// s.engine.GET("/health", s.healthHandler)

	// API versioning base path
	v1 := s.engine.Group("/api/v1")
	{
		// API routes will be added here
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"time":    time.Now().UTC(),
			})
		})
	}

	s.logger.Info("Default routes setup completed")
}

// Initialize sets up the server with middleware and routes
func (s *Server) Initialize() {
	s.logger.Info("Initializing HTTP server")

	// Setup middleware
	s.SetupMiddleware()

	// Setup routes
	s.SetupRoutes()

	s.logger.Info("HTTP server initialization completed")
}
