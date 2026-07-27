package repository

import (
	"testing"
)

func TestLogMultiCommitDAG(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create 3 commits
	messages := []string{"Commit 1", "Commit 2", "Commit 3"}
	for _, msg := range messages {
		if err := CreateCommit(msg); err != nil {
			t.Fatalf("CreateCommit(%q) failed: %v", msg, err)
		}
	}

	// Walk log
	commits, err := Log()
	if err != nil {
		t.Fatalf("Log() failed: %v", err)
	}

	if len(commits) != 3 {
		t.Fatalf("expected 3 commits in log, got %d", len(commits))
	}

	// Verify reverse chronological order (newest to oldest)
	if commits[0].Message != "Commit 3" {
		t.Errorf("expected latest commit to be 'Commit 3', got %q", commits[0].Message)
	}
	if commits[1].Message != "Commit 2" {
		t.Errorf("expected middle commit to be 'Commit 2', got %q", commits[1].Message)
	}
	if commits[2].Message != "Commit 1" {
		t.Errorf("expected initial commit to be 'Commit 1', got %q", commits[2].Message)
	}
}
