package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LogFormat represents the log format type
type LogFormat string

const (
	// FormatJSON uses JSON formatting
	FormatJSON LogFormat = "json"
	// FormatText uses text formatting
	FormatText LogFormat = "text"
)

// LogLevel represents the log level
type LogLevel string

const (
	// LogLevelDebug is the debug level
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo is the info level
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn is the warn level
	LogLevelWarn LogLevel = "warn"
	// LogLevelError is the error level
	LogLevelError LogLevel = "error"
	// LogLevelFatal is the fatal level
	LogLevelFatal LogLevel = "fatal"
	// LogLevelPanic is the panic level
	LogLevelPanic LogLevel = "panic"
)

// LogConfig holds logger configuration
type LogConfig struct {
	Format      LogFormat `json:"format" yaml:"format"`
	Level       LogLevel  `json:"level" yaml:"level"`
	TimeFormat  string    `json:"timeFormat" yaml:"timeFormat"`
	EnableColor bool      `json:"enableColor" yaml:"enableColor"`
	Output      string    `json:"output" yaml:"output"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Format:      FormatJSON,
		Level:       LogLevelInfo,
		TimeFormat:  time.RFC3339,
		EnableColor: false,
		Output:      "stdout",
	}
}

// DevelopmentConfig returns development logger configuration
func DevelopmentConfig() *LogConfig {
	return &LogConfig{
		Format:      FormatText,
		Level:       LogLevelDebug,
		TimeFormat:  "2006-01-02 15:04:05",
		EnableColor: true,
		Output:      "stdout",
	}
}

// ProductionConfig returns production logger configuration
func ProductionConfig() *LogConfig {
	return &LogConfig{
		Format:      FormatJSON,
		Level:       LogLevelInfo,
		TimeFormat:  time.RFC3339,
		EnableColor: false,
		Output:      "stdout",
	}
}

// JSONFormatter implements logrus.Formatter for JSON logging
type JSONFormatter struct {
	TimeFormat string `json:"timeFormat"`
}

// Format formats a log entry as JSON
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Create log entry
	data := make(map[string]interface{})
	data["timestamp"] = entry.Time.Format(f.TimeFormat)
	data["level"] = strings.ToUpper(entry.Level.String())
	data["message"] = entry.Message

	// Add fields
	for k, v := range entry.Data {
		data[k] = v
	}

	// Add file and line number if available
	if entry.HasCaller() {
		data["file"] = entry.Caller.File
		data["line"] = entry.Caller.Line
		data["function"] = entry.Caller.Function
	}

	return json.Marshal(data)
}

// TextFormatter implements logrus.Formatter for text logging
type TextFormatter struct {
	TimeFormat  string `json:"timeFormat"`
	EnableColor bool   `json:"enableColor"`
}

// Format formats a log entry as text
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimeFormat)
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	// Add color if enabled
	if f.EnableColor {
		var colorCode string
		switch entry.Level {
		case logrus.DebugLevel:
			colorCode = "\033[36m" // Cyan
		case logrus.InfoLevel:
			colorCode = "\033[32m" // Green
		case logrus.WarnLevel:
			colorCode = "\033[33m" // Yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			colorCode = "\033[31m" // Red
		default:
			colorCode = "\033[0m" // Reset
		}
		level = fmt.Sprintf("%s%s\033[0m", colorCode, level)
	}

	// Build the log message
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s] %s: %s", timestamp, level, message))

	// Add fields if any
	if len(entry.Data) > 0 {
		builder.WriteString(" |")
		for k, v := range entry.Data {
			builder.WriteString(fmt.Sprintf(" %s=%v", k, v))
		}
	}

	// Add file and line number if available
	if entry.HasCaller() {
		builder.WriteString(fmt.Sprintf(" | %s:%d", entry.Caller.File, entry.Caller.Line))
	}

	builder.WriteString("\n")
	return []byte(builder.String()), nil
}

// ConfigureLogger configures a logrus logger with the given configuration
func ConfigureLogger(config *LogConfig) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	switch config.Format {
	case FormatJSON:
		logger.SetFormatter(&JSONFormatter{
			TimeFormat: config.TimeFormat,
		})
	case FormatText:
		logger.SetFormatter(&TextFormatter{
			TimeFormat:  config.TimeFormat,
			EnableColor: config.EnableColor,
		})
	default:
		logger.SetFormatter(&JSONFormatter{
			TimeFormat: config.TimeFormat,
		})
	}

	// Set output
	switch strings.ToLower(config.Output) {
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "stdout":
		fallthrough
	default:
		logger.SetOutput(os.Stdout)
	}

	// Enable caller reporting
	logger.SetReportCaller(true)

	return logger
}

// ParseLogLevel parses a log level string
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	case "panic":
		return LogLevelPanic
	default:
		return LogLevelInfo
	}
}

// ParseLogFormat parses a log format string
func ParseLogFormat(format string) LogFormat {
	switch strings.ToLower(format) {
	case "json":
		return FormatJSON
	case "text":
		return FormatText
	default:
		return FormatJSON
	}
}
