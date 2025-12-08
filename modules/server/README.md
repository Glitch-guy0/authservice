# HTTP Server Module

This module provides the HTTP server implementation using the Gin web framework with comprehensive middleware, graceful shutdown, and health check integration.

## Features

- **Gin Web Server**: High-performance HTTP server with Gin framework
- **Graceful Shutdown**: Proper shutdown handling with connection management
- **Middleware Support**: Request logging, recovery, CORS, and security middleware
- **Health Check Integration**: Built-in health check registration with AppContext
- **Configuration Management**: Flexible server configuration with defaults
- **Metrics Collection**: Request metrics and performance tracking
- **Dependency Injection**: Full integration with AppContext

## Architecture

### Core Components

1. **Server** (`server.go`): Main HTTP server implementation
2. **Middleware** (`middleware/`): Request processing middleware
3. **Shutdown** (`shutdown.go`): Graceful shutdown management
4. **Tests** (`server_test.go`): Comprehensive test suite

### Integration Points

- **AppContext**: Dependency injection and configuration
- **Logger**: Structured logging integration
- **Health Checks**: Server health monitoring
- **Configuration**: Server configuration management

## Usage

### Basic Server Setup

```go
package main

import (
    "github.com/Glitch-guy0/authService/modules/core"
    "github.com/Glitch-guy0/authService/modules/core/logger"
    "github.com/Glitch-guy0/authService/modules/server"
)

func main() {
    // Create dependencies
    log := logger.New()
    appCtx := core.NewAppContext(log, map[string]interface{}{})
    
    // Create server with default configuration
    srv := server.NewServerWithDefaults(appCtx)
    
    // Initialize server (middleware + routes)
    srv.Initialize()
    
    // Start server
    if err := srv.Start(); err != nil {
        log.Fatal("Failed to start server", "error", err)
    }
}
```

### Custom Configuration

```go
config := server.ServerConfig{
    Host:         "0.0.0.0",
    Port:         8080,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
    Mode:         "release", // "debug", "release", "test"
}

srv := server.NewServer(appCtx, config)
```

### Middleware Setup

```go
// Request logging middleware
loggerConfig := middleware.DefaultLoggerConfig()
loggerMiddleware := middleware.NewLoggerMiddleware(log, loggerConfig)
srv.GetEngine().Use(loggerMiddleware.Middleware())

// Recovery middleware
recoveryConfig := middleware.DefaultRecoveryConfig()
recoveryMiddleware := middleware.NewRecoveryMiddleware(log, recoveryConfig)
srv.GetEngine().Use(recoveryMiddleware.Middleware())

// CORS middleware
corsConfig := middleware.DefaultCORSConfig()
corsMiddleware := middleware.NewCORSMiddleware(log, corsConfig)
srv.GetEngine().Use(corsMiddleware.Middleware())
```

### Custom Routes

```go
// Setup routes after server initialization
engine := srv.GetEngine()

// API versioning
v1 := engine.Group("/api/v1")
{
    v1.GET("/users", getUsersHandler)
    v1.POST("/users", createUserHandler)
    v1.PUT("/users/:id", updateUserHandler)
}

// Health check endpoint
engine.GET("/health", func(c *gin.Context) {
    health := srv.GetAppContext().GetOverallHealth()
    c.JSON(200, health)
})
```

## Configuration

### Server Configuration Options

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| Host | string | "0.0.0.0" | Server bind address |
| Port | int | 8080 | Server port |
| ReadTimeout | time.Duration | 15s | Request read timeout |
| WriteTimeout | time.Duration | 15s | Response write timeout |
| IdleTimeout | time.Duration | 60s | Connection idle timeout |
| Mode | string | "debug" | Gin mode ("debug", "release", "test") |

### Middleware Configuration

#### Logger Middleware

