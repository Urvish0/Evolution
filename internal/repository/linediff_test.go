package repository

import (
	"os"
	"testing"
)

func TestComputeLineDiffNoChange(t *testing.T) {
	old := "line 1\nline 2\nline 3"
	new := "line 1\nline 2\nline 3"

	diff := ComputeLineDiff(old, new)

	for _, d := range diff {
		if d.Op != DiffEqual {
			t.Errorf("expected all lines DiffEqual, got op=%d for %q", d.Op, d.Content)
		}
	}
}

func TestComputeLineDiffInsert(t *testing.T) {
	old := "line 1\nline 3"
	new := "line 1\nline 2\nline 3"

	diff := ComputeLineDiff(old, new)

	insertCount := 0
	for _, d := range diff {
		if d.Op == DiffInsert {
			insertCount++
			if d.Content != "line 2" {
				t.Errorf("expected inserted line 'line 2', got %q", d.Content)
			}
		}
	}

	if insertCount != 1 {
		t.Errorf("expected 1 insertion, got %d", insertCount)
	}
}

func TestComputeLineDiffDelete(t *testing.T) {
	old := "line 1\nline 2\nline 3"
	new := "line 1\nline 3"

	diff := ComputeLineDiff(old, new)

	deleteCount := 0
	for _, d := range diff {
		if d.Op == DiffDelete {
			deleteCount++
			if d.Content != "line 2" {
				t.Errorf("expected deleted line 'line 2', got %q", d.Content)
			}
		}
	}

	if deleteCount != 1 {
		t.Errorf("expected 1 deletion, got %d", deleteCount)
	}
}

func TestComputeLineDiffModifiedLine(t *testing.T) {
	old := "System Prompt: Be concise"
	new := "System Prompt: Be thorough and cite sources"

	diff := ComputeLineDiff(old, new)

	var hasDelete, hasInsert bool
	for _, d := range diff {
		if d.Op == DiffDelete && d.Content == "System Prompt: Be concise" {
			hasDelete = true
		}
		if d.Op == DiffInsert && d.Content == "System Prompt: Be thorough and cite sources" {
			hasInsert = true
		}
	}

	if !hasDelete {
		t.Errorf("expected deletion of old line")
	}
	if !hasInsert {
		t.Errorf("expected insertion of new line")
	}
}

func TestComputeLineDiffEmpty(t *testing.T) {
	diff := ComputeLineDiff("", "new content")

	insertCount := 0
	for _, d := range diff {
		if d.Op == DiffInsert {
			insertCount++
		}
	}

	if insertCount != 1 {
		t.Errorf("expected 1 insertion for empty->content diff, got %d", insertCount)
	}
}

func TestGetContentDiffIntegration(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create, stage, and commit a file
	if err := os.WriteFile("prompt.txt", []byte("System Prompt: Be concise"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := StagePath("prompt.txt"); err != nil {
		t.Fatalf("StagePath() failed: %v", err)
	}

	if err := CreateCommit("initial prompt"); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// Modify the file
	if err := os.WriteFile("prompt.txt", []byte("System Prompt: Be thorough and cite sources"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	// Get content diff
	diffLines, err := GetContentDiff("prompt.txt")
	if err != nil {
		t.Fatalf("GetContentDiff() failed: %v", err)
	}

	var hasDelete, hasInsert bool
	for _, d := range diffLines {
		if d.Op == DiffDelete {
			hasDelete = true
		}
		if d.Op == DiffInsert {
			hasInsert = true
		}
	}

	if !hasDelete || !hasInsert {
		t.Errorf("expected both deletions and insertions in content diff, got delete=%v insert=%v", hasDelete, hasInsert)
	}
}
