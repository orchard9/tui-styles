package o9b

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/pkg/types"
)

// SyncOptions contains options for syncing components
type SyncOptions struct {
	O9BPath     string
	Component   string
	All         bool
	DryRun      bool
	ProjectName string
	ProjectSlug string
}

// SyncComponents synchronizes O9B components to the project
func SyncComponents(opts SyncOptions) error {
	// Load or create manifest
	manifest, err := LoadManifest()
	if err != nil {
		// Create new manifest
		manifest, err = CreateManifest(opts.ProjectName, opts.ProjectSlug)
		if err != nil {
			return fmt.Errorf("failed to create manifest: %w", err)
		}
		if !opts.DryRun {
			if err := SaveManifest(manifest); err != nil {
				return fmt.Errorf("failed to save manifest: %w", err)
			}
			fmt.Println("✅ Created O9B manifest")
		}
	}

	// Determine what to sync
	var componentsToSync []types.ComponentType
	if opts.All {
		componentsToSync = []types.ComponentType{
			types.ComponentCLI,
			types.ComponentDocs,
			types.ComponentClaude,
			types.ComponentQualityConfig,
			types.ComponentReferenceCode,
		}
	} else if opts.Component != "" {
		componentsToSync = []types.ComponentType{types.ComponentType(opts.Component)}
	} else {
		return fmt.Errorf("must specify --component or --all")
	}

	// Sync each component
	for _, component := range componentsToSync {
		if err := syncComponent(component, manifest, opts); err != nil {
			return fmt.Errorf("failed to sync %s: %w", component, err)
		}
	}

	// Update manifest
	if !opts.DryRun {
		manifest.Components.CLI = detectCLI(manifest.ProjectSlug)
		manifest.Components.Docs = detectDocs()
		manifest.Components.Claude = detectClaude()
		manifest.Components.QualityConfigs = detectQualityConfigs()
		manifest.Components.ReferenceCode = detectReferenceCode()

		if err := SaveManifest(manifest); err != nil {
			return fmt.Errorf("failed to save manifest: %w", err)
		}
	}

	return nil
}

