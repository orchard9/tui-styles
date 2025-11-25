package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orchard9/tui-styles/tools/tui-styles-cli/internal/o9b"
	"github.com/spf13/cobra"
)

var o9bCmd = &cobra.Command{
	Use:   "o9b",
	Short: "O9B project upgrade and management commands",
	Long: `Commands for managing O9B components in your project.

These commands help you:
  - Check what O9B components are installed
  - Install missing components
  - Upgrade existing components
  - Track O9B versions`,
}

var o9bCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check what O9B components are installed",
	Long: `Analyze the current project and report which O9B components are installed.

This will check for:
  - Project CLI (tools/<project>-cli)
  - Documentation (CLAUDE.md, DESIGN_SYSTEM.md, CODING_GUIDELINES.md)
  - Claude Code integration (.claude/)
  - Quality configs (eslint, prettier, golangci-lint, clippy)
  - Reference code (reference-code/)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		o9bPath, _ := cmd.Flags().GetString("o9b-path")
		return runO9BCheck(o9bPath)
	},
}

var o9bSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Install or upgrade O9B components",
	Long: `Install missing O9B components or upgrade existing ones.

Examples:
  # Install specific component
  tui-styles-cli o9b sync --component cli

  # Install all missing components
  tui-styles-cli o9b sync --all

  # Dry run to see what would be installed
  tui-styles-cli o9b sync --all --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		o9bPath, _ := cmd.Flags().GetString("o9b-path")
		component, _ := cmd.Flags().GetString("component")
		all, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		projectName, _ := cmd.Flags().GetString("project-name")
		projectSlug, _ := cmd.Flags().GetString("project-slug")

		return runO9BSync(o9bPath, component, all, dryRun, projectName, projectSlug)
	},
}

