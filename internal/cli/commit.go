package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	message        string
	commitTags     []string
	commitMetaPairs map[string]string
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new Intelligence Commit",
	Long:  "Creates an immutable commit snapshot. Supports tags (--tag) and execution metadata (--meta key=value).",
	Run: func(cmd *cobra.Command, args []string) {
		if message == "" {
			fmt.Println("Commit message is required.")
			return
		}

		if err := repository.CreateCommitWithOpts(message, commitTags, commitMetaPairs); err != nil {
			fmt.Printf("Error creating commit: %v\n", err)
			return
		}

		fmt.Println("Commit created successfully.")
	},
}

func init() {
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (required)")
	commitCmd.Flags().StringSliceVarP(&commitTags, "tag", "t", []string{}, "Tag labels for this commit (e.g. --tag production --tag v1.0)")
	commitCmd.Flags().StringToStringVar(&commitMetaPairs, "meta", map[string]string{}, "Execution metadata key=value pairs (e.g. --meta model=gpt-4o)")

	_ = commitCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(commitCmd)
}
