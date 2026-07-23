package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Evolution Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Evolution v%s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
