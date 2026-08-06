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

	if commit.Metadata.Environment["os"] == "" {
		t.Errorf("expected metadata environment os to be populated")
	}
}

func TestIntelligenceCommitArtifactsSerialization(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	commit := NewCommit("Intelligence commit with artifacts")
	commit.Artifacts["prompts"] = []Artifact{
		PromptArtifact{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypePrompt, Name: "sys-prompt", Hash: "1111111111111111111111111111111111111111111111111111111111111111", Path: "prompts/sys.txt"},
			Role:         "system",
			Format:       "text",
		},
	}
	commit.Artifacts["model_config"] = []Artifact{
		ModelConfigArtifact{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypeModelConfig, Name: "gpt4-cfg", Hash: "2222222222222222222222222222222222222222222222222222222222222222", Path: "config/model.json"},
			Model:        "gpt-4o",
			Provider:     "openai",
			Temperature:  0.7,
		},
	}

	if err := WriteCommit(commit); err != nil {
		t.Fatalf("WriteCommit() failed: %v", err)
	}

	loadedCommit, err := LoadCommit(commit.ID)
	if err != nil {
		t.Fatalf("LoadCommit() failed: %v", err)
	}

	if len(loadedCommit.Artifacts["prompts"]) != 1 {
		t.Fatalf("expected 1 prompt artifact, got %d", len(loadedCommit.Artifacts["prompts"]))
	}

	promptArt, ok := loadedCommit.Artifacts["prompts"][0].(PromptArtifact)
	if !ok {
		t.Fatalf("type assertion to PromptArtifact failed")
	}

	if promptArt.Role != "system" {
		t.Errorf("expected role 'system', got %q", promptArt.Role)
	}

	if len(loadedCommit.Artifacts["model_config"]) != 1 {
		t.Fatalf("expected 1 model_config artifact, got %d", len(loadedCommit.Artifacts["model_config"]))
	}

	modelArt, ok := loadedCommit.Artifacts["model_config"][0].(ModelConfigArtifact)
	if !ok {
		t.Fatalf("type assertion to ModelConfigArtifact failed")
	}

	if modelArt.Model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", modelArt.Model)
	}
}

func TestCommitWithOptsTagsAndMeta(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	tags := []string{"production", "v1.0"}
	meta := map[string]string{
		"model":       "gpt-4o",
		"temperature": "0.2",
	}

	if err := CreateCommitWithOpts("Tagged commit", tags, meta); err != nil {
		t.Fatalf("CreateCommitWithOpts() failed: %v", err)
	}

	repo, _ := OpenRepository()
	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		t.Fatalf("LoadCommit() failed: %v", err)
	}

	if len(commit.Metadata.Tags) != 2 || commit.Metadata.Tags[0] != "production" {
		t.Errorf("expected tags [production v1.0], got %v", commit.Metadata.Tags)
	}

	if commit.Metadata.Environment["model"] != "gpt-4o" {
		t.Errorf("expected custom meta model 'gpt-4o', got %q", commit.Metadata.Environment["model"])
	}
}
