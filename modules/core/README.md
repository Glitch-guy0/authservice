# Application Context

The `modules/core` package provides the central application context that manages dependencies, configuration, and runtime state for the authService application.

## Overview

The `AppContext` serves as the single source of truth for application-wide dependencies and configuration. It follows the dependency injection pattern where components declare their dependencies through constructors rather than accessing global state.

## Core Components

### AppContext

The main struct that holds:
- **Configuration**: Application settings and environment variables
- **Dependencies**: Database, cache, logger, tracer, meter, broker connections
- **Runtime State**: Start time, shutdown signals, health status
- **Thread Safety**: Mutex protection for concurrent access

### Health Management

- **HealthStatus**: Represents component health with status, message, and timestamps
- **HealthChecker**: Interface for components to report their health
- **HealthManager**: Manages health checks for all registered components

### Shutdown Management

- **ShutdownHandler**: Represents a function that handles graceful shutdown
- **ShutdownManager**: Coordinates shutdown of all components in priority order
- **Graceful Shutdown**: Ensures clean termination with timeouts

## Usage

### Basic Usage

```go
package main

import (
    "github.com/Glitch-guy0/authService/modules/core"
    "github.com/Glitch-guy0/authService/modules/logger"
)

func main() {
    // Create logger
    log := logger.New()
    
    // Create application context
    config := map[string]interface{}{
        "app_name":    "authService",
        "version":     "1.0.0",
        "environment": "production",
    }
    
    appCtx := core.NewAppContext(log, config)
    
    // Use the context throughout your application
    logger := appCtx.GetLogger()
    logger.Info("Application started")
    
    // Graceful shutdown
    appCtx.Shutdown(30 * time.Second)
}
```

### Dependency Injection

```go
// Service with dependencies injected through constructor
type UserService struct {
    db     interface{}
    cache  interface{}
    logger logger.Logger
}

func NewUserService(appCtx *core.AppContext) *UserService {
    return &UserService{
        db:     appCtx.GetDatabase(),
        cache:  appCtx.GetCache(),
        logger: appCtx.GetLogger(),
    }
}
```

### Health Checks

```go
// Implement HealthChecker for your components
type DatabaseService struct {
    db *sql.DB
}

func (ds *DatabaseService) CheckHealth() core.HealthStatus {
    err := ds.db.Ping()
    if err != nil {
        return core.HealthStatus{
            Status:    core.StatusUnhealthy,
            Message:   fmt.Sprintf("Database ping failed: %v", err),
            Timestamp: time.Now(),
        }
    }
    
    return core.HealthStatus{
        Status:    core.StatusHealthy,
        Message:   "Database connection healthy",
        Timestamp: time.Now(),
    }
}

func (ds *DatabaseService) Name() string {
    return "database"
}

// Register the health checker
func main() {
    // ... setup appCtx ...
    
    dbService := NewDatabaseService(appCtx)
    appCtx.RegisterHealthChecker(dbService)
}
```

### Custom Shutdown Handlers

```go
func main() {
    // ... setup appCtx ...
    
    // Register custom shutdown handler
    appCtx.RegisterShutdownHandler("database", func(ctx context.Context) error {
        return db.Close()
    }, 10*time.Second, 10) // 10s timeout, priority 10
    
    // Or use simple registration
    appCtx.AddShutdownHandler("cache", func(ctx context.Context) error {
        return cache.Close()
    })
}
```

## Configuration

The AppContext accepts a `map[string]interface{}` for configuration. This allows flexible configuration management that can be loaded from:

- Environment variables
- Configuration files (YAML, JSON)
- Command-line arguments
- Remote configuration services

### Default Configuration

```go
// Use built-in defaults
appCtx := core.NewAppContextWithDefaults(logger)

// Default config includes:
// - app_name: "authService"
// - version: "1.0.0"
// - environment: "development"
// - debug: true
```

## Health Status

Components can have three health states:

- **healthy**: Component is functioning normally
- **degraded**: Component is functioning but with reduced capability
- **unhealthy**: Component is not functioning or has critical issues

### Health Status Structure

```go
type HealthStatus struct {
    Status    string    `json:"status"`    // "healthy", "unhealthy", "degraded"
    Message   string    `json:"message"`   // Optional message
    LastCheck time.Time `json:"lastCheck"`
    Timestamp time.Time `json:"timestamp"`
}
```

## Shutdown Process

The shutdown manager handles graceful shutdown in priority order:

1. **Lower priority numbers** shutdown first
2. **Each handler** gets its own timeout
3. **Overall timeout** prevents hanging
4. **Errors are logged** but don't stop other handlers

### Shutdown Priority Examples

```go
// High priority (shutdown first)
appCtx.RegisterShutdownHandler("http-server", shutdownServer, 5*time.Second, 1)
appCtx.RegisterShutdownHandler("database", closeDatabase, 10*time.Second, 2)

// Lower priority (shutdown later)
appCtx.RegisterShutdownHandler("cache", closeCache, 5*time.Second, 10)
appCtx.RegisterShutdownHandler("metrics", flushMetrics, 2*time.Second, 20)
```

## Thread Safety

The AppContext is designed to be thread-safe:

- **Mutex protection** for all shared state
- **Read-write locks** for configuration and health status
- **Atomic operations** for shutdown signaling
- **Copy-on-read** for configuration maps

## Best Practices

1. **Constructor Injection**: Always inject dependencies through constructors
2. **Interface Segregation**: Depend on interfaces, not concrete types
3. **Health Checks**: Implement health checks for all critical components
4. **Graceful Shutdown**: Register shutdown handlers for all resources
5. **Configuration**: Use environment-specific configuration
6. **Logging**: Use the provided logger for consistent formatting

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./modules/core/...

# Run with coverage
go test -cover ./modules/core/...

# Run benchmarks
go test -bench=. ./modules/core/...
```

## Integration with Other Modules

The AppContext integrates with:

- **Logger**: Provides structured logging with context
- **Configuration**: Loads and validates configuration
- **Database**: Manages database connections and pools
- **Cache**: Handles cache connections and strategies
- **HTTP Server**: Manages server lifecycle and graceful shutdown
- **Metrics**: Provides metrics collection and reporting

## Future Enhancements

Planned improvements include:

1. **Configuration Validation**: Automatic validation of configuration schemas
2. **Dependency Graph**: Visual representation of component dependencies
3. **Hot Reload**: Configuration changes without restart
4. **Circuit Breaker**: Automatic health-based circuit breaking
5. **Metrics Integration**: Built-in metrics for health and performance

## Troubleshooting

### Common Issues

1. **Deadlocks**: Ensure proper lock ordering and avoid nested locks
2. **Memory Leaks**: Close all resources in shutdown handlers
3. **Race Conditions**: Use the provided thread-safe methods
4. **Timeout Issues**: Set appropriate timeouts for shutdown handlers

### Debugging

Enable debug logging to troubleshoot issues:

```go
config := map[string]interface{}{
    "debug": true,
}
appCtx := core.NewAppContext(logger, config)
```
