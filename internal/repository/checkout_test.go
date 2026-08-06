package repository

import (
	"os"
	"testing"
)

func TestCheckoutTreeRestoration(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a file on main branch
	_ = os.MkdirAll("prompts", 0755)
	_ = os.WriteFile("prompts/v1.txt", []byte("Prompt V1 on main"), 0644)
	if err := CreateCommit("Main commit V1"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// 2. Create feature branch and switch to it
	if err := CreateBranch("feature"); err != nil {
		t.Fatalf("CreateBranch() failed: %v", err)
	}
	if err := CheckoutBranch("feature"); err != nil {
		t.Fatalf("CheckoutBranch() failed: %v", err)
	}

	// 3. Modify file on feature branch and commit
	_ = os.WriteFile("prompts/v1.txt", []byte("Prompt V2 on feature branch"), 0644)
	if err := CreateCommit("Feature commit V2"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// Read content on feature branch
	content, _ := os.ReadFile("prompts/v1.txt")
	if string(content) != "Prompt V2 on feature branch" {
		t.Errorf("expected 'Prompt V2 on feature branch', got %q", string(content))
	}

	// 4. Switch back to main -> workspace file should be restored to V1!
	if err := CheckoutBranch("main"); err != nil {
		t.Fatalf("CheckoutBranch('main') failed: %v", err)
	}

	restoredContent, _ := os.ReadFile("prompts/v1.txt")
	if string(restoredContent) != "Prompt V1 on main" {
		t.Errorf("expected restored content 'Prompt V1 on main', got %q", string(restoredContent))
	}
}

func TestRestoreFileFromHEAD(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create file and commit
	_ = os.MkdirAll("config", 0755)
	_ = os.WriteFile("config/model.json", []byte(`{"model":"gpt-4o"}`), 0644)
	_ = CreateCommit("Initial model config")

	// 2. Mutate file locally on disk without committing
	_ = os.WriteFile("config/model.json", []byte(`{"model":"corrupted"}`), 0644)

	// 3. Restore file from HEAD
	if err := RestoreFileFromHEAD("config/model.json"); err != nil {
		t.Fatalf("RestoreFileFromHEAD() failed: %v", err)
	}

	// 4. Verify file content restored
	restored, _ := os.ReadFile("config/model.json")
	if string(restored) != `{"model":"gpt-4o"}` {
		t.Errorf("expected '{\"model\":\"gpt-4o\"}', got %q", string(restored))
	}
}

func TestCheckoutCommitDetachedHEAD(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	_ = os.WriteFile("test.txt", []byte("State 1"), 0644)
	_ = CreateCommit("Commit 1")
	repo1, _ := OpenRepository()
	c1ID := repo1.Branch.Head

	_ = os.WriteFile("test.txt", []byte("State 2"), 0644)
	_ = CreateCommit("Commit 2")

	// Checkout Commit 1 by short ID
	shortID := c1ID[:8]
	if err := CheckoutCommit(shortID); err != nil {
		t.Fatalf("CheckoutCommit(%s) failed: %v", shortID, err)
	}

	content, _ := os.ReadFile("test.txt")
	if string(content) != "State 1" {
		t.Errorf("expected 'State 1', got %q", string(content))
	}
}
