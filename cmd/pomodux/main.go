package main

import (
	"fmt"
	"os"

	"github.com/rsmacapinlac/pomodux/internal/cli"
	"github.com/rsmacapinlac/pomodux/internal/timer"
)

func main() {
	defer timer.ShutdownGlobalTimer() // Ensure clean shutdown
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
