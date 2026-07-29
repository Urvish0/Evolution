package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileChange represents the status of a single file relative to the last commit.
type FileChange struct {
	Path   string
	Status ChangeStatus
}

// ChangeStatus represents what happened to a file.
type ChangeStatus string

const (
	StatusNew      ChangeStatus = "new"
	StatusModified ChangeStatus = "modified"
	StatusDeleted  ChangeStatus = "deleted"
	StatusStaged   ChangeStatus = "staged"
)

// WorkingTreeStatus holds the complete diff between workspace, index, and HEAD.
type WorkingTreeStatus struct {
	Staged    []FileChange
	Modified  []FileChange
	Untracked []FileChange
	Deleted   []FileChange
}

// CompareWorkingTree compares the current workspace files against the last committed
// tree snapshot to detect new, modified, and deleted files.
func CompareWorkingTree() (WorkingTreeStatus, error) {
	var wts WorkingTreeStatus

	// 1. Collect all committed file hashes from HEAD tree
	committedFiles := make(map[string]string) // path -> blob hash
	repo, err := OpenRepository()
	if err == nil && repo.Branch.Head != "" {
		commit, err := LoadCommit(repo.Branch.Head)
		if err == nil && commit.TreeHash != "" {
			collectTreeEntries(commit.TreeHash, "", committedFiles)
		}
	}

	// 2. Collect staged file hashes from index
	stagedFiles := make(map[string]string) // path -> blob hash
	idx, err := LoadIndex()
	if err == nil {
		for path, entry := range idx.Entries {
			stagedFiles[path] = entry.Hash
		}
	}

	// 3. Walk workspace and collect current file hashes
	workspaceFiles := make(map[string]string) // path -> blob hash
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := filepath.ToSlash(path)

		// Skip .evolution directory and hidden files/dirs
		if relPath == RepositoryDir || strings.HasPrefix(relPath, RepositoryDir+"/") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasPrefix(filepath.Base(relPath), ".") {
			if info.IsDir() && relPath != "." {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil // Skip unreadable files
			}
			hash := HashContent(ObjectTypeBlob, content)
			workspaceFiles[relPath] = hash
		}

		return nil
	})
	if err != nil {
		return wts, fmt.Errorf("walking workspace: %w", err)
	}

	// 4. Classify each file

	// 4a. Staged files (in index but not yet committed, or hash differs from HEAD)
	for path, stagedHash := range stagedFiles {
		committedHash, inCommit := committedFiles[path]
		if !inCommit || stagedHash != committedHash {
			wts.Staged = append(wts.Staged, FileChange{Path: path, Status: StatusStaged})
		}
	}

	// 4b. Modified files (in workspace with different hash than committed, and NOT staged)
	for path, workspaceHash := range workspaceFiles {
		committedHash, inCommit := committedFiles[path]
		_, inStaged := stagedFiles[path]

		if inCommit && workspaceHash != committedHash && !inStaged {
			wts.Modified = append(wts.Modified, FileChange{Path: path, Status: StatusModified})
		}
	}

	// 4c. Untracked files (in workspace but not in committed tree AND not staged)
	for path := range workspaceFiles {
		_, inCommit := committedFiles[path]
		_, inStaged := stagedFiles[path]

		if !inCommit && !inStaged {
			wts.Untracked = append(wts.Untracked, FileChange{Path: path, Status: StatusNew})
		}
	}

	// 4d. Deleted files (in committed tree but not in workspace)
	for path := range committedFiles {
		if _, exists := workspaceFiles[path]; !exists {
			wts.Deleted = append(wts.Deleted, FileChange{Path: path, Status: StatusDeleted})
		}
	}

	return wts, nil
}

// collectTreeEntries recursively flattens a Tree object into a map of path -> blob hash.
func collectTreeEntries(treeHash string, prefix string, result map[string]string) {
	tree, err := ReadTree(treeHash)
	if err != nil {
		return
	}

	for _, entry := range tree.Entries {
		fullPath := entry.Name
		if prefix != "" {
			fullPath = prefix + "/" + entry.Name
		}

		if entry.Type == ObjectTypeTree {
			collectTreeEntries(entry.Hash, fullPath, result)
		} else {
			result[fullPath] = entry.Hash
		}
	}
}

// IsClean returns true if there are no changes in the working tree.
func (wts WorkingTreeStatus) IsClean() bool {
	return len(wts.Staged) == 0 && len(wts.Modified) == 0 && len(wts.Untracked) == 0 && len(wts.Deleted) == 0
}
