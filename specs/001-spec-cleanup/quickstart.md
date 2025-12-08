# Quickstart – Controller Pattern & Health Module

## 1. Prerequisites
- Go 1.25.1 toolchain installed
- `gin-gonic` and repo dependencies (`go mod download`)
- Access to `core.AppContext` initialization (see `cmd/app/main.go`)

## 2. Health Module Structure

The health module is now located at `modules/core/health/` with the following structure:
- `controller.go`: Handles HTTP routing and delegates to handlers
- `handler.go`: Implements the actual health check logic
- `health.go`: Contains health check interfaces and models

## 3. Module Controller Pattern

Each module follows this controller pattern:

```go
// Controller handles HTTP requests for a specific domain
package health

type Controller struct {
    appCtx *core.AppContext
}

// NewController creates a new controller instance
func NewController(appCtx *core.AppContext) *Controller {
    return &Controller{
        appCtx: appCtx,
    }
}

// RegisterRoutes sets up all routes for this controller
func (c *Controller) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("", c.HealthCheck)
    router.GET("/check", c.HealthCheck)
    router.GET("/live", c.LivenessProbe)
    router.GET("/ready", c.ReadinessProbe)
}
```

Key points:
- Controllers are stateless; all dependencies are injected via `appCtx`
- Each controller is responsible for its own route registration
- Handlers are kept separate from routing logic

## 4. Registering Controllers with the Server

In `cmd/app/main.go`, register the health controller:

```go
// Initialize core dependencies
appCtx := core.NewAppContext(...)

// Create router
r := gin.Default()

// Register health routes
healthController := health.NewController(appCtx)
healthRoutes := r.Group("/health")
healthController.RegisterRoutes(healthRoutes)
```

## 5. Health Check Endpoints

The health module provides these endpoints:
- `GET /health` or `/health/check` - Full health check
- `GET /health/live` - Liveness probe (is the app running?)
- `GET /health/ready` - Readiness probe (can handle traffic?)

## 6. Extending Health Checks

To add custom health checks:

1. Implement the `HealthChecker` interface:

```go
type HealthChecker interface {
    Check() HealthStatus
    Name() string
}
```

2. Register your checker during app initialization:

```go
healthService := health.NewHealthService()
healthService.RegisterChecker("database", &DatabaseHealthChecker{db: appCtx.DB})
```

## 7. Verification

1. Run tests:
   ```bash
   go test ./...
   ```

2. Test endpoints:
   ```bash
   # Basic health check
   curl http://localhost:8080/health
   
   # Readiness check
   curl http://localhost:8080/health/ready
   
   # Liveness check
   curl http://localhost:8080/health/live
   ```

3. Run linters:
   ```bash
   golangci-lint run
   ```

## 8. Next Steps
- Add more health checks as needed (database, cache, external services)
- Implement circuit breakers for critical dependencies
- Add metrics and logging to health endpoints
- Coordinate config module relocation under `modules/core` per spec follow-up.
