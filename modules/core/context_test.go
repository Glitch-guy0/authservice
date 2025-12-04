package core

import (
	"testing"
	"time"

	"github.com/Glitch-guy0/authService/modules/logger"
)

func TestNewAppContext(t *testing.T) {
	// Create a mock logger
	log := logger.New()
	config := map[string]interface{}{
		"app_name": "test-app",
		"version":  "1.0.0",
	}

	ctx := NewAppContext(log, config)

	if ctx == nil {
		t.Fatal("NewAppContext returned nil")
	}

	if ctx.GetLogger() != log {
		t.Error("Logger not set correctly")
	}

	retrievedConfig := ctx.GetConfig()
	if retrievedConfig["app_name"] != "test-app" {
		t.Error("Config not set correctly")
	}

	if ctx.GetStartTime().IsZero() {
		t.Error("Start time not set")
	}

	if ctx.IsShutdown() {
		t.Error("AppContext should not be shutdown initially")
	}
}

func TestNewAppContextWithDefaults(t *testing.T) {
	log := logger.New()
	ctx := NewAppContextWithDefaults(log)

	if ctx == nil {
		t.Fatal("NewAppContextWithDefaults returned nil")
	}

	config := ctx.GetConfig()
	if config["app_name"] != "authService" {
		t.Error("Default app_name not set correctly")
	}

	if config["version"] != "1.0.0" {
		t.Error("Default version not set correctly")
	}
}

func TestAppContextConfig(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	// Test SetConfig and GetConfig
	ctx.SetConfig("test_key", "test_value")
	config := ctx.GetConfig()

	if config["test_key"] != "test_value" {
		t.Error("SetConfig/GetConfig not working correctly")
	}

	// Test that GetConfig returns a copy
	config["external_key"] = "external_value"
	originalConfig := ctx.GetConfig()
	if _, exists := originalConfig["external_key"]; exists {
		t.Error("GetConfig should return a copy, not reference to original")
	}
}

func TestAppContextDependencies(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	// Test database
	db := &struct{}{}
	ctx.SetDatabase(db)
	if ctx.GetDatabase() != db {
		t.Error("Database not set/get correctly")
	}

	// Test cache
	cache := &struct{}{}
	ctx.SetCache(cache)
	if ctx.GetCache() != cache {
		t.Error("Cache not set/get correctly")
	}

	// Test tracer
	tracer := &struct{}{}
	ctx.SetTracer(tracer)
	if ctx.GetTracer() != tracer {
		t.Error("Tracer not set/get correctly")
	}

	// Test meter
	meter := &struct{}{}
	ctx.SetMeter(meter)
	if ctx.GetMeter() != meter {
		t.Error("Meter not set/get correctly")
	}

	// Test broker
	broker := &struct{}{}
	ctx.SetBroker(broker)
	if ctx.GetBroker() != broker {
		t.Error("Broker not set/get correctly")
	}
}

func TestAppContextUptime(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	startTime := ctx.GetStartTime()
	if startTime.IsZero() {
		t.Error("Start time should not be zero")
	}

	// Wait a bit and check uptime
	time.Sleep(10 * time.Millisecond)
	uptime := ctx.GetUptime()
	if uptime < 10*time.Millisecond {
		t.Error("Uptime should be at least 10ms")
	}
}

func TestAppContextShutdown(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	if ctx.IsShutdown() {
		t.Error("AppContext should not be shutdown initially")
	}

	// Test shutdown channel
	shutdownChan := ctx.GetShutdownChannel()
	if shutdownChan == nil {
		t.Error("Shutdown channel should not be nil")
	}

	// Trigger shutdown
	err := ctx.Shutdown(1 * time.Second)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	if !ctx.IsShutdown() {
		t.Error("AppContext should be shutdown after calling Shutdown")
	}

	// Test that shutdown channel is closed
	select {
	case <-shutdownChan:
		// Channel should be closed
	default:
		t.Error("Shutdown channel should be closed")
	}
}

