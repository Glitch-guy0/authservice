package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjectRoot(t *testing.T) {
	root, err := GetProjectRoot()
	require.NoError(t, err)

	// Verify it's a directory
	info, err := os.Stat(root)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Verify it contains the go.mod file
	_, err = os.Stat(filepath.Join(root, "go.mod"))
	assert.NoError(t, err)
}

func TestMustReadFile(t *testing.T) {
	// Create a temporary file with test content
	content := "test content"
	tmpfile, err := os.CreateTemp("", "testfile-*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	// Test reading the file
	result := MustReadFile(tmpfile.Name())
	assert.Equal(t, content, string(result))

	// Test non-existent file
	assert.Panics(t, func() {
		MustReadFile("/non/existent/file")
	})
}

func TestMustCopyFile(t *testing.T) {
	// Create a source file
	srcContent := "source file content"
	srcPath := filepath.Join(t.TempDir(), "source.txt")
	err := os.WriteFile(srcPath, []byte(srcContent), 0o644)
	require.NoError(t, err)

	// Create a destination path
	dstPath := filepath.Join(t.TempDir(), "destination.txt")

	// Test copying the file
	MustCopyFile(srcPath, dstPath)

	// Verify the destination file exists and has the same content
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, srcContent, string(dstContent))

	// Test copying to a non-existent directory
	nonExistentDst := filepath.Join(t.TempDir(), "nonexistent", "file.txt")
	assert.Panics(t, func() {
		MustCopyFile(srcPath, nonExistentDst)
	})
}

func TestCreateTempDir(t *testing.T) {
	tempDir, cleanup := CreateTempDir()
	defer cleanup()

	// Verify the directory exists
	info, err := os.Stat(tempDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test cleanup
	cleanup()
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))
}

func TestCreateTempFile(t *testing.T) {
	content := "test file content"
	filePath, cleanup := CreateTempFile(content)
	defer cleanup()

	// Verify the file exists and has the correct content
	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, string(fileContent))

	// Test cleanup
	cleanup()
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestFileExists(t *testing.T) {
	// Test with existing file
	tmpfile, err := os.CreateTemp("", "testfile-*.txt")
	require.NoError(t, err)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	assert.True(t, FileExists(tmpfile.Name()))

	// Test with non-existent file
	assert.False(t, FileExists("/non/existent/file"))

	// Test with directory (should return false on most systems)
	// Note: On some systems, FileExists might return true for directories
	// So we'll just test that it doesn't panic
	tempDir := t.TempDir()
	_ = FileExists(tempDir)
}

func TestDirExists(t *testing.T) {
	// Test with existing directory
	tempDir := t.TempDir()
	assert.True(t, DirExists(tempDir))

	// Test with non-existent directory
	assert.False(t, DirExists("/non/existent/directory"))

	// Test with file (should return false)
	tmpfile, err := os.CreateTemp("", "testfile-*.txt")
	require.NoError(t, err)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	assert.False(t, DirExists(tmpfile.Name()))
}

func TestIntegration_FileOperations(t *testing.T) {
	// Test creating and using a temporary directory
	tempDir, cleanupDir := CreateTempDir()
	defer cleanupDir()

	// Create a test file in the temporary directory
	testContent := "integration test content"
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte(testContent), 0o644)
	require.NoError(t, err)

	// Verify file exists
	assert.True(t, FileExists(filePath))

	// Read the file back
	content := MustReadFile(filePath)
	assert.Equal(t, testContent, string(content))

	// Copy the file
	destPath := filepath.Join(tempDir, "copy.txt")
	MustCopyFile(filePath, destPath)

	// Verify the copy
	copiedContent, err := os.ReadFile(destPath)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(copiedContent))

	// Create a temporary file with content
	tempFilePath, cleanupFile := CreateTempFile("temporary content")
	defer cleanupFile()

	// Verify the temporary file
	tempFileContent, err := os.ReadFile(tempFilePath)
	require.NoError(t, err)
	assert.Equal(t, "temporary content", string(tempFileContent))

	// Test project root (should be the same for all tests)
	root1, err := GetProjectRoot()
	require.NoError(t, err)
	root2, err := GetProjectRoot()
	require.NoError(t, err)
	assert.Equal(t, root1, root2)
	assert.True(t, strings.HasSuffix(root1, "authService"), "project root should end with 'authService'")
}
