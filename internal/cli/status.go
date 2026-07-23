package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
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

		fmt.Println("Repository : initialized")
		fmt.Printf("Branch     : %s\n", status.Branch)
		fmt.Printf("Commits    : %d\n", status.Commits)

		if status.Clean {
			fmt.Println("Working Tree : clean")
		} else {
			fmt.Println("Working Tree : dirty")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
