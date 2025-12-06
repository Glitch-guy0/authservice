package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected string
	}{
		{
			name:     "without underlying error",
			err:      &AppError{Code: ErrCodeValidation, Message: "invalid input", Op: "test"},
			expected: "test: invalid input",
		},
		{
			name:     "with underlying error",
			err:      &AppError{Code: ErrCodeInternal, Message: "internal error", Op: "test", Err: errors.New("underlying error")},
			expected: "test: internal error: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	original := errors.New("original error")
	wrapped := Wrap(original, ErrCodeValidation, "validation failed", "test")

	if wrapped == nil {
		t.Fatal("Wrap() returned nil")
	}

	appErr, ok := wrapped.(*AppError)
	if !ok {
		t.Fatal("Wrap() did not return an AppError")
	}

	if appErr.Code != ErrCodeValidation {
		t.Errorf("expected code %v, got %v", ErrCodeValidation, appErr.Code)
	}

	if appErr.Message != "validation failed" {
		t.Errorf("expected message 'validation failed', got '%s'", appErr.Message)
	}

	if appErr.Op != "test" {
		t.Errorf("expected op 'test', got '%s'", appErr.Op)
	}

	if appErr.Err != original {
		t.Error("wrapped error does not contain the original error")
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     ErrorCode
		expected bool
	}{
		{
			name:     "direct match",
			err:      New(ErrCodeValidation, "invalid input", "test"),
			code:     ErrCodeValidation,
			expected: true,
		},
		{
			name:     "no match",
			err:      New(ErrCodeValidation, "invalid input", "test"),
			code:     ErrCodeInternal,
			expected: false,
		},
		{
			name:     "wrapped error match",
			err:      Wrap(errors.New("underlying"), ErrCodeNotFound, "not found", "test"),
			code:     ErrCodeNotFound,
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			code:     ErrCodeInternal,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.code); got != tt.expected {
				t.Errorf("Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "validation error",
			err:            New(ErrCodeValidation, "invalid email", "test"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   string(ErrCodeValidation),
		},
		{
			name:           "unauthorized error",
			err:            New(ErrCodeUnauthorized, "invalid token", "test"),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   string(ErrCodeUnauthorized),
		},
		{
			name:           "forbidden error",
			err:            New(ErrCodeForbidden, "access denied", "test"),
			expectedStatus: http.StatusForbidden,
			expectedCode:   string(ErrCodeForbidden),
		},
		{
			name:           "not found error",
			err:            New(ErrCodeNotFound, "user not found", "test"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   string(ErrCodeNotFound),
		},
		{
			name:           "internal error",
			err:            New(ErrCodeInternal, "database error", "test"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   string(ErrCodeInternal),
		},
		{
			name:           "non-AppError",
			err:            errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   string(ErrCodeInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, httpErr := Format(tt.err)

			if status != tt.expectedStatus {
				t.Errorf("Format() status = %v, want %v", status, tt.expectedStatus)
			}

			if httpErr.Code != tt.expectedCode {
				t.Errorf("Format() code = %v, want %v", httpErr.Code, tt.expectedCode)
			}

			// For internal errors, the message should be generic
			if tt.expectedStatus == http.StatusInternalServerError && httpErr.Message == "" {
				t.Error("Format() message should not be empty for internal errors")
			}
		})
	}
}

func TestJSON(t *testing.T) {
	// This would typically be tested with an integration test using httptest
	// For now, we'll just test that it doesn't panic with a proper writer
	err := New(ErrCodeValidation, "invalid input", "test")

	// Use httptest.NewRecorder instead of nil
	w := httptest.NewRecorder()
	JSON(w, err) // We can't easily test the writer output in a unit test

	// Verify the response was written
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
