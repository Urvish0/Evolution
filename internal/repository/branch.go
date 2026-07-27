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
