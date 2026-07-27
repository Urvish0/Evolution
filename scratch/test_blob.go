package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Urvish0/evolution/internal/repository"
)

func main() {
	// 1. Ensure repo is initialized
	if !repository.Exists() {
		if err := repository.Init(); err != nil {
			fmt.Printf("Error initializing repo: %v\n", err)
			return
		}
		fmt.Println("Initialized repository for testing.")
	}

	// 2. Write a system prompt blob
	promptContent := []byte("System Prompt: You are a legal research assistant AI.")
	hash, err := repository.WriteBlob(promptContent)
	if err != nil {
		fmt.Printf("Error writing blob: %v\n", err)
		return
	}

	fmt.Printf("   Blob created successfully!\n")
	fmt.Printf("   SHA-256 Hash: %s\n", hash)
	fmt.Printf("   Object Path:  .evolution/objects/%s/%s\n\n", hash[:2], hash[2:])

	// 3. Test reading it back
	readBack, err := repository.ReadBlob(hash)
	if err != nil {
		fmt.Printf("Error reading blob: %v\n", err)
		return
	}

	fmt.Printf("  Read Blob Content Back:\n   %q\n\n", string(readBack))

	// 4. Test deduplication
	hashDup, _ := repository.WriteBlob(promptContent)
	if hashDup == hash {
		fmt.Println("Deduplication working! Re-writing same content returned same hash without duplicating files.")
	}

	// Clean up test object folder afterwards
	dir := filepath.Join(repository.RepositoryDir, repository.ObjectsDir, hash[:2])
	_ = os.RemoveAll(dir)
}
