---
description: Verify completed work meets acceptance criteria and quality standards
---

I need to verify completed work. Use the tui-styles-verifier agent to:

1. Review acceptance criteria in task files
2. Check code implementation matches requirements
3. Verify tests exist and pass
4. Validate documentation is updated
5. Mark task as complete or needs_review

**Technical Note**: This command uses:
- `tui-styles-cli task get <id>` to read acceptance criteria
- `tui-styles-cli task update <id> --status complete` to mark verified tasks
- `tui-styles-cli task update <id> --status needs_review` if issues found
- `tui-styles-cli comment create <id>` to document verification results

Which task(s) should I verify? Examples:

-   `001` (specific task ID)
-   `roadmap/milestone-6/phase-1-quick-wins/*_needs_testing.md` (all tasks needing testing)
-   Leave blank to verify all tasks with status needs_testing
