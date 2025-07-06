package timer

import (
	"fmt"
	"sync"
	"time"
)

// Timer implements the TimerEngine interface.
type Timer struct {
	mu           sync.Mutex
	status       TimerStatus
	startTime    time.Time
	duration     time.Duration
	elapsed      time.Duration
	stateManager *StateManager
}

// NewTimer creates a new Timer instance.
func NewTimer() *Timer {
	stateManager, err := NewStateManager()
	if err != nil {
		// If we can't create state manager, create timer without persistence
		return &Timer{
			status: StatusIdle,
		}
	}

	timer := &Timer{
		status:       StatusIdle,
		stateManager: stateManager,
	}

	// Load existing state if available
	if state, err := stateManager.LoadState(); err == nil {
		timer.status = state.Status
		timer.duration = state.Duration
		timer.startTime = state.StartTime
		timer.elapsed = state.Elapsed
	}

	return timer
}

// Start begins the timer for the specified duration.
func (t *Timer) Start(duration time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status == StatusRunning {
		return fmt.Errorf("timer already running")
	}
	if duration <= 0 {
		return fmt.Errorf("invalid duration")
	}
	t.duration = duration
	t.startTime = time.Now()
	t.elapsed = 0
	t.status = StatusRunning
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
	// TODO: Emit EventTimerStarted
	return nil
}

// Stop stops the timer.
func (t *Timer) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status == StatusIdle {
		return fmt.Errorf("timer not running")
	}
	t.status = StatusIdle
	t.elapsed = 0
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
	// TODO: Emit EventTimerStopped
	return nil
}

// Pause pauses the timer.
func (t *Timer) Pause() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status != StatusRunning {
		return fmt.Errorf("timer not running")
	}
	t.elapsed += time.Since(t.startTime)
	t.status = StatusPaused
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
	// TODO: Emit EventTimerPaused
	return nil
}

// Resume resumes a paused timer.
func (t *Timer) Resume() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status != StatusPaused {
		return fmt.Errorf("timer not paused")
	}
	t.startTime = time.Now()
	t.status = StatusRunning
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
	// TODO: Emit EventTimerResumed
	return nil
}

// GetStatus returns the current timer status.
func (t *Timer) GetStatus() TimerStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Check if timer has completed
	if t.status == StatusRunning {
		elapsed := t.elapsed + time.Since(t.startTime)
		if elapsed >= t.duration {
			t.status = StatusCompleted
			t.elapsed = t.duration
			// Save state when timer completes
			if t.stateManager != nil {
				t.stateManager.SaveState(t)
			}
		}
	}
	return t.status
}

// GetProgress returns the progress as a float64 between 0 and 1.
func (t *Timer) GetProgress() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.duration == 0 {
		return 0
	}
	var elapsed time.Duration
	if t.status == StatusRunning {
		elapsed = t.elapsed + time.Since(t.startTime)
		// Check if timer has completed
		if elapsed >= t.duration {
			t.status = StatusCompleted
			t.elapsed = t.duration
			// Save state when timer completes
			if t.stateManager != nil {
				t.stateManager.SaveState(t)
			}
		}
	} else {
		elapsed = t.elapsed
	}
	progress := float64(elapsed) / float64(t.duration)
	if progress > 1 {
		progress = 1
	}
	return progress
}

// GetDuration returns the timer duration
func (t *Timer) GetDuration() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.duration
}

// GetElapsed returns the elapsed time
func (t *Timer) GetElapsed() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	var elapsed time.Duration
	if t.status == StatusRunning {
		elapsed = t.elapsed + time.Since(t.startTime)
	} else {
		elapsed = t.elapsed
	}
	return elapsed
}

// Reset resets a completed timer to idle state
func (t *Timer) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status != StatusCompleted {
		return fmt.Errorf("can only reset completed timer")
	}
	t.status = StatusIdle
	t.elapsed = 0
	t.duration = 0
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
	return nil
}
