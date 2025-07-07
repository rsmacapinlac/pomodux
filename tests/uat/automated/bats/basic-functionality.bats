#!/usr/bin/env bats

# BATS test suite for Pomodux Basic Functionality
# Requires: bats-core (https://github.com/bats-core/bats-core)
# Run with: bats tests/uat/automated/bats/basic-functionality.bats

setup() {
    # Setup test environment
    export APP_BINARY="${BATS_TEST_DIRNAME}/../../../../bin/pomodux"
    export CONFIG_DIR="${HOME}/.config/pomodux"
    export TEST_CONFIG="${BATS_TEST_DIRNAME}/../../../fixtures/configs/test-config.json"
    
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
    # Restore original config
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
    else
        rm -f "${CONFIG_DIR}/config.json"
    fi
    
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
}

@test "version command should succeed and show app name" {
    run "$APP_BINARY" version
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Pomodux" ]]
}

@test "version command should show version information" {
    run "$APP_BINARY" version
    [ "$status" -eq 0 ]
    [[ "$output" =~ [0-9]+\.[0-9]+\.[0-9]+ ]] || [[ "$output" =~ "dev" ]]
}

@test "help command should succeed" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
}

@test "help command should show start command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "start" ]]
}





@test "help command should show break command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "break" ]]
}

@test "help command should show long-break command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "long-break" ]]
}

@test "help command should show config command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "config" ]]
}

@test "help command should show history command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "history" ]]
}

@test "help command should show completion command" {
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "completion" ]]
}

@test "completion bash command should succeed" {
    run "$APP_BINARY" completion bash
    [ "$status" -eq 0 ]
}

@test "completion bash should contain complete command" {
    run "$APP_BINARY" completion bash
    [ "$status" -eq 0 ]
    [[ "$output" =~ "complete" ]]
}

@test "completion zsh command should succeed" {
    run "$APP_BINARY" completion zsh
    [ "$status" -eq 0 ]
}

@test "completion zsh should contain compdef command" {
    run "$APP_BINARY" completion zsh
    [ "$status" -eq 0 ]
    [[ "$output" =~ "compdef" ]]
}



 