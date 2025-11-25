package o9b

import (
	"fmt"
	"os"
	"time"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
	"gopkg.in/yaml.v3"
)

const ManifestFile = ".o9b-manifest.yaml"
const CurrentO9BVersion = "1.0.0"

// LoadManifest loads the O9B manifest from the current directory
func LoadManifest() (*types.O9BManifest, error) {
	data, err := os.ReadFile(ManifestFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("manifest not found (run 'o9b init' first)")
		}
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest types.O9BManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// SaveManifest saves the O9B manifest to the current directory
func SaveManifest(manifest *types.O9BManifest) error {
	manifest.LastUpdated = time.Now()

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(ManifestFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// CreateManifest creates a new manifest by detecting existing components
func CreateManifest(projectName, projectSlug string) (*types.O9BManifest, error) {
	manifest := &types.O9BManifest{
		O9BVersion:  CurrentO9BVersion,
		ProjectName: projectName,
		ProjectSlug: projectSlug,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Components: types.ManifestComponents{
			Docs: make(map[string]types.ComponentStatus),
		},
	}

	// Detect existing components
	manifest.Components.CLI = detectCLI(projectSlug)
	manifest.Components.Docs = detectDocs()
	manifest.Components.Claude = detectClaude()
	manifest.Components.QualityConfigs = detectQualityConfigs()
	manifest.Components.ReferenceCode = detectReferenceCode()

	return manifest, nil
}

// detectCLI checks if project CLI is installed
func detectCLI(projectSlug string) types.ComponentStatus {
	cliPath := fmt.Sprintf("tools/%s-cli", projectSlug)
	if _, err := os.Stat(cliPath); err == nil {
		return types.ComponentStatus{Installed: true, Version: CurrentO9BVersion}
	}
	return types.ComponentStatus{Installed: false}
}

// detectDocs checks which documentation files exist
func detectDocs() map[string]types.ComponentStatus {
	docs := make(map[string]types.ComponentStatus)

	docFiles := []string{"CLAUDE.md", "DESIGN_SYSTEM.md", "CODING_GUIDELINES.md"}
	for _, doc := range docFiles {
		if _, err := os.Stat(doc); err == nil {
			docs[doc] = types.ComponentStatus{Installed: true, Version: CurrentO9BVersion}
		} else {
			docs[doc] = types.ComponentStatus{Installed: false}
		}
	}

	return docs
}

// detectClaude checks if .claude directory exists
func detectClaude() types.ComponentStatus {
	if _, err := os.Stat(".claude"); err == nil {
		return types.ComponentStatus{Installed: true, Version: CurrentO9BVersion}
	}
	return types.ComponentStatus{Installed: false}
}

// detectQualityConfigs checks if quality config files exist
func detectQualityConfigs() types.ComponentStatus {
	configs := []string{
		"eslint.config.mjs",
		".prettierrc.json",
		".golangci.yml",
		"clippy.toml",
	}

	foundAny := false
	for _, config := range configs {
		if _, err := os.Stat(config); err == nil {
			foundAny = true
			break
		}
	}

	if foundAny {
		return types.ComponentStatus{Installed: true, Version: CurrentO9BVersion}
	}
	return types.ComponentStatus{Installed: false}
}

// detectReferenceCode checks if reference-code directory exists
func detectReferenceCode() types.ComponentStatus {
	if _, err := os.Stat("reference-code"); err == nil {
		return types.ComponentStatus{Installed: true, Version: CurrentO9BVersion}
	}
	return types.ComponentStatus{Installed: false}
}
