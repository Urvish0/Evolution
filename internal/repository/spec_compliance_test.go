package repository

import (
	"encoding/json"
	"os"
	"testing"
)

// TestSpecComplianceRequiredFields validates that the manifest enforces v0.2 required fields.
func TestSpecComplianceRequiredFields(t *testing.T) {
	tests := []struct {
		name      string
		manifest  Manifest
		expectErr bool
	}{
		{
			name:      "valid manifest with version and name",
			manifest:  Manifest{Version: "0.2.0", Name: "test-ai"},
			expectErr: false,
		},
		{
			name:      "missing version",
			manifest:  Manifest{Version: "", Name: "test-ai"},
			expectErr: true,
		},
		{
			name:      "missing name",
			manifest:  Manifest{Version: "0.2.0", Name: ""},
			expectErr: true,
		},
		{
			name:      "both missing",
			manifest:  Manifest{Version: "", Name: ""},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateManifest(tt.manifest)
			if tt.expectErr && err == nil {
				t.Errorf("expected validation error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

// TestSpecComplianceModelConfigRequired validates model_config.model is required when model_config is present.
func TestSpecComplianceModelConfigRequired(t *testing.T) {
	m := Manifest{
		Version: "0.2.0",
		Name:    "test-ai",
		Artifacts: ManifestArtifacts{
			ModelConfig: &ModelConfigArtifact{
				BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypeModelConfig, Name: "mc", Path: "config/m.json"},
				Model:        "", // Missing required field
			},
		},
	}

	if err := ValidateManifest(m); err == nil {
		t.Errorf("expected validation error for empty model_config.model, got nil")
	}
}

// TestSpecComplianceAutoHash validates that hash is auto-computed when empty and file exists.
func TestSpecComplianceAutoHash(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create a real file for the prompt artifact
	_ = os.MkdirAll("prompts", 0755)
	_ = os.WriteFile("prompts/system.txt", []byte("You are a legal research AI."), 0644)

	m := NewDefaultManifest("hash-test", "Testing auto-hash computation")
	m.Artifacts.Prompts = []PromptArtifact{
		{
			BaseArtifact: BaseArtifact{
				ArtifactType: ArtifactTypePrompt,
				Name:         "sys-prompt",
				Hash:         "", // Empty — should be auto-computed
				Path:         "prompts/system.txt",
			},
			Role:   "system",
			Format: "text",
		},
	}

	artMap := m.ToArtifactMap()

	prompts, ok := artMap["prompts"]
	if !ok || len(prompts) == 0 {
		t.Fatalf("expected prompts in artifact map")
	}

	hash := prompts[0].GetHash()
	if hash == "" {
		t.Errorf("expected auto-computed hash, got empty string")
	}

	// Verify hash is a valid 64-char SHA-256 hex string
	if len(hash) != 64 {
		t.Errorf("expected 64-char SHA-256 hash, got %d chars: %s", len(hash), hash)
	}
}

// TestSpecComplianceJSONSchemaStructure validates manifest serializes to JSON matching the schema shape.
func TestSpecComplianceJSONSchemaStructure(t *testing.T) {
	m := NewDefaultManifest("schema-test", "Testing JSON Schema compliance")
	m.Artifacts.Prompts = []PromptArtifact{
		{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypePrompt, Name: "p1", Hash: "", Path: "prompts/p1.txt"},
			Role:         "system",
			Format:       "text",
		},
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal manifest: %v", err)
	}

	// Unmarshal back into a generic map to check shape
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal manifest JSON: %v", err)
	}

	// Check required root fields
	if _, ok := raw["version"]; !ok {
		t.Errorf("missing 'version' field in JSON output")
	}
	if _, ok := raw["name"]; !ok {
		t.Errorf("missing 'name' field in JSON output")
	}

	// Check artifacts is an object with prompts array
	artifacts, ok := raw["artifacts"].(map[string]any)
	if !ok {
		t.Fatalf("'artifacts' is not an object")
	}

	prompts, ok := artifacts["prompts"].([]any)
	if !ok {
		t.Fatalf("'artifacts.prompts' is not an array")
	}

	if len(prompts) != 1 {
		t.Errorf("expected 1 prompt, got %d", len(prompts))
	}

	// Check prompt has required artifact fields
	prompt, ok := prompts[0].(map[string]any)
	if !ok {
		t.Fatalf("prompt is not an object")
	}

	for _, field := range []string{"type", "name", "path"} {
		if _, ok := prompt[field]; !ok {
			t.Errorf("prompt missing required field %q", field)
		}
	}
}
