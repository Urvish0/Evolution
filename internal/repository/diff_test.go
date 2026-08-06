package repository

import (
	"os"
	"testing"
)

func TestCompareWorkingTreeUntracked(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create a file but don't commit — it should be untracked
	if err := os.WriteFile("untracked.txt", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	wts, err := CompareWorkingTree()
	if err != nil {
		t.Fatalf("CompareWorkingTree() failed: %v", err)
	}

	if len(wts.Untracked) == 0 {
		t.Errorf("expected untracked files, got none")
	}

	found := false
	for _, f := range wts.Untracked {
		if f.Path == "untracked.txt" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'untracked.txt' in untracked files")
	}
}

func TestCompareWorkingTreeModified(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create, stage, and commit a file
	if err := os.WriteFile("prompt.txt", []byte("version 1"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := StagePath("prompt.txt"); err != nil {
		t.Fatalf("StagePath() failed: %v", err)
	}

	if err := CreateCommit("initial"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// Modify the file without staging
	if err := os.WriteFile("prompt.txt", []byte("version 2"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	wts, err := CompareWorkingTree()
	if err != nil {
		t.Fatalf("CompareWorkingTree() failed: %v", err)
	}

	if len(wts.Modified) == 0 {
		t.Errorf("expected modified files, got none")
	}
}

func TestCompareWorkingTreeDeleted(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create, stage, and commit a file
	if err := os.WriteFile("to_delete.txt", []byte("will be deleted"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := StagePath("to_delete.txt"); err != nil {
		t.Fatalf("StagePath() failed: %v", err)
	}

	if err := CreateCommit("add file"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// Delete the file from workspace
	if err := os.Remove("to_delete.txt"); err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	wts, err := CompareWorkingTree()
	if err != nil {
		t.Fatalf("CompareWorkingTree() failed: %v", err)
	}

	if len(wts.Deleted) == 0 {
		t.Errorf("expected deleted files, got none")
	}
}

func TestCompareWorkingTreeClean(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create, stage, and commit a file
	if err := os.WriteFile("stable.txt", []byte("no changes"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := StagePath("stable.txt"); err != nil {
		t.Fatalf("StagePath() failed: %v", err)
	}

	if err := CreateCommit("stable commit"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// File is unchanged — tree should be clean
	wts, err := CompareWorkingTree()
	if err != nil {
		t.Fatalf("CompareWorkingTree() failed: %v", err)
	}

	if !wts.IsClean() {
		t.Errorf("expected clean working tree, got staged=%d modified=%d untracked=%d deleted=%d",
			len(wts.Staged), len(wts.Modified), len(wts.Untracked), len(wts.Deleted))
	}
}

func TestCrossRevisionDiff(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	_ = os.WriteFile("p.txt", []byte("Prompt V1"), 0644)
	_ = StagePath("p.txt")
	_ = CreateCommit("Commit V1")
	repo1, _ := OpenRepository()
	c1ID := repo1.Branch.Head

	_ = CreateBranch("exp")
	_ = CheckoutBranch("exp")
	_ = os.WriteFile("p.txt", []byte("Prompt V2 Exp"), 0644)
	_ = StagePath("p.txt")
	_ = CreateCommit("Commit V2 Exp")
	repo2, _ := OpenRepository()
	c2ID := repo2.Branch.Head

	// Test ResolveRevisionToCommit
	c1, err := ResolveRevisionToCommit("main")
	if err != nil || c1.ID != c1ID {
		t.Errorf("failed to resolve main branch to commit")
	}

	c2, err := ResolveRevisionToCommit(c2ID[:8])
	if err != nil || c2.ID != c2ID {
		t.Errorf("failed to resolve short commit hash %s", c2ID[:8])
	}

	// Test GetRevisionDiff between main and exp branches
	diffStr, err := GetRevisionDiff("main", "exp")
	if err != nil {
		t.Fatalf("GetRevisionDiff() failed: %v", err)
	}

	if !os.IsNotExist(err) && (len(diffStr) == 0) {
		t.Errorf("expected cross-branch diff output")
	}
}
