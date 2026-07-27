package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows the history",
	Run: func(cmd *cobra.Command, args []string) {
		commits, err := repository.Log()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, commit := range commits {
			fmt.Printf("commit %s\n", commit.ID)
			if commit.Author != "" {
				fmt.Printf("Author: %s\n", commit.Author)
			}
			fmt.Printf("Date:   %s\n\n", commit.Timestamp)
			fmt.Printf("    %s\n\n", commit.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
