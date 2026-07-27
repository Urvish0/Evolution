package repository

import "fmt"

func Log() ([]Commit, error) {

	repo, err := OpenRepository()
	if err != nil {
		return nil, fmt.Errorf("opening repository for log: %w", err)
	}

	if repo.Branch.Head == "" {
		return []Commit{}, nil
	}

	var commits []Commit
	currentID := repo.Branch.Head

	for currentID != "" {
		commit, err := LoadCommit(currentID)
		if err != nil {
			return nil, fmt.Errorf("loading commit %s: %w", currentID[:8], err)
		}

		commits = append(commits, commit)
		currentID = commit.Parent
	}

	return commits, nil

}
