package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionRecord represents a completed timer session
type SessionRecord struct {
	Type      SessionType   `json:"type"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Completed bool          `json:"completed"`
}

// HistoryManager handles session history persistence
type HistoryManager struct {
	historyFile string
	mu          sync.Mutex
}

// NewHistoryManager creates a new history manager
func NewHistoryManager() (*HistoryManager, error) {
	stateDir, err := getStateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	historyFile := filepath.Join(stateDir, "session_history.json")
	return &HistoryManager{historyFile: historyFile}, nil
}

// AddSession adds a completed session to history
func (hm *HistoryManager) AddSession(session SessionRecord) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Load existing history
	history, err := hm.loadHistory()
	if err != nil {
		return fmt.Errorf("failed to load history: %w", err)
	}

	// Add new session to the beginning (most recent first)
	history = append([]SessionRecord{session}, history...)

	// Keep only last 100 sessions to prevent file from growing too large
	if len(history) > 100 {
		history = history[:100]
	}

	// Save updated history
	return hm.saveHistory(history)
}

// GetLastSession returns the most recent completed session
func (hm *HistoryManager) GetLastSession() (*SessionRecord, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	history, err := hm.loadHistory()
	if err != nil {
		return nil, fmt.Errorf("failed to load history: %w", err)
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no session history found")
	}

	return &history[0], nil
}

// GetRecentSessions returns the most recent N completed sessions
func (hm *HistoryManager) GetRecentSessions(count int) ([]SessionRecord, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	history, err := hm.loadHistory()
	if err != nil {
		return nil, fmt.Errorf("failed to load history: %w", err)
	}

	if len(history) == 0 {
		return []SessionRecord{}, nil
	}

	// Return up to 'count' sessions, or all if less than count
	if count > len(history) {
		count = len(history)
	}

	return history[:count], nil
}

// loadHistory loads session history from file
func (hm *HistoryManager) loadHistory() ([]SessionRecord, error) {
	if _, err := os.Stat(hm.historyFile); os.IsNotExist(err) {
		return []SessionRecord{}, nil
	}

	data, err := os.ReadFile(hm.historyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var history []SessionRecord
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	return history, nil
}

// saveHistory saves session history to file
func (hm *HistoryManager) saveHistory(history []SessionRecord) error {
	// Ensure state directory exists
	stateDir := filepath.Dir(hm.historyFile)
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(hm.historyFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}
