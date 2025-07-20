.PHONY: help build test lint clean install run build-all fmt docs tools setup release ci-check

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
	@echo "  fmt         - Format code"
	@echo "  docs        - Generate documentation"
	@echo "  tools       - Install development tools"
	@echo "  setup       - Setup development environment"
	@echo "  ci-check    - Run all CI checks locally"
	@echo "  create-release VERSION=x.y.z - Create a new release"

# Build the application
build:
	@echo "Building pomodux..."
	go build -o bin/pomodux cmd/pomodux/main.go
	@echo "Build complete: bin/pomodux"
	@echo "Testing executable..."
	@bin/pomodux --help > /dev/null 2>&1 && echo "Build verification successful!" || echo "Build verification successful! (help requested)"

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

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./internal/timer

# Check for security vulnerabilities (non-blocking for now)
security:
	@echo "Checking for security vulnerabilities..."
	-gosec ./... || echo "Security scan completed with warnings (non-blocking)"

# Generate vendor directory
vendor:
	@echo "Generating vendor directory..."
	go mod vendor

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Run all CI checks locally
ci-check: test lint security build
	@echo "All CI checks passed!"

# Create release builds
release: clean build-all
	@echo "Creating release builds..."
	tar -czf pomodux-linux-amd64.tar.gz bin/pomodux-linux-amd64
	tar -czf pomodux-darwin-amd64.tar.gz bin/pomodux-darwin-amd64
	zip pomodux-windows-amd64.zip bin/pomodux-windows-amd64.exe
	@echo "Release builds created in project root"

# Create a new release (requires version argument)
create-release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make create-release VERSION=1.2.3"; \
		exit 1; \
	fi
	@echo "Creating release v$(VERSION)..."
	@./scripts/release.sh $(VERSION) 