package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	replayExportFlag            string
	replayExecutionFlag         string
	replayCompareExecutionsFlag []string
)

var replayCmd = &cobra.Command{
	Use:   "replay [commit1] [commit2]",
	Short: "Reconstruct, compare, and replay historical AI system states and executions",
	Long:  "Reconstructs prompts, model config, retrieval, memory, tools, and policies for any commit, compares two executions side-by-side, or exports as manifest.",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Case 1: --compare-executions flag with 2 execution IDs
		if len(replayCompareExecutionsFlag) == 2 {
			id1, id2 := replayCompareExecutionsFlag[0], replayCompareExecutionsFlag[1]
			comp, err := repository.CompareExecutions(id1, id2)
			if err != nil {
				fmt.Printf("Error comparing executions: %v\n", err)
				return
			}

			fmt.Printf("%s=== Execution Side-by-Side Comparison ===%s\n", colorCyan, colorReset)
			fmt.Printf("Execution 1: %s%s%s (Commit: %s)\n", colorYellow, comp.Exec1.ID[:8], colorReset, comp.Exec1.CommitID[:8])
			fmt.Printf("Execution 2: %s%s%s (Commit: %s)\n\n", colorYellow, comp.Exec2.ID[:8], colorReset, comp.Exec2.CommitID[:8])

			fmt.Printf("%sMetrics Comparison:%s\n", colorCyan, colorReset)
			fmt.Printf("  Prompt Tokens:     %d vs %d (%+d)\n", comp.Exec1.Tokens.PromptTokens, comp.Exec2.Tokens.PromptTokens, comp.PromptTokenDelta)
			fmt.Printf("  Completion Tokens: %d vs %d (%+d)\n", comp.Exec1.Tokens.CompletionTokens, comp.Exec2.Tokens.CompletionTokens, comp.CompTokenDelta)
			fmt.Printf("  Total Tokens:      %d vs %d (%+d)\n", comp.Exec1.Tokens.TotalTokens, comp.Exec2.Tokens.TotalTokens, comp.TotalTokenDelta)
			fmt.Printf("  Duration:          %dms vs %dms (%+dms)\n\n", comp.Exec1.DurationMs, comp.Exec2.DurationMs, comp.DurationDeltaMs)

			fmt.Printf("%sOutput Text Differences:%s\n", colorCyan, colorReset)
			for _, line := range comp.OutputDiffLines {
				switch line.Op {
				case repository.DiffInsert:
					fmt.Printf("%s+%s%s\n", colorGreen, line.Content, colorReset)
				case repository.DiffDelete:
					fmt.Printf("%s-%s%s\n", colorRed, line.Content, colorReset)
				case repository.DiffEqual:
					fmt.Printf(" %s\n", line.Content)
				}
			}
			return
		}

		// Case 2: 2 positional arguments provided (evo replay commit1 commit2)
		if len(args) == 2 {
			rev1, rev2 := args[0], args[1]
			s1, s2, changes, err := repository.CompareCommitReplays(rev1, rev2)
			if err != nil {
				fmt.Printf("Error comparing commit replays: %v\n", err)
				return
			}

			fmt.Printf("%s=== Comparing Reconstructed Intelligence States ===%s\n", colorCyan, colorReset)
			fmt.Printf("Revision 1: %s%s%s - %s\n", colorYellow, s1.CommitID[:8], colorReset, s1.CommitMsg)
			fmt.Printf("Revision 2: %s%s%s - %s\n\n", colorYellow, s2.CommitID[:8], colorReset, s2.CommitMsg)

			if len(changes) == 0 {
				fmt.Println("No artifact changes detected between revisions.")
				return
			}

			for _, c := range changes {
				fmt.Printf("%s[%s]%s %s (%s)\n", colorGreen, c.Category, colorReset, c.ArtifactName, c.Action)
				if c.Details != "" {
					fmt.Printf("  Details: %s\n", c.Details)
				}
			}
			return
		}

		// Case 3: Single execution replay (--execution <id>)
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

		// Case 4: Single commit / active branch replay
		var target string
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
	replayCmd.Flags().StringSliceVar(&replayCompareExecutionsFlag, "compare-executions", []string{}, "Compare two execution runs side-by-side by ID (e.g. --compare-executions id1,id2)")
	rootCmd.AddCommand(replayCmd)
}
