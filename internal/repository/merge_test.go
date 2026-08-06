package repository

import (
	"os"
	"strings"
	"testing"
)

func TestFindLowestCommonAncestor(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Root commit
	_ = CreateCommit("Root commit C0")
	repo0, _ := OpenRepository()
	c0 := repo0.Branch.Head

	// Branch 1 (Main)
	_ = CreateCommit("Main commit C1")
	repo1, _ := OpenRepository()
	c1 := repo1.Branch.Head

	// Branch 2 (Feature) starting from C0
	_ = CreateBranch("feature")
	fBranch, _ := LoadBranch("feature")
	fBranch.Head = c0
	_ = UpdateBranch(fBranch)
	_ = CheckoutBranch("feature")
	_ = CreateCommit("Feature commit C2")
	repo2, _ := OpenRepository()
	c2 := repo2.Branch.Head

	// Find LCA between C1 and C2 -> should be C0!
	lca, err := FindLowestCommonAncestor(c1, c2)
	if err != nil {
		t.Fatalf("FindLowestCommonAncestor() failed: %v", err)
	}

	if lca != c0 {
		t.Errorf("expected LCA to be %s, got %s", c0, lca)
	}
}

func TestFastForwardMerge(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	_ = CreateCommit("Initial commit")
	_ = CreateBranch("feature")
	_ = CheckoutBranch("feature")
	_ = os.WriteFile("feature.txt", []byte("Feature content"), 0644)
	_ = StagePath("feature.txt")
	_ = CreateCommit("Feature commit")

	// Switch back to main and merge feature -> should be Fast-Forward!
	_ = CheckoutBranch("main")
	result, err := MergeBranch("feature")
	if err != nil {
		t.Fatalf("MergeBranch() failed: %v", err)
	}

	if !result.IsFastForward {
		t.Errorf("expected fast-forward merge")
	}

	if _, err := os.Stat("feature.txt"); os.IsNotExist(err) {
		t.Errorf("expected feature.txt to be merged into working tree")
	}
}

func TestThreeWayMergeNoConflict(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Base state
	_ = os.WriteFile("file1.txt", []byte("Base 1"), 0644)
	_ = os.WriteFile("file2.txt", []byte("Base 2"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Base commit")

	// Feature branch edits file2.txt
	_ = CreateBranch("feature")
	_ = CheckoutBranch("feature")
	_ = os.WriteFile("file2.txt", []byte("Feature 2 Edit"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Feature edit file2")

	// Main branch edits file1.txt
	_ = CheckoutBranch("main")
	_ = os.WriteFile("file1.txt", []byte("Main 1 Edit"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Main edit file1")

	// Merge feature into main
	result, err := MergeBranch("feature")
	if err != nil {
		t.Fatalf("MergeBranch() failed: %v", err)
	}

	if result.HasConflicts {
		t.Fatalf("expected clean merge without conflicts")
	}

	c1, _ := os.ReadFile("file1.txt")
	c2, _ := os.ReadFile("file2.txt")

	if string(c1) != "Main 1 Edit" {
		t.Errorf("expected 'Main 1 Edit', got %q", string(c1))
	}
	if string(c2) != "Feature 2 Edit" {
		t.Errorf("expected 'Feature 2 Edit', got %q", string(c2))
	}
}

func TestThreeWayMergeWithConflict(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Base state
	_ = os.WriteFile("prompt.txt", []byte("System Prompt: Base version"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Base prompt")

	// Feature branch edits prompt.txt
	_ = CreateBranch("feature")
	_ = CheckoutBranch("feature")
	_ = os.WriteFile("prompt.txt", []byte("System Prompt: Feature version"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Feature prompt edit")

	// Main branch edits prompt.txt differently
	_ = CheckoutBranch("main")
	_ = os.WriteFile("prompt.txt", []byte("System Prompt: Main version"), 0644)
	_ = StagePath(".")
	_ = CreateCommit("Main prompt edit")

	// Merge feature into main -> should detect CONFLICT!
	result, err := MergeBranch("feature")
	if err != nil {
		t.Fatalf("MergeBranch() failed: %v", err)
	}

	if !result.HasConflicts {
		t.Fatalf("expected merge conflict")
	}

	content, _ := os.ReadFile("prompt.txt")
	strContent := string(content)

	if !strings.Contains(strContent, "<<<<<<< OURS") || !strings.Contains(strContent, ">>>>>>> THEIRS") {
		t.Errorf("expected conflict markers in prompt.txt, got:\n%s", strContent)
	}
}
