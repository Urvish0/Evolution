package repository

import (
	"fmt"
	"os"
)

func Init() error {
	if err := ensureNotInitialized(); err != nil {
		return fmt.Errorf("checking initialization: %w", err)
	}

	if err := createRepositoryDirectory(); err != nil {
		return fmt.Errorf("creating repository directories: %w", err)
	}

	config := NewConfig()

	if err := writeConfig(config); err != nil {
		return fmt.Errorf("writing initial config: %w", err)
	}

	if err := writeHead(); err != nil {
		return fmt.Errorf("writing HEAD: %w", err)
	}

	branch := NewBranch(DefaultBranch)

	if err := writeBranch(branch); err != nil {
		return fmt.Errorf("writing initial branch: %w", err)
	}

	return nil
}

func Exists() bool {
	info, err := os.Stat(RepositoryDir)
	if err != nil {
		return false
	}

	return info.IsDir()
}
