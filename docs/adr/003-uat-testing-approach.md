---
status: approved
type: technical
---

# ADR 003: UAT Testing Approach for Command Line Applications

## 1. Context / Background

### 1.1 Current State
Pomodux is a command-line application that requires comprehensive User Acceptance Testing (UAT) to validate functionality, user experience, and real-world usage scenarios. The current UAT script is a basic bash script with echo statements, but we need to establish a standardized approach for CLI testing that ensures reliability, maintainability, and comprehensive coverage.

### 1.2 Requirements
- Comprehensive testing of all CLI commands and features
- Validation of user workflows and real-world usage scenarios
- Clear documentation of expected behaviors and outputs
- Reproducible test execution across different environments
- Integration with the 4-gate approval process
- Support for both manual and automated testing approaches

## 2. Decision

**Selected Solution:** Hybrid UAT Testing Approach with Structured Bash Scripts and Test Frameworks

### 2.1 Rationale
After researching industry best practices for CLI testing, we've identified that a hybrid approach combining structured bash scripts with dedicated testing frameworks provides the best balance of simplicity, reliability, and comprehensive coverage. This approach leverages the strengths of both manual testing (for user experience validation) and automated testing (for regression testing and CI/CD integration).

## 3. Solutions Considered

### 3.1 Option A: Pure Bash Scripts with Manual Validation
- **Pros:**
  - Simple to implement and understand
  - No external dependencies
  - Easy to modify and extend
  - Clear visibility of test steps and expected outputs
- **Cons:**
  - Manual validation required for each test
  - No automated pass/fail detection
  - Difficult to integrate with CI/CD
  - Prone to human error in result interpretation
  - No structured reporting or metrics

### 3.2 Option B: Dedicated Testing Frameworks (BATS, shunit2)
- **Pros:**
  - Automated test execution and validation
  - Built-in assertions and test helpers
  - Integration with CI/CD pipelines
  - Structured reporting and metrics
  - Professional testing practices
- **Cons:**
  - Additional dependencies and complexity
  - Steeper learning curve
  - May be overkill for simple CLI applications
  - Less intuitive for non-technical stakeholders

### 3.3 Option C: Hybrid Approach (Selected)
- **Pros:**
  - Combines benefits of both approaches
  - Structured bash scripts for manual UAT
  - Automated test framework for regression testing
  - Clear separation of concerns
  - Flexible and extensible
  - Supports both manual validation and automated execution
- **Cons:**
  - Requires maintaining two test suites
  - Slightly more complex setup
  - Need to ensure consistency between approaches

## 4. Consequences

### 4.1 Positive
- **Comprehensive Coverage**: Both manual and automated testing ensure thorough validation
- **User Experience Focus**: Manual UAT validates real-world usage scenarios
- **Regression Protection**: Automated tests prevent regressions in CI/CD
- **Clear Documentation**: Structured scripts serve as living documentation
- **Flexible Execution**: Tests can be run manually or automatically
- **Professional Standards**: Follows industry best practices for CLI testing

### 4.2 Negative
- **Maintenance Overhead**: Two test suites require maintenance
- **Setup Complexity**: Initial setup requires more configuration
- **Learning Curve**: Team needs to understand both approaches
- **Consistency Challenges**: Ensuring both test suites stay in sync

### 4.3 Risks
- **Test Drift**: Manual and automated tests may diverge over time
- **Maintenance Burden**: Increased complexity may lead to neglected tests
- **Tool Dependency**: Additional testing frameworks add dependencies

**Mitigation Strategies**:
- Regular review and synchronization of test suites
- Clear documentation and training for team members
- Automated validation of test consistency
- Gradual migration path from manual to automated testing

## 5. Component Information

### 5.1 Repository Links
- **BATS (Bash Automated Testing System)**: https://github.com/bats-core/bats-core
- **shunit2**: https://github.com/kward/shunit2
- **ShellSpec**: https://github.com/shellspec/shellspec
- **Documentation**: https://github.com/bats-core/bats-core/wiki

### 5.2 Maintenance Status
- **Last Updated**: 2025-01-27
- **Active Development**: Yes (ongoing)
- **Community Support**: Strong community support for BATS and shunit2
- **Version Compatibility**: Compatible with current Go and shell environments

### 5.3 Integration Verification
- **Compatibility Tested**: Yes, with current Pomodux CLI structure
- **Existing Component Impact**: Minimal impact on existing codebase
- **Migration Path**: Gradual migration from current manual scripts

