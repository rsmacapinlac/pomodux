---
status: approved
type: technical
---

# Persistent Timer with Keypress Controls

## 1. Context / Background

### 1.1 Current State
The initial timer implementation used separate CLI commands for timer control (`pause`, `resume`, `stop`, `status`). This approach had several limitations:

- Users needed to remember and type separate commands for timer control
- Timer state was only checked when commands were executed
- No real-time progress display or interactive feedback
- Session recording only happened when manually stopped
- Poor user experience with multiple command invocations

### 1.2 Requirements
- Interactive timer control without leaving the session
- Real-time progress display with visual feedback
- Automatic session recording on completion
- Intuitive user interface for timer management
- Persistent timer state across process restarts
- Simplified command structure

## 2. Decision

**Selected Solution:** Persistent timer sessions with interactive keypress controls

### 2.1 Rationale
The persistent timer approach provides a superior user experience by keeping users in an interactive session with real-time feedback. Keypress controls are more intuitive than separate commands, and the persistent session ensures automatic session recording and state management. This design simplifies the command structure while providing richer functionality.

## 3. Solutions Considered

### 3.1 Option A: Separate CLI Commands (Original)
- **Pros:**
  - Simple command structure
  - Familiar CLI pattern
  - Easy to implement
- **Cons:**
  - Poor user experience
  - No real-time feedback
  - Manual session recording
  - Multiple command invocations required
  - No visual progress indication

### 3.2 Option B: Persistent Timer with Keypress Controls (Selected)
- **Pros:**
  - Superior user experience
  - Real-time progress display
  - Interactive controls
  - Automatic session recording
  - Visual feedback with progress bars
  - Simplified command structure
  - Persistent state management
- **Cons:**
  - More complex implementation
  - Requires raw terminal mode
  - Platform-specific terminal handling
  - More complex testing requirements

### 3.3 Option C: Real-time Event-Driven Engine (Planned but not implemented)
- **Pros:**
  - Event-driven architecture
  - Background monitoring
  - Plugin system foundation
- **Cons:**
  - Over-engineered for current needs
  - Complex implementation
  - Higher resource usage
  - Not aligned with user experience goals

## 4. Consequences

### 4.1 Positive
- Excellent user experience with interactive controls
- Real-time progress display provides immediate feedback
- Automatic session recording ensures no lost data
- Simplified command structure (only start, break, long-break)
- Persistent state across process restarts
- Visual progress bars and time remaining display
- Intuitive keypress controls (p, r, q, s, Ctrl+C)

### 4.2 Negative
- More complex implementation with raw terminal mode
- Platform-specific terminal handling required
- Testing interactive features is more challenging
- Requires proper terminal state cleanup
- More complex error handling for terminal operations

### 4.3 Risks
- **Terminal Compatibility**: Mitigated through proper raw mode handling and graceful fallbacks
- **Platform Differences**: Mitigated through cross-platform terminal library usage
- **State Cleanup**: Mitigated through proper defer statements and signal handling
- **User Learning Curve**: Mitigated through clear on-screen instructions and help system

## 5. Component Information

### 5.1 Repository Links
- **GitHub Repository**: https://github.com/ritchie46/pomodux
- **Documentation**: [Internal timer engine documentation]
- **Release Notes**: [Release 0.2.1 notes]

### 5.2 Maintenance Status
- **Last Updated**: 2025-01-27
- **Active Development**: Yes - Part of Release 0.2.1
- **Community Support**: N/A (Internal component)
- **Version Compatibility**: Compatible with existing CLI commands

### 5.3 Integration Verification
- **Compatibility Tested**: Yes - Maintains existing CLI API
- **Existing Component Impact**: Minimal - Timer interface remains unchanged
- **Migration Path**: Seamless - No user-facing changes required

## 6. Technical Design

### 6.1 Core Components

#### Persistent Timer Session
- **Session Management**: Continuous timer session until user interaction or completion
- **State Persistence**: Timer state maintained across process restarts
- **Progress Display**: Real-time progress bars and time remaining
- **Session Recording**: Automatic history recording on completion

