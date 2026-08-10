package repository

import (
	"os"
	"testing"
)

func TestReplayAndStateReconstruction(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a prompt artifact and manifest
	_ = os.MkdirAll("prompts", 0755)
	_ = os.WriteFile("prompts/sys.txt", []byte("Replay prompt content"), 0644)

	m := NewDefaultManifest("replay-ai", "Testing state reconstruction")
	m.Artifacts.Prompts = []PromptArtifact{
		{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypePrompt, Name: "p1", Path: "prompts/sys.txt"},
			Role:         "system",
			Format:       "text",
		},
	}
	_ = SaveManifest(m, ManifestFileName)

	// 2. Create commit
	if err := CreateCommit("Replay test commit"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}
	repo, _ := OpenRepository()
	commitID := repo.Branch.Head

	// 3. Reconstruct state from commit
	state, err := ReconstructCommitState(commitID[:8])
	if err != nil {
		t.Fatalf("ReconstructCommitState() failed: %v", err)
	}

	if state.CommitID != commitID {
		t.Errorf("expected commitID %s, got %s", commitID, state.CommitID)
	}
	if len(state.Manifest.Artifacts.Prompts) != 1 {
		t.Errorf("expected 1 prompt artifact in reconstructed state, got %d", len(state.Manifest.Artifacts.Prompts))
	}

	// 4. Test ExportReconstructedState
	exportPath := "exported_manifest.json"
	if err := ExportReconstructedState(state, exportPath); err != nil {
		t.Fatalf("ExportReconstructedState() failed: %v", err)
	}
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Errorf("expected exported manifest file to exist")
	}

	// 5. Test ReplayExecution
	tokens := TokenUsage{PromptTokens: 50, CompletionTokens: 20}
	exec, _ := RecordExecution("input", "output", 150, tokens, nil)

	reconstructedState, loadedExec, err := ReplayExecution(exec.ID[:8])
	if err != nil {
		t.Fatalf("ReplayExecution() failed: %v", err)
	}

	if loadedExec.ID != exec.ID {
		t.Errorf("expected execution ID %s, got %s", exec.ID, loadedExec.ID)
	}
	if reconstructedState.CommitID != commitID {
		t.Errorf("expected reconstructed commit ID %s, got %s", commitID, reconstructedState.CommitID)
	}
}
