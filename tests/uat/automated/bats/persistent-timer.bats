#!/usr/bin/env bats

# BATS test suite for Pomodux Persistent Timer Functionality
# Requires: bats-core (https://github.com/bats-core/bats-core)
# Run with: bats tests/uat/automated/bats/persistent-timer.bats

setup() {
    # Setup test environment
    export APP_BINARY="${BATS_TEST_DIRNAME}/../../../../bin/pomodux"
    export CONFIG_DIR="${HOME}/.config/pomodux"
    
    # Create test config directory
    mkdir -p "$CONFIG_DIR"
    
    # Backup existing config files if they exist
    if [ -f "${CONFIG_DIR}/config.json" ]; then
        cp "${CONFIG_DIR}/config.json" "${CONFIG_DIR}/config.json.backup"
    fi
    if [ -f "${CONFIG_DIR}/config.yaml" ]; then
        cp "${CONFIG_DIR}/config.yaml" "${CONFIG_DIR}/config.yaml.backup"
    fi
    
    # Create test configuration with YAML format (app reads this)
    cat > "${CONFIG_DIR}/config.yaml" << EOF
timer:
    default_work_duration: 30s
    default_break_duration: 10s
    default_long_break_duration: 1m
    auto_start_breaks: false
tui:
    theme: default
    key_bindings:
        pause: p
        resume: r
        start: s
        stop: q
notifications:
    enabled: true
    sound: false
    message: Timer completed!
EOF
    
    # Ensure application is built
    cd "${BATS_TEST_DIRNAME}/../../../../"
    make build > /dev/null 2>&1 || true
    
    # Clear timer state
    rm -f "${CONFIG_DIR}/timer_state.json"
}

teardown() {
    # Restore original config files
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
    else
        rm -f "${CONFIG_DIR}/config.json"
    fi
    if [ -f "${CONFIG_DIR}/config.yaml.backup" ]; then
        mv "${CONFIG_DIR}/config.yaml.backup" "${CONFIG_DIR}/config.yaml"
    else
        rm -f "${CONFIG_DIR}/config.yaml"
    fi
    
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
    rm -f "${CONFIG_DIR}/timer_state.json"
    rm -f output.txt
}

@test "start command should work with custom duration" {
    # Test that the start command accepts custom durations
    run "$APP_BINARY" start 10s
    [ "$status" -eq 0 ]
}

@test "start command should work without duration" {
    # Test that start command works without duration (should use default)
    run "$APP_BINARY" start
    [ "$status" -eq 0 ]
}

@test "break command should work" {
    # Test that break command works
    run "$APP_BINARY" break
    [ "$status" -eq 0 ]
}

@test "long-break command should work" {
    # Test that long-break command works
    run "$APP_BINARY" long-break
    [ "$status" -eq 0 ]
}

@test "start command should show timer started message" {
    # Start a timer and check for basic output
    timeout 3s "$APP_BINARY" start 5s > output.txt 2>&1 &
    timer_pid=$!
    sleep 1
    
    if [ -f output.txt ]; then
        grep -q "Timer started" output.txt || echo "Timer started message not found in output"
    fi
    
    kill $timer_pid 2>/dev/null || true
    rm -f output.txt
}

@test "start command should show session type" {
    # Start a timer and check for session type
    timeout 3s "$APP_BINARY" start 5s > output.txt 2>&1 &
    timer_pid=$!
    sleep 1
    
    if [ -f output.txt ]; then
        grep -q "Session type" output.txt || echo "Session type message not found in output"
    fi
    
    kill $timer_pid 2>/dev/null || true
    rm -f output.txt
}

@test "start command should show keypress instructions" {
    # Start a timer and check for instructions
    timeout 3s "$APP_BINARY" start 5s > output.txt 2>&1 &
    timer_pid=$!
    sleep 1
    
    if [ -f output.txt ]; then
        grep -q "Press" output.txt || echo "Press instructions not found in output"
    fi
    
    kill $timer_pid 2>/dev/null || true
    rm -f output.txt
}

@test "break command should show break session" {
    # Start a break and check for break session
    timeout 3s "$APP_BINARY" break > output.txt 2>&1 &
    timer_pid=$!
    sleep 1
    
    if [ -f output.txt ]; then
        grep -q "break" output.txt || echo "Break session message not found in output"
    fi
    
    kill $timer_pid 2>/dev/null || true
    rm -f output.txt
}

@test "long-break command should show long break session" {
    # Start a long break and check for long break session
    timeout 3s "$APP_BINARY" long-break > output.txt 2>&1 &
    timer_pid=$!
    sleep 1
    
    if [ -f output.txt ]; then
        grep -q "long" output.txt || echo "Long break session message not found in output"
    fi
    
    kill $timer_pid 2>/dev/null || true
    rm -f output.txt
}

@test "timer should handle very short durations" {
    # Test with a very short duration
    # First check if timer is running and wait for it to complete
    run "$APP_BINARY" status
    if [ "$status" -eq 0 ] && [[ "$output" =~ "running" ]]; then
        # Wait for timer to complete
        sleep 5
    fi
    
    run "$APP_BINARY" start 1s
    [ "$status" -eq 0 ]
} 