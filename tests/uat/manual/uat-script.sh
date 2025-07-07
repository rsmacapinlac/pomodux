#!/bin/bash
# UAT Test Suite for Pomodux
# Version: 1.0
# Purpose: Comprehensive user acceptance testing
# Based on ADR 003: UAT Testing Approach for Command Line Applications

set -euo pipefail  # Strict error handling

# Test configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_BINARY="${PROJECT_ROOT}/pomodux"
CONFIG_DIR="${HOME}/.config/pomodux"
TEST_RESULTS="${SCRIPT_DIR}/uat-results.log"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test utilities
log_test() {
    local level="$1"
    local message="$2"
    local timestamp="$(date '+%Y-%m-%d %H:%M:%S')"
    
    case "$level" in
        "INFO")
            echo -e "${BLUE}[${timestamp}] INFO: ${message}${NC}" | tee -a "$TEST_RESULTS"
            ;;
        "PASS")
            echo -e "${GREEN}[${timestamp}] ✅ PASS: ${message}${NC}" | tee -a "$TEST_RESULTS"
            PASS_COUNT=$((PASS_COUNT + 1))
            ;;
        "FAIL")
            echo -e "${RED}[${timestamp}] ❌ FAIL: ${message}${NC}" | tee -a "$TEST_RESULTS"
            FAIL_COUNT=$((FAIL_COUNT + 1))
            ;;
        "WARN")
            echo -e "${YELLOW}[${timestamp}] ⚠️  WARN: ${message}${NC}" | tee -a "$TEST_RESULTS"
            ;;
        *)
            echo -e "[${timestamp}] ${message}" | tee -a "$TEST_RESULTS"
            ;;
    esac
}

assert_exit_code() {
    local expected="$1"
    local actual="$2"
    local test_name="$3"
    TEST_COUNT=$((TEST_COUNT + 1))
    
    if [ "$actual" -eq "$expected" ]; then
        log_test "PASS" "$test_name"
    else
        log_test "FAIL" "$test_name (expected exit code $expected, got $actual)"
        return 1
    fi
}

assert_output_contains() {
    local expected="$1"
    local output="$2"
    local test_name="$3"
    TEST_COUNT=$((TEST_COUNT + 1))
    
    if echo "$output" | grep -q "$expected"; then
        log_test "PASS" "$test_name"
    else
        log_test "FAIL" "$test_name (expected '$expected' in output)"
        return 1
    fi
}

assert_output_not_contains() {
    local unexpected="$1"
    local output="$2"
    local test_name="$3"
    TEST_COUNT=$((TEST_COUNT + 1))
    
    if ! echo "$output" | grep -q "$unexpected"; then
        log_test "PASS" "$test_name"
    else
        log_test "FAIL" "$test_name (unexpected '$unexpected' in output)"
        return 1
    fi
}

assert_file_exists() {
    local file_path="$1"
    local test_name="$2"
    TEST_COUNT=$((TEST_COUNT + 1))
    
    if [ -f "$file_path" ]; then
        log_test "PASS" "$test_name"
    else
        log_test "FAIL" "$test_name (file '$file_path' does not exist)"
        return 1
    fi
}

# Test setup and cleanup
setup_test_environment() {
    log_test "INFO" "Setting up test environment..."
    
    # Ensure we're in the project directory
    cd "$PROJECT_ROOT"
    
    # Build the application
    log_test "INFO" "Building Pomodux application..."
    if ! make build > /dev/null 2>&1; then
        log_test "FAIL" "Failed to build application"
        exit 1
    fi
    
    # Create config directory
    mkdir -p "$CONFIG_DIR"
    
    # Backup existing config
    if [ -f "${CONFIG_DIR}/config.json" ]; then
        cp "${CONFIG_DIR}/config.json" "${CONFIG_DIR}/config.json.backup"
        log_test "INFO" "Backed up existing configuration"
    fi
    
    # Create test config
    cat > "${CONFIG_DIR}/config.json" << EOF
{
    "default_work_duration": "1m",
    "default_break_duration": "30s",
    "default_long_break_duration": "2m"
}
EOF
    log_test "INFO" "Created test configuration"
}

cleanup_test_environment() {
    log_test "INFO" "Cleaning up test environment..."
    
    # Stop any running timer
    if [ -f "$APP_BINARY" ]; then
        "$APP_BINARY" stop > /dev/null 2>&1 || true
    fi
    
    # Restore original config
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
        log_test "INFO" "Restored original configuration"
    else
        rm -f "${CONFIG_DIR}/config.json"
    fi
    
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
    
    log_test "INFO" "Cleanup completed"
}

