package repository

import (
	"os"
	"testing"
)

func TestManifestLifecycle(t *testing.T) {
	setupTestRepo(t)

	// 1. Init Manifest
	if err := InitManifest("test-ai", "test description"); err != nil {
		t.Fatalf("InitManifest() failed: %v", err)
	}

	if !ManifestExists() {
		t.Fatalf("expected ManifestExists() to be true")
	}

	// 2. Load Manifest
	m, err := LoadManifest(ManifestFileName)
	if err != nil {
		t.Fatalf("LoadManifest() failed: %v", err)
	}

	if m.Name != "test-ai" {
		t.Errorf("expected name 'test-ai', got %q", m.Name)
	}
	if m.Version != SpecVersionV10 {
		t.Errorf("expected version %q, got %q", SpecVersionV10, m.Version)
	}

	// 3. Validate Manifest
	if err := ValidateManifest(m); err != nil {
		t.Fatalf("ValidateManifest() failed: %v", err)
	}

	// 4. Test invalid manifest
	invalidM := Manifest{Version: "", Name: ""}
	if err := ValidateManifest(invalidM); err == nil {
		t.Errorf("expected validation error for empty manifest, got nil")
	}
}

func TestManifestAutoCommitIntegration(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create manifest with a prompt artifact
	m := NewDefaultManifest("auto-commit-ai", "testing manifest auto-commit")
	m.Artifacts.Prompts = []PromptArtifact{
		{
			BaseArtifact: BaseArtifact{
				ArtifactType: ArtifactTypePrompt,
				Name:         "system-prompt",
				Hash:         "b671051cadc3b0806ab998726cb95055a713983b1e94605e9bc3d88689d21a0e",
				Path:         "prompts/system.txt",
			},
			Role:   "system",
			Format: "text",
		},
	}

	if err := SaveManifest(m, ManifestFileName); err != nil {
		t.Fatalf("SaveManifest() failed: %v", err)
	}

	// Create commit — should auto-detect manifest
	if err := CreateCommit("Commit with manifest"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	repo, _ := OpenRepository()
	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		t.Fatalf("LoadCommit() failed: %v", err)
	}

	if len(commit.Artifacts["prompts"]) != 1 {
		t.Fatalf("expected 1 prompt artifact in commit, got %d", len(commit.Artifacts["prompts"]))
	}

	if len(commit.Artifacts["model_config"]) != 1 {
		t.Fatalf("expected 1 model_config artifact in commit, got %d", len(commit.Artifacts["model_config"]))
	}

	// Cleanup manifest file from working directory
	_ = os.Remove(ManifestFileName)
}
