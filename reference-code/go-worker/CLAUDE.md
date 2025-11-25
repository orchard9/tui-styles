# CLAUDE.md

> **📝 Helpful Tips Directive**: Whenever we make a mistake during development and figure out the correct approach, add a simple tip to the "Helpful Tips" section at the bottom of this file.

**Project**: /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker  
**Version**: 0.1.0  
**Description**: HTTP API built with go-core-http-toolkit with postgres, redis

## Forge Project Management

> **📌 Project Management**: This project is managed by .forge - use the forge CLI tools below for all task and project management operations. Do not manually edit .forge directories.

### Environment Setup

```bash
# Required environment variable
export FORGE_SERVER=http://localhost:52502

# Helpful aliases
alias ft='forge-cli list-tasks'
alias ft-todo='forge-cli list-tasks --status=todo'
alias ft-wip='forge-cli list-tasks --status=in-progress'
alias fp-up='forge-project start --open'
alias fp-down='forge-project stop'
```

### Task Management Commands

**⚠️ CRITICAL**: Forge tasks must be detailed implementation blueprints with specific file paths, line numbers, and code examples in implementation_notes.

- `forge-cli create-task "<title>" [--definition="..."] [--criteria="..."] [--deps="task-id"]` - Create new tasks when starting work
- `forge-cli list-tasks [--status=todo|in-progress|in-review|completed] [--format=json]` - List tasks filtered by status
- `forge-cli get-task <task-id> [--detailed]` - View task details before starting work
- `forge-cli move-task <task-id> <status>` - Move tasks through workflow (todo→in-progress→in-review→completed)
- `forge-cli delete-task <task-id> [--force] [--reason="..."]` - Delete obsolete/duplicate tasks

### Project Management Commands

- `forge-project start [--open]` - Start forge server and web UI when beginning development session
- `forge-project stop` - Stop all forge services when done working
- `forge-project status [--json]` - Check service health
- `forge-project restart [--open]` - Restart services after configuration changes

### Forge Development Workflow

1. **Start Development Session**:
   ```bash
   forge-project start --open  # Starts server + web UI + opens browser
   ```

2. **Find Available Work**:
   ```bash
   forge-cli list-tasks --status=todo  # See available tasks
   forge-cli get-task <task-id> --detailed  # Review task details
   ```

3. **Start Working on Task**:
   ```bash
   forge-cli move-task <task-id> in-progress  # Mark as started
   git checkout -b task/<task-id>  # Create branch for work
   # Follow implementation_notes step-by-step for exact file paths and code
   ```

4. **Complete Task**:
   - Follow all implementation_notes steps exactly
   - Verify all acceptance criteria are met
   - **MANDATORY**: Run `make ci` and fix ALL issues that arise
   - Commit all changes with descriptive message
   ```bash
   forge-cli move-task <task-id> completed  # Mark as done
   ```

5. **End Development Session**:
   ```bash
   forge-project stop  # Clean shutdown
   ```

## Available Commands

- `make help` - Show all available commands
- `make build` - Build the application binary
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make lint` - Run code linting
- `make fmt` - Format code
- `make ci` - Run all CI checks (lint + test + build)
- `make clean` - Clean build artifacts
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make dev` - Run application in development mode
- `make watch` - Run application with hot reload

## Development Rules

- **Testing**: All code must have tests with >80% coverage
- **Linting**: Code must pass `make lint`
- **CI**: `make ci` must pass before any commit
- **Pre-commit Hooks**: We use local pre-commit hooks for quality checks
- **NO GitHub Hooks/CI**: We never use GitHub hooks or GitHub CI - all checks run locally
- **Dependencies**: Use only approved toolkit dependencies
- **Documentation**: Less is more - clean as we go, keep documentation minimal and efficient
- **Communication**: Hyper-efficient communication - focus on essential information only

## 📝 Helpful Tips

*This section contains lessons learned from mistakes made during development.*

1. **Always run make ci before committing**
   - **Mistake**: Committing code without running full CI checks
   - **Correct approach**: Always run `make ci` and fix all issues
   - **Why**: Prevents broken builds and maintains code quality

