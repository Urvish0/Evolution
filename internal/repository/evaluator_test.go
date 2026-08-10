package repository

import (
	"testing"
)

func TestEvaluationFramework(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	_ = CreateCommit("Commit for evaluation framework test")
	exec, err := RecordExecution("What is corporate law?", "Corporate law governs company formation and contracts.", 300, TokenUsage{PromptTokens: 40, CompletionTokens: 20}, nil)
	if err != nil {
		t.Fatalf("RecordExecution() failed: %v", err)
	}

	state, _ := ReconstructCommitState(exec.CommitID)

	res, err := EvaluateExecution(exec, state, nil)
	if err != nil {
		t.Fatalf("EvaluateExecution() failed: %v", err)
	}

	if res.ExecutionID != exec.ID {
		t.Errorf("expected execution ID %s, got %s", exec.ID, res.ExecutionID)
	}
	if len(res.Scores) != 4 {
		t.Errorf("expected 4 evaluator scores, got %d", len(res.Scores))
	}
	if res.OverallScore <= 0.0 {
		t.Errorf("expected positive overall score, got %.2f", res.OverallScore)
	}

	// Test LoadEvaluationResult by short ID
	loaded, err := LoadEvaluationResult(res.ID[:8])
	if err != nil {
		t.Fatalf("LoadEvaluationResult(%s) failed: %v", res.ID[:8], err)
	}
	if loaded.OverallScore != res.OverallScore {
		t.Errorf("expected overall score %.2f, got %.2f", res.OverallScore, loaded.OverallScore)
	}

	// Test ListEvaluationResults
	list, err := ListEvaluationResults()
	if err != nil {
		t.Fatalf("ListEvaluationResults() failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 evaluation result in list, got %d", len(list))
	}
}

func TestCommitEvaluationAndComparison(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Commit 1 & Execution
	_ = CreateCommit("Commit 1 for evaluation comparison")
	r1, _ := OpenRepository()
	c1ID := r1.Branch.Head
	_, _ = RecordExecution("Q1", "Fast response", 250, TokenUsage{PromptTokens: 20, CompletionTokens: 10, TotalTokens: 30}, nil)

	// 2. Commit 2 & Execution
	_ = CreateCommit("Commit 2 for evaluation comparison")
	r2, _ := OpenRepository()
	c2ID := r2.Branch.Head
	_, _ = RecordExecution("Q1", "Slower response with more details", 850, TokenUsage{PromptTokens: 50, CompletionTokens: 40, TotalTokens: 90}, nil)

	report1, err := EvaluateCommit(c1ID[:8])
	if err != nil {
		t.Fatalf("EvaluateCommit(c1) failed: %v", err)
	}
	if report1.ExecutionCount != 1 {
		t.Errorf("expected 1 execution for commit 1, got %d", report1.ExecutionCount)
	}

	report2, err := EvaluateCommit(c2ID[:8])
	if err != nil {
		t.Fatalf("EvaluateCommit(c2) failed: %v", err)
	}
	if report2.ExecutionCount != 1 {
		t.Errorf("expected 1 execution for commit 2, got %d", report2.ExecutionCount)
	}

	comp, err := CompareCommitEvaluations(c1ID[:8], c2ID[:8])
	if err != nil {
		t.Fatalf("CompareCommitEvaluations() failed: %v", err)
	}

	if comp.Report1.CommitID != c1ID {
		t.Errorf("expected Report1 commit %s, got %s", c1ID, comp.Report1.CommitID)
	}
	if comp.Report2.CommitID != c2ID {
		t.Errorf("expected Report2 commit %s, got %s", c2ID, comp.Report2.CommitID)
	}
}
