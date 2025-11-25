---
description: Analyze roadmap progress and recommend next actions for any project
---

Show me current progress and recommend what to work on next. Use the tui-styles-progress-analyst agent to:

1. Scan roadmap structure with `tui-styles-cli milestone list` and count tasks by status
2. Identify blockers using `tui-styles-cli project blockers`
3. Analyze dependencies and calculate critical path
4. Recommend next 2-3 tasks based on dependencies, risk, and momentum
5. Estimate milestone completion timeline using actual velocity from `tui-styles-cli project stats`

**Technical Note**: This command uses:
- `tui-styles-cli milestone list` to show all milestones with progress percentages
- `tui-styles-cli milestone tasks <name>` to list tasks by phase for specific milestones
- `tui-styles-cli project blockers` to find all blocked tasks requiring attention
- `tui-styles-cli project stats` to calculate completion rates and velocity
- `tui-styles-cli task list --filter <status>` to filter tasks by status (pending, ready, in_progress, etc.)

Which milestone should I analyze?

-   `milestone-6` (specific milestone)
-   Leave blank to analyze all active milestones
-   `milestone-6 milestone-7` (multiple milestones)
