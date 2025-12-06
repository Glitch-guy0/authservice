# Error Handling Package

This package provides a consistent way to handle and report errors throughout the application. It includes custom error types, error formatting, and HTTP error responses.

## Overview

The error handling system is built around the `AppError` type, which includes:
- A standardized error code
- A human-readable message
- The operation that caused the error
- The underlying error (if any)

## Error Codes

The following error codes are defined:

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INTERNAL_ERROR` | 500 | Unexpected internal server error |
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `NOT_FOUND` | 404 | Requested resource not found |
| `UNAUTHORIZED` | 401 | Authentication required or invalid credentials |
| `FORBIDDEN` | 403 | Insufficient permissions |

## Usage

### Creating Errors

```go
import "github.com/Glitch-guy0/authService/pkg/errors"

// Simple error
err := errors.New(errors.ErrCodeValidation, "Email is required", "user.Create")

// Error with formatting
err := errors.Errorf(errors.ErrCodeValidation, "user.Validate", "Invalid email format: %s", email)

// Wrap an existing error
if err := someOperation(); err != nil {
    return errors.Wrap(err, errors.ErrCodeInternal, "Failed to process user", "user.Process")
}
```

### Checking Error Types

```go
if errors.Is(err, errors.ErrCodeNotFound) {
    // Handle not found
}

// Get the underlying AppError
if appErr, ok := err.(*errors.AppError); ok {
    log.Printf("Error code: %s, Message: %s", appErr.Code, appErr.Message)
}
```

### HTTP Error Responses

In HTTP handlers, use the `JSON` function to send standardized error responses:

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.GetUser(r.URL.Query().Get("id"))
    if err != nil {
        if errors.Is(err, errors.ErrCodeNotFound) {
            errors.JSON(w, errors.New(errors.ErrCodeNotFound, "User not found", "user.Get"))
            return
        }
        errors.JSON(w, err)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}
```

## Middleware

The `middleware` package includes an `ErrorHandler` that can be used to handle errors in a consistent way across all HTTP handlers:

```go
import "github.com/Glitch-guy0/authService/modules/server/middleware"

router := mux.NewRouter()
// Apply error handling middleware
router.Use(middleware.ErrorHandler)
```

## Best Practices

1. **Always include an operation name** when creating errors to make debugging easier.
2. **Wrap underlying errors** to maintain the error chain.
3. **Use specific error codes** to allow callers to handle different error cases appropriately.
4. **Keep error messages user-friendly** but include enough detail for debugging.
5. **Log errors with context** including the operation and any relevant data.

## Testing

Run the tests with:

```bash
go test -v ./pkg/errors/...
```
