// Package config provides configuration management for the service.
package config

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// loadFromGSM fetches secret from Google Secret Manager if CONFIG_FROM_GSM_SECRET is set
// and populates environment variables from the .env formatted content.
// Environment variables already set take precedence over GSM values.
func loadFromGSM(ctx context.Context) error {
	secretName := os.Getenv("CONFIG_FROM_GSM_SECRET")
	if secretName == "" {
		// No GSM secret configured, skip
		return nil
	}

	// Create Secret Manager client
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("create secret manager client: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Build the resource name for the secret version
	name := fmt.Sprintf("%s/versions/latest", secretName)

	// Access the secret version
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("access secret version %s: %w", name, err)
	}

	// Parse .env format content
	content := string(result.Payload.Data)
	if err := parseAndSetEnvVars(content); err != nil {
		return fmt.Errorf("parse env vars from GSM secret: %w", err)
	}

	return nil
}

// parseAndSetEnvVars parses .env format content and sets environment variables.
// Existing environment variables are not overwritten.
func parseAndSetEnvVars(content string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env var format at line %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already in environment (env vars take precedence)
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set env var %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan secret content: %w", err)
	}

	return nil
}
