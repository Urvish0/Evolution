package repository

import "fmt"

type Repository struct {
	Config Config
	Branch Branch
}

func OpenRepository() (Repository, error) {
	config, err := LoadConfig()
	if err != nil {
		return Repository{}, fmt.Errorf("loading config: %w", err)
	}

	currentBranch, err := GetCurrentBranchName()
	if err != nil {
		return Repository{}, fmt.Errorf("reading current branch name: %w", err)
	}

	branch, err := LoadBranch(currentBranch)
	if err != nil {
		return Repository{}, fmt.Errorf("loading branch %s: %w", currentBranch, err)
	}

	return Repository{
		Config: config,
		Branch: branch,
	}, nil

}
