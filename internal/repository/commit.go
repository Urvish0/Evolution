package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Commit struct {
	ID        string `json:"id"`
	Parent    string `json:"parent"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}

func NewCommit(message string) Commit {
	return Commit{
		ID:        uuid.New().String(),
		Parent:    "",
		Message:   message,
		Author:    "",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
func LoadCommit(id string) (Commit, error) {
	var commit Commit

	data, err := os.ReadFile(filepath.Join(RepositoryDir, CommitsDir, id+".json"))
	if err != nil {
		return commit, fmt.Errorf("reading commit: %w", err)
	}

	err = json.Unmarshal(data, &commit)
	if err != nil {
		return commit, fmt.Errorf("parsing commit: %w", err)
	}

	return commit, nil
}

func WriteCommit(commit Commit) error {
	data, err := json.MarshalIndent(commit, "", " ")
	if err != nil {
		return fmt.Errorf("writing commit: %w", err)
	}

	commitPath := filepath.Join(RepositoryDir, CommitsDir, commit.ID+".json")
	return os.WriteFile(commitPath, data, 0644)
}
func CreateCommit(message string) error {

	repo, err := OpenRepository()
	if err != nil {
		return fmt.Errorf("opening repo for commit: %w", err)
	}

	commit := NewCommit(message)
	commit.Parent = repo.Branch.Head

	userCfg, _ := LoadUserConfig()
	if userCfg.Name != "" {
		if userCfg.Email != "" {
			commit.Author = fmt.Sprintf("%s <%s>", userCfg.Name, userCfg.Email)
		} else {
			commit.Author = userCfg.Name
		}
	}

	if err := WriteCommit(commit); err != nil {
		return fmt.Errorf("writing commit: %w", err)
	}

	branch := repo.Branch
	branch.Head = commit.ID

	if err := UpdateBranch(branch); err != nil {
		return fmt.Errorf("updating branch for commit: %w", err)
	}

	return nil

}
