package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/task"
	"github.com/spf13/cobra"
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage milestones",
	Long:  `Commands for managing milestones in the roadmap.`,
}

var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all milestones",
	Long:  `List all milestones in the roadmap directory with status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listMilestones()
	},
}

var milestoneInfoCmd = &cobra.Command{
	Use:   "info <milestone-name>",
	Short: "Get milestone information",
	Long: `Display detailed information about a milestone.

Examples:
  tui-styles-cli milestone info milestone-2
  tui-styles-cli milestone info 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return showMilestoneInfo(args[0])
	},
}

var milestoneTasksCmd = &cobra.Command{
	Use:   "tasks <milestone-name>",
	Short: "List all tasks in a milestone",
	Long: `List all tasks in a milestone with their current status.

Examples:
  tui-styles-cli milestone tasks milestone-2
  tui-styles-cli milestone tasks 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listMilestoneTasks(args[0])
	},
}

var milestoneUpdateCmd = &cobra.Command{
	Use:   "update <milestone-name>",
	Short: "Update milestone metadata",
	Long: `Update milestone status, goals, or success criteria.

Opens the milestone index.md in your default editor.

Examples:
  tui-styles-cli milestone update milestone-2
  tui-styles-cli milestone update 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateMilestone(args[0])
	},
}

var milestoneCompleteCmd = &cobra.Command{
	Use:   "complete <milestone-name>",
	Short: "Mark milestone as complete",
	Long: `Mark a milestone as complete and update its status in index.md.

Examples:
  tui-styles-cli milestone complete milestone-2
  tui-styles-cli milestone complete 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return completeMilestone(args[0])
	},
}

var milestoneCreateCmd = &cobra.Command{
	Use:   "create <milestone-name> <title>",
	Short: "Create a new milestone",
	Long: `Create a new milestone with the standard structure.

Examples:
  tui-styles-cli milestone create milestone-6 "Auth Service Implementation"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createMilestone(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(milestoneCmd)
	milestoneCmd.AddCommand(milestoneListCmd)
	milestoneCmd.AddCommand(milestoneInfoCmd)
	milestoneCmd.AddCommand(milestoneTasksCmd)
	milestoneCmd.AddCommand(milestoneUpdateCmd)
	milestoneCmd.AddCommand(milestoneCompleteCmd)
	milestoneCmd.AddCommand(milestoneCreateCmd)
}

func listMilestones() error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	entries, err := os.ReadDir(roadmapDir)
	if err != nil {
		return fmt.Errorf("failed to read roadmap directory: %w", err)
	}

	fmt.Println("Milestones:")
	fmt.Println()

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "milestone-") {
			continue
		}

		milestonePath := filepath.Join(roadmapDir, entry.Name())
		indexPath := filepath.Join(milestonePath, "index.md")

		if _, err := os.Stat(indexPath); err != nil {
			fmt.Printf("  ⚠️  %s (missing index.md)\n", entry.Name())
			continue
		}

		// Read status from index.md
		content, _ := os.ReadFile(indexPath)
		lines := strings.Split(string(content), "\n")

		status := "Unknown"
		title := entry.Name()

		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				title = strings.TrimPrefix(line, "# ")
			}
			if strings.HasPrefix(line, "**Status**:") {
				status = strings.TrimSpace(strings.TrimPrefix(line, "**Status**:"))
			}
		}

		// Get task stats
		stats := countTasksInMilestone(milestonePath)
		completion := 0
		if stats["total"] > 0 {
			completion = (stats["complete"] * 100) / stats["total"]
		}

		statusIcon := "🔵"
		switch strings.ToLower(status) {
		case "complete":
			statusIcon = "✅"
		case "in progress":
			statusIcon = "🟡"
		case "blocked":
			statusIcon = "🔴"
		}

		fmt.Printf("  %s %s\n", statusIcon, entry.Name())
		fmt.Printf("     Title: %s\n", title)
		fmt.Printf("     Status: %s\n", status)
		fmt.Printf("     Progress: %d%% (%d/%d tasks)\n", completion, stats["complete"], stats["total"])
		if stats["blocked"] > 0 {
			fmt.Printf("     ⚠️  %d blocked tasks\n", stats["blocked"])
		}
		fmt.Println()
	}

	return nil
}

