package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsmacapinlac/pomodux/internal/cli"
	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/logger"
	"github.com/rsmacapinlac/pomodux/internal/timer"
)

// Version will be set during build time via ldflags
var Version = "dev"

func main() {
	// Get config file path from flag (will be parsed during Execute)
	cfgFile := ""
	for i, arg := range os.Args[1:] {
		if arg == "--config" && i+1 < len(os.Args[1:]) {
			cfgFile = os.Args[1:][i+1]
			break
		} else if strings.HasPrefix(arg, "--config=") {
			cfgFile = strings.TrimPrefix(arg, "--config=")
			break
		}
	}

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
		// Check if this is a help request (which is not an error)
		if err.Error() == "pflag: help requested" {
			os.Exit(0)
		}
		logger.Error("Application error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
