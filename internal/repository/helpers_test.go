package repository

import (
	"os"
	"testing"
)

// setupTestRepo creates an isolated temporary directory and changes the working directory to it.
// When the test finishes, Go automatically cleans up the directory and restores the original working directory.
func setupTestRepo(t *testing.T) string {
	t.Helper()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory to tempDir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	return tempDir
}
