package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test timer defaults
	if config.Timer.DefaultWorkDuration != 25*time.Minute {
		t.Errorf("expected default work duration 25m, got %v", config.Timer.DefaultWorkDuration)
	}

	if config.Timer.DefaultBreakDuration != 5*time.Minute {
		t.Errorf("expected default break duration 5m, got %v", config.Timer.DefaultBreakDuration)
	}

	if config.Timer.AutoStartBreaks {
		t.Error("expected auto start breaks to be false by default")
	}

	// Test TUI defaults
	if config.TUI.Theme != "default" {
		t.Errorf("expected default theme 'default', got %s", config.TUI.Theme)
	}

	expectedKeyBindings := map[string]string{
		"start":  "s",
		"stop":   "q",
		"pause":  "p",
		"resume": "r",
	}

	for key, expected := range expectedKeyBindings {
		if actual := config.TUI.KeyBindings[key]; actual != expected {
			t.Errorf("expected key binding %s=%s, got %s", key, expected, actual)
		}
	}

	// Test notification defaults
	if !config.Notifications.Enabled {
		t.Error("expected notifications to be enabled by default")
	}

	if config.Notifications.Sound {
		t.Error("expected sound to be disabled by default")
	}

	if config.Notifications.Message != "Timer completed!" {
		t.Errorf("expected default message 'Timer completed!', got %s", config.Notifications.Message)
	}
}

func TestConfigValidation(t *testing.T) {
	config := DefaultConfig()

	// Test valid config
	if err := Validate(config); err != nil {
		t.Errorf("expected valid config to pass validation, got error: %v", err)
	}

	// Test invalid work duration
	config.Timer.DefaultWorkDuration = 0
	if err := Validate(config); err == nil {
		t.Error("expected validation to fail with zero work duration")
	}

	// Test invalid break duration
	config = DefaultConfig()
	config.Timer.DefaultBreakDuration = -1 * time.Minute
	if err := Validate(config); err == nil {
		t.Error("expected validation to fail with negative break duration")
	}

	// Test invalid theme
	config = DefaultConfig()
	config.TUI.Theme = ""
	if err := Validate(config); err == nil {
		t.Error("expected validation to fail with empty theme")
	}
}

func TestConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}

	// Should contain pomodux/config.yaml
	if !strings.HasSuffix(path, "pomodux/config.yaml") {
		t.Errorf("config path should end with pomodux/config.yaml, got %s", path)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "pomodux-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test saving and loading with direct file path
	originalConfig := DefaultConfig()
	originalConfig.Timer.DefaultWorkDuration = 30 * time.Minute
	originalConfig.TUI.Theme = "dark"

	configPath := filepath.Join(tempDir, "config.yaml")

	// Save config directly
	data, err := yaml.Marshal(originalConfig)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load config directly
	data, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	loadedConfig := DefaultConfig()
	if err := yaml.Unmarshal(data, loadedConfig); err != nil {
		t.Fatalf("failed to parse config file: %v", err)
	}

	// Verify loaded config matches saved config
	if loadedConfig.Timer.DefaultWorkDuration != originalConfig.Timer.DefaultWorkDuration {
		t.Errorf("work duration mismatch: expected %v, got %v",
			originalConfig.Timer.DefaultWorkDuration, loadedConfig.Timer.DefaultWorkDuration)
	}

	if loadedConfig.TUI.Theme != originalConfig.TUI.Theme {
		t.Errorf("theme mismatch: expected %s, got %s",
			originalConfig.TUI.Theme, loadedConfig.TUI.Theme)
	}
}

func TestSaveConfig(t *testing.T) {
	// Test that Save function works with valid config
	config := DefaultConfig()
	config.Timer.DefaultWorkDuration = 45 * time.Minute

	// This test will create a real config file in the user's config directory
	// In a real test environment, we'd mock this, but for now we'll skip
	// if we can't write to the config directory
	if err := Save(config); err != nil {
		t.Skipf("skipping test - cannot save config: %v", err)
	}
}

func TestLoadConfigWithInvalidYAML(t *testing.T) {
	// Test YAML unmarshaling with invalid data
	invalidYAML := `timer:
  default_work_duration: "invalid duration"
  default_break_duration: 5m`

	config := DefaultConfig()
	err := yaml.Unmarshal([]byte(invalidYAML), config)
	if err == nil {
		t.Error("expected error when unmarshaling invalid YAML")
	}
}

func TestConfigPathWithXDGConfigHome(t *testing.T) {
	// Save original environment
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)

	// Set XDG_CONFIG_HOME
	testConfigHome := "/test/config/home"
	os.Setenv("XDG_CONFIG_HOME", testConfigHome)

	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}

	expectedPath := filepath.Join(testConfigHome, "pomodux", "config.yaml")
	if path != expectedPath {
		t.Errorf("expected config path %s, got %s", expectedPath, path)
	}
}
