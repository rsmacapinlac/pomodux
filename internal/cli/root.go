package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "pomodux",
		Short: "A terminal-based timer and Pomodoro application",
		Long: `Pomodux is a powerful terminal timer application that helps you 
manage your time effectively with work sessions, breaks, and Pomodoro technique.

Features:
  • Start work timers with custom durations
  • Pomodoro technique support
  • Break timer management
  • Session tracking and statistics
  • Rich terminal user interface
  • Plugin system for extensibility`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pomodux/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO: Implement configuration loading
	// This will be implemented when we add the configuration system
}
