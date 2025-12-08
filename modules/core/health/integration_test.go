package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/core/logger"
)

func TestHealthIntegration_FullServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create application context
	log := logger.New()
	appCtx := core.NewAppContextWithDefaults(log)

	// Create test router with health endpoints
	router := gin.New()

	// Add health endpoints like the real server does
	healthHandler := NewHealthHandler(appCtx)
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/live", healthHandler.LivenessProbe)
	router.GET("/health/ready", healthHandler.ReadinessProbe)

	// Test main health endpoint
	t.Run("Main Health Endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, StatusHealthy, response.Status)
		assert.NotZero(t, response.Timestamp)
		assert.Equal(t, "1.0.0", response.Version.Version)
		assert.NotEmpty(t, response.Checks)
		assert.NotEmpty(t, response.Uptime)

		// Verify default checkers
		checkNames := make(map[string]bool)
		for _, check := range response.Checks {
			checkNames[check.Name] = true
			assert.NotEmpty(t, check.Duration)
			assert.NotZero(t, check.Timestamp)
		}

		assert.True(t, checkNames["server"])
		assert.True(t, checkNames["database"])
		assert.True(t, checkNames["logger"])
	})

	// Test liveness probe
	t.Run("Liveness Probe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/live", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "ok", response["status"])
		assert.NotNil(t, response["timestamp"])
	})

	// Test readiness probe
	t.Run("Readiness Probe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "ready", response["status"])
		assert.NotNil(t, response["timestamp"])
	})
}

func TestHealthIntegration_WithCustomCheckers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create application context
	log := logger.New()
	appCtx := core.NewAppContextWithDefaults(log)

	// Create health handler
	healthHandler := NewHealthHandler(appCtx)

	// Add custom checkers
	slowChecker := NewBasicHealthChecker("slow-service", func(ctx context.Context) Check {
		// Simulate slow service
		time.Sleep(100 * time.Millisecond)
		return Check{
			Name:      "slow-service",
			Status:    StatusDegraded,
			Message:   "Service responding slowly",
			Timestamp: time.Now(),
		}
	})

	failingChecker := NewBasicHealthChecker("failing-service", func(ctx context.Context) Check {
		return Check{
			Name:      "failing-service",
			Status:    StatusUnhealthy,
			Message:   "Service is down",
			Timestamp: time.Now(),
		}
	})

	healthHandler.healthService.RegisterChecker(slowChecker)
	healthHandler.healthService.RegisterChecker(failingChecker)

	// Create router
	router := gin.New()
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/live", healthHandler.LivenessProbe)
	router.GET("/health/ready", healthHandler.ReadinessProbe)

	// Test health endpoint with custom checkers
	t.Run("Health with Custom Checkers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be unhealthy due to failing checker
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, StatusUnhealthy, response.Status)
		assert.Equal(t, 5, len(response.Checks)) // 3 default + 2 custom

		// Find custom checkers
		foundSlow := false
		foundFailing := false
		for _, check := range response.Checks {
			if check.Name == "slow-service" {
				foundSlow = true
				assert.Equal(t, StatusDegraded, check.Status)
				assert.Equal(t, "Service responding slowly", check.Message)
				assert.Contains(t, check.Duration, "ms")
			}
			if check.Name == "failing-service" {
				foundFailing = true
				assert.Equal(t, StatusUnhealthy, check.Status)
				assert.Equal(t, "Service is down", check.Message)
			}
		}
		assert.True(t, foundSlow)
		assert.True(t, foundFailing)
	})

	// Test readiness probe should be not ready
	t.Run("Readiness with Failing Checker", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "not ready", response["status"])
	})
}

func TestHealthIntegration_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create application context
	log := logger.New()
	appCtx := core.NewAppContextWithDefaults(log)

	// Create health handler
	healthHandler := NewHealthHandler(appCtx)

	// Add a checker that simulates some work
	workChecker := NewBasicHealthChecker("work-service", func(ctx context.Context) Check {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return Check{
			Name:      "work-service",
			Status:    StatusHealthy,
			Message:   "Work completed",
			Timestamp: time.Now(),
		}
	})

	healthHandler.healthService.RegisterChecker(workChecker)

	// Create router
	router := gin.New()
	router.GET("/health", healthHandler.HealthCheck)

	// Test concurrent requests
	t.Run("Concurrent Health Checks", func(t *testing.T) {
		const numRequests = 10
		results := make(chan int, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				req, _ := http.NewRequest("GET", "/health", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				results <- w.Code
			}()
		}

		// Collect results
		successCount := 0
		for i := 0; i < numRequests; i++ {
			statusCode := <-results
			if statusCode == http.StatusOK {
				successCount++
			}
		}

		// All requests should succeed
		assert.Equal(t, numRequests, successCount)
	})
}

func TestHealthIntegration_VersionInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create application context
	log := logger.New()
	appCtx := core.NewAppContextWithDefaults(log)

	// Create health handler
	healthHandler := NewHealthHandler(appCtx)

	// Create router
	router := gin.New()
	router.GET("/health", healthHandler.HealthCheck)

	// Test version information in health response
	t.Run("Version Information", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify version info
		assert.Equal(t, "1.0.0", response.Version.Version)
		assert.Equal(t, "none", response.Version.Commit)
		assert.NotEmpty(t, response.Version.BuildTime)
		assert.NotEmpty(t, response.Version.GoVersion)
		assert.Contains(t, response.Version.GoVersion, "go")
	})
}

func TestHealthIntegration_UptimeTracking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create application context
	log := logger.New()
	appCtx := core.NewAppContextWithDefaults(log)

	// Create health handler
	healthHandler := NewHealthHandler(appCtx)

	// Create router
	router := gin.New()
	router.GET("/health", healthHandler.HealthCheck)

	// Test uptime tracking
	t.Run("Uptime Tracking", func(t *testing.T) {
		// First request
		req1, _ := http.NewRequest("GET", "/health", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		var response1 HealthResponse
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		require.NoError(t, err)

		firstUptime := response1.Uptime
		assert.NotEmpty(t, firstUptime)

		// Wait a bit
		time.Sleep(100 * time.Millisecond)

		// Second request
		req2, _ := http.NewRequest("GET", "/health", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response2 HealthResponse
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		require.NoError(t, err)

		secondUptime := response2.Uptime
		assert.NotEmpty(t, secondUptime)

		// Uptime should have increased
		assert.NotEqual(t, firstUptime, secondUptime)

		// Both should be valid duration strings
		assert.Contains(t, firstUptime, "s")
		assert.Contains(t, secondUptime, "s")
	})
}
