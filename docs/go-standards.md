# Go Development Standards for Pomodux

> **Note:** This document establishes Go development standards and practices for the Pomodux project, following the approved Go language selection in [ADR 001](adr/001-programming-language-selection.md).

## 1.0 Project Structure

### 1.1 Standard Go Project Layout

Pomodux follows the **Standard Go Project Layout** pattern:

```
pomodux/
├── cmd/                    # Main applications
│   └── pomodux/           # Main application entry point
│       └── main.go        # Application entry point
├── internal/              # Private application and library code
│   ├── timer/             # Timer engine and session management
│   ├── tui/               # Terminal UI components and rendering
│   ├── cli/               # CLI command implementations
│   ├── plugins/           # Plugin system and Lua integration
│   └── config/            # Configuration management and XDG compliance
├── pkg/                   # Library code that's ok to use by external applications
│   ├── notifications/     # Desktop notification system
│   └── events/            # Event system and hook management
├── configs/               # Configuration file templates and defaults
├── scripts/               # Build, install, and analysis scripts
├── build/                 # Packaging and CI/CD configurations
├── test/                  # Integration tests and test data
├── docs/                  # Project documentation
├── tools/                 # Development and build tools
├── examples/              # Usage examples and tutorials
├── githooks/              # Git hooks for development workflow
├── README.md
├── LICENSE
├── Makefile
├── go.mod
└── go.sum
```

### 1.2 Package Organization Principles

