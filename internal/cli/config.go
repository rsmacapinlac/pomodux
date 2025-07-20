package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Pomodux configuration",
	Long:  `View, edit, or reset your Pomodux configuration file with enhanced validation and templates.`,
}

var (
	configJSON     bool
	configValidate bool
	configBackup   string
	configRestore  string
	configTemplate string
	configNested   string
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if configValidate {
			if err := validateConfig(cfg); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}
			fmt.Println("✅ Configuration is valid")
		}

		if configJSON {
			// Convert to JSON for proper output
			jsonData, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			// Pretty print configuration
			showConfigPretty(cfg)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value using nested paths. Examples:
  pomodux config set timer.default_work_duration 25m
  pomodux config set timer.default_break_duration 5m
  pomodux config set timer.default_long_break_duration 15m`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Handle nested configuration paths
		if configNested != "" {
			key = configNested + "." + key
		}

		// Parse and set the value
		if err := setConfigValue(cfg, key, value); err != nil {
			return fmt.Errorf("failed to set config value: %w", err)
		}

		// Validate the updated configuration
		if err := validateConfig(cfg); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		// Save the updated config
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Updated %s to %s\n", key, value)
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration in your default editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := getConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		// Validate editor is a safe executable
		if !isSafeExecutable(editor) {
			return fmt.Errorf("unsafe editor executable: %s", editor)
		}
		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to reset config: %w", err)
		}
		fmt.Println("✅ Configuration reset to defaults.")
		return nil
	},
}

var configBackupCmd = &cobra.Command{
	Use:   "backup [path]",
	Short: "Backup current configuration",
	Long:  `Backup current configuration to a file. If no path is specified, uses ~/.config/pomodux/config.yaml.backup`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := configBackup
		if len(args) > 0 {
			backupPath = args[0]
		}
		if backupPath == "" {
			configDir, err := getConfigDir()
			if err != nil {
				return fmt.Errorf("failed to get config directory: %w", err)
			}
			backupPath = filepath.Join(configDir, "config.yaml.backup")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := config.SaveToPath(cfg, backupPath); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}

		fmt.Printf("✅ Configuration backed up to: %s\n", backupPath)
		return nil
	},
}

var configRestoreCmd = &cobra.Command{
	Use:   "restore [path]",
	Short: "Restore configuration from backup",
	Long:  `Restore configuration from a backup file. If no path is specified, uses ~/.config/pomodux/config.yaml.backup`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		restorePath := configRestore
		if len(args) > 0 {
			restorePath = args[0]
		}
		if restorePath == "" {
			configDir, err := getConfigDir()
			if err != nil {
				return fmt.Errorf("failed to get config directory: %w", err)
			}
			restorePath = filepath.Join(configDir, "config.yaml.backup")
		}

		// Load the backup configuration
		backupCfg, err := config.LoadFromPath(restorePath)
		if err != nil {
			return fmt.Errorf("failed to load backup config: %w", err)
		}

		// Validate the backup configuration
		if err := validateConfig(backupCfg); err != nil {
			return fmt.Errorf("backup configuration validation failed: %w", err)
		}

		// Save as current configuration
		if err := config.Save(backupCfg); err != nil {
			return fmt.Errorf("failed to restore config: %w", err)
		}

		fmt.Printf("✅ Configuration restored from: %s\n", restorePath)
		return nil
	},
}

var configTemplateCmd = &cobra.Command{
	Use:   "template [name]",
	Short: "Apply configuration template",
	Long:  `Apply a predefined configuration template. Available templates: default, productivity, short-breaks, long-breaks`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]

		var templateCfg *config.Config
		switch templateName {
		case "default":
			templateCfg = config.DefaultConfig()
		case "productivity":
			templateCfg = getProductivityTemplate()
		case "short-breaks":
			templateCfg = getShortBreaksTemplate()
		case "long-breaks":
			templateCfg = getLongBreaksTemplate()
		default:
			return fmt.Errorf("unknown template: %s. Available templates: default, productivity, short-breaks, long-breaks", templateName)
		}

		if err := config.Save(templateCfg); err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}

		fmt.Printf("✅ Applied %s configuration template\n", templateName)
		return nil
	},
}

func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func getConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, "pomodux"), nil
}

func setConfigValue(cfg *config.Config, key, value string) error {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "timer":
		if len(parts) != 2 {
			return fmt.Errorf("invalid timer configuration key: %s", key)
		}

		duration, err := parseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}

		switch parts[1] {
		case "default_work_duration":
			cfg.Timer.DefaultWorkDuration = duration
		case "default_break_duration":
			cfg.Timer.DefaultBreakDuration = duration
		case "default_long_break_duration":
			cfg.Timer.DefaultLongBreakDuration = duration
		default:
			return fmt.Errorf("unknown timer setting: %s", parts[1])
		}
	case "logging":
		if len(parts) != 3 {
			return fmt.Errorf("logging configuration requires exactly 2 parts: logging.<setting> <value>")
		}
		switch parts[1] {
		case "level":
			// Validate log level
			validLevels := []string{"debug", "info", "warn", "error"}
			isValid := false
			for _, level := range validLevels {
				if value == level {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid log level: %s (valid: %s)", value, strings.Join(validLevels, ", "))
			}
			cfg.Logging.Level = value
		case "format":
			// Validate format
			if value != "text" && value != "json" {
				return fmt.Errorf("invalid log format: %s (valid: text, json)", value)
			}
			cfg.Logging.Format = value
		case "output":
			// Validate output
			if value != "console" && value != "file" && value != "both" {
				return fmt.Errorf("invalid log output: %s (valid: console, file, both)", value)
			}
			cfg.Logging.Output = value
		case "log_file":
			cfg.Logging.LogFile = value
		case "show_caller":
			if value == "true" {
				cfg.Logging.ShowCaller = true
			} else if value == "false" {
				cfg.Logging.ShowCaller = false
			} else {
				return fmt.Errorf("show_caller must be true or false")
			}
		default:
			return fmt.Errorf("unknown logging setting: %s", parts[1])
		}
	default:
		return fmt.Errorf("unknown configuration section: %s", parts[0])
	}

	return nil
}

func validateConfig(cfg *config.Config) error {
	var errors []string

	// Validate timer durations
	if cfg.Timer.DefaultWorkDuration <= 0 {
		errors = append(errors, "default_work_duration must be positive")
	}
	if cfg.Timer.DefaultBreakDuration <= 0 {
		errors = append(errors, "default_break_duration must be positive")
	}
	if cfg.Timer.DefaultLongBreakDuration <= 0 {
		errors = append(errors, "default_long_break_duration must be positive")
	}

	// Validate reasonable ranges
	if cfg.Timer.DefaultWorkDuration < 1*time.Minute {
		errors = append(errors, "default_work_duration should be at least 1 minute")
	}
	if cfg.Timer.DefaultWorkDuration > 4*time.Hour {
		errors = append(errors, "default_work_duration should not exceed 4 hours")
	}
	if cfg.Timer.DefaultBreakDuration < 30*time.Second {
		errors = append(errors, "default_break_duration should be at least 30 seconds")
	}
	if cfg.Timer.DefaultBreakDuration > 30*time.Minute {
		errors = append(errors, "default_break_duration should not exceed 30 minutes")
	}
	if cfg.Timer.DefaultLongBreakDuration < 5*time.Minute {
		errors = append(errors, "default_long_break_duration should be at least 5 minutes")
	}
	if cfg.Timer.DefaultLongBreakDuration > 2*time.Hour {
		errors = append(errors, "default_long_break_duration should not exceed 2 hours")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation errors:\n  %s", strings.Join(errors, "\n  "))
	}

	return nil
}

func showConfigPretty(cfg *config.Config) {
	fmt.Println("Pomodux Configuration:")
	fmt.Println("======================")
	fmt.Printf("Timer Settings:\n")
	fmt.Printf("  Default Work Duration:     %s\n", formatDuration(cfg.Timer.DefaultWorkDuration))
	fmt.Printf("  Default Break Duration:    %s\n", formatDuration(cfg.Timer.DefaultBreakDuration))
	fmt.Printf("  Default Long Break Duration: %s\n", formatDuration(cfg.Timer.DefaultLongBreakDuration))
	fmt.Printf("\nLogging Settings:\n")
	fmt.Printf("  Level:      %s\n", cfg.Logging.Level)
	fmt.Printf("  Format:     %s\n", cfg.Logging.Format)
	fmt.Printf("  Output:     %s\n", cfg.Logging.Output)
	if cfg.Logging.LogFile != "" {
		fmt.Printf("  Log File:   %s\n", cfg.Logging.LogFile)
	}
	fmt.Printf("  Show Caller: %t\n", cfg.Logging.ShowCaller)
}

func getProductivityTemplate() *config.Config {
	return &config.Config{
		Timer: struct {
			DefaultWorkDuration      time.Duration `yaml:"default_work_duration"`
			DefaultBreakDuration     time.Duration `yaml:"default_break_duration"`
			DefaultLongBreakDuration time.Duration `yaml:"default_long_break_duration"`
			AutoStartBreaks          bool          `yaml:"auto_start_breaks"`
		}{
			DefaultWorkDuration:      25 * time.Minute,
			DefaultBreakDuration:     5 * time.Minute,
			DefaultLongBreakDuration: 15 * time.Minute,
			AutoStartBreaks:          false,
		},
		TUI: struct {
			Theme       string            `yaml:"theme"`
			KeyBindings map[string]string `yaml:"key_bindings"`
		}{
			Theme: "default",
			KeyBindings: map[string]string{
				"start":  "s",
				"stop":   "q",
				"pause":  "p",
				"resume": "r",
			},
		},
		Notifications: struct {
			Enabled bool   `yaml:"enabled"`
			Sound   bool   `yaml:"sound"`
			Message string `yaml:"message"`
		}{
			Enabled: true,
			Sound:   false,
			Message: "Timer completed!",
		},
	}
}

func getShortBreaksTemplate() *config.Config {
	return &config.Config{
		Timer: struct {
			DefaultWorkDuration      time.Duration `yaml:"default_work_duration"`
			DefaultBreakDuration     time.Duration `yaml:"default_break_duration"`
			DefaultLongBreakDuration time.Duration `yaml:"default_long_break_duration"`
			AutoStartBreaks          bool          `yaml:"auto_start_breaks"`
		}{
			DefaultWorkDuration:      45 * time.Minute,
			DefaultBreakDuration:     3 * time.Minute,
			DefaultLongBreakDuration: 10 * time.Minute,
			AutoStartBreaks:          false,
		},
		TUI: struct {
			Theme       string            `yaml:"theme"`
			KeyBindings map[string]string `yaml:"key_bindings"`
		}{
			Theme: "default",
			KeyBindings: map[string]string{
				"start":  "s",
				"stop":   "q",
				"pause":  "p",
				"resume": "r",
			},
		},
		Notifications: struct {
			Enabled bool   `yaml:"enabled"`
			Sound   bool   `yaml:"sound"`
			Message string `yaml:"message"`
		}{
			Enabled: true,
			Sound:   false,
			Message: "Timer completed!",
		},
	}
}

func getLongBreaksTemplate() *config.Config {
	return &config.Config{
		Timer: struct {
			DefaultWorkDuration      time.Duration `yaml:"default_work_duration"`
			DefaultBreakDuration     time.Duration `yaml:"default_break_duration"`
			DefaultLongBreakDuration time.Duration `yaml:"default_long_break_duration"`
			AutoStartBreaks          bool          `yaml:"auto_start_breaks"`
		}{
			DefaultWorkDuration:      90 * time.Minute,
			DefaultBreakDuration:     15 * time.Minute,
			DefaultLongBreakDuration: 30 * time.Minute,
			AutoStartBreaks:          false,
		},
		TUI: struct {
			Theme       string            `yaml:"theme"`
			KeyBindings map[string]string `yaml:"key_bindings"`
		}{
			Theme: "default",
			KeyBindings: map[string]string{
				"start":  "s",
				"stop":   "q",
				"pause":  "p",
				"resume": "r",
			},
		},
		Notifications: struct {
			Enabled bool   `yaml:"enabled"`
			Sound   bool   `yaml:"sound"`
			Message string `yaml:"message"`
		}{
			Enabled: true,
			Sound:   false,
			Message: "Timer completed!",
		},
	}
}

// parseDuration parses a duration string in various formats
func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// isSafeExecutable checks if the given executable is safe to run
func isSafeExecutable(editor string) bool {
	// List of safe editors
	safeEditors := []string{
		"vi", "vim", "nano", "emacs", "code", "subl", "atom",
		"gedit", "kate", "mousepad", "leafpad", "geany",
	}

	// Check if editor is in the safe list
	for _, safe := range safeEditors {
		if editor == safe {
			return true
		}
	}

	// For other editors, check if they exist in PATH and are not in dangerous locations
	if strings.Contains(editor, "/") {
		// Absolute or relative path - be more restrictive
		return false
	}

	// Check if executable exists in PATH
	_, err := exec.LookPath(editor)
	return err == nil
}

func init() {
	// Add flags to commands
	configShowCmd.Flags().BoolVar(&configJSON, "json", false, "Output as JSON")
	configShowCmd.Flags().BoolVar(&configValidate, "validate", false, "Validate configuration")
	configSetCmd.Flags().StringVar(&configNested, "nested", "", "Use nested configuration path")

	// Add subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configBackupCmd)
	configCmd.AddCommand(configRestoreCmd)
	configCmd.AddCommand(configTemplateCmd)
	rootCmd.AddCommand(configCmd)
}
