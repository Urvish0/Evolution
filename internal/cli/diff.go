package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between working tree and HEAD",
	Long:  "Shows line-by-line content differences for modified tracked files between the working tree and the last commit.",
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.Exists() {
			fmt.Println("Error: Not an Evolution repository (run 'evo init' first)")
			return
		}

		wts, err := repository.CompareWorkingTree()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// evo diff only shows content diffs for modified tracked files
		// Use 'evo status' for untracked/staged/deleted file lists
		if len(wts.Modified) == 0 {
			fmt.Println("No modified files.")
			return
		}

		for _, f := range wts.Modified {
			diffLines, err := repository.GetContentDiff(f.Path)
			if err != nil {
				fmt.Printf("Error diffing %s: %v\n", f.Path, err)
				continue
			}

			// Count additions and deletions
			adds, dels := 0, 0
			for _, d := range diffLines {
				if d.Op == repository.DiffInsert {
					adds++
				} else if d.Op == repository.DiffDelete {
					dels++
				}
			}

			// File header (like git diff)
			fmt.Printf("%s--- a/%s%s\n", colorRed, f.Path, colorReset)
			fmt.Printf("%s+++ b/%s%s\n", colorGreen, f.Path, colorReset)
			fmt.Printf("%s@@ -%d +%d @@%s\n", colorCyan, dels, adds, colorReset)

			// Content lines
			for _, line := range diffLines {
				switch line.Op {
				case repository.DiffDelete:
					fmt.Printf("%s-%s%s\n", colorRed, line.Content, colorReset)
				case repository.DiffInsert:
					fmt.Printf("%s+%s%s\n", colorGreen, line.Content, colorReset)
				case repository.DiffEqual:
					fmt.Printf(" %s\n", line.Content)
				}
			}

			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
