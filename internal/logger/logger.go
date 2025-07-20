package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Logger is the global logger instance
	Logger *logrus.Logger
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Config holds logging configuration
type Config struct {
	Level      LogLevel `yaml:"level"`
	Format     string   `yaml:"format"` // "text" or "json"
	Output     string   `yaml:"output"` // "console", "file", or "both"
	LogFile    string   `yaml:"log_file"`
	ShowCaller bool     `yaml:"show_caller"`
}

// DefaultConfig returns default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      LogLevelInfo,
		Format:     "text",
		Output:     "file", // Changed from "console" to "file"
		LogFile:    "",
		ShowCaller: false,
	}
}

// Init initializes the global logger with the given configuration
func Init(config *Config) error {
	// Set log level first to validate it
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		Logger = nil
		return fmt.Errorf("invalid log level %s: %w", config.Level, err)
	}

	Logger = logrus.New()
	Logger.SetLevel(level)

	// Set formatter
	var formatter logrus.Formatter
	switch config.Format {
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	case "text":
		fallthrough
	default:
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			DisableColors:   false,
		}
	}

	Logger.SetFormatter(formatter)

	// Set output
	switch config.Output {
	case "file":
		if config.LogFile == "" {
			config.LogFile = getDefaultLogFile()
		}
		config.LogFile = expandHome(config.LogFile)
		// Ensure parent directory exists
		dir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", config.LogFile, err)
		}
		Logger.SetOutput(file)
	case "both":
		if config.LogFile == "" {
			config.LogFile = getDefaultLogFile()
		}
		config.LogFile = expandHome(config.LogFile)
		// Ensure parent directory exists
		dir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}
		_, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", config.LogFile, err)
		}
		Logger.SetOutput(os.Stdout) // For now, just use stdout for both
		// TODO: Implement multi-writer for both console and file
	case "console":
		fallthrough
	default:
		Logger.SetOutput(os.Stdout)
	}

	// Set caller info
	if config.ShowCaller {
		Logger.SetReportCaller(true)
	}

	return nil
}

// ConfigFromMainConfig converts the main application config to logger config
func ConfigFromMainConfig(mainConfig interface{}) (*Config, error) {
	// This is a simplified conversion - in a real implementation,
	// you would use reflection or a more sophisticated approach
	// For now, we'll return a default config
	return DefaultConfig(), nil
}

// getDefaultLogFile returns the default log file path
func getDefaultLogFile() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "./pomodux.log"
		}
		stateHome = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(stateHome, "pomodux", "pomodux.log")
}

// expandHome expands ~ to the user's home directory in a path
func expandHome(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// Debug logs a debug message
func Debug(msg string, fields ...map[string]interface{}) {
	if Logger != nil {
		if len(fields) > 0 {
			Logger.WithFields(fields[0]).Debug(msg)
		} else {
			Logger.Debug(msg)
		}
	}
}

// Info logs an info message
func Info(msg string, fields ...map[string]interface{}) {
	if Logger != nil {
		if len(fields) > 0 {
			Logger.WithFields(fields[0]).Info(msg)
		} else {
			Logger.Info(msg)
		}
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...map[string]interface{}) {
	if Logger != nil {
		if len(fields) > 0 {
			Logger.WithFields(fields[0]).Warn(msg)
		} else {
			Logger.Warn(msg)
		}
	}
}

// Error logs an error message
func Error(msg string, err error, fields ...map[string]interface{}) {
	if Logger != nil {
		logEntry := Logger.WithError(err)
		if len(fields) > 0 {
			logEntry = logEntry.WithFields(fields[0])
		}
		logEntry.Error(msg)
	}
}

// WithField creates a new log entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithField(key, value)
	}
	return nil
}

// WithFields creates a new log entry with multiple fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithFields(fields)
	}
	return nil
}
