package repository

type Repository struct {
	Config Config
	Branch Branch
}

func OpenRepository() (Repository, error) {
	config, err := LoadConfig()
	if err != nil {
		return Repository{}, err
	}

	branch, err := LoadBranch(config.DefaultBranch)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		Config: config,
		Branch: branch,
	}, nil

}
