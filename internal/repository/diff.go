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

// GetCommittedFileHash returns the blob hash for a file path from the HEAD commit tree.
func GetCommittedFileHash(filePath string) (string, error) {
	repo, err := OpenRepository()
	if err != nil {
		return "", fmt.Errorf("opening repository: %w", err)
	}

	if repo.Branch.Head == "" {
		return "", fmt.Errorf("no commits yet")
	}

	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		return "", fmt.Errorf("loading HEAD commit: %w", err)
	}

	if commit.TreeHash == "" {
		return "", fmt.Errorf("HEAD commit has no tree hash")
	}

	committedFiles := make(map[string]string)
	collectTreeEntries(commit.TreeHash, "", committedFiles)

	hash, exists := committedFiles[filePath]
	if !exists {
		return "", fmt.Errorf("file %s not found in HEAD commit", filePath)
	}

	return hash, nil
}

// GetContentDiff produces a line-by-line diff for a file between its committed version and current disk version.
func GetContentDiff(filePath string) ([]DiffLine, error) {
	// Read old content from committed blob
	committedHash, err := GetCommittedFileHash(filePath)
	if err != nil {
		return nil, err
	}

	oldContent, err := ReadBlob(committedHash)
	if err != nil {
		return nil, fmt.Errorf("reading committed blob for %s: %w", filePath, err)
	}

	// Read new content from disk
	newContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading workspace file %s: %w", filePath, err)
	}

	return ComputeLineDiff(string(oldContent), string(newContent)), nil
}

// ResolveRevisionToCommit resolves a branch name, full commit UUID, or short 8-char commit prefix into a Commit struct.
func ResolveRevisionToCommit(rev string) (Commit, error) {
	var commit Commit

	// 1. Try branch name
	branch, err := LoadBranch(rev)
	if err == nil && branch.Head != "" {
		return LoadCommit(branch.Head)
	}

	// 2. Try exact commit ID
	commit, err = LoadCommit(rev)
	if err == nil {
		return commit, nil
	}

	// 3. Try short commit ID prefix match
	if len(rev) < 36 {
		commitsDir := filepath.Join(RepositoryDir, CommitsDir)
		entries, err := os.ReadDir(commitsDir)
		if err == nil {
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), rev) {
					id := strings.TrimSuffix(entry.Name(), ".json")
					return LoadCommit(id)
				}
			}
		}
	}

	return commit, fmt.Errorf("revision %q not found", rev)
}

// GetRevisionDiff returns a formatted unified diff comparing two revisions (branches or commits).
func GetRevisionDiff(rev1, rev2 string) (string, error) {
	c1, err := ResolveRevisionToCommit(rev1)
	if err != nil {
		return "", fmt.Errorf("resolving revision %s: %w", rev1, err)
	}

	c2, err := ResolveRevisionToCommit(rev2)
	if err != nil {
		return "", fmt.Errorf("resolving revision %s: %w", rev2, err)
	}

	files1, _ := GetFlatTreeEntries(c1.TreeHash, "")
	files2, _ := GetFlatTreeEntries(c2.TreeHash, "")

	allFiles := make(map[string]bool)
	for f := range files1 {
		allFiles[f] = true
	}
	for f := range files2 {
		allFiles[f] = true
	}

	var buf strings.Builder
	for file := range allFiles {
		h1 := files1[file]
		h2 := files2[file]

		if h1 == h2 {
			continue // No change
		}

		var content1, content2 []byte
		if h1 != "" {
			content1, _ = ReadBlob(h1)
		}
		if h2 != "" {
			content2, _ = ReadBlob(h2)
		}

		diffLines := ComputeLineDiff(string(content1), string(content2))
		buf.WriteString(fmt.Sprintf("--- a/%s (%s)\n", file, rev1[:min(8, len(rev1))]))
		buf.WriteString(fmt.Sprintf("+++ b/%s (%s)\n", file, rev2[:min(8, len(rev2))]))
		buf.WriteString("@@ -1 +1 @@\n")

		for _, dl := range diffLines {
			switch dl.Op {
			case DiffInsert:
				buf.WriteString(fmt.Sprintf("+%s\n", dl.Content))
			case DiffDelete:
				buf.WriteString(fmt.Sprintf("-%s\n", dl.Content))
			case DiffEqual:
				buf.WriteString(fmt.Sprintf(" %s\n", dl.Content))
			}
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
