package repository

import (
	"bytes"
	"testing"
)

func TestArtifactSaveAndLoad(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 1. Create a Prompt Artifact
	prompt := PromptArtifact{
		BaseArtifact: BaseArtifact{
			ArtifactType: ArtifactTypePrompt,
			Name:         "system-prompt",
			Hash:         "b671051cadc3b0806ab998726cb95055a713983b1e94605e9bc3d88689d21a0e",
			Path:         "prompts/system.txt",
			Description:  "Main system persona prompt",
		},
		Role:   "system",
		Format: "text",
	}

	// 2. Save Artifact to disk
	if err := SaveArtifact(prompt); err != nil {
		t.Fatalf("SaveArtifact() failed: %v", err)
	}

	// 3. Load Artifact back
	loaded, err := LoadArtifact(ArtifactTypePrompt, prompt.Hash)
	if err != nil {
		t.Fatalf("LoadArtifact() failed: %v", err)
	}

	if loaded.Type() != ArtifactTypePrompt {
		t.Errorf("expected type %q, got %q", ArtifactTypePrompt, loaded.Type())
	}
	if loaded.GetName() != "system-prompt" {
		t.Errorf("expected name %q, got %q", "system-prompt", loaded.GetName())
	}

	// Type assertion test (Polymorphism in Go)
	promptLoaded, ok := loaded.(PromptArtifact)
	if !ok {
		t.Fatalf("type assertion to PromptArtifact failed")
	}

	if promptLoaded.Role != "system" {
		t.Errorf("expected role 'system', got %q", promptLoaded.Role)
	}
}

func TestPolymorphicArtifactSlice(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Polymorphic slice holding different concrete types implementing Artifact interface
	artifacts := []Artifact{
		PromptArtifact{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypePrompt, Name: "p1", Hash: "1111111111111111111111111111111111111111111111111111111111111111", Path: "p1.txt"},
			Role:         "system",
			Format:       "text",
		},
		ModelConfigArtifact{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypeModelConfig, Name: "m1", Hash: "2222222222222222222222222222222222222222222222222222222222222222", Path: "m1.json"},
			Model:        "gpt-4o",
			Provider:     "openai",
			Temperature:  0.7,
		},
		PolicyArtifact{
			BaseArtifact: BaseArtifact{ArtifactType: ArtifactTypePolicy, Name: "pol1", Hash: "3333333333333333333333333333333333333333333333333333333333333333", Path: "pol1.json"},
			Enforcement:  "strict",
		},
	}

	for _, a := range artifacts {
		if err := SaveArtifact(a); err != nil {
			t.Fatalf("SaveArtifact(%s) failed: %v", a.GetName(), err)
		}
	}

	for _, a := range artifacts {
		loaded, err := LoadArtifact(a.Type(), a.GetHash())
		if err != nil {
			t.Fatalf("LoadArtifact(%s, %s) failed: %v", a.Type(), a.GetHash()[:8], err)
		}

		if loaded.GetName() != a.GetName() {
			t.Errorf("expected name %s, got %s", a.GetName(), loaded.GetName())
		}

		serializedOrig, _ := a.Serialize()
		serializedLoaded, _ := loaded.Serialize()

		if !bytes.Equal(serializedOrig, serializedLoaded) {
			t.Errorf("mismatched serialization for %s", a.GetName())
		}
	}
}
