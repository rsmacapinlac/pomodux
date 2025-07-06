# Pomodux

A powerful terminal-based timer and Pomodoro application built in Go, designed for productivity and time management.

## 🚀 Current Status

**Release 0.1.0** ✅ **COMPLETE** - Foundation and Core Timer Engine released on 2025-07-26  
**Release 0.2.0** 🔄 **IN PLANNING** - CLI Interface & Basic Functionality (Gate 1 Pending)

## 📋 Features

### ✅ Released (v0.1.0)
- **Core Timer Engine**: Robust timer with state management
- **Basic CLI Interface**: Start, stop, and status commands
- **State Persistence**: Timer state survives process restarts
- **Configuration System**: XDG-compliant configuration management
- **Cross-Platform**: Linux, macOS, and Windows support
- **Comprehensive Testing**: 80%+ test coverage with TDD approach

### 🔄 Planned (v0.2.0)
- **Enhanced Timer Control**: Pause/resume functionality
- **Pomodoro Technique**: Dedicated work/break commands
- **Tab Completion**: Shell completion for all commands
- **Configuration Commands**: View and edit configuration
- **Session History**: Track and display session statistics

### 📋 Future Releases
- **v0.3.0**: Terminal UI (TUI) with Bubbletea
- **v0.4.0**: Plugin system and advanced features

## 🛠️ Installation

### Prerequisites
- Go 1.21 or later
- Git

### Build from Source
```bash
git clone https://github.com/yourusername/pomodux.git
cd pomodux
make build
```

### Install Binary
```bash
# Download the latest release binary for your platform
# Add to your PATH
```

## 🎯 Quick Start

### Basic Timer Usage
```bash
# Start a 25-minute timer
pomodux start 25m

# Check timer status
pomodux status

# Stop the timer
pomodux stop
```

### Supported Duration Formats
- `25m` - 25 minutes
- `1h30m` - 1 hour 30 minutes
- `1500s` - 1500 seconds
- `1.5h` - 1.5 hours

## 📁 Project Structure

```
pomodux/
├── cmd/pomodux/          # Main application entry point
├── internal/
│   ├── cli/             # CLI commands and interface
│   ├── config/          # Configuration management
│   └── timer/           # Core timer engine
├── docs/                # Documentation and ADRs
├── .github/             # GitHub templates and workflows
└── Makefile            # Build and development tasks
```

## 🧪 Development

### Prerequisites
- Go 1.21+
- golangci-lint
- Make

### Development Setup
```bash
# Clone and setup
git clone https://github.com/yourusername/pomodux.git
cd pomodux

# Install dependencies
go mod download

# Run tests
make test

# Build binary
make build

# Run linter
make lint
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/timer/...
```

## 📚 Documentation

- **[Technical Specifications](docs/technical_specifications.md)** - Detailed technical documentation
- **[Implementation Roadmap](docs/implementation-roadmap.md)** - Development roadmap and releases
- **[Release Notes](RELEASE_NOTES.md)** - Complete release history
- **[Release Management](docs/release-management.md)** - Release process and documentation standards
- **[Architecture Decisions](docs/adr/)** - Architecture Decision Records (ADRs)
- **[Development Setup](docs/development-setup.md)** - Development environment setup

## 🤝 Contributing

### Development Process
Pomodux follows a structured 4-gate approval process:
1. **Gate 1**: Work Plan Approval
2. **Gate 2**: Development Completion
3. **Gate 3**: Testing/Validation
4. **Gate 4**: Final Approval (Releases)

### Issue Management
- Use GitHub issue templates for bug reports and feature requests
- Follow TDD (Test-Driven Development) approach
- Reference the backlog in `docs/backlog.md`
- Link issues to appropriate release milestones

### Code Standards
- Follow Go best practices and standards
- Maintain 80%+ test coverage
- Use golangci-lint for code quality
- Document all public APIs

## 📊 Quality Metrics

### Release 0.1.0
- **Test Coverage**: 80%+ overall, 95%+ critical paths
- **Performance**: < 2 second startup time
- **Memory Usage**: < 50MB during operation
- **Cross-Platform**: Linux, macOS, Windows support
- **Documentation**: Complete technical and user documentation

## 🔧 Configuration

Pomodux uses XDG-compliant configuration:
- **Linux/macOS**: `~/.config/pomodux/config.yaml`
- **Windows**: `%APPDATA%\pomodux\config.yaml`

### Default Configuration
```yaml
timer:
  default_duration: 25m
  auto_start: false
  notifications: true

cli:
  output_format: text
  verbose: false
```

## 🐛 Known Issues

### Release 0.1.0
- None - all issues resolved and closed

## 📈 Roadmap

### Release 0.2.0 (In Planning)
- Enhanced CLI features (pause/resume, Pomodoro)
- Tab completion for all commands
- Configuration management commands
- Session history and statistics

### Release 0.3.0 (Planned)
- Terminal UI (TUI) with Bubbletea
- Interactive menus and themes
- Real-time progress visualization
- Keyboard shortcuts

### Release 0.4.0 (Planned)
- Lua-based plugin system
- Advanced features and integrations
- Custom workflows and automation
- Extended configuration options

## 📄 License

[License information to be added]

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Following [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- Inspired by the Pomodoro Technique

---

**Note**: Pomodux is actively developed following a structured release process. For the latest updates, check the [releases page](https://github.com/yourusername/pomodux/releases) and [issue tracker](https://github.com/yourusername/pomodux/issues).