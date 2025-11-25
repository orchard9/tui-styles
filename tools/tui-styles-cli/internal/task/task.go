package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/filelock"
	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
	"gopkg.in/yaml.v3"
)

var (
	// Task filename pattern: NNN_task_name_status.md
	taskFilePattern = regexp.MustCompile(`^(\d{3})_(.+)_(pending|ready|in_progress|blocked|complete|needs_testing|needs_review)\.md$`)
)

// FindTaskFile finds a task file by ID in the roadmap directory
func FindTaskFile(taskID string) (string, error) {
	// Pad task ID to 3 digits if needed
	if len(taskID) < 3 {
		taskID = fmt.Sprintf("%03s", taskID)
	}

	// Search in roadmap directory
	roadmapDir := "roadmap"
	if _, err := os.Stat(roadmapDir); os.IsNotExist(err) {
		return "", fmt.Errorf("roadmap directory not found")
	}

	var foundFile string
	err := filepath.Walk(roadmapDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		if strings.HasPrefix(filename, taskID+"_") && strings.HasSuffix(filename, ".md") {
			foundFile = path
			return filepath.SkipAll // Stop searching
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundFile == "" {
		return "", fmt.Errorf("task file not found for ID: %s", taskID)
	}

	return foundFile, nil
}

// ParseTaskFile parses a task file and extracts frontmatter (supports both YAML and markdown formats)
func ParseTaskFile(filePath string) (*types.TaskFrontmatter, string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read task file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Check if file uses YAML frontmatter format
	if len(lines) >= 3 && lines[0] == "---" {
		return parseYAMLFrontmatter(lines)
	}

	// Otherwise, parse as markdown format
	return parseMarkdownFormat(filePath, lines)
}

// parseYAMLFrontmatter parses YAML frontmatter format
func parseYAMLFrontmatter(lines []string) (*types.TaskFrontmatter, string, error) {
	// Find end of frontmatter
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, "", fmt.Errorf("invalid task file format: unclosed frontmatter")
	}

	// Parse YAML frontmatter
	frontmatterYAML := strings.Join(lines[1:endIdx], "\n")
	var frontmatter types.TaskFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
		return nil, "", fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Get body content (everything after frontmatter)
	bodyContent := strings.Join(lines[endIdx+1:], "\n")

	return &frontmatter, bodyContent, nil
}

// parseMarkdownFormat parses markdown header format used by existing tasks
func parseMarkdownFormat(filePath string, lines []string) (*types.TaskFrontmatter, string, error) {
	frontmatter := &types.TaskFrontmatter{}
	bodyStartIdx := 0

	// Extract task ID from filename
	filename := filepath.Base(filePath)
	matches := taskFilePattern.FindStringSubmatch(filename)
	if len(matches) >= 2 {
		frontmatter.TaskID = matches[1]
	}

	// Extract status from filename
	if len(matches) >= 4 {
		frontmatter.Status = types.TaskStatus(matches[3])
	}

	// Parse markdown headers for metadata
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Extract title from first heading
		if strings.HasPrefix(line, "# Task ") && frontmatter.Title == "" {
			frontmatter.Title = strings.TrimPrefix(line, "# Task ")
			// Remove task number prefix if present
			if idx := strings.Index(frontmatter.Title, ": "); idx != -1 {
				frontmatter.Title = frontmatter.Title[idx+2:]
			}
		}

		// Parse **Key**: value format
		if strings.HasPrefix(line, "**Status**:") {
			statusVal := strings.TrimSpace(strings.TrimPrefix(line, "**Status**:"))
			if statusVal != "" {
				frontmatter.Status = types.TaskStatus(statusVal)
			}
		} else if strings.HasPrefix(line, "**Confidence**:") {
			confVal := strings.TrimSpace(strings.TrimPrefix(line, "**Confidence**:"))
			if confVal != "" && confVal != "TBD" && !strings.Contains(confVal, "assign during") {
				// Extract percentage if present
				var conf int
				fmt.Sscanf(confVal, "%d", &conf)
				frontmatter.Confidence = conf
			}
		} else if strings.HasPrefix(line, "**Assigned Agent**:") || strings.HasPrefix(line, "**Assigned To**:") {
			assignedVal := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "**Assigned Agent**:"), "**Assigned To**:"))
			if assignedVal != "" && assignedVal != "none" {
				frontmatter.AssignedTo = assignedVal
			}
		}

		// Stop parsing headers when we hit the first major section
		if strings.HasPrefix(line, "## Purpose") || strings.HasPrefix(line, "## Acceptance Criteria") {
			bodyStartIdx = i
			break
		}
	}

	// Get body content (everything from first section onwards)
	bodyContent := strings.Join(lines[bodyStartIdx:], "\n")

	return frontmatter, bodyContent, nil
}

