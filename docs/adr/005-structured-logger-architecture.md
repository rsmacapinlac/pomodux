---
status: proposed
version: 0.1
created: 2025-07-13
type: technical
---

# ADR 005: Structured Logger Architecture for Pomodux

## 1. Context / Background

Pomodux currently uses ad-hoc logging via `fmt.Printf` and similar calls, with no log levels, structure, or configuration. As the application grows (core, CLI, TUI, plugin system), robust logging is essential for maintainability, debugging, and user support.

## 2. Motivation
- Inconsistent and unstructured logging makes debugging difficult
- No log levels (DEBUG/INFO/WARN/ERROR)
- No structured fields or context
- No configuration for log output, format, or verbosity
- No separation of user-facing output from diagnostics
- Need for better observability as Pomodux grows

## 3. Requirements
- **Log Levels:** DEBUG, INFO, WARN, ERROR
- **Structured Fields:** Attach context (component, plugin, event, etc.)
- **Configurable Output:** Console, file, or both
- **Format:** Human-readable text and machine-parsable JSON
- **Configurable via config file** (XDG-compliant)
- **Minimal Overhead:** No impact on timer accuracy or performance
- **Cross-Platform:** No OS-specific dependencies
- **Easy Integration:** Usable in all parts of the codebase

## 4. Decision

**Selected Solution:**
- Use a mature Go logging library (**logrus**)
- Add a new `internal/logger` package
- Expose a global logger instance with helper functions for each log level
- Add a `logging` section to the config file for user control
- Replace all `fmt.Printf`/`fmt.Fprintf` debug/info/warn/error calls with logger calls
- Provide a simple API for plugins (future work)

## 5. Alternatives Considered
- **Keep using fmt.Printf:** Not scalable, no log levels, no structure
- **Use Go’s standard log package:** Lacks structured fields and log levels
- **Write a custom logger:** Reinvents the wheel, more maintenance

## 6. Risks
- **Migration:** Need to update all existing debug/info/error output
- **User confusion:** Must clearly separate user-facing output from logs
- **Performance:** Must ensure logging does not block timer or CLI responsiveness

## 7. Proof-of-Concept Plan
- Implement a minimal logger using logrus
- Integrate with config for log level and output
- Replace a few key debug/info/error prints in core and CLI
- Validate log output in both text and JSON modes
- Test log file output and rotation (if needed)

## 8. Integration Plan
- Add logger initialization to main application startup
- Update config system to support logging section
- Gradually migrate all code to use the logger
- Document usage and configuration in the developer guide

## 9. Gate 0 Status
- **Approved** for implementation (2025-07-13)

---
**Last Updated:** 2025-07-13 