---
description: Update project state based on external changes or decisions
---

I need to update project state. Use the tui-styles-state-updater agent to:

1. Scan for external changes (new requirements, tech decisions, blockers resolved)
2. Update affected task statuses and priorities
3. Add comments documenting state changes
4. Recalculate milestone timelines
5. Identify new risks or dependencies

**Technical Note**: This command uses:
- `tui-styles-cli task list` to find affected tasks
- `tui-styles-cli task update <id> --status <new-status>` to update states
- `tui-styles-cli comment create <id>` to document changes
- `tui-styles-cli milestone info <name>` to update milestone metadata

What has changed that requires state update?

-   New technical decision
-   External dependency resolved
-   Requirement change
-   Blocker resolved
-   Priority shift
