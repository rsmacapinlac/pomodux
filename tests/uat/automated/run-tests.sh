#!/bin/bash
# Automated Test Runner for Pomodux
# Runs all BATS tests and generates reports

set -uo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
BATS_DIR="${SCRIPT_DIR}/bats"
INTEGRATION_DIR="${SCRIPT_DIR}/integration"
REPORTS_DIR="${PROJECT_ROOT}/tests/reports"
TIMESTAMP="$(date '+%Y%m%d_%H%M%S')"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Ensure variables are treated as integers
declare -i TOTAL_TESTS
declare -i PASSED_TESTS
declare -i FAILED_TESTS
declare -i SKIPPED_TESTS

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_bats_installation() {
    if ! command -v bats > /dev/null 2>&1; then
        log_error "BATS is not installed. Please install bats-core first."
        log_info "Installation instructions: https://github.com/bats-core/bats-core#installation"
        exit 1
    fi
    
    log_success "BATS is installed: $(bats --version)"
}

setup_test_environment() {
    log_info "Setting up test environment..."
    
    # Create reports directory
    mkdir -p "$REPORTS_DIR"
    
    # Ensure application is built
    cd "$PROJECT_ROOT"
    if ! make build > /dev/null 2>&1; then
        log_error "Failed to build application"
        exit 1
    fi
    
    log_success "Test environment setup complete"
}

