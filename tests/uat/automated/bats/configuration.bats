#!/usr/bin/env bats

# BATS test suite for Pomodux Configuration Management
# Requires: bats-core (https://github.com/bats-core/bats-core)
# Run with: bats tests/uat/automated/bats/configuration.bats

setup() {
    # Setup test environment
    export APP_BINARY="${BATS_TEST_DIRNAME}/../../../../bin/pomodux"
    export CONFIG_DIR="${HOME}/.config/pomodux"
    
    # Create test config directory
    mkdir -p "$CONFIG_DIR"
    
    # Backup existing config if it exists
    if [ -f "${CONFIG_DIR}/config.json" ]; then
        cp "${CONFIG_DIR}/config.json" "${CONFIG_DIR}/config.json.backup"
    fi
    
    # Create test configuration
    cat > "${CONFIG_DIR}/config.json" << EOF
{
    "default_work_duration": "1m",
    "default_break_duration": "30s",
    "default_long_break_duration": "2m"
}
EOF
    
    # Ensure application is built
    cd "${BATS_TEST_DIRNAME}/../../../../"
    make build > /dev/null 2>&1 || true
}

teardown() {
    # Stop any running timer
    if [ -f "$APP_BINARY" ]; then
        "$APP_BINARY" stop > /dev/null 2>&1 || true
    fi
    
    # Restore original config
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
    else
        rm -f "${CONFIG_DIR}/config.json"
    fi
    
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
}

@test "config show command should succeed" {
    run "$APP_BINARY" config show
    [ "$status" -eq 0 ]
}

@test "config show should display work duration" {
    run "$APP_BINARY" config show --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DefaultWorkDuration" ]]
}

@test "config show should display break duration" {
    run "$APP_BINARY" config show --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DefaultBreakDuration" ]]
}

@test "config show should display long break duration" {
    run "$APP_BINARY" config show --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DefaultLongBreakDuration" ]]
}

@test "config set command should succeed with valid work duration" {
    run "$APP_BINARY" config set timer.default_work_duration 25m
    [ "$status" -eq 0 ]
}

@test "config set command should succeed with valid break duration" {
    run "$APP_BINARY" config set timer.default_break_duration 5m
    [ "$status" -eq 0 ]
}

@test "config set command should succeed with valid long break duration" {
    run "$APP_BINARY" config set timer.default_long_break_duration 15m
    [ "$status" -eq 0 ]
}

@test "config set should fail with invalid work duration" {
    run "$APP_BINARY" config set timer.default_work_duration invalid
    [ "$status" -eq 1 ]
}

@test "config set should fail with invalid break duration" {
    run "$APP_BINARY" config set timer.default_break_duration invalid
    [ "$status" -eq 1 ]
}

@test "config set should fail with invalid long break duration" {
    run "$APP_BINARY" config set timer.default_long_break_duration invalid
    [ "$status" -eq 1 ]
}

@test "config set should fail with unknown setting" {
    run "$APP_BINARY" config set unknown_setting 10m
    [ "$status" -eq 1 ]
}

@test "config set should accept various duration formats" {
    # Test minutes format
    run "$APP_BINARY" config set timer.default_work_duration 25m
    [ "$status" -eq 0 ]
    
    # Test seconds format
    run "$APP_BINARY" config set timer.default_break_duration 300s
    [ "$status" -eq 0 ]
    
    # Test numeric minutes (with explicit unit)
    run "$APP_BINARY" config set timer.default_long_break_duration 15m
    [ "$status" -eq 0 ]
}

@test "config show should display JSON format" {
    run "$APP_BINARY" config show --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ "{" ]]
    [[ "$output" =~ "}" ]]
}

@test "config set should handle edge case durations" {
    # Test very short duration
    run "$APP_BINARY" config set timer.default_break_duration 1s
    [ "$status" -eq 0 ]
    
    # Test very long duration
    run "$APP_BINARY" config set timer.default_work_duration 120m
    [ "$status" -eq 0 ]
} 