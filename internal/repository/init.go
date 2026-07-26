package repository

import "os"

func Exists() bool {
	info, err := os.Stat(RepositoryDir)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func Init() error {
	if err := ensureNotInitialized(); err != nil {
		return err
	}

	if err := createRepositoryDirectory(); err != nil {
		return err
	}

	config := NewConfig()

	if err := writeConfig(config); err != nil {
		return err
	}

	if err := writeHead(); err != nil {
		return err
	}

	branch := NewBranch(DefaultBranch)

	if err := writeBranch(branch); err != nil {
		return err
	}

	return nil
}
