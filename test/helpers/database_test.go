package helpers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseHelper_Connect(t *testing.T) {
	t.Run("successful connection to mock database", func(t *testing.T) {
		h := NewMockDatabaseHelper(t)
		defer func() {
			err := os.RemoveAll(h.tmpDir)
			require.NoError(t, err)
		}()

		db := h.GetDB()
		require.NotNil(t, db)

		// Test connection
		err := db.Ping()
		require.NoError(t, err)
	})

	t.Run("successful connection with config", func(t *testing.T) {
		h := NewDatabaseHelper(t, DatabaseConfig{
			Driver:   "sqlite3",
			Database: ":memory:",
		})

		h.Connect()
		defer h.db.Close()

		err := h.db.Ping()
		require.NoError(t, err)
	})
}

func TestDatabaseHelper_TableOperations(t *testing.T) {
	t.Run("create mock table", func(t *testing.T) {
		h := NewMockDatabaseHelper(t)
		defer func() {
			err := os.RemoveAll(h.tmpDir)
			require.NoError(t, err)
		}()

		tableName := "test_table"
		schema := `id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		value INTEGER DEFAULT 0`

		h.CreateMockTable(tableName, schema)

		// Verify table exists by inserting data
		db := h.GetDB()
		_, err := db.Exec("INSERT INTO "+tableName+" (name, value) VALUES (?, ?)", "test", 123)
		require.NoError(t, err)

		// Verify data was inserted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestDatabaseHelper_TestData(t *testing.T) {
	h := NewMockDatabaseHelper(t)
	defer func() {
		err := os.RemoveAll(h.tmpDir)
		require.NoError(t, err)
	}()

	tableName := "test_data"
	schema := `id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	value INTEGER`

	h.CreateMockTable(tableName, schema)

	// Insert test data directly
	db := h.GetDB()
	_, err := db.Exec("INSERT INTO "+tableName+" (name, value) VALUES (?, ?), (?, ?)",
		"test1", 1, "test2", 2)
	require.NoError(t, err)

	// Verify data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Verify data values
	var name string
	var value int
	err = db.QueryRow("SELECT name, value FROM "+tableName+" WHERE name = ?", "test1").Scan(&name, &value)
	require.NoError(t, err)
	assert.Equal(t, "test1", name)
	assert.Equal(t, 1, value)
}

func TestDatabaseHelper_RealDatabase(t *testing.T) {
	t.Run("test with real database connection", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "db-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")
		h := NewDatabaseHelper(t, DatabaseConfig{
			Driver:   "sqlite3",
			Database: dbPath,
		})

		h.Connect()
		defer h.db.Close()

		// Test connection
		err = h.db.Ping()
		require.NoError(t, err)

		// Test WaitForConnection
		h.WaitForConnection(5 * time.Second)

		// Create a test table
		tableName := "test_table"
		h.CreateTestTable(tableName, `id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)

		// Insert test data using the helper
		columns := []string{"name"}
		data := [][]interface{}{
			{"test1"},
			{"test2"},
		}
		h.InsertTestData(tableName, columns, data)

		// Verify data was inserted
		count := h.CountRows(tableName)
		assert.Equal(t, 2, count)

		// Test clearing the table (using DELETE since TRUNCATE isn't supported in SQLite)
		_, err = h.db.Exec("DELETE FROM " + tableName)
		require.NoError(t, err)
		count = h.CountRows(tableName)
		assert.Equal(t, 0, count)

		// Clean up
		h.DropTable(tableName)

		// Verify table was dropped
		_, err = h.db.Exec("SELECT * FROM " + tableName)
		require.Error(t, err)
	})
}

func TestMockDatabaseHelper(t *testing.T) {
	h := NewMockDatabaseHelper(t)
	defer func() {
		err := os.RemoveAll(h.tmpDir)
		require.NoError(t, err)
	}()

	tableName := "mock_test"
	schema := `id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL`

	h.CreateMockTable(tableName, schema)

	// Test inserting data
	db := h.GetDB()
	_, err := db.Exec("INSERT INTO "+tableName+" (name) VALUES (?), (?)", "test1", "test2")
	require.NoError(t, err)

	// Verify data
	rows, err := db.Query("SELECT name FROM " + tableName)
	require.NoError(t, err)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		require.NoError(t, err)
		names = append(names, name)
	}

	assert.ElementsMatch(t, []string{"test1", "test2"}, names)
}

func TestDatabaseHelper_Transaction(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "db-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	h := NewDatabaseHelper(t, DatabaseConfig{
		Driver:   "sqlite3",
		Database: dbPath,
	})

	h.Connect()
	defer h.db.Close()

	// Create a test table
	tableName := "test_table"
	h.CreateTestTable(tableName, `id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)

	// Insert test data
	columns := []string{"name"}
	data := [][]interface{}{
		{"test1"},
		{"test2"},
	}
	h.InsertTestData(tableName, columns, data)

	// Test transaction
	tx, err := h.db.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO "+tableName+" (name) VALUES (?)", "tx_test")
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Verify the transaction was committed
	var count int
	err = h.db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}
