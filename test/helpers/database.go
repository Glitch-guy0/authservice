// Package helpers provides database testing utilities
package helpers

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// DatabaseHelper provides database testing utilities
type DatabaseHelper struct {
	t      *testing.T
	db     *sql.DB
	config DatabaseConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// NewDatabaseHelper creates a new database helper
func NewDatabaseHelper(t *testing.T, config DatabaseConfig) *DatabaseHelper {
	return &DatabaseHelper{
		t:      t,
		config: config,
	}
}

// Connect connects to the database
func (h *DatabaseHelper) Connect() {
	dsn := h.buildDSN()
	db, err := sql.Open(h.config.Driver, dsn)
	require.NoError(h.t, err)

	// Test the connection
	err = db.Ping()
	require.NoError(h.t, err)

	h.db = db

	h.t.Cleanup(func() {
		if h.db != nil {
			h.db.Close()
		}
	})
}

// buildDSN builds the database connection string
func (h *DatabaseHelper) buildDSN() string {
	switch h.config.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			h.config.Host, h.config.Port, h.config.Username, h.config.Password, h.config.Database)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			h.config.Username, h.config.Password, h.config.Host, h.config.Port, h.config.Database)
	case "sqlite3":
		return h.config.Database
	default:
		require.Fail(h.t, fmt.Sprintf("Unsupported database driver: %s", h.config.Driver))
		return ""
	}
}

// GetDB returns the database connection
func (h *DatabaseHelper) GetDB() *sql.DB {
	require.NotNil(h.t, h.db, "Database not connected. Call Connect() first.")
	return h.db
}

// CreateTestTable creates a test table
func (h *DatabaseHelper) CreateTestTable(tableName string, schema string) {
	_, err := h.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, schema))
	require.NoError(h.t, err)

	h.t.Cleanup(func() {
		h.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	})
}

// DropTable drops a table
func (h *DatabaseHelper) DropTable(tableName string) {
	_, err := h.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	require.NoError(h.t, err)
}

// TruncateTable truncates a table
func (h *DatabaseHelper) TruncateTable(tableName string) {
	_, err := h.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	require.NoError(h.t, err)
}

// InsertTestData inserts test data into a table
func (h *DatabaseHelper) InsertTestData(tableName string, columns []string, data [][]interface{}) {
	for _, row := range data {
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			joinStrings(columns, ", "),
			placeholders(len(row)))

		_, err := h.db.Exec(query, row...)
		require.NoError(h.t, err)
	}
}

// CountRows returns the number of rows in a table
func (h *DatabaseHelper) CountRows(tableName string) int {
	var count int
	err := h.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	require.NoError(h.t, err)
	return count
}

// WaitForConnection waits for database to be available
func (h *DatabaseHelper) WaitForConnection(timeout time.Duration) {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			require.Fail(h.t, "Database connection timeout")
		case <-ticker.C:
			if h.db.Ping() == nil {
				return
			}
		}
	}
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// placeholders generates SQL placeholders for the given count
func placeholders(count int) string {
	if count <= 0 {
		return ""
	}

	result := "?"
	for i := 1; i < count; i++ {
		result += ", ?"
	}
	return result
}

// MockDatabaseHelper provides a mock database for testing
type MockDatabaseHelper struct {
	t      *testing.T
	db     *sql.DB
	tmpDir string
}

// NewMockDatabaseHelper creates a new mock database helper using SQLite
func NewMockDatabaseHelper(t *testing.T) *MockDatabaseHelper {
	helper := &TestHelper{t: t}
	tmpDir := helper.CreateTestDir()

	dbPath := fmt.Sprintf("%s/test.db", tmpDir)
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)

	mock := &MockDatabaseHelper{
		t:      t,
		db:     db,
		tmpDir: tmpDir,
	}

	t.Cleanup(func() {
		mock.db.Close()
	})

	return mock
}

// GetDB returns the mock database connection
func (m *MockDatabaseHelper) GetDB() *sql.DB {
	return m.db
}

// CreateMockTable creates a mock table for testing
func (m *MockDatabaseHelper) CreateMockTable(tableName string, schema string) {
	_, err := m.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, schema))
	require.NoError(m.t, err)
}
