# Pomodux UAT Test Suite

This directory contains comprehensive User Acceptance Testing (UAT) scripts and automated tests for the Pomodux application, implementing the testing strategy defined in [ADR 003: UAT Testing Approach](docs/adr/003-uat-testing-approach.md).

## Overview

The UAT test suite provides two complementary testing approaches:

1. **Manual UAT Scripts** - For user experience validation and manual testing
2. **Automated BATS Tests** - For regression testing and CI/CD integration

## Directory Structure

```
tests/uat/
├── manual/
│   └── uat-script.sh              # Enhanced manual UAT script
├── automated/
│   ├── bats/
│   │   ├── basic-functionality.bats
│   │   ├── timer-operations.bats
│   │   ├── configuration.bats
│   │   ├── session-types.bats
│   │   └── session-history.bats
│   ├── integration/
│   │   └── end-to-end.bats
│   └── run-tests.sh               # Test runner script
├── fixtures/
│   └── configs/                   # Test configuration files
└── reports/                       # Test reports and results
```

## Manual UAT Testing

### Enhanced UAT Script

The enhanced manual UAT script provides comprehensive testing with automated validation:

```bash
# Run the enhanced UAT script
./tests/uat/manual/uat-script.sh
```

**Features:**
- ✅ Automated test execution and validation
- ✅ Comprehensive error handling and logging
- ✅ Test result reporting and metrics
- ✅ Proper test environment setup and cleanup
- ✅ Structured test organization
- ✅ Integration with release management process

**Test Categories:**
- Basic functionality testing
- Configuration management
- Timer operations (start, pause, resume, stop)
- Session types (work, break, long break)
- Session history
- Error handling
- Completion commands
- Force flag functionality

### Simple UAT Script

For basic manual testing with clear explanations:

```bash
# Use the simple UAT script from docs/uat-script.md
# Copy the bash section and run it manually
```

## Automated Testing with BATS

### Prerequisites

Install BATS (Bash Automated Testing System):

```bash
# On macOS with Homebrew
brew install bats-core

# On Ubuntu/Debian
sudo apt-get install bats

# On Arch Linux
sudo pacman -S bats

# Or install from source
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

### Running Automated Tests

#### Run All Tests
```bash
# Run all automated tests with reporting
./tests/uat/automated/run-tests.sh
```

#### Run Individual Test Suites
```bash
# Basic functionality tests
bats tests/uat/automated/bats/basic-functionality.bats

# Timer operations tests
bats tests/uat/automated/bats/timer-operations.bats

# Configuration management tests
bats tests/uat/automated/bats/configuration.bats

# Session types tests
bats tests/uat/automated/bats/session-types.bats

# Session history tests
bats tests/uat/automated/bats/session-history.bats

