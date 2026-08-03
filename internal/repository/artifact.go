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
