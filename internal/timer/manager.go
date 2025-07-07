package timer

import (
	"sync"
)

// Global timer instance
var (
	globalTimer          *Timer
	timerOnce            sync.Once
	globalStateManager   *StateManager
	globalHistoryManager *HistoryManager
)

// GetGlobalTimer returns the global timer instance.
// This ensures all CLI commands use the same timer.
func GetGlobalTimer() *Timer {
	timerOnce.Do(func() {
		// Create state manager
		stateManager, err := NewStateManager()
		if err != nil {
			// If we can't create state manager, create timer without persistence
			globalTimer = NewTimer()
			return
		}
		globalStateManager = stateManager

		// Create history manager
		historyManager, err := NewHistoryManager()
		if err != nil {
			// If we can't create history manager, create timer without history
			globalTimer = NewTimerWithManagers(stateManager, nil)
			return
		}
		globalHistoryManager = historyManager

		// Create timer with both managers
		globalTimer = NewTimerWithManagers(stateManager, historyManager)
	})

	return globalTimer
}

// ShutdownGlobalTimer gracefully shuts down the global timer.
// This should be called when the application exits.
func ShutdownGlobalTimer() {
	// Nothing to do for the simplified timer
}
