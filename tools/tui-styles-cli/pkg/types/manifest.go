package types

import "time"

// O9BManifest tracks what O9B components are installed in a project
type O9BManifest struct {
	O9BVersion  string             `yaml:"o9b_version"`
	ProjectName string             `yaml:"project_name"`
	ProjectSlug string             `yaml:"project_slug"`
	Created     time.Time          `yaml:"created"`
	LastUpdated time.Time          `yaml:"last_updated"`
	Components  ManifestComponents `yaml:"components"`
}

// ManifestComponents tracks individual component status
type ManifestComponents struct {
	CLI            ComponentStatus            `yaml:"cli"`
	Docs           map[string]ComponentStatus `yaml:"docs"`
	Claude         ComponentStatus            `yaml:"claude"`
	QualityConfigs ComponentStatus            `yaml:"quality_configs"`
	ReferenceCode  ComponentStatus            `yaml:"reference_code"`
}

// ComponentStatus tracks if a component is installed and its version
type ComponentStatus struct {
	Installed bool   `yaml:"installed"`
	Version   string `yaml:"version,omitempty"`
}

// ComponentType represents types of O9B components
type ComponentType string

const (
	ComponentCLI           ComponentType = "cli"
	ComponentDocs          ComponentType = "docs"
	ComponentClaude        ComponentType = "claude"
	ComponentQualityConfig ComponentType = "quality-configs"
	ComponentReferenceCode ComponentType = "reference-code"
)

// ComponentCheck represents the result of checking a component
type ComponentCheck struct {
	Name      string
	Installed bool
	Version   string
	Path      string
	Status    string // "ok", "missing", "outdated"
	Message   string
}
