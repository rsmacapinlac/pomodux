package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the Pomodux configuration structure
type Config struct {
	Timer struct {
		DefaultWorkDuration      time.Duration `yaml:"default_work_duration"`
		DefaultBreakDuration     time.Duration `yaml:"default_break_duration"`
		DefaultLongBreakDuration time.Duration `yaml:"default_long_break_duration"`
		AutoStartBreaks          bool          `yaml:"auto_start_breaks"`
	} `yaml:"timer"`

	TUI struct {
		Theme       string            `yaml:"theme"`
		KeyBindings map[string]string `yaml:"key_bindings"`
	} `yaml:"tui"`

	Notifications struct {
		Enabled bool   `yaml:"enabled"`
		Sound   bool   `yaml:"sound"`
		Message string `yaml:"message"`
	} `yaml:"notifications"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	config := &Config{}

	// Timer defaults
	config.Timer.DefaultWorkDuration = 25 * time.Minute
	config.Timer.DefaultBreakDuration = 5 * time.Minute
	config.Timer.DefaultLongBreakDuration = 15 * time.Minute
	config.Timer.AutoStartBreaks = false

	// TUI defaults
	config.TUI.Theme = "default"
	config.TUI.KeyBindings = map[string]string{
		"start":  "s",
		"stop":   "q",
		"pause":  "p",
		"resume": "r",
	}

	// Notification defaults
	config.Notifications.Enabled = true
	config.Notifications.Sound = false
	config.Notifications.Message = "Timer completed!"

	return config
}

// Load loads configuration from the default XDG location
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// If config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := Save(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := Validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves configuration to the default XDG location
func Save(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the XDG-compliant configuration file path
func getConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configHome = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configHome, "pomodux", "config.yaml"), nil
}

// Validate validates the configuration values
func Validate(config *Config) error {
	if config.Timer.DefaultWorkDuration <= 0 {
		return fmt.Errorf("default work duration must be positive")
	}

	if config.Timer.DefaultBreakDuration <= 0 {
		return fmt.Errorf("default break duration must be positive")
	}

	if config.Timer.DefaultLongBreakDuration <= 0 {
		return fmt.Errorf("default long break duration must be positive")
	}

	if config.TUI.Theme == "" {
		return fmt.Errorf("theme cannot be empty")
	}

	return nil
}
