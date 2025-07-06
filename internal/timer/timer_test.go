package timer

import (
	"testing"
	"time"
)

func TestTimerLifecycle(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	timer := NewTimer()
	duration := 100 * time.Millisecond

	t.Run("Start", func(t *testing.T) {
		err := timer.Start(duration)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if timer.GetStatus() != StatusRunning {
			t.Errorf("expected status running, got %v", timer.GetStatus())
		}
	})

	t.Run("Pause", func(t *testing.T) {
		err := timer.Pause()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if timer.GetStatus() != StatusPaused {
			t.Errorf("expected status paused, got %v", timer.GetStatus())
		}
	})

	t.Run("Resume", func(t *testing.T) {
		err := timer.Resume()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if timer.GetStatus() != StatusRunning {
			t.Errorf("expected status running, got %v", timer.GetStatus())
		}
	})

	t.Run("Stop", func(t *testing.T) {
		err := timer.Stop()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if timer.GetStatus() != StatusIdle {
			t.Errorf("expected status idle, got %v", timer.GetStatus())
		}
	})
}

func TestTimerErrors(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	t.Run("StartAlreadyRunning", func(t *testing.T) {
		// Clear state before this test
		if stateManager, err := NewStateManager(); err == nil {
			stateManager.ClearState()
		}
		timer := NewTimer()
		timer.Start(1 * time.Minute)

		err := timer.Start(2 * time.Minute)
		if err == nil {
			t.Error("expected error when starting already running timer")
		}
	})

	t.Run("StartInvalidDuration", func(t *testing.T) {
		// Clear state before this test
		if stateManager, err := NewStateManager(); err == nil {
			stateManager.ClearState()
		}
		timer := NewTimer()

		err := timer.Start(0)
		if err == nil {
			t.Error("expected error when starting with zero duration")
		}

		err = timer.Start(-1 * time.Minute)
		if err == nil {
			t.Error("expected error when starting with negative duration")
		}
	})

	t.Run("PauseNotRunning", func(t *testing.T) {
		// Clear state before this test
		if stateManager, err := NewStateManager(); err == nil {
			stateManager.ClearState()
		}
		timer := NewTimer()

		err := timer.Pause()
		if err == nil {
			t.Error("expected error when pausing idle timer")
		}
	})

	t.Run("ResumeNotPaused", func(t *testing.T) {
		// Clear state before this test
		if stateManager, err := NewStateManager(); err == nil {
			stateManager.ClearState()
		}
		timer := NewTimer()

		err := timer.Resume()
		if err == nil {
			t.Error("expected error when resuming idle timer")
		}

		timer.Start(1 * time.Minute)
		err = timer.Resume()
		if err == nil {
			t.Error("expected error when resuming running timer")
		}
	})

	t.Run("StopNotRunning", func(t *testing.T) {
		// Clear state before this test
		if stateManager, err := NewStateManager(); err == nil {
			stateManager.ClearState()
		}
		timer := NewTimer()

		err := timer.Stop()
		if err == nil {
			t.Error("expected error when stopping idle timer")
		}
	})
}

func TestTimerProgress(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	t.Run("ProgressIdle", func(t *testing.T) {
		timer := NewTimer()
		progress := timer.GetProgress()
		if progress != 0 {
			t.Errorf("expected 0 progress for idle timer, got %f", progress)
		}
	})

	t.Run("ProgressRunning", func(t *testing.T) {
		timer := NewTimer()
		duration := 100 * time.Millisecond
		timer.Start(duration)

		// Wait a bit and check progress
		time.Sleep(50 * time.Millisecond)
		progress := timer.GetProgress()

		if progress <= 0 || progress >= 1 {
			t.Errorf("expected progress between 0 and 1, got %f", progress)
		}
	})

	t.Run("ProgressPaused", func(t *testing.T) {
		timer := NewTimer()
		duration := 100 * time.Millisecond
		timer.Start(duration)
		time.Sleep(50 * time.Millisecond)
		timer.Pause()

		progress1 := timer.GetProgress()
		time.Sleep(10 * time.Millisecond)
		progress2 := timer.GetProgress()

		if progress1 != progress2 {
			t.Errorf("progress should not change when paused: %f vs %f", progress1, progress2)
		}
	})

	t.Run("ProgressCompleted", func(t *testing.T) {
		timer := NewTimer()
		duration := 10 * time.Millisecond
		timer.Start(duration)
		time.Sleep(20 * time.Millisecond) // Wait longer than duration

		progress := timer.GetProgress()
		if progress != 1.0 {
			t.Errorf("expected 1.0 progress for completed timer, got %f", progress)
		}
	})
}

func TestNewTimer(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	timer := NewTimer()
	if timer.GetStatus() != StatusIdle {
		t.Errorf("new timer should be idle, got %v", timer.GetStatus())
	}
	if timer.GetProgress() != 0 {
		t.Errorf("new timer should have 0 progress, got %f", timer.GetProgress())
	}
}
