package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetGlobalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting user home directory: %w", err)
	}

	evoHome := filepath.Join(homeDir, ".evolution")
	if err := os.MkdirAll(evoHome, 0755); err != nil {
		return "", fmt.Errorf("creating global config directory: %w", err)
	}

	return filepath.Join(evoHome, "config.json"), nil
}

func LoadUserConfig() (UserConfig, error) {
	var cfg UserConfig

	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading global user config: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing global user config: %w", err)
	}

	return cfg, nil
}

func SaveUserConfig(cfg UserConfig) error {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing global user config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing global user config: %w", err)
	}

	return nil
}

func SetUserSetting(key, value string) error {
	cfg, err := LoadUserConfig()
	if err != nil {
		return err
	}

	switch key {
	case "user.name":
		cfg.Name = value
	case "user.email":
		cfg.Email = value
	default:
		return fmt.Errorf("unknown config key %q (valid keys: user.name, user.email)", key)
	}

	return SaveUserConfig(cfg)
}
