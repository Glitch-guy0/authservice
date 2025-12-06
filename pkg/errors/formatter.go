// Package errors provides custom error types and error handling utilities for the application.
package errors

import (
	"encoding/json"
	"net/http"
)

// HTTPError represents a standardized HTTP error response
type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Format formats an error into an HTTP error response
func Format(err error) (int, HTTPError) {
	if err == nil {
		return http.StatusOK, HTTPError{}
	}

	// Handle AppError
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Code {
		case ErrCodeValidation:
			return http.StatusBadRequest, HTTPError{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			}
		case ErrCodeUnauthorized:
			return http.StatusUnauthorized, HTTPError{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			}
		case ErrCodeForbidden:
			return http.StatusForbidden, HTTPError{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			}
		case ErrCodeNotFound:
			return http.StatusNotFound, HTTPError{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			}
		default:
			return http.StatusInternalServerError, HTTPError{
				Code:    string(ErrCodeInternal),
				Message: "Internal server error",
			}
		}
	}

	// Default to internal server error for unknown error types
	return http.StatusInternalServerError, HTTPError{
		Code:    string(ErrCodeInternal),
		Message: "An unexpected error occurred",
	}
}

// JSON writes an error response as JSON
func JSON(w http.ResponseWriter, err error) {
	statusCode, httpErr := Format(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Only write the error body if it's not a 200 OK
	if statusCode != http.StatusOK {
		json.NewEncoder(w).Encode(httpErr)
	}
}
