package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show recent timer session(s)",
	Long:  `Show recent completed timer sessions.`,
	RunE:  runHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	historyManager, err := timer.NewHistoryManager()
	if err != nil {
		return fmt.Errorf("failed to create history manager: %w", err)
	}

	sessions, err := historyManager.GetRecentSessions(10)
	if err != nil {
		return fmt.Errorf("failed to get session history: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No session history found.")
		return nil
	}

	fmt.Printf("Recent completed sessions (%d total):\n\n", len(sessions))
	for i, session := range sessions {
		fmt.Printf("%d. %s Session\n", i+1, session.Type)
		fmt.Printf("   Duration: %s\n", session.Duration)
		fmt.Printf("   Start: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   End: %s\n", session.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Actual Duration: %s\n", session.EndTime.Sub(session.StartTime))
		fmt.Printf("   Completed: %t\n\n", session.Completed)
	}

	return nil
}

func getStateDir() (string, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		stateHome = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(stateHome, "pomodux"), nil
}
