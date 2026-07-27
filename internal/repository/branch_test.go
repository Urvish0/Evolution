package repository

import (
	"testing"
)

func TestBranchOperations(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Check current branch is "main"
	curr, err := GetCurrentBranchName()
	if err != nil {
		t.Fatalf("GetCurrentBranchName() failed: %v", err)
	}
	if curr != "main" {
		t.Errorf("expected initial branch to be 'main', got %q", curr)
	}

	// 2. Create new branch "experiment"
	if err := CreateBranch("experiment"); err != nil {
		t.Fatalf("CreateBranch('experiment') failed: %v", err)
	}

	// 3. List branches (should contain "main" and "experiment")
	branches, err := ListBranches()
	if err != nil {
		t.Fatalf("ListBranches() failed: %v", err)
	}
	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}

	// 4. Checkout "experiment"
	if err := CheckoutBranch("experiment"); err != nil {
		t.Fatalf("CheckoutBranch('experiment') failed: %v", err)
	}

	curr, _ = GetCurrentBranchName()
	if curr != "experiment" {
		t.Errorf("expected current branch to be 'experiment', got %q", curr)
	}

	// 5. Deleting active branch should fail
	if err := DeleteBranch("experiment"); err == nil {
		t.Errorf("expected deleting active branch to fail, but got nil")
	}

	// 6. Checkout "main" and delete "experiment"
	_ = CheckoutBranch("main")
	if err := DeleteBranch("experiment"); err != nil {
		t.Fatalf("DeleteBranch('experiment') failed: %v", err)
	}

	branches, _ = ListBranches()
	if len(branches) != 1 {
		t.Errorf("expected 1 branch after deletion, got %d", len(branches))
	}
}