```go
config := middleware.LoggerConfig{
    SkipPaths:         []string{"/health", "/metrics"},
    LogRequestBody:    true,
    LogResponseBody:   false,
    MaxBodySize:       1024 * 1024, // 1MB
    EnableSecurityLog: true,
}
```

#### Recovery Middleware

```go
config := middleware.RecoveryConfig{
    EnableStackTrace: true,
    SkipPaths:        []string{"/health"},
    StackSize:        50, // Max stack frames to log
}
```

#### CORS Middleware

```go
config := middleware.CORSConfig{
    AllowedOrigins: []string{
        "http://localhost:3000",
        "https://yourdomain.com",
    },
    AllowedMethods: []string{
        "GET", "POST", "PUT", "DELETE", "OPTIONS",
    },
    AllowedHeaders: []string{
        "Origin", "Content-Type", "Authorization",
    },
    AllowCredentials: true,
    MaxAge:          12 * time.Hour,
}
```

## Graceful Shutdown

The server provides comprehensive graceful shutdown support:

### Automatic Shutdown Registration

The server automatically registers its shutdown handler with AppContext:

```go
// This is done automatically during server creation
srv.GetAppContext().RegisterShutdownHandler("gin-server", 
    func(ctx context.Context) error {
        return srv.Shutdown(ctx)
    }, 
    30*time.Second, // timeout
    1,              // priority (high)
)
```

### Manual Shutdown

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Error("Server shutdown failed", "error", err)
}
```

### Enhanced Shutdown with Connection Management

```go
shutdownManager := server.NewEnhancedShutdown(srv)

// Wait for active connections before shutdown
err := shutdownManager.GracefulShutdownWithConnections(ctx)
```

## Health Checks

The server integrates with the AppContext health check system:

### Health Check Registration

```go
// Automatically registers health checker
srv.RegisterHealthChecker()

// Health status is available via AppContext
health := srv.GetAppContext().GetHealthStatus()
serverHealth := health["http-server"]
```

### Health Check Implementation

The server health checker monitors:
- Server running status
- HTTP server responsiveness
- Connection health

## Metrics and Monitoring

### Built-in Metrics

The server provides built-in metrics collection:

```go
// Get server metrics
metrics := server.GetMetrics()

// Available metrics
metrics.StartTime         // Server start time
metrics.TotalRequests     // Total request count
metrics.ActiveConnections // Current active connections
metrics.TotalConnections  // Total connections handled
metrics.Errors            // Error count
metrics.AvgResponseTime   // Average response time
metrics.GetUptime()       // Server uptime
```

### Metrics Middleware

```go
// Add metrics collection to all routes
metrics := server.NewServerMetrics()
srv.GetEngine().Use(metrics.MetricsMiddleware())
```

## Security Features

### Request ID Tracking

```go
// Automatic request ID generation
srv.GetEngine().Use(middleware.RequestIDMiddleware())

// Access request ID in handlers
requestID := c.MustGet("request_id")
```

### Security Logging

```go
// Security event logging
securityConfig := middleware.DefaultRecoveryConfig()
securityMiddleware := middleware.NewSecurityLoggerMiddleware(log, securityConfig)
srv.GetEngine().Use(securityMiddleware.Middleware())
```

### CORS Security

```go
// Production CORS with security features
corsConfig := middleware.ProductionCORSConfig([]string{
    "https://yourdomain.com",
})
corsMiddleware := middleware.NewCORSSecurityMiddleware(log, corsConfig, 60)
srv.GetEngine().Use(corsMiddleware.Middleware())
```

## Testing

### Running Tests

```bash
# Run all server tests
go test ./modules/server/...

# Run specific test
go test ./modules/server/ -run TestServerConfig

# Run with verbose output
go test ./modules/server/... -v

