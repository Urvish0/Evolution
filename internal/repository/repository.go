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

	branch, err := LoadBranch(config.DefaultBranch)
	if err != nil {
		return Repository{}, fmt.Errorf("loading branch %s: %w", config.DefaultBranch, err)
	}

	return Repository{
		Config: config,
		Branch: branch,
	}, nil

}
