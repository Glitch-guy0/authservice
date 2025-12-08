package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadEnv loads environment variables from .env file if it exists
func LoadEnv() error {
	// Look for .env file in the current directory and 5 upper directories (parent directories)
	for i := 0; i < 5; i++ {
		dotenv := ".env"
		if i > 0 {
			dotenv = strings.Repeat("../", i) + ".env"
		}

		if _, err := os.Stat(dotenv); err == nil {
			viper.SetConfigFile(dotenv)
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading .env file: %w", err)
			}
			break
		}
	}

	// Set default values
	setDefaults()

	// Bind environment variables
	if err := bindEnvVars(); err != nil {
		return fmt.Errorf("error binding environment variables: %w", err)
	}

	// Set the global Config variable
	Config = viper.GetViper()

	return nil
}

// setDefaults sets default configuration values (fallback values)
func setDefaults() {
	// Server defaults
	viper.SetDefault("env", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.timeout.read", 15)
	viper.SetDefault("server.timeout.write", 15)
	viper.SetDefault("server.timeout.idle", 60)
	viper.SetDefault("server.debug", false)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.file_logging.enabled", false)
	viper.SetDefault("log.file_logging.filename", "app.log")
	viper.SetDefault("log.file_logging.max_size", 100) // MB
	viper.SetDefault("log.file_logging.max_age", 30)   // days

	// Database defaults
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5) // minutes
	viper.SetDefault("database.log_queries", false)
}

// bindEnvVars binds environment variables to configuration keys
func bindEnvVars() error {
	// Collect any errors that occur during binding
	var bindErrors []error

	// Helper function to bind env var and collect errors
	bind := func(key, envVar string) {
		if err := viper.BindEnv(key, envVar); err != nil {
			bindErrors = append(bindErrors, fmt.Errorf("failed to bind %s to %s: %w", envVar, key, err))
		}
	}

	// Server
	bind("env", "APP_ENV")
	bind("server.port", "SERVER_PORT")
	bind("server.timeout.read", "SERVER_READ_TIMEOUT")
	bind("server.timeout.write", "SERVER_WRITE_TIMEOUT")
	bind("server.timeout.idle", "SERVER_IDLE_TIMEOUT")
	bind("server.debug", "SERVER_DEBUG")

	// Logging
	bind("log.level", "LOG_LEVEL")
	bind("log.format", "LOG_FORMAT")
	bind("log.file_logging.enabled", "LOG_FILE_ENABLED")
	bind("log.file_logging.filename", "LOG_FILENAME")
	bind("log.file_logging.max_size", "LOG_MAX_SIZE")
	bind("log.file_logging.max_age", "LOG_MAX_AGE")

	// Database
	bind("database.driver", "DB_DRIVER")
	bind("database.dsn", "DB_DSN")
	bind("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	bind("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	bind("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	bind("database.log_queries", "DB_LOG_QUERIES")

	// Return a combined error if any bindings failed
	if len(bindErrors) > 0 {
		errMsgs := make([]string, len(bindErrors))
		for i, err := range bindErrors {
			errMsgs[i] = err.Error()
		}
		return fmt.Errorf("failed to bind environment variables: %s", strings.Join(errMsgs, "; "))
	}

	return nil
}

// GetEnv returns the current environment
func GetEnv() string {
	if Config != nil {
		return strings.ToLower(Config.GetString("env"))
	}
	return strings.ToLower(viper.GetString("env"))
}

// IsProduction checks if the current environment is production
func IsProduction() bool {
	return GetEnv() == "production"
}

// IsDevelopment checks if the current environment is development
func IsDevelopment() bool {
	return GetEnv() == "development"
}

// IsTesting checks if the current environment is testing
func IsTesting() bool {
	return GetEnv() == "test"
}
