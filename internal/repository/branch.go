package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Branch struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

func NewBranch(name string) Branch {
	return Branch{
		Name: DefaultBranch,
		Head: "",
	}
}

func LoadBranch(name string) (Branch, error) {
	var branch Branch

	data, err := os.ReadFile(filepath.Join(RepositoryDir, BranchesDir, DefaultBranch))

	if err != nil {
		return branch, err
	}

	err = json.Unmarshal(data, &branch)
	if err != nil {
		return branch, err
	}

	return branch, nil
}

func writeBranch(branch Branch) error {
	data, err := json.MarshalIndent(branch, "", " ")
	if err != nil {
		return err
	}

	branchPath := filepath.Join(RepositoryDir, BranchesDir, DefaultBranch)
	return os.WriteFile(branchPath, []byte(data), 0644)
}

func UpdateBranch(branch Branch) error {
	data, err := json.MarshalIndent(branch, "", " ")
	if err != nil {
		return err
	}

	branchPath := filepath.Join(RepositoryDir, BranchesDir, branch.Name)
	return os.WriteFile(branchPath, data, 0644)
}
