package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.3.0"
	BuildDate = "unknown"
	Commit    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pomodux version: %s\nBuild date: %s\nCommit: %s\n", Version, BuildDate, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
