package cli

import (
	"fmt"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var longBreakCmd = &cobra.Command{
	Use:   "long-break",
	Short: "Start a 15-minute long break session",
	Long:  `Start a long break session (15 minutes by default, configurable in future releases).`,
	RunE:  runLongBreak,
}

func init() {
	rootCmd.AddCommand(longBreakCmd)
}

func runLongBreak(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		cmd.PrintErrln("Warning: Failed to load configuration, using 15 minutes as default:", err)
		duration := 15 * time.Minute
		t := timer.GetGlobalTimer()
		if t.GetStatus() == timer.StatusRunning {
			cmd.PrintErrln("Timer already running.")
			return fmt.Errorf("timer already running")
		}
		return t.StartPersistent(duration, timer.SessionTypeLongBreak)
	}
	duration := cfg.Timer.DefaultLongBreakDuration
	t := timer.GetGlobalTimer()
	if t.GetStatus() == timer.StatusRunning {
		cmd.PrintErrln("Timer already running.")
		return fmt.Errorf("timer already running")
	}
	return t.StartPersistent(duration, timer.SessionTypeLongBreak)
}
