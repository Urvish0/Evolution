package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Branch struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

func NewBranch(name string) Branch {
	return Branch{
		Name: name,
		Head: "",
	}
}

func LoadBranch(name string) (Branch, error) {
	var branch Branch

	data, err := os.ReadFile(filepath.Join(RepositoryDir, BranchesDir, name))

	if err != nil {
		return branch, fmt.Errorf("reading branch: %w", err)
	}

	err = json.Unmarshal(data, &branch)
	if err != nil {
		return branch, fmt.Errorf("parsing branch: %w", err)
	}

	return branch, nil
}

func writeBranch(branch Branch) error {
	data, err := json.MarshalIndent(branch, "", " ")
	if err != nil {
		return fmt.Errorf("writing branch: %w", err)
	}

	branchPath := filepath.Join(RepositoryDir, BranchesDir, branch.Name)
	return os.WriteFile(branchPath, []byte(data), 0644)
}

func UpdateBranch(branch Branch) error {
	data, err := json.MarshalIndent(branch, "", " ")
	if err != nil {
		return fmt.Errorf("updating branch: %w", err)
	}

	branchPath := filepath.Join(RepositoryDir, BranchesDir, branch.Name)
	return os.WriteFile(branchPath, data, 0644)
}

func GetCurrentBranchName() (string, error) {
	headPath := filepath.Join(RepositoryDir, HeadFile)
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "",fmt.Errorf("reading HEAD file: %w", err)
	}
	return string(data), nil
}

func ListBranches() ([]Branch, error) {
	branchesDir := filepath.Join(RepositoryDir, BranchesDir)
	entries, err := os.ReadDir(branchesDir)
	if err != nil {
		return nil, fmt.Errorf("reading branches directory: %w", err)
	}

	var branches []Branch
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		branch, err := LoadBranch(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("loading branch %s: %w", entry.Name(), err)
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// CreateBranch creates a new branch pointing to the current HEAD commit.
func CreateBranch(name string) error {
	branchPath := filepath.Join(RepositoryDir, BranchesDir, name)
	if _, err := os.Stat(branchPath); err == nil {
		return fmt.Errorf("branch %q already exists", name)
	}

	repo, err := OpenRepository()
	if err != nil {
		return fmt.Errorf("opening repository: %w", err)
	}

	branch := NewBranch(name)
	branch.Head = repo.Branch.Head

	return writeBranch(branch)
}

// CheckoutBranch switches the active branch by updating .evolution/HEAD and restoring workspace files.
func CheckoutBranch(name string) error {
	branch, err := LoadBranch(name)
	if err != nil {
		return fmt.Errorf("cannot checkout %q: %w", name, err)
	}

	headPath := filepath.Join(RepositoryDir, HeadFile)
	if err := os.WriteFile(headPath, []byte(name), 0644); err != nil {
		return fmt.Errorf("updating HEAD file: %w", err)
	}

	// Restore workspace files from branch HEAD commit
	if branch.Head != "" {
		commit, err := LoadCommit(branch.Head)
		if err == nil && commit.TreeHash != "" {
			_ = RestoreTreeToWorkspace(commit.TreeHash, ".")
		}
	}

	return nil
}

// CheckoutCommit switches to detached HEAD mode pointing directly to a commit ID.
func CheckoutCommit(commitID string) error {
	// Support short commit ID prefixes (e.g. first 8 chars)
	realID := commitID
	if len(commitID) < 36 {
		// Search commits directory for matching prefix
		commitsDir := filepath.Join(RepositoryDir, CommitsDir)
		entries, err := os.ReadDir(commitsDir)
		if err == nil {
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), commitID) {
					realID = strings.TrimSuffix(entry.Name(), ".json")
					break
				}
			}
		}
	}

	commit, err := LoadCommit(realID)
	if err != nil {
		return fmt.Errorf("checkout commit %s failed: %w", commitID, err)
	}

	headPath := filepath.Join(RepositoryDir, HeadFile)
	if err := os.WriteFile(headPath, []byte(realID), 0644); err != nil {
		return fmt.Errorf("updating HEAD file: %w", err)
	}

	if commit.TreeHash != "" {
		_ = RestoreTreeToWorkspace(commit.TreeHash, ".")
	}

	return nil
}

// DeleteBranch removes a branch file. Cannot delete the active branch.
func DeleteBranch(name string) error {
	currentBranch, err := GetCurrentBranchName()
	if err != nil {
		return fmt.Errorf("checking active branch: %w", err)
	}

	if currentBranch == name {
		return fmt.Errorf("cannot delete active branch %q", name)
	}

	branchPath := filepath.Join(RepositoryDir, BranchesDir, name)
	if err := os.Remove(branchPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("branch %q does not exist", name)
		}
		return fmt.Errorf("deleting branch %q: %w", name, err)
	}

	return nil
}

// BranchDetails contains rich metadata for a branch.
type BranchDetails struct {
	Name              string `json:"name"`
	IsActive          bool   `json:"is_active"`
	HeadCommitID      string `json:"head_commit_id"`
	LastCommitMessage string `json:"last_commit_message"`
	LastCommitDate    string `json:"last_commit_date"`
	CommitCount       int    `json:"commit_count"`
}

// GetBranchDetails retrieves rich metadata for a given branch by traversing its commit DAG.
func GetBranchDetails(b Branch) (BranchDetails, error) {
	currentBranch, _ := GetCurrentBranchName()

	details := BranchDetails{
		Name:         b.Name,
		IsActive:     b.Name == currentBranch,
		HeadCommitID: b.Head,
	}

	if b.Head == "" {
		return details, nil
	}

	// Traversal to count commits and extract latest commit info
	currID := b.Head
	count := 0
	for currID != "" {
		commit, err := LoadCommit(currID)
		if err != nil {
			break
		}

		if count == 0 {
			details.LastCommitMessage = commit.Message
			details.LastCommitDate = commit.Timestamp
		}

		count++
		currID = commit.Parent
	}

	details.CommitCount = count
	return details, nil
}

// RenameBranch renames an existing branch. Updates .evolution/HEAD if active branch is renamed.
func RenameBranch(oldName, newName string) error {
	if newName == "" {
		return fmt.Errorf("new branch name cannot be empty")
	}

	oldPath := filepath.Join(RepositoryDir, BranchesDir, oldName)
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("branch %q does not exist", oldName)
	}

	newPath := filepath.Join(RepositoryDir, BranchesDir, newName)
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("branch %q already exists", newName)
	}

	branch, err := LoadBranch(oldName)
	if err != nil {
		return fmt.Errorf("loading branch %q: %w", oldName, err)
	}

	branch.Name = newName
	if err := writeBranch(branch); err != nil {
		return fmt.Errorf("writing renamed branch %q: %w", newName, err)
	}

	_ = os.Remove(oldPath)

	// If old branch was the active branch, update HEAD pointer
	currentBranch, _ := GetCurrentBranchName()
	if currentBranch == oldName {
		_ = CheckoutBranch(newName)
	}

	return nil
}