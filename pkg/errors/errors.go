// Package errors provides custom error types and error handling utilities for the application.
package errors

import "fmt"

// ErrorCode represents a specific error type in the application.
type ErrorCode string

// Common error codes
const (
	// ErrCodeInternal represents internal server errors
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	// ErrCodeValidation represents validation errors
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	// ErrCodeNotFound represents resource not found errors
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeUnauthorized represents authentication/authorization errors
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrCodeForbidden represents permission denied errors
	ErrCodeForbidden ErrorCode = "FORBIDDEN"
)

// AppError represents an application error with a code and message
type AppError struct {
	// Code is the error code that can be used by clients to handle specific error cases
	Code ErrorCode `json:"code"`
	// Message is a human-readable error message
	Message string `json:"message"`
	// Op is the operation that caused the error (e.g., "user.Create")
	Op string `json:"-"`
	// Err is the underlying error that triggered this error
	Err error `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with the given code, message and operation
func New(code ErrorCode, message, op string) error {
	return &AppError{
		Code:    code,
		Message: message,
		Op:      op,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message, op string) error {
	if err == nil {
		return nil
	}

	// If it's already an AppError, just update the operation
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:    appErr.Code,
			Message: message,
			Op:      op,
			Err:     appErr.Err,
		}
	}

	return &AppError{
		Code:    code,
		Message: message,
		Op:      op,
		Err:     err,
	}
}

// Errorf creates a new error with the given code, operation and formatted message
func Errorf(code ErrorCode, op, format string, args ...interface{}) error {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Op:      op,
	}
}

// Is checks if the error is of a specific error code
func Is(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}

	// Check wrapped errors
	for err != nil {
		if appErr, ok := err.(*AppError); ok && appErr.Code == code {
			return true
		}
		err = Unwrap(err)
	}

	return false
}

// Unwrap is a helper function that calls errors.Unwrap
func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(interface{ Unwrap() error }); ok {
		return appErr.Unwrap()
	}

	return nil
}
