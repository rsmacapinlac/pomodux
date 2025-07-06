package cli

import (
	"fmt"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current timer status",
	Long: `Show the current status of the timer including:
  • Current timer state (idle, running, paused, completed)
  • Progress percentage
  • Time remaining (if running)
  • Session information`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	t := timer.NewTimer()
	status := t.GetStatus()
	progress := t.GetProgress()

	fmt.Printf("Timer Status: %s\n", status)

	if status == timer.StatusIdle {
		fmt.Println("No active timer session")
		fmt.Println("Use 'pomodux start [duration]' to begin a timer")
	} else if status == timer.StatusCompleted {
		fmt.Printf("Progress: %.1f%%\n", progress*100)
		fmt.Println("Time Remaining: Completed")
		fmt.Println("Timer session completed! Use 'pomodux start [duration]' to begin a new session.")
	} else {
		fmt.Printf("Progress: %.1f%%\n", progress*100)
		duration := t.GetDuration()
		elapsed := t.GetElapsed()
		if duration > 0 {
			remaining := duration - elapsed
			if remaining > 0 {
				fmt.Printf("Time Remaining: %s\n", formatDuration(remaining))
			} else {
				fmt.Println("Time Remaining: Completed")
			}
		}
	}

	return nil
}