## 6. Implementation Strategy

### 6.1 Phase 1: Enhanced Manual UAT Scripts
**Duration**: Immediate implementation
**Deliverables**:
- Structured bash scripts with clear test organization
- Comprehensive test scenarios covering all features
- Clear documentation of expected outputs
- Integration with release management process

**Implementation Details**:
```bash
#!/bin/bash
# UAT Test Suite for Pomodux
# Version: 1.0
# Purpose: Comprehensive user acceptance testing

set -euo pipefail  # Strict error handling

# Test configuration
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_BINARY="${TEST_DIR}/../pomodux"
CONFIG_DIR="${HOME}/.config/pomodux"
TEST_RESULTS="${TEST_DIR}/uat-results.log"

# Test utilities
log_test() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$TEST_RESULTS"
}

assert_exit_code() {
    local expected=$1
    local actual=$2
    local test_name=$3
    
    if [ "$actual" -eq "$expected" ]; then
        log_test "✅ PASS: $test_name"
    else
        log_test "❌ FAIL: $test_name (expected $expected, got $actual)"
        return 1
    fi
}

assert_output_contains() {
    local expected=$1
    local output=$2
    local test_name=$3
    
    if echo "$output" | grep -q "$expected"; then
        log_test "✅ PASS: $test_name"
    else
        log_test "❌ FAIL: $test_name (expected '$expected' in output)"
        return 1
    fi
}

# Test setup and cleanup
setup_test_environment() {
    log_test "Setting up test environment..."
    mkdir -p "$CONFIG_DIR"
    # Backup existing config
    if [ -f "${CONFIG_DIR}/config.json" ]; then
        cp "${CONFIG_DIR}/config.json" "${CONFIG_DIR}/config.json.backup"
    fi
}

cleanup_test_environment() {
    log_test "Cleaning up test environment..."
    # Restore original config
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
    fi
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
}

# Test scenarios
test_basic_functionality() {
    log_test "=== Testing Basic Functionality ==="
    
    # Test version command
    output=$("$APP_BINARY" version 2>&1)
    assert_exit_code 0 $? "Version command should succeed"
    assert_output_contains "pomodux" "$output" "Version output should contain app name"
    
    # Test help command
    output=$("$APP_BINARY" --help 2>&1)
    assert_exit_code 0 $? "Help command should succeed"
    assert_output_contains "start" "$output" "Help should show start command"
}

test_timer_operations() {
    log_test "=== Testing Timer Operations ==="
    
    # Test start command
    output=$("$APP_BINARY" start 1m 2>&1)
    assert_exit_code 0 $? "Start command should succeed"
    
    # Test status command
    output=$("$APP_BINARY" status 2>&1)
    assert_exit_code 0 $? "Status command should succeed"
    assert_output_contains "running" "$output" "Status should show running timer"
    
    # Test stop command
    output=$("$APP_BINARY" stop 2>&1)
    assert_exit_code 0 $? "Stop command should succeed"
}

# Main test execution
main() {
    log_test "Starting Pomodux UAT Test Suite"
    
    setup_test_environment
    
    # Run test scenarios
    test_basic_functionality
    test_timer_operations
    # Add more test scenarios...
    
    cleanup_test_environment
    
    log_test "UAT Test Suite completed"
    log_test "Results saved to: $TEST_RESULTS"
}

# Execute main function
main "$@"
```

### 6.2 Phase 2: Automated Test Framework Integration
**Duration**: Future implementation (Release 0.3.0+)
**Deliverables**:
- BATS test suite for automated testing
- CI/CD integration
- Automated regression testing
- Test reporting and metrics

**Implementation Details**:
```bash
#!/usr/bin/env bats

# BATS test suite for Pomodux
# Requires: bats-core (https://github.com/bats-core/bats-core)

setup() {
    # Setup test environment
    export APP_BINARY="${BATS_TEST_DIRNAME}/../pomodux"
    export CONFIG_DIR="${HOME}/.config/pomodux"
    
    # Create test config
    mkdir -p "$CONFIG_DIR"
    cat > "${CONFIG_DIR}/config.json" << EOF
{
    "default_work_duration": "1m",
    "default_break_duration": "30s",
    "default_long_break_duration": "2m"
}
EOF
}

teardown() {
    # Cleanup test environment
    rm -f "${CONFIG_DIR}/session_history.json"
}

@test "version command should succeed" {
    run "$APP_BINARY" version
    [ "$status" -eq 0 ]
    [[ "$output" =~ "pomodux" ]]
}

@test "help command should show available commands" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "start" ]]
    [[ "$output" =~ "stop" ]]
    [[ "$output" =~ "status" ]]
}

@test "start command should create running timer" {
    run "$APP_BINARY" start 1m
    [ "$status" -eq 0 ]
    
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "running" ]]
}

@test "stop command should stop running timer" {
    # Start a timer
    "$APP_BINARY" start 1m > /dev/null 2>&1
    
    # Stop the timer
    run "$APP_BINARY" stop
    [ "$status" -eq 0 ]
    
    # Verify timer is stopped
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "idle" ]]
}
```

