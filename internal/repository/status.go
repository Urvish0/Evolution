package repository

import "fmt"

// Status represents the current state of an Evolution repository.
type Status struct {
	Initialized bool
	Branch      string
	Commits     int
	Clean       bool
	WorkingTree WorkingTreeStatus
}

// GetStatus returns the full repository status including working tree changes.
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

	wts, err := CompareWorkingTree()
	if err != nil {
		return Status{}, fmt.Errorf("comparing working tree: %w", err)
	}

	return Status{
		Initialized: true,
		Branch:      repo.Branch.Name,
		Commits:     count,
		Clean:       wts.IsClean(),
		WorkingTree: wts,
	}, nil
}
