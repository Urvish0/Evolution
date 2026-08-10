package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	replayExportFlag    string
	replayExecutionFlag string
)

var replayCmd = &cobra.Command{
	Use:   "replay [commit_id|branch]",
	Short: "Reconstruct and replay historical AI system state from an Intelligence Commit",
	Long:  "Reconstructs prompts, model config, retrieval, memory, tools, and policies for any commit, export as manifest, or verify against a recorded execution.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var target string

		// Handle --execution flag
		if replayExecutionFlag != "" {
			state, exec, err := repository.ReplayExecution(replayExecutionFlag)
			if err != nil {
				fmt.Printf("Error replaying execution: %v\n", err)
				return
			}

			fmt.Printf("%s=== Replaying Execution %s ===%s\n", colorCyan, exec.ID[:8], colorReset)
			fmt.Printf("Commit:    %s%s%s\n", colorYellow, exec.CommitID[:8], colorReset)
			fmt.Printf("Timestamp: %s\n", exec.Timestamp)
			fmt.Printf("Duration:  %dms\n", exec.DurationMs)
			fmt.Printf("Tokens:    %d total (%d prompt, %d completion)\n\n",
				exec.Tokens.TotalTokens, exec.Tokens.PromptTokens, exec.Tokens.CompletionTokens)

			fmt.Printf("%s[INPUTS]%s\n%s\n\n", colorGreen, colorReset, exec.Inputs)
			fmt.Printf("%s[RECORDED OUTPUTS]%s\n%s\n\n", colorGreen, colorReset, exec.Outputs)

			fmt.Printf("%s=== Reconstructed AI System State ===%s\n", colorCyan, colorReset)
			printReconstructedState(state)

			if replayExportFlag != "" {
				if err := repository.ExportReconstructedState(state, replayExportFlag); err != nil {
					fmt.Printf("Error exporting manifest: %v\n", err)
					return
				}
				fmt.Printf("\nExported reconstructed manifest to '%s'\n", replayExportFlag)
			}
			return
		}

		if len(args) == 1 {
			target = args[0]
		} else {
			currentBranch, err := repository.GetCurrentBranchName()
			if err != nil {
				fmt.Printf("Error getting current branch: %v\n", err)
				return
			}
			target = currentBranch
		}

		state, err := repository.ReconstructCommitState(target)
		if err != nil {
			fmt.Printf("Error reconstructing state: %v\n", err)
			return
		}

		printReconstructedState(state)

		if replayExportFlag != "" {
			if err := repository.ExportReconstructedState(state, replayExportFlag); err != nil {
				fmt.Printf("Error exporting manifest: %v\n", err)
				return
			}
			fmt.Printf("\nExported reconstructed manifest to '%s'\n", replayExportFlag)
		}
	},
}

func printReconstructedState(state *repository.ReconstructedState) {
	fmt.Printf("Commit ID: %s%s%s\n", colorYellow, state.CommitID[:8], colorReset)
	fmt.Printf("Author:    %s\n", state.Author)
	fmt.Printf("Message:   %s\n", state.CommitMsg)

	if len(state.Tags) > 0 {
		fmt.Printf("Tags:      %v\n", state.Tags)
	}

	manifest := state.Manifest
	fmt.Printf("\n%sAttached Artifacts (%d types):%s\n", colorCyan, countArtifactTypes(manifest), colorReset)

	if manifest.Artifacts.ModelConfig != nil {
		mc := manifest.Artifacts.ModelConfig
		fmt.Printf("  %s[model_config]%s %s (model: %s, provider: %s, temp: %.2f)\n",
			colorGreen, colorReset, mc.Name, mc.Model, mc.Provider, mc.Temperature)
	}

	for _, p := range manifest.Artifacts.Prompts {
		fmt.Printf("  %s[prompt]%s       %s (%s, format: %s) -> path: %s\n",
			colorGreen, colorReset, p.Name, p.Role, p.Format, p.Path)
	}

	for _, m := range manifest.Artifacts.Memory {
		fmt.Printf("  %s[memory]%s       %s (strategy: %s, max_tokens: %d)\n",
			colorGreen, colorReset, m.Name, m.Strategy, m.MaxTokens)
	}

	for _, r := range manifest.Artifacts.Retrieval {
		fmt.Printf("  %s[retrieval]%s    %s (source: %s, top_k: %d)\n",
			colorGreen, colorReset, r.Name, r.Source, r.TopK)
	}

	for _, t := range manifest.Artifacts.Tools {
		fmt.Printf("  %s[tool]%s         %s (provider: %s)\n",
			colorGreen, colorReset, t.Name, t.Provider)
	}

	for _, pol := range manifest.Artifacts.Policies {
		fmt.Printf("  %s[policy]%s       %s (enforcement: %s)\n",
			colorGreen, colorReset, pol.Name, pol.Enforcement)
	}
}

func countArtifactTypes(m repository.Manifest) int {
	count := 0
	if m.Artifacts.ModelConfig != nil {
		count++
	}
	if len(m.Artifacts.Prompts) > 0 {
		count++
	}
	if len(m.Artifacts.Memory) > 0 {
		count++
	}
	if len(m.Artifacts.Retrieval) > 0 {
		count++
	}
	if len(m.Artifacts.Tools) > 0 {
		count++
	}
	if len(m.Artifacts.Policies) > 0 {
		count++
	}
	return count
}

func init() {
	replayCmd.Flags().StringVarP(&replayExportFlag, "export", "e", "", "Export reconstructed state as a manifest file (e.g. evolution.manifest.json)")
	replayCmd.Flags().StringVar(&replayExecutionFlag, "execution", "", "Replay and compare state against a recorded execution ID")
	rootCmd.AddCommand(replayCmd)
}