func showMilestoneInfo(milestoneName string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	milestonePath := filepath.Join("roadmap", milestoneName)
	if _, err := os.Stat(milestonePath); os.IsNotExist(err) {
		return fmt.Errorf("milestone not found: %s", milestoneName)
	}

	// Read index.md
	indexPath := filepath.Join(milestonePath, "index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read milestone index: %w", err)
	}

	fmt.Printf("📋 Milestone: %s\n", milestoneName)
	fmt.Printf("Path: %s\n", milestonePath)
	fmt.Println()

	// Count phases and tasks
	entries, err := os.ReadDir(milestonePath)
	if err == nil {
		phaseCount := 0
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), "phase-") {
				phaseCount++
			}
		}
		stats := countTasksInMilestone(milestonePath)

		fmt.Printf("Phases: %d\n", phaseCount)
		fmt.Printf("Tasks: %d total (%d complete, %d in progress, %d blocked)\n\n",
			stats["total"], stats["complete"], stats["in_progress"], stats["blocked"])
	}

	// Display full index.md content
	fmt.Println("Index Content:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(string(content))

	return nil
}

func listMilestoneTasks(milestoneName string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	milestonePath := filepath.Join("roadmap", milestoneName)
	if _, err := os.Stat(milestonePath); os.IsNotExist(err) {
		return fmt.Errorf("milestone not found: %s", milestoneName)
	}

	fmt.Printf("📋 Tasks in %s\n", milestoneName)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Walk through phases
	entries, err := os.ReadDir(milestonePath)
	if err != nil {
		return fmt.Errorf("failed to read milestone directory: %w", err)
	}

	totalTasks := 0

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "phase-") {
			continue
		}

		phasePath := filepath.Join(milestonePath, entry.Name())
		fmt.Printf("Phase: %s\n", entry.Name())
		fmt.Println(strings.Repeat("-", 80))

		// List tasks in this phase
		phaseEntries, err := os.ReadDir(phasePath)
		if err != nil {
			continue
		}

		for _, taskEntry := range phaseEntries {
			if taskEntry.IsDir() || !strings.HasSuffix(taskEntry.Name(), ".md") {
				continue
			}

			filename := taskEntry.Name()

			// Extract task number and status from filename
			parts := strings.Split(filename, "_")
			if len(parts) < 2 {
				continue
			}

			taskNum := parts[0]
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
			}

			// Read title from file
			taskPath := filepath.Join(phasePath, filename)
			frontmatter, _, _ := task.ParseTaskFile(taskPath)

			statusIcon := getStatusIcon(status)
			title := filename
			if frontmatter != nil && frontmatter.Title != "" {
				title = frontmatter.Title
			}

			fmt.Printf("  %s [%s] %s\n", statusIcon, taskNum, title)
			totalTasks++
		}

		fmt.Println()
	}

	fmt.Printf("Total tasks: %d\n", totalTasks)

	return nil
}

func updateMilestone(milestoneName string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	milestonePath := filepath.Join("roadmap", milestoneName)
	if _, err := os.Stat(milestonePath); os.IsNotExist(err) {
		return fmt.Errorf("milestone not found: %s", milestoneName)
	}

	indexPath := filepath.Join(milestonePath, "index.md")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return fmt.Errorf("index.md not found in %s", milestoneName)
	}

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	fmt.Printf("Opening %s in %s...\n", indexPath, editor)

	cmd := exec.Command(editor, indexPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	fmt.Println("✓ Milestone updated")

	return nil
}

