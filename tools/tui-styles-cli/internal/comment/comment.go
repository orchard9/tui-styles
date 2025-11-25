package comment

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/filelock"
	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/task"
	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
	"gopkg.in/yaml.v3"
)

// CreateComment creates a flagged comment on a task file
func CreateComment(taskFileOrID, commentText, author string, needsAddressing bool) error {
	// Find task file if ID provided
	taskFile := taskFileOrID
	if !strings.HasSuffix(taskFileOrID, ".md") {
		found, err := task.FindTaskFile(taskFileOrID)
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

	// Create comment with frontmatter
	comment := types.Comment{
		Author:          author,
		Timestamp:       time.Now(),
		NeedsAddressing: needsAddressing,
		Content:         commentText,
	}

	// Generate comment block
	commentBlock := generateCommentBlock(&comment)

	// Append comment to file
	f, err := os.OpenFile(taskFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open task file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// Add separator and comment
	writer.WriteString("\n---\n\n")
	writer.WriteString(commentBlock)

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to write comment: %w", err)
	}

	status := "✓"
	if needsAddressing {
		status = "⚠"
	}
	fmt.Printf("%s Comment added to %s\n", status, filepath.Base(taskFile))
	fmt.Printf("  Author: %s\n", author)
	fmt.Printf("  Needs Addressing: %v\n", needsAddressing)

	return nil
}

// generateCommentBlock creates a formatted comment block with YAML frontmatter
func generateCommentBlock(comment *types.Comment) string {
	var builder strings.Builder

	// Write comment frontmatter
	builder.WriteString("<!--\n")

	frontmatter := map[string]interface{}{
		"author":           comment.Author,
		"timestamp":        comment.Timestamp.Format(time.RFC3339),
		"needs_addressing": comment.NeedsAddressing,
	}

	yamlData, _ := yaml.Marshal(frontmatter)
	builder.WriteString(string(yamlData))
	builder.WriteString("-->\n\n")

	// Write comment content
	builder.WriteString(comment.Content)
	builder.WriteString("\n")

	return builder.String()
}

// ListComments lists all comments in a task file
func ListComments(taskFileOrID string) error {
	// Find task file if ID provided
	taskFile := taskFileOrID
	if !strings.HasSuffix(taskFileOrID, ".md") {
		found, err := task.FindTaskFile(taskFileOrID)
		if err != nil {
			return err
		}
		taskFile = found
	}

	content, err := os.ReadFile(taskFile)
	if err != nil {
		return fmt.Errorf("failed to read task file: %w", err)
	}

	comments := parseComments(string(content))

	if len(comments) == 0 {
		fmt.Printf("No comments found in %s\n", filepath.Base(taskFile))
		return nil
	}

	fmt.Printf("Comments in %s:\n\n", filepath.Base(taskFile))
	for i, comment := range comments {
		status := "✓"
		if comment.NeedsAddressing {
			status = "⚠ NEEDS ADDRESSING"
		}
		fmt.Printf("[%d] %s\n", i+1, status)
		fmt.Printf("    Author: %s\n", comment.Author)
		fmt.Printf("    Time: %s\n", comment.Timestamp.Format(time.RFC3339))
		fmt.Printf("    Content:\n")
		// Indent content
		for _, line := range strings.Split(strings.TrimSpace(comment.Content), "\n") {
			fmt.Printf("      %s\n", line)
		}
		fmt.Println()
	}

	return nil
}

// parseComments extracts all comments from task file content
func parseComments(content string) []types.Comment {
	var comments []types.Comment

	// Split by HTML comment blocks
	parts := strings.Split(content, "<!--\n")

	for _, part := range parts[1:] { // Skip first part (before first comment)
		endIdx := strings.Index(part, "-->")
		if endIdx == -1 {
			continue
		}

		// Extract YAML frontmatter
		yamlContent := part[:endIdx]
		var frontmatter map[string]interface{}
		if err := yaml.Unmarshal([]byte(yamlContent), &frontmatter); err != nil {
			continue
		}

		// Extract comment content (after -->)
		remaining := part[endIdx+3:]
		nextComment := strings.Index(remaining, "<!--")
		var contentText string
		if nextComment == -1 {
			contentText = strings.TrimSpace(remaining)
		} else {
			contentText = strings.TrimSpace(remaining[:nextComment])
		}

		// Build comment
		comment := types.Comment{
			Content: contentText,
		}

		if author, ok := frontmatter["author"].(string); ok {
			comment.Author = author
		}

		if needsAddr, ok := frontmatter["needs_addressing"].(bool); ok {
			comment.NeedsAddressing = needsAddr
		}

		if timestamp, ok := frontmatter["timestamp"].(string); ok {
			if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
				comment.Timestamp = t
			}
		}

		comments = append(comments, comment)
	}

	return comments
}