func TestAppContextContext(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	appCtx := ctx.Context()
	if appCtx == nil {
		t.Error("Context should not be nil")
	}

	// Test that context is cancelled on shutdown
	go func() {
		time.Sleep(10 * time.Millisecond)
		ctx.Shutdown(1 * time.Second)
	}()

	select {
	case <-appCtx.Done():
		// Context should be cancelled
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled after shutdown")
	}
}

func TestAppContextHealth(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	// Test initial health status
	status := ctx.GetOverallHealth()
	if status.Status != StatusHealthy {
		t.Errorf("Initial status should be healthy, got %s", status.Status)
	}

	// Test updating health status
	testStatus := HealthStatus{
		Status:    StatusDegraded,
		Message:   "Test degraded status",
		Timestamp: time.Now(),
	}
	ctx.UpdateHealthStatus("test-component", testStatus)

	retrievedStatus := ctx.GetHealthStatus()
	if retrievedStatus["test-component"].Status != StatusDegraded {
		t.Error("Health status not updated correctly")
	}

	// Test overall health with degraded component
	overall := ctx.GetOverallHealth()
	if overall.Status != StatusDegraded {
		t.Errorf("Overall status should be degraded, got %s", overall.Status)
	}

	// Test IsHealthy
	if ctx.IsHealthy() {
		t.Error("IsHealthy should return false when component is degraded")
	}
}

func TestAppContextClone(t *testing.T) {
	log := logger.New()
	config := map[string]interface{}{
		"test": "value",
	}
	ctx := NewAppContext(log, config)

	// Set some dependencies
	ctx.SetDatabase(&struct{}{})
	ctx.UpdateHealthStatus("test", HealthStatus{
		Status:    StatusHealthy,
		Message:   "Test",
		Timestamp: time.Now(),
	})

	// Clone the context
	clone := ctx.Clone()
	if clone == nil {
		t.Fatal("Clone returned nil")
	}

	// Test that clone has same values
	if clone.GetConfig()["test"] != "value" {
		t.Error("Clone should have same config")
	}

	if clone.GetDatabase() == nil {
		t.Error("Clone should have same database")
	}

	cloneStatus := clone.GetHealthStatus()
	if cloneStatus["test"].Status != StatusHealthy {
		t.Error("Clone should have same health status")
	}

	// Test that clone is independent (modify original)
	ctx.SetConfig("modified", "true")
	cloneConfig := clone.GetConfig()
	if _, exists := cloneConfig["modified"]; exists {
		t.Error("Clone should be independent of original")
	}
}

func TestHealthStatusConstants(t *testing.T) {
	if StatusHealthy != "healthy" {
		t.Error("StatusHealthy constant incorrect")
	}

	if StatusUnhealthy != "unhealthy" {
		t.Error("StatusUnhealthy constant incorrect")
	}

	if StatusDegraded != "degraded" {
		t.Error("StatusDegraded constant incorrect")
	}
}

func TestAppContextHealthSummary(t *testing.T) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	// Add a health status
	ctx.UpdateHealthStatus("test-component", HealthStatus{
		Status:    StatusHealthy,
		Message:   "Test component is healthy",
		Timestamp: time.Now(),
	})

	summary := ctx.GetHealthSummary()

	// Check summary structure
	if summary["overall"] == nil {
		t.Error("Summary should contain overall status")
	}

	if summary["components"] == nil {
		t.Error("Summary should contain components status")
	}

	if summary["timestamp"] == nil {
		t.Error("Summary should contain timestamp")
	}

	if summary["uptime"] == nil {
		t.Error("Summary should contain uptime")
	}
}

// Benchmark tests
func BenchmarkNewAppContext(b *testing.B) {
	log := logger.New()
	config := map[string]interface{}{
		"app_name": "bench-app",
		"version":  "1.0.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewAppContext(log, config)
	}
}

func BenchmarkAppContextGetConfig(b *testing.B) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.GetConfig()
	}
}

func BenchmarkAppContextUpdateHealthStatus(b *testing.B) {
	log := logger.New()
	ctx := NewAppContext(log, map[string]interface{}{})

	status := HealthStatus{
		Status:    StatusHealthy,
		Message:   "Benchmark test",
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.UpdateHealthStatus("bench-component", status)
	}
}