# Integration tests
bats tests/uat/automated/integration/end-to-end.bats
```

### Test Coverage

The automated tests cover:

#### Unit Tests
- **Basic Functionality**: Version, help, completion commands
- **Timer Operations**: Start, stop, pause, resume, status
- **Configuration Management**: Show, set, validation
- **Session Types**: Work, break, long break sessions
- **Session History**: Recording, display, persistence

#### Integration Tests
- **End-to-End Workflows**: Complete user scenarios
- **Real-time Timer Completion**: Auto-completion behavior
- **State Persistence**: Cross-process state management
- **Error Handling**: Edge cases and error conditions

## Test Reports

### Automated Test Reports

Test reports are generated in the `tests/reports/` directory:

- **TAP Format**: Individual test suite results
- **Summary Report**: Comprehensive test summary with metrics
- **Timestamped**: Each run creates timestamped reports

### Report Structure

```
tests/reports/
├── test-summary-20250127_143022.md    # Summary report
├── basic-functionality-20250127_143022.tap
├── timer-operations-20250127_143022.tap
├── configuration-20250127_143022.tap
├── session-types-20250127_143022.tap
├── session-history-20250127_143022.tap
└── end-to-end-20250127_143022.tap
```

## Test Configuration

### Test Environment Setup

All tests automatically:
- Build the application
- Set up test configuration
- Backup existing configuration
- Clean up after completion

### Test Configuration Files

Test configurations use short durations for quick execution:

```json
{
    "default_work_duration": "1m",
    "default_break_duration": "30s",
    "default_long_break_duration": "2m"
}
```

### Test Isolation

Each test runs in isolation:
- Configuration is backed up and restored
- Session history is cleaned between tests
- No interference between test runs

## Integration with Release Process

### Gate 3 Integration

The UAT test suite integrates with the 4-gate approval process:

1. **Manual UAT**: Required for Gate 3 approval
2. **Automated Tests**: Run as part of CI/CD pipeline
3. **Test Results**: Documented in release documentation
4. **Issue Tracking**: Test failures tracked as GitHub issues

### Release Testing Workflow

1. **Development Phase**: Automated tests run on every commit
2. **UAT Phase**: Manual UAT scripts executed by stakeholders
3. **Release Phase**: Both test suites run for final validation
4. **Post-Release**: Automated tests run for regression protection

## Quality Standards

### Test Coverage Requirements

- **Feature Coverage**: 100% of user-facing features
- **Command Coverage**: 100% of CLI commands
- **Error Scenario Coverage**: All documented error conditions
- **Workflow Coverage**: All documented user workflows

### Test Quality Standards

- **Clarity**: Tests must be clear and understandable
- **Reliability**: Tests must be deterministic and repeatable
- **Maintainability**: Tests must be easy to update and extend
- **Documentation**: All tests must be well-documented

### Performance Standards

- **Execution Time**: Manual UAT should complete within 30 minutes
- **Resource Usage**: Tests should not consume excessive resources
- **Reliability**: Tests should have 95%+ pass rate in stable environments

## Troubleshooting

### Common Issues

#### BATS Not Found
```bash
# Install BATS
brew install bats-core  # macOS
sudo apt-get install bats  # Ubuntu/Debian
```

#### Test Failures
1. Check if application is built: `make build`
2. Check configuration: `./pomodux config show`
3. Check for running timers: `./pomodux status`
4. Stop any running timers: `./pomodux stop`

#### Permission Issues
```bash
# Make scripts executable
chmod +x tests/uat/manual/uat-script.sh
chmod +x tests/uat/automated/run-tests.sh
```

### Debug Mode

Run tests with verbose output:

```bash
# Manual UAT with debug output
./tests/uat/manual/uat-script.sh 2>&1 | tee uat-debug.log

# BATS with verbose output
bats --verbose tests/uat/automated/bats/basic-functionality.bats
```

## Future Enhancements

### Planned Improvements

1. **Visual Testing**: Screenshot comparison for TUI components
2. **Performance Testing**: Automated performance benchmarking
3. **Cross-Platform Testing**: Multi-platform test execution
4. **Continuous Integration**: GitHub Actions integration
5. **Test Parallelization**: Parallel execution of independent tests

### Advanced Features

- **Test Distribution**: Distributed testing across multiple environments
- **Accessibility Testing**: Automated accessibility validation
- **Security Testing**: Automated security vulnerability scanning
- **Advanced Reporting**: Enhanced test analytics and metrics

## 🔗 Related Documentation

- **[Release Management](docs/release-management.md)** - Release process and approval gates
- **[Backlog](docs/backlog/)** - Planning and requirements (current & future work)
- **[Releases](docs/releases/)** - Historical release documentation

## Contributing

When adding new tests:

1. **Follow Naming Conventions**: Use descriptive test names
2. **Maintain Isolation**: Ensure tests don't interfere with each other
3. **Add Documentation**: Document test purpose and expected behavior
4. **Update Coverage**: Ensure new features are covered by tests
5. **Run All Tests**: Verify no regressions in existing tests

## Support

For issues with the test suite:

1. Check the troubleshooting section above
2. Review the test reports for specific failure details
3. Check the [ADR 003](docs/adr/003-uat-testing-approach.md) for testing strategy
4. Create a GitHub issue with detailed error information 