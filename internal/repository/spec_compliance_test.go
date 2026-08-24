package repository

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

// TestSpecComplianceRequiredFields validates that the manifest enforces v1.0 required fields.
func TestSpecComplianceRequiredFields(t *testing.T) {
	tests := []struct {
		name      string
		manifest  Manifest
		expectErr bool
	}{
		{
			name:      "valid manifest with version and name",
			manifest:  Manifest{Version: "1.0.0", Name: "test-ai"},
			expectErr: false,
		},
		{
			name:      "missing version",
			manifest:  Manifest{Version: "", Name: "test-ai"},
			expectErr: true,
		},
		{
			name:      "missing name",
			manifest:  Manifest{Version: "1.0.0", Name: ""},
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
		Version: "1.0.0",
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

// TestSpecComplianceExecutionSchema validates that the Execution struct serializes
// to JSON matching the v1.0 execution schema definition.
func TestSpecComplianceExecutionSchema(t *testing.T) {
	exec := NewExecution(
		"commit-123",
		"What is contract law?",
		"Contract law governs agreements.",
		250,
		TokenUsage{PromptTokens: 40, CompletionTokens: 25, TotalTokens: 65},
		map[string]string{"model": "gpt-4o"},
	)

	data, err := json.MarshalIndent(exec, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal Execution: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal Execution JSON: %v", err)
	}

	// Validate all required fields per v1.0 spec Section 4.2
	requiredFields := []string{"id", "commit_id", "inputs", "outputs", "duration_ms", "tokens", "timestamp"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("execution missing required field %q", field)
		}
	}

	// Validate tokens sub-object required fields
	tokens, ok := raw["tokens"].(map[string]any)
	if !ok {
		t.Fatalf("'tokens' is not an object")
	}
	for _, field := range []string{"prompt_tokens", "completion_tokens", "total_tokens"} {
		if _, ok := tokens[field]; !ok {
			t.Errorf("tokens missing required field %q", field)
		}
	}

	// Validate timestamp is RFC 3339
	ts, ok := raw["timestamp"].(string)
	if !ok {
		t.Fatalf("'timestamp' is not a string")
	}
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Errorf("timestamp %q is not RFC 3339: %v", ts, err)
	}
}

// TestSpecComplianceEvaluationSchema validates that the EvaluationResult struct serializes
// to JSON matching the v1.0 evaluation schema definition.
func TestSpecComplianceEvaluationSchema(t *testing.T) {
	result := EvaluationResult{
		ID:           "eval-001",
		CommitID:     "commit-123",
		ExecutionID:  "exec-456",
		OverallScore: 0.92,
		Scores: map[string]EvaluationScore{
			"performance": {Name: "performance", Score: 1.0, Unit: "latency_score", Details: "Fast"},
			"safety":      {Name: "safety", Score: 0.85, Unit: "pass_rate", Details: "Minor flag"},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal EvaluationResult: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal EvaluationResult JSON: %v", err)
	}

	// Validate all required fields per v1.0 spec Section 5.2
	requiredFields := []string{"id", "commit_id", "execution_id", "overall_score", "scores", "timestamp"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("evaluation missing required field %q", field)
		}
	}

	// Validate scores is a map of evaluator -> score object
	scores, ok := raw["scores"].(map[string]any)
	if !ok {
		t.Fatalf("'scores' is not an object")
	}

	for evalName, scoreRaw := range scores {
		score, ok := scoreRaw.(map[string]any)
		if !ok {
			t.Errorf("score for %q is not an object", evalName)
			continue
		}
		// name and score are required per spec Section 5.3
		for _, field := range []string{"name", "score"} {
			if _, ok := score[field]; !ok {
				t.Errorf("score %q missing required field %q", evalName, field)
			}
		}
	}

	// Validate overall_score range
	os, ok := raw["overall_score"].(float64)
	if !ok {
		t.Fatalf("'overall_score' is not a number")
	}
	if os < 0.0 || os > 1.0 {
		t.Errorf("overall_score %f out of [0.0, 1.0] range", os)
	}
}

// TestSpecComplianceRegressionRules validates the regression rule structure
// matches the v1.0 spec Section 6 definition.
func TestSpecComplianceRegressionRules(t *testing.T) {
	// Test that RegressionRule serializes correctly
	rule := RegressionRule{
		MinScore:      0.80,
		MaxDrop:       0.05,
		RequireSafety: true,
	}

	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal RegressionRule: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal RegressionRule JSON: %v", err)
	}

	// Validate fields per spec Section 6.1
	if ms, ok := raw["min_score"].(float64); !ok || ms != 0.80 {
		t.Errorf("expected min_score 0.80, got %v", raw["min_score"])
	}
	if md, ok := raw["max_drop"].(float64); !ok || md != 0.05 {
		t.Errorf("expected max_drop 0.05, got %v", raw["max_drop"])
	}
	if rs, ok := raw["require_safety"].(bool); !ok || !rs {
		t.Errorf("expected require_safety true, got %v", raw["require_safety"])
	}
}
