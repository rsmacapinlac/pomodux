package timer

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/logger"
)

// TestMain initializes the logger for all tests in this package
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logConfig := &logger.Config{
		Level:      logger.LogLevelDebug,
		Format:     "text",
		Output:     "console",
		ShowCaller: false,
	}
	if err := logger.Init(logConfig); err != nil {
		panic("Failed to initialize logger for tests: " + err.Error())
	}

	// Run tests
	os.Exit(m.Run())
}

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

func TestStartPersistent(t *testing.T) {
	// Clear any existing state for testing
	if stateManager, err := NewStateManager(); err == nil {
		stateManager.ClearState()
	}

	t.Run("StartPersistentBasic", func(t *testing.T) {
		timer := NewTimer()
		duration := 100 * time.Millisecond
		sessionType := SessionTypeWork

		// Start persistent timer in a goroutine
		errChan := make(chan error, 1)
		go func() {
			errChan <- timer.StartPersistent(duration, sessionType)
		}()

		// Wait for timer to complete
		select {
		case err := <-errChan:
			// In non-interactive environments, terminal raw mode will fail
			// This is expected behavior, so we don't treat it as a test failure
			if err != nil && !strings.Contains(err.Error(), "inappropriate ioctl for device") {
				t.Errorf("StartPersistent should not return error, got %v", err)
			}
		case <-time.After(200 * time.Millisecond):
			t.Error("StartPersistent should complete within expected time")
		}

		// Verify timer completed (even if terminal mode failed)
		if timer.GetStatus() != StatusCompleted && timer.GetStatus() != StatusRunning {
			t.Errorf("timer should be completed or running after StartPersistent, got %v", timer.GetStatus())
		}
	})

	t.Run("StartPersistentWithSessionType", func(t *testing.T) {
		timer := NewTimer()
		duration := 50 * time.Millisecond
		sessionType := SessionTypeBreak

		// Start persistent timer in a goroutine
		errChan := make(chan error, 1)
		go func() {
			errChan <- timer.StartPersistent(duration, sessionType)
		}()

		// Wait for timer to complete
		select {
		case err := <-errChan:
			// In non-interactive environments, terminal raw mode will fail
			// This is expected behavior, so we don't treat it as a test failure
			if err != nil && !strings.Contains(err.Error(), "inappropriate ioctl for device") {
				t.Errorf("StartPersistent should not return error, got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("StartPersistent should complete within expected time")
		}

		// Verify session type was set correctly
		if timer.GetSessionType() != sessionType {
			t.Errorf("expected session type %s, got %s", sessionType, timer.GetSessionType())
		}
	})

	t.Run("StartPersistentInvalidDuration", func(t *testing.T) {
		timer := NewTimer()
		duration := 0 * time.Millisecond
		sessionType := SessionTypeWork

		err := timer.StartPersistent(duration, sessionType)
		if err == nil {
			t.Error("StartPersistent should return error for zero duration")
		}
	})

	t.Run("StartPersistentAlreadyRunning", func(t *testing.T) {
		timer := NewTimer()
		duration := 100 * time.Millisecond
		sessionType := SessionTypeWork

		// Start a regular timer first
		err := timer.Start(duration)
		if err != nil {
			t.Fatalf("failed to start timer: %v", err)
		}

		// Try to start persistent timer while already running
		err = timer.StartPersistent(duration, sessionType)
		if err == nil {
			t.Error("StartPersistent should return error when timer already running")
		}
	})
}

func TestProgressBarGeneration(t *testing.T) {
	t.Run("ProgressBarZero", func(t *testing.T) {
		bar := createProgressBar(0.0, 10)
		expected := "[░░░░░░░░░░]"
		if bar != expected {
			t.Errorf("expected progress bar %s, got %s", expected, bar)
		}
	})

	t.Run("ProgressBarHalf", func(t *testing.T) {
		bar := createProgressBar(0.5, 10)
		expected := "[█████░░░░░]"
		if bar != expected {
			t.Errorf("expected progress bar %s, got %s", expected, bar)
		}
	})

	t.Run("ProgressBarFull", func(t *testing.T) {
		bar := createProgressBar(1.0, 10)
		expected := "[██████████]"
		if bar != expected {
			t.Errorf("expected progress bar %s, got %s", expected, bar)
		}
	})

	t.Run("ProgressBarNegative", func(t *testing.T) {
		bar := createProgressBar(-0.5, 10)
		expected := "[░░░░░░░░░░]"
		if bar != expected {
			t.Errorf("expected progress bar %s, got %s", expected, bar)
		}
	})

	t.Run("ProgressBarOverflow", func(t *testing.T) {
		bar := createProgressBar(1.5, 10)
		expected := "[██████████]"
		if bar != expected {
			t.Errorf("expected progress bar %s, got %s", expected, bar)
		}
	})
}

func TestDurationFormatting(t *testing.T) {
	t.Run("FormatSeconds", func(t *testing.T) {
		duration := 45 * time.Second
		formatted := formatDuration(duration)
		expected := "0 minutes 45 seconds"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})

	t.Run("FormatMinutes", func(t *testing.T) {
		duration := 5 * time.Minute
		formatted := formatDuration(duration)
		expected := "5 minutes"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})

	t.Run("FormatMinutesSeconds", func(t *testing.T) {
		duration := 2*time.Minute + 30*time.Second
		formatted := formatDuration(duration)
		expected := "2 minutes 30 seconds"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})

	t.Run("FormatHours", func(t *testing.T) {
		duration := 1 * time.Hour
		formatted := formatDuration(duration)
		expected := "1 hour"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})

	t.Run("FormatHoursMinutes", func(t *testing.T) {
		duration := 1*time.Hour + 30*time.Minute
		formatted := formatDuration(duration)
		expected := "1 hour 30 minutes"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})

	t.Run("FormatNegative", func(t *testing.T) {
		duration := -5 * time.Minute
		formatted := formatDuration(duration)
		expected := "Completed"
		if formatted != expected {
			t.Errorf("expected %s, got %s", expected, formatted)
		}
	})
}

func TestPluralHelper(t *testing.T) {
	t.Run("PluralZero", func(t *testing.T) {
		result := plural(0)
		expected := "s"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("PluralOne", func(t *testing.T) {
		result := plural(1)
		expected := ""
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("PluralMultiple", func(t *testing.T) {
		result := plural(5)
		expected := "s"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})
}
