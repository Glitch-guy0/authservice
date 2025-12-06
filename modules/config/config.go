package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration
var Config *viper.Viper

// Init initializes the configuration
func Init(configPath string) error {
	// Initialize Viper
	v := viper.New()

	// Set configuration file name and type
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	v.AddConfigPath(configPath)           // Custom config path
	v.AddConfigPath("./configs")          // Default config directory
	v.AddConfigPath("$HOME/.authservice") // User's home directory
	v.AddConfigPath(".")                  // Current directory

	// Enable environment variable support
	automaticEnv(v)

	// Read in config file
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config into the global Config variable
	Config = v

	return nil
}

// automaticEnv sets up automatic environment variable binding
func automaticEnv(v *viper.Viper) {
	// Enable environment variable support
	v.AutomaticEnv()

	// Set environment variable prefix
	v.SetEnvPrefix("AUTH")

	// Bind environment variables
	v.BindEnv("env")
	v.BindEnv("log.level")
	v.BindEnv("log.format")
	v.BindEnv("server.port")
	v.BindEnv("server.timeout.read")
	v.BindEnv("server.timeout.write")
	v.BindEnv("server.timeout.idle")
}

// GetString is a wrapper around viper's GetString
func GetString(key string) string {
	return Config.GetString(key)
}

// GetInt is a wrapper around viper's GetInt
func GetInt(key string) int {
	return Config.GetInt(key)
}

// GetBool is a wrapper around viper's GetBool
func GetBool(key string) bool {
	return Config.GetBool(key)
}

// GetStringMapString is a wrapper around viper's GetStringMapString
func GetStringMapString(key string) map[string]string {
	return Config.GetStringMapString(key)
}

// UnmarshalKey is a wrapper around viper's UnmarshalKey
func UnmarshalKey(key string, rawVal interface{}) error {
	return Config.UnmarshalKey(key, rawVal)
}
