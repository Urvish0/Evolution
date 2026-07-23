package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func ensureNotInitialized() error {
	_, err := os.Stat(RepositoryDir)

	if err == nil {
		return errors.New("Evolution is already initialized in this directory.")
	}

	if !os.IsNotExist(err) {
		return err
	}

	return nil
}
func createRepositoryDirectory() error {
	return os.Mkdir(RepositoryDir, 0755)
}

func writeConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(RepositoryDir, ConfigFile)

	return os.WriteFile(configPath, data, 0644)
}
func writeHead() error {
	headPath := filepath.Join(RepositoryDir, HeadFile)

	return os.WriteFile(headPath, []byte(DefaultBranch), 0644)
}
