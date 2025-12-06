package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/Glitch-guy0/authService/modules/config"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (string, func())
		expectedErr bool
	}{
		{
			name: "valid config file",
			setup: func() (string, func()) {
				tempDir, err := os.MkdirTemp("", "config-test")
				require.NoError(t, err)

				configContent := `
env: test
server:
  port: 8080
  timeout:
    read: 5
    write: 10
    idle: 60
  debug: true
log:
  level: debug
  format: json
`
				err = os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(configContent), 0644)
				require.NoError(t, err)

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedErr: false,
		},
		{
			name: "invalid config file",
			setup: func() (string, func()) {
				tempDir, err := os.MkdirTemp("", "config-test")
				require.NoError(t, err)

				// Invalid YAML
				configContent := `invalid: yaml: file`
				err = os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(configContent), 0644)
				require.NoError(t, err)

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configDir, cleanup := tc.setup()
			defer cleanup()

			err := config.Init(configDir)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config.Config)
			}
		})
	}
}

func TestEnvironmentHelpers(t *testing.T) {
	// Setup
	v := viper.New()
	v.Set("env", "test")
	config.Config = v

	// Test environment helpers
	assert.True(t, config.IsTesting())
	assert.False(t, config.IsDevelopment())
	assert.False(t, config.IsProduction())

	// Change environment
	v.Set("env", "development")
	assert.True(t, config.IsDevelopment())

	v.Set("env", "production")
	assert.True(t, config.IsProduction())
}

func TestConfigGetters(t *testing.T) {
	// Setup
	v := viper.New()
	v.Set("test.string", "test")
	v.Set("test.int", 42)
	v.Set("test.bool", true)
	v.Set("test.map", map[string]string{"key": "value"})

	config.Config = v

	// Test getters
	assert.Equal(t, "test", config.GetString("test.string"))
	assert.Equal(t, 42, config.GetInt("test.int"))
	assert.True(t, config.GetBool("test.bool"))
	assert.Equal(t, map[string]string{"key": "value"}, config.GetStringMapString("test.map"))

	// Test UnmarshalKey
	var result map[string]string
	err := config.UnmarshalKey("test.map", &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (string, func())
		envVars     map[string]string
		expectedErr bool
	}{
		{
			name: "valid .env file",
			setup: func() (string, func()) {
				tempDir, err := os.MkdirTemp("", "env-test")
				require.NoError(t, err)

				envContent := `
APP_ENV=test
SERVER_PORT=3000
LOG_LEVEL=debug
DB_DSN=postgres://user:pass@localhost:5432/testdb
`
				err = os.WriteFile(filepath.Join(tempDir, ".env"), []byte(envContent), 0644)
				require.NoError(t, err)

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, cleanup := tc.setup()
			defer cleanup()

			// Change to the temp directory
			oldDir, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(oldDir)

			err = os.Chdir(tempDir)
			require.NoError(t, err)

			err = config.LoadEnv()
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	// Setup a valid config
	validConfig := &config.AppConfig{
		Environment: "development",
		Server: config.ServerConfig{
			Port: 8080,
			Timeout: struct {
				Read  int `mapstructure:"read"`
				Write int `mapstructure:"write"`
				Idle  int `mapstructure:"idle"`
			}{
				Read:  15,
				Write: 15,
				Idle:  60,
			},
			Debug: false,
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "json",
		},
		Database: config.DatabaseConfig{
			Driver:          "postgres",
			DSN:             "postgres://user:pass@localhost:5432/db",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5,
			LogQueries:      false,
		},
	}

	// Test valid config
	t.Run("valid config", func(t *testing.T) {
		viper.Reset()
		v := viper.New()
		v.Set("env", validConfig.Environment)
		v.Set("server", validConfig.Server)
		v.Set("log", validConfig.Log)
		v.Set("database", validConfig.Database)
		config.Config = v

		err := config.ValidateConfig()
		assert.NoError(t, err)
	})

	// Test invalid environment
	t.Run("invalid environment", func(t *testing.T) {
		viper.Reset()
		v := viper.New()
		v.Set("env", "invalid")
		config.Config = v

		err := config.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid configuration")
	})

	// Test invalid database driver
	t.Run("invalid database driver", func(t *testing.T) {
		viper.Reset()
		v := viper.New()
		v.Set("env", "development")
		v.Set("database.driver", "invalid")
		config.Config = v

		err := config.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid configuration")
	})

	// Test JWT secret too short
	t.Run("database validation", func(t *testing.T) {
		viper.Reset()
		v := viper.New()
		v.Set("env", "development")
		v.Set("database.driver", "invalid")
		config.Config = v

		err := config.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid configuration")
	})
}

func TestEnvironmentVariables(t *testing.T) {
	// Set up test environment variables
	t.Setenv("APP_ENV", "test")
	t.Setenv("SERVER_PORT", "3000")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DB_DSN", "postgres://test:test@localhost:5432/testdb")

	// Initialize config
	err := config.LoadEnv()
	require.NoError(t, err)

	// Verify environment variables were loaded correctly
	assert.Equal(t, "test", config.GetString("env"))
	assert.Equal(t, 3000, config.GetInt("server.port"))
	assert.Equal(t, "debug", config.GetString("log.level"))
	assert.Equal(t, "postgres://test:test@localhost:5432/testdb", config.GetString("database.dsn"))
}

func TestConfigReload(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "config-reload-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "config.yaml")

	// Initial config
	initialConfig := `
env: development
server:
  port: 8080
  timeout:
    read: 15
    write: 15
    idle: 60
  debug: true
`
	err = os.WriteFile(configFile, []byte(initialConfig), 0644)
	require.NoError(t, err)

	// Initialize config
	err = config.Init(tempDir)
	require.NoError(t, err)

	// Verify initial config
	assert.Equal(t, 8080, config.GetInt("server.port"))
	assert.True(t, config.GetBool("server.debug"))

	// Update config file
	updatedConfig := `
env: development
server:
  port: 9000
  timeout:
    read: 15
    write: 15
    idle: 60
  debug: false
`
	err = os.WriteFile(configFile, []byte(updatedConfig), 0644)
	require.NoError(t, err)

	// Force config reload
	err = config.Config.ReadInConfig()
	require.NoError(t, err)

	// Verify updated config
	assert.Equal(t, 9000, config.GetInt("server.port"))
	assert.False(t, config.GetBool("server.debug"))
}

func TestConfigDefaults(t *testing.T) {
	// Reset viper and load defaults
	viper.Reset()
	config.LoadEnv()

	// Test default values
	assert.Equal(t, "development", config.GetString("env"))
	assert.Equal(t, 8080, config.GetInt("server.port"))
	assert.Equal(t, 15, config.GetInt("server.timeout.read"))
	assert.Equal(t, 15, config.GetInt("server.timeout.write"))
	assert.Equal(t, 60, config.GetInt("server.timeout.idle"))
	assert.Equal(t, "info", config.GetString("log.level"))
	assert.Equal(t, "json", config.GetString("log.format"))
	assert.Equal(t, "postgres", config.GetString("database.driver"))
	assert.Equal(t, 25, config.GetInt("database.max_open_conns"))
	assert.Equal(t, 5, config.GetInt("database.max_idle_conns"))
	assert.Equal(t, 5, config.GetInt("database.conn_max_lifetime"))
}
