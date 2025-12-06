# Health Check Module

This module provides comprehensive health monitoring capabilities for the auth service, including liveness and readiness probes, version information, and extensible health checks.

## Features

- **Health Status Monitoring**: Track service health with configurable checkers
- **Liveness/Readiness Probes**: Kubernetes-ready probe endpoints
- **Version Information**: Build-time version injection and reporting
- **Extensible Architecture**: Easy to add custom health checkers
- **Structured Responses**: JSON responses with detailed status information
- **Performance Metrics**: Check execution timing and uptime tracking

## Endpoints

### `/health`
Main health check endpoint that returns comprehensive service health status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-12-05T16:50:00.554221+05:30",
  "version": {
    "version": "dev",
    "commit": "none",
    "build_time": "2025-12-05T16:49:55+05:30",
    "go_version": "go1.25.1"
  },
  "checks": [
    {
      "name": "server",
      "status": "healthy",
      "message": "Server is running",
      "duration": "417ns",
      "timestamp": "2025-12-05T16:50:00.554215+05:30"
    }
  ],
  "uptime": "4.712780208s"
}
```

**Status Codes:**
- `200 OK`: Service is healthy or degraded
- `503 Service Unavailable`: Service is unhealthy

### `/health/live`
Liveness probe endpoint. Always returns `200 OK` if the service is running.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-12-05T16:50:02.134313+05:30"
}
```

### `/health/ready`
Readiness probe endpoint. Returns service readiness based on health check results.

**Response (Ready):**
```json
{
  "status": "ready",
  "timestamp": "2025-12-05T16:50:05.590129+05:30"
}
```

**Response (Not Ready):**
```json
{
  "status": "not ready",
  "timestamp": "2025-12-05T16:50:05.590129+05:30"
}
```

**Status Codes:**
- `200 OK`: Service is ready (all checks healthy)
- `503 Service Unavailable`: Service is not ready (some checks unhealthy)

## Health Status Values

- `healthy`: All systems operational
- `degraded`: Service functional but with performance issues
- `unhealthy`: Service has critical failures

## Architecture

### Core Components

#### HealthService
Main service that manages health checkers and provides health status aggregation.

```go
type HealthService struct {
    appCtx          *core.AppContext
    logger          logger.Logger
    config          HealthCheckConfig
    startTime       time.Time
    checkers        []HealthChecker
    versionProvider *version.VersionProvider
    mu              sync.RWMutex
}
```

#### HealthHandler
HTTP handler for health check endpoints.

```go
type HealthHandler struct {
    healthService *HealthService
    logger        logger.Logger
}
```

#### VersionProvider
Provides build-time version information.

```go
type VersionProvider struct {
    appCtx  *core.AppContext
    logger  logger.Logger
    version VersionInfo
}
```

### Health Checkers

#### Default Checkers
- **server**: Basic server health check
- **database**: Database connectivity (placeholder for future implementation)
- **logger**: Logger operational status

#### Custom Checkers
Implement the `HealthChecker` interface to add custom health checks:

```go
type HealthChecker interface {
    Name() string
    Check(ctx context.Context) Check
}

type CustomChecker struct{}

func (c *CustomChecker) Name() string {
    return "custom-service"
}

func (c *CustomChecker) Check(ctx context.Context) Check {
    // Implement your health check logic
    return Check{
        Name:      "custom-service",
        Status:    StatusHealthy,
        Message:   "Service is operational",
        Timestamp: time.Now(),
    }
}
```

#### Basic Health Checker
For simple health checks, use the `BasicHealthChecker`:

```go
checker := NewBasicHealthChecker("service-name", func(ctx context.Context) Check {
    // Your check logic
    return Check{
        Name:      "service-name",
        Status:    StatusHealthy,
        Message:   "Service is healthy",
        Timestamp: time.Now(),
    }
})
```

## Configuration

### HealthCheckConfig
```go
type HealthCheckConfig struct {
    Enabled          bool          `json:"enabled"`
    CheckInterval    time.Duration `json:"check_interval"`
    Timeout          time.Duration `json:"timeout"`
    FailureThreshold int          `json:"failure_threshold"`
}
```

