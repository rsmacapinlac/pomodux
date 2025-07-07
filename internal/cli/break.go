package cli

import (
	"fmt"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Start a 5-minute break session",
	Long:  `Start a short break session (5 minutes by default, configurable in future releases).`,
	RunE:  runBreak,
}

func init() {
	rootCmd.AddCommand(breakCmd)
}

func runBreak(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		cmd.PrintErrln("Warning: Failed to load configuration, using 5 minutes as default:", err)
		duration := 5 * time.Minute
		t := timer.GetGlobalTimer()
		if t.GetStatus() == timer.StatusRunning {
			cmd.PrintErrln("Timer already running.")
			return fmt.Errorf("timer already running")
		}
		return t.StartPersistent(duration, timer.SessionTypeBreak)
	}
	duration := cfg.Timer.DefaultBreakDuration
	t := timer.GetGlobalTimer()
	if t.GetStatus() == timer.StatusRunning {
		cmd.PrintErrln("Timer already running.")
		return fmt.Errorf("timer already running")
	}
	return t.StartPersistent(duration, timer.SessionTypeBreak)
}
