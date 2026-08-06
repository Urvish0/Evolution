package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Discard working directory changes and restore file from HEAD commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		if err := repository.RestoreFileFromHEAD(filePath); err != nil {
			fmt.Printf("Error restoring file: %v\n", err)
			return
		}

		fmt.Printf("Restored '%s' from HEAD commit\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
