package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "evo",
	Short: "Evolution CLI",
	Long:  "Evolution is an AI engineering platform for versioning, replaying, and evaluating AI systems.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Evolution!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
