package repository

type Status struct {
	Initialized bool
	Branch      string
	Commits     int
	Clean       bool
}

func GetStatus() (Status, error) {
	config, err := LoadConfig()
	if err != nil {
		return Status{}, err
	}

	return Status{
		Initialized: true,
		Branch:      config.DefaultBranch,
		Commits:     0,
		Clean:       true,
	}, nil
}
