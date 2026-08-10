package repository

import (
	"testing"
)

func TestExecutionSaveAndLoad(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a commit to link against
	_ = CreateCommit("Commit for execution test")
	repo, _ := OpenRepository()
	headID := repo.Branch.Head

	// 2. Record execution
	tokens := TokenUsage{PromptTokens: 100, CompletionTokens: 50}
	meta := map[string]string{"session_id": "sess_123"}

	exec, err := RecordExecution("What is corporate law?", "Corporate law governs company formation and contracts.", 450, tokens, meta)
	if err != nil {
		t.Fatalf("RecordExecution() failed: %v", err)
	}

	if exec.CommitID != headID {
		t.Errorf("expected exec.CommitID %s, got %s", headID, exec.CommitID)
	}
	if exec.Tokens.TotalTokens != 150 {
		t.Errorf("expected total tokens 150, got %d", exec.Tokens.TotalTokens)
	}

	// 3. Load execution back by short ID
	shortID := exec.ID[:8]
	loaded, err := LoadExecution(shortID)
	if err != nil {
		t.Fatalf("LoadExecution(%s) failed: %v", shortID, err)
	}

	if loaded.Inputs != exec.Inputs {
		t.Errorf("expected input %q, got %q", exec.Inputs, loaded.Inputs)
	}

	// 4. List executions
	list, err := ListExecutions()
	if err != nil {
		t.Fatalf("ListExecutions() failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 execution in list, got %d", len(list))
	}
}
