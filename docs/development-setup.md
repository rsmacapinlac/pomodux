# Development Environment Setup for Pomodux

> **Note:** This guide helps developers set up their environment for contributing to the Pomodux project.

## 1.0 Prerequisites

### 1.1 Required Software

- **Go 1.21+**: [Download from golang.org](https://golang.org/dl/)
- **Git**: [Download from git-scm.com](https://git-scm.com/)
- **Make**: Usually pre-installed on Linux/macOS, [download for Windows](https://www.gnu.org/software/make/)
- **golangci-lint**: [Installation instructions](https://golangci-lint.run/usage/install/)

### 1.2 Optional Software

- **Docker**: For containerized development and testing
- **Visual Studio Code**: Recommended IDE with Go extension
- **GoLand**: JetBrains IDE for Go development

## 2.0 Environment Setup

### 2.1 Go Environment Configuration

Set up your Go environment variables:

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export GOPATH=$HOME/go
export GOROOT=/usr/local/go  # Adjust path as needed
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

### 2.2 Project Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/rsmacapinlac/pomodux.git
   cd pomodux
   ```

2. **Initialize Go module**:
   ```bash
   go mod init github.com/yourusername/pomodux
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   go mod tidy
   ```

4. **Install development tools**:
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install additional tools (optional)
   go install golang.org/x/tools/cmd/goimports@latest
   go install golang.org/x/tools/cmd/godoc@latest
   ```

### 2.3 IDE Configuration

#### Visual Studio Code

1. **Install Go extension**:
   - Open VS Code
   - Go to Extensions (Ctrl+Shift+X)
   - Search for "Go" by Google
   - Install the extension

2. **Configure Go tools**:
   - Open Command Palette (Ctrl+Shift+P)
   - Run "Go: Install/Update Tools"
   - Select all tools and install

3. **Recommended settings** (add to `.vscode/settings.json`):
   ```json
   {
     "go.formatTool": "gofmt",
     "go.lintTool": "golangci-lint",
     "go.lintFlags": ["--fast"],
     "go.testFlags": ["-v"],
     "go.buildTags": "",
     "go.toolsManagement.autoUpdate": true,
     "go.useLanguageServer": true,
     "go.gopath": "",
     "go.goroot": "",
     "go.gocodeAutoBuild": false
   }
   ```

#### GoLand

1. **Configure Go SDK**:
   - Go to File → Settings → Languages & Frameworks → Go
   - Set GOROOT to your Go installation path
   - Set GOPATH to your Go workspace

2. **Enable code inspection**:
   - Go to File → Settings → Editor → Inspections
   - Enable Go-specific inspections

## 3.0 Project Structure Setup

### 3.1 Create Directory Structure

Run the following commands to create the standard Go project structure:

```bash
# Create main application directory
mkdir -p cmd/pomodux

# Create internal packages
mkdir -p internal/timer
mkdir -p internal/tui
mkdir -p internal/cli
mkdir -p internal/plugins
mkdir -p internal/config

# Create public packages
mkdir -p pkg/notifications
mkdir -p pkg/events

# Create supporting directories
mkdir -p configs
mkdir -p scripts
mkdir -p build
mkdir -p test
mkdir -p tools
mkdir -p examples
mkdir -p githooks
mkdir -p bin
```

### 3.2 Initialize Core Files

1. **Create main.go**:
   ```bash
   touch cmd/pomodux/main.go
   ```

2. **Create package files**:
   ```bash
   touch internal/timer/timer.go
   touch internal/tui/tui.go
   touch internal/cli/cli.go
   touch internal/plugins/plugins.go
   touch internal/config/config.go
   touch pkg/notifications/notifications.go
   touch pkg/events/events.go
   ```

3. **Create configuration files**:
   ```bash
   touch .golangci.yml
   touch Makefile
   touch .gitignore
   ```

## 4.0 Development Workflow

### 4.1 Git Hooks Setup

1. **Configure git hooks**:
   ```bash
   git config core.hooksPath githooks
   chmod +x githooks/pre-commit
   chmod +x githooks/pre-push
   ```

2. **Create pre-commit hook**:
   ```bash
   cat > githooks/pre-commit << 'EOF'
   #!/bin/bash
   set -e
   
   echo "Running Go formatting..."
   gofmt -s -w .
   
   echo "Running linter..."
   golangci-lint run
   
   echo "Running tests..."
   go test ./...
   
   echo "Checking for uncommitted changes..."
   if ! git diff --quiet; then
       echo "Code formatting changed files. Please review and commit changes."
       exit 1
   fi
   EOF
   ```

3. **Create pre-push hook**:
   ```bash
   cat > githooks/pre-push << 'EOF'
   #!/bin/bash
   set -e
   
   echo "Running full test suite..."
   go test -v ./...
   
   echo "Running linter..."
   golangci-lint run
   
   echo "Building application..."
   go build -o bin/pomodux cmd/pomodux/main.go
   EOF
   ```

### 4.2 Makefile Setup

Create a comprehensive Makefile for development tasks:

```makefile
.PHONY: help build test lint clean install run build-all

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  clean       - Clean build artifacts"
	@echo "  install     - Install dependencies"
	@echo "  run         - Run the application"
	@echo "  build-all   - Build for all platforms"

# Build the application
build:
	@echo "Building pomodux..."
	go build -o bin/pomodux cmd/pomodux/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the application
run:
	@echo "Running pomodux..."
	go run cmd/pomodux/main.go

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/pomodux-linux-amd64 cmd/pomodux/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/pomodux-darwin-amd64 cmd/pomodux/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/pomodux-windows-amd64.exe cmd/pomodux/main.go

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Install development tools
tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/cmd/godoc@latest

# Setup development environment
setup: tools install
	@echo "Setting up development environment..."
	git config core.hooksPath githooks
	chmod +x githooks/pre-commit
	chmod +x githooks/pre-push
	@echo "Development environment setup complete!"
```

### 4.3 Linting Configuration

Create `.golangci.yml` for comprehensive linting:

```yaml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor/
    - third_party/
    - test/

linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell
    - gosec
    - goimports
    - gocyclo
    - dupl
    - goconst

linters-settings:
  golint:
    min-confidence: 0.8
  gocyclo:
    min-complexity: 15
  dupl:
    threshold: 100
  goconst:
    min-len: 3
    min-occurrences: 3

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gocyclo
        - dupl
```

## 5.0 Testing Setup

### 5.1 Test Configuration

1. **Create test helpers**:
   ```bash
   mkdir -p test/helpers
   touch test/helpers/helpers.go
   ```

2. **Create integration test directory**:
   ```bash
   mkdir -p test/integration
   touch test/integration/timer_test.go
   ```

### 5.2 Test Utilities

Create `test/helpers/helpers.go`:

```go
package helpers

import (
    "os"
    "path/filepath"
    "testing"
)

// TempDir creates a temporary directory for testing
func TempDir(t *testing.T) string {
    t.Helper()
    dir, err := os.MkdirTemp("", "pomodux-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    t.Cleanup(func() {
        os.RemoveAll(dir)
    })
    return dir
}

// TempFile creates a temporary file for testing
func TempFile(t *testing.T, dir, pattern string) *os.File {
    t.Helper()
    file, err := os.CreateTemp(dir, pattern)
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    t.Cleanup(func() {
        file.Close()
        os.Remove(file.Name())
    })
    return file
}

// WriteTestConfig writes a test configuration file
func WriteTestConfig(t *testing.T, dir, content string) string {
    t.Helper()
    configPath := filepath.Join(dir, "config.yaml")
    err := os.WriteFile(configPath, []byte(content), 0644)
    if err != nil {
        t.Fatalf("Failed to write test config: %v", err)
    }
    return configPath
}
```

## 6.0 Debugging Setup

### 6.1 VS Code Debug Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/pomodux/main.go",
            "args": ["--help"]
        },
        {
            "name": "Test Current File",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${file}"
        },
        {
            "name": "Test Package",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/internal/timer"
        }
    ]
}
```

### 6.2 Profiling Setup

Add profiling targets to Makefile:

```makefile
# Profile CPU usage
profile-cpu:
	@echo "Profiling CPU usage..."
	go test -cpuprofile=cpu.prof -bench=. ./internal/timer
	go tool pprof cpu.prof

# Profile memory usage
profile-mem:
	@echo "Profiling memory usage..."
	go test -memprofile=mem.prof -bench=. ./internal/timer
	go tool pprof mem.prof

# Profile HTTP server (if applicable)
profile-http:
	@echo "Starting HTTP profiler..."
	go tool pprof -http=:8080 cpu.prof
```

## 7.0 Continuous Integration

### 7.1 GitHub Actions Setup

Create `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go ${{ matrix.go-version }}
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}

    - name: Install dependencies
      run: go mod download

    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.txt
        flags: unittests
        name: codecov-umbrella
