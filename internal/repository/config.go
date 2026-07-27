package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	Version       int    `json:"version"`
	RepositoryID  string `json:"repository_id"`
	CreatedAt     string `json:"created_at"`
	DefaultBranch string `json:"default_branch"`
}

func NewConfig() Config {
	return Config{
		Version:       1,
		RepositoryID:  uuid.New().String(),
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		DefaultBranch: DefaultBranch,
	}
}
func LoadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile(filepath.Join(RepositoryDir, ConfigFile))
	if err != nil {
		return config, fmt.Errorf("reading config: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("parsing config: %w", err)
	}

	return config, nil
}
func writeConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("serealizing config: %w", err)
	}

	configPath := filepath.Join(RepositoryDir, ConfigFile)

	return os.WriteFile(configPath, data, 0644)
}
