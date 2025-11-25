package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  `Commands for managing the overall project structure and status.`,
}

var projectStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall project status",
	Long: `Display comprehensive project status including:
  - Active milestones and completion %
  - Tasks by status (pending, ready, in_progress, blocked, complete)
  - Recent activity
  - Blockers requiring attention`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showProjectStatus()
	},
}

var projectStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show project statistics",
	Long:  `Display detailed statistics across all milestones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showProjectStats()
	},
}

var projectBlockersCmd = &cobra.Command{
	Use:   "blockers",
	Short: "List all blocked tasks",
	Long:  `Show all tasks with status 'blocked' across all milestones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listBlockers()
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectStatsCmd)
	projectCmd.AddCommand(projectBlockersCmd)
}

func showProjectStatus() error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	fmt.Println("🎯 Project Project Status")
	fmt.Println("================================")
	fmt.Println()

	// Get all milestones
	entries, err := os.ReadDir(roadmapDir)
	if err != nil {
		return fmt.Errorf("failed to read roadmap: %w", err)
	}

	totalTasks := 0
	completeTasks := 0
	inProgressTasks := 0
	blockedTasks := 0
	readyTasks := 0
	pendingTasks := 0

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "milestone-") {
			continue
		}

		milestonePath := filepath.Join(roadmapDir, entry.Name())
		stats := countTasksInMilestone(milestonePath)

		totalTasks += stats["total"]
		completeTasks += stats["complete"]
		inProgressTasks += stats["in_progress"]
		blockedTasks += stats["blocked"]
		readyTasks += stats["ready"]
		pendingTasks += stats["pending"]

		// Show milestone summary
		completion := 0
		if stats["total"] > 0 {
			completion = (stats["complete"] * 100) / stats["total"]
		}

		status := "🔵"
		if completion == 100 {
			status = "✅"
		} else if stats["in_progress"] > 0 {
			status = "🟡"
		} else if stats["blocked"] > 0 {
			status = "🔴"
		}

		fmt.Printf("%s %s: %d%% complete (%d/%d tasks)\n",
			status, entry.Name(), completion, stats["complete"], stats["total"])
	}

	fmt.Println()
	fmt.Println("Overall Statistics")
	fmt.Println("------------------")
	fmt.Printf("Total Tasks: %d\n", totalTasks)
	fmt.Printf("  ✅ Complete:     %d (%d%%)\n", completeTasks, percentage(completeTasks, totalTasks))
	fmt.Printf("  🟡 In Progress:  %d (%d%%)\n", inProgressTasks, percentage(inProgressTasks, totalTasks))
	fmt.Printf("  🟢 Ready:        %d (%d%%)\n", readyTasks, percentage(readyTasks, totalTasks))
	fmt.Printf("  ⏳ Pending:      %d (%d%%)\n", pendingTasks, percentage(pendingTasks, totalTasks))
	fmt.Printf("  🔴 Blocked:      %d (%d%%)\n", blockedTasks, percentage(blockedTasks, totalTasks))

	if blockedTasks > 0 {
		fmt.Println()
		fmt.Printf("⚠️  %d tasks are blocked and need attention!\n", blockedTasks)
		fmt.Println("   Run 'tui-styles-cli project blockers' for details")
	}

	return nil
}

func showProjectStats() error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	fmt.Println("📊 Project Project Statistics")
	fmt.Println("=================================")
	fmt.Println()

	entries, err := os.ReadDir(roadmapDir)
	if err != nil {
		return fmt.Errorf("failed to read roadmap: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "milestone-") {
			continue
		}

		milestonePath := filepath.Join(roadmapDir, entry.Name())
		stats := countTasksInMilestone(milestonePath)

		if stats["total"] == 0 {
			continue
		}

		completion := (stats["complete"] * 100) / stats["total"]

		fmt.Printf("\n%s (%d%% complete)\n", entry.Name(), completion)
		fmt.Println(strings.Repeat("-", len(entry.Name())+20))
		fmt.Printf("  Total:       %d\n", stats["total"])
		fmt.Printf("  Complete:    %d\n", stats["complete"])
		fmt.Printf("  In Progress: %d\n", stats["in_progress"])
		fmt.Printf("  Ready:       %d\n", stats["ready"])
		fmt.Printf("  Pending:     %d\n", stats["pending"])
		fmt.Printf("  Blocked:     %d\n", stats["blocked"])
		fmt.Printf("  Needs Review: %d\n", stats["needs_review"])
		fmt.Printf("  Needs Testing: %d\n", stats["needs_testing"])
	}

	return nil
}

func listBlockers() error {
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return fmt.Errorf("roadmap directory not found")
	}

	fmt.Println("🔴 Blocked Tasks")
	fmt.Println("================")
	fmt.Println()

	blockedCount := 0

	err := filepath.Walk(roadmapDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, "_blocked.md") {
			return nil
		}

		blockedCount++
		relativePath := strings.TrimPrefix(path, roadmapDir+"/")
		filename := filepath.Base(path)

		// Extract task number
		taskNum := "???"
		if len(filename) >= 3 {
			taskNum = filename[:3]
		}

		// Get milestone from path
		parts := strings.Split(relativePath, "/")
		milestone := "unknown"
		if len(parts) > 0 {
			milestone = parts[0]
		}

		fmt.Printf("[%s] Task %s - %s\n", milestone, taskNum, filename)
		fmt.Printf("    Path: %s\n", path)
		fmt.Println()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan for blockers: %w", err)
	}

	if blockedCount == 0 {
		fmt.Println("✅ No blocked tasks!")
	} else {
		fmt.Printf("Total blocked tasks: %d\n", blockedCount)
		fmt.Println("\nUse 'tui-styles-cli comment list <task>' to see blocker details")
	}

	return nil
}

func countTasksInMilestone(milestonePath string) map[string]int {
	stats := map[string]int{
		"total":         0,
		"complete":      0,
		"in_progress":   0,
		"blocked":       0,
		"ready":         0,
		"pending":       0,
		"needs_review":  0,
		"needs_testing": 0,
	}

	filepath.Walk(milestonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)

		// Only count task files (NNN_*.md pattern)
		if len(filename) < 3 || filename[0] < '0' || filename[0] > '9' {
			return nil
		}

		stats["total"]++

		if strings.Contains(filename, "_complete.md") {
			stats["complete"]++
		} else if strings.Contains(filename, "_in_progress.md") {
			stats["in_progress"]++
		} else if strings.Contains(filename, "_blocked.md") {
			stats["blocked"]++
		} else if strings.Contains(filename, "_ready.md") {
			stats["ready"]++
		} else if strings.Contains(filename, "_pending.md") {
			stats["pending"]++
		} else if strings.Contains(filename, "_needs_review.md") {
			stats["needs_review"]++
		} else if strings.Contains(filename, "_needs_testing.md") {
			stats["needs_testing"]++
		}

		return nil
	})

	return stats
}

func percentage(part, total int) int {
	if total == 0 {
		return 0
	}
	return (part * 100) / total
}
