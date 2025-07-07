# Pomodux Technical Specifications

> **Note:** The implementation language for this project is Go, as approved in [ADR 001](adr/001-programming-language-selection.md).

## 1.0 Architecture Overview

### 1.1 System Architecture

The Pomodux terminal timer application follows a modular architecture with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │   TUI Layer     │    │   Core Engine   │
│                 │    │                 │    │                 │
│ - Argument      │    │ - Terminal UI   │    │ - Timer Logic   │
│   Parsing       │    │ - Event Loop    │    │ - Session Mgmt  │
│ - Command       │    │ - Input/Output  │    │ - State Mgmt    │
│   Execution     │    │ - Rendering     │    │ - Notifications │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Plugin System │
                    │                 │
                    │ - Lua Runtime   │
                    │ - Hook System   │
                    │ - API Interface │
                    └─────────────────┘
```

### 1.2 Core Components

#### Timer Engine
- **Purpose**: Core timing logic and session management
- **Responsibilities**:
  - Timer state management
  - Session tracking
  - Progress calculation
  - Persistent timer sessions
  - Interactive keypress controls
  - Automatic session recording
  - Real-time progress display

#### TUI Renderer
- **Purpose**: Terminal user interface rendering
- **Responsibilities**:
  - Screen layout management
  - Progress bar rendering
  - Input handling
  - Terminal resize handling

#### CLI Interface
- **Purpose**: Command-line argument processing
- **Responsibilities**:
  - Argument parsing
  - Command routing
  - Tab completion
  - Help system

#### Plugin System
- **Purpose**: Extensibility through Lua scripts
- **Responsibilities**:
  - Lua runtime management
  - Hook system
  - Plugin API
  - Plugin lifecycle management

## 2.0 Component Specifications

### 2.1 Timer Engine

#### 2.1.1 Core Timer Logic
- **State Management**: Idle, Running, Paused, Completed states
- **Duration Handling**: Support for various time formats and durations
- **Progress Tracking**: Real-time progress calculation and updates
- **Persistent Sessions**: Timer sessions that run continuously until user interaction
- **Interactive Controls**: Keypress-based control system (p, r, q, s, Ctrl+C)
- **Automatic Recording**: Sessions automatically recorded upon completion
- **Real-time Display**: Continuous progress updates with progress bars

#### 2.1.2 Session Management
- **Session Types**: Work sessions, break sessions, long break sessions
- **Session History**: Persistent storage of completed sessions with automatic recording
- **Statistics**: Session duration, completion rates, productivity metrics
- **Session Configuration**: Customizable session parameters
- **Automatic Recording**: Sessions automatically recorded upon completion
- **Session State**: Persistent timer state across process restarts
- **Session Recorder**: Dedicated component for handling session recording and history management

### 2.2 TUI Renderer

#### 2.2.1 Display Components
- **Progress Visualization**: Real-time progress bars and indicators
- **Status Display**: Current timer state and session information
- **Menu System**: Interactive navigation and command selection
- **Theme System**: Configurable colors, fonts, and layouts

#### 2.2.2 Input Handling
- **Keyboard Events**: Handle all keyboard input and shortcuts
- **Resize Events**: Adapt layout to terminal size changes
- **Focus Management**: Navigate between UI components
- **Event Loop**: Manage UI updates and user interactions

### 2.3 CLI Interface

#### 2.3.1 Command Structure
- **Root Commands**: Main application commands (start, break, long-break)
- **Subcommands**: Specialized operations (config, history, stats)
- **Flags and Options**: Command-line parameters and configuration
- **Help System**: Comprehensive help and usage documentation

#### 2.3.2 Integration Features
- **Shell Integration**: Tab completion and command history
- **Script Support**: Command-line automation and scripting
- **Output Formats**: Text, JSON, and structured output options
- **Error Handling**: Clear error messages and exit codes

### 2.4 Plugin System

#### 2.4.1 Lua Runtime
- **Sandboxed Environment**: Secure plugin execution
- **API Access**: Controlled access to application features
- **Resource Management**: Memory and CPU usage limits
- **Error Isolation**: Plugin errors don't affect main application

#### 2.4.2 Hook System
- **Event Types**: Timer events, UI events, system events
- **Hook Registration**: Plugin hook registration and management
- **Event Propagation**: Reliable event delivery to plugins
- **Hook Lifecycle**: Hook loading, execution, and cleanup

## 3.0 Data Flow Architecture

### 3.1 Persistent Timer Flow
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

### 3.2 Configuration Management
- **Configuration Sources**: File-based, environment variables, command-line
- **Configuration Hierarchy**: Default → User → Command-line precedence
- **Configuration Validation**: Schema validation and error reporting
- **Configuration Persistence**: Reliable saving and loading of settings

### 3.3 Persistent Timer System
- **Session Persistence**: Timer state maintained across process restarts
- **Interactive Controls**: Keypress-based control system during sessions
- **Real-time Display**: Continuous progress updates with progress bars
- **Session Recording**: Automatic history recording on timer completion
- **State Management**: Persistent timer state with background monitoring
- **Keypress Handling**: Raw terminal mode for reliable input handling

### 3.4 State Management
- **Application State**: Global application state and settings
- **Session State**: Current timer and session information
- **UI State**: Interface state and user preferences
- **Plugin State**: Plugin-specific state and configuration
- **Timer State**: Persistent timer state with interactive controls
- **History State**: Session history and statistics management

## 4.0 Keypress Control System

### 4.1 Control Mapping
- **'p'**: Pause current timer session
- **'r'**: Resume paused timer session
- **'q'**: Stop timer session and exit
- **'s'**: Stop timer session and exit (alternative)
- **Ctrl+C**: Emergency exit from timer session

### 4.2 Input Handling
- **Raw Terminal Mode**: Terminal put into raw mode for reliable keypress detection
- **Non-blocking Input**: Keypress handling doesn't block timer updates
- **Graceful Cleanup**: Terminal state restored on exit
- **Cross-platform**: Works across different terminal environments

### 4.3 Session Flow
1. **Start**: User starts timer with `pomodux start [duration]`
2. **Persistent Session**: Timer runs continuously with real-time display
3. **Interactive Control**: User can pause, resume, or stop via keypresses
4. **Completion**: Timer completes automatically or user stops manually
5. **Recording**: Session automatically recorded to history
6. **Exit**: User exits session and returns to command line

## 5.0 Session Types

### 5.1 Work Sessions
- **Command**: `pomodux start [duration]`
- **Default Duration**: 25 minutes (configurable)
- **Purpose**: Focused work periods
- **Recording**: Automatically recorded on completion

### 5.2 Break Sessions
- **Command**: `pomodux break`
- **Default Duration**: 5 minutes (configurable)
- **Purpose**: Short breaks between work sessions
- **Recording**: Automatically recorded on completion

### 5.3 Long Break Sessions
- **Command**: `pomodux long-break`
- **Default Duration**: 15 minutes (configurable)
- **Purpose**: Extended breaks after multiple work sessions
- **Recording**: Automatically recorded on completion

## 6.0 Data Persistence

### 6.1 Session History
- **Storage**: JSON file in XDG state directory
- **Format**: Structured session records with metadata
- **Access**: Via `pomodux history` command
- **Statistics**: Automatic calculation of session metrics

### 6.2 Timer State
- **Storage**: JSON file in XDG state directory
- **Purpose**: Persist timer state across process restarts
- **Recovery**: Automatic state restoration on startup
- **Cleanup**: State cleared when timer completes or stops

### 6.3 Configuration
- **Storage**: YAML file in XDG config directory
- **Validation**: Schema-based configuration validation
- **Defaults**: Sensible defaults for all settings
- **Migration**: Automatic configuration migration between versions

 