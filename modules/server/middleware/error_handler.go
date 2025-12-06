// Package middleware provides HTTP middleware for the server
package middleware

import (
	"log"
	"net/http"

	"github.com/Glitch-guy0/authService/pkg/errors"
)

// ErrorHandler is a middleware that handles errors returned by HTTP handlers
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response recorder to capture the response
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// If no error occurred and status is OK, return early
		if rw.status < http.StatusBadRequest {
			return
		}

		// Handle the error based on status code
		switch rw.status {
		case http.StatusNotFound:
			err := errors.New(
				errors.ErrCodeNotFound,
				"The requested resource was not found",
				r.URL.Path,
			)
			errors.JSON(w, err)
		case http.StatusMethodNotAllowed:
			err := errors.New(
				errors.ErrCodeValidation,
				"Method not allowed",
				r.Method+" "+r.URL.Path,
			)
			errors.JSON(w, err)
		default:
			// If we have a custom error in the context, use it
			if err := r.Context().Value("error"); err != nil {
				if e, ok := err.(error); ok {
					errors.JSON(w, e)
					return
				}
			}

			// Otherwise, use a generic error
			err := errors.New(
				errors.ErrCodeInternal,
				"An unexpected error occurred",
				r.URL.Path,
			)
			errors.JSON(w, err)
		}
	})
}

// responseWriter is a wrapper around http.ResponseWriter that records the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if WriteHeader hasn't been called
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// Recoverer is a middleware that recovers from panics and returns a 500 error
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr) // Let the server handle this
				}

				// Log the panic
				log.Printf("panic: %v", rvr)

				// Return a 500 error
				err := errors.New(
					errors.ErrCodeInternal,
					"The server encountered a problem and could not complete your request",
					r.URL.Path,
				)
				errors.JSON(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
