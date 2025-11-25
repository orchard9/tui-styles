package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/task"
	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
	"github.com/spf13/cobra"
)

var (
	taskStatus string
	taskFilter string
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  `Commands for managing tasks in the roadmap.`,
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update <task-id-or-file> --status <status>",
	Short: "Update task status",
	Long: `Update the status of a task and rename the file accordingly.

Valid statuses: pending, ready, in_progress, blocked, complete, needs_testing, needs_review, needs_human_verification

Examples:
  tui-styles-cli task update 001 --status ready
  tui-styles-cli task update roadmap/milestone-2/phase-1/001_task_pending.md --status in_progress`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskFileOrID := args[0]

		if taskStatus == "" {
			return fmt.Errorf("--status flag is required")
		}

		status := types.TaskStatus(taskStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status: %s (valid: %v)", taskStatus, types.ValidStatuses())
		}

		return task.UpdateTaskStatus(taskFileOrID, status)
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get <task-id-or-file>",
	Short: "Get task information",
	Long: `Retrieve and display information about a task.

Examples:
  tui-styles-cli task get 001
  tui-styles-cli task get roadmap/milestone-2/phase-1/001_task_ready.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskFileOrID := args[0]
		return task.GetTaskInfo(taskFileOrID)
	},
}

var taskEditCmd = &cobra.Command{
	Use:   "edit <task-id-or-file>",
	Short: "Edit a task file",
	Long: `Open a task file in your default editor.

Uses the $EDITOR environment variable, defaults to vim.

Examples:
  tui-styles-cli task edit 001
  tui-styles-cli task edit roadmap/milestone-2/phase-1/001_task_ready.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return editTask(args[0])
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create <milestone> <phase> <task-number> <title>",
	Short: "Create a new task",
	Long: `Create a new task file with standard template.

Examples:
  tui-styles-cli task create milestone-3 phase-1-quick-wins 010 "Implement JWT auth"
  tui-styles-cli task create 3 phase-2-core-features 015 "Add user profile endpoint"`,
	Args: cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createTask(args[0], args[1], args[2], args[3])
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long: `List all tasks across all milestones with optional filtering.

Examples:
  tui-styles-cli task list
  tui-styles-cli task list --filter blocked
  tui-styles-cli task list --filter in_progress`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listAllTasks()
	},
}

var taskSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for tasks by title or content",
	Long: `Search for tasks matching a query string.

Examples:
  tui-styles-cli task search authentication
  tui-styles-cli task search "health check"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return searchTasks(args[0])
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskCmd.AddCommand(taskGetCmd)
	taskCmd.AddCommand(taskEditCmd)
	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskSearchCmd)

	taskUpdateCmd.Flags().StringVar(&taskStatus, "status", "", "New status (required)")
	taskUpdateCmd.MarkFlagRequired("status")

	taskListCmd.Flags().StringVar(&taskFilter, "filter", "", "Filter by status (pending, ready, in_progress, blocked, complete)")
}

func editTask(taskFileOrID string) error {
	// Find task file if ID provided
	taskFile := taskFileOrID
	if !strings.HasSuffix(taskFileOrID, ".md") {
		found, err := task.FindTaskFile(taskFileOrID)
		if err != nil {
			return err
		}
		taskFile = found
	}

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	fmt.Printf("Opening %s in %s...\n", taskFile, editor)

	cmd := exec.Command(editor, taskFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	fmt.Println("✓ Task updated")

	return nil
}

func createTask(milestoneName, phaseName, taskNum, title string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	// Pad task number to 3 digits
	if len(taskNum) < 3 {
		taskNum = fmt.Sprintf("%03s", taskNum)
	}

	milestonePath := filepath.Join("roadmap", milestoneName)
	if _, err := os.Stat(milestonePath); os.IsNotExist(err) {
		return fmt.Errorf("milestone not found: %s", milestoneName)
	}

	phasePath := filepath.Join(milestonePath, phaseName)
	if _, err := os.Stat(phasePath); os.IsNotExist(err) {
		return fmt.Errorf("phase not found: %s", phaseName)
	}

	// Create task filename
	taskSlug := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	taskSlug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, taskSlug)

	filename := fmt.Sprintf("%s_%s_pending.md", taskNum, taskSlug)
	taskPath := filepath.Join(phasePath, filename)

	// Check if task already exists
	if _, err := os.Stat(taskPath); err == nil {
		return fmt.Errorf("task already exists: %s", taskPath)
	}

	// Create task file with template
	taskContent := fmt.Sprintf(`# Task %s: %s

**Status**: pending
**Phase**: %s
**Milestone**: %s
**Dependencies**: []
**Assigned Agent**: none
**Confidence**: TBD (assign during pre-planning)

## Purpose

