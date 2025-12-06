// Package testutils provides common utilities for testing
package testutils

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// GetProjectRoot returns the absolute path to the project root directory
func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	// Navigate up to the project root (assuming this file is in test/testutils/)
	return filepath.Abs(filepath.Join(filepath.Dir(filename), "../.."))
}

// MustReadFile reads a file and returns its content, or panics on error
func MustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

// MustCopyFile copies a file from src to dst, or panics on error
func MustCopyFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		panic(err)
	}
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir() (string, func()) {
	tmpDir, err := os.MkdirTemp("", "authservice_test_*")
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// CreateTempFile creates a temporary file with the given content
func CreateTempFile(content string) (string, func()) {
	tmpFile, err := os.CreateTemp("", "authservice_test_*.tmp")
	if err != nil {
		panic(err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		panic(err)
	}
	tmpFile.Close()

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
