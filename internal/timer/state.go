package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State represents the persistent timer state
type State struct {
	Status    TimerStatus   `json:"status"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
	Elapsed   time.Duration `json:"elapsed"`
}

// StateManager handles persistent timer state
type StateManager struct {
	stateFile string
	mu        sync.Mutex
}

// NewStateManager creates a new state manager
func NewStateManager() (*StateManager, error) {
	stateDir, err := getStateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	stateFile := filepath.Join(stateDir, "timer_state.json")
	return &StateManager{stateFile: stateFile}, nil
}

// SaveState saves the current timer state to file
func (sm *StateManager) SaveState(timer *Timer) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state := State{
		Status:    timer.status,
		Duration:  timer.duration,
		StartTime: timer.startTime,
		Elapsed:   timer.elapsed,
	}

	// Ensure state directory exists
	stateDir := filepath.Dir(sm.stateFile)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(sm.stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// LoadState loads timer state from file
func (sm *StateManager) LoadState() (*State, error) {
	if _, err := os.Stat(sm.stateFile); os.IsNotExist(err) {
		return &State{Status: StatusIdle}, nil
	}

	data, err := os.ReadFile(sm.stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &state, nil
}

// ClearState removes the state file
func (sm *StateManager) ClearState() error {
	if _, err := os.Stat(sm.stateFile); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to clear
	}

	return os.Remove(sm.stateFile)
}

// getStateDir returns the XDG-compliant state directory
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
