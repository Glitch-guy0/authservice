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
	bindEnvs()

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

// bindEnvs binds environment variables to configuration keys
func bindEnvs() {
	// Server
	_ = viper.BindEnv("env", "APP_ENV")
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.timeout.read", "SERVER_READ_TIMEOUT")
	_ = viper.BindEnv("server.timeout.write", "SERVER_WRITE_TIMEOUT")
	_ = viper.BindEnv("server.timeout.idle", "SERVER_IDLE_TIMEOUT")
	_ = viper.BindEnv("server.debug", "SERVER_DEBUG")

	// Logging
	_ = viper.BindEnv("log.level", "LOG_LEVEL")
	_ = viper.BindEnv("log.format", "LOG_FORMAT")
	_ = viper.BindEnv("log.file_logging.enabled", "LOG_FILE_ENABLED")
	_ = viper.BindEnv("log.file_logging.filename", "LOG_FILENAME")
	_ = viper.BindEnv("log.file_logging.max_size", "LOG_MAX_SIZE")
	_ = viper.BindEnv("log.file_logging.max_age", "LOG_MAX_AGE")

	// Database
	_ = viper.BindEnv("database.driver", "DB_DRIVER")
	_ = viper.BindEnv("database.dsn", "DB_DSN")
	_ = viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	_ = viper.BindEnv("database.log_queries", "DB_LOG_QUERIES")
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