print_test_summary() {
    echo ""
    echo "=========================================="
    echo "UAT Test Summary"
    echo "=========================================="
    echo "Total Tests: $TEST_COUNT"
    echo "Passed: $PASS_COUNT"
    echo "Failed: $FAIL_COUNT"
    echo "Success Rate: $(( (PASS_COUNT * 100) / TEST_COUNT ))%"
    echo ""
    echo "Results saved to: $TEST_RESULTS"
    echo "=========================================="
    
    if [ $FAIL_COUNT -gt 0 ]; then
        exit 1
    fi
}

# Test scenarios
test_basic_functionality() {
    log_test "INFO" "=== Testing Basic Functionality ==="
    
    # Test version command
    log_test "INFO" "Testing version command..."
    output=$("$APP_BINARY" version 2>&1)
    assert_exit_code 0 $? "Version command should succeed"
    assert_output_contains "pomodux" "$output" "Version output should contain app name"
    
    # Test help command
    log_test "INFO" "Testing help command..."
    output=$("$APP_BINARY" --help 2>&1)
    assert_exit_code 0 $? "Help command should succeed"
    assert_output_contains "start" "$output" "Help should show start command"
    assert_output_contains "stop" "$output" "Help should show stop command"
    assert_output_contains "status" "$output" "Help should show status command"
}

test_configuration_management() {
    log_test "INFO" "=== Testing Configuration Management ==="
    
    # Test config show command
    log_test "INFO" "Testing config show command..."
    output=$("$APP_BINARY" config show 2>&1)
    assert_exit_code 0 $? "Config show command should succeed"
    assert_output_contains "default_work_duration" "$output" "Config should show work duration"
    
    # Test config set command
    log_test "INFO" "Testing config set command..."
    output=$("$APP_BINARY" config set default_work_duration 2m 2>&1)
    assert_exit_code 0 $? "Config set command should succeed"
    
    # Verify config change
    output=$("$APP_BINARY" config show 2>&1)
    assert_output_contains "2m" "$output" "Config should show updated work duration"
    
    # Test invalid config value
    log_test "INFO" "Testing invalid config value..."
    output=$("$APP_BINARY" config set default_work_duration invalid 2>&1)
    assert_exit_code 1 $? "Config set with invalid value should fail"
    assert_output_contains "error" "$output" "Should show error for invalid config"
}

test_timer_operations() {
    log_test "INFO" "=== Testing Timer Operations ==="
    
    # Test start command
    log_test "INFO" "Testing start command..."
    output=$("$APP_BINARY" start 1m 2>&1)
    assert_exit_code 0 $? "Start command should succeed"
    
    # Test status command with running timer
    log_test "INFO" "Testing status command with running timer..."
    output=$("$APP_BINARY" status 2>&1)
    assert_exit_code 0 $? "Status command should succeed"
    assert_output_contains "running" "$output" "Status should show running timer"
    assert_output_contains "remaining" "$output" "Status should show remaining time"
    
    # Test pause command
    log_test "INFO" "Testing pause command..."
    output=$("$APP_BINARY" pause 2>&1)
    assert_exit_code 0 $? "Pause command should succeed"
    
    # Test status command with paused timer
    output=$("$APP_BINARY" status 2>&1)
    assert_output_contains "paused" "$output" "Status should show paused timer"
    
    # Test resume command
    log_test "INFO" "Testing resume command..."
    output=$("$APP_BINARY" resume 2>&1)
    assert_exit_code 0 $? "Resume command should succeed"
    
    # Test status command with resumed timer
    output=$("$APP_BINARY" status 2>&1)
    assert_output_contains "running" "$output" "Status should show running timer after resume"
    
    # Test stop command
    log_test "INFO" "Testing stop command..."
    output=$("$APP_BINARY" stop 2>&1)
    assert_exit_code 0 $? "Stop command should succeed"
    
    # Test status command with stopped timer
    output=$("$APP_BINARY" status 2>&1)
    assert_output_contains "idle" "$output" "Status should show idle after stop"
}

test_session_types() {
    log_test "INFO" "=== Testing Session Types ==="
    
    # Test break command
    log_test "INFO" "Testing break command..."
    output=$("$APP_BINARY" break 2>&1)
    assert_exit_code 0 $? "Break command should succeed"
    
    output=$("$APP_BINARY" status 2>&1)
    assert_output_contains "break" "$output" "Status should show break session"
    
    "$APP_BINARY" stop > /dev/null 2>&1
    
    # Test long-break command
    log_test "INFO" "Testing long-break command..."
    output=$("$APP_BINARY" long-break 2>&1)
    assert_exit_code 0 $? "Long-break command should succeed"
    
    output=$("$APP_BINARY" status 2>&1)
    assert_output_contains "long break" "$output" "Status should show long break session"
    
    "$APP_BINARY" stop > /dev/null 2>&1
}

