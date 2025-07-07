package timer

import "time"

// TimerEngine defines the core timer interface
// for Pomodux timer operations.
type TimerEngine interface {
	Start(duration time.Duration) error
	Stop() error
	Pause() error
	Resume() error
	GetStatus() TimerStatus
	GetProgress() float64
	GetSessionType() SessionType
}

// TimerStatus represents the state of the timer.
type TimerStatus string

const (
	StatusIdle      TimerStatus = "idle"
	StatusRunning   TimerStatus = "running"
	StatusPaused    TimerStatus = "paused"
	StatusCompleted TimerStatus = "completed"
)

// SessionType represents the type of timer session.
type SessionType string

const (
	SessionTypeWork      SessionType = "work"
	SessionTypeBreak     SessionType = "break"
	SessionTypeLongBreak SessionType = "long-break"
)
