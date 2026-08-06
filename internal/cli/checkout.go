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
		target := args[0]

		if createBranchFlag {
			if err := repository.CreateBranch(target); err != nil {
				fmt.Printf("Error creating branch: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", target)
		}

		if err := repository.CheckoutBranch(target); err == nil {
			fmt.Printf("Switched to branch '%s'\n", target)
			return
		}

		// Fallback to checking out a specific commit (detached HEAD)
		if err := repository.CheckoutCommit(target); err == nil {
			fmt.Printf("HEAD is now at commit '%s' (detached HEAD)\n", target)
			return
		}

		fmt.Printf("Error: branch or commit '%s' not found\n", target)
	},
}

func init() {
	checkoutCmd.Flags().BoolVarP(&createBranchFlag, "create", "b", false, "Create new branch and switch to it")
	rootCmd.AddCommand(checkoutCmd)
}
