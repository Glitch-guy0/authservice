package config

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var validate *validator.Validate

// ValidateConfig validates the application configuration
func ValidateConfig() error {
	if validate == nil {
		validate = validator.New()
		registerCustomValidations(validate)
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Perform basic validation
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Perform advanced validation
	if err := validateDatabaseConfig(&cfg); err != nil {
		return err
	}

	return nil
}

// registerCustomValidations registers custom validation functions
func registerCustomValidations(v *validator.Validate) {
	_ = v.RegisterValidation("env", validateEnv)
	_ = v.RegisterValidation("log_level", validateLogLevel)
	_ = v.RegisterValidation("log_format", validateLogFormat)
	_ = v.RegisterValidation("db_driver", validateDBDriver)
}

// validateEnv validates the environment value
func validateEnv(fl validator.FieldLevel) bool {
	env := fl.Field().String()
	switch env {
	case "development", "staging", "production":
		return true
	default:
		return false
	}
}

// validateLogLevel validates the log level
func validateLogLevel(fl validator.FieldLevel) bool {
	level := fl.Field().String()
	switch level {
	case "debug", "info", "warn", "error", "fatal", "panic":
		return true
	default:
		return false
	}
}

// validateLogFormat validates the log format
func validateLogFormat(fl validator.FieldLevel) bool {
	format := fl.Field().String()
	switch format {
	case "json", "text":
		return true
	default:
		return false
	}
}

// validateDBDriver validates the database driver
func validateDBDriver(fl validator.FieldLevel) bool {
	driver := fl.Field().String()
	switch driver {
	case "postgres", "mysql", "sqlite":
		return true
	default:
		return false
	}
}

// validateDatabaseConfig performs advanced validation for database configuration
func validateDatabaseConfig(cfg *AppConfig) error {
	switch cfg.Database.Driver {
	case "postgres":
		if !strings.Contains(cfg.Database.DSN, "postgres") {
			return fmt.Errorf("invalid DSN for PostgreSQL")
		}
	case "mysql":
		if !strings.HasPrefix(cfg.Database.DSN, "mysql://") {
			return fmt.Errorf("invalid DSN for MySQL")
		}
	case "sqlite":
		if _, err := os.Stat(cfg.Database.DSN); os.IsNotExist(err) {
			// Check if we can create the SQLite database file
			f, err := os.Create(cfg.Database.DSN)
			if err != nil {
				return fmt.Errorf("failed to create SQLite database file: %w", err)
			}
			f.Close()
			os.Remove(cfg.Database.DSN) // Remove the test file
		}
	}

	// Validate connection pool settings
	if cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return fmt.Errorf("max_idle_conns cannot be greater than max_open_conns")
	}

	if cfg.Database.ConnMaxLifetime < 1 || cfg.Database.ConnMaxLifetime > 30 {
		return fmt.Errorf("conn_max_lifetime must be between 1 and 30 minutes")
	}

	return nil
}

// validateURL validates a URL string
func validateURL(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	if urlStr == "" {
		return false
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}

// validateDuration validates a duration string
func validateDuration(fl validator.FieldLevel) bool {
	durationStr := fl.Field().String()
	if durationStr == "" {
		return false
	}

	_, err := time.ParseDuration(durationStr)
	return err == nil
}

// validateRegexp validates a regular expression string
func validateRegexp(fl validator.FieldLevel) bool {
	regexStr := fl.Field().String()
	if regexStr == "" {
		return false
	}

	_, err := regexp.Compile(regexStr)
	return err == nil
}
