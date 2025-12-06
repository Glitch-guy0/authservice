// Package testutils provides mock utilities for testing
package testutils

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) WithField(key string, value interface{}) interface{} {
	args := m.Called(key, value)
	return args.Get(0)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) interface{} {
	args := m.Called(fields)
	return args.Get(0)
}

func (m *MockLogger) WithError(err error) interface{} {
	args := m.Called(err)
	return args.Get(0)
}

// MockConfig is a mock implementation of the Config interface
type MockConfig struct {
	mock.Mock
	data map[string]interface{}
}

func NewMockConfig() *MockConfig {
	return &MockConfig{
		data: make(map[string]interface{}),
	}
}

func (m *MockConfig) Set(key string, value interface{}) {
	m.data[key] = value
}

func (m *MockConfig) Get(key string) interface{} {
	if val, ok := m.data[key]; ok {
		return val
	}
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockConfig) GetString(key string) string {
	if val, ok := m.data[key].(string); ok {
		return val
	}
	args := m.Called(key)
	return args.String(0)
}

func (m *MockConfig) GetInt(key string) int {
	if val, ok := m.data[key].(int); ok {
		return val
	}
	args := m.Called(key)
	return args.Int(0)
}

func (m *MockConfig) GetBool(key string) bool {
	if val, ok := m.data[key].(bool); ok {
		return val
	}
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockConfig) IsSet(key string) bool {
	_, ok := m.data[key]
	if ok {
		return true
	}
	args := m.Called(key)
	return args.Bool(0)
}

// TestContext provides a test context with timeout and cancellation
type TestContext struct {
	context.Context
	Cancel context.CancelFunc
}

// NewTestContext creates a new test context with timeout
func NewTestContext(t *testing.T, timeout time.Duration) *TestContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	t.Cleanup(func() {
		cancel()
	})

	return &TestContext{
		Context: ctx,
		Cancel:  cancel,
	}
}

// CaptureOutput captures stdout/stderr for testing
type CaptureOutput struct {
	oldStdout *os.File
	oldStderr *os.File
	readPipe  *os.File
	writePipe *os.File
}

// NewCaptureOutput creates a new output capturer
func NewCaptureOutput() *CaptureOutput {
	readPipe, writePipe, _ := os.Pipe()
	return &CaptureOutput{
		oldStdout: os.Stdout,
		oldStderr: os.Stderr,
		readPipe:  readPipe,
		writePipe: writePipe,
	}
}

// Start begins capturing output
func (c *CaptureOutput) Start() {
	os.Stdout = c.writePipe
	os.Stderr = c.writePipe
}

// Stop stops capturing and returns the captured output
func (c *CaptureOutput) Stop() string {
	c.writePipe.Close()
	os.Stdout = c.oldStdout
	os.Stderr = c.oldStderr

	data, _ := io.ReadAll(c.readPipe)
	c.readPipe.Close()

	return string(data)
}
