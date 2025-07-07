package timer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

// Timer represents a timer instance
type Timer struct {
	mu             sync.Mutex
	status         TimerStatus
	sessionType    SessionType
	startTime      time.Time
	duration       time.Duration
	elapsed        time.Duration
	stateManager   *StateManager
	historyManager *HistoryManager
}

// NewTimer creates a new Timer instance.
func NewTimer() *Timer {
	return &Timer{
		status: StatusIdle,
	}
}

// NewTimerWithManagers creates a new Timer instance with state and history managers.
func NewTimerWithManagers(stateManager *StateManager, historyManager *HistoryManager) *Timer {
	timer := &Timer{
		status:         StatusIdle,
		stateManager:   stateManager,
		historyManager: historyManager,
	}

	// Load existing state if available
	if state, err := stateManager.LoadState(); err == nil {
		timer.status = state.Status
		timer.sessionType = state.SessionType
		timer.duration = state.Duration
		timer.startTime = state.StartTime
		timer.elapsed = state.Elapsed
	}

	return timer
}

// Start begins the timer for the specified duration.
func (t *Timer) Start(duration time.Duration) error {
	return t.StartWithType(duration, SessionTypeWork)
}

// StartWithType begins the timer for the specified duration and session type.
func (t *Timer) StartWithType(duration time.Duration, sessionType SessionType) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status == StatusRunning {
		return fmt.Errorf("timer already running")
	}
	if duration <= 0 {
		return fmt.Errorf("invalid duration")
	}
	t.duration = duration
	t.sessionType = sessionType
	t.startTime = time.Now()
	t.elapsed = 0
	t.status = StatusRunning

	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}

	return nil
}

// Stop stops the timer.
func (t *Timer) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.status == StatusIdle {
		return fmt.Errorf("timer not running")
	}

	// Record session in history if we have a history manager
	if t.historyManager != nil && t.sessionType != "" {
		session := SessionRecord{
			Type:      t.sessionType,
			Duration:  t.duration,
			StartTime: t.startTime,
			EndTime:   time.Now(),
			Completed: t.status == StatusCompleted,
		}
		// Don't check error here as history recording shouldn't prevent stopping
		t.historyManager.AddSession(session)
	}

	t.status = StatusIdle
	t.elapsed = 0
	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}
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
			// Record session in history immediately when timer completes
			if t.historyManager != nil && t.sessionType != "" {
				session := SessionRecord{
					Type:      t.sessionType,
					Duration:  t.duration,
					StartTime: t.startTime,
					EndTime:   time.Now(),
					Completed: true,
				}
				// Record session immediately
				t.historyManager.AddSession(session)
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

	// Record session in history if we have a history manager
	if t.historyManager != nil && t.sessionType != "" {
		session := SessionRecord{
			Type:      t.sessionType,
			Duration:  t.duration,
			StartTime: t.startTime,
			EndTime:   time.Now(),
			Completed: true,
		}
		// Don't check error here as history recording shouldn't prevent resetting
		t.historyManager.AddSession(session)
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

// GetSessionType returns the session type of the timer
func (t *Timer) GetSessionType() SessionType {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.sessionType
}

// GetStartTime returns the start time of the timer
func (t *Timer) GetStartTime() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.startTime
}

// SetHistoryManager sets the history manager for the timer.
func (t *Timer) SetHistoryManager(historyManager *HistoryManager) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.historyManager = historyManager
}

