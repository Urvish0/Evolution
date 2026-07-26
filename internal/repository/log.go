package repository

func Log() ([]Commit, error) {

	repo, err := OpenRepository()
	if err != nil {
		return nil, err
	}

	if repo.Branch.Head == "" {
		return []Commit{}, nil
	}

	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		return nil, err
	}

	return []Commit{commit}, nil

}