// UpdateTaskStatus updates the status of a task file
func UpdateTaskStatus(taskFileOrID string, newStatus types.TaskStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid status: %s (valid: %v)", newStatus, types.ValidStatuses())
	}

	// Find task file if ID provided
	taskFile := taskFileOrID
	if !strings.HasSuffix(taskFileOrID, ".md") {
		found, err := FindTaskFile(taskFileOrID)
		if err != nil {
			return err
		}
		taskFile = found
	}

	// Acquire lock
	lock := filelock.NewLock(taskFile)
	if err := lock.Acquire(5 * time.Second); err != nil {
		return err
	}
	defer lock.Release()

	// Read original file content
	originalContent, err := os.ReadFile(taskFile)
	if err != nil {
		return fmt.Errorf("failed to read task file: %w", err)
	}

	// Parse existing file
	frontmatter, bodyContent, err := ParseTaskFile(taskFile)
	if err != nil {
		return err
	}

	// Update status and timestamp
	oldStatus := frontmatter.Status
	frontmatter.Status = newStatus
	frontmatter.UpdatedAt = time.Now()

	// Generate new filename
	filename := filepath.Base(taskFile)
	dir := filepath.Dir(taskFile)

	matches := taskFilePattern.FindStringSubmatch(filename)
	if len(matches) < 4 {
		return fmt.Errorf("invalid task filename format: %s", filename)
	}

	taskNum := matches[1]
	taskName := matches[2]
	newFilename := fmt.Sprintf("%s_%s_%s.md", taskNum, taskName, newStatus)
	newPath := filepath.Join(dir, newFilename)

	// Determine if original file uses YAML frontmatter
	usesYAMLFrontmatter := strings.HasPrefix(string(originalContent), "---\n")

	// Write updated content to new file
	if err := writeTaskFile(newPath, frontmatter, bodyContent, usesYAMLFrontmatter); err != nil {
		return err
	}

	// Remove old file if different
	if taskFile != newPath {
		if err := os.Remove(taskFile); err != nil {
			// Try to remove new file to maintain consistency
			os.Remove(newPath)
			return fmt.Errorf("failed to remove old task file: %w", err)
		}
	}

	fmt.Printf("✓ Task %s status updated: %s → %s\n", frontmatter.TaskID, oldStatus, newStatus)
	fmt.Printf("  Renamed: %s → %s\n", filepath.Base(taskFile), newFilename)

	return nil
}

// GetTaskInfo retrieves and displays task information
func GetTaskInfo(taskFileOrID string) error {
	// Find task file if ID provided
	taskFile := taskFileOrID
	if !strings.HasSuffix(taskFileOrID, ".md") {
		found, err := FindTaskFile(taskFileOrID)
		if err != nil {
			return err
		}
		taskFile = found
	}

	frontmatter, _, err := ParseTaskFile(taskFile)
	if err != nil {
		return err
	}

	// Display task information
	fmt.Printf("Task ID: %s\n", frontmatter.TaskID)
	fmt.Printf("Title: %s\n", frontmatter.Title)
	fmt.Printf("Status: %s\n", frontmatter.Status)
	if frontmatter.Confidence > 0 {
		fmt.Printf("Confidence: %d%%\n", frontmatter.Confidence)
	}
	if frontmatter.AssignedTo != "" {
		fmt.Printf("Assigned To: %s\n", frontmatter.AssignedTo)
	}
	if len(frontmatter.Dependencies) > 0 {
		fmt.Printf("Dependencies: %v\n", frontmatter.Dependencies)
	}
	if !frontmatter.CreatedAt.IsZero() {
		fmt.Printf("Created: %s\n", frontmatter.CreatedAt.Format(time.RFC3339))
	}
	if !frontmatter.UpdatedAt.IsZero() {
		fmt.Printf("Updated: %s\n", frontmatter.UpdatedAt.Format(time.RFC3339))
	}
	fmt.Printf("File: %s\n", taskFile)

	return nil
}

// writeTaskFile writes a task file with frontmatter and body
func writeTaskFile(filePath string, frontmatter *types.TaskFrontmatter, bodyContent string, useYAMLFrontmatter bool) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create task file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	if useYAMLFrontmatter {
		// Write YAML frontmatter format
		writer.WriteString("---\n")
		yamlData, err := yaml.Marshal(frontmatter)
		if err != nil {
			return fmt.Errorf("failed to marshal frontmatter: %w", err)
		}
		writer.Write(yamlData)
		writer.WriteString("---\n")
		writer.WriteString(bodyContent)
	} else {
		// Write markdown format (update status inline)
		lines := strings.Split(bodyContent, "\n")
		headerWritten := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			// Update status line
			if strings.HasPrefix(trimmed, "**Status**:") {
				writer.WriteString(fmt.Sprintf("**Status**: %s\n", frontmatter.Status))
				headerWritten = true
			} else if strings.HasPrefix(trimmed, "**Confidence**:") && frontmatter.Confidence > 0 {
				writer.WriteString(fmt.Sprintf("**Confidence**: %d%%\n", frontmatter.Confidence))
				headerWritten = true
			} else {
				writer.WriteString(line + "\n")
				if !headerWritten && strings.HasPrefix(trimmed, "**Status**") {
					headerWritten = true
				}
			}
		}
	}

	return writer.Flush()
}
