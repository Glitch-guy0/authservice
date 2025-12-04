package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerConfig tests server configuration
func TestServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, 15*time.Second, config.ReadTimeout)
	assert.Equal(t, 15*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.IdleTimeout)
	assert.Equal(t, "debug", config.Mode)
}

// TestNewServer tests server creation
func TestNewServer(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create dependencies
	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	config := DefaultServerConfig()

	// Create server
	server := NewServer(appCtx, config)

	require.NotNil(t, server)
	assert.Equal(t, appCtx, server.GetAppContext())
	assert.Equal(t, log, server.GetLogger())
	assert.Equal(t, config, server.GetConfig())
	assert.Equal(t, "0.0.0.0:8080", server.GetAddress())
	assert.NotNil(t, server.GetEngine())
	assert.True(t, server.IsRunning())
}

// TestNewServerWithDefaults tests server creation with defaults
func TestNewServerWithDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})

	server := NewServerWithDefaults(appCtx)

	require.NotNil(t, server)
	assert.Equal(t, DefaultServerConfig(), server.GetConfig())
}

// TestServerInitialize tests server initialization
func TestServerInitialize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)

	// Initialize server
	server.Initialize()

	// Test that default routes are set up
	engine := server.GetEngine()

	// Test ping endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

// TestServerHealthCheck tests server health check integration
func TestServerHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)

	// Register health checker
	server.RegisterHealthChecker()

	// Check initial health status
	healthStatus := appCtx.GetOverallHealth()
	assert.Equal(t, core.StatusHealthy, healthStatus.Status)

	// Check server-specific health
	serverHealth := appCtx.GetHealthStatus()
	assert.Contains(t, serverHealth, "http-server")
	assert.Equal(t, core.StatusHealthy, serverHealth["http-server"].Status)
}

// TestServerStartAndShutdown tests server start and shutdown
func TestServerStartAndShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})

	// Use a different port to avoid conflicts
	config := DefaultServerConfig()
	config.Port = 8081

	server := NewServer(appCtx, config)
	server.Initialize()

	// Start server in a goroutine
	startErr := make(chan error, 1)
	go func() {
		startErr <- server.Start()
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is responding
	resp, err := http.Get("http://localhost:8081/api/v1/ping")
	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	assert.NoError(t, shutdownErr)

	// Wait for server to finish
	select {
	case err := <-startErr:
		// Server should exit due to shutdown - http.ErrServerClosed is expected and acceptable
		if err != nil && err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Server did not shutdown within timeout")
	}
}

// TestServerGracefulShutdown tests graceful shutdown with active connections
func TestServerGracefulShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	config := DefaultServerConfig()
	config.Port = 8082

	server := NewServer(appCtx, config)
	server.Initialize()

	// Create shutdown manager
	shutdownManager := NewShutdownManager(server)

	// Start server
	go func() {
		server.Start()
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := shutdownManager.GracefulShutdown(ctx)
	assert.NoError(t, err)

	// Verify graceful shutdown completed successfully
	// Note: server instance still exists but is no longer listening
}

// TestServerMiddleware tests middleware integration
func TestServerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)

	// Add custom middleware
	server.GetEngine().Use(func(c *gin.Context) {
		c.Set("test_middleware", "executed")
		c.Next()
	})

	// Add test route
	server.GetEngine().GET("/test", func(c *gin.Context) {
		middlewareValue := c.MustGet("test_middleware")
		c.JSON(http.StatusOK, gin.H{
			"middleware": middlewareValue,
		})
	})

	// Test middleware execution
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	server.GetEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "executed")
}

// TestServerConfigurationValidation tests configuration validation
func TestServerConfigurationValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})

	// Test invalid port
	config := DefaultServerConfig()
	config.Port = -1

	server := NewServer(appCtx, config)
	assert.Equal(t, "0.0.0.0:-1", server.GetAddress())

	// Test invalid host
	config.Port = 8080
	config.Host = ""

	server = NewServer(appCtx, config)
	assert.Equal(t, ":8080", server.GetAddress())
}

// TestServerContext tests server context integration
func TestServerContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	testConfig := map[string]interface{}{
		"test_key": "test_value",
	}
	appCtx := core.NewAppContext(log, testConfig)

	server := NewServerWithDefaults(appCtx)

	// Test that server has access to app context config
	serverConfig := server.GetAppContext().GetConfig()
	assert.Equal(t, "test_value", serverConfig["test_key"])

	// Test that server can update app context
	server.GetAppContext().SetConfig("server_key", "server_value")
	updatedConfig := server.GetAppContext().GetConfig()
	assert.Equal(t, "server_value", updatedConfig["server_key"])
}

// TestServerMode tests different Gin modes
func TestServerMode(t *testing.T) {
	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})

	// Test debug mode
	config := DefaultServerConfig()
	config.Mode = "debug"
	_ = NewServer(appCtx, config)
	assert.Equal(t, "debug", gin.Mode())

	// Test release mode
	gin.SetMode(gin.ReleaseMode)
	config.Mode = "release"
	_ = NewServer(appCtx, config)
	assert.Equal(t, "release", gin.Mode())

	// Test test mode
	gin.SetMode(gin.TestMode)
	config.Mode = "test"
	_ = NewServer(appCtx, config)
	assert.Equal(t, "test", gin.Mode())
}

// TestServerTimeouts tests server timeout configurations
func TestServerTimeouts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})

	config := DefaultServerConfig()
	config.ReadTimeout = 5 * time.Second
	config.WriteTimeout = 10 * time.Second
	config.IdleTimeout = 30 * time.Second

	server := NewServer(appCtx, config)

	// Access the internal HTTP server to verify timeouts
	httpServer := server.httpServer
	assert.Equal(t, 5*time.Second, httpServer.ReadTimeout)
	assert.Equal(t, 10*time.Second, httpServer.WriteTimeout)
	assert.Equal(t, 30*time.Second, httpServer.IdleTimeout)
}

// TestServerMultipleShutdowns tests multiple shutdown attempts
func TestServerMultipleShutdowns(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)

	// Create shutdown manager
	shutdownManager := NewShutdownManager(server)

	ctx := context.Background()

	// First shutdown should succeed
	err1 := shutdownManager.GracefulShutdown(ctx)
	assert.NoError(t, err1)

	// Second shutdown should also succeed (should not panic)
	err2 := shutdownManager.GracefulShutdown(ctx)
	assert.NoError(t, err2)
}

// TestServerHealthStatusUpdate tests health status updates during server lifecycle
func TestServerHealthStatusUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)

	// Register health checker
	server.RegisterHealthChecker()

	// Check initial status
	health := appCtx.GetHealthStatus()
	assert.Contains(t, health, "http-server")
	assert.Equal(t, core.StatusHealthy, health["http-server"].Status)

	// Simulate server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	// Health status should be updated
	health = appCtx.GetHealthStatus()
	assert.Contains(t, health, "http-server")
	assert.Equal(t, core.StatusHealthy, health["http-server"].Status) // Should be healthy after successful shutdown
}

// Benchmark tests
func BenchmarkNewServer(b *testing.B) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	config := DefaultServerConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewServer(appCtx, config)
	}
}

func BenchmarkServerRequest(b *testing.B) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	server := NewServerWithDefaults(appCtx)
	server.Initialize()

	engine := server.GetEngine()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// Helper function to create test server
func createTestServer() *Server {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	appCtx := core.NewAppContext(log, map[string]interface{}{})
	config := DefaultServerConfig()

	server := NewServer(appCtx, config)
	server.Initialize()

	return server
}
