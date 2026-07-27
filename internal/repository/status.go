package repository

import "fmt"

type Status struct {
	Initialized bool
	Branch      string
	Commits     int
	Clean       bool
}

func GetStatus() (Status, error) {
	repo, err := OpenRepository()
	if err != nil {
		return Status{}, fmt.Errorf("opening repository for status: %w", err)
	}

	commits, err := Log()
	if err != nil {
		return Status{}, fmt.Errorf("loading commit history: %w", err)
	}

	count := len(commits)

	return Status{
		Initialized: true,
		Branch:      repo.Branch.Name,
		Commits:     count,
		Clean:       true,
	}, nil
}
