---
status: approved
type: technical
---

# ADR 004: Plugin System Architecture for Pomodux

## 1. Context / Background

### 1.1 Motivation
Pomodux aims to be a highly extensible terminal timer application. Users and developers have requested the ability to extend Pomodux with custom features, notifications, integrations, and analytics. A plugin system is needed to:
- Enable user and community-driven feature development
- Allow for rapid prototyping and experimentation
- Support integrations with external tools and services
- Provide advanced customization (notifications, statistics, themes, etc.)

### 1.2 Requirements
- **Extensibility:** Allow third-party plugins to add/modify behavior
- **Security:** Sandbox plugin execution to prevent malicious actions
- **Performance:** Minimal impact on timer accuracy and resource usage
- **Cross-Platform:** Work on Linux, macOS, and Windows
- **Event-Driven:** Plugins should react to timer/session events
- **Ease of Use:** Simple API for plugin authors

## 2. Decision

**Selected Solution:** Lua-based plugin system using [gopher-lua](https://github.com/yuin/gopher-lua)

### 2.1 Rationale
- **Lua is lightweight, embeddable, and widely used for plugin systems** (e.g., Neovim)
- **gopher-lua** is a pure Go implementation, ensuring easy integration and cross-platform support
- **Event-driven architecture** allows plugins to subscribe to timer lifecycle events (start, pause, resume, complete, stop)
- **Sandboxing** is feasible with Lua, limiting plugin access to only the provided API
- **Proven approach:** Many successful projects use Lua for extensibility

### 2.2 Key Design Points
- **Plugin Lifecycle:** Plugins are loaded at startup or on demand, and can be enabled/disabled/unloaded at runtime
- **Event Hooks:** Plugins register Lua functions for specific timer events (e.g., `timer_started`, `timer_completed`)
- **API Surface:** Plugins can access timer/session data, send notifications, and log output
- **Isolation:** Each plugin runs in its own Lua state
- **Performance:** Event dispatch is asynchronous to avoid blocking the main timer loop

## 3. Alternatives Considered

### 3.1 No Plugin System
- **Pros:** Simpler codebase, less maintenance
- **Cons:** No extensibility, all features must be built-in, less community involvement

### 3.2 Go Plugins (native Go plugin package)
- **Pros:** Native performance, type safety
- **Cons:** Not cross-platform, complex build/distribution, unsafe for untrusted code

### 3.3 Scripting with Python/JS
- **Pros:** Popular languages, rich ecosystems
- **Cons:** Heavyweight, more dependencies, harder to sandbox, larger attack surface

### 3.4 External Process Plugins (IPC)
- **Pros:** Strong isolation, language agnostic
- **Cons:** Complex IPC, higher resource usage, more difficult for simple use cases

## 4. Consequences

### 4.1 Positive
- **Extensible:** Users can add new features without modifying core code
- **Safe:** Sandboxed Lua environment limits plugin capabilities
- **Cross-Platform:** Works on all supported OSes
- **Community-Friendly:** Low barrier for plugin authors
- **Maintainable:** Clear separation between core and extensions

### 4.2 Negative
- **Learning Curve:** Plugin authors must learn Lua
- **Debugging:** Errors in plugins may be harder to trace
- **Resource Usage:** Each plugin has its own Lua state (minimal, but nonzero overhead)
- **API Stability:** Need to maintain a stable plugin API

### 4.3 Risks
- **Security:** Plugins could attempt to escape sandbox (mitigated by strict API)
- **Performance:** Poorly written plugins could impact performance (mitigated by async event dispatch)
- **Complexity:** Plugin manager adds architectural complexity

## 5. Implementation and Validation Results

The plugin system has been successfully implemented and integrated into Pomodux:

### 5.1 Core Implementation
- **PluginManager** loads plugins from configurable directory, dispatches events, and manages plugin lifecycle
- **Event System** supports 5 timer events: `timer_started`, `timer_completed`, `timer_paused`, `timer_resumed`, `timer_stopped`
- **Lua Integration** using gopher-lua with proper sandboxing and error handling
- **Plugin API** provides registration, hook management, logging, and notification capabilities

### 5.2 Example Plugins
- **Mako Notification Plugin**: Sends desktop notifications using `notify-send` for all timer events
- **Debug Events Plugin**: Prints debug messages for all timer events to verify system functionality
- **Plugin Architecture**: Self-registering plugins with event hooks and configuration

### 5.3 Validation Results
- **Event Dispatch**: All 5 timer events properly emitted and processed by plugins
- **Plugin Loading**: Plugins load from `~/.config/pomodux/plugins/` directory
- **Hook Execution**: Lua hooks execute correctly with proper event data
- **Concurrency**: Event processing is asynchronous and non-blocking
- **Error Handling**: Plugin errors are caught and logged without affecting timer operation

### 5.4 Performance Characteristics
- **Event Dispatch**: Non-blocking with minimal latency
- **Plugin Loading**: Fast startup with lazy loading
- **Memory Usage**: Each plugin has isolated Lua state with minimal overhead
- **Timer Accuracy**: No measurable impact on timer precision

### 5.5 Integration Status
- **Timer Engine**: Fully integrated with event emission on all lifecycle changes
- **CLI Commands**: Plugin system available in all timer operations
- **Configuration**: Plugin directory configurable via XDG-compliant config system

## 6. Integration Plan
- **Timer Engine:** ✅ Fully integrated with event emission on all lifecycle changes
- **CLI:** ✅ Plugin system available in all timer operations
- **Documentation:** ✅ Plugin authoring guide and API reference available
- **Configuration:** ✅ XDG-compliant plugin directory configuration

## 7. Status
- **✅ Implemented and Validated** in Release 0.3.0
- **✅ All 5 timer events working correctly**
- **✅ Example plugins (notification, debug) functional**
- **✅ Plugin system ready for community use**
- **See also:**
  - [Technical Specifications](../technical_specifications.md)
  - [Current Release Documentation](../current-release.md)
  - [Implementation Roadmap](../implementation-roadmap.md)

## 8. References
- [gopher-lua](https://github.com/yuin/gopher-lua)
- [Neovim Plugin System](https://neovim.io/)
- [Pomodux Plugin System Implementation](../internal/plugin/)
- [Example Plugins](../plugins/)

## 9. Lessons Learned

### 9.1 Implementation Insights
- **Event Timing**: Added small delay in timer completion to ensure event processing completes before application exit
- **Plugin Registration Order**: Hooks must be registered before plugin registration for proper association
- **Debug Output**: Comprehensive debug logging essential for troubleshooting plugin system issues
- **Concurrency**: Asynchronous event processing prevents blocking of main timer loop

### 9.2 Validation Approach
- **Debug Plugin**: Created comprehensive debug plugin to verify all 5 timer events
- **Real-world Testing**: Tested with actual notification plugin to validate practical use cases
- **Event Coverage**: Verified all timer lifecycle events (start, pause, resume, stop, complete) work correctly

## Related Documentation

- [Release Management](../release-management.md)
- [Backlog](../backlog/)

---
**Last Updated:** 2025-07-13 