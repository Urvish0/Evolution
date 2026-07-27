package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout <branch-name>",
	Short: "Switch to a specified branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]

		if err := repository.CheckoutBranch(branchName); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Switched to branch '%s'\n", branchName)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
