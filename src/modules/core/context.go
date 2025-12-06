package core

import (
	"context"
	"sync"
	"time"

	"github.com/Glitch-guy0/authService/src/modules/logger"
)

// AppContext holds application-wide dependencies and configuration
type AppContext struct {
	// Core dependencies
	Logger logger.Logger
	Config map[string]interface{}

	// Infrastructure placeholders (to be implemented later)
	Database interface{}
	Cache    interface{}
	Tracer   interface{}
	Meter    interface{}
	Broker   interface{}

	// Runtime state
	startTime time.Time
	shutdown  chan struct{}
	mu        sync.RWMutex

	// Health tracking
	healthStatus map[string]HealthStatus
	healthMu     sync.RWMutex
}

// HealthStatus represents the health status of a component
type HealthStatus struct {
	Status    string    `json:"status"`  // "healthy", "unhealthy", "degraded"
	Message   string    `json:"message"` // Optional message
	LastCheck time.Time `json:"lastCheck"`
	Timestamp time.Time `json:"timestamp"`
}

// ComponentStatus constants
const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
	StatusDegraded  = "degraded"
)

// NewAppContext creates a new application context
func NewAppContext(logger logger.Logger, config map[string]interface{}) *AppContext {
	ctx := &AppContext{
		Logger:       logger,
		Config:       config,
		startTime:    time.Now(),
		shutdown:     make(chan struct{}),
		healthStatus: make(map[string]HealthStatus),
	}

	// Set initial health status for the context itself
	ctx.UpdateHealthStatus("appcontext", HealthStatus{
		Status:    StatusHealthy,
		Message:   "Application context initialized",
		Timestamp: time.Now(),
	})

	return ctx
}

// NewAppContextWithDefaults creates a new application context with default configuration
func NewAppContextWithDefaults(logger logger.Logger) *AppContext {
	defaultConfig := map[string]interface{}{
		"app_name":    "authService",
		"version":     "1.0.0",       // will be replaced when viper is initialized
		"environment": "development", // will be replaced when viper is initialized
		"debug":       true,          // will be replaced when viper is initialized
	}

	return NewAppContext(logger, defaultConfig)
}

// NewAppContextFromEnv creates a new application context from environment configuration
func NewAppContextFromEnv(logger logger.Logger) *AppContext {
	// Load configuration from environment (placeholder for now)
	// This would typically use viper or similar configuration management

	return NewAppContextWithDefaults(logger)
}

// Clone creates a shallow copy of the application context
func (ac *AppContext) Clone() *AppContext {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	clone := &AppContext{
		Logger:       ac.Logger,
		Config:       ac.GetConfig(),
		Database:     ac.Database,
		Cache:        ac.Cache,
		Tracer:       ac.Tracer,
		Meter:        ac.Meter,
		Broker:       ac.Broker,
		startTime:    ac.startTime,
		shutdown:     make(chan struct{}),
		healthStatus: ac.GetHealthStatus(),
	}

	return clone
}

// GetStartTime returns the application start time
func (ac *AppContext) GetStartTime() time.Time {
	return ac.startTime
}

// GetUptime returns the application uptime
func (ac *AppContext) GetUptime() time.Duration {
	return time.Since(ac.startTime)
}

// GetLogger returns the logger instance
func (ac *AppContext) GetLogger() logger.Logger {
	return ac.Logger
}

// GetConfig returns the configuration map
func (ac *AppContext) GetConfig() map[string]interface{} {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	// Return a copy to prevent external modification
	config := make(map[string]interface{})
	for k, v := range ac.Config {
		config[k] = v
	}
	return config
}

// SetConfig updates a configuration value
func (ac *AppContext) SetConfig(key string, value interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Config[key] = value
}

// GetDatabase returns the database connection (placeholder)
func (ac *AppContext) GetDatabase() interface{} {
	return ac.Database
}

// SetDatabase sets the database connection
func (ac *AppContext) SetDatabase(db interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Database = db
}

// GetCache returns the cache connection (placeholder)
func (ac *AppContext) GetCache() interface{} {
	return ac.Cache
}

// SetCache sets the cache connection
func (ac *AppContext) SetCache(cache interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Cache = cache
}

// GetTracer returns the tracer instance (placeholder)
func (ac *AppContext) GetTracer() interface{} {
	return ac.Tracer
}

// SetTracer sets the tracer instance
func (ac *AppContext) SetTracer(tracer interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Tracer = tracer
}

// GetMeter returns the meter instance (placeholder)
func (ac *AppContext) GetMeter() interface{} {
	return ac.Meter
}

// SetMeter sets the meter instance
func (ac *AppContext) SetMeter(meter interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Meter = meter
}

// GetBroker returns the message broker instance (placeholder)
func (ac *AppContext) GetBroker() interface{} {
	return ac.Broker
}

// SetBroker sets the message broker instance
func (ac *AppContext) SetBroker(broker interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.Broker = broker
}

// IsShutdown checks if the application is shutting down
func (ac *AppContext) IsShutdown() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	return ac.shutdown == nil
}

// GetShutdownChannel returns the shutdown channel
func (ac *AppContext) GetShutdownChannel() <-chan struct{} {
	return ac.shutdown
}

// Context returns a Go context that cancels on shutdown
func (ac *AppContext) Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ac.shutdown:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}

// UpdateHealthStatus updates the health status of a component
func (ac *AppContext) UpdateHealthStatus(component string, status HealthStatus) {
	ac.healthMu.Lock()
	defer ac.healthMu.Unlock()

	status.Timestamp = time.Now()
	ac.healthStatus[component] = status
}

// GetHealthStatus returns the health status of all components
func (ac *AppContext) GetHealthStatus() map[string]HealthStatus {
	ac.healthMu.RLock()
	defer ac.healthMu.RUnlock()

	// Return a copy to prevent external modification
	status := make(map[string]HealthStatus)
	for k, v := range ac.healthStatus {
		status[k] = v
	}
	return status
}

// GetOverallHealth returns the overall health status
func (ac *AppContext) GetOverallHealth() HealthStatus {
	ac.healthMu.RLock()
	defer ac.healthMu.RUnlock()

	if len(ac.healthStatus) == 0 {
		return HealthStatus{
			Status:    StatusHealthy,
			Message:   "No components registered",
			Timestamp: time.Now(),
		}
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, status := range ac.healthStatus {
		switch status.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	overallStatus := StatusHealthy
	message := "All components healthy"

	if hasUnhealthy {
		overallStatus = StatusUnhealthy
		message = "Some components are unhealthy"
	} else if hasDegraded {
		overallStatus = StatusDegraded
		message = "Some components are degraded"
	}

	return HealthStatus{
		Status:    overallStatus,
		Message:   message,
		Timestamp: time.Now(),
	}
}