#### Keypress Control System
- **Control Mapping**: 'p' (pause), 'r' (resume), 'q'/'s' (stop), Ctrl+C (exit)
- **Raw Terminal Mode**: Terminal put into raw mode for reliable keypress detection
- **Non-blocking Input**: Keypress handling doesn't block timer updates
- **Graceful Cleanup**: Terminal state restored on exit

#### Real-time Display
- **Progress Bar**: Visual progress indication with Unicode block characters
- **Time Remaining**: Formatted time display (hours:minutes:seconds)
- **Session Type**: Display current session type (work, break, long-break)
- **Control Instructions**: On-screen help for available controls

### 6.2 Architecture Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Timer Start │───▶│ Persistent  │───▶│ Keypress    │
│ Command     │    │ Session     │    │ Handler     │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                           ▼                   ▼
                    ┌─────────────┐    ┌─────────────┐
                    │ Progress    │    │ Session     │
                    │ Display     │    │ Recording   │
                    └─────────────┘    └─────────────┘
                              │                   │
                              ▼                   ▼
                       ┌─────────────┐    ┌─────────────┐
                       │ Real-time   │    │ History     │
                       │ Updates     │    │ Persistence │
                       └─────────────┘    └─────────────┘
```

### 6.3 Implementation Details

#### Persistent Timer Session
```go
// Persistent timer session
type PersistentTimer struct {
    timer    *Timer
    session  *Session
    display  *ProgressDisplay
    controls *KeypressHandler
}

func (pt *PersistentTimer) Start() error {
    // Start persistent session with keypress controls
    // Real-time progress display
    // Automatic session recording
}
```

#### Keypress Controls
```go
// Keypress control mapping
const (
    KeyPause  = 'p'
    KeyResume = 'r'
    KeyStop   = 'q'
    KeyStopAlt = 's'
    KeyExit   = 3 // Ctrl+C
)

// Keypress handler
type KeypressHandler struct {
    timer    *Timer
    done     chan bool
    controls map[rune]func()
}
```

#### Session Flow
1. **Start**: User starts timer with `pomodux start [duration]`
2. **Persistent Session**: Timer runs continuously with real-time display
3. **Interactive Control**: User can pause, resume, or stop via keypresses
4. **Completion**: Timer completes automatically or user stops manually
5. **Recording**: Session automatically recorded to history
6. **Exit**: User exits session and returns to command line

### 6.4 Backward Compatibility
- All existing CLI commands continue to work
- Timer interface remains unchanged
- Configuration system unaffected
- State file format unchanged
- No breaking changes to user workflow

### 6.5 Testing Strategy
- Unit tests for timer logic
- Integration tests for keypress handling
- Terminal compatibility tests
- State persistence tests
- Progress display tests
- Session recording tests

## 7. User Experience

### 7.1 Interactive Controls
- **Pause/Resume**: 'p' to pause, 'r' to resume
- **Stop**: 'q' or 's' to stop and exit
- **Emergency Exit**: Ctrl+C for immediate exit
- **Visual Feedback**: Clear indication of current state

### 7.2 Progress Display
- **Progress Bar**: Visual progress with Unicode blocks
- **Time Remaining**: Formatted time display
- **Session Type**: Clear indication of work/break session
- **Instructions**: On-screen help for controls

### 7.3 Session Management
- **Automatic Recording**: Sessions recorded on completion
- **State Persistence**: Timer state maintained across restarts
- **History Integration**: Seamless integration with history system
- **Statistics**: Automatic calculation of session metrics

## 8. Future Considerations

### 8.1 Potential Enhancements
- Sound notifications on completion
- Desktop notifications integration
- Advanced progress visualization
- Custom keypress mappings
- Session templates and presets

### 8.2 Plugin System Integration
- Event hooks for timer state changes
- Plugin API for timer control
- Custom display components
- Enhanced notification system

### 8.3 Performance Optimization
- Memory usage optimization for long sessions
- CPU usage optimization for real-time display
- Cross-platform performance tuning
- Terminal compatibility improvements 