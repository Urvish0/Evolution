package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile(filepath.Join(RepositoryDir, ConfigFile))
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
