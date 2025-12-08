package core

import (
	"sync"
	"time"

	"github.com/Glitch-guy0/authService/modules/core/logger"
)

// HealthChecker interface for components that can report their health
type HealthChecker interface {
	CheckHealth() HealthStatus
	Name() string
}

// HealthManager manages health checks for application components
type HealthManager struct {
	checkers map[string]HealthChecker
	mu       sync.RWMutex
	logger   logger.Logger
}

// NewHealthManager creates a new health manager
func NewHealthManager(logger logger.Logger) *HealthManager {
	return &HealthManager{
		checkers: make(map[string]HealthChecker),
		logger:   logger,
	}
}

// RegisterChecker registers a health checker
func (hm *HealthManager) RegisterChecker(checker HealthChecker) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.checkers[checker.Name()] = checker
	hm.logger.Info("Registered health checker", "name", checker.Name())
}

// UnregisterChecker unregisters a health checker
func (hm *HealthManager) UnregisterChecker(name string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	delete(hm.checkers, name)
	hm.logger.Info("Unregistered health checker", "name", name)
}

// CheckAll performs health checks on all registered components
func (hm *HealthManager) CheckAll() map[string]HealthStatus {
	hm.mu.RLock()
	checkers := make(map[string]HealthChecker)
	for name, checker := range hm.checkers {
		checkers[name] = checker
	}
	hm.mu.RUnlock()

	results := make(map[string]HealthStatus)

	for name, checker := range checkers {
		status := checker.CheckHealth()
		status.LastCheck = time.Now()
		results[name] = status

		hm.logger.Debug("Health check completed", "component", name, "status", status.Status)
	}

	return results
}

// CheckHealth performs health check for a specific component
func (hm *HealthManager) CheckHealth(name string) (HealthStatus, bool) {
	hm.mu.RLock()
	checker, exists := hm.checkers[name]
	hm.mu.RUnlock()

	if !exists {
		return HealthStatus{
			Status:    StatusUnhealthy,
			Message:   "Component not found",
			Timestamp: time.Now(),
		}, false
	}

	status := checker.CheckHealth()
	status.LastCheck = time.Now()

	hm.logger.Debug("Health check completed", "component", name, "status", status.Status)

	return status, true
}

// GetCheckerNames returns the names of all registered checkers
func (hm *HealthManager) GetCheckerNames() []string {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	names := make([]string, 0, len(hm.checkers))
	for name := range hm.checkers {
		names = append(names, name)
	}

	return names
}

// RegisterHealthChecker registers a health checker with the AppContext
func (ac *AppContext) RegisterHealthChecker(checker HealthChecker) {
	// Create health manager if it doesn't exist
	if ac.healthStatus == nil {
		ac.healthStatus = make(map[string]HealthStatus)
	}

	// Perform initial health check
	status := checker.CheckHealth()
	status.Timestamp = time.Now()
	status.LastCheck = time.Now()

	// Store the health status
	ac.UpdateHealthStatus(checker.Name(), status)

	ac.GetLogger().Info("Health checker registered", "name", checker.Name(), "status", status.Status)
}

// UnregisterHealthChecker unregisters a health checker
func (ac *AppContext) UnregisterHealthChecker(name string) {
	ac.healthMu.Lock()
	defer ac.healthMu.Unlock()

	delete(ac.healthStatus, name)
	ac.GetLogger().Info("Health checker unregistered", "name", name)
}

// PerformHealthChecks performs health checks on all registered components
func (ac *AppContext) PerformHealthChecks() map[string]HealthStatus {
	ac.healthMu.RLock()
	checkers := make(map[string]HealthChecker)

	// Convert health status to checkers (simplified approach)
	for name := range ac.healthStatus {
		checkers[name] = &defaultHealthChecker{
			name:   name,
			status: ac.healthStatus[name],
		}
	}
	ac.healthMu.RUnlock()

	results := make(map[string]HealthStatus)

	for name, checker := range checkers {
		status := checker.CheckHealth()
		status.LastCheck = time.Now()
		results[name] = status

		// Update the stored health status
		ac.UpdateHealthStatus(name, status)
	}

	return results
}

// defaultHealthChecker is a simple implementation of HealthChecker
type defaultHealthChecker struct {
	name   string
	status HealthStatus
}

func (d *defaultHealthChecker) CheckHealth() HealthStatus {
	return d.status
}

func (d *defaultHealthChecker) Name() string {
	return d.name
}

// SimpleHealthChecker creates a simple health checker from a function
type SimpleHealthChecker struct {
	componentName string
	checkFunc     func() HealthStatus
}

func NewSimpleHealthChecker(name string, checkFunc func() HealthStatus) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		componentName: name,
		checkFunc:     checkFunc,
	}
}

func (s *SimpleHealthChecker) CheckHealth() HealthStatus {
	if s.checkFunc != nil {
		return s.checkFunc()
	}
	return HealthStatus{
		Status:    StatusHealthy,
		Message:   "No health check implemented",
		Timestamp: time.Now(),
	}
}

func (s *SimpleHealthChecker) Name() string {
	return s.componentName
}

// AddHealthCheck adds a simple health check function
func (ac *AppContext) AddHealthCheck(name string, checkFunc func() HealthStatus) {
	checker := NewSimpleHealthChecker(name, checkFunc)
	ac.RegisterHealthChecker(checker)
}

// IsHealthy returns true if all components are healthy
func (ac *AppContext) IsHealthy() bool {
	overall := ac.GetOverallHealth()
	return overall.Status == StatusHealthy
}

// GetHealthSummary returns a summary of all component health
func (ac *AppContext) GetHealthSummary() map[string]interface{} {
	status := ac.GetHealthStatus()
	overall := ac.GetOverallHealth()

	summary := map[string]interface{}{
		"overall":    overall,
		"components": status,
		"timestamp":  time.Now(),
		"uptime":     ac.GetUptime().String(),
	}

	return summary
}
