---
description: Pre-plan pending tasks to achieve 70%+ confidence before implementation
---

I need to pre-plan pending tasks. Use the tui-styles-preplanner agent to:

1. Analyze pending tasks for clarity, size, and technical feasibility
2. Research technical approaches in codebase and documentation
3. Assign confidence scores (0-100) based on implementation readiness
4. Move tasks to ready (confidence ≥70%) or blocked (confidence <70%)
5. Create flagged comments for blockers requiring decisions

**Technical Note**: This command uses:
- `tui-styles-cli task get <id>` to read task details and acceptance criteria
- `tui-styles-cli task update <id> --status ready` to promote tasks with ≥70% confidence
- `tui-styles-cli task update <id> --status blocked` to block tasks with <70% confidence
- `tui-styles-cli comment create <id>` to add flagged comments explaining blockers

Which phase or milestone should I pre-plan? Examples:

-   `roadmap/milestone-6/phase-1-quick-wins` (specific phase)
-   `roadmap/milestone-6` (entire milestone)
-   `roadmap/milestone-6/phase-1-quick-wins/*_pending.md` (specific tasks)
