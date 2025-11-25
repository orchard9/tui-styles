# Tui Styles CLI - Quick Reference

## Installation

```bash
# Local install (recommended, no sudo)
cd tools/tui-styles-cli
make install

# System-wide install (requires sudo)
make install-system

# Or just build
make build  # Binary in bin/tui-styles-cli

# Using install script
./install.sh         # System-wide (sudo)
./install.sh --local # Local (no sudo)
```

## Common Commands

### Task Management

```bash
# Get task info
tui-styles-cli task get 001
tui-styles-cli task get roadmap/milestone-2/phase-1/001_task_pending.md

# Update task status
tui-styles-cli task update 001 --status ready
tui-styles-cli task update 002 --status in_progress
tui-styles-cli task update 003 --status blocked
tui-styles-cli task update 004 --status complete
```

### Comments

```bash
# Create a flagged comment (needs addressing)
tui-styles-cli comment create 001 "## BLOCKER: Decision needed

**Issue**: Library choice unclear

**Options**:
1. Option A
2. Option B" --author project-preplanner --needs-addressing

# Create a regular comment
tui-styles-cli comment create 002 "Decision approved" --author human

# List all comments on a task
tui-styles-cli comment list 001
```

### Milestones

```bash
# List all milestones
tui-styles-cli milestone list

# Get milestone info
tui-styles-cli milestone info milestone-2
tui-styles-cli milestone info 2  # Shorthand
```

## Valid Task Statuses

- `pending` - Task defined but not analyzed
- `ready` - Pre-planned and ready for work
- `in_progress` - Currently being worked on
- `blocked` - Blocked by dependencies
- `needs_review` - Code review required
- `needs_testing` - Verification required
- `complete` - Fully done

## Status Flow

```
pending → ready → in_progress → needs_review → needs_testing → complete
              ↓         ↓
           blocked   blocked
```

## AI Agent Usage

### Pre-planning Agent

```bash
# After analyzing, move to ready
tui-styles-cli task update 001 --status ready

# If blocked, create comment
tui-styles-cli task update 002 --status blocked
tui-styles-cli comment create 002 "## Question: Config format unclear
..." --author project-preplanner --needs-addressing
```

### Implementation Agent

```bash
# Start work
tui-styles-cli task update 001 --status in_progress

# Mark for review
tui-styles-cli task update 001 --status needs_review

# Complete
tui-styles-cli task update 001 --status complete
```

### Alignment Agent

```bash
# Block ambiguous task
tui-styles-cli task update 018 --status blocked
tui-styles-cli comment create 018 "## Alignment Issue: Missing criteria
..." --author project-milestone-aligner --needs-addressing
```

## File Formats Supported

### YAML Frontmatter (Future)
```markdown
---
task_id: "001"
title: "Task Name"
status: pending
confidence: 85
---

## Content here
```

### Markdown Format (Current)
```markdown
# Task 001: Task Name

**Status**: pending
**Confidence**: 85%
**Assigned Agent**: none

## Purpose
...
```

Both formats are supported! The CLI auto-detects and preserves the original format.

## Troubleshooting

**Lock timeout**: Another process is using the file
```bash
# Remove stale locks (>5 minutes old)
find roadmap -name "*.lock" -mmin +5 -delete
```

**Task not found**: Ensure you're in project root
```bash
cd /path/to/project
tui-styles-cli task get 001
```

**Invalid status**: Use underscores not hyphens
```bash
# ✓ Correct
tui-styles-cli task update 001 --status in_progress

# ✗ Wrong
tui-styles-cli task update 001 --status in-progress
```
