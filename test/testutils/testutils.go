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
