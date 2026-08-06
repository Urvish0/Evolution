package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ManifestFileName = "evolution.manifest.json"
	SpecVersionV01   = "0.1.0"
)

// ManifestArtifacts holds arrays of typed artifacts matching the Spec v0.1 schema.
type ManifestArtifacts struct {
	Prompts     []PromptArtifact     `json:"prompts,omitempty"`
	Memory      []MemoryArtifact     `json:"memory,omitempty"`
	Retrieval   []RetrievalArtifact  `json:"retrieval,omitempty"`
	Tools       []ToolArtifact       `json:"tools,omitempty"`
	ModelConfig *ModelConfigArtifact `json:"model_config,omitempty"`
	Policies    []PolicyArtifact     `json:"policies,omitempty"`
}

// Manifest represents the root evolution.manifest.json structure.
type Manifest struct {
	Version     string            `json:"version"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Artifacts   ManifestArtifacts `json:"artifacts,omitempty"`
	Metadata    map[string]any    `json:"metadata,omitempty"`
}

// NewDefaultManifest initializes a starter Manifest struct conforming to Spec v0.1.
func NewDefaultManifest(name, description string) Manifest {
	if name == "" {
		name = "ai-intelligence"
	}
	if description == "" {
		description = "AI system powered by Evolution version control"
	}

	return Manifest{
		Version:     SpecVersionV01,
		Name:        name,
		Description: description,
		Artifacts: ManifestArtifacts{
			Prompts:   []PromptArtifact{},
			Memory:    []MemoryArtifact{},
			Retrieval: []RetrievalArtifact{},
			Tools:     []ToolArtifact{},
			ModelConfig: &ModelConfigArtifact{
				BaseArtifact: BaseArtifact{
					ArtifactType: ArtifactTypeModelConfig,
					Name:         "primary-model",
					Path:         "config/model.json",
				},
				Model:       "gpt-4o",
				Provider:    "openai",
				Temperature: 0.7,
				MaxTokens:   4096,
			},
			Policies: []PolicyArtifact{},
		},
		Metadata: make(map[string]any),
	}
}

// InitManifest creates a starter evolution.manifest.json in the current working directory.
func InitManifest(name, description string) error {
	if ManifestExists() {
		return fmt.Errorf("manifest file %s already exists", ManifestFileName)
	}

	m := NewDefaultManifest(name, description)
	return SaveManifest(m, ManifestFileName)
}

// ManifestExists returns true if evolution.manifest.json exists in the specified directory.
func ManifestExists() bool {
	_, err := os.Stat(ManifestFileName)
	return err == nil
}

// LoadManifest reads and parses a manifest file from path.
func LoadManifest(path string) (Manifest, error) {
	var m Manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return m, fmt.Errorf("reading manifest file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("parsing manifest JSON: %w", err)
	}

	return m, nil
}

// SaveManifest writes a Manifest struct to disk in pretty-printed JSON.
func SaveManifest(m Manifest, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing manifest file %s: %w", path, err)
	}

	return nil
}

// ValidateManifest checks if a Manifest conforms to Intelligence Manifest Specification v0.1.
func ValidateManifest(m Manifest) error {
	if m.Version == "" {
		return fmt.Errorf("validation error: missing required field 'version'")
	}
	if m.Name == "" {
		return fmt.Errorf("validation error: missing required field 'name'")
	}

	// Validate ModelConfig if present
	if m.Artifacts.ModelConfig != nil {
		mc := m.Artifacts.ModelConfig
		if mc.Model == "" {
			return fmt.Errorf("validation error: model_config missing required 'model' field")
		}
	}

	return nil
}

// ToArtifactMap converts a Manifest's typed artifacts into a map[string][]Artifact for Commit.
func (m Manifest) ToArtifactMap() map[string][]Artifact {
	artMap := make(map[string][]Artifact)

	if len(m.Artifacts.Prompts) > 0 {
		var list []Artifact
		for _, p := range m.Artifacts.Prompts {
			list = append(list, p)
		}
		artMap["prompts"] = list
	}

	if len(m.Artifacts.Memory) > 0 {
		var list []Artifact
		for _, mem := range m.Artifacts.Memory {
			list = append(list, mem)
		}
		artMap["memory"] = list
	}

	if len(m.Artifacts.Retrieval) > 0 {
		var list []Artifact
		for _, r := range m.Artifacts.Retrieval {
			list = append(list, r)
		}
		artMap["retrieval"] = list
	}

	if len(m.Artifacts.Tools) > 0 {
		var list []Artifact
		for _, t := range m.Artifacts.Tools {
			list = append(list, t)
		}
		artMap["tools"] = list
	}

	if m.Artifacts.ModelConfig != nil {
		artMap["model_config"] = []Artifact{*m.Artifacts.ModelConfig}
	}

	if len(m.Artifacts.Policies) > 0 {
		var list []Artifact
		for _, pol := range m.Artifacts.Policies {
			list = append(list, pol)
		}
		artMap["policies"] = list
	}

	return artMap
}
