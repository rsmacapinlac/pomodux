package timer

// EventType represents the type of timer event.
type EventType string

const (
	EventTimerStarted   EventType = "timer_started"
	EventTimerPaused    EventType = "timer_paused"
	EventTimerResumed   EventType = "timer_resumed"
	EventTimerCompleted EventType = "timer_completed"
	EventTimerStopped   EventType = "timer_stopped"
)