func completeMilestone(milestoneName string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	milestonePath := filepath.Join("roadmap", milestoneName)
	if _, err := os.Stat(milestonePath); os.IsNotExist(err) {
		return fmt.Errorf("milestone not found: %s", milestoneName)
	}

	indexPath := filepath.Join(milestonePath, "index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index.md: %w", err)
	}

	// Update status to Complete
	lines := strings.Split(string(content), "\n")
	updated := false

	for i, line := range lines {
		if strings.HasPrefix(line, "**Status**:") {
			lines[i] = "**Status**: Complete"
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("could not find **Status**: line in index.md")
	}

	// Make file writable before writing
	if err := os.Chmod(indexPath, 0644); err != nil {
		return fmt.Errorf("failed to make file writable: %w", err)
	}

	// Write back
	if err := os.WriteFile(indexPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write index.md: %w", err)
	}

	fmt.Printf("✅ Milestone %s marked as Complete\n", milestoneName)

	// Show summary
	stats := countTasksInMilestone(milestonePath)
	if stats["total"] > 0 {
		completion := (stats["complete"] * 100) / stats["total"]
		fmt.Printf("   Tasks: %d/%d complete (%d%%)\n", stats["complete"], stats["total"], completion)

		if stats["complete"] < stats["total"] {
			fmt.Printf("   ⚠️  Warning: Not all tasks are complete!\n")
			fmt.Printf("   - In Progress: %d\n", stats["in_progress"])
			fmt.Printf("   - Pending: %d\n", stats["pending"])
			fmt.Printf("   - Blocked: %d\n", stats["blocked"])
		}
	}

	return nil
}

func createMilestone(milestoneName, title string) error {
	// Normalize milestone name
	if !strings.HasPrefix(milestoneName, "milestone-") {
		milestoneName = "milestone-" + milestoneName
	}

	milestonePath := filepath.Join("roadmap", milestoneName)

	// Check if milestone already exists
	if _, err := os.Stat(milestonePath); err == nil {
		return fmt.Errorf("milestone already exists: %s", milestoneName)
	}

	// Create milestone directory
	if err := os.MkdirAll(milestonePath, 0755); err != nil {
		return fmt.Errorf("failed to create milestone directory: %w", err)
	}

	// Create standard phases
	phases := []string{
		"phase-1-quick-wins",
		"phase-2-core-features",
		"phase-3-polish",
	}

	for _, phase := range phases {
		phasePath := filepath.Join(milestonePath, phase)
		if err := os.MkdirAll(phasePath, 0755); err != nil {
			return fmt.Errorf("failed to create phase directory: %w", err)
		}
	}

	// Create index.md template
	indexPath := filepath.Join(milestonePath, "index.md")
	indexContent := fmt.Sprintf(`# %s

**Status**: Pending
**Owner**: TBD
**Duration**: TBD
**Dependencies**: []

## Goals

[Define 2-3 high-level goals for this milestone]

1. Goal 1
2. Goal 2
3. Goal 3

## Success Criteria

- [ ] Success criterion 1
- [ ] Success criterion 2
- [ ] Success criterion 3

## Scope

**In Scope**:
- Feature/capability 1
- Feature/capability 2

**Out of Scope**:
- Feature/capability that belongs in a later milestone

## Technical Approach

[High-level technical approach and architecture decisions]

## Risks

1. **Risk 1**: Description and mitigation strategy
2. **Risk 2**: Description and mitigation strategy

## Timeline

**Phase 1 - Quick Wins**: X days
**Phase 2 - Core Features**: Y days
**Phase 3 - Polish**: Z days

**Total**: N days

## Tasks by Phase

### Phase 1: Quick Wins
[Tasks will be added here]

### Phase 2: Core Features
[Tasks will be added here]

### Phase 3: Polish
[Tasks will be added here]
`, title)

	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.md: %w", err)
	}

	fmt.Printf("✅ Created milestone: %s\n", milestoneName)
	fmt.Printf("   Path: %s\n", milestonePath)
	fmt.Printf("   Phases created: %d\n", len(phases))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit the milestone: tui-styles-cli milestone update %s\n", milestoneName)
	fmt.Println("  2. Add tasks to the phases")
	fmt.Println("  3. Update goals and success criteria")

	return nil
}

func getStatusIcon(status string) string {
	switch status {
	case "complete":
		return "✅"
	case "in_progress":
		return "🟡"
	case "blocked":
		return "🔴"
	case "ready":
		return "🟢"
	case "needs_review":
		return "👀"
	case "needs_testing":
		return "🧪"
	case "needs_human_verification":
		return "🔍"
	case "pending":
		return "⏳"
	default:
		return "❓"
	}
}
