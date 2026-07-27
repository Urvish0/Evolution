package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// CheckoutBranch switches the active branch by updating the .evolution/HEAD file.
func CheckoutBranch(name string) error {
	_, err := LoadBranch(name)
	if err != nil {
		return fmt.Errorf("cannot checkout %q: %w", name, err)
	}

	headPath := filepath.Join(RepositoryDir, HeadFile)
	if err := os.WriteFile(headPath, []byte(name), 0644); err != nil {
		return fmt.Errorf("updating HEAD file: %w", err)
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