```

## 8.0 Troubleshooting

### 8.1 Common Issues

#### Go Module Issues
```bash
# Clear module cache
go clean -modcache

# Reset go.mod
rm go.mod go.sum
go mod init github.com/yourusername/pomodux
go mod tidy
```

#### Linting Issues
```bash
# Update golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run with verbose output
golangci-lint run --verbose
```

#### Test Issues
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/timer -run TestTimerStart

# Run tests with race detection
go test -race ./...
```

### 8.2 Performance Issues

#### Memory Leaks
```bash
# Run with memory profiling
go test -memprofile=mem.prof ./internal/timer
go tool pprof mem.prof
```

#### CPU Performance
```bash
# Run benchmarks
go test -bench=. ./internal/timer

# Run with CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./internal/timer
go tool pprof cpu.prof
```

## 9.0 Next Steps

After completing the setup:

1. **Review the Go standards document**: `docs/go-standards.md`
2. **Start with a simple implementation**: Begin with the timer engine
3. **Write tests first**: Follow TDD practices
4. **Use the Makefile**: Leverage the provided build targets
5. **Contribute to documentation**: Keep docs up to date

---

## References

### Project Documentation
- **[Go Standards Document](docs/go-standards.md)** - Comprehensive Go development standards and implementation guidelines
- **[Cursor Rules](.cursor/rules/go-development.mdc)** - AI assistant guidance and project-specific patterns
- **[Implementation Roadmap](docs/implementation-roadmap.md)** - Development phases and priorities
- **[Architecture Decision Records](docs/adr/)** - Project architectural decisions and rationale

### External Resources
- [Go Installation Guide](https://golang.org/doc/install)
- [Go Modules Documentation](https://golang.org/doc/modules)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TableDrivenTests)

### How These Files Work Together
- **This file (Development Setup)**: Provides environment setup and development workflow guidance
- **[Go Standards](docs/go-standards.md)**: Provides detailed implementation standards and examples
- **[Cursor Rules](.cursor/rules/go-development.mdc)**: Provides AI assistants with project-specific patterns
- **Together**: Create a complete development ecosystem for both human developers and AI assistants
