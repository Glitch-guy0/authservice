package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/core/logger"
	"github.com/gin-gonic/gin"
)

// ShutdownManager manages graceful shutdown of the Gin server
type ShutdownManager struct {
	server         *Server
	logger         logger.Logger
	shutdownMu     sync.Mutex
	isShuttingDown bool
}

// NewShutdownManager creates a new shutdown manager for the server
func NewShutdownManager(server *Server) *ShutdownManager {
	return &ShutdownManager{
		server: server,
		logger: server.GetLogger(),
	}
}

// GracefulShutdown performs graceful shutdown of the HTTP server
func (sm *ShutdownManager) GracefulShutdown(ctx context.Context) error {
	sm.shutdownMu.Lock()
	defer sm.shutdownMu.Unlock()

	if sm.isShuttingDown {
		sm.logger.Info("Shutdown already in progress")
		return nil
	}

	sm.isShuttingDown = true
	sm.logger.Info("Starting graceful shutdown of HTTP server")

	// Update server health status
	sm.server.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusDegraded,
		Message:   "HTTP server shutting down",
		Timestamp: time.Now(),
	})

	// Create a context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	err := sm.server.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		sm.logger.Error("Error during HTTP server shutdown", "error", err)

		// Force close if graceful shutdown fails
		sm.logger.Warn("Forcing HTTP server close")
		closeErr := sm.server.httpServer.Close()
		if closeErr != nil {
			sm.logger.Error("Error during forced server close", "error", closeErr)
		}

		// Update health status
		sm.server.appCtx.UpdateHealthStatus("server", core.HealthStatus{
			Status:    core.StatusUnhealthy,
			Message:   fmt.Sprintf("Forced shutdown: %v", err),
			Timestamp: time.Now(),
		})

		return err
	}

	sm.logger.Info("HTTP server shutdown successfully")

	// Update final health status
	sm.server.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusHealthy,
		Message:   "HTTP server shutdown successfully",
		Timestamp: time.Now(),
	})

	return nil
}

// WaitForShutdown waits for shutdown signals and triggers graceful shutdown
func (sm *ShutdownManager) WaitForShutdown() {
	sm.logger.Info("Waiting for shutdown signals")

	// Wait for AppContext shutdown signal
	select {
	case <-sm.server.appCtx.GetShutdownChannel():
		sm.logger.Info("Shutdown signal received, initiating graceful shutdown")

		// Create context for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Perform graceful shutdown
		if err := sm.GracefulShutdown(ctx); err != nil {
			sm.logger.Error("Graceful shutdown failed", "error", err)
		}

	case <-sm.server.appCtx.Context().Done():
		sm.logger.Info("Application context canceled, shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sm.GracefulShutdown(ctx)
	}
}

// ConnectionManager manages active connections during shutdown
type ConnectionManager struct {
	activeConnections map[string]*http.Request
	mu                sync.RWMutex
	logger            logger.Logger
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(logger logger.Logger) *ConnectionManager {
	return &ConnectionManager{
		activeConnections: make(map[string]*http.Request),
		logger:            logger,
	}
}

// AddConnection tracks a new active connection
func (cm *ConnectionManager) AddConnection(id string, req *http.Request) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.activeConnections[id] = req
	cm.logger.Debug("Connection added", "connectionID", id, "total", len(cm.activeConnections))
}

// RemoveConnection removes a connection from tracking
func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.activeConnections, id)
	cm.logger.Debug("Connection removed", "connectionID", id, "total", len(cm.activeConnections))
}

// GetActiveConnections returns the number of active connections
func (cm *ConnectionManager) GetActiveConnections() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.activeConnections)
}