func syncComponent(component types.ComponentType, manifest *types.O9BManifest, opts SyncOptions) error {
	switch component {
	case types.ComponentCLI:
		return syncCLI(manifest, opts)
	case types.ComponentDocs:
		return syncDocs(manifest, opts)
	case types.ComponentClaude:
		return syncClaude(manifest, opts)
	case types.ComponentQualityConfig:
		return syncQualityConfigs(manifest, opts)
	case types.ComponentReferenceCode:
		return syncReferenceCode(manifest, opts)
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
}

func syncCLI(manifest *types.O9BManifest, opts SyncOptions) error {
	cliSource := filepath.Join(opts.O9BPath, "tui-styles-cli")
	cliDest := fmt.Sprintf("tools/%s-cli", manifest.ProjectSlug)

	if _, err := os.Stat(cliDest); err == nil {
		fmt.Printf("⚠️  CLI already exists at %s (skipping)\n", cliDest)
		return nil
	}

	fmt.Printf("📦 Installing project CLI to %s\n", cliDest)

	if opts.DryRun {
		fmt.Println("   [DRY RUN] Would copy and configure CLI")
		return nil
	}

	// Create tools directory
	if err := os.MkdirAll("tools", 0755); err != nil {
		return fmt.Errorf("failed to create tools directory: %w", err)
	}

	// Copy tui-styles-cli
	if err := copyDir(cliSource, cliDest); err != nil {
		return fmt.Errorf("failed to copy CLI: %w", err)
	}

	// Update go.mod - need to replace both the base path and the CLI suffix
	goModPath := filepath.Join(cliDest, "go.mod")
	if err := replaceInFile(goModPath, "github.com/orchard9/tui-styles/tools/tui-styles-cli", fmt.Sprintf("github.com/orchard9/%s/tools/%s-cli", manifest.ProjectSlug, manifest.ProjectSlug)); err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// Update import paths in all Go files
	importOldPath := "github.com/orchard9/tui-styles/tools/tui-styles-cli"
	importNewPath := fmt.Sprintf("github.com/orchard9/%s/tools/%s-cli", manifest.ProjectSlug, manifest.ProjectSlug)

	// Update all Go files and docs
	if err := filepath.Walk(cliDest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		baseName := filepath.Base(path)

		// For Go files, update import paths
		if ext == ".go" {
			if err := replaceInFile(path, importOldPath, importNewPath); err != nil {
				return err
			}
		}

		// For all file types, replace tui-styles-cli with {slug}-cli
		if ext == ".go" || ext == ".md" || baseName == "Makefile" {
			if err := replaceInFile(path, "tui-styles-cli", fmt.Sprintf("%s-cli", manifest.ProjectSlug)); err != nil {
				return err
			}
		}

		// Only replace "project" → slug in non-Go files (docs and Makefile)
		// Go files use "project" as variable names which shouldn't be changed
		if ext == ".md" || baseName == "Makefile" {
			// Also replace in command descriptions and help text
			if err := replaceInFile(path, "Project CLI", fmt.Sprintf("%s CLI", capitalizeFirst(manifest.ProjectSlug))); err != nil {
				return err
			}
			if err := replaceInFile(path, "Project Project", manifest.ProjectName); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to update CLI files: %w", err)
	}

	fmt.Printf("✅ CLI installed at %s\n", cliDest)
	fmt.Printf("   Run: cd %s && make build\n", cliDest)

	return nil
}

func syncDocs(manifest *types.O9BManifest, opts SyncOptions) error {
	templates := filepath.Join(opts.O9BPath, "templates")
	docs := []string{"CLAUDE.md", "DESIGN_SYSTEM.md", "CODING_GUIDELINES.md"}

	for _, doc := range docs {
		// Skip if exists
		if _, err := os.Stat(doc); err == nil {
			fmt.Printf("⚠️  %s already exists (skipping)\n", doc)
			continue
		}

		sourcePath := filepath.Join(templates, doc+".template")
		fmt.Printf("📝 Installing %s\n", doc)

		if opts.DryRun {
			fmt.Printf("   [DRY RUN] Would copy %s\n", doc)
			continue
		}

		// Read template
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", doc, err)
		}

		// Replace placeholders
		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_NAME}}", manifest.ProjectName)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_SLUG}}", manifest.ProjectSlug)

		// Write file
		if err := os.WriteFile(doc, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", doc, err)
		}

		fmt.Printf("✅ %s installed\n", doc)
	}

	return nil
}

func syncClaude(manifest *types.O9BManifest, opts SyncOptions) error {
	claudeSource := filepath.Join(opts.O9BPath, ".claude")
	claudeDest := ".claude"

	fmt.Printf("🤖 Syncing Claude Code integration\n")

	// Create .claude directory if it doesn't exist
	if _, err := os.Stat(claudeDest); os.IsNotExist(err) {
		if !opts.DryRun {
			if err := os.MkdirAll(claudeDest, 0755); err != nil {
				return fmt.Errorf("failed to create .claude directory: %w", err)
			}
		}
	}

	// Sync agents/ subdirectory
	if err := syncClaudeSubdir(claudeSource, claudeDest, "agents", manifest, opts); err != nil {
		return fmt.Errorf("failed to sync agents: %w", err)
	}

	// Sync commands/ subdirectory
	if err := syncClaudeSubdir(claudeSource, claudeDest, "commands", manifest, opts); err != nil {
		return fmt.Errorf("failed to sync commands: %w", err)
	}

	// Copy README if it doesn't exist
	readmeSrc := filepath.Join(claudeSource, "README.md.template")
	readmeDest := filepath.Join(claudeDest, "README.md")
	if _, err := os.Stat(readmeDest); os.IsNotExist(err) {
		fmt.Println("   📄 Installing README.md")
		if !opts.DryRun {
			content, err := os.ReadFile(readmeSrc)
			if err == nil {
				contentStr := string(content)
				contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_NAME}}", manifest.ProjectName)
				contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_SLUG}}", manifest.ProjectSlug)
				if err := os.WriteFile(readmeDest, []byte(contentStr), 0644); err != nil {
					return fmt.Errorf("failed to write README: %w", err)
				}
			}
		}
	}

	fmt.Println("✅ Claude Code integration synced")

	return nil
}

