---
description: Validate milestone alignment with goals and enforce quality standards
---

I need to align a milestone. Use the tui-styles-milestone-aligner agent to:

1. Review all tasks for alignment with milestone goals
2. Check tasks have clear acceptance criteria
3. Verify task sizes are appropriate (not too large/small)
4. Flag ambiguous tasks as blocked with specific questions
5. Ensure dependencies are properly tracked

**Technical Note**: This command uses:
- `tui-styles-cli milestone info <name>` to review goals and success criteria
- `tui-styles-cli milestone tasks <name>` to list all tasks in milestone
- `tui-styles-cli task update <id> --status blocked` to flag misaligned tasks
- `tui-styles-cli comment create <id> --needs-addressing` to document alignment issues

Which milestone should I align?

-   `milestone-6` (specific milestone)
-   Leave blank to align the next active milestone
