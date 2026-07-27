package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func ensureNotInitialized() error {
	_, err := os.Stat(RepositoryDir)

	if err == nil {
		return errors.New("Evolution is already initialized in this directory.")
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("checking repository directory: %w", err)
	}

	return nil
}
func createRepositoryDirectory() error {
	directories := []string{
		RepositoryDir,
		filepath.Join(RepositoryDir, BranchesDir),
		filepath.Join(RepositoryDir, CommitsDir),
		filepath.Join(RepositoryDir, ObjectsDir),
		filepath.Join(RepositoryDir, ArtifactsDir),
	}

	for _, dir := range directories {
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}

func writeHead() error {
	headPath := filepath.Join(RepositoryDir, HeadFile)

	return os.WriteFile(headPath, []byte(DefaultBranch), 0644)
}
