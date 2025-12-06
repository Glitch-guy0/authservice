package config

// AppConfig represents the application configuration structure
type AppConfig struct {
	// Environment specifies the current environment (development, staging, production)
	Environment string `mapstructure:"env"`

	// Server configuration
	Server ServerConfig `mapstructure:"server"`

	// Log configuration
	Log LogConfig `mapstructure:"log"`

	// Database configuration
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	// Port to listen on
	Port int `mapstructure:"port"`

	// Timeout settings
	Timeout struct {
		// Read timeout in seconds
		Read int `mapstructure:"read"`
		// Write timeout in seconds
		Write int `mapstructure:"write"`
		// Idle timeout in seconds
		Idle int `mapstructure:"idle"`
	} `mapstructure:"timeout"`

	// Enable/disable debug mode
	Debug bool `mapstructure:"debug"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	// Log level (debug, info, warn, error, fatal, panic)
	Level string `mapstructure:"level"`

	// Log format (json, text)
	Format string `mapstructure:"format"`

	// Enable/disable file logging
	FileLogging struct {
		Enabled  bool   `mapstructure:"enabled"`
		Filename string `mapstructure:"filename"`
		MaxSize  int    `mapstructure:"max_size"` // in MB
		MaxAge   int    `mapstructure:"max_age"`  // in days
	} `mapstructure:"file_logging"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// Database driver (postgres, mysql, sqlite)
	Driver string `mapstructure:"driver"`

	// Database connection string
	DSN string `mapstructure:"dsn"`

	// Maximum number of open connections
	MaxOpenConns int `mapstructure:"max_open_conns"`

	// Maximum number of idle connections
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// Connection maximum lifetime in minutes
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`

	// Enable/disable SQL query logging
	LogQueries bool `mapstructure:"log_queries"`
}
