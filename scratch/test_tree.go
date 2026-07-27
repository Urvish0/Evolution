package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Urvish0/evolution/internal/repository"
)

func main() {
	// 1. Setup isolated test environment in a temp folder
	tempDir, err := os.MkdirTemp("", "evo-tree-test-*")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// 2. Initialize Evolution repo inside temp folder
	if err := repository.Init(); err != nil {
		fmt.Printf("Error initializing repo: %v\n", err)
		return
	}

	// 3. Create a nested workspace structure:
	// workspace/
	// ├── system_prompt.txt
	// └── tools/
	//     └── search_tool.json
	_ = os.WriteFile("system_prompt.txt", []byte("System Prompt: You are a helpful legal assistant AI."), 0644)
	_ = os.MkdirAll("tools", 0755)
	_ = os.WriteFile(filepath.Join("tools", "search_tool.json"), []byte(`{"name":"search_web","description":"Search Google"}`), 0644)

	// 4. Build Merkle Tree for the entire workspace
	rootTreeHash, err := repository.BuildTreeFromDirectory(".")
	if err != nil {
		fmt.Printf("Error building tree: %v\n", err)
		return
	}

	fmt.Printf(" Root Merkle Tree Hash: %s\n", rootTreeHash)
	fmt.Printf("   Object Path:           .evolution/objects/%s/%s\n\n", rootTreeHash[:2], rootTreeHash[2:])

	// 5. Inspect the Root Tree entries
	rootTree, err := repository.ReadTree(rootTreeHash)
	if err != nil {
		fmt.Printf("Error reading root tree: %v\n", err)
		return
	}

	fmt.Println(" Root Tree Entries:")
	for _, entry := range rootTree.Entries {
		fmt.Printf("   [%s] %-5s %s  %s\n", entry.Mode, entry.Type, entry.Hash[:8], entry.Name)
	}

	// 6. Find and inspect the 'tools' child tree
	for _, entry := range rootTree.Entries {
		if entry.Name == "tools" {
			childTree, _ := repository.ReadTree(entry.Hash)
			fmt.Printf("\n Subdirectory 'tools/' Tree Entries (Hash: %s):\n", entry.Hash[:8])
			for _, child := range childTree.Entries {
				fmt.Printf("   [%s] %-5s %s  %s\n", child.Mode, child.Type, child.Hash[:8], child.Name)
			}
		}
	}
}
