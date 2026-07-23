package repository

import (
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
