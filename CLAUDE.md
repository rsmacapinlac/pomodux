# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
```bash
make build          # Build the application to bin/pomodux
make run             # Run the application directly with go run
go run cmd/pomodux/main.go  # Alternative way to run
```

### Testing
```bash
make test           # Run all tests
make test-coverage  # Run tests with coverage report (generates coverage.html)
go test ./internal/timer/...  # Run specific package tests
```

### Code Quality
```bash
make lint           # Run golangci-lint
make fmt            # Format code with gofmt
```

### Development Setup
```bash
make setup          # Install dev tools and configure git hooks
make tools          # Install development tools only
make install        # Install dependencies
```

### Other Useful Commands
```bash
make clean          # Remove build artifacts
make build-all      # Build for multiple platforms
make vendor         # Generate vendor directory
```

## Project Architecture

### Core Components

**Timer Engine** (`internal/timer/`): Central timer functionality with state management, persistence, and session tracking. Uses interfaces for extensibility and supports work/break/long-break session types.

**CLI Interface** (`internal/cli/`): Cobra-based command structure with commands for start, stop, pause, resume, status, history, and break management. Supports shell completion and configuration.

**Configuration** (`internal/config/`): XDG-compliant configuration management with YAML support. Handles default settings and validation.

**Logger** (`internal/logger/`): Structured logging system with configurable levels, formats, and outputs. Supports file rotation and caller information.

**Plugin System** (`internal/plugin/`): Lua-based plugin architecture for extensibility. Includes event system and plugin lifecycle management.

### Application Flow

1. **main.go**: Loads configuration, initializes logger, ensures clean timer shutdown
2. **CLI Commands**: Handle user input and coordinate with timer engine
3. **Timer Engine**: Manages state persistence, session tracking, and timer operations
4. **Configuration**: Provides settings throughout the application
5. **Plugins**: Extend functionality through Lua scripts

### Key Design Patterns

- **Interface-based design**: Timer engine uses interfaces for testability
- **State persistence**: Timer state survives process restarts
- **Plugin architecture**: Lua-based extensibility system
- **XDG compliance**: Follows Linux desktop standards for configuration
- **TDD approach**: Test-driven development with 80%+ coverage requirement

### Dependencies

- **Cobra**: CLI framework for commands and completion
- **Lua**: Plugin system scripting language  
- **testify**: Testing framework with assertions
- **yaml.v3**: Configuration file parsing
- **term**: Terminal control and keypress handling

### Project Structure Notes

- `cmd/pomodux/`: Application entry point
- `internal/`: Private application code (not importable by external projects)
- `plugins/`: Example Lua plugin scripts
- `tests/`: UAT tests and fixtures
- `docs/`: Comprehensive documentation including ADRs and release management

## Development Standards

- Follow Go best practices and project layout standards
- Maintain 80%+ test coverage
- Use TDD approach for new features
- Document all public APIs
- Follow XDG Base Directory Specification for configuration
- Use structured logging with appropriate levels
- Implement proper error handling and validation