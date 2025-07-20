package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogLevelConstants tests the log level constants
func TestLogLevelConstants(t *testing.T) {
	assert.Equal(t, LogLevel("debug"), LogLevelDebug)
	assert.Equal(t, LogLevel("info"), LogLevelInfo)
	assert.Equal(t, LogLevel("warn"), LogLevelWarn)
	assert.Equal(t, LogLevel("error"), LogLevelError)
}

// TestDefaultConfig tests the default configuration
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, LogLevelInfo, config.Level)
	assert.Equal(t, "text", config.Format)
	assert.Equal(t, "file", config.Output)
	assert.Equal(t, "", config.LogFile)
	assert.False(t, config.ShowCaller)
}

// TestInitWithExtendedLogLevels tests logger initialization with extended log levels
func TestInitWithExtendedLogLevels(t *testing.T) {
	tests := []struct {
		name      string
		level     LogLevel
		expectErr bool
	}{
		{"debug level", LogLevelDebug, false},
		{"info level", LogLevelInfo, false},
		{"warn level", LogLevelWarn, false},
		{"error level", LogLevelError, false},
		{"invalid level", LogLevel("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Level:      tt.level,
				Format:     "text",
				Output:     "console",
				LogFile:    "",
				ShowCaller: false,
			}

			err := Init(config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, Logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, Logger)
			}
		})
	}
}

// TestTraceLogging tests trace level logging
func TestTraceLogging(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	config := &Config{
		Level:      LogLevelDebug,
		Format:     "text",
		Output:     "console",
		LogFile:    "",
		ShowCaller: false,
	}

	err := Init(config)
	require.NoError(t, err)
	require.NotNil(t, Logger)

	// Set the logger output to our buffer
	Logger.SetOutput(&buf)

	// Test debug logging
	Debug("Debug message")
	Info("Info message")

	output := buf.String()

	// Verify debug message is logged
	assert.Contains(t, output, "Debug message")
	assert.Contains(t, output, "Info message")
}

// TestLogLevelValidation tests log level validation
func TestLogLevelValidation(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		expectErr bool
	}{
		{"debug", "debug", false},
		{"info", "info", false},
		{"warn", "warn", false},
		{"error", "error", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
		{"mixed case", "Debug", false}, // logrus is case-insensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Level:      LogLevel(tt.level),
				Format:     "text",
				Output:     "console",
				LogFile:    "",
				ShowCaller: false,
			}

			err := Init(config)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLogLevelPersistence tests that log level persists across logger usage
func TestLogLevelPersistence(t *testing.T) {
	config := &Config{
		Level:      LogLevelDebug,
		Format:     "text",
		Output:     "console",
		LogFile:    "",
		ShowCaller: false,
	}

	err := Init(config)
	require.NoError(t, err)
	require.NotNil(t, Logger)

	// Verify logger level is set correctly
	assert.Equal(t, logrus.DebugLevel, Logger.GetLevel())
}

// TestExtendedLogLevelsWithFields tests extended log levels with fields
func TestExtendedLogLevelsWithFields(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	config := &Config{
		Level:      LogLevelDebug,
		Format:     "text",
		Output:     "console",
		LogFile:    "",
		ShowCaller: false,
	}

	err := Init(config)
	require.NoError(t, err)
	require.NotNil(t, Logger)

	// Set the logger output to our buffer
	Logger.SetOutput(&buf)

	// Test logging with fields
	WithField("component", "timer").Debug("Debug with field")
	WithFields(map[string]interface{}{
		"component": "cli",
		"action":    "start",
	}).Info("Info with fields")

	output := buf.String()

	// Verify fields are included
	assert.Contains(t, output, "component=timer")
	assert.Contains(t, output, "component=cli")
	assert.Contains(t, output, "action=start")
}

// TestLogLevelFiltering tests that log levels filter correctly
func TestLogLevelFiltering(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	config := &Config{
		Level:      LogLevelWarn, // Only warn and above
		Format:     "text",
		Output:     "console",
		LogFile:    "",
		ShowCaller: false,
	}

	err := Init(config)
	require.NoError(t, err)
	require.NotNil(t, Logger)

	// Set the logger output to our buffer
	Logger.SetOutput(&buf)

	// Test different log levels
	Debug("Debug message - should not appear")
	Info("Info message - should not appear")
	Warn("Warn message - should appear")
	Error("Error message", nil)

	output := buf.String()

	// Verify filtering
	assert.NotContains(t, output, "Debug message")
	assert.NotContains(t, output, "Info message")
	assert.Contains(t, output, "Warn message")
	assert.Contains(t, output, "Error message")
}

// TestLogLevelHierarchy tests the log level hierarchy
func TestLogLevelHierarchy(t *testing.T) {
	tests := []struct {
		name           string
		level          LogLevel
		shouldLogInfo  bool
		shouldLogWarn  bool
		shouldLogError bool
	}{
		{"debug", LogLevelDebug, true, true, true},
		{"info", LogLevelInfo, true, true, true},
		{"warn", LogLevelWarn, false, true, true},
		{"error", LogLevelError, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			config := &Config{
				Level:      tt.level,
				Format:     "text",
				Output:     "console",
				LogFile:    "",
				ShowCaller: false,
			}

			err := Init(config)
			require.NoError(t, err)
			require.NotNil(t, Logger)

			// Set the logger output to our buffer
			Logger.SetOutput(&buf)

			Info("info message for " + string(tt.level))
			Warn("warn message for " + string(tt.level))
			Error("error message for "+string(tt.level), nil)

			output := buf.String()

			// Check expected log levels
			if tt.shouldLogInfo {
				assert.Contains(t, output, "info message for "+string(tt.level))
			} else {
				assert.NotContains(t, output, "info message for "+string(tt.level))
			}

			if tt.shouldLogWarn {
				assert.Contains(t, output, "warn message for "+string(tt.level))
			} else {
				assert.NotContains(t, output, "warn message for "+string(tt.level))
			}

			if tt.shouldLogError {
				assert.Contains(t, output, "error message for "+string(tt.level))
			} else {
				assert.NotContains(t, output, "error message for "+string(tt.level))
			}
		})
	}
}

// TestLoggerNilSafety tests that logger functions handle nil logger gracefully
func TestLoggerNilSafety(t *testing.T) {
	// Set Logger to nil
	Logger = nil

	// These should not panic
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message", nil)
	WithField("key", "value")
	WithFields(map[string]interface{}{"key": "value"})
}

// TestLogLevelStringConversion tests log level string conversion
func TestLogLevelStringConversion(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
	}{
		{"debug", LogLevelDebug},
		{"info", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"error", LogLevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Level:      tt.level,
				Format:     "text",
				Output:     "console",
				LogFile:    "",
				ShowCaller: false,
			}

			err := Init(config)
			require.NoError(t, err)
			require.NotNil(t, Logger)

			// Verify the level was set correctly
			expectedLevel, _ := logrus.ParseLevel(string(tt.level))
			assert.Equal(t, expectedLevel, Logger.GetLevel())
		})
	}
}
