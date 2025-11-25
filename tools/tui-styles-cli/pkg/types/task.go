package types

import "time"

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	StatusPending                TaskStatus = "pending"
	StatusReady                  TaskStatus = "ready"
	StatusInProgress             TaskStatus = "in_progress"
	StatusBlocked                TaskStatus = "blocked"
	StatusComplete               TaskStatus = "complete"
	StatusNeedsTesting           TaskStatus = "needs_testing"
	StatusNeedsReview            TaskStatus = "needs_review"
	StatusNeedsHumanVerification TaskStatus = "needs_human_verification"
)

// ValidStatuses returns all valid task statuses
func ValidStatuses() []TaskStatus {
	return []TaskStatus{
		StatusPending,
		StatusReady,
		StatusInProgress,
		StatusBlocked,
		StatusComplete,
		StatusNeedsTesting,
		StatusNeedsReview,
		StatusNeedsHumanVerification,
	}
}

// IsValid checks if a status is valid
func (s TaskStatus) IsValid() bool {
	for _, valid := range ValidStatuses() {
		if s == valid {
			return true
		}
	}
	return false
}

// TaskFrontmatter represents the YAML frontmatter of a task file
type TaskFrontmatter struct {
	TaskID       string     `yaml:"task_id"`
	Title        string     `yaml:"title"`
	Status       TaskStatus `yaml:"status"`
	Confidence   int        `yaml:"confidence,omitempty"`
	AssignedTo   string     `yaml:"assigned_to,omitempty"`
	Dependencies []string   `yaml:"dependencies,omitempty"`
	CreatedAt    time.Time  `yaml:"created_at,omitempty"`
	UpdatedAt    time.Time  `yaml:"updated_at,omitempty"`
}

// Comment represents a flagged comment on a task
type Comment struct {
	Author          string    `yaml:"author"`
	Timestamp       time.Time `yaml:"timestamp"`
	NeedsAddressing bool      `yaml:"needs_addressing"`
	Content         string    `yaml:"-"` // Content comes after frontmatter
}
