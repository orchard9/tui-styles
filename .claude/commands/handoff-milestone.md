---
description: Prepare milestone handoff documentation for human developers
---

I need to create a handoff document. Use the tui-styles-handoff-writer agent to:

1. Summarize milestone goals and what was accomplished
2. Document key technical decisions and their rationale
3. List completed tasks with links to implementation
4. Note any blockers or outstanding questions
5. Create handoff.md in milestone directory

**Technical Note**: This command uses:
- `tui-styles-cli milestone info <name>` to get milestone details
- `tui-styles-cli milestone tasks <name>` to list all tasks
- `tui-styles-cli comment list <id>` to gather technical decisions

Which milestone should I document for handoff?

-   `milestone-6` (specific milestone)
-   Leave blank to document the most recently completed milestone