### Default Configuration
```go
func DefaultHealthCheckConfig() HealthCheckConfig {
    return HealthCheckConfig{
        Enabled:          true,
        CheckInterval:    30 * time.Second,
        Timeout:          5 * time.Second,
        FailureThreshold: 3,
    }
}
```

## Version Information

### Build-Time Variables
The version provider uses build-time variables that can be set during compilation:

```bash
go build -ldflags \
  "-X github.com/Glitch-guy0/authService/src/modules/version.version=1.0.0 \
   -X github.com/Glitch-guy0/authService/src/modules/version.commit=abc123 \
   -X github.com/Glitch-guy0/authService/src/modules/version.buildTime=2025-12-05T16:49:55Z \
   -X github.com/Glitch-guy0/authService/src/modules/version.buildUser=builder \
   -X github.com/Glitch-guy0/authService/src/modules/version.buildHost=build-server"
```

### Version Types
```go
type Version struct {
    Version   string    `json:"version"`
    Commit    string    `json:"commit,omitempty"`
    BuildTime time.Time `json:"build_time"`
    GoVersion string    `json:"go_version"`
    BuildUser string    `json:"build_user,omitempty"`
    BuildHost string    `json:"build_host,omitempty"`
    Dirty     bool      `json:"dirty,omitempty"`
    Tags      []string  `json:"tags,omitempty"`
}
```

## Usage

### Basic Setup
```go
// Create application context
appCtx := core.NewAppContextWithDefaults(logger)

// Create health handler
healthHandler := health.NewHealthHandler(appCtx)

// Register endpoints
router := gin.New()
router.GET("/health", healthHandler.HealthCheck)
router.GET("/health/live", healthHandler.LivenessProbe)
router.GET("/health/ready", healthHandler.ReadinessProbe)
```

### Adding Custom Checkers
```go
// Register custom checker
healthHandler.GetHealthService().RegisterChecker(
    NewBasicHealthChecker("database", func(ctx context.Context) Check {
        // Check database connectivity
        if dbIsHealthy() {
            return Check{
                Name:      "database",
                Status:    StatusHealthy,
                Message:   "Database connection OK",
                Timestamp: time.Now(),
            }
        }
        return Check{
            Name:      "database",
            Status:    StatusUnhealthy,
            Message:   "Database connection failed",
            Timestamp: time.Now(),
        }
    }),
)
```

## Kubernetes Integration

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

### Startup Probe
```yaml
startupProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 30
```

## Testing

### Unit Tests
Run unit tests for the health module:
```bash
go test ./src/modules/api/health/... -v
```

### Integration Tests
Run integration tests:
```bash
go test ./src/modules/api/health/... -v -tags=integration
```

### Test Coverage
Generate test coverage report:
```bash
go test ./src/modules/api/health/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices

1. **Keep Health Checks Fast**: Health checks should complete quickly (under 5 seconds)
2. **Be Specific**: Use descriptive names and messages for health checks
3. **Handle Timeouts**: Implement proper timeout handling for external dependencies
4. **Monitor Performance**: Track check execution times for performance issues
5. **Graceful Degradation**: Use `degraded` status for non-critical issues
6. **Consistent Responses**: Maintain consistent response format across all endpoints

## Troubleshooting

### Common Issues

1. **Slow Health Checks**: Check if external dependencies are causing delays
2. **False Positives**: Ensure health check logic accurately reflects service state
3. **Memory Leaks**: Monitor for memory leaks in long-running health checks
4. **Timeout Issues**: Adjust timeout values for slow external services

### Debugging

Enable debug logging to troubleshoot health check issues:
```go
logger.SetLevel(logger.DEBUG)
```

Check individual health check responses:
```bash
curl -s http://localhost:8080/health | jq '.checks[]'
```

## Future Enhancements

- [ ] Database health check implementation
- [ ] Redis health check
- [ ] External API dependency checks
- [ ] Metrics collection (Prometheus)
- [ ] Health check history tracking
- [ ] Alerting integration
- [ ] Configuration-based health checks
