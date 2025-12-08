# Logger Module

A simple logging module that provides structured logging with different log levels.

## Usage

```go
import "github.com/Glitch-guy0/authService/internal/logger"

// Create a new logger instance
log := logger.New()

// Use the different log methods
log.Info("Application started")
log.Warn("This is a warning")
log.Error("Something went wrong")
log.Critical("Critical error occurred")

// Get the underlying logger if needed
stdLogger := log.Create()
stdLogger.Printf("Standard logger output")
```

## Public Methods

- **create()**: Returns the underlying standard logger instance
- **info(msg, args...)**: Logs informational messages with [INFO] prefix
- **warn(msg, args...)**: Logs warning messages with [WARN] prefix
- **error(msg, args...)**: Logs error messages with [ERROR] prefix
- **critical(msg, args...)**: Logs critical messages with [CRITICAL] prefix

## Features

- Structured logging with level prefixes
- Printf-style formatting support
- File and line number information
- Thread-safe operations
- Simple interface for easy testing and mocking