test_session_history() {
    log_test "INFO" "=== Testing Session History ==="
    
    # Start and stop a session to create history
    log_test "INFO" "Creating test session for history..."
    "$APP_BINARY" start 1m > /dev/null 2>&1
    sleep 2
    "$APP_BINARY" stop > /dev/null 2>&1
    
    # Test history command
    log_test "INFO" "Testing history command..."
    output=$("$APP_BINARY" history 2>&1)
    assert_exit_code 0 $? "History command should succeed"
    assert_output_contains "work" "$output" "History should show work session"
    
    # Verify history file exists
    assert_file_exists "${CONFIG_DIR}/session_history.json" "Session history file should exist"
}

test_error_handling() {
    log_test "INFO" "=== Testing Error Handling ==="
    
    # Test pause without running timer
    log_test "INFO" "Testing pause without running timer..."
    output=$("$APP_BINARY" pause 2>&1)
    assert_exit_code 1 $? "Pause without running timer should fail"
    assert_output_contains "error" "$output" "Should show error for invalid pause"
    
    # Test resume without paused timer
    log_test "INFO" "Testing resume without paused timer..."
    output=$("$APP_BINARY" resume 2>&1)
    assert_exit_code 1 $? "Resume without paused timer should fail"
    assert_output_contains "error" "$output" "Should show error for invalid resume"
    
    # Test stop without running timer
    log_test "INFO" "Testing stop without running timer..."
    output=$("$APP_BINARY" stop 2>&1)
    assert_exit_code 1 $? "Stop without running timer should fail"
    assert_output_contains "error" "$output" "Should show error for invalid stop"
    
    # Test start with invalid duration
    log_test "INFO" "Testing start with invalid duration..."
    output=$("$APP_BINARY" start invalid 2>&1)
    assert_exit_code 1 $? "Start with invalid duration should fail"
    assert_output_contains "error" "$output" "Should show error for invalid duration"
}

test_completion_commands() {
    log_test "INFO" "=== Testing Completion Commands ==="
    
    # Test bash completion
    log_test "INFO" "Testing bash completion..."
    output=$("$APP_BINARY" completion bash 2>&1)
    assert_exit_code 0 $? "Bash completion command should succeed"
    assert_output_contains "complete" "$output" "Completion should contain complete command"
    
    # Test zsh completion
    log_test "INFO" "Testing zsh completion..."
    output=$("$APP_BINARY" completion zsh 2>&1)
    assert_exit_code 0 $? "Zsh completion command should succeed"
    assert_output_contains "compdef" "$output" "Completion should contain compdef command"
}

test_force_flag() {
    log_test "INFO" "=== Testing Force Flag ==="
    
    # Start a timer
    "$APP_BINARY" start 1m > /dev/null 2>&1
    
    # Test start without force (should fail)
    log_test "INFO" "Testing start without force flag..."
    output=$("$APP_BINARY" start 1m 2>&1)
    assert_exit_code 1 $? "Start without force should fail when timer is running"
    assert_output_contains "error" "$output" "Should show error for existing timer"
    
    # Test start with force (should succeed)
    log_test "INFO" "Testing start with force flag..."
    output=$("$APP_BINARY" start --force 1m 2>&1)
    assert_exit_code 0 $? "Start with force should succeed"
    
    "$APP_BINARY" stop > /dev/null 2>&1
}

# Main test execution
main() {
    echo "=========================================="
    echo "Pomodux UAT Test Suite"
    echo "=========================================="
    echo "Starting comprehensive user acceptance testing..."
    echo ""
    
    # Initialize test results file
    echo "Pomodux UAT Test Results - $(date)" > "$TEST_RESULTS"
    echo "==========================================" >> "$TEST_RESULTS"
    
    # Setup test environment
    setup_test_environment
    
    # Run test scenarios
    test_basic_functionality
    test_configuration_management
    test_timer_operations
    test_session_types
    test_session_history
    test_error_handling
    test_completion_commands
    test_force_flag
    
    # Cleanup test environment
    cleanup_test_environment
    
    # Print test summary
    print_test_summary
}

# Handle script interruption
trap cleanup_test_environment EXIT

# Execute main function
main "$@" 