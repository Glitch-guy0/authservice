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

func TestNewHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	require.NotNil(t, handler)
	assert.NotNil(t, handler.healthService)
	assert.NotNil(t, handler.logger)
}

func TestHealthHandler_HealthCheck_Healthy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	router := gin.New()
	router.GET("/health", handler.HealthCheck)

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
}

func TestHealthHandler_HealthCheck_WithUnhealthyChecker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	// Add an unhealthy checker
	unhealthyChecker := NewBasicHealthChecker("unhealthy-service", func(ctx context.Context) Check {
		return Check{
			Name:      "unhealthy-service",
			Status:    StatusUnhealthy,
			Message:   "Service is down",
			Timestamp: time.Now(),
		}
	})
	handler.healthService.RegisterChecker(unhealthyChecker)

	router := gin.New()
	router.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusUnhealthy, response.Status)
	assert.NotEmpty(t, response.Checks)

	// Find the unhealthy check
	found := false
	for _, check := range response.Checks {
		if check.Name == "unhealthy-service" {
			found = true
			assert.Equal(t, StatusUnhealthy, check.Status)
			assert.Equal(t, "Service is down", check.Message)
			break
		}
	}
	assert.True(t, found, "Unhealthy check should be present in response")
}

func TestHealthHandler_HealthCheck_Degraded(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	// Add a degraded checker
	degradedChecker := NewBasicHealthChecker("degraded-service", func(ctx context.Context) Check {
		return Check{
			Name:      "degraded-service",
			Status:    StatusDegraded,
			Message:   "Service is slow",
			Timestamp: time.Now(),
		}
	})
	handler.healthService.RegisterChecker(degradedChecker)

	router := gin.New()
	router.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Degraded should still return 200
	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusDegraded, response.Status)
	assert.NotEmpty(t, response.Checks)
}

func TestHealthHandler_LivenessProbe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	router := gin.New()
	router.GET("/health/live", handler.LivenessProbe)

	req, _ := http.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_ReadinessProbe_Ready(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	router := gin.New()
	router.GET("/health/ready", handler.ReadinessProbe)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ready", response["status"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_ReadinessProbe_NotReady(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	// Add an unhealthy checker to make service not ready
	unhealthyChecker := NewBasicHealthChecker("unhealthy-service", func(ctx context.Context) Check {
		return Check{
			Name:      "unhealthy-service",
			Status:    StatusUnhealthy,
			Message:   "Service is down",
			Timestamp: time.Now(),
		}
	})
	handler.healthService.RegisterChecker(unhealthyChecker)

	router := gin.New()
	router.GET("/health/ready", handler.ReadinessProbe)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "not ready", response["status"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_GetHealthService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	healthService := handler.GetHealthService()
	require.NotNil(t, healthService)
	assert.Equal(t, handler.healthService, healthService)
}

func TestHealthHandler_DefaultCheckers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logger.New()
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	handler := NewHealthHandler(appCtx)

	// Verify default checkers are registered
	healthService := handler.GetHealthService()
	ctx := context.Background()
	health := healthService.GetHealth(ctx)

	// Should have server, database, and logger checkers
	assert.Equal(t, 3, len(health.Checks))

	checkNames := make(map[string]bool)
	for _, check := range health.Checks {
		checkNames[check.Name] = true
	}

	assert.True(t, checkNames["server"])
	assert.True(t, checkNames["database"])
	assert.True(t, checkNames["logger"])
}
