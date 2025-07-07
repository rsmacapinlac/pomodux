package cli

import (
	"fmt"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused timer",
	Long:  `Resume a paused timer. The timer will continue counting down from where it was paused. Use 'pomodux pause' to pause again if needed.`,
	RunE:  runResume,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func runResume(cmd *cobra.Command, args []string) error {
	t := timer.GetGlobalTimer()

	status := t.GetStatus()
	if status != timer.StatusPaused {
		cmd.PrintErrln("Cannot resume timer: timer is not paused (current status:", status, ")")
		return fmt.Errorf("cannot resume timer: timer is not paused (current status: %v)", status)
	}

	if err := t.Resume(); err != nil {
		cmd.PrintErrln("Failed to resume timer:", err)
		return fmt.Errorf("failed to resume timer: %w", err)
	}

	fmt.Println("Timer resumed.")
	return nil
}
