// Package helpers provides common test helpers and utilities
package helpers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// RequireNoError asserts that an error is nil
func (h *TestHelper) RequireNoError(err error) {
	require.NoError(h.t, err)
}

// AssertNoError asserts that an error is nil
func (h *TestHelper) AssertNoError(err error) {
	assert.NoError(h.t, err)
}

// RequireEqual asserts that two values are equal
func (h *TestHelper) RequireEqual(expected, actual interface{}) {
	require.Equal(h.t, expected, actual)
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual interface{}) {
	assert.Equal(h.t, expected, actual)
}

// RequireNotNil asserts that a value is not nil
func (h *TestHelper) RequireNotNil(value interface{}) {
	require.NotNil(h.t, value)
}

// AssertNotNil asserts that a value is not nil
func (h *TestHelper) AssertNotNil(value interface{}) {
	assert.NotNil(h.t, value)
}

// RandomString generates a random string of specified length
func (h *TestHelper) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomEmail generates a random email address
func (h *TestHelper) RandomEmail() string {
	return fmt.Sprintf("%s@test.com", h.RandomString(8))
}

// CreateTestDir creates a temporary test directory
func (h *TestHelper) CreateTestDir() string {
	tmpDir, err := os.MkdirTemp("", "authservice_test_*")
	require.NoError(h.t, err)

	h.t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return tmpDir
}

// CreateTestFile creates a temporary test file with content
func (h *TestHelper) CreateTestFile(content string) string {
	tmpFile, err := os.CreateTemp("", "authservice_test_*.tmp")
	require.NoError(h.t, err)

	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	require.NoError(h.t, err)

	h.t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

// WriteJSONToFile writes JSON data to a file
func (h *TestHelper) WriteJSONToFile(filename string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	require.NoError(h.t, err)

	err = os.WriteFile(filename, jsonData, 0644)
	require.NoError(h.t, err)
}

// ReadJSONFromFile reads JSON data from a file
func (h *TestHelper) ReadJSONFromFile(filename string, target interface{}) {
	data, err := os.ReadFile(filename)
	require.NoError(h.t, err)

	err = json.Unmarshal(data, target)
	require.NoError(h.t, err)
}

// FileExists checks if a file exists
func (h *TestHelper) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func (h *TestHelper) DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetProjectRoot returns the project root directory
func (h *TestHelper) GetProjectRoot() string {
	dir, err := os.Getwd()
	require.NoError(h.t, err)

	// Navigate up until we find go.mod
	for {
		if h.FileExists(filepath.Join(dir, "go.mod")) {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	require.Fail(h.t, "Could not find project root (go.mod)")
	return ""
}

// SetupTestEnvironment sets up a test environment with proper cleanup
func (h *TestHelper) SetupTestEnvironment() {
	// Ensure we're in the project root
	projectRoot := h.GetProjectRoot()
	err := os.Chdir(projectRoot)
	require.NoError(h.t, err)

	// Set test environment variables
	os.Setenv("GIN_MODE", "test")
	os.Setenv("LOG_LEVEL", "debug")

	h.t.Cleanup(func() {
		// Clean up environment variables
		os.Unsetenv("GIN_MODE")
		os.Unsetenv("LOG_LEVEL")
	})
}

// WaitForCondition waits for a condition to be true or timeout
func (h *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration, message string) {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			require.Fail(h.t, fmt.Sprintf("Condition not met within timeout: %s", message))
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}