// syncClaudeSubdir syncs a subdirectory within .claude (agents/ or commands/)
func syncClaudeSubdir(claudeSource, claudeDest, subdir string, manifest *types.O9BManifest, opts SyncOptions) error {
	sourcePath := filepath.Join(claudeSource, subdir)
	destPath := filepath.Join(claudeDest, subdir)

	// Check if source directory exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil // Skip if source doesn't exist
	}

	// Create destination directory if needed
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if !opts.DryRun {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create %s directory: %w", subdir, err)
			}
		}
	}

	// Read source files
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read %s directory: %w", subdir, err)
	}

	addedCount := 0
	skippedCount := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		srcFile := filepath.Join(sourcePath, entry.Name())

		// Determine destination filename
		destFileName := entry.Name()

		// For agents, rename wsp- prefix to project slug
		if subdir == "agents" && strings.HasPrefix(entry.Name(), "wsp-") {
			destFileName = strings.Replace(entry.Name(), "wsp-", manifest.ProjectSlug+"-", 1)
		}

		destFile := filepath.Join(destPath, destFileName)

		// Skip if file already exists
		if _, err := os.Stat(destFile); err == nil {
			skippedCount++
			continue
		}

		// Read source file
		content, err := os.ReadFile(srcFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		// Replace placeholders
		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_NAME}}", manifest.ProjectName)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_SLUG}}", manifest.ProjectSlug)

		if opts.DryRun {
			fmt.Printf("   [DRY RUN] Would add %s/%s\n", subdir, destFileName)
		} else {
			if err := os.WriteFile(destFile, []byte(contentStr), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", destFileName, err)
			}
			addedCount++
		}
	}

	if addedCount > 0 || opts.DryRun {
		fmt.Printf("   ✅ %s: added %d files", subdir, addedCount)
		if skippedCount > 0 {
			fmt.Printf(" (skipped %d existing)", skippedCount)
		}
		fmt.Println()
	} else if skippedCount > 0 {
		fmt.Printf("   ⚠️  %s: all files already exist (skipped %d)\n", subdir, skippedCount)
	}

	return nil
}

func syncQualityConfigs(manifest *types.O9BManifest, opts SyncOptions) error {
	configsSource := filepath.Join(opts.O9BPath, "configs")

	configs := map[string]string{
		"eslint/eslint.config.mjs":  "eslint.config.mjs",
		"prettier/.prettierrc.json": ".prettierrc.json",
		"go/.golangci.yml":          ".golangci.yml",
		"rust/clippy.toml":          "clippy.toml",
	}

	fmt.Println("🔧 Installing quality configs")

	for source, dest := range configs {
		// Skip if exists
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("   ⚠️  %s already exists (skipping)\n", dest)
			continue
		}

		sourcePath := filepath.Join(configsSource, source)
		if _, err := os.Stat(sourcePath); err != nil {
			// Config doesn't exist in O9B, skip
			continue
		}

		if opts.DryRun {
			fmt.Printf("   [DRY RUN] Would copy %s\n", dest)
			continue
		}

		// Copy file
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", source, err)
		}

		// Replace placeholders
		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_NAME}}", manifest.ProjectName)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_SLUG}}", manifest.ProjectSlug)

		if err := os.WriteFile(dest, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", dest, err)
		}

		fmt.Printf("   ✅ %s installed\n", dest)
	}

	return nil
}

func syncReferenceCode(manifest *types.O9BManifest, opts SyncOptions) error {
	refSource := filepath.Join(opts.O9BPath, "reference-code")
	refDest := "reference-code"

	if _, err := os.Stat(refDest); err == nil {
		fmt.Printf("⚠️  reference-code/ already exists (skipping)\n")
		return nil
	}

	fmt.Printf("📚 Installing reference code\n")

	if opts.DryRun {
		fmt.Println("   [DRY RUN] Would copy reference-code/")
		return nil
	}

	if err := copyDir(refSource, refDest); err != nil {
		return fmt.Errorf("failed to copy reference-code: %w", err)
	}

	fmt.Println("✅ Reference code installed")

	return nil
}

// Helper functions

func copyDir(src, dst string) error {
	// Use cp -r for efficiency
	cmd := exec.Command("cp", "-r", src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cp command failed: %w", err)
	}
	return nil
}

func replaceInFile(path, old, new string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	newContent := strings.ReplaceAll(string(content), old, new)
	return os.WriteFile(path, []byte(newContent), 0644)
}

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	// Handle kebab-case by capitalizing each word
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
