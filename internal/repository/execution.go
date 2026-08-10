package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TokenUsage captures token consumption metrics for an AI execution run.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Execution represents a single recorded AI execution run linked to an Intelligence Commit.
type Execution struct {
	ID         string            `json:"id"`
	CommitID   string            `json:"commit_id"`
	Inputs     string            `json:"inputs"`
	Outputs    string            `json:"outputs"`
	DurationMs int64             `json:"duration_ms"`
	Tokens     TokenUsage        `json:"tokens"`
	Timestamp  string            `json:"timestamp"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// NewExecution initializes an Execution struct with a new UUID and current timestamp.
func NewExecution(commitID, inputs, outputs string, durationMs int64, tokens TokenUsage, meta map[string]string) Execution {
	if tokens.TotalTokens == 0 {
		tokens.TotalTokens = tokens.PromptTokens + tokens.CompletionTokens
	}
	if meta == nil {
		meta = make(map[string]string)
	}

	return Execution{
		ID:         uuid.New().String(),
		CommitID:   commitID,
		Inputs:     inputs,
		Outputs:    outputs,
		DurationMs: durationMs,
		Tokens:     tokens,
		Timestamp:  time.Now().Format(time.RFC3339),
		Metadata:   meta,
	}
}

// GetExecutionPath returns the disk path for a given execution ID (.evolution/executions/<id>.json).
func GetExecutionPath(id string) string {
	return filepath.Join(RepositoryDir, ExecutionsDir, id+".json")
}

// SaveExecution writes an Execution struct to disk as pretty-printed JSON.
func SaveExecution(exec Execution) error {
	dir := filepath.Join(RepositoryDir, ExecutionsDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating executions directory: %w", err)
	}

	data, err := json.MarshalIndent(exec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling execution: %w", err)
	}

	filePath := GetExecutionPath(exec.ID)
	return os.WriteFile(filePath, data, 0644)
}

// LoadExecution reads and parses an Execution struct from disk by ID or short prefix.
func LoadExecution(id string) (Execution, error) {
	var exec Execution
	realID := id

	// Handle short prefix matching
	if len(id) < 36 {
		dir := filepath.Join(RepositoryDir, ExecutionsDir)
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), id) {
					realID = strings.TrimSuffix(entry.Name(), ".json")
					break
				}
			}
		}
	}

	filePath := GetExecutionPath(realID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return exec, fmt.Errorf("reading execution %s: %w", id, err)
	}

	if err := json.Unmarshal(data, &exec); err != nil {
		return exec, fmt.Errorf("parsing execution JSON: %w", err)
	}

	return exec, nil
}

// ListExecutions lists all recorded executions sorted by timestamp descending.
func ListExecutions() ([]Execution, error) {
	dir := filepath.Join(RepositoryDir, ExecutionsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Execution{}, nil
		}
		return nil, fmt.Errorf("reading executions directory: %w", err)
	}

	var list []Execution
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		exec, err := LoadExecution(id)
		if err == nil {
			list = append(list, exec)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Timestamp > list[j].Timestamp
	})

	return list, nil
}

// RecordExecution creates and persists an Execution record automatically bound to current HEAD commit ID.
func RecordExecution(inputs, outputs string, durationMs int64, tokens TokenUsage, meta map[string]string) (Execution, error) {
	repo, err := OpenRepository()
	if err != nil {
		return Execution{}, fmt.Errorf("opening repository: %w", err)
	}

	if repo.Branch.Head == "" {
		return Execution{}, fmt.Errorf("cannot record execution: repository has no commits")
	}

	exec := NewExecution(repo.Branch.Head, inputs, outputs, durationMs, tokens, meta)
	if err := SaveExecution(exec); err != nil {
		return exec, err
	}

	return exec, nil
}
