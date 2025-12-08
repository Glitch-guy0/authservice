package health

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/core/logger"
)

// MockLogger for testing
type MockLogger struct {
	// Simple mock logger implementation
}

func (m *MockLogger) Create() *logrus.Logger {
	return nil // Return nil for simplicity in tests
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) Fatal(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) Panic(msg string, args ...interface{}) {
	// No-op for tests
}

func (m *MockLogger) WithField(key string, value interface{}) logger.Logger {
	return m
}

func (m *MockLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return m
}

// MockHealthChecker for testing
type MockHealthChecker struct {
	name  string
	check Check
}

func (m *MockHealthChecker) Name() string {
	return m.name
}

func (m *MockHealthChecker) Check(ctx context.Context) Check {
	return m.check
}

func TestNewHealthService(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	config := DefaultHealthCheckConfig()
	service := NewHealthService(appCtx, config)

	require.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.NotNil(t, service.versionProvider)
	assert.Equal(t, 0, len(service.checkers))
}

func TestHealthService_RegisterChecker(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	service := NewHealthService(appCtx, DefaultHealthCheckConfig())

	checker := &MockHealthChecker{
		name: "test-checker",
		check: Check{
			Name:      "test-checker",
			Status:    StatusHealthy,
			Message:   "Test check passed",
			Timestamp: time.Now(),
		},
	}

	service.RegisterChecker(checker)

	assert.Equal(t, 1, len(service.checkers))
	assert.Equal(t, checker, service.checkers[0])
}

func TestHealthService_GetHealth_NoCheckers(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	service := NewHealthService(appCtx, DefaultHealthCheckConfig())
	ctx := context.Background()

	health := service.GetHealth(ctx)

	assert.Equal(t, StatusHealthy, health.Status)
	assert.NotZero(t, health.Timestamp)
	assert.Equal(t, "1.0.0", health.Version.Version)
	assert.Empty(t, health.Checks)
	assert.NotEmpty(t, health.Uptime)
}

func TestHealthService_GetHealth_WithCheckers(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	service := NewHealthService(appCtx, DefaultHealthCheckConfig())

	healthyChecker := &MockHealthChecker{
		name: "healthy-service",
		check: Check{
			Name:      "healthy-service",
			Status:    StatusHealthy,
			Message:   "Service is healthy",
			Timestamp: time.Now(),
		},
	}

	unhealthyChecker := &MockHealthChecker{
		name: "unhealthy-service",
		check: Check{
			Name:      "unhealthy-service",
			Status:    StatusUnhealthy,
			Message:   "Service is unhealthy",
			Timestamp: time.Now(),
		},
	}

	degradedChecker := &MockHealthChecker{
		name: "degraded-service",
		check: Check{
			Name:      "degraded-service",
			Status:    StatusDegraded,
			Message:   "Service is degraded",
			Timestamp: time.Now(),
		},
	}

	service.RegisterChecker(healthyChecker)
	service.RegisterChecker(unhealthyChecker)
	service.RegisterChecker(degradedChecker)

	ctx := context.Background()
	health := service.GetHealth(ctx)

	// Overall status should be unhealthy when any checker is unhealthy
	assert.Equal(t, StatusUnhealthy, health.Status)
	assert.Equal(t, 3, len(health.Checks))

	// Check that all checkers were executed
	checkNames := make(map[string]bool)
	for _, check := range health.Checks {
		checkNames[check.Name] = true
		assert.NotEmpty(t, check.Duration)
	}

	assert.True(t, checkNames["healthy-service"])
	assert.True(t, checkNames["unhealthy-service"])
	assert.True(t, checkNames["degraded-service"])
}

func TestHealthService_GetHealth_AllHealthy(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	service := NewHealthService(appCtx, DefaultHealthCheckConfig())

	checker1 := &MockHealthChecker{
		name: "service1",
		check: Check{
			Name:      "service1",
			Status:    StatusHealthy,
			Message:   "Service1 is healthy",
			Timestamp: time.Now(),
		},
	}

	checker2 := &MockHealthChecker{
		name: "service2",
		check: Check{
			Name:      "service2",
			Status:    StatusHealthy,
			Message:   "Service2 is healthy",
			Timestamp: time.Now(),
		},
	}

	service.RegisterChecker(checker1)
	service.RegisterChecker(checker2)

	ctx := context.Background()
	health := service.GetHealth(ctx)

	assert.Equal(t, StatusHealthy, health.Status)
	assert.Equal(t, 2, len(health.Checks))
}

func TestHealthService_GetHealth_DegradedOnly(t *testing.T) {
	mockLogger := &MockLogger{}
	appCtx := &core.AppContext{
		Logger: mockLogger,
		Config: make(map[string]interface{}),
	}

	service := NewHealthService(appCtx, DefaultHealthCheckConfig())

	checker := &MockHealthChecker{
		name: "degraded-service",
		check: Check{
			Name:      "degraded-service",
			Status:    StatusDegraded,
			Message:   "Service is degraded",
			Timestamp: time.Now(),
		},
	}

	service.RegisterChecker(checker)

	ctx := context.Background()
	health := service.GetHealth(ctx)

	assert.Equal(t, StatusDegraded, health.Status)
	assert.Equal(t, 1, len(health.Checks))
}

func TestDefaultHealthCheckConfig(t *testing.T) {
	config := DefaultHealthCheckConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, 30*time.Second, config.CheckInterval)
	assert.Equal(t, 5*time.Second, config.Timeout)
	assert.Equal(t, 3, config.FailureThreshold)
}

func TestNewBasicHealthChecker(t *testing.T) {
	name := "test-checker"
	checkFunc := func(ctx context.Context) Check {
		return Check{
			Name:      name,
			Status:    StatusHealthy,
			Message:   "Test check",
			Timestamp: time.Now(),
		}
	}

	checker := NewBasicHealthChecker(name, checkFunc)

	require.NotNil(t, checker)
	assert.Equal(t, name, checker.Name())

	ctx := context.Background()
	result := checker.Check(ctx)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Equal(t, "Test check", result.Message)
}
