package repository

type Status struct {
	Initialized bool
	Branch      string
	Commits     int
	Clean       bool
}

func GetStatus() (Status, error) {
	repo, err := OpenRepository()
	if err != nil {
		return Status{}, err
	}

	return Status{
		Initialized: true,
		Branch:      repo.Branch.Name,
		Commits:     0,
		Clean:       true,
	}, nil
}
