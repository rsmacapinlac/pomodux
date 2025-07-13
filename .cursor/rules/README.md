# Pomodux Cursor Rules - Consolidated Structure

## Overview

This directory contains the consolidated cursor rules for the Pomodux project. The rules have been reorganized into 4 focused files with consistent command patterns, maintaining all context and quality while improving usability.

## Consolidated Structure

### 1. **`release-process.mdc`** - Release Process and Development Workflow
**Purpose**: Manages the release process and development workflow with consistent command patterns

**Primary Commands**:
- `Design a new feature` - Plan and architect new features with proper validation
- `Organize a release` - Plan and organize features into a release
- `Implement a release` - Execute development with proper gates and quality controls
- `Execute a retrospective` - Analyze release outcomes and drive continuous improvement

**Consolidated From**:
- `development-process.mdc` - Development process and approval gates
- `retrospective.mdc` - Post-release analysis and improvement

### 2. **`architecture-standards.mdc`** - Architecture Guidelines and Decisions
**Purpose**: Manages architectural decisions and component selection

**Primary Commands**:
- `Review architecture decisions` - Audit existing architectural decisions and identify constraints
- `Propose architectural changes` - Create comprehensive architecture proposals for new features or changes
- `Validate component selection` - Evaluate third-party components and libraries for architectural fit

**Consolidated From**:
- `architecture.mdc` - Architecture decision records and guidelines
- `pre-development-architecture.mdc` - Pre-development architecture review process

### 3. **`code-standards.mdc`** - Go Development Standards
**Purpose**: Language-specific development standards and patterns

**Primary Commands**:
- `Apply Go standards` - Apply Go-specific development patterns and conventions
- `Validate code quality` - Ensure code meets quality standards and requirements
- `Review security patterns` - Ensure code follows security and performance best practices

**Consolidated From**:
- `go-development.mdc` - Go development standards and principles

### 4. **`project-management.mdc`** - Documentation and Project Management
**Purpose**: Manages documentation organization and project status tracking

**Primary Commands**:
- `Organize documentation` - Maintain proper documentation structure and organization
- `Update project status` - Maintain current project status and progress tracking

**Consolidated From**:
- `docs-folder.mdc` - Documentation organization rules

## Key Improvements

### Consistent Command Patterns
All rules now use consistent command patterns that are easy to remember and use:

- **Design** - For planning and ideation phases
- **Organize** - For structure and organization tasks
- **Implement** - For execution and development work
- **Execute** - For analysis and review processes
- **Review** - For validation and assessment activities
- **Validate** - For verification and testing tasks
- **Apply** - For standards and pattern application
- **Update** - For status and progress tracking

### Maintained Quality and Context
- **All original content preserved**: No information was lost during consolidation
- **Enhanced organization**: Related concepts are now grouped logically
- **Improved cross-references**: Clear integration between different rule areas
- **Consistent structure**: All files follow the same organizational pattern

### Reduced Complexity
- **From 7 files to 4**: 43% reduction in file count
- **Clearer responsibilities**: Each file has a focused, well-defined purpose
- **Easier navigation**: Fewer files to search through
- **Better integration**: Related concepts are co-located

## Usage Guidelines

### When to Use Each File

**Use `release-process.mdc` when**:
- Starting new features or releases
- Planning development work
- Executing development with gates
- Conducting retrospectives
- Managing the overall development lifecycle

**Use `architecture-standards.mdc` when**:
- Making architectural decisions
- Selecting third-party components
- Reviewing existing architecture
- Proposing architectural changes
- Validating component compatibility

**Use `code-standards.mdc` when**:
- Writing Go code
- Applying development patterns
- Validating code quality
- Reviewing security and performance
- Following Go conventions

**Use `project-management.mdc` when**:
- Organizing documentation
- Tracking project status
- Maintaining documentation structure
- Coordinating project workflow

### Command Usage Examples

```bash
# Planning phases
"Design a new feature for enhanced logging"
"Organize a release for version 0.3.1"

# Development phases  
"Implement a release following the 4-gate process"
"Apply Go standards to the timer engine"

# Review phases
"Review architecture decisions for the plugin system"
"Validate code quality for the new feature"

# Management phases
"Organize documentation for the current release"
"Update project status with current progress"

# Analysis phases
"Execute a retrospective for release 0.3.0"
```

## Integration with Development Workflow

### Typical Development Flow
1. **`Design a new feature`** (release-process.mdc)
2. **`Review architecture decisions`** (architecture-standards.mdc)
3. **`Propose architectural changes`** (architecture-standards.mdc)
4. **`Organize a release`** (release-process.mdc)
5. **`Organize documentation`** (project-management.mdc)
6. **`Implement a release`** (release-process.mdc)
7. **`Apply Go standards`** (code-standards.mdc)
8. **`Validate code quality`** (code-standards.mdc)
9. **`Update project status`** (project-management.mdc)
10. **`Execute a retrospective`** (release-process.mdc)

### Quality Assurance Integration
- All commands include quality checklists
- TDD requirements are integrated throughout
- Documentation standards are maintained
- Architecture compliance is enforced
- Process gates are clearly defined

## File Relationships

### How Files Work Together
- **`release-process.mdc`**: Handles development workflow and release execution
- **`project-management.mdc`**: Handles documentation organization and status tracking
- **`architecture-standards.mdc`**: Handles architectural decisions and component selection
- **`code-standards.mdc`**: Handles Go development patterns and quality standards

### Integration Points
- Release process uses project management for documentation organization
- Architecture standards inform both release process and code standards
- Code standards support release process quality requirements
- Project management tracks progress across all other areas

## Migration Notes

### What Changed
- **File consolidation**: 7 files → 4 files
- **Command standardization**: Consistent command patterns across all files
- **Content reorganization**: Related concepts grouped together
- **Cross-reference updates**: All internal references updated
- **New release-process.mdc**: Focused on release workflow and development process

### What Stayed the Same
- **All original content**: No information was lost
- **Quality standards**: All quality requirements maintained
- **Process requirements**: All approval gates and processes preserved
- **Documentation standards**: All documentation requirements maintained

### Benefits Achieved
- **Easier to use**: Consistent command patterns
- **Better organization**: Logical grouping of related concepts
- **Reduced complexity**: Fewer files to manage
- **Improved integration**: Clear relationships between different areas
- **Maintained quality**: All standards and requirements preserved

## References

### Project Documentation
- **[Release Process](mdc:.cursor/rules/release-process.mdc)** - Release workflow and development process
- **[Architecture Standards](mdc:.cursor/rules/architecture-standards.mdc)** - Architecture guidelines and decisions
- **[Code Standards](mdc:.cursor/rules/code-standards.mdc)** - Go development standards
- **[Project Management](mdc:.cursor/rules/project-management.mdc)** - Documentation and project management

### External Resources
- [Cursor Rules Documentation](mdc:https:/cursor.sh/docs/rules)
- [Cursor AI Best Practices](mdc:https:/cursor.sh/docs/ai)

---

**Note**: This consolidated structure maintains all the quality and context of the original rules while providing a more streamlined and consistent user experience. Use the primary commands to guide all development activities. 