// StartPersistent starts a timer and blocks until completion, with live progress and keypress controls.
func (t *Timer) StartPersistent(duration time.Duration, sessionType SessionType) error {
	if err := t.StartWithType(duration, sessionType); err != nil {
		return err
	}

	fmt.Printf("Timer started for %v\n", duration)
	fmt.Printf("Session type: %s\n", sessionType)
	fmt.Println("Press 'p' to pause, 'r' to resume, 'q'/'s' to stop, Ctrl+C to exit.")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set terminal raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopChan := make(chan struct{}, 1)
	pauseChan := make(chan struct{}, 1)
	resumeChan := make(chan struct{}, 1)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Remove stdinR, only use stdinW for unblocking
	_, stdinW, _ := os.Pipe()
	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			b, err := r.ReadByte()
			if err != nil {
				continue
			}
			if b == 3 { // Ctrl+C
				stopChan <- struct{}{}
				return
			}
			switch b {
			case 'p':
				pauseChan <- struct{}{}
			case 'r':
				resumeChan <- struct{}{}
			case 'q', 's':
				stopChan <- struct{}{}
			}
		}
	}()

	go func() {
		<-interruptChan
		cancel()
		stdinW.Close() // Unblock keypress goroutine if blocked
	}()

	startTime := time.Now()
	endTime := startTime.Add(duration)
	paused := false
	var pauseStart time.Time
	var totalPaused time.Duration

	for {
		select {
		case <-ctx.Done():
			fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
			fmt.Println("Timer stopped by user (signal).")
			t.Stop()
			return nil
		case <-stopChan:
			fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
			fmt.Println("Timer stopped.")
			t.Stop()
			return nil
		case <-pauseChan:
			if !paused {
				pauseStart = time.Now()
				paused = true
				t.Pause()
				fmt.Print("\r⏸️  PAUSED - Press 'r' to resume" + strings.Repeat(" ", 50))
			}
		case <-resumeChan:
			if paused {
				t.Resume()
				totalPaused += time.Since(pauseStart)
				paused = false
				fmt.Print("\r▶️  RESUMED" + strings.Repeat(" ", 50))
			}
		default:
			if !paused {
				remaining := endTime.Add(totalPaused).Sub(time.Now())
				if remaining < 0 {
					remaining = 0
				}
				// Calculate progress
				elapsed := duration - remaining
				progress := elapsed.Seconds() / duration.Seconds()
				if progress > 1 {
					progress = 1
				}
				// Create progress bar
				progressBar := createProgressBar(progress, 30)
				percentage := int(progress * 100)
				// Display timer with progress bar
				fmt.Printf("\r%s %3d%% %s | %s",
					progressBar,
					percentage,
					formatDuration(remaining),
					sessionType)
				if remaining <= 0 {
					break
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
		if !paused && time.Now().After(endTime.Add(totalPaused)) {
			break
		}
	}

	fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
	t.handleCompletion()
	return nil
}

// handleCompletion records the session and sends notifications when timer completes.
func (t *Timer) handleCompletion() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Update status to completed
	t.status = StatusCompleted
	t.elapsed = t.duration

	// Save state
	if t.stateManager != nil {
		t.stateManager.SaveState(t)
	}

	// Record session in history immediately
	if t.historyManager != nil && t.sessionType != "" {
		session := SessionRecord{
			Type:      t.sessionType,
			Duration:  t.duration,
			StartTime: t.startTime,
			EndTime:   time.Now(),
			Completed: true,
		}
		t.historyManager.AddSession(session)
	}

	// Send notification
	sendNotification(t.sessionType, "completed")

	fmt.Print("Timer completed! Session recorded.")
}

// sendNotification sends a system notification based on session type.
func sendNotification(sessionType SessionType, status string) {
	var title, message string

	switch sessionType {
	case SessionTypeWork:
		title = "Work Session Complete"
		message = "Your work session has finished. Time for a break!"
	case SessionTypeBreak:
		title = "Break Complete"
		message = "Break time is over. Ready to work?"
	case SessionTypeLongBreak:
		title = "Long Break Complete"
		message = "Long break finished. Time to get back to work!"
	default:
		title = "Timer Complete"
		message = "Your timer session has finished."
	}

	// Try to send system notification
	cmd := exec.Command("notify-send", title, message)
	cmd.Run() // Ignore errors - notification is optional
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "Completed"
	}
	if d >= time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%d hour%s", hours, plural(hours))
		}
		return fmt.Sprintf("%d hour%s %d minute%s", hours, plural(hours), minutes, plural(minutes))
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if seconds == 0 {
		return fmt.Sprintf("%d minute%s", minutes, plural(minutes))
	}
	return fmt.Sprintf("%d minute%s %d second%s", minutes, plural(minutes), seconds, plural(seconds))
}

// createProgressBar creates a visual progress bar
func createProgressBar(progress float64, width int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(float64(width) * progress)
	empty := width - filled

	bar := "["
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	bar += "]"

	return bar
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
