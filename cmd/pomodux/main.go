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
	// Initialize CLI to get flags
	rootCmd := cli.GetRootCmd()
	rootCmd.ParseFlags(os.Args[1:])

	// Get config file path from flag
	cfgFile, _ := rootCmd.Flags().GetString("config")

	// Load configuration
	var cfg *config.Config
	var err error

	if cfgFile != "" {
		// Load from specified config file
		cfg, err = config.LoadFromPath(cfgFile)
	} else {
		// Load from default location
		cfg, err = config.Load()
	}

	if err != nil {
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
