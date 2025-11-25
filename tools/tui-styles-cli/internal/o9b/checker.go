package o9b

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
)

// CheckProject analyzes the current project for O9B components
func CheckProject(o9bPath string) ([]types.ComponentCheck, error) {
	var checks []types.ComponentCheck

	// Get manifest or detect components
	manifest, err := LoadManifest()
	if err != nil {
		// No manifest, detect components
		manifest, _ = CreateManifest("", "")
	}

	// Check CLI
	checks = append(checks, checkCLI(manifest, o9bPath))

	// Check docs
	checks = append(checks, checkDocs(manifest, o9bPath)...)

	// Check Claude
	checks = append(checks, checkClaude(manifest, o9bPath))

	// Check quality configs
	checks = append(checks, checkQualityConfigs(manifest, o9bPath))

	// Check reference code
	checks = append(checks, checkReferenceCode(manifest, o9bPath))

	return checks, nil
}

func checkCLI(manifest *types.O9BManifest, o9bPath string) types.ComponentCheck {
	cliPath := fmt.Sprintf("tools/%s-cli", manifest.ProjectSlug)
	check := types.ComponentCheck{
		Name: "Project CLI",
		Path: cliPath,
	}

	if _, err := os.Stat(cliPath); err == nil {
		check.Installed = true
		check.Version = manifest.Components.CLI.Version
		check.Status = "ok"
		check.Message = fmt.Sprintf("✅ Installed at %s", cliPath)
	} else {
		check.Installed = false
		check.Status = "missing"
		check.Message = fmt.Sprintf("❌ NOT INSTALLED - would be at %s", cliPath)
	}

	return check
}

func checkDocs(manifest *types.O9BManifest, o9bPath string) []types.ComponentCheck {
	var checks []types.ComponentCheck

	docs := []string{"CLAUDE.md", "DESIGN_SYSTEM.md", "CODING_GUIDELINES.md"}
	for _, doc := range docs {
		check := types.ComponentCheck{
			Name: doc,
			Path: doc,
		}

		if _, err := os.Stat(doc); err == nil {
			check.Installed = true
			if docStatus, ok := manifest.Components.Docs[doc]; ok {
				check.Version = docStatus.Version
			}
			check.Status = "ok"
			check.Message = "✅ Installed"
		} else {
			check.Installed = false
			check.Status = "missing"
			check.Message = "❌ MISSING"
		}

		checks = append(checks, check)
	}

	return checks
}

func checkClaude(manifest *types.O9BManifest, o9bPath string) types.ComponentCheck {
	check := types.ComponentCheck{
		Name: "Claude Code Integration",
		Path: ".claude/",
	}

	if stat, err := os.Stat(".claude"); err == nil && stat.IsDir() {
		// Count agents
		agentCount := 0
		agentPath := filepath.Join(".claude", "agents")
		if entries, err := os.ReadDir(agentPath); err == nil {
			for _, entry := range entries {
				if filepath.Ext(entry.Name()) == ".md" {
					agentCount++
				}
			}
		}

		check.Installed = true
		check.Version = manifest.Components.Claude.Version
		check.Status = "ok"
		check.Message = fmt.Sprintf("✅ Installed with %d agents", agentCount)
	} else {
		check.Installed = false
		check.Status = "missing"
		check.Message = "❌ NOT INSTALLED"
	}

	return check
}

func checkQualityConfigs(manifest *types.O9BManifest, o9bPath string) types.ComponentCheck {
	check := types.ComponentCheck{
		Name: "Quality Configs",
		Path: ".",
	}

	configs := map[string]string{
		"eslint.config.mjs": "ESLint",
		".prettierrc.json":  "Prettier",
		".golangci.yml":     "golangci-lint",
		"clippy.toml":       "Clippy",
	}

	var found []string
	var missing []string

	for file, name := range configs {
		if _, err := os.Stat(file); err == nil {
			found = append(found, name)
		} else {
			missing = append(missing, name)
		}
	}

	if len(found) > 0 {
		check.Installed = true
		check.Version = manifest.Components.QualityConfigs.Version
		check.Status = "partial"
		if len(missing) == 0 {
			check.Status = "ok"
			check.Message = fmt.Sprintf("✅ All configs installed: %v", found)
		} else {
			check.Message = fmt.Sprintf("⚠️  Partial - Found: %v, Missing: %v", found, missing)
		}
	} else {
		check.Installed = false
		check.Status = "missing"
		check.Message = "❌ NOT INSTALLED"
	}

	return check
}

func checkReferenceCode(manifest *types.O9BManifest, o9bPath string) types.ComponentCheck {
	check := types.ComponentCheck{
		Name: "Reference Code",
		Path: "reference-code/",
	}

	if stat, err := os.Stat("reference-code"); err == nil && stat.IsDir() {
		// Count reference projects
		refCount := 0
		if entries, err := os.ReadDir("reference-code"); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					refCount++
				}
			}
		}

		check.Installed = true
		check.Version = manifest.Components.ReferenceCode.Version
		check.Status = "ok"
		check.Message = fmt.Sprintf("✅ Installed with %d reference projects", refCount)
	} else {
		check.Installed = false
		check.Status = "missing"
		check.Message = "❌ NOT INSTALLED"
	}

	return check
}
