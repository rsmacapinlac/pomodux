package timer

import (
	"testing"
	"time"
)

func TestManagerSingleton(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	// Test that GetManager returns the same instance
	manager1 := GetManager()
	manager2 := GetManager()

	if manager1 != manager2 {
		t.Error("GetManager should return the same instance")
	}
}

func TestManagerTimerLifecycle(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	manager := GetManager()
	duration := 100 * time.Millisecond

	t.Run("Start", func(t *testing.T) {
		err := manager.StartTimer(duration)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if manager.GetStatus() != StatusRunning {
			t.Errorf("expected status running, got %v", manager.GetStatus())
		}
	})

	t.Run("Status", func(t *testing.T) {
		status := manager.GetStatus()
		if status != StatusRunning {
			t.Errorf("expected status running, got %v", status)
		}

		progress := manager.GetProgress()
		if progress <= 0 || progress >= 1 {
			t.Errorf("expected progress between 0 and 1, got %f", progress)
		}
	})

	t.Run("Stop", func(t *testing.T) {
		err := manager.StopTimer()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if manager.GetStatus() != StatusIdle {
			t.Errorf("expected status idle, got %v", manager.GetStatus())
		}
	})
}

func TestManagerStatePersistence(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	manager := GetManager()

	// Start a timer
	duration := 50 * time.Millisecond
	err := manager.StartTimer(duration)
	if err != nil {
		t.Fatalf("failed to start timer: %v", err)
	}

	// Verify timer is running
	if manager.GetStatus() != StatusRunning {
		t.Errorf("expected status running, got %v", manager.GetStatus())
	}

	// Get a new manager instance (simulating different command)
	manager2 := GetManager()

	// Verify the same timer state is accessible
	if manager2.GetStatus() != StatusRunning {
		t.Errorf("expected status running in new manager, got %v", manager2.GetStatus())
	}

	// Stop timer using second manager
	err = manager2.StopTimer()
	if err != nil {
		t.Fatalf("failed to stop timer: %v", err)
	}

	// Verify both managers see the stopped state
	if manager.GetStatus() != StatusIdle {
		t.Errorf("expected status idle in original manager, got %v", manager.GetStatus())
	}
	if manager2.GetStatus() != StatusIdle {
		t.Errorf("expected status idle in new manager, got %v", manager2.GetStatus())
	}
}

func TestManagerGetTimer(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	manager := GetManager()

	// Get timer instance
	timer := manager.GetTimer()
	if timer == nil {
		t.Error("GetTimer should not return nil")
	}

	// Verify it's the same timer instance
	timer2 := manager.GetTimer()
	if timer != timer2 {
		t.Error("GetTimer should return the same timer instance")
	}
}
