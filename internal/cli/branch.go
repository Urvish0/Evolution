package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	deleteBranchFlag string
	newBranchFlag    string
	renameBranchFlag string
)

var branchCmd = &cobra.Command{
	Use:   "branch [name]",
	Short: "List, create, rename, or delete branches",
	Long:  "List all local branches with rich metadata, create a new branch (-n / --new), rename (-m / --move), or delete (-d / --delete).",
	Args:  cobra.MaximumNArgs(2),
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

		// Case 2: Rename branch (-m <new> or -m <old> <new> or 2 args: evo branch -m old new)
		if renameBranchFlag != "" || len(args) == 2 {
			var oldName, newName string
			if len(args) == 2 {
				oldName, newName = args[0], args[1]
			} else if len(args) == 1 {
				oldName, _ = repository.GetCurrentBranchName()
				newName = args[0]
			} else {
				oldName, _ = repository.GetCurrentBranchName()
				newName = renameBranchFlag
			}

			if err := repository.RenameBranch(oldName, newName); err != nil {
				fmt.Printf("Error renaming branch: %v\n", err)
				return
			}
			fmt.Printf("Renamed branch '%s' to '%s'\n", oldName, newName)
			return
		}

		// Case 3: New branch flag (-n <name>)
		if newBranchFlag != "" {
			if err := repository.CreateBranch(newBranchFlag); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", newBranchFlag)
			return
		}

		// Case 4: Positional argument provided (evo branch <name>)
		if len(args) == 1 {
			branchName := args[0]
			if err := repository.CreateBranch(branchName); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Created branch '%s'\n", branchName)
			return
		}

		// Case 5: No arguments or flags -> list all branches with rich details
		branches, err := repository.ListBranches()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		for _, b := range branches {
			details, err := repository.GetBranchDetails(b)
			if err != nil {
				details.Name = b.Name
			}

			prefix := "  "
			color := colorReset
			if details.IsActive {
				prefix = "* "
				color = colorGreen
			}

			headStr := "no commits"
			if details.HeadCommitID != "" {
				shortHash := details.HeadCommitID
				if len(shortHash) > 8 {
					shortHash = shortHash[:8]
				}
				headStr = shortHash
			}

			msgStr := ""
			if details.LastCommitMessage != "" {
				msgStr = fmt.Sprintf(" - %s", details.LastCommitMessage)
			}

			fmt.Printf("%s%s%-18s%s %s[%s]%s (%d commits)%s\n",
				color, prefix, details.Name, colorReset,
				colorYellow, headStr, colorReset,
				details.CommitCount, msgStr)
		}
	},
}

func init() {
	branchCmd.Flags().StringVarP(&newBranchFlag, "new", "n", "", "Create a new branch")
	branchCmd.Flags().StringVarP(&deleteBranchFlag, "delete", "d", "", "Delete a branch")
	branchCmd.Flags().StringVarP(&renameBranchFlag, "move", "m", "", "Rename a branch")
	rootCmd.AddCommand(branchCmd)
}
