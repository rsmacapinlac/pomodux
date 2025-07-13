package main

import (
	"fmt"
	"os"

	"github.com/rsmacapinlac/pomodux/internal/cli"
	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/logger"
	"github.com/rsmacapinlac/pomodux/internal/timer"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Error loading config", err)
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logConfig := &logger.Config{
		Level:      logger.LogLevel(cfg.Logging.Level),
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		LogFile:    cfg.Logging.LogFile,
		ShowCaller: cfg.Logging.ShowCaller,
	}
	if err := logger.Init(logConfig); err != nil {
		logger.Error("Error initializing logger", err)
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	defer timer.ShutdownGlobalTimer() // Ensure clean shutdown
	if err := cli.Execute(); err != nil {
		logger.Error("Application error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
