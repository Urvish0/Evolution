package repository

import (
	"testing"
)

func TestCommitOperations(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a commit
	if err := CreateCommit("First commit"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// 2. Open repo and check HEAD is non-empty
	repo, err := OpenRepository()
	if err != nil {
		t.Fatalf("OpenRepository() failed: %v", err)
	}

	if repo.Branch.Head == "" {
		t.Fatalf("expected repo.Branch.Head to be non-empty UUID")
	}

	// 3. Load commit by ID
	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		t.Fatalf("LoadCommit() failed: %v", err)
	}

	if commit.Message != "First commit" {
		t.Errorf("expected message 'First commit', got %q", commit.Message)
	}
}
