# Tui Styles CLI

**The brain of the Project project** - A comprehensive command-line tool for project management, development environment control, and infrastructure operations.

## Features

### 🎯 Project Management
- **Project Status**: Real-time dashboard with milestone completion, task breakdown, and blocker identification
- **Task Management**: Update statuses with automatic file renaming and validation
- **Comments & Collaboration**: Flagged comments with YAML frontmatter for team communication
- **Milestone Tracking**: List, inspect, and analyze milestone progress
- **Statistics & Analytics**: Detailed reporting across all milestones

### 🛠️ Development Environment
- **Status Checking**: Verify required tools, running services, and port availability
- **Automated Setup**: One-command environment configuration and service startup
- **Diagnostics**: Comprehensive troubleshooting with actionable suggestions
- **Service Control**: Start/stop all development services via Docker Compose

### 🚀 Infrastructure & Operations
- **DNS Management**: Manage DNS records (placeholder for Route53, Cloudflare, etc.)
- **Deployments**: Deploy to staging and production with safety checks
- **Health Monitoring**: Check service health across environments
- **Secrets Management**: Secure secret storage and retrieval (placeholder for Vault, AWS Secrets Manager, etc.)

See [FEATURES.md](FEATURES.md) for complete documentation.

## Installation

### Quick Install (Recommended)

No sudo required! Installs to `~/.local/bin`:

```bash
cd tools/tui-styles-cli
make install
```

The installer will check if `~/.local/bin` is in your PATH and provide instructions if needed.

### System-wide Install (Optional)

Requires sudo, installs to `/usr/local/bin`:

```bash
cd tools/tui-styles-cli
make install-system
```

### Manual Build

```bash
cd tools/tui-styles-cli
make build
# Binary will be in bin/tui-styles-cli
```

## Usage

All commands should be run from the project root directory where the `roadmap/` folder exists.

### Task Commands

#### Update Task Status

Update a task's status and automatically rename the file:

```bash
# Using task ID (searches roadmap/)
tui-styles-cli task update 001 --status ready

# Using full file path
tui-styles-cli task update roadmap/milestone-2/phase-1/001_task_pending.md --status in_progress
```

**Valid Statuses:**
- `pending` - Task is defined but not analyzed
- `ready` - Task is pre-planned and ready for implementation
- `in_progress` - Task is currently being worked on
- `blocked` - Task is blocked by dependencies or decisions
- `complete` - Task is fully completed
- `needs_testing` - Implementation done, needs verification
- `needs_review` - Implementation done, needs code review

**Status Flow:**
```
pending → ready → in_progress → needs_review → needs_testing → complete
              ↓         ↓
           blocked   blocked
```

#### Get Task Information

Display task metadata and current status:

```bash
# Using task ID
tui-styles-cli task get 001

# Using file path
tui-styles-cli task get roadmap/milestone-2/phase-1/001_task_ready.md
```

**Output Example:**
```
Task ID: 001
Title: Setup Project Structure
Status: ready
Confidence: 95%
Assigned To: project-preplanner
Dependencies: []
Created: 2025-11-19T10:30:00Z
Updated: 2025-11-19T15:45:00Z
File: roadmap/milestone-2/phase-1/001_task_ready.md
```

### Comment Commands

#### Create Flagged Comment

Add a comment to a task file with proper frontmatter:

```bash
tui-styles-cli comment create 001 "## BLOCKER: Library Choice

**Issue**: Multiple options available for HTTP routing library.

**Options**:
1. Chi (lightweight, idiomatic)
2. Gorilla Mux (feature-rich)

**Recommendation**: Choose Chi based on guidelines
" --author project-preplanner --needs-addressing
```

**Comment Format:**
The CLI automatically generates properly formatted comments:

```markdown
<!--
author: project-preplanner
timestamp: 2025-11-19T16:20:00Z
needs_addressing: true
-->

## BLOCKER: Library Choice

**Issue**: Multiple options available for HTTP routing library.
...
```

#### List Comments

View all comments on a task:

```bash
tui-styles-cli comment list 001
```

**Output Example:**
```
Comments in 001_task_blocked.md:

[1] ⚠ NEEDS ADDRESSING
    Author: project-preplanner
    Time: 2025-11-19T16:20:00Z
    Content:
      ## BLOCKER: Library Choice

      **Issue**: Multiple options available for HTTP routing library.
      ...

[2] ✓
    Author: human-developer
    Time: 2025-11-19T17:30:00Z
    Content:
      Decision made: Use Chi framework.
```

### Milestone Commands

