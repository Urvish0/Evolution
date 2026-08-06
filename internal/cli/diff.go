package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [rev1] [rev2]",
	Short: "Show changes between working tree and HEAD, or between revisions/branches",
	Long:  "Compares content line-by-line.\nNo args: working tree vs HEAD\n1 arg: HEAD vs revision\n2 args: revision1 vs revision2",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.Exists() {
			fmt.Println("Error: Not an Evolution repository (run 'evo init' first)")
			return
		}

		// Case 1: Two revisions provided (evo diff rev1 rev2)
		if len(args) == 2 {
			diffOutput, err := repository.GetRevisionDiff(args[0], args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if diffOutput == "" {
				fmt.Println("No differences found.")
				return
			}
			fmt.Print(diffOutput)
			return
		}

		// Case 2: One revision provided (evo diff rev1) -> compare HEAD vs rev1
		if len(args) == 1 {
			currentBranch, _ := repository.GetCurrentBranchName()
			diffOutput, err := repository.GetRevisionDiff(currentBranch, args[0])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if diffOutput == "" {
				fmt.Println("No differences found.")
				return
			}
			fmt.Print(diffOutput)
			return
		}

		// Case 3: 0 args -> Working tree vs HEAD
		wts, err := repository.CompareWorkingTree()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

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

			adds, dels := 0, 0
			for _, d := range diffLines {
				if d.Op == repository.DiffInsert {
					adds++
				} else if d.Op == repository.DiffDelete {
					dels++
				}
			}

			fmt.Printf("%s--- a/%s%s\n", colorRed, f.Path, colorReset)
			fmt.Printf("%s+++ b/%s%s\n", colorGreen, f.Path, colorReset)
			fmt.Printf("%s@@ -%d +%d @@%s\n", colorCyan, dels, adds, colorReset)

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
