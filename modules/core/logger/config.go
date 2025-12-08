package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ConfigManager manages logger configuration
type ConfigManager struct {
	config *LogConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: DefaultConfig(),
	}
}

// NewConfigManagerWithConfig creates a new configuration manager with initial config
func NewConfigManagerWithConfig(config *LogConfig) *ConfigManager {
	return &ConfigManager{
		config: config,
	}
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *LogConfig {
	return cm.config
}

// SetConfig sets the configuration
func (cm *ConfigManager) SetConfig(config *LogConfig) {
	cm.config = config
}

// LoadFromEnv loads configuration from environment variables
func (cm *ConfigManager) LoadFromEnv() {
	// Load log format
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cm.config.Format = ParseLogFormat(format)
	}

	// Load log level
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cm.config.Level = ParseLogLevel(level)
	}

	// Load time format
	if timeFormat := os.Getenv("LOG_TIME_FORMAT"); timeFormat != "" {
		cm.config.TimeFormat = timeFormat
	}

	// Load color setting
	if color := os.Getenv("LOG_ENABLE_COLOR"); color != "" {
		if enableColor, err := strconv.ParseBool(color); err == nil {
			cm.config.EnableColor = enableColor
		}
	}

	// Load output
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		cm.config.Output = output
	}
}

// LoadFromEnvWithPrefix loads configuration from environment variables with a prefix
func (cm *ConfigManager) LoadFromEnvWithPrefix(prefix string) {
	// Load log format
	if format := os.Getenv(prefix + "_FORMAT"); format != "" {
		cm.config.Format = ParseLogFormat(format)
	}

	// Load log level
	if level := os.Getenv(prefix + "_LEVEL"); level != "" {
		cm.config.Level = ParseLogLevel(level)
	}

	// Load time format
	if timeFormat := os.Getenv(prefix + "_TIME_FORMAT"); timeFormat != "" {
		cm.config.TimeFormat = timeFormat
	}

	// Load color setting
	if color := os.Getenv(prefix + "_ENABLE_COLOR"); color != "" {
		if enableColor, err := strconv.ParseBool(color); err == nil {
			cm.config.EnableColor = enableColor
		}
	}

	// Load output
	if output := os.Getenv(prefix + "_OUTPUT"); output != "" {
		cm.config.Output = output
	}
}

// SetFormat sets the log format
func (cm *ConfigManager) SetFormat(format LogFormat) {
	cm.config.Format = format
}

// SetLevel sets the log level
func (cm *ConfigManager) SetLevel(level LogLevel) {
	cm.config.Level = level
}

// SetTimeFormat sets the time format
func (cm *ConfigManager) SetTimeFormat(timeFormat string) {
	cm.config.TimeFormat = timeFormat
}

// SetEnableColor sets whether to enable colors
func (cm *ConfigManager) SetEnableColor(enableColor bool) {
	cm.config.EnableColor = enableColor
}

// SetOutput sets the output destination
func (cm *ConfigManager) SetOutput(output string) {
	cm.config.Output = output
}

// IsLevelEnabled checks if the given log level is enabled
func (cm *ConfigManager) IsLevelEnabled(level LogLevel) bool {
	currentLevel, _ := logrus.ParseLevel(string(cm.config.Level))
	checkLevel, _ := logrus.ParseLevel(string(level))
	return checkLevel >= currentLevel
}

// GetLogrusLevel returns the logrus log level
func (cm *ConfigManager) GetLogrusLevel() logrus.Level {
	level, _ := logrus.ParseLevel(string(cm.config.Level))
	return level
}

// Validate validates the configuration
func (cm *ConfigManager) Validate() error {
	// Validate format
	validFormats := []LogFormat{FormatJSON, FormatText}
	formatValid := false
	for _, validFormat := range validFormats {
		if cm.config.Format == validFormat {
			formatValid = true
			break
		}
	}
	if !formatValid {
		return fmt.Errorf("invalid log format: %s", cm.config.Format)
	}

	// Validate level
	validLevels := []LogLevel{
		LogLevelDebug, LogLevelInfo, LogLevelWarn,
		LogLevelError, LogLevelFatal, LogLevelPanic,
	}
	levelValid := false
	for _, validLevel := range validLevels {
		if cm.config.Level == validLevel {
			levelValid = true
			break
		}
	}
	if !levelValid {
		return fmt.Errorf("invalid log level: %s", cm.config.Level)
	}

	// Validate time format
	if cm.config.TimeFormat == "" {
		return fmt.Errorf("time format cannot be empty")
	}

	// Validate output
	validOutputs := []string{"stdout", "stderr"}
	outputValid := false
	for _, validOutput := range validOutputs {
		if strings.ToLower(cm.config.Output) == validOutput {
			outputValid = true
			break
		}
	}
	if !outputValid {
		return fmt.Errorf("invalid log output: %s", cm.config.Output)
	}

	return nil
}

// Clone creates a copy of the configuration manager
func (cm *ConfigManager) Clone() *ConfigManager {
	// Deep copy the config
	newConfig := &LogConfig{
		Format:      cm.config.Format,
		Level:       cm.config.Level,
		TimeFormat:  cm.config.TimeFormat,
		EnableColor: cm.config.EnableColor,
		Output:      cm.config.Output,
	}

	return NewConfigManagerWithConfig(newConfig)
}

// String returns a string representation of the configuration
func (cm *ConfigManager) String() string {
	return fmt.Sprintf(
		"LogConfig{Format: %s, Level: %s, TimeFormat: %s, EnableColor: %t, Output: %s}",
		cm.config.Format,
		cm.config.Level,
		cm.config.TimeFormat,
		cm.config.EnableColor,
		cm.config.Output,
	)
}

// Environment variable names
const (
	EnvLogFormat      = "LOG_FORMAT"
	EnvLogLevel       = "LOG_LEVEL"
	EnvLogTimeFormat  = "LOG_TIME_FORMAT"
	EnvLogEnableColor = "LOG_ENABLE_COLOR"
	EnvLogOutput      = "LOG_OUTPUT"
)

// GetEnvDefaults returns a map of default environment variables
func GetEnvDefaults() map[string]string {
	return map[string]string{
		EnvLogFormat:      string(FormatJSON),
		EnvLogLevel:       string(LogLevelInfo),
		EnvLogTimeFormat:  time.RFC3339,
		EnvLogEnableColor: "false",
		EnvLogOutput:      "stdout",
	}
}

// PrintEnvHelp prints help for environment variables
func PrintEnvHelp() {
	fmt.Println("Logger Environment Variables:")
	fmt.Printf("  %s\tLog format (json|text) [default: %s]\n", EnvLogFormat, FormatJSON)
	fmt.Printf("  %s\t\tLog level (debug|info|warn|error|fatal|panic) [default: %s]\n", EnvLogLevel, LogLevelInfo)
	fmt.Printf("  %s\tTime format string [default: %s]\n", EnvLogTimeFormat, time.RFC3339)
	fmt.Printf("  %s\tEnable color output (true|false) [default: false]\n", EnvLogEnableColor)
	fmt.Printf("  %s\t\tOutput destination (stdout|stderr) [default: stdout]\n", EnvLogOutput)
}
