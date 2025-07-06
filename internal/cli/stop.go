package cli

import (
	"fmt"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current timer",
	Long: `Stop the currently running timer and reset to idle state.

This command will immediately stop any running timer session.`,
	RunE: runStopTimer,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStopTimer(cmd *cobra.Command, args []string) error {
	manager := timer.GetManager()
	status := manager.GetStatus()

	if status == timer.StatusIdle {
		fmt.Println("No active timer to stop")
		fmt.Println("Use 'pomodux start [duration]' to begin a timer")
		return nil
	}

	if err := manager.StopTimer(); err != nil {
		return fmt.Errorf("failed to stop timer: %w", err)
	}

	fmt.Println("Timer stopped")
	return nil
}
