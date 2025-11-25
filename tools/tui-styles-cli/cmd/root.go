package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tui-styles-cli",
	Short: "Project CLI - Task and roadmap management tool",
	Long: `Project CLI is a command-line tool for managing tasks and roadmap items.

It provides commands for:
  - Updating task statuses with automatic file renaming
  - Creating flagged comments on tasks
  - Retrieving task information
  - Managing milestones

The CLI handles file locking for concurrent access and ensures proper
YAML frontmatter formatting in task files.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
