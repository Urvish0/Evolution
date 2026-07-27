package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTreeWriteAndRead(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	tree := Tree{
		Entries: []TreeEntry{
			{Mode: ModeFile, Type: ObjectTypeBlob, Hash: "b671051cadc3b0806ab998726cb95055a713983b1e94605e9bc3d88689d21a0e", Name: "prompt.txt"},
			{Mode: ModeFile, Type: ObjectTypeBlob, Hash: "a42f891b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f", Name: "model.json"},
		},
	}

	// 1. Write tree
	treeHash, err := WriteTree(tree)
	if err != nil {
		t.Fatalf("WriteTree() failed: %v", err)
	}

	if len(treeHash) != 64 {
		t.Fatalf("expected 64-character tree hash, got %d", len(treeHash))
	}

	// 2. Read tree back
	loadedTree, err := ReadTree(treeHash)
	if err != nil {
		t.Fatalf("ReadTree() failed: %v", err)
	}

	if len(loadedTree.Entries) != 2 {
		t.Fatalf("expected 2 tree entries, got %d", len(loadedTree.Entries))
	}

	// Sorted alphabetically by name: model.json first, prompt.txt second
	if loadedTree.Entries[0].Name != "model.json" {
		t.Errorf("expected first entry 'model.json', got %q", loadedTree.Entries[0].Name)
	}
	if loadedTree.Entries[1].Name != "prompt.txt" {
		t.Errorf("expected second entry 'prompt.txt', got %q", loadedTree.Entries[1].Name)
	}
}

func TestBuildTreeFromDirectoryRecursive(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create a nested test directory structure:
	// workspace/
	// ├── main_prompt.txt
	// └── tools/
	//     └── search_tool.json
	workDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(workDir, "main_prompt.txt"), []byte("System prompt v1"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	toolsDir := filepath.Join(workDir, "tools")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		t.Fatalf("failed to create tools dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(toolsDir, "search_tool.json"), []byte(`{"name":"search"}`), 0644); err != nil {
		t.Fatalf("failed to write tool file: %v", err)
	}

	// Build Merkle Tree from workDir
	rootTreeHash, err := BuildTreeFromDirectory(workDir)
	if err != nil {
		t.Fatalf("BuildTreeFromDirectory() failed: %v", err)
	}

	if len(rootTreeHash) != 64 {
		t.Fatalf("expected 64-char root tree hash, got %d", len(rootTreeHash))
	}

	// Load root tree back and inspect entries
	rootTree, err := ReadTree(rootTreeHash)
	if err != nil {
		t.Fatalf("ReadTree(rootTreeHash) failed: %v", err)
	}

	if len(rootTree.Entries) != 2 {
		t.Fatalf("expected 2 entries in root tree (main_prompt.txt and tools), got %d", len(rootTree.Entries))
	}

	// Check tools subdirectory entry is type "tree"
	var toolsEntry *TreeEntry
	for _, e := range rootTree.Entries {
		if e.Name == "tools" {
			toolsEntry = &e
			break
		}
	}

	if toolsEntry == nil {
		t.Fatalf("expected 'tools' entry in root tree")
	}

	if toolsEntry.Type != ObjectTypeTree {
		t.Errorf("expected 'tools' entry type 'tree', got %q", toolsEntry.Type)
	}

	// Load child tree for "tools"
	childTree, err := ReadTree(toolsEntry.Hash)
	if err != nil {
		t.Fatalf("ReadTree(toolsEntry.Hash) failed: %v", err)
	}

	if len(childTree.Entries) != 1 || childTree.Entries[0].Name != "search_tool.json" {
		t.Errorf("expected child tree to contain 'search_tool.json'")
	}
}