### 6.3 Phase 3: Integration and Optimization
**Duration**: Ongoing improvement
**Deliverables**:
- Unified test reporting
- Performance benchmarking
- Cross-platform testing
- Advanced test scenarios

## 7. Test Organization Structure

### 7.1 Directory Structure
```
tests/
├── uat/
│   ├── manual/
│   │   ├── uat-script.sh          # Main UAT script
│   │   ├── test-scenarios.md      # Test scenario documentation
│   │   └── expected-outputs.md    # Expected output documentation
│   └── automated/
│       ├── bats/
│       │   ├── basic-functionality.bats
│       │   ├── timer-operations.bats
│       │   └── configuration.bats
│       └── integration/
│           └── end-to-end.bats
├── fixtures/
│   ├── configs/                   # Test configuration files
│   └── data/                      # Test data files
└── reports/
    ├── uat-results.log           # Manual test results
    └── bats-results.xml          # Automated test results
```

### 7.2 Test Categories

#### Manual UAT Tests
- **User Experience Tests**: Validate intuitive usage and clear feedback
- **Workflow Tests**: Test complete user workflows and scenarios
- **Error Handling Tests**: Validate error messages and recovery
- **Configuration Tests**: Test configuration management features
- **Performance Tests**: Validate response times and resource usage

#### Automated Tests
- **Unit Tests**: Test individual command functionality
- **Integration Tests**: Test command interactions and workflows
- **Regression Tests**: Ensure no regressions in existing functionality
- **Cross-Platform Tests**: Validate behavior across different platforms

## 8. Integration with Release Process

### 8.1 Gate 3 Integration
- **Manual UAT**: Required for Gate 3 approval
- **Automated Tests**: Run as part of CI/CD pipeline
- **Test Results**: Documented in release documentation
- **Issue Tracking**: Test failures tracked as GitHub issues

### 8.2 Test Execution Workflow
1. **Development Phase**: Automated tests run on every commit
2. **UAT Phase**: Manual UAT scripts executed by stakeholders
3. **Release Phase**: Both test suites run for final validation
4. **Post-Release**: Automated tests run for regression protection

## 9. Quality Standards

### 9.1 Test Coverage Requirements
- **Feature Coverage**: 100% of user-facing features
- **Command Coverage**: 100% of CLI commands
- **Error Scenario Coverage**: All documented error conditions
- **Workflow Coverage**: All documented user workflows

### 9.2 Test Quality Standards
- **Clarity**: Tests must be clear and understandable
- **Reliability**: Tests must be deterministic and repeatable
- **Maintainability**: Tests must be easy to update and extend
- **Documentation**: All tests must be well-documented

### 9.3 Performance Standards
- **Execution Time**: Manual UAT should complete within 30 minutes
- **Resource Usage**: Tests should not consume excessive resources
- **Reliability**: Tests should have 95%+ pass rate in stable environments

## 10. Future Considerations

### 10.1 Scalability
- **Test Parallelization**: Parallel execution of independent tests
- **Test Distribution**: Distributed testing across multiple environments
- **Performance Testing**: Automated performance benchmarking

### 10.2 Advanced Features
- **Visual Testing**: Screenshot comparison for TUI components
- **Accessibility Testing**: Automated accessibility validation
- **Security Testing**: Automated security vulnerability scanning

### 10.3 Tool Evolution
- **Framework Migration**: Potential migration to more advanced testing frameworks
- **CI/CD Integration**: Enhanced integration with modern CI/CD platforms
- **Reporting Enhancement**: Advanced test reporting and analytics

---

**References**:
- [BATS Core Documentation](https://github.com/bats-core/bats-core/wiki)
- [Shell Testing Best Practices](https://github.com/kward/shunit2)
- [CLI Testing Patterns](https://clig.dev/#testing)
- [Go Testing Standards](https://golang.org/doc/code.html#Testing) 