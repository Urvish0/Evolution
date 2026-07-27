package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <path...>",
	Short: "Add file contents to the staging area (index)",
	Long:  "Stages specified files or directories for the next Intelligence Commit.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.Exists() {
			fmt.Println("Error: Not an Evolution repository (run 'evo init' first)")
			return
		}

		for _, path := range args {
			if err := repository.StagePath(path); err != nil {
				fmt.Printf("Error staging %s: %v\n", path, err)
				return
			}
			fmt.Printf("Staged '%s'\n", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
