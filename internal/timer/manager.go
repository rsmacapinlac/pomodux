package timer

import (
	"sync"
	"time"
)

// Manager provides a singleton timer instance for the application
type Manager struct {
	timer *Timer
	mu    sync.RWMutex
}

var (
	manager *Manager
	once    sync.Once
)

// GetManager returns the singleton timer manager instance
func GetManager() *Manager {
	once.Do(func() {
		manager = &Manager{
			timer: NewTimer(),
		}
	})
	return manager
}

// GetTimer returns the current timer instance
func (m *Manager) GetTimer() *Timer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timer
}

// StartTimer starts the timer with the given duration
func (m *Manager) StartTimer(duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timer.Start(duration)
}

// StopTimer stops the current timer
func (m *Manager) StopTimer() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timer.Stop()
}

// PauseTimer pauses the current timer
func (m *Manager) PauseTimer() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timer.Pause()
}

// ResumeTimer resumes the current timer
func (m *Manager) ResumeTimer() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timer.Resume()
}

// GetStatus returns the current timer status
func (m *Manager) GetStatus() TimerStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timer.GetStatus()
}

// GetProgress returns the current timer progress
func (m *Manager) GetProgress() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timer.GetProgress()
}
