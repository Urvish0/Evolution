package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	deleteBranchFlag string
	newBranchFlag    string
)

var branchCmd = &cobra.Command{
	Use:   "branch [name]",
	Short: "List, create, or delete branches",
	Long:  "With no arguments or flags, lists all local branches.\nUse -n / --new <name> (or evo branch <name>) to create a new branch.\nUse -d / --delete <name> to delete a branch.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Case 1: Delete branch flag (-d <name>)
		if deleteBranchFlag != "" {
			if err := repository.DeleteBranch(deleteBranchFlag); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Deleted branch '%s'\n", deleteBranchFlag)
			return
		}

		// Case 2: New branch flag (-n <name>)
		if newBranchFlag != "" {
			if err := repository.CreateBranch(newBranchFlag); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", newBranchFlag)
			return
		}

		// Case 3: Positional argument provided (evo branch <name>)
		if len(args) == 1 {
			branchName := args[0]
			if err := repository.CreateBranch(branchName); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", branchName)
			return
		}

		// Case 4: No arguments or flags -> list all branches
		branches, err := repository.ListBranches()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		currentBranch, err := repository.GetCurrentBranchName()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		for _, b := range branches {
			if b.Name == currentBranch {
				fmt.Printf("* %s\n", b.Name)
			} else {
				fmt.Printf("  %s\n", b.Name)
			}
		}
	},
}

func init() {
	branchCmd.Flags().StringVarP(&newBranchFlag, "new", "n", "", "Create a new branch")
	branchCmd.Flags().StringVarP(&deleteBranchFlag, "delete", "d", "", "Delete a branch")
	rootCmd.AddCommand(branchCmd)
}
