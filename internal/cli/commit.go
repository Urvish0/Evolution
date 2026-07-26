package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var message string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new commit",
	Run: func(cmd *cobra.Command, args []string) {
		if message == "" {
			fmt.Println("Commit message is required.")
			return
		}

		if err := repository.CreateCommit(message); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Commit created successfully.")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(
		&message,
		"message",
		"m",
		"",
		"Commit message",
	)
}
