package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReconstructedState holds the complete reconstructed operational state of an AI system from an Intelligence Commit.
type ReconstructedState struct {
	CommitID    string            `json:"commit_id"`
	CommitMsg   string            `json:"commit_message"`
	Author      string            `json:"author"`
	Timestamp   string            `json:"timestamp"`
	Manifest    Manifest          `json:"manifest"`
	Environment map[string]string `json:"environment,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

// ReconstructCommitState reconstructs the full operational state (all 6 artifact types + manifest) of a target commit.
func ReconstructCommitState(target string) (*ReconstructedState, error) {
	commit, err := ResolveRevisionToCommit(target)
	if err != nil {
		return nil, fmt.Errorf("resolving commit %s: %w", target, err)
	}

	manifest := NewDefaultManifest("reconstructed-intelligence", fmt.Sprintf("Reconstructed state from commit %s", commit.ID[:8]))
	manifest.Artifacts = ManifestArtifacts{}

	// Hydrate artifacts map attached to commit
	if commit.Artifacts != nil {
		for category, artList := range commit.Artifacts {
			for _, art := range artList {
				switch category {
				case "prompts", ArtifactTypePrompt:
					if p, ok := art.(PromptArtifact); ok {
						manifest.Artifacts.Prompts = append(manifest.Artifacts.Prompts, p)
					} else if p, ok := art.(*PromptArtifact); ok {
						manifest.Artifacts.Prompts = append(manifest.Artifacts.Prompts, *p)
					}
				case ArtifactTypeMemory:
					if m, ok := art.(MemoryArtifact); ok {
						manifest.Artifacts.Memory = append(manifest.Artifacts.Memory, m)
					} else if m, ok := art.(*MemoryArtifact); ok {
						manifest.Artifacts.Memory = append(manifest.Artifacts.Memory, *m)
					}
				case ArtifactTypeRetrieval:
					if r, ok := art.(RetrievalArtifact); ok {
						manifest.Artifacts.Retrieval = append(manifest.Artifacts.Retrieval, r)
					} else if r, ok := art.(*RetrievalArtifact); ok {
						manifest.Artifacts.Retrieval = append(manifest.Artifacts.Retrieval, *r)
					}
				case "tools", ArtifactTypeTool:
					if t, ok := art.(ToolArtifact); ok {
						manifest.Artifacts.Tools = append(manifest.Artifacts.Tools, t)
					} else if t, ok := art.(*ToolArtifact); ok {
						manifest.Artifacts.Tools = append(manifest.Artifacts.Tools, *t)
					}
				case ArtifactTypeModelConfig:
					if mc, ok := art.(ModelConfigArtifact); ok {
						manifest.Artifacts.ModelConfig = &mc
					} else if mc, ok := art.(*ModelConfigArtifact); ok {
						manifest.Artifacts.ModelConfig = mc
					}
				case "policies", ArtifactTypePolicy:
					if pol, ok := art.(PolicyArtifact); ok {
						manifest.Artifacts.Policies = append(manifest.Artifacts.Policies, pol)
					} else if pol, ok := art.(*PolicyArtifact); ok {
						manifest.Artifacts.Policies = append(manifest.Artifacts.Policies, *pol)
					}
				}
			}
		}
	}

	state := &ReconstructedState{
		CommitID:    commit.ID,
		CommitMsg:   commit.Message,
		Author:      commit.Author,
		Timestamp:   commit.Timestamp,
		Manifest:    manifest,
		Environment: commit.Metadata.Environment,
		Tags:        commit.Metadata.Tags,
	}

	return state, nil
}

// ExportReconstructedState exports the reconstructed Manifest to a target file path as pretty-printed JSON.
func ExportReconstructedState(state *ReconstructedState, targetPath string) error {
	data, err := json.MarshalIndent(state.Manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	return os.WriteFile(targetPath, data, 0644)
}

// ReplayExecution loads a recorded execution and reconstructs the Intelligence Commit state active during that run.
func ReplayExecution(executionID string) (*ReconstructedState, Execution, error) {
	exec, err := LoadExecution(executionID)
	if err != nil {
		return nil, Execution{}, fmt.Errorf("loading execution %s: %w", executionID, err)
	}

	state, err := ReconstructCommitState(exec.CommitID)
	if err != nil {
		return nil, exec, fmt.Errorf("reconstructing commit state for execution: %w", err)
	}

	return state, exec, nil
}

// ExecutionComparison holds metrics and text differences between two execution runs.
type ExecutionComparison struct {
	Exec1            Execution  `json:"execution_1"`
	Exec2            Execution  `json:"execution_2"`
	PromptTokenDelta int        `json:"prompt_token_delta"`
	CompTokenDelta   int        `json:"comp_token_delta"`
	TotalTokenDelta  int        `json:"total_token_delta"`
	DurationDeltaMs  int64      `json:"duration_delta_ms"`
	OutputDiffLines  []DiffLine `json:"output_diff_lines"`
}

// CompareExecutions compares two recorded execution runs side-by-side.
func CompareExecutions(id1, id2 string) (*ExecutionComparison, error) {
	e1, err := LoadExecution(id1)
	if err != nil {
		return nil, fmt.Errorf("loading execution 1 (%s): %w", id1, err)
	}

	e2, err := LoadExecution(id2)
	if err != nil {
		return nil, fmt.Errorf("loading execution 2 (%s): %w", id2, err)
	}

	diffLines := ComputeLineDiff(e1.Outputs, e2.Outputs)

	comp := &ExecutionComparison{
		Exec1:            e1,
		Exec2:            e2,
		PromptTokenDelta: e2.Tokens.PromptTokens - e1.Tokens.PromptTokens,
		CompTokenDelta:   e2.Tokens.CompletionTokens - e1.Tokens.CompletionTokens,
		TotalTokenDelta:  e2.Tokens.TotalTokens - e1.Tokens.TotalTokens,
		DurationDeltaMs:  e2.DurationMs - e1.DurationMs,
		OutputDiffLines:  diffLines,
	}

	return comp, nil
}

// CompareCommitReplays reconstructs both commit states and compares their intelligence artifacts.
func CompareCommitReplays(target1, target2 string) (*ReconstructedState, *ReconstructedState, []SemanticDiffChange, error) {
	s1, err := ReconstructCommitState(target1)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reconstructing target 1 (%s): %w", target1, err)
	}

	s2, err := ReconstructCommitState(target2)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reconstructing target 2 (%s): %w", target2, err)
	}

	changes, err := CompareCommitArtifacts(s1.CommitID, s2.CommitID)
	if err != nil {
		return s1, s2, nil, fmt.Errorf("comparing artifacts: %w", err)
	}

	return s1, s2, changes, nil
}
