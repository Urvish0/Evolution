package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ArtifactType constants matching the Intelligence Manifest Specification v0.1.
const (
	ArtifactTypePrompt      = "prompt"
	ArtifactTypeMemory      = "memory"
	ArtifactTypeRetrieval   = "retrieval"
	ArtifactTypeTool        = "tool"
	ArtifactTypeModelConfig = "model_config"
	ArtifactTypePolicy      = "policy"
)

// Artifact is the fundamental interface every AI intelligence component implements.
type Artifact interface {
	Type() string
	GetName() string
	GetHash() string
	GetPath() string
	Serialize() ([]byte, error)
}

// BaseArtifact contains common fields shared by all artifact types.
type BaseArtifact struct {
	ArtifactType string `json:"type"`
	Name         string `json:"name"`
	Hash         string `json:"hash"`
	Path         string `json:"path"`
	Description  string `json:"description,omitempty"`
}

func (b BaseArtifact) Type() string    { return b.ArtifactType }
func (b BaseArtifact) GetName() string { return b.Name }
func (b BaseArtifact) GetHash() string { return b.Hash }
func (b BaseArtifact) GetPath() string { return b.Path }
func (b BaseArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

// --- 1. Prompt Artifact ---

type PromptArtifact struct {
	BaseArtifact
	Role   string `json:"role"`   // system, user, assistant, few_shot
	Format string `json:"format"` // text, template, jinja2, mustache
}

func (p PromptArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// --- 2. Memory Artifact ---

type MemoryArtifact struct {
	BaseArtifact
	Strategy  string `json:"strategy"`   // buffer_window, summary, vector, graph
	MaxTokens int    `json:"max_tokens"` // token budget limit
}

func (m MemoryArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// --- 3. Retrieval Artifact ---

type RetrievalArtifact struct {
	BaseArtifact
	Source    string `json:"source"`     // pinecone, chroma, weaviate, local, elasticsearch
	ChunkSize int    `json:"chunk_size"` // token chunk size
	TopK      int    `json:"top_k"`      // number of search results
}

func (r RetrievalArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// --- 4. Tool Artifact ---

type ToolArtifact struct {
	BaseArtifact
	Provider     string `json:"provider"`      // e.g. google, custom, bash
	AuthRequired bool   `json:"auth_required"` // whether authentication is required
}

func (t ToolArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

// --- 5. ModelConfig Artifact ---

type ModelConfigArtifact struct {
	BaseArtifact
	Model       string  `json:"model"`       // e.g. gpt-4o, claude-3.5-sonnet
	Provider    string  `json:"provider"`    // openai, anthropic, google, local
	Temperature float64 `json:"temperature"` // sampling temperature
	MaxTokens   int     `json:"max_tokens"`  // max output tokens
	TopP        float64 `json:"top_p"`       // nucleus sampling
}

func (mc ModelConfigArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(mc, "", "  ")
}

// --- 6. Policy Artifact ---

type PolicyArtifact struct {
	BaseArtifact
	Enforcement string `json:"enforcement"` // strict, warn, log
}

func (p PolicyArtifact) Serialize() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// --- Storage Helper Functions ---

// GetArtifactPath returns disk location .evolution/artifacts/<type>/<hash>.json
func GetArtifactPath(artType, hash string) string {
	return filepath.Join(RepositoryDir, ArtifactsDir, artType, hash+".json")
}

// SaveArtifact writes the serialized artifact JSON under .evolution/artifacts/<type>/<hash>.json
func SaveArtifact(a Artifact) error {
	data, err := a.Serialize()
	if err != nil {
		return fmt.Errorf("serializing artifact %s: %w", a.GetName(), err)
	}

	// Compute hash if not set
	hash := a.GetHash()
	if hash == "" {
		hash = HashRaw(data)
	}

	artDir := filepath.Join(RepositoryDir, ArtifactsDir, a.Type())
	if err := os.MkdirAll(artDir, 0755); err != nil {
		return fmt.Errorf("creating artifact directory %s: %w", artDir, err)
	}

	filePath := GetArtifactPath(a.Type(), hash)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("writing artifact file %s: %w", filePath, err)
	}

	return nil
}

// LoadArtifact reads an artifact from disk and unmarshals it into the specific typed struct.
func LoadArtifact(artType, hash string) (Artifact, error) {
	filePath := GetArtifactPath(artType, hash)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading artifact %s/%s: %w", artType, hash[:8], err)
	}

	switch artType {
	case ArtifactTypePrompt:
		var a PromptArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	case ArtifactTypeMemory:
		var a MemoryArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	case ArtifactTypeRetrieval:
		var a RetrievalArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	case ArtifactTypeTool:
		var a ToolArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	case ArtifactTypeModelConfig:
		var a ModelConfigArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	case ArtifactTypePolicy:
		var a PolicyArtifact
		err = json.Unmarshal(data, &a)
		return a, err
	default:
		return nil, fmt.Errorf("unknown artifact type %q", artType)
	}
}

// SemanticDiffChange describes a high-level AI component difference between two commits.
type SemanticDiffChange struct {
	Category     string `json:"category"`
	ArtifactName string `json:"artifact_name"`
	Action       string `json:"action"` // "added", "removed", "modified"
	Details      string `json:"details"`
}

// RegisterArtifactInManifest adds or updates an artifact entry inside evolution.manifest.json.
func RegisterArtifactInManifest(artType, name, path string) error {
	var m Manifest
	if ManifestExists() {
		var err error
		m, err = LoadManifest(ManifestFileName)
		if err != nil {
			return err
		}
	} else {
		m = NewDefaultManifest("ai-intelligence", "AI system state")
	}

	base := BaseArtifact{
		ArtifactType: artType,
		Name:         name,
		Path:         path,
	}

	switch artType {
	case ArtifactTypePrompt:
		// Check if entry already exists, update or append
		found := false
		for i, p := range m.Artifacts.Prompts {
			if p.Name == name {
				m.Artifacts.Prompts[i].Path = path
				found = true
				break
			}
		}
		if !found {
			m.Artifacts.Prompts = append(m.Artifacts.Prompts, PromptArtifact{
				BaseArtifact: base,
				Role:         "system",
				Format:       "text",
			})
		}

	case ArtifactTypeModelConfig:
		m.Artifacts.ModelConfig = &ModelConfigArtifact{
			BaseArtifact: base,
			Model:        "gpt-4o",
			Provider:     "openai",
			Temperature:  0.7,
			MaxTokens:    4096,
		}

	case ArtifactTypeTool:
		found := false
		for i, t := range m.Artifacts.Tools {
			if t.Name == name {
				m.Artifacts.Tools[i].Path = path
				found = true
				break
			}
		}
		if !found {
			m.Artifacts.Tools = append(m.Artifacts.Tools, ToolArtifact{
				BaseArtifact: base,
				Provider:     "custom",
				AuthRequired: false,
			})
		}

	case ArtifactTypeMemory:
		found := false
		for i, mem := range m.Artifacts.Memory {
			if mem.Name == name {
				m.Artifacts.Memory[i].Path = path
				found = true
				break
			}
		}
		if !found {
			m.Artifacts.Memory = append(m.Artifacts.Memory, MemoryArtifact{
				BaseArtifact: base,
				Strategy:     "buffer_window",
				MaxTokens:    4096,
			})
		}

	case ArtifactTypeRetrieval:
		found := false
		for i, r := range m.Artifacts.Retrieval {
			if r.Name == name {
				m.Artifacts.Retrieval[i].Path = path
				found = true
				break
			}
		}
		if !found {
			m.Artifacts.Retrieval = append(m.Artifacts.Retrieval, RetrievalArtifact{
				BaseArtifact: base,
				Source:       "local",
				ChunkSize:    512,
				TopK:         5,
			})
		}

	case ArtifactTypePolicy:
		found := false
		for i, pol := range m.Artifacts.Policies {
			if pol.Name == name {
				m.Artifacts.Policies[i].Path = path
				found = true
				break
			}
		}
		if !found {
			m.Artifacts.Policies = append(m.Artifacts.Policies, PolicyArtifact{
				BaseArtifact: base,
				Enforcement:  "strict",
			})
		}

	default:
		return fmt.Errorf("invalid artifact type %q", artType)
	}

	return SaveManifest(m, ManifestFileName)
}

// GetHeadCommitArtifacts returns the artifacts map attached to the current HEAD commit.
func GetHeadCommitArtifacts() (map[string][]Artifact, error) {
	repo, err := OpenRepository()
	if err != nil {
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	if repo.Branch.Head == "" {
		return nil, fmt.Errorf("no commits in repository")
	}

	commit, err := LoadCommit(repo.Branch.Head)
	if err != nil {
		return nil, fmt.Errorf("loading HEAD commit: %w", err)
	}

	return commit.Artifacts, nil
}

// CompareCommitArtifacts performs a high-level semantic diff between two commits.
func CompareCommitArtifacts(commitID1, commitID2 string) ([]SemanticDiffChange, error) {
	c1, err := LoadCommit(commitID1)
	if err != nil {
		return nil, fmt.Errorf("loading commit %s: %w", commitID1[:8], err)
	}

	c2, err := LoadCommit(commitID2)
	if err != nil {
		return nil, fmt.Errorf("loading commit %s: %w", commitID2[:8], err)
	}

	var changes []SemanticDiffChange

	allCategories := map[string]bool{
		"prompts": true, "model_config": true, "tools": true,
		"memory": true, "retrieval": true, "policies": true,
	}

	for category := range allCategories {
		list1 := c1.Artifacts[category]
		list2 := c2.Artifacts[category]

		map1 := make(map[string]Artifact)
		for _, a := range list1 {
			map1[a.GetName()] = a
		}

		map2 := make(map[string]Artifact)
		for _, a := range list2 {
			map2[a.GetName()] = a
		}

		// Check for added or modified artifacts in commit 2
		for name, a2 := range map2 {
			a1, exists := map1[name]
			if !exists {
				changes = append(changes, SemanticDiffChange{
					Category:     category,
					ArtifactName: name,
					Action:       "added",
					Details:      fmt.Sprintf("New %s registered at path '%s'", category, a2.GetPath()),
				})
			} else {
				data1, _ := a1.Serialize()
				data2, _ := a2.Serialize()
				if string(data1) != string(data2) {
					changes = append(changes, SemanticDiffChange{
						Category:     category,
						ArtifactName: name,
						Action:       "modified",
						Details:      fmt.Sprintf("Configuration changed for '%s'", name),
					})
				}
			}
		}

		// Check for removed artifacts in commit 2
		for name := range map1 {
			if _, exists := map2[name]; !exists {
				changes = append(changes, SemanticDiffChange{
					Category:     category,
					ArtifactName: name,
					Action:       "removed",
					Details:      fmt.Sprintf("%s '%s' was removed", category, name),
				})
			}
		}
	}

	return changes, nil
}
