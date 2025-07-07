package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Pomodux configuration",
	Long:  `View, edit, or reset your Pomodux configuration file.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Convert to JSON for proper output
		jsonData, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		fmt.Println(string(jsonData))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: default_work_duration, default_break_duration, default_long_break_duration`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Parse duration value using the existing parseDuration function
		duration, err := parseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}

		// Update the appropriate field
		switch key {
		case "default_work_duration":
			cfg.Timer.DefaultWorkDuration = duration
		case "default_break_duration":
			cfg.Timer.DefaultBreakDuration = duration
		case "default_long_break_duration":
			cfg.Timer.DefaultLongBreakDuration = duration
		default:
			return fmt.Errorf("unknown setting: %s", key)
		}

		// Save the updated config
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Updated %s to %s\n", key, formatDuration(duration))
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
		fmt.Println("Configuration reset to defaults.")
		return nil
	},
}

func getConfigPath() (string, error) {
	// Use the unexported getConfigPath from config package via a wrapper here
	// (since it's not exported, duplicate logic)
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configHome = homeDir + "/.config"
	}
	return configHome + "/pomodux/config.yaml", nil
}

// parseDuration parses a duration string in various formats
func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