run_test_suite() {
    local test_dir="$1"
    local suite_name="$2"
    local report_file="$3"
    
    log_info "Running $suite_name tests..."
    
    if [ ! -d "$test_dir" ]; then
        log_warning "Test directory not found: $test_dir"
        return 0
    fi
    
    # Find all .bats files in the directory
    local test_files=($(find "$test_dir" -name "*.bats" -type f))
    
    if [ ${#test_files[@]} -eq 0 ]; then
        log_warning "No .bats files found in $test_dir"
        return 0
    fi
    
    # Run tests with BATS
    local bats_output
    bats_output=$(bats --tap "$test_dir" 2>&1)
    local bats_exit_code=$?
    
    # Always save the output
    echo "$bats_output" > "$report_file"
    
    # Count tests - ensure we get valid numbers
    local test_count=$(echo "$bats_output" | grep -c "^ok [0-9]" || echo "0")
    local fail_count=$(echo "$bats_output" | grep -c "^not ok [0-9]" || echo "0")
    
    # Ensure variables are numbers, else set to 0
    [[ "$test_count" =~ ^[0-9]+$ ]] || test_count=0
    [[ "$fail_count" =~ ^[0-9]+$ ]] || fail_count=0
    
    # Calculate total count safely
    local total_count=$((test_count + fail_count))
    
    # Update global counters
    TOTAL_TESTS=$((TOTAL_TESTS + total_count))
    PASSED_TESTS=$((PASSED_TESTS + test_count))
    FAILED_TESTS=$((FAILED_TESTS + fail_count))
    
    if [ $bats_exit_code -eq 0 ]; then
        log_success "$suite_name tests passed"
        return 0
    else
        log_error "$suite_name tests failed"
        return 1
    fi
}

# Only run relevant BATS test files
BATS_TEST_FILES=(
  "${BATS_DIR}/basic-functionality.bats"
  "${BATS_DIR}/configuration.bats"
  "${BATS_DIR}/persistent-timer.bats"
)

run_all_bats_tests() {
  for test_file in "${BATS_TEST_FILES[@]}"; do
    if [ -f "$test_file" ]; then
      suite_name=$(basename "$test_file" .bats)
      report_file="$REPORTS_DIR/${suite_name}-${TIMESTAMP}.tap"
      log_info "Running $suite_name tests..."
      bats_output=$(bats --tap "$test_file" 2>&1)
      bats_exit_code=$?
      echo "$bats_output" > "$report_file"
      test_count=$(echo "$bats_output" | grep -c "^ok [0-9]" || echo "0")
      fail_count=$(echo "$bats_output" | grep -c "^not ok [0-9]" || echo "0")
      [[ "$test_count" =~ ^[0-9]+$ ]] || test_count=0
      [[ "$fail_count" =~ ^[0-9]+$ ]] || fail_count=0
      total_count=$((test_count + fail_count))
      TOTAL_TESTS=$((TOTAL_TESTS + total_count))
      PASSED_TESTS=$((PASSED_TESTS + test_count))
      FAILED_TESTS=$((FAILED_TESTS + fail_count))
      if [ $bats_exit_code -eq 0 ]; then
        log_success "$suite_name tests passed"
      else
        log_error "$suite_name tests failed"
      fi
    fi
  done
}

generate_summary_report() {
    local summary_file="$REPORTS_DIR/test-summary-${TIMESTAMP}.md"
    
    log_info "Generating summary report..."
    
    # Calculate success rate safely
    local success_rate="0%"
    if [ "$TOTAL_TESTS" -gt 0 ] && [ "$PASSED_TESTS" -ge 0 ]; then
        success_rate="$(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%"
    fi
    
    cat > "$summary_file" << EOF
# Pomodux Automated Test Summary

**Generated**: $(date)
**Timestamp**: $TIMESTAMP

## Test Results

- **Total Tests**: $TOTAL_TESTS
- **Passed**: $PASSED_TESTS
- **Failed**: $FAILED_TESTS
- **Skipped**: $SKIPPED_TESTS
- **Success Rate**: $success_rate

## Test Suites

### Unit Tests
- **Basic Functionality**: [basic-functionality.bats]($REPORTS_DIR/basic-functionality-${TIMESTAMP}.tap)
- **Configuration Management**: [configuration.bats]($REPORTS_DIR/configuration-${TIMESTAMP}.tap)
- **Persistent Timer**: [persistent-timer.bats]($REPORTS_DIR/persistent-timer-${TIMESTAMP}.tap)

### Integration Tests
- (none)

## Test Coverage

The automated tests cover:

- ✅ Basic CLI functionality (version, help, completion)
- ✅ Configuration management (show, set, validation)
- ✅ Persistent timer functionality
- ✅ Error handling and edge cases

## Test Environment

- **Application**: Pomodux CLI
- **Test Framework**: BATS (Bash Automated Testing System)
- **Configuration**: Test-specific configs with short durations
- **Isolation**: Each test runs in isolated environment

## Notes

- All tests use short durations (1m, 30s, 2m) for quick execution
- Tests automatically clean up after themselves
- Configuration is backed up and restored for each test
- Session history is cleaned between tests

EOF
    
    log_success "Summary report generated: $summary_file"
}

print_final_summary() {
    echo ""
    echo "=========================================="
    echo "Automated Test Summary"
    echo "=========================================="
    echo "Total Tests: $TOTAL_TESTS"
    echo "Passed: $PASSED_TESTS"
    echo "Failed: $FAILED_TESTS"
    echo "Skipped: $SKIPPED_TESTS"
    
    # Calculate success rate safely
    local success_rate="0%"
    if [ "$TOTAL_TESTS" -gt 0 ] && [ "$PASSED_TESTS" -ge 0 ]; then
        success_rate="$(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%"
    fi
    echo "Success Rate: $success_rate"
    
    echo ""
    echo "Reports saved to: $REPORTS_DIR"
    echo "=========================================="
    
    if [ "$FAILED_TESTS" -gt 0 ]; then
        log_error "Some tests failed. Check the reports for details."
        exit 1
    else
        log_success "All tests passed!"
    fi
}

# Main execution
main() {
    echo "=========================================="
    echo "Pomodux Automated Test Runner"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_bats_installation
    
    # Setup environment
    setup_test_environment
    
    # Run test suites
    local overall_success=true
    
    # Run all relevant BATS tests
    run_all_bats_tests
    
    # Generate summary report
    generate_summary_report
    
    # Print final summary
    print_final_summary
}

# Execute main function
main "$@" 