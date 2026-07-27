package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexStageAndBuildTree(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a test file
	fileContent := []byte("staged prompt content")
	if err := os.WriteFile("staged_prompt.txt", fileContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// 2. Stage file
	if err := StagePath("staged_prompt.txt"); err != nil {
		t.Fatalf("StagePath() failed: %v", err)
	}

	// 3. Load index and verify entry
	idx, err := LoadIndex()
	if err != nil {
		t.Fatalf("LoadIndex() failed: %v", err)
	}

	if len(idx.Entries) != 1 {
		t.Fatalf("expected 1 staged index entry, got %d", len(idx.Entries))
	}

	entry, exists := idx.Entries["staged_prompt.txt"]
	if !exists {
		t.Fatalf("expected 'staged_prompt.txt' entry in index")
	}

	if len(entry.Hash) != 64 {
		t.Errorf("expected 64-char blob hash, got %d", len(entry.Hash))
	}

	// 4. Build Merkle Tree from index
	treeHash, err := BuildTreeFromIndex(idx)
	if err != nil {
		t.Fatalf("BuildTreeFromIndex() failed: %v", err)
	}

	if len(treeHash) != 64 {
		t.Errorf("expected 64-char tree hash, got %d", len(treeHash))
	}

	// 5. Verify commit uses staged index and clears it
	if err := CreateCommit("Commit with index"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	idxAfter, _ := LoadIndex()
	if len(idxAfter.Entries) != 0 {
		t.Errorf("expected index to be cleared after commit, got %d entries", len(idxAfter.Entries))
	}
}

func TestStageDirectoryRecursive(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create subfolder with file
	subDir := "prompts"
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subfolder: %v", err)
	}

	if err := os.WriteFile(filepath.Join(subDir, "system.txt"), []byte("system prompt"), 0644); err != nil {
		t.Fatalf("failed to write subfolder file: %v", err)
	}

	// Stage directory
	if err := StagePath(subDir); err != nil {
		t.Fatalf("StagePath(%q) failed: %v", subDir, err)
	}

	idx, err := LoadIndex()
	if err != nil {
		t.Fatalf("LoadIndex() failed: %v", err)
	}

	if len(idx.Entries) != 1 {
		t.Fatalf("expected 1 staged entry for subfolder file, got %d", len(idx.Entries))
	}
}
