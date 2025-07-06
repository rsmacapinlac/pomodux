package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [duration]",
	Short: "Start a work timer",
	Long: `Start a work timer for the specified duration or use default.

Examples:
  pomodux start          # Start with default duration (25 minutes)
  pomodux start 30m      # Start a 30-minute timer
  pomodux start 1h       # Start a 1-hour timer
  pomodux start 1500s    # Start a 1500-second timer`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStartTimer,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStartTimer(cmd *cobra.Command, args []string) error {
	var duration time.Duration
	var err error

	if len(args) == 0 {
		// Default to 25 minutes (Pomodoro technique)
		duration = 25 * time.Minute
	} else {
		duration, err = parseDuration(args[0])
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
	}

	// Use file-based state persistence
	stateManager, err := timer.NewStateManager()
	if err != nil {
		return fmt.Errorf("failed to initialize state manager: %w", err)
	}

	// Load or create timer from state
	t := timer.NewTimer()
	if t.GetStatus() == timer.StatusRunning {
		return fmt.Errorf("timer already running")
	}
	if err := t.Start(duration); err != nil {
		return fmt.Errorf("failed to start timer: %w", err)
	}
	if err := stateManager.SaveState(t); err != nil {
		return fmt.Errorf("failed to save timer state: %w", err)
	}

	fmt.Printf("Timer started for %s\n", formatDuration(duration))
	return nil
}

// parseDuration parses a duration string in various formats
func parseDuration(s string) (time.Duration, error) {
	s = strings.ToLower(s)
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}
	if minutes, err := strconv.Atoi(s); err == nil {
		return time.Duration(minutes) * time.Minute, nil
	}
	return 0, fmt.Errorf("unable to parse duration: %s", s)
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	if d >= time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%d hour%s", hours, plural(hours))
		}
		return fmt.Sprintf("%d hour%s %d minute%s", hours, plural(hours), minutes, plural(minutes))
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if seconds == 0 {
		return fmt.Sprintf("%d minute%s", minutes, plural(minutes))
	}
	return fmt.Sprintf("%d minute%s %d second%s", minutes, plural(minutes), seconds, plural(seconds))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
