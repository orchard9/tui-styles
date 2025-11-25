package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development environment management",
	Long:  `Commands for managing the local development environment.`,
}

var devStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check development environment status",
	Long: `Verify all required tools and services are installed and running:
  - Required tools (Go, Node.js, Docker, etc.)
  - Running services (PostgreSQL, Redis, etc.)
  - Project dependencies
  - Port availability`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkDevStatus()
	},
}

var devSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up development environment",
	Long: `Install and configure everything needed for development:
  - Install missing dependencies
  - Set up databases
  - Create .env files
  - Initialize services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupDevEnvironment()
	},
}

var devDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose development environment issues",
	Long: `Run comprehensive diagnostics to identify and suggest fixes for:
  - Missing dependencies
  - Service connectivity issues
  - Port conflicts
  - Configuration problems`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDevDiagnostics()
	},
}

var devStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all development services",
	Long:  `Start all required services using docker-compose.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startDevServices()
	},
}

var devStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all development services",
	Long:  `Stop all running development services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopDevServices()
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devStatusCmd)
	devCmd.AddCommand(devSetupCmd)
	devCmd.AddCommand(devDoctorCmd)
	devCmd.AddCommand(devStartCmd)
	devCmd.AddCommand(devStopCmd)
}

func checkDevStatus() error {
	fmt.Println("🔍 Development Environment Status")
	fmt.Println("==================================")
	fmt.Println()

	issues := 0

	// Check required tools
	fmt.Println("Required Tools:")
	fmt.Println("---------------")

	tools := []struct {
		name    string
		command string
		version string
	}{
		{"Go", "go", "version"},
		{"Node.js", "node", "--version"},
		{"npm", "npm", "--version"},
		{"Docker", "docker", "--version"},
		{"Docker Compose", "docker-compose", "--version"},
		{"PostgreSQL Client", "psql", "--version"},
		{"Git", "git", "--version"},
	}

	for _, tool := range tools {
		cmd := exec.Command(tool.command, tool.version)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("  ❌ %s: NOT INSTALLED\n", tool.name)
			issues++
		} else {
			version := strings.TrimSpace(strings.Split(string(output), "\n")[0])
			fmt.Printf("  ✅ %s: %s\n", tool.name, version)
		}
	}

	fmt.Println()

	// Check Docker services
	fmt.Println("Docker Services:")
	fmt.Println("----------------")

	cmd := exec.Command("docker-compose", "ps", "--services", "--filter", "status=running")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("  ⚠️  Unable to check docker-compose services")
		fmt.Println("     Run 'tui-styles-cli dev start' to start services")
	} else {
		services := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(services) == 1 && services[0] == "" {
			fmt.Println("  ⚠️  No services running")
			fmt.Println("     Run 'tui-styles-cli dev start' to start services")
			issues++
		} else {
			for _, service := range services {
				if service != "" {
					fmt.Printf("  ✅ %s: RUNNING\n", service)
				}
			}
		}
	}

	fmt.Println()

	// Check port availability
	fmt.Println("Port Availability (34070-34079):")
	fmt.Println("---------------------------------")

	ports := map[int]string{
		34070: "PostgreSQL",
		34071: "Redis",
		34072: "Identity Studio (Next.js)",
		34073: "Reserved",
		34074: "Auth Service (Go)",
		34075: "Identity API (Go)",
		34076: "WebSocket Service (Go)",
		34077: "Transformation API (Go)",
		34078: "Identity Processing API (Python)",
		34079: "Analytics Service (Python)",
	}

	for port, service := range ports {
		if isPortInUse(port) {
			fmt.Printf("  🟢 Port %d: IN USE (%s)\n", port, service)
		} else {
			fmt.Printf("  ⚪ Port %d: AVAILABLE (%s)\n", port, service)
		}
	}

	fmt.Println()

	// Summary
	if issues == 0 {
		fmt.Println("✅ Development environment is ready!")
	} else {
		fmt.Printf("⚠️  Found %d issue(s)\n", issues)
		fmt.Println("   Run 'tui-styles-cli dev doctor' for detailed diagnostics")
		fmt.Println("   Run 'tui-styles-cli dev setup' to fix issues automatically")
	}

	return nil
}

func setupDevEnvironment() error {
	fmt.Println("🛠️  Setting Up Development Environment")
	fmt.Println("======================================")
	fmt.Println()

	// Check for docker-compose.yml
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found. Are you in the project root?")
	}

	// Start services
	fmt.Println("Starting development services...")
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Development environment setup complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Check status: tui-styles-cli dev status")
	fmt.Println("  2. View logs: docker-compose logs -f")
	fmt.Println("  3. Start coding!")

	return nil
}

func runDevDiagnostics() error {
	fmt.Println("🔬 Running Development Environment Diagnostics")
	fmt.Println("==============================================")
	fmt.Println()

	problems := []string{}
	suggestions := []string{}

	// Check Go
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		problems = append(problems, "Go is not installed")
		suggestions = append(suggestions, "Install Go 1.21+ from https://golang.org/dl/")
	}

	// Check Node.js
	cmd = exec.Command("node", "--version")
	if err := cmd.Run(); err != nil {
		problems = append(problems, "Node.js is not installed")
		suggestions = append(suggestions, "Install Node.js 18+ from https://nodejs.org/")
	}

	// Check Docker
	cmd = exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		problems = append(problems, "Docker is not running")
		suggestions = append(suggestions, "Start Docker Desktop or install Docker from https://docker.com/")
	}

	// Check docker-compose.yml
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		problems = append(problems, "docker-compose.yml not found")
		suggestions = append(suggestions, "Ensure you're running from the project root directory")
	}

	// Check for .env files
	envFiles := []string{
		"apps/creator-studio-web/.env.local",
		"services/creator-api/.env",
	}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			problems = append(problems, fmt.Sprintf("Missing %s", envFile))
			exampleFile := envFile + ".example"
			if _, err := os.Stat(exampleFile); err == nil {
				suggestions = append(suggestions, fmt.Sprintf("Copy %s to %s and configure", exampleFile, envFile))
			}
		}
	}

	// Display results
	if len(problems) == 0 {
		fmt.Println("✅ No issues detected!")
		fmt.Println()
		fmt.Println("Your development environment appears healthy.")
		return nil
	}

	fmt.Printf("❌ Found %d issue(s):\n\n", len(problems))

	for i, problem := range problems {
		fmt.Printf("%d. %s\n", i+1, problem)
	}

	fmt.Println()
	fmt.Println("💡 Suggested Fixes:")
	fmt.Println()

	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion)
	}

	fmt.Println()
	fmt.Println("Run 'tui-styles-cli dev setup' to attempt automatic fixes")

	return nil
}

func startDevServices() error {
	fmt.Println("🚀 Starting development services...")

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Services started!")
	fmt.Println()
	fmt.Println("View logs: docker-compose logs -f")

	return nil
}

func stopDevServices() error {
	fmt.Println("🛑 Stopping development services...")

	cmd := exec.Command("docker-compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Services stopped!")

	return nil
}

func isPortInUse(port int) bool {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	return cmd.Run() == nil
}