var o9bInitCmd = &cobra.Command{
	Use:   "init <project-name> <project-slug>",
	Short: "Initialize O9B manifest for existing project",
	Long: `Create an O9B manifest file by detecting existing components.

This is useful for projects that were created before O9B manifest tracking
or were partially scaffolded.

Example:
  tui-styles-cli o9b init "Peach Platform" "peach"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runO9BInit(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(o9bCmd)

	// Add subcommands
	o9bCmd.AddCommand(o9bCheckCmd)
	o9bCmd.AddCommand(o9bSyncCmd)
	o9bCmd.AddCommand(o9bInitCmd)

	// Flags for check
	o9bCheckCmd.Flags().String("o9b-path", "", "Path to O9B directory (default: auto-detect)")

	// Flags for sync
	o9bSyncCmd.Flags().String("o9b-path", "", "Path to O9B directory (default: auto-detect)")
	o9bSyncCmd.Flags().String("component", "", "Specific component to sync (cli, docs, claude, quality-configs, reference-code)")
	o9bSyncCmd.Flags().Bool("all", false, "Sync all components")
	o9bSyncCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
	o9bSyncCmd.Flags().String("project-name", "", "Project name (required if no manifest exists)")
	o9bSyncCmd.Flags().String("project-slug", "", "Project slug (required if no manifest exists)")
}

func runO9BCheck(o9bPath string) error {
	// Auto-detect O9B path if not provided
	if o9bPath == "" {
		detected, err := detectO9BPath()
		if err != nil {
			return fmt.Errorf("could not detect O9B path: %w (use --o9b-path to specify)", err)
		}
		o9bPath = detected
	}

	// Verify O9B path exists
	if _, err := os.Stat(o9bPath); err != nil {
		return fmt.Errorf("O9B path not found: %s", o9bPath)
	}

	// Run checks
	checks, err := o9b.CheckProject(o9bPath)
	if err != nil {
		return fmt.Errorf("failed to check project: %w", err)
	}

	// Load manifest for project info
	manifest, _ := o9b.LoadManifest()
	projectName := "Current Project"
	if manifest != nil {
		projectName = manifest.ProjectName
	}

	// Display results
	fmt.Printf("🔍 O9B Component Check - %s\n", projectName)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	installedCount := 0
	missingCount := 0

	for _, check := range checks {
		fmt.Printf("%-30s %s\n", check.Name+":", check.Message)
		if check.Installed {
			installedCount++
		} else {
			missingCount++
		}
	}

	// Summary
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Summary: %d installed, %d missing\n", installedCount, missingCount)

	if missingCount > 0 {
		fmt.Println()
		fmt.Println("💡 To install missing components:")
		fmt.Println("   tui-styles-cli o9b sync --all")
		fmt.Println()
		fmt.Println("   Or install specific components:")
		fmt.Println("   tui-styles-cli o9b sync --component cli")
		fmt.Println("   tui-styles-cli o9b sync --component docs")
		fmt.Println("   tui-styles-cli o9b sync --component claude")
	}

	return nil
}

func runO9BSync(o9bPath, component string, all, dryRun bool, projectName, projectSlug string) error {
	// Auto-detect O9B path if not provided
	if o9bPath == "" {
		detected, err := detectO9BPath()
		if err != nil {
			return fmt.Errorf("could not detect O9B path: %w (use --o9b-path to specify)", err)
		}
		o9bPath = detected
	}

	// Verify O9B path exists
	if _, err := os.Stat(o9bPath); err != nil {
		return fmt.Errorf("O9B path not found: %s", o9bPath)
	}

	// Check if manifest exists, if not require project name/slug
	manifest, err := o9b.LoadManifest()
	if err != nil {
		if projectName == "" || projectSlug == "" {
			return fmt.Errorf("no manifest found - please provide --project-name and --project-slug")
		}
	} else {
		// Use manifest values if not provided
		if projectName == "" {
			projectName = manifest.ProjectName
		}
		if projectSlug == "" {
			projectSlug = manifest.ProjectSlug
		}
	}

	// Validate flags
	if !all && component == "" {
		return fmt.Errorf("must specify either --component or --all")
	}

	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No changes will be made")
		fmt.Println()
	}

	// Run sync
	opts := o9b.SyncOptions{
		O9BPath:     o9bPath,
		Component:   component,
		All:         all,
		DryRun:      dryRun,
		ProjectName: projectName,
		ProjectSlug: projectSlug,
	}

	fmt.Printf("🔄 Syncing O9B components for %s\n", projectName)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	if err := o9b.SyncComponents(opts); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✅ Sync complete!")

	if !dryRun {
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   1. Review installed components")
		fmt.Println("   2. Customize templates for your project")
		fmt.Println("   3. Build CLI if installed: cd tools/<project>-cli && make build")
	}

	return nil
}

func runO9BInit(projectName, projectSlug string) error {
	// Check if manifest already exists
	if _, err := os.Stat(o9b.ManifestFile); err == nil {
		return fmt.Errorf("manifest already exists at %s", o9b.ManifestFile)
	}

	fmt.Printf("🔧 Creating O9B manifest for %s (%s)\n", projectName, projectSlug)
	fmt.Println()

	// Create manifest
	manifest, err := o9b.CreateManifest(projectName, projectSlug)
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}

	// Save manifest
	if err := o9b.SaveManifest(manifest); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("✅ Manifest created at %s\n", o9b.ManifestFile)
	fmt.Println()
	fmt.Println("Detected components:")

	// Show what was detected
	if manifest.Components.CLI.Installed {
		fmt.Println("  ✅ CLI")
	}
	for doc, status := range manifest.Components.Docs {
		if status.Installed {
			fmt.Printf("  ✅ %s\n", doc)
		}
	}
	if manifest.Components.Claude.Installed {
		fmt.Println("  ✅ Claude Code")
	}
	if manifest.Components.QualityConfigs.Installed {
		fmt.Println("  ✅ Quality Configs")
	}
	if manifest.Components.ReferenceCode.Installed {
		fmt.Println("  ✅ Reference Code")
	}

	fmt.Println()
	fmt.Println("💡 Run 'tui-styles-cli o9b check' to see detailed component status")

	return nil
}

// detectO9BPath attempts to auto-detect the O9B directory
func detectO9BPath() (string, error) {
	// Common locations to check
	possiblePaths := []string{
		"/Users/jordanwashburn/Workspace/orchard9/o9b",
		"../o9b",
		"../../o9b",
		filepath.Join(os.Getenv("HOME"), "Workspace/orchard9/o9b"),
	}

	for _, path := range possiblePaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				// Verify it's actually O9B by checking for key files
				if _, err := os.Stat(filepath.Join(absPath, "VERSION")); err == nil {
					return absPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("O9B directory not found in common locations")
}
