---
description: Autonomously execute a milestone from planning through completion
---

I need to ship a milestone. Use the tui-styles-milestone-shipper agent to:

1. Validate milestone is properly aligned and planned
2. Execute tasks in dependency order
3. Run tests and verification after each task
4. Handle blockers by creating flagged comments
5. Track progress and update milestone status

**Technical Note**: This command uses:
- `tui-styles-cli milestone info <name>` to review milestone structure
- `tui-styles-cli task list` to find ready tasks
- `tui-styles-cli task update <id> --status in_progress` when starting work
- `tui-styles-cli task update <id> --status complete` when task verified
- `tui-styles-cli milestone complete <name>` when all tasks done

Which milestone should I ship?

-   `milestone-6` (specific milestone)
-   Leave blank and I'll recommend the next priority milestone
