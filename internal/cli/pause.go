package cli

import (
	"fmt"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the currently running timer",
	Long: `Pause the currently running timer. The timer will stop counting down
but retain its current progress. Use 'pomodux resume' to continue the timer.`,
	RunE: runPause,
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}

func runPause(cmd *cobra.Command, args []string) error {
	t := timer.GetGlobalTimer()

	status := t.GetStatus()
	if status != timer.StatusRunning {
		cmd.PrintErrln("Cannot pause timer: timer is not running (current status:", status, ")")
		return fmt.Errorf("cannot pause timer: timer is not running (current status: %v)", status)
	}

	if err := t.Pause(); err != nil {
		cmd.PrintErrln("Failed to pause timer:", err)
		return fmt.Errorf("failed to pause timer: %w", err)
	}
	fmt.Println("Timer paused")
	return nil
}
