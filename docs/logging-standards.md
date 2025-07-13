# Logging Standards for Pomodux

> **Status**: ✅ **ACTIVE** - Implemented in Release 0.3.0  
> **Last Updated**: 2025-01-27  
> **Related**: [Go Standards](go-standards.md), [Release 0.3.0](releases/release-0.3.0.md), [ADR 005](adr/005-structured-logger-architecture.md)

## Overview

This document defines the logging standards and best practices for the Pomodux project. All logging in Pomodux uses structured logging with the Logrus library, providing consistent, configurable, and informative logging across the application.

The logging architecture was designed and approved through [ADR 005: Structured Logger Architecture](adr/005-structured-logger-architecture.md), which established the use of Logrus as the logging library and defined the core requirements for structured logging, log levels, and configuration.

## Core Principles

### 1. Structured Logging
- **Library**: Use [Logrus](https://github.com/sirupsen/logrus) for all logging
- **Fields**: Always use structured fields (key-value pairs) for context
- **Consistency**: Maintain consistent field names across the application
- **Searchability**: Use fields that enable log searching and filtering

### 2. Log Levels
Use the following standard log levels consistently:

- **Debug**: Detailed information for debugging and development
  - Use for: Function entry/exit, variable values, detailed flow
  - Example: `log.WithField("duration", duration).Debug("Timer duration set")`

- **Info**: General operational events and normal application flow
  - Use for: Timer start/stop, configuration changes, user actions
  - Example: `log.WithField("timer_id", timerID).Info("Timer started")`

- **Warn**: Unexpected situations that don't cause failure
  - Use for: Deprecated features, configuration issues, recoverable errors
  - Example: `log.WithField("config_path", path).Warn("Using default configuration")`

- **Error**: Errors that require attention and may cause issues
  - Use for: File I/O errors, network failures, invalid states
  - Example: `log.WithError(err).Error("Failed to save timer state")`

## Implementation Standards

### Logger Initialization
Always use the provided logger package for initialization:

```go
import "pomodux/internal/logger"

// Get the configured logger instance
log := logger.Get()

// Use structured logging with fields
log.WithFields(logrus.Fields{
    "timer_id": timerID,
    "duration": duration,
    "session_type": "work",
}).Info("Timer session started")
```

### Field Naming Conventions
Use consistent field names across the application:

- **IDs**: `timer_id`, `session_id`, `user_id`
- **Durations**: `duration`, `elapsed`, `remaining`
- **Timestamps**: `started_at`, `completed_at`, `created_at`
- **Types**: `session_type`, `event_type`, `config_type`
- **Paths**: `config_path`, `log_path`, `data_path`
- **Status**: `status`, `state`, `phase`

### Error Logging
Always include error context when logging errors:

```go
// Good: Include error and context
if err != nil {
    log.WithError(err).WithFields(logrus.Fields{
        "timer_id": timerID,
        "operation": "start",
    }).Error("Failed to start timer")
    return fmt.Errorf("start timer: %w", err)
}

// Avoid: Plain error logging without context
if err != nil {
    log.Error("Error occurred") // ❌ Too vague
    return err
}
```

## Configuration

### Log Configuration Structure
Logging is configured through the application configuration:

```yaml
logging:
  level: info                    # debug, info, warn, error
  output: file                   # file, console, both
  file: "~/.config/pomodux/pomodux.log"  # Log file path
```

### Configuration Features
- **Automatic Directory Creation**: Log directories are created automatically
- **Path Expansion**: `~` is expanded to user home directory
- **Default File Output**: Preserves terminal UX by defaulting to file output
- **Configurable Levels**: Runtime log level adjustment
- **Multiple Outputs**: Support for file, console, or both

## Best Practices

### 1. Contextual Logging
Always include relevant context in log entries:

```go
// Good: Rich context for debugging
log.WithFields(logrus.Fields{
    "timer_id": timer.ID,
    "duration": timer.Duration,
    "session_type": timer.SessionType,
    "user_id": timer.UserID,
}).Info("Timer session completed")

// Avoid: Minimal context
log.Info("Timer completed") // ❌ Not enough context
```

### 2. Performance Considerations
- **Avoid Log Spam**: Don't log in tight loops or performance-critical paths
- **Conditional Logging**: Use log level checks for expensive operations
- **Structured Data**: Use fields instead of string concatenation

```go
// Good: Conditional logging for performance
if log.GetLevel() >= logrus.DebugLevel {
    log.WithField("timer_state", timer.State()).Debug("Timer state updated")
}

// Avoid: Always logging in loops
for _, timer := range timers {
    log.Debug("Processing timer") // ❌ Could be expensive
}
```

### 3. Security and Privacy
- **No Sensitive Data**: Never log passwords, tokens, or personal information
- **Sanitized Paths**: Log file paths, not full file contents
- **User Consent**: Consider privacy implications of logged data

### 4. Testing
- **Test Coverage**: Ensure critical log paths are covered by tests
- **Mock Logging**: Use test hooks or mocks to avoid log pollution
- **Log Verification**: Test that appropriate logs are generated

```go
// Example: Testing log output
func TestTimerStart(t *testing.T) {
    // Capture log output for testing
    var logOutput bytes.Buffer
    log.SetOutput(&logOutput)
    
    timer.Start(25 * time.Minute)
    
    // Verify log contains expected information
    if !strings.Contains(logOutput.String(), "Timer started") {
        t.Error("Expected log message not found")
    }
}
```

## Common Patterns

### Timer Operations
```go
func (t *Timer) Start(duration time.Duration) error {
    log := logger.Get()
    
    log.WithFields(logrus.Fields{
        "timer_id": t.ID,
        "duration": duration,
        "session_type": t.SessionType,
    }).Info("Starting timer session")
    
    if err := t.validateDuration(duration); err != nil {
        log.WithError(err).WithField("duration", duration).Error("Invalid timer duration")
        return fmt.Errorf("validate duration: %w", err)
    }
    
    // ... implementation ...
    
    log.WithField("timer_id", t.ID).Info("Timer session started successfully")
    return nil
}
```

### Configuration Changes
```go
func (c *Config) SetLogLevel(level string) error {
    log := logger.Get()
    
    log.WithField("new_level", level).Info("Changing log level")
    
    if err := c.validateLogLevel(level); err != nil {
        log.WithError(err).WithField("level", level).Error("Invalid log level")
        return fmt.Errorf("validate log level: %w", err)
    }
    
    // ... implementation ...
    
    log.WithField("level", level).Info("Log level changed successfully")
    return nil
}
```

### Error Handling
```go
func (m *Manager) LoadConfiguration(path string) error {
    log := logger.Get()
    
    log.WithField("config_path", path).Debug("Loading configuration")
    
    data, err := os.ReadFile(path)
    if err != nil {
        log.WithError(err).WithField("config_path", path).Error("Failed to read configuration file")
        return fmt.Errorf("read config file: %w", err)
    }
    
    if err := m.parseConfiguration(data); err != nil {
        log.WithError(err).WithField("config_path", path).Error("Failed to parse configuration")
        return fmt.Errorf("parse config: %w", err)
    }
    
    log.WithField("config_path", path).Info("Configuration loaded successfully")
    return nil
}
```

## Integration with Other Systems

### Plugin System
When implementing plugin hooks, ensure logging includes plugin context:

```go
func (pm *PluginManager) ExecuteHook(event Event) error {
    log := logger.Get()
    
    log.WithFields(logrus.Fields{
        "event_type": event.Type,
        "plugin_count": len(pm.plugins),
    }).Debug("Executing plugin hooks")
    
    for _, plugin := range pm.plugins {
        log.WithField("plugin_name", plugin.Name).Debug("Executing plugin hook")
        
        if err := plugin.Execute(event); err != nil {
            log.WithError(err).WithField("plugin_name", plugin.Name).Error("Plugin hook failed")
            // Continue with other plugins
        }
    }
    
    return nil
}
```

### CLI Commands
CLI commands should log their execution and results:

```go
func (c *StartCmd) Run(cmd *cobra.Command, args []string) error {
    log := logger.Get()
    
    log.WithField("command", "start").Debug("Executing start command")
    
    // ... command implementation ...
    
    log.WithFields(logrus.Fields{
        "command": "start",
        "duration": duration,
        "session_type": sessionType,
    }).Info("Start command completed successfully")
    
    return nil
}
```

## Troubleshooting

### Common Issues

1. **Missing Log Output**
   - Check log level configuration
   - Verify log file path and permissions
   - Ensure logger is properly initialized

2. **Performance Impact**
   - Review log level usage in performance-critical paths
   - Use conditional logging for expensive operations
   - Consider log output destination (file vs console)

3. **Inconsistent Logging**
   - Follow field naming conventions
   - Use structured logging consistently
   - Include appropriate context in all log entries

### Debugging Tips

1. **Temporary Debug Logging**
   ```go
   // Add temporary debug logging for troubleshooting
   log.WithField("variable_name", value).Debug("Debug information")
   ```

2. **Log Level Adjustment**
   ```bash
   # Temporarily increase log level for debugging
   pomodux config set logging.level debug
   ```

3. **Log File Analysis**
   ```bash
   # View recent logs
   tail -f ~/.config/pomodux/pomodux.log
   
   # Search for specific events
   grep "Timer started" ~/.config/pomodux/pomodux.log
   ```

## References

- **[ADR 005: Structured Logger Architecture](adr/005-structured-logger-architecture.md)** - Architectural decision record for logging system design
- **[Logrus Documentation](https://github.com/sirupsen/logrus)** - Official Logrus documentation
- **[Go Standards](go-standards.md)** - General Go development standards
- **[Release 0.3.0](releases/release-0.3.0.md)** - Structured logging implementation details
- **[Configuration Guide](development-setup.md)** - Logging configuration setup

---

**Note**: These logging standards ensure consistent, maintainable, and useful logging across the Pomodux application. All developers should follow these guidelines when adding or modifying logging in the codebase. 