// WaitForConnections waits for all active connections to complete
func (cm *ConnectionManager) WaitForConnections(timeout time.Duration) error {
	cm.logger.Info("Waiting for active connections to complete", "count", cm.GetActiveConnections())

	start := time.Now()
	for {
		active := cm.GetActiveConnections()
		if active == 0 {
			cm.logger.Info("All connections completed")
			return nil
		}

		if time.Since(start) > timeout {
			cm.logger.Warn("Timeout waiting for connections", "remaining", active)
			return fmt.Errorf("timeout waiting for %d active connections", active)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// ServerMetrics tracks server performance metrics
type ServerMetrics struct {
	StartTime         time.Time     `json:"startTime"`
	TotalRequests     int64         `json:"totalRequests"`
	ActiveConnections int           `json:"activeConnections"`
	TotalConnections  int64         `json:"totalConnections"`
	Errors            int64         `json:"errors"`
	AvgResponseTime   time.Duration `json:"avgResponseTime"`
	mu                sync.RWMutex
}

// NewServerMetrics creates new server metrics
func NewServerMetrics() *ServerMetrics {
	return &ServerMetrics{
		StartTime: time.Now(),
	}
}

// RecordRequest records a request in metrics
func (sm *ServerMetrics) RecordRequest(duration time.Duration, isError bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.TotalRequests++
	if isError {
		sm.Errors++
	}

	// Calculate moving average for response time
	if sm.TotalRequests == 1 {
		sm.AvgResponseTime = duration
	} else {
		sm.AvgResponseTime = time.Duration(
			(int64(sm.AvgResponseTime)*int64(sm.TotalRequests-1) + int64(duration)) / int64(sm.TotalRequests),
		)
	}
}

// IncrementConnections increments active connection count
func (sm *ServerMetrics) IncrementConnections() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.ActiveConnections++
	sm.TotalConnections++
}

// DecrementConnections decrements active connection count
func (sm *ServerMetrics) DecrementConnections() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.ActiveConnections > 0 {
		sm.ActiveConnections--
	}
}

// GetMetrics returns current metrics
func (sm *ServerMetrics) GetMetrics() ServerMetrics {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return ServerMetrics{
		StartTime:         sm.StartTime,
		TotalRequests:     sm.TotalRequests,
		ActiveConnections: sm.ActiveConnections,
		TotalConnections:  sm.TotalConnections,
		Errors:            sm.Errors,
		AvgResponseTime:   sm.AvgResponseTime,
	}
}

// GetUptime returns server uptime
func (sm *ServerMetrics) GetUptime() time.Duration {
	return time.Since(sm.StartTime)
}

// Middleware for tracking metrics
func (sm *ServerMetrics) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment active connections
		sm.IncrementConnections()
		defer sm.DecrementConnections()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start)
		isError := c.Writer.Status() >= 400
		sm.RecordRequest(duration, isError)
	}
}

// HealthChecker implements server health checking
type HealthChecker struct {
	server *Server
}

// NewHealthChecker creates a new health checker for the server
func NewHealthChecker(server *Server) *HealthChecker {
	return &HealthChecker{
		server: server,
	}
}

// CheckHealth implements the core.HealthChecker interface
func (hc *HealthChecker) CheckHealth() core.HealthStatus {
	if hc.server.IsRunning() {
		return core.HealthStatus{
			Status:    core.StatusHealthy,
			Message:   "HTTP server is running",
			Timestamp: time.Now(),
		}
	}

	return core.HealthStatus{
		Status:    core.StatusUnhealthy,
		Message:   "HTTP server is not running",
		Timestamp: time.Now(),
	}
}

// Name returns the name of this health checker
func (hc *HealthChecker) Name() string {
	return "http-server"
}

// RegisterHealthChecker registers the server health checker with AppContext
func (s *Server) RegisterHealthChecker() {
	healthChecker := NewHealthChecker(s)
	s.appCtx.RegisterHealthChecker(healthChecker)
}

// EnhancedShutdown provides additional shutdown features
type EnhancedShutdown struct {
	*ShutdownManager
	connectionManager *ConnectionManager
	metrics           *ServerMetrics
}

// NewEnhancedShutdown creates an enhanced shutdown manager
func NewEnhancedShutdown(server *Server) *EnhancedShutdown {
	return &EnhancedShutdown{
		ShutdownManager:   NewShutdownManager(server),
		connectionManager: NewConnectionManager(server.GetLogger()),
		metrics:           NewServerMetrics(),
	}
}

// GracefulShutdownWithConnections performs graceful shutdown with connection management
func (es *EnhancedShutdown) GracefulShutdownWithConnections(ctx context.Context) error {
	es.logger.Info("Starting enhanced graceful shutdown")

	// Update server health status
	es.server.appCtx.UpdateHealthStatus("server", core.HealthStatus{
		Status:    core.StatusDegraded,
		Message:   "Enhanced shutdown in progress",
		Timestamp: time.Now(),
	})

	// Wait for active connections to complete (with timeout)
	connTimeout := 10 * time.Second
	if err := es.connectionManager.WaitForConnections(connTimeout); err != nil {
		es.logger.Warn("Not all connections completed gracefully", "error", err)
	}

	// Perform standard graceful shutdown
	return es.GracefulShutdown(ctx)
}

// GetConnectionManager returns the connection manager
func (es *EnhancedShutdown) GetConnectionManager() *ConnectionManager {
	return es.connectionManager
}

// GetMetrics returns the server metrics
func (es *EnhancedShutdown) GetMetrics() *ServerMetrics {
	return es.metrics
}
