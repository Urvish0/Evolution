package repository

import (
	"encoding/json"
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
		return commit, err
	}

	err = json.Unmarshal(data, &commit)
	if err != nil {
		return commit, err
	}

	return commit, nil
}

func WriteCommit(commit Commit) error {
	data, err := json.MarshalIndent(commit, "", " ")
	if err != nil {
		return err
	}

	commitPath := filepath.Join(RepositoryDir, CommitsDir, commit.ID+".json")
	return os.WriteFile(commitPath, data, 0644)
}
func CreateCommit(message string) error {

	repo, err := OpenRepository()
	if err != nil {
		return err
	}

	commit := NewCommit(message)
	commit.Parent = repo.Branch.Head

	if err := WriteCommit(commit); err != nil {
		return err
	}

	branch := repo.Branch
	branch.Head = commit.ID

	if err := UpdateBranch(branch); err != nil {
		return err
	}

	return nil

}
