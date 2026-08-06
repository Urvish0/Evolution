package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// CommitMetadata captures execution context and environment metadata.
type CommitMetadata struct {
	Environment map[string]string `json:"environment,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

// Commit represents an immutable Intelligence Commit snapshot.
type Commit struct {
	ID        string                `json:"id"`
	Parent    string                `json:"parent"`
	Message   string                `json:"message"`
	Author    string                `json:"author"`
	Timestamp string                `json:"timestamp"`
	TreeHash  string                `json:"tree_hash"`
	Artifacts map[string][]Artifact `json:"artifacts,omitempty"`
	Metadata  CommitMetadata        `json:"metadata,omitempty"`
}

// NewCommit initializes a new empty Commit struct.
func NewCommit(message string) Commit {
	return Commit{
		ID:        uuid.New().String(),
		Parent:    "",
		Message:   message,
		Author:    "",
		Timestamp: time.Now().Format(time.RFC3339),
		TreeHash:  "",
		Artifacts: make(map[string][]Artifact),
		Metadata: CommitMetadata{
			Environment: map[string]string{
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
				"go_version": runtime.Version(),
			},
			Tags: []string{},
		},
	}
}

// Intermediate DTO for JSON serialization of interface types
type commitDTO struct {
	ID        string                     `json:"id"`
	Parent    string                     `json:"parent"`
	Message   string                     `json:"message"`
	Author    string                     `json:"author"`
	Timestamp string                     `json:"timestamp"`
	TreeHash  string                     `json:"tree_hash"`
	Artifacts map[string][]json.RawMessage `json:"artifacts,omitempty"`
	Metadata  CommitMetadata             `json:"metadata,omitempty"`
}

// Custom MarshalJSON to handle interface slice serialization in Commit.
func (c Commit) MarshalJSON() ([]byte, error) {
	dto := commitDTO{
		ID:        c.ID,
		Parent:    c.Parent,
		Message:   c.Message,
		Author:    c.Author,
		Timestamp: c.Timestamp,
		TreeHash:  c.TreeHash,
		Artifacts: make(map[string][]json.RawMessage),
		Metadata:  c.Metadata,
	}

	for category, list := range c.Artifacts {
		for _, art := range list {
			bytes, err := art.Serialize()
			if err != nil {
				return nil, fmt.Errorf("marshaling artifact %s: %w", art.GetName(), err)
			}
			dto.Artifacts[category] = append(dto.Artifacts[category], json.RawMessage(bytes))
		}
	}

	return json.MarshalIndent(dto, "", "  ")
}

// Custom UnmarshalJSON to dynamically deserialize concrete Artifact structs.
func (c *Commit) UnmarshalJSON(data []byte) error {
	var dto commitDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}

	c.ID = dto.ID
	c.Parent = dto.Parent
	c.Message = dto.Message
	c.Author = dto.Author
	c.Timestamp = dto.Timestamp
	c.TreeHash = dto.TreeHash
	c.Metadata = dto.Metadata
	c.Artifacts = make(map[string][]Artifact)

	for category, rawList := range dto.Artifacts {
		for _, raw := range rawList {
			// Inspect 'type' field from raw JSON
			var base BaseArtifact
			if err := json.Unmarshal(raw, &base); err != nil {
				return fmt.Errorf("unmarshaling artifact base: %w", err)
			}

			var art Artifact
			switch base.ArtifactType {
			case ArtifactTypePrompt:
				var a PromptArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			case ArtifactTypeMemory:
				var a MemoryArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			case ArtifactTypeRetrieval:
				var a RetrievalArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			case ArtifactTypeTool:
				var a ToolArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			case ArtifactTypeModelConfig:
				var a ModelConfigArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			case ArtifactTypePolicy:
				var a PolicyArtifact
				if err := json.Unmarshal(raw, &a); err != nil {
					return err
				}
				art = a
			default:
				art = base
			}

			c.Artifacts[category] = append(c.Artifacts[category], art)
		}
	}

	return nil
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
	data, err := json.MarshalIndent(commit, "", "  ")
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

	// Capture workspace snapshot Merkle Tree hash
	var treeHash string
	idx, err := LoadIndex()
	if err == nil && len(idx.Entries) > 0 {
		treeHash, err = BuildTreeFromIndex(idx)
		if err != nil {
			return fmt.Errorf("building tree from index: %w", err)
		}
	} else {
		// Fallback: build tree from directory if index is empty
		treeHash, err = BuildTreeFromDirectory(".")
		if err != nil {
			return fmt.Errorf("creating workspace snapshot tree: %w", err)
		}
	}
	commit.TreeHash = treeHash

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

	// Clear index after successful commit if it was used
	if len(idx.Entries) > 0 {
		_ = SaveIndex(NewIndex())
	}

	return nil
}
