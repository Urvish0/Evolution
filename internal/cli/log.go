package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	onelineFlag bool
	limitFlag   int
)

const (
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorReset  = "\033[0m"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows the commit history",
	Long:  "Displays historical commits. Supports --oneline format, -n <count> limit, and short commit IDs.",
	Run: func(cmd *cobra.Command, args []string) {
		commits, err := repository.Log()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(commits) == 0 {
			fmt.Println("No commits yet.")
			return
		}

		// Apply limit if specified
		if limitFlag > 0 && limitFlag < len(commits) {
			commits = commits[:limitFlag]
		}

		for _, commit := range commits {
			shortID := commit.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			tagStr := ""
			if len(commit.Metadata.Tags) > 0 {
				for _, t := range commit.Metadata.Tags {
					tagStr += fmt.Sprintf(" %s[%s]%s", colorGreen, t, colorReset)
				}
			}

			if onelineFlag {
				fmt.Printf("%s%s%s%s %s\n", colorYellow, shortID, colorReset, tagStr, commit.Message)
			} else {
				fmt.Printf("%scommit %s%s%s\n", colorYellow, shortID, colorReset, tagStr)
				if commit.TreeHash != "" {
					shortTree := commit.TreeHash
					if len(shortTree) > 8 {
						shortTree = shortTree[:8]
					}
					fmt.Printf("Tree:   %s%s%s\n", colorGray, shortTree, colorReset)
				}
				if commit.Author != "" {
					fmt.Printf("Author: %s%s%s\n", colorCyan, commit.Author, colorReset)
				}
				fmt.Printf("Date:   %s%s%s\n", colorGray, commit.Timestamp, colorReset)

				if len(commit.Metadata.Environment) > 0 {
					metaStr := ""
					for k, v := range commit.Metadata.Environment {
						metaStr += fmt.Sprintf("%s=%s ", k, v)
					}
					fmt.Printf("Meta:   %s%s%s\n", colorGray, metaStr, colorReset)
				}

				fmt.Printf("\n    %s\n\n", commit.Message)
			}
		}
	},
}

func init() {
	logCmd.Flags().BoolVar(&onelineFlag, "oneline", false, "Display commits in compact single-line format")
	logCmd.Flags().IntVarP(&limitFlag, "limit", "n", 0, "Limit the number of commits shown")
	rootCmd.AddCommand(logCmd)
}
