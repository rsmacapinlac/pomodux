# Architecture Decision Records (ADR) Index

This directory contains Architecture Decision Records documenting key technical decisions made during Pomodux development. Each ADR captures the context, decision, rationale, and consequences of important architectural choices.

## 📁 Directory Structure

```
docs/adr/
├── README.md                           # This index file
├── 001-programming-language-selection.md  # Programming language choice
├── 002-persistent-timer-design.md         # Persistent timer architecture
├── 003-uat-testing-approach.md            # User acceptance testing strategy
├── 004-plugin-system-architecture.md      # Plugin system design
├── 005-structured-logger-architecture.md  # Structured logging approach
└── 006-lumberjack-log-rotation-library.md # Log rotation library selection
```

## 📋 ADR Index

### Core Architecture Decisions
- **[ADR 001](001-programming-language-selection.md)** - Programming Language Selection ✅ **APPROVED**
  - **Component**: Project Foundation
  - **Decision**: Go language selected for cross-platform CLI development
  - **Impact**: Establishes development environment and toolchain

- **[ADR 002](002-persistent-timer-design.md)** - Persistent Timer Design ✅ **APPROVED**
  - **Component**: Timer Engine
  - **Decision**: File-based persistence with JSON state storage
  - **Impact**: Enables timer state recovery across application restarts

- **[ADR 003](003-uat-testing-approach.md)** - UAT Testing Approach ✅ **APPROVED**
  - **Component**: Quality Assurance
  - **Decision**: Automated UAT with shell scripts and BATS framework
  - **Impact**: Ensures user-facing functionality works as expected

### System Architecture Decisions
- **[ADR 004](004-plugin-system-architecture.md)** - Plugin System Architecture ✅ **APPROVED**
  - **Component**: Plugin System
  - **Decision**: Lua-based plugin system with event hooks
  - **Impact**: Enables extensibility and community contributions

- **[ADR 005](005-structured-logger-architecture.md)** - Structured Logger Architecture 🔄 **PROPOSED**
  - **Component**: Logging System
  - **Decision**: Logrus-based structured logging with component filtering
  - **Impact**: Provides comprehensive logging for debugging and monitoring

- **[ADR 006](006-lumberjack-log-rotation-library.md)** - Lumberjack Log Rotation Library ✅ **APPROVED**
  - **Component**: Logging System
  - **Decision**: Lumberjack library for automatic log file rotation
  - **Impact**: Prevents log files from growing too large

## 🔗 Related Documentation

- **[Release Management](../release-management.md)** - ADR process and approval gates
- **[Technical Specifications](../technical_specifications.md)** - Technical architecture overview
- **[Documentation Standards](../documentation-standards.md)** - Documentation guidelines

## 📖 ADR Process

### When to Create an ADR
Create an ADR when making decisions about:
- **Technology Selection**: Programming languages, libraries, frameworks
- **Architecture Patterns**: System design, data flow, component interactions
- **Integration Approaches**: How components communicate and integrate
- **Quality Assurance**: Testing strategies, validation approaches
- **Operational Concerns**: Logging, monitoring, deployment strategies

### ADR Approval Process
1. **Gate 0 Review**: ADR presented during architecture review phase
2. **Stakeholder Approval**: Technical approach approved by stakeholders
3. **Implementation**: ADR guides development implementation
4. **Validation**: ADR decisions validated during release process

### ADR Status Indicators
- **✅ APPROVED** - Decision approved and implemented
- **🔄 PROPOSED** - Decision proposed, pending approval
- **❌ DEPRECATED** - Decision superseded or no longer relevant

## 📋 Maintenance Requirements

This index **MUST** be updated when:
- New ADRs are created
- ADR status changes (e.g., from PROPOSED to APPROVED)
- ADRs are deprecated or superseded
- ADR content is significantly modified

### Quality Checks
- [ ] All ADR files are listed in index
- [ ] All links are valid and current
- [ ] Status indicators reflect current state
- [ ] Component assignments are accurate
- [ ] Impact descriptions are clear and current

---

**Note**: This directory serves as the authoritative record of architectural decisions. Each ADR provides context for understanding why specific technical choices were made and their implications for the project.

**Last Updated:** 2025-01-27 