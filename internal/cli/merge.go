package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <branch>",
	Short: "Merge another branch into the active branch",
	Long:  "Performs a three-way merge from target branch into HEAD. Handles fast-forward merges, auto-merges, and injects conflict markers if conflicts exist.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetBranch := args[0]

		result, err := repository.MergeBranch(targetBranch)
		if err != nil {
			fmt.Printf("Error merging branch: %v\n", err)
			return
		}

		if result.IsFastForward {
			fmt.Printf("Fast-forward merged branch '%s'\n", targetBranch)
			return
		}

		if result.HasConflicts {
			fmt.Printf("⚠️  Merge conflict in %d file(s):\n", len(result.ConflictedFiles))
			for _, file := range result.ConflictedFiles {
				fmt.Printf("  %s%s%s\n", colorRed, file, colorReset)
			}
			fmt.Println("\nConflict markers (<<<<<<< / ======= / >>>>>>>) have been injected into conflicted files.")
			fmt.Println("Resolve conflicts manually, then stage and commit.")
			return
		}

		fmt.Printf("Successfully merged '%s' (Merge commit: %s)\n", targetBranch, result.MergeCommitID[:8])
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