#### List Milestones

Show all milestones in the roadmap:

```bash
tui-styles-cli milestone list
```

**Output:**
```
Milestones:
  ✓ milestone-1
  ✓ milestone-2
  ✓ milestone-3
  ⚠ milestone-4 (missing index.md)
```

#### Get Milestone Info

Display milestone details:

```bash
# Using milestone name
tui-styles-cli milestone info milestone-2

# Using milestone number
tui-styles-cli milestone info 2
```

**Output:**
```
Milestone: milestone-2
Path: roadmap/milestone-2

Phases: 3

Index Preview:
---
# Milestone 2: Creator API Foundation
...
```

## File Locking

The CLI uses file-level locking to prevent concurrent modifications:

- **Automatic Lock Acquisition**: Each operation acquires a `.lock` file
- **Timeout**: 5-second timeout for lock acquisition
- **Automatic Cleanup**: Locks are released after operation completes
- **Stale Lock Detection**: Locks older than 5 minutes are automatically cleaned

### Lock File Format

Lock files are created as `<taskfile>.lock`:

```
<PID>
<timestamp>
```

Example: `001_task_pending.md.lock`

## Architecture

### Project Structure

```
tui-styles-cli/
├── cmd/                    # Cobra command definitions
│   ├── root.go            # Root command
│   ├── task.go            # Task commands
│   ├── comment.go         # Comment commands
│   └── milestone.go       # Milestone commands
├── internal/
│   ├── task/              # Task management logic
│   │   └── task.go
│   ├── comment/           # Comment management logic
│   │   └── comment.go
│   └── filelock/          # File locking implementation
│       └── filelock.go
├── pkg/
│   └── types/             # Shared type definitions
│       └── task.go
├── main.go                # Entry point
├── Makefile               # Build automation
└── README.md              # This file
```

### Design Principles

1. **CLI-First for Consistency**: All task/comment modifications go through the CLI to ensure:
   - Valid YAML frontmatter formatting
   - Proper file locking for concurrent access
   - Consistent status transitions
   - Automated file renaming

2. **Atomic Operations**: File updates use atomic operations:
   - Create new file with updated content
   - Remove old file only after new file succeeds
   - Rollback on errors

3. **Validation**: All inputs are validated before processing:
   - Status values must be valid
   - Task files must have proper frontmatter
   - Required flags are enforced

4. **Error Handling**: Clear error messages with actionable guidance

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for build automation)

### Build

```bash
make build
```

Binary will be in `bin/tui-styles-cli`.

### Test

```bash
make test
```

Run with coverage:

```bash
make coverage
```

### Format Code

```bash
make fmt
```

### Lint

```bash
make lint
```

Requires [golangci-lint](https://golangci-lint.run/usage/install/).

### Clean

```bash
make clean
```

## Integration with AI Agents

The CLI is designed to be used by AI agents for automated task management:

### project-preplanner Agent

```bash
# Move task to ready after pre-planning
tui-styles-cli task update 001 --status ready

# Block task and create comment
tui-styles-cli task update 002 --status blocked
tui-styles-cli comment create 002 "## Question: Required Configuration Fields

Cannot define config file format without knowing what agents need to configure.

**Options**:
1. Minimal config (port, env)
2. Full config (all settings)

**Recommendation**: Start minimal
" --author project-preplanner --needs-addressing
```

### project-milestone-aligner Agent

```bash
# Block ambiguous task
tui-styles-cli task update 018 --status blocked
tui-styles-cli comment create 018 "## Alignment Issue: Missing Implementation Details

Task lacks specific acceptance criteria.
" --author project-milestone-aligner --needs-addressing
```

### Implementation Agents

```bash
# Start work
tui-styles-cli task update 001 --status in_progress

# Mark for review
tui-styles-cli task update 001 --status needs_review

# Complete task
tui-styles-cli task update 001 --status complete
```

## Troubleshooting

### Lock Acquisition Timeout

**Error:** `timeout acquiring lock for 001_task_pending.md.lock`

**Solution:** Another process is modifying the file. Wait and retry, or manually remove stale `.lock` files.

### Task File Not Found

**Error:** `task file not found for ID: 001`

**Solution:** Ensure you're running from the project root where `roadmap/` exists, and the task ID is correct.

### Invalid Status

**Error:** `invalid status: in-progress (valid: [pending ready in_progress blocked complete needs_testing needs_review])`

**Solution:** Use underscores, not hyphens: `in_progress` not `in-progress`.

## License

Part of the Project project.

## Support

For issues or questions, create an issue in the Project repository.
