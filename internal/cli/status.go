package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show repository status",
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.Exists() {
			fmt.Println("Not an Evolution repository.")
			return
		}

		status, err := repository.GetStatus()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("On branch %s%s%s\n", colorCyan, status.Branch, colorReset)
		fmt.Printf("Commits: %d\n\n", status.Commits)

		if status.Clean {
			fmt.Println("nothing to commit, working tree clean")
			return
		}

		wts := status.WorkingTree

		if len(wts.Staged) > 0 {
			fmt.Println("Changes to be committed:")
			for _, f := range wts.Staged {
				fmt.Printf("  %s%s:   %s%s\n", colorGreen, f.Status, f.Path, colorReset)
			}
			fmt.Println()
		}

		if len(wts.Modified) > 0 {
			fmt.Println("Changes not staged for commit:")
			for _, f := range wts.Modified {
				fmt.Printf("  %s%s:   %s%s\n", colorRed, f.Status, f.Path, colorReset)
			}
			fmt.Println()
		}

		if len(wts.Deleted) > 0 {
			fmt.Println("Deleted files:")
			for _, f := range wts.Deleted {
				fmt.Printf("  %s%s:   %s%s\n", colorRed, f.Status, f.Path, colorReset)
			}
			fmt.Println()
		}

		if len(wts.Untracked) > 0 {
			fmt.Println("Untracked files:")
			for _, f := range wts.Untracked {
				fmt.Printf("  %s%s%s\n", colorRed, f.Path, colorReset)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
