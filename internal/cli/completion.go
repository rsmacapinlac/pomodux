package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for bash, zsh, or fish.

Examples:
  pomodux completion bash > /etc/bash_completion.d/pomodux
  pomodux completion zsh > "${fpath[1]}/_pomodux"
  pomodux completion fish > ~/.config/fish/completions/pomodux.fish
`,
	Args: cobra.ExactArgs(1),
	RunE: runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	default:
		cmd.PrintErrln("Unsupported shell:", args[0])
		return fmt.Errorf("unsupported shell: %s", args[0])
	}
}
