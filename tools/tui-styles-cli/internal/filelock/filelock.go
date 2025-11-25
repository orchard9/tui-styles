package filelock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Lock represents a file lock
type Lock struct {
	lockFile string
	acquired bool
}

// NewLock creates a new file lock
func NewLock(filePath string) *Lock {
	lockFile := filePath + ".lock"
	return &Lock{
		lockFile: lockFile,
		acquired: false,
	}
}

// Acquire attempts to acquire the lock with retries
func (l *Lock) Acquire(timeout time.Duration) error {
	start := time.Now()
	for {
		// Try to create lock file exclusively
		f, err := os.OpenFile(l.lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			// Successfully acquired lock
			defer f.Close()
			// Write PID to lock file
			fmt.Fprintf(f, "%d\n%s\n", os.Getpid(), time.Now().Format(time.RFC3339))
			l.acquired = true
			return nil
		}

		// Check if timeout exceeded
		if time.Since(start) >= timeout {
			return fmt.Errorf("timeout acquiring lock for %s", filepath.Base(l.lockFile))
		}

		// Wait before retrying
		time.Sleep(100 * time.Millisecond)
	}
}

// Release releases the lock
func (l *Lock) Release() error {
	if !l.acquired {
		return nil
	}

	err := os.Remove(l.lockFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	l.acquired = false
	return nil
}

// CleanStaleLocks removes lock files older than the specified duration
func CleanStaleLocks(dir string, maxAge time.Duration) error {
	pattern := filepath.Join(dir, "*.lock")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, lockFile := range matches {
		info, err := os.Stat(lockFile)
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			os.Remove(lockFile) // Ignore errors for stale lock removal
		}
	}

	return nil
}