[Describe what this task accomplishes and why it's needed]

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Technical Approach

[Describe how this will be implemented]

**Files to Create/Modify**:
- file1.go
- file2.go

**Dependencies**:
- Library/service dependencies

## Testing Strategy

**Unit Tests**:
- Test case 1
- Test case 2

**Integration Tests**:
- Integration scenario 1

## Notes

[Any additional context, links, or considerations]
`, taskNum, title, phaseName, milestoneName)

	if err := os.WriteFile(taskPath, []byte(taskContent), 0644); err != nil {
		return fmt.Errorf("failed to create task file: %w", err)
	}

	fmt.Printf("✅ Created task: %s\n", filename)
	fmt.Printf("   Path: %s\n", taskPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit the task: tui-styles-cli task edit %s\n", taskNum)
	fmt.Printf("  2. Pre-plan the task: tui-styles-cli /preplan %s\n", taskNum)
	fmt.Printf("  3. Move to ready: tui-styles-cli task update %s --status ready\n", taskNum)

	return nil
}

func listAllTasks() error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	fmt.Println("📋 All Tasks")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	totalTasks := 0
	filteredTasks := 0

	err := filepath.Walk(roadmapDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") || !strings.Contains(path, "phase-") {
			return nil
		}

		filename := filepath.Base(path)

		// Extract task number and status
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return nil
		}

		taskNum := parts[0]
		if len(taskNum) != 3 || taskNum[0] < '0' || taskNum[0] > '9' {
			return nil
		}

		totalTasks++

		// Determine status
		status := "unknown"
		if strings.Contains(filename, "_complete.md") {
			status = "complete"
		} else if strings.Contains(filename, "_in_progress.md") {
			status = "in_progress"
		} else if strings.Contains(filename, "_blocked.md") {
			status = "blocked"
		} else if strings.Contains(filename, "_ready.md") {
			status = "ready"
		} else if strings.Contains(filename, "_pending.md") {
			status = "pending"
		} else if strings.Contains(filename, "_needs_review.md") {
			status = "needs_review"
		} else if strings.Contains(filename, "_needs_testing.md") {
			status = "needs_testing"
		} else if strings.Contains(filename, "_needs_human_verification.md") {
			status = "needs_human_verification"
		} else if strings.Contains(filename, "_needs_review.md") {
			status = "needs_review"
		} else if strings.Contains(filename, "_needs_testing.md") {
			status = "needs_testing"
		} else if strings.Contains(filename, "_needs_human_verification.md") {
			status = "needs_human_verification"
		}

		// Apply filter
		if taskFilter != "" && status != taskFilter {
			return nil
		}

		filteredTasks++

		// Parse task file
		frontmatter, _, _ := task.ParseTaskFile(path)

		statusIcon := getStatusIcon(status)
		title := filename
		if frontmatter != nil && frontmatter.Title != "" {
			title = frontmatter.Title
		}

		// Get milestone and phase from path
		pathParts := strings.Split(path, string(filepath.Separator))
		milestone := "unknown"
		phase := "unknown"
		for i, part := range pathParts {
			if strings.HasPrefix(part, "milestone-") {
				milestone = part
				if i+1 < len(pathParts) && strings.HasPrefix(pathParts[i+1], "phase-") {
					phase = pathParts[i+1]
				}
			}
		}

		fmt.Printf("%s [%s] %s\n", statusIcon, taskNum, title)
		fmt.Printf("   %s / %s\n", milestone, phase)
		fmt.Printf("   Path: %s\n", path)
		fmt.Println()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk roadmap: %w", err)
	}

	if taskFilter != "" {
		fmt.Printf("Showing %d of %d tasks (filter: %s)\n", filteredTasks, totalTasks, taskFilter)
	} else {
		fmt.Printf("Total tasks: %d\n", totalTasks)
	}

	return nil
}

func searchTasks(query string) error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	query = strings.ToLower(query)

	fmt.Printf("🔍 Searching for: %s\n", query)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	matchCount := 0

	err := filepath.Walk(roadmapDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") || !strings.Contains(path, "phase-") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentLower := strings.ToLower(string(content))

		// Check if query matches
		if !strings.Contains(contentLower, query) {
			return nil
		}

		matchCount++

		// Parse task file
		frontmatter, _, _ := task.ParseTaskFile(path)

		filename := filepath.Base(path)
		parts := strings.Split(filename, "_")
		taskNum := "???"
		if len(parts) > 0 {
			taskNum = parts[0]
		}

		title := filename
		if frontmatter != nil && frontmatter.Title != "" {
			title = frontmatter.Title
		}

		status := "unknown"
		if strings.Contains(filename, "_complete.md") {
			status = "complete"
		} else if strings.Contains(filename, "_in_progress.md") {
			status = "in_progress"
		} else if strings.Contains(filename, "_blocked.md") {
			status = "blocked"
		} else if strings.Contains(filename, "_ready.md") {
			status = "ready"
		} else if strings.Contains(filename, "_pending.md") {
			status = "pending"
		} else if strings.Contains(filename, "_needs_review.md") {
			status = "needs_review"
		} else if strings.Contains(filename, "_needs_testing.md") {
			status = "needs_testing"
		} else if strings.Contains(filename, "_needs_human_verification.md") {
			status = "needs_human_verification"
		}

		statusIcon := getStatusIcon(status)

		fmt.Printf("%s [%s] %s\n", statusIcon, taskNum, title)
		fmt.Printf("   Path: %s\n", path)
		fmt.Println()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to search tasks: %w", err)
	}

	fmt.Printf("Found %d matching tasks\n", matchCount)

	return nil
}
