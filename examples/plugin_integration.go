package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/logger"
	"github.com/rsmacapinlac/pomodux/internal/plugin"
	"github.com/rsmacapinlac/pomodux/internal/timer"
)

func main() {
	logger.Info("=== Pomodux Plugin System Demo ===")
	fmt.Println("=== Pomodux Plugin System Demo ===")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", err)
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	pluginsDir := cfg.Plugins.Directory
	logger.Info("Using plugins directory", map[string]interface{}{"plugins_dir": pluginsDir})
	fmt.Println("Using plugins directory:", pluginsDir)

	// Create plugin manager
	pm := plugin.NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Load all plugins from the plugins directory
	logger.Info("Loading plugins from directory", map[string]interface{}{"plugins_dir": pluginsDir})
	fmt.Println("Loading plugins from:", pluginsDir)
	if err := pm.LoadPlugins(); err != nil {
		logger.Warn("Failed to load plugins", map[string]interface{}{"error": err.Error()})
		fmt.Printf("Warning: Failed to load plugins: %v\n", err)
	}

	// List loaded plugins
	plugins := pm.ListPlugins()
	logger.Info("Loaded plugins", map[string]interface{}{"count": len(plugins)})
	fmt.Printf("Loaded %d plugins:\n", len(plugins))
	for _, p := range plugins {
		logger.Info("Plugin loaded", map[string]interface{}{"name": p.Name, "version": p.Version, "author": p.Author})
		fmt.Printf("  - %s v%s by %s: %s\n",
			p.Name, p.Version, p.Author, p.Description)
	}

	// Create timer with plugin manager
	t := timer.NewTimerWithPluginManager(pm)

	// Demo: Start a short work session
	logger.Info("Starting 10-second work session")
	fmt.Println("\n=== Starting 10-second work session ===")

	// Start the timer
	if err := t.StartWithType(10*time.Second, timer.SessionTypeWork); err != nil {
		logger.Error("Error starting timer", err)
		fmt.Printf("Error starting timer: %v\n", err)
		return
	}

	// Simulate some timer interactions
	logger.Info("Pausing timer after 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Pausing timer...")
	if err := t.Pause(); err != nil {
		logger.Error("Error pausing timer", err)
		fmt.Printf("Error pausing timer: %v\n", err)
	}

	logger.Info("Resuming timer after 2 seconds")
	time.Sleep(2 * time.Second)
	fmt.Println("Resuming timer...")
	if err := t.Resume(); err != nil {
		logger.Error("Error resuming timer", err)
		fmt.Printf("Error resuming timer: %v\n", err)
	}

	// Wait for completion
	logger.Info("Waiting for timer completion (8 seconds)")
	time.Sleep(8 * time.Second)

	// Check if timer completed
	status := t.GetStatus()
	if status == timer.StatusCompleted {
		logger.Info("Timer completed")
		fmt.Println("Timer completed!")
	} else {
		logger.Info("Stopping timer before completion")
		fmt.Println("Stopping timer...")
		if err := t.Stop(); err != nil {
			logger.Error("Error stopping timer", err)
			fmt.Printf("Error stopping timer: %v\n", err)
		}
	}

	// Demo: Start a break session
	logger.Info("Starting 5-second break session")
	fmt.Println("\n=== Starting 5-second break session ===")

	if err := t.StartWithType(5*time.Second, timer.SessionTypeBreak); err != nil {
		logger.Error("Error starting timer", err)
		fmt.Printf("Error starting timer: %v\n", err)
		return
	}

	// Wait for completion
	time.Sleep(6 * time.Second)

	status = t.GetStatus()
	if status == timer.StatusCompleted {
		logger.Info("Break completed")
		fmt.Println("Break completed!")
	} else {
		logger.Info("Stopping break before completion")
		fmt.Println("Stopping break...")
		if err := t.Stop(); err != nil {
			logger.Error("Error stopping break", err)
			fmt.Printf("Error stopping break: %v\n", err)
		}
	}

	logger.Info("Plugin Demo Complete")
	fmt.Println("\n=== Plugin Demo Complete ===")
	fmt.Println("Check the console output above to see plugin events being logged.")
	fmt.Println("The plugins should have received events for:")
	fmt.Println("  - timer_started (when sessions began)")
	fmt.Println("  - timer_paused (when work session was paused)")
	fmt.Println("  - timer_resumed (when work session was resumed)")
	fmt.Println("  - timer_completed (when sessions finished)")
	fmt.Println("  - timer_stopped (if sessions were stopped early)")
}
