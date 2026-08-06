package repository

import (
	"os"
	"testing"
)

func TestArtifactRegistrationAndListing(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Register a prompt artifact
	if err := RegisterArtifactInManifest(ArtifactTypePrompt, "sys-prompt", "prompts/sys.txt"); err != nil {
		t.Fatalf("RegisterArtifactInManifest() failed: %v", err)
	}

	if !ManifestExists() {
		t.Fatalf("expected evolution.manifest.json to exist")
	}

	// 2. Create a commit (should auto-detect manifest)
	if err := CreateCommit("Added sys-prompt artifact"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// 3. Inspect HEAD artifacts
	artifacts, err := GetHeadCommitArtifacts()
	if err != nil {
		t.Fatalf("GetHeadCommitArtifacts() failed: %v", err)
	}

	if len(artifacts["prompts"]) != 1 {
		t.Fatalf("expected 1 prompt artifact in HEAD commit, got %d", len(artifacts["prompts"]))
	}

	if artifacts["prompts"][0].GetName() != "sys-prompt" {
		t.Errorf("expected artifact name 'sys-prompt', got %q", artifacts["prompts"][0].GetName())
	}

	// Cleanup manifest
	_ = os.Remove(ManifestFileName)
}

func TestSemanticArtifactDiff(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Commit 1: Register prompt artifact v1
	_ = RegisterArtifactInManifest(ArtifactTypePrompt, "sys-prompt-v1", "prompts/v1.txt")
	_ = CreateCommit("Commit 1: v1 prompt")
	repo1, _ := OpenRepository()
	commit1ID := repo1.Branch.Head

	// Commit 2: Register second prompt artifact v2
	_ = RegisterArtifactInManifest(ArtifactTypePrompt, "sys-prompt-v2", "prompts/v2.txt")
	_ = CreateCommit("Commit 2: v2 prompt added")
	repo2, _ := OpenRepository()
	commit2ID := repo2.Branch.Head

	// Compare commit 1 vs commit 2
	changes, err := CompareCommitArtifacts(commit1ID, commit2ID)
	if err != nil {
		t.Fatalf("CompareCommitArtifacts() failed: %v", err)
	}

	if len(changes) == 0 {
		t.Fatalf("expected semantic diff changes between commit 1 and commit 2")
	}

	foundPromptAdded := false
	for _, c := range changes {
		if c.ArtifactName == "sys-prompt-v2" && c.Action == "added" {
			foundPromptAdded = true
		}
	}

	if !foundPromptAdded {
		t.Errorf("expected 'sys-prompt-v2' added change in semantic diff")
	}

	_ = os.Remove(ManifestFileName)
}