# Run benchmarks
go test ./modules/server/... -bench=.
```

### Test Coverage

The test suite covers:
- Server configuration and initialization
- Middleware functionality
- Graceful shutdown behavior
- Health check integration
- Error handling and recovery
- CORS and security features

## Best Practices

### 1. Always Initialize Server

```go
srv := server.NewServerWithDefaults(appCtx)
srv.Initialize() // Always call this
```

### 2. Use AppContext for Dependencies

```go
// Good: Use AppContext for all dependencies
log := srv.GetAppContext().GetLogger()
config := srv.GetAppContext().GetConfig()

// Bad: Create new instances
log := logger.New() // Don't do this
```

### 3. Configure Middleware Appropriately

```go
// Development: Permissive configuration
if gin.Mode() == gin.DebugMode {
    corsConfig := middleware.DevelopmentCORSConfig()
    loggerConfig := middleware.DefaultLoggerConfig()
    loggerConfig.Debug = true
}

// Production: Secure configuration
if gin.Mode() == gin.ReleaseMode {
    corsConfig := middleware.ProductionCORSConfig(allowedOrigins)
    recoveryConfig := middleware.DefaultRecoveryConfig()
    recoveryConfig.EnableStackTrace = false
}
```

### 4. Handle Shutdown Gracefully

```go
// Always handle shutdown signals
srv := server.NewServerWithDefaults(appCtx)
srv.Initialize()

// Start server in goroutine
go func() {
    if err := srv.Start(); err != nil && err != http.ErrServerClosed {
        log.Fatal("Server failed", "error", err)
    }
}()

// Wait for shutdown signal
<-srv.GetAppContext().GetShutdownChannel()

// Graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Error("Shutdown failed", "error", err)
}
```

### 5. Monitor Health and Metrics

```go
// Regular health checks
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        health := srv.GetAppContext().GetOverallHealth()
        if health.Status != core.StatusHealthy {
            log.Warn("Server health degraded", "status", health.Status)
        }
        
        metrics := srv.GetMetrics()
        log.Info("Server metrics",
            "requests", metrics.TotalRequests,
            "errors", metrics.Errors,
            "avgResponseTime", metrics.AvgResponseTime,
        )
    }
}()
```

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   - Change server port in configuration
   - Check for other processes using the port

2. **CORS Issues**
   - Verify allowed origins configuration
   - Check preflight request handling

3. **Graceful Shutdown Timeout**
   - Increase shutdown timeout
   - Check for long-running requests

4. **Memory Leaks**
   - Monitor connection count
   - Ensure proper connection cleanup

### Debug Mode

Enable debug mode for detailed logging:

```go
config := server.DefaultServerConfig()
config.Mode = "debug"
srv := server.NewServer(appCtx, config)
```

### Health Check Failures

Check server health status:

```go
health := srv.GetAppContext().GetHealthStatus()
for name, status := range health {
    if status.Status != core.StatusHealthy {
        log.Error("Component unhealthy",
            "component", name,
            "status", status.Status,
            "message", status.Message,
        )
    }
}
```

## Integration Examples

### With Database

```go
// Add database to AppContext
db, err := sql.Open("postgres", dsn)
if err != nil {
    return err
}

appCtx.SetDatabase(db)

// Use in handlers
func getUserHandler(c *gin.Context) {
    db := c.MustGet("appCtx").(*core.AppContext).GetDatabase()
    // Use database...
}
```

### With Authentication

```go
// Auth middleware
func authMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    // Validate token...
    c.Next()
}

srv.GetEngine().Use(authMiddleware)
```

### With Rate Limiting

```go
// Rate limiting middleware
func rateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(100, 200) // 100 req/s, burst 200
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}

srv.GetEngine().Use(rateLimitMiddleware())
```

## Performance Considerations

### Connection Pooling

- Configure appropriate timeouts
- Monitor connection metrics
- Use connection limits

### Memory Usage

- Limit request body size
- Use streaming for large responses
- Monitor memory metrics

### Caching

- Implement response caching
- Use middleware for cache headers
- Cache static assets

This server module provides a solid foundation for building scalable, maintainable web applications with Go and Gin.
