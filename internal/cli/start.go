package cli

import (
	"fmt"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [duration]",
	Short: "Start a work timer",
	Long: `Start a work timer for the specified duration.
	
Examples:
  pomodux start 25m          # Start a 25-minute work session
  pomodux start 1h30m        # Start a 1 hour 30 minute session
  pomodux start 45s          # Start a 45-second session
  
If no duration is specified, uses the default work duration from config.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var duration time.Duration
		var err error

		if len(args) > 0 {
			duration, err = time.ParseDuration(args[0])
			if err != nil {
				return fmt.Errorf("invalid duration: %v", err)
			}
		} else {
			// Use default work duration from config
			cfg, err := config.Load()
			if err != nil {
				cmd.PrintErrln("Warning: Failed to load configuration, using 25 minutes as default:", err)
				duration = 25 * time.Minute
			} else {
				duration = cfg.Timer.DefaultWorkDuration
			}
		}

		if duration <= 0 {
			return fmt.Errorf("duration must be positive")
		}

		// Get the global timer
		t := timer.GetGlobalTimer()

		// Check if timer is already running
		if t.GetStatus() == timer.StatusRunning {
			return fmt.Errorf("timer already running")
		}

		// Start the persistent timer (this will block until completion)
		return t.StartPersistent(duration, timer.SessionTypeWork)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