#### Internal Packages (`internal/`)
- **timer/**: Core timer logic, session management, state handling
- **tui/**: Terminal UI components, rendering, input handling
- **cli/**: Command implementations, argument parsing, help system
- **plugins/**: Lua runtime, hook system, plugin API
- **config/**: Configuration loading, XDG compliance, defaults

#### Public Packages (`pkg/`)
- **notifications/**: Cross-platform desktop notification system
- **events/**: Event system for plugin hooks and internal communication

## 2.0 Code Organization Standards

### 2.1 Package Naming Conventions

- Use **short, concise names** that are clear and descriptive
- Avoid **underscores** in package names
- Use **singular form** for package names
- Examples: `timer`, `tui`, `cli`, `plugins`, `config`

### 2.2 File Naming Conventions

- Use **snake_case** for file names
- Group related functionality in the same file
- Separate interfaces from implementations when beneficial
- Examples: `timer_engine.go`, `session_manager.go`, `cli_commands.go`

### 2.3 Import Organization

Organize imports in the following order:

```go
package main

import (
    // Standard library imports
    "fmt"
    "os"
    "time"

    // Third-party imports
    "github.com/spf13/cobra"
    "github.com/charmbracelet/bubbletea"

    // Internal imports
    "github.com/yourusername/pomodux/internal/timer"
    "github.com/yourusername/pomodux/internal/tui"
)
```

### 2.4 Interface Design Principles

- **Keep interfaces small** (interface segregation principle)
- **Define interfaces where they are used**, not where they are implemented
- **Use descriptive names** for interfaces
- **Document public interfaces** with examples

Example:
```go
// TimerEngine defines the core timer functionality
type TimerEngine interface {
    Start(duration time.Duration) error
    Stop() error
    Pause() error
    Resume() error
    GetStatus() TimerStatus
}
```

## 3.0 Coding Standards

### 3.1 Naming Conventions

#### Variables and Functions
- Use **camelCase** for variables and functions
- Use **PascalCase** for exported names
- Use **descriptive names** that explain purpose
- Avoid **abbreviations** unless widely understood

```go
// Good
func StartWorkTimer(duration time.Duration) error
func calculateProgress() float64
var isTimerRunning bool

// Avoid
func StartWT(d time.Duration) error
func calcProg() float64
var running bool
```

#### Constants
- Use **PascalCase** for exported constants
- Use **camelCase** for unexported constants
- Group related constants together

```go
const (
    DefaultWorkDuration = 25 * time.Minute
    DefaultBreakDuration = 5 * time.Minute
    MaxSessionDuration = 4 * time.Hour
)
```

### 3.2 Error Handling

- **Always check errors** returned by functions
- **Wrap errors** with context using `fmt.Errorf`
- **Use custom error types** for specific error conditions
- **Return errors early** to avoid deep nesting

```go
// Good
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }
    
    return &config, nil
}
```

### 3.3 Testing Standards

#### Test File Organization
- Create test files with `_test.go` suffix
- Place tests in the same package as the code being tested
- Use **table-driven tests** for multiple test cases

```go
func TestTimerStart(t *testing.T) {
    tests := []struct {
        name     string
        duration time.Duration
        wantErr  bool
    }{
        {"valid duration", 25 * time.Minute, false},
        {"zero duration", 0, true},
        {"negative duration", -1 * time.Minute, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            timer := NewTimer()
            err := timer.Start(tt.duration)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Timer.Start() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Test Coverage
- Aim for **80%+ test coverage** for critical packages
- Focus on **business logic** and **error conditions**
- Use **integration tests** for complex workflows

### 3.4 Documentation Standards

#### Package Documentation
- Every package should have a **package comment**
- Explain the **purpose** and **usage** of the package
- Include **examples** for complex packages

```go
// Package timer provides core timer functionality for the Pomodux application.
// It includes session management, progress tracking, and event emission.
//
// Example usage:
//
//	timer := timer.New()
//	err := timer.Start(25 * time.Minute)
//	if err != nil {
//	    log.Fatal(err)
//	}
package timer
```

#### Function Documentation
- Document **exported functions** with clear descriptions
- Include **parameter descriptions** and **return value explanations**
- Provide **usage examples** for complex functions

```go
// Start begins a new timer session with the specified duration.
// The timer will emit events when started, paused, resumed, and completed.
// Returns an error if the timer is already running or if the duration is invalid.
func (t *Timer) Start(duration time.Duration) error {
    // Implementation...
}
```

## 4.0 Development Tools and Workflow

### 4.1 Required Tools

#### Core Tools
- **Go 1.21+**: Minimum Go version for the project
- **gofmt**: Code formatting (built into Go)
- **golangci-lint**: Comprehensive linting
- **go test**: Testing framework
- **go mod**: Dependency management

#### Development Tools
- **Make**: Build automation
- **git**: Version control
- **pre-commit hooks**: Code quality enforcement

### 4.2 Code Quality Enforcement

#### Linting Configuration
Create `.golangci.yml` in the project root:

```yaml
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

linters-settings:
  golint:
    min-confidence: 0.8

run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor/
    - third_party/
```

#### Pre-commit Hooks
Create `.githooks/pre-commit`:

```bash
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
```

### 4.3 Build and Development Workflow

#### Makefile Targets
Create `Makefile` with common development tasks:

```makefile
.PHONY: build test lint clean install

# Build the application
build:
	go build -o bin/pomodux cmd/pomodux/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
install:
	go mod download
	go mod tidy

# Run the application
run:
	go run cmd/pomodux/main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/pomodux-linux-amd64 cmd/pomodux/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/pomodux-darwin-amd64 cmd/pomodux/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/pomodux-windows-amd64.exe cmd/pomodux/main.go
```

## 5.0 Dependency Management

### 5.1 Module Configuration

#### go.mod Structure
```go
module github.com/yourusername/pomodux

go 1.21

require (
    github.com/spf13/cobra v1.7.0
    github.com/charmbracelet/bubbletea v0.24.0
    github.com/spf13/viper v1.16.0
    github.com/yuin/gopher-lua v1.1.0
)

require (
    // Transitive dependencies will be added automatically
)
```

#### Dependency Guidelines
- **Pin major versions** for stability
- **Use go mod tidy** to clean up dependencies
- **Review dependencies** regularly for security updates
- **Prefer well-maintained libraries** with active development

### 5.2 Vendor Management
- Use `go mod vendor` for reproducible builds
- Include vendor directory in version control
- Update vendor dependencies regularly

## 6.0 Performance and Best Practices

### 6.1 Performance Guidelines
- **Profile critical paths** using `go test -bench`
- **Use sync.Pool** for frequently allocated objects
- **Avoid unnecessary allocations** in hot paths
- **Use buffered channels** when appropriate

### 6.2 Concurrency Patterns
- **Use goroutines** for I/O-bound operations
- **Use channels** for communication between goroutines
- **Implement proper cancellation** with context
- **Avoid goroutine leaks** by ensuring proper cleanup

### 6.3 Memory Management
- **Use sync.Pool** for frequently allocated objects
- **Avoid memory leaks** in long-running applications
- **Profile memory usage** in development
- **Use pprof** for performance analysis

## 7.0 Security Considerations

### 7.1 Input Validation
- **Validate all user inputs** before processing
- **Use secure defaults** for configuration
- **Sanitize file paths** and URLs
- **Implement proper error handling** without information disclosure

### 7.2 Configuration Security
- **Use environment variables** for sensitive configuration
- **Validate configuration** at startup
- **Use secure file permissions** for config files
- **Implement configuration encryption** if needed

## 8.0 Cross-Platform Considerations

### 8.1 Platform-Specific Code
- **Use build tags** for platform-specific implementations
- **Test on all target platforms** regularly
- **Handle platform differences** gracefully
- **Use Go's cross-compilation** for builds

### 8.2 File System Operations
- **Use path/filepath** for cross-platform path handling
- **Respect XDG standards** on Linux
- **Handle Windows path separators** correctly
- **Use os.UserConfigDir()** for configuration paths

## 9.0 Monitoring and Observability

### 9.1 Logging Standards
- **Use structured logging** with consistent fields
- **Include context** in log messages
- **Use appropriate log levels** (debug, info, warn, error)
- **Avoid logging sensitive information**

### 9.2 Metrics and Monitoring
- **Expose metrics** for key operations
- **Use Prometheus** for metrics collection
- **Monitor application health** with health checks
- **Track performance metrics** in production

## 10.0 Documentation and Examples

### 10.1 Code Examples
- **Include examples** in package documentation
- **Create example programs** in `examples/` directory
- **Document complex workflows** with step-by-step guides
- **Provide usage examples** for all public APIs

### 10.2 API Documentation
- **Generate API documentation** using `go doc`
- **Include OpenAPI specs** if applicable
- **Document breaking changes** clearly
- **Maintain changelog** for releases

## 11.0 Test-Driven Development (TDD) Examples

### 11.1 TDD Process Examples

#### 11.1.1 Basic TDD Cycle

**For every code change, follow the Red-Green-Refactor cycle:**

1. **Write failing test first (Red)**:
   ```go
   func TestNewFeature_Behavior(t *testing.T) {
       // Test the desired behavior
       result := NewFeature()
       assert.Equal(t, expected, result)
   }
   ```

2. **Write minimal code to pass (Green)**:
   ```go
   func NewFeature() string {
       return "expected"
   }
   ```

3. **Refactor while keeping tests passing**:
   ```go
   func NewFeature() string {
       return calculateExpectedValue()
   }
   ```

### 11.2 Component-Specific Examples

#### 11.2.1 Timer Engine Examples

**Timer Creation and State**

**Step 1: Write Test (Red)**
```go
func TestTimerEngine_NewTimer(t *testing.T) {
    engine := NewTimerEngine()
    assert.Equal(t, StatusIdle, engine.GetStatus())
    assert.Equal(t, 0.0, engine.GetProgress())
}
```

**Step 2: Write Code (Green)**
```go
type TimerEngine struct {
    status   TimerStatus
    progress float64
}

func NewTimerEngine() *TimerEngine {
    return &TimerEngine{
        status:   StatusIdle,
        progress: 0.0,
    }
}

func (t *TimerEngine) GetStatus() TimerStatus {
    return t.status
}

func (t *TimerEngine) GetProgress() float64 {
    return t.progress
}
```

**Step 3: Refactor**
```go
type TimerEngine struct {
    state *TimerState
}

type TimerState struct {
    status   TimerStatus
    progress float64
}

func NewTimerEngine() *TimerEngine {
    return &TimerEngine{
        state: &TimerState{
            status:   StatusIdle,
            progress: 0.0,
        },
    }
}

func (t *TimerEngine) GetStatus() TimerStatus {
    return t.state.status
}

func (t *TimerEngine) GetProgress() float64 {
    return t.state.progress
}
```

**Timer Start Functionality**

**Step 1: Write Test (Red)**
```go
func TestTimerEngine_Start(t *testing.T) {
    engine := NewTimerEngine()
    duration := 25 * time.Minute
    
    err := engine.Start(duration)
    assert.NoError(t, err)
    assert.Equal(t, StatusRunning, engine.GetStatus())
}
```

**Step 2: Write Code (Green)**
```go
func (t *TimerEngine) Start(duration time.Duration) error {
    if duration <= 0 {
        return errors.New("duration must be positive")
    }
    
    t.state.status = StatusRunning
    return nil
}
```

**Step 3: Refactor**
```go
func (t *TimerEngine) Start(duration time.Duration) error {
    if err := t.validateDuration(duration); err != nil {
        return err
    }
    
    t.state.status = StatusRunning
    return nil
}

func (t *TimerEngine) validateDuration(duration time.Duration) error {
    if duration <= 0 {
        return errors.New("duration must be positive")
    }
    return nil
}
```

#### 11.2.2 CLI Command Examples

**Start Command with Duration**

**Step 1: Write Test (Red)**
```go
func TestStartCommand_WithDuration(t *testing.T) {
    cmd := NewStartCommand()
    cmd.SetArgs([]string{"25m"})
    
    err := cmd.Execute()
    assert.NoError(t, err)
    // Verify timer was started with correct duration
}
```

**Step 2: Write Code (Green)**
```go
var startCmd = &cobra.Command{
    Use:   "start [duration]",
    Short: "Start a work timer",
    RunE:  runStartTimer,
}

func runStartTimer(cmd *cobra.Command, args []string) error {
    duration := "25m" // default
    if len(args) > 0 {
        duration = args[0]
    }
    
    // Parse duration and start timer
    return nil
}
```

**Step 3: Refactor**
```go
func runStartTimer(cmd *cobra.Command, args []string) error {
    duration, err := parseDuration(args)
    if err != nil {
        return err
    }
    
    return startTimer(duration)
}

func parseDuration(args []string) (time.Duration, error) {
    if len(args) == 0 {
        return 25 * time.Minute, nil
    }
    
    return time.ParseDuration(args[0])
}

func startTimer(duration time.Duration) error {
    engine := NewTimerEngine()
    return engine.Start(duration)
}
```

#### 11.2.3 Configuration Management Examples

**Configuration Loading**

**Step 1: Write Test (Red)**
```go
func TestConfigManager_LoadConfig(t *testing.T) {
    manager := NewConfigManager()
    
    // Create temporary config file
    configData := `{
        "work_duration": "25m",
        "break_duration": "5m",
        "long_break_duration": "15m"
    }`
    
    tmpFile := createTempConfigFile(t, configData)
    defer os.Remove(tmpFile.Name())
    
    config, err := manager.LoadConfig(tmpFile.Name())
    
    assert.NoError(t, err)
    assert.Equal(t, 25*time.Minute, config.WorkDuration)
    assert.Equal(t, 5*time.Minute, config.BreakDuration)
    assert.Equal(t, 15*time.Minute, config.LongBreakDuration)
}
```

**Step 2: Write Code (Green)**
```go
type ConfigManager struct{}

type Config struct {
    WorkDuration      time.Duration `json:"work_duration"`
    BreakDuration     time.Duration `json:"break_duration"`
    LongBreakDuration time.Duration `json:"long_break_duration"`
}

func (cm *ConfigManager) LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }
    
    return &config, nil
}
```

**Step 3: Refactor**
```go
func (cm *ConfigManager) LoadConfig(path string) (*Config, error) {
    data, err := cm.readConfigFile(path)
    if err != nil {
        return nil, err
    }
    
    return cm.parseConfig(data)
}

func (cm *ConfigManager) readConfigFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    return data, nil
}

func (cm *ConfigManager) parseConfig(data []byte) (*Config, error) {
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }
    return &config, nil
}
```

### 11.3 Test Organization Examples

#### 11.3.1 Table-Driven Tests

```go
func TestTimerEngine_Start(t *testing.T) {
    tests := []struct {
        name     string
        duration time.Duration
        wantErr  bool
        wantStatus TimerStatus
    }{
        {
            name:        "valid work duration",
            duration:    25 * time.Minute,
            wantErr:     false,
            wantStatus:  StatusRunning,
        },
        {
            name:        "zero duration",
            duration:    0,
            wantErr:     true,
            wantStatus:  StatusIdle,
        },
        {
            name:        "negative duration",
            duration:    -1 * time.Minute,
            wantErr:     true,
            wantStatus:  StatusIdle,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := NewTimerEngine()
            
            err := engine.Start(tt.duration)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("TimerEngine.Start() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if engine.GetStatus() != tt.wantStatus {
                t.Errorf("TimerEngine.GetStatus() = %v, want %v", engine.GetStatus(), tt.wantStatus)
            }
        })
    }
}
```

#### 11.3.2 Test Helpers

```go
// test/helpers/config.go
package test

import (
    "os"
    "testing"
)

// CreateTempConfigFile creates a temporary configuration file for testing
func CreateTempConfigFile(t *testing.T, content string) *os.File {
    t.Helper()
    
    file, err := os.CreateTemp("", "pomodux_config_*.json")
    require.NoError(t, err)
    
    _, err = file.WriteString(content)
    require.NoError(t, err)
    
    return file
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration) {
    t.Helper()
    
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        if condition() {
            return
        }
        time.Sleep(10 * time.Millisecond)
    }
    
    t.Fatal("condition not met within timeout")
}
```

#### 11.3.3 Mock Objects

```go
// internal/timer/mocks.go
package timer

// MockTimerEngine provides a mock implementation for testing
type MockTimerEngine struct {
    status   TimerStatus
    progress float64
    err      error
}

func (m *MockTimerEngine) Start(duration time.Duration) error {
    if m.err != nil {
        return m.err
    }
    m.status = StatusRunning
    return nil
}

func (m *MockTimerEngine) GetStatus() TimerStatus {
    return m.status
}

func (m *MockTimerEngine) GetProgress() float64 {
    return m.progress
}
```

### 11.4 Integration Test Examples

#### 11.4.1 CLI Integration Test

```go
func TestCLI_StartCommand_Integration(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    // Create command and run
    cmd := NewStartCommand()
    cmd.SetArgs([]string{"25m"})
    err := cmd.Execute()
    
    // Restore stdout and read output
    w.Close()
    os.Stdout = oldStdout
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    assert.NoError(t, err)
    assert.Contains(t, output, "Timer started for 25m0s")
}
```

### 11.5 Performance Test Examples

#### 11.5.1 Benchmark Tests

```go
func BenchmarkTimerEngine_Start(b *testing.B) {
    engine := NewTimerEngine()
    duration := 25 * time.Minute
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.Start(duration)
        engine.Stop()
    }
}
```

### 11.6 Test Coverage Examples

#### 11.6.1 Coverage Requirements

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# Check coverage by package
go test -cover ./internal/timer
go test -cover ./internal/cli
go test -cover ./internal/config

# View coverage in browser
go tool cover -html=coverage.out -o coverage.html
```

#### 11.6.2 Coverage Targets

- **Overall Coverage**: Minimum 80% for all packages
- **Critical Paths**: Minimum 95% for timer engine and core functionality
- **Public APIs**: 100% coverage for all exported functions
- **Integration Tests**: Test component interactions and workflows

---

## References

### Project Documentation
- **[Cursor Rules](.cursor/rules/go-development.mdc)** - AI assistant guidance and project-specific patterns
- **[Development Setup Guide](docs/development-setup.md)** - Environment setup and development workflow
- **[Implementation Roadmap](docs/implementation-roadmap.md)** - Development phases and priorities
- **[Architecture Decision Records](docs/adr/)** - Project architectural decisions and rationale
- **[Technical Specifications](docs/technical_specifications.md)** - System architecture and component design

### External Resources
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Best Practices](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TableDrivenTests)

### How These Files Work Together
- **This file (Go Standards)**: Provides human developers with detailed implementation standards and examples
- **[Cursor Rules](.cursor/rules/go-development.mdc)**: Provides AI assistants with project-specific patterns and high-level guidance
- **Together**: Ensure consistent code quality and project architecture across both AI and human contributions 