package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var createBranchFlag bool

var checkoutCmd = &cobra.Command{
	Use:   "checkout [-b] <branch-name>",
	Short: "Switch to a specified branch (or create and switch with -b)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]

		if createBranchFlag {
			if err := repository.CreateBranch(branchName); err != nil {
				fmt.Printf("Error creating branch: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", branchName)
		}

		if err := repository.CheckoutBranch(branchName); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Switched to branch '%s'\n", branchName)
	},
}

func init() {
	checkoutCmd.Flags().BoolVarP(&createBranchFlag, "create", "b", false, "Create new branch and switch to it")
	rootCmd.AddCommand(checkoutCmd)
}
