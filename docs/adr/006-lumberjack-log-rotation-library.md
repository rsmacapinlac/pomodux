# ADR 006: Lumberjack Library for Log Rotation

> **Status**: ✅ **APPROVED**  
> **Component**: Logging System  
> **Related ADRs**: [ADR 005 - Structured Logger Architecture](005-structured-logger-architecture.md)

## Context

Future releases may require automatic log file rotation functionality to prevent log files from growing too large and consuming excessive disk space. The current logging system uses logrus for structured logging but lacks built-in rotation capabilities. We need to select a log rotation library that integrates well with logrus while maintaining timer accuracy and cross-platform compatibility.

## Decision

**We will use the lumberjack library for log rotation functionality.**

## Rationale

### Why Lumberjack?

1. **Logrus Integration**: Lumberjack is specifically designed to work with logrus and is the recommended solution in logrus documentation
2. **Simplicity**: Provides a simple, configurable writer that can be used as logrus output
3. **Cross-Platform**: Works consistently across Linux, macOS, and Windows
4. **Performance**: Non-blocking operations that won't impact timer accuracy
5. **Mature**: Well-established library with active maintenance
6. **Configuration**: Supports size-based and time-based rotation with compression

### Alternative Considered

**Custom Implementation**: Building log rotation from scratch
- **Pros**: Full control over implementation
- **Cons**: Significant development time, potential for bugs, maintenance burden
- **Decision**: Rejected due to complexity and maintenance overhead

## Consequences

### Positive Consequences

1. **Rapid Implementation**: Lumberjack provides immediate log rotation capabilities
2. **Proven Reliability**: Mature library with extensive usage in production
3. **Minimal Code Changes**: Simple integration with existing logrus setup
4. **Rich Configuration**: Comprehensive rotation options (size, time, compression, retention)
5. **Performance**: Optimized for high-performance logging scenarios

### Negative Consequences

1. **External Dependency**: Adds lumberjack as a new dependency
2. **Library Maintenance**: Dependent on lumberjack's continued maintenance
3. **Configuration Complexity**: Additional configuration options for users

### Mitigation Strategies

1. **Dependency Management**: Pin lumberjack version and monitor for updates
2. **Configuration Defaults**: Provide sensible defaults to minimize complexity
3. **Documentation**: Clear documentation for rotation configuration
4. **Testing**: Comprehensive testing of rotation functionality

## Implementation Details

### Integration Approach

```go
// Current logrus setup
Logger.SetOutput(file)

// New lumberjack integration
lumberjackLogger := &lumberjack.Logger{
    Filename:   config.LogFile,
    MaxSize:    config.Rotation.MaxSize,    // megabytes
    MaxAge:     config.Rotation.MaxAge,     // days
    MaxBackups: config.Rotation.MaxBackups,
    Compress:   config.Rotation.Compress,
    LocalTime:  config.Rotation.LocalTime,
}
Logger.SetOutput(lumberjackLogger)
```

### Configuration Structure

```yaml
logging:
  # Existing configuration
  level: info
  format: text
  output: file
  log_file: "~/.config/pomodux/pomodux.log"
  show_caller: false
  
  # New rotation configuration
  rotation:
    enabled: true
    max_size: "10MB"              # Maximum file size before rotation
    max_age: "7d"                 # Maximum age of rotated files
    max_backups: 5                # Maximum number of rotated files to keep
    compress: true                # Compress rotated files
    local_time: true              # Use local time for rotation
    rotate_on_startup: false      # Rotate on application startup
```

### Dependencies

```go
// go.mod addition
require (
    github.com/natefinch/lumberjack v2.0.0+incompatible
)
```

## Compliance with Existing ADRs

### ADR 001: Programming Language Selection
- **Compliance**: ✅ Lumberjack is a Go library, maintaining Go ecosystem consistency
- **Impact**: No impact on cross-platform compatibility

### ADR 002: Persistent Timer Design
- **Compliance**: ✅ Lumberjack operations are non-blocking and won't impact timer accuracy
- **Impact**: Maintains timer precision requirements

### ADR 004: Plugin System Architecture
- **Compliance**: ✅ No changes to plugin API or event logging
- **Impact**: Plugins continue to work without modification

### ADR 005: Structured Logger Architecture
- **Compliance**: ✅ Extends existing logrus-based architecture
- **Impact**: Maintains structured logging capabilities while adding rotation

## Testing Strategy

### Unit Tests
- [ ] Test lumberjack configuration parsing
- [ ] Test rotation trigger conditions (size, time)
- [ ] Test compression and retention policies
- [ ] Test error handling and edge cases

### Integration Tests
- [ ] Test lumberjack integration with logrus
- [ ] Test rotation during active logging
- [ ] Test file system operations and permissions
- [ ] Test configuration validation

### Performance Tests
- [ ] Test rotation impact on logging performance
- [ ] Test rotation impact on timer accuracy
- [ ] Test memory usage during rotation operations
- [ ] Test disk I/O performance with compression

### Cross-Platform Tests
- [ ] Test rotation on Linux, macOS, and Windows
- [ ] Test file path handling across platforms
- [ ] Test file locking behavior
- [ ] Test permission handling

## Migration Strategy

### Backward Compatibility
- **Existing Configurations**: Continue to work with default rotation settings
- **Log Files**: Existing log files remain accessible
- **CLI Commands**: No changes to existing logging commands

### Configuration Migration
- **Automatic**: New rotation settings default to sensible values
- **Optional**: Users can disable rotation entirely
- **Gradual**: Configuration can be updated incrementally

## Success Criteria

### Functional Requirements
- [ ] Log files rotate automatically based on size and time
- [ ] Rotated files are compressed and retained according to policy
- [ ] No log data is lost during rotation
- [ ] Rotation works consistently across all platforms

### Performance Requirements
- [ ] Rotation operations don't impact timer accuracy (< 1ms deviation)
- [ ] Logging performance remains within acceptable limits (< 5ms per entry)
- [ ] Memory usage doesn't increase significantly during rotation
- [ ] Disk I/O performance is acceptable with compression

### Quality Requirements
- [ ] Comprehensive test coverage for rotation functionality
- [ ] Clear documentation for rotation configuration
- [ ] Error handling for all failure scenarios
- [ ] Monitoring and alerting for rotation issues

## Future Considerations

### Potential Enhancements
1. **Advanced Rotation Strategies**: Custom rotation schedules
2. **Remote Logging**: Integration with remote log aggregation
3. **Metrics Integration**: Rotation statistics and monitoring
4. **Custom Compression**: Alternative compression algorithms

### Maintenance Considerations
1. **Version Updates**: Regular updates to lumberjack library
2. **Security**: Monitor for security vulnerabilities
3. **Performance**: Continuous performance monitoring
4. **Compatibility**: Ensure compatibility with future logrus versions

## References

- [Lumberjack GitHub Repository](https://github.com/natefinch/lumberjack)
- [Logrus Documentation](https://github.com/sirupsen/logrus)
- [ADR 005: Structured Logger Architecture](005-structured-logger-architecture.md)


---

**Note**: This ADR represents the architectural decision to use lumberjack for log rotation functionality. Implementation details and configuration options will be documented in the release documentation. 