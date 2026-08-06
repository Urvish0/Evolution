package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MergeResult describes the outcome of a branch merge operation.
type MergeResult struct {
	IsFastForward  bool     `json:"is_fast_forward"`
	HasConflicts   bool     `json:"has_conflicts"`
	ConflictedFiles []string `json:"conflicted_files,omitempty"`
	MergeCommitID  string   `json:"merge_commit_id,omitempty"`
}

// FindLowestCommonAncestor finds the shared ancestor commit ID where two commit DAGs diverged.
func FindLowestCommonAncestor(commitID1, commitID2 string) (string, error) {
	if commitID1 == "" || commitID2 == "" {
		return "", nil
	}
	if commitID1 == commitID2 {
		return commitID1, nil
	}

	ancestors1 := make(map[string]bool)
	curr := commitID1
	for curr != "" {
		ancestors1[curr] = true
		c, err := LoadCommit(curr)
		if err != nil {
			break
		}
		curr = c.Parent
	}

	curr = commitID2
	for curr != "" {
		if ancestors1[curr] {
			return curr, nil
		}
		c, err := LoadCommit(curr)
		if err != nil {
			break
		}
		curr = c.Parent
	}

	return "", nil
}

// ThreeWayLineMerge performs line-by-line three-way content merge (Base, Ours, Theirs).
// If a line is modified differently in Ours vs Theirs, conflict markers are injected.
func ThreeWayLineMerge(base, ours, theirs []byte, oursLabel, theirsLabel string) ([]byte, bool) {
	baseLines := splitLines(string(base))
	oursLines := splitLines(string(ours))
	theirsLines := splitLines(string(theirs))

	// Fast checks
	if bytes.Equal(ours, theirs) {
		return ours, false
	}
	if bytes.Equal(base, ours) {
		return theirs, false
	}
	if bytes.Equal(base, theirs) {
		return ours, false
	}

	// Line-by-line comparison
	var result []string
	hasConflict := false

	maxLen := max(len(baseLines), max(len(oursLines), len(theirsLines)))

	for i := 0; i < maxLen; i++ {
		var baseLine, oursLine, theirsLine string
		var hasBase, hasOurs, hasTheirs bool

		if i < len(baseLines) {
			baseLine = baseLines[i]
			hasBase = true
		}
		if i < len(oursLines) {
			oursLine = oursLines[i]
			hasOurs = true
		}
		if i < len(theirsLines) {
			theirsLine = theirsLines[i]
			hasTheirs = true
		}

		if hasOurs && hasTheirs && oursLine == theirsLine {
			result = append(result, oursLine)
		} else if hasBase && hasOurs && baseLine == oursLine && hasTheirs {
			result = append(result, theirsLine)
		} else if hasBase && hasTheirs && baseLine == theirsLine && hasOurs {
			result = append(result, oursLine)
		} else if !hasBase && hasOurs && !hasTheirs {
			result = append(result, oursLine)
		} else if !hasBase && !hasOurs && hasTheirs {
			result = append(result, theirsLine)
		} else {
			// Conflict detected!
			hasConflict = true
			result = append(result, fmt.Sprintf("<<<<<<< OURS (%s)", oursLabel))
			if hasOurs {
				result = append(result, oursLine)
			}
			result = append(result, "=======")
			if hasTheirs {
				result = append(result, theirsLine)
			}
			result = append(result, fmt.Sprintf(">>>>>>> THEIRS (%s)", theirsLabel))
		}
	}

	output := strings.Join(result, "\n")
	if len(result) > 0 && !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	return []byte(output), hasConflict
}

// GetFlatTreeEntries recursively flattens a Merkle tree into a map of relative file path -> blob hash.
func GetFlatTreeEntries(treeHash, prefix string) (map[string]string, error) {
	result := make(map[string]string)
	if treeHash == "" {
		return result, nil
	}

	tree, err := ReadTree(treeHash)
	if err != nil {
		return nil, fmt.Errorf("reading tree %s: %w", treeHash[:8], err)
	}

	for _, entry := range tree.Entries {
		relPath := filepath.Join(prefix, entry.Name)
		if entry.Type == ObjectTypeBlob {
			result[relPath] = entry.Hash
		} else if entry.Type == ObjectTypeTree {
			subMap, err := GetFlatTreeEntries(entry.Hash, relPath)
			if err != nil {
				return nil, err
			}
			for k, v := range subMap {
				result[k] = v
			}
		}
	}

	return result, nil
}

// MergeBranch merges targetBranchName into the current active branch.
func MergeBranch(targetBranchName string) (*MergeResult, error) {
	currentBranchName, err := GetCurrentBranchName()
	if err != nil {
		return nil, fmt.Errorf("getting current branch: %w", err)
	}

	if currentBranchName == targetBranchName {
		return nil, fmt.Errorf("cannot merge branch %q into itself", targetBranchName)
	}

	oursBranch, err := LoadBranch(currentBranchName)
	if err != nil {
		return nil, fmt.Errorf("loading current branch %q: %w", currentBranchName, err)
	}

	theirsBranch, err := LoadBranch(targetBranchName)
	if err != nil {
		return nil, fmt.Errorf("loading target branch %q: %w", targetBranchName, err)
	}

	if theirsBranch.Head == "" {
		return &MergeResult{IsFastForward: true}, nil
	}

	if oursBranch.Head == "" {
		// Fast-forward: ours has no commits, just update ours head to theirs head
		oursBranch.Head = theirsBranch.Head
		if err := UpdateBranch(oursBranch); err != nil {
			return nil, err
		}
		_ = CheckoutBranch(currentBranchName)
		return &MergeResult{IsFastForward: true}, nil
	}

	// Find Lowest Common Ancestor
	lcaID, err := FindLowestCommonAncestor(oursBranch.Head, theirsBranch.Head)
	if err != nil {
		return nil, fmt.Errorf("finding lowest common ancestor: %w", err)
	}

	// Check Fast-Forward cases
	if lcaID == oursBranch.Head {
		// Fast-forward merge! Ours is direct ancestor of theirs
		oursBranch.Head = theirsBranch.Head
		if err := UpdateBranch(oursBranch); err != nil {
			return nil, err
		}
		_ = CheckoutBranch(currentBranchName)
		return &MergeResult{IsFastForward: true}, nil
	}

	if lcaID == theirsBranch.Head {
		// Already up to date!
		return &MergeResult{IsFastForward: true}, nil
	}

	// 3-Way Merge required
	lcaCommit, _ := LoadCommit(lcaID)
	oursCommit, _ := LoadCommit(oursBranch.Head)
	theirsCommit, _ := LoadCommit(theirsBranch.Head)

	baseFiles, _ := GetFlatTreeEntries(lcaCommit.TreeHash, "")
	oursFiles, _ := GetFlatTreeEntries(oursCommit.TreeHash, "")
	theirsFiles, _ := GetFlatTreeEntries(theirsCommit.TreeHash, "")

	allFiles := make(map[string]bool)
	for f := range baseFiles {
		allFiles[f] = true
	}
	for f := range oursFiles {
		allFiles[f] = true
	}
	for f := range theirsFiles {
		allFiles[f] = true
	}

	var conflictedFiles []string
	hasConflicts := false

	for file := range allFiles {
		baseHash := baseFiles[file]
		oursHash := oursFiles[file]
		theirsHash := theirsFiles[file]

		if oursHash == theirsHash {
			continue // No change or identical change
		}

		var baseContent, oursContent, theirsContent []byte
		if baseHash != "" {
			baseContent, _ = ReadBlob(baseHash)
		}
		if oursHash != "" {
			oursContent, _ = ReadBlob(oursHash)
		}
		if theirsHash != "" {
			theirsContent, _ = ReadBlob(theirsHash)
		}

		merged, conflict := ThreeWayLineMerge(baseContent, oursContent, theirsContent, currentBranchName, targetBranchName)

		// Write merged/conflicted content to working tree file
		dir := filepath.Dir(file)
		if dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		_ = os.WriteFile(file, merged, 0644)
		_ = StagePath(file)

		if conflict {
			hasConflicts = true
			conflictedFiles = append(conflictedFiles, file)
		}
	}

	result := &MergeResult{
		IsFastForward:   false,
		HasConflicts:    hasConflicts,
		ConflictedFiles: conflictedFiles,
	}

	if !hasConflicts {
		// Create automatic merge commit
		msg := fmt.Sprintf("Merge branch '%s' into %s", targetBranchName, currentBranchName)
		if err := CreateCommit(msg); err != nil {
			return nil, fmt.Errorf("creating merge commit: %w", err)
		}
		repo, _ := OpenRepository()
		result.MergeCommitID = repo.Branch.Head
	}

	return result, nil
}
