package cmd

import (
	"fmt"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/comment"
	"github.com/spf13/cobra"
)

var (
	commentAuthor          string
	commentNeedsAddressing bool
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage task comments",
	Long:  `Commands for creating and managing comments on tasks.`,
}

var commentCreateCmd = &cobra.Command{
	Use:   "create <task-id-or-file> <comment-text>",
	Short: "Create a flagged comment on a task",
	Long: `Create a comment on a task with proper YAML frontmatter.

The comment will be appended to the task file with metadata including:
  - Author
  - Timestamp (automatically set)
  - Needs addressing flag

Examples:
  tui-styles-cli comment create 001 "## BLOCKER: Library Choice

**Issue**: Multiple options available.

**Options**:
1. Option A
2. Option B" --author project-preplanner --needs-addressing

  tui-styles-cli comment create roadmap/milestone-2/phase-1/002_config_blocked.md "Decision needed" --author human --needs-addressing`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskFileOrID := args[0]
		commentText := args[1]

		if commentAuthor == "" {
			return fmt.Errorf("--author flag is required")
		}

		return comment.CreateComment(taskFileOrID, commentText, commentAuthor, commentNeedsAddressing)
	},
}

var commentListCmd = &cobra.Command{
	Use:   "list <task-id-or-file>",
	Short: "List all comments on a task",
	Long: `Display all comments on a task file.

Examples:
  tui-styles-cli comment list 001
  tui-styles-cli comment list roadmap/milestone-2/phase-1/001_task_blocked.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskFileOrID := args[0]
		return comment.ListComments(taskFileOrID)
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentCreateCmd)
	commentCmd.AddCommand(commentListCmd)

	commentCreateCmd.Flags().StringVar(&commentAuthor, "author", "", "Comment author (required)")
	commentCreateCmd.Flags().BoolVar(&commentNeedsAddressing, "needs-addressing", false, "Mark comment as needing attention")
	commentCreateCmd.MarkFlagRequired("author")
}
