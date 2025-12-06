package health

import (
	"context"
	"time"
)

// HealthStatus represents the health status of the service
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
	StatusDegraded  HealthStatus = "degraded"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    HealthStatus `json:"status"`
	Timestamp time.Time    `json:"timestamp"`
	Version   VersionInfo  `json:"version"`
	Checks    []Check      `json:"checks,omitempty"`
	Uptime    string       `json:"uptime,omitempty"`
}

// Check represents an individual health check
type Check struct {
	Name      string       `json:"name"`
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	Duration  string       `json:"duration,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// VersionInfo represents version information
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	GoVersion string `json:"go_version,omitempty"`
}

// HealthCheckConfig represents configuration for health checks
type HealthCheckConfig struct {
	Enabled          bool          `json:"enabled"`
	CheckInterval    time.Duration `json:"check_interval"`
	Timeout          time.Duration `json:"timeout"`
	FailureThreshold int           `json:"failure_threshold"`
}

// ServiceHealth represents the health of a specific service
type ServiceHealth struct {
	Name        string                 `json:"name"`
	Status      HealthStatus           `json:"status"`
	LastChecked time.Time              `json:"last_checked"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// HealthChecker interface for implementing health checks
type HealthChecker interface {
	Name() string
	Check(ctx context.Context) Check
}

// DefaultHealthCheckConfig returns default health check configuration
func DefaultHealthCheckConfig() HealthCheckConfig {
	return HealthCheckConfig{
		Enabled:          true,
		CheckInterval:    30 * time.Second,
		Timeout:          5 * time.Second,
		FailureThreshold: 3,
	}
}
