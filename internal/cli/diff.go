package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between working tree and HEAD",
	Long:  "Compares the current workspace files against the last committed snapshot to show new, modified, and deleted files.",
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

		if wts.IsClean() {
			fmt.Println("No changes detected. Working tree matches HEAD.")
			return
		}

		totalChanges := len(wts.Staged) + len(wts.Modified) + len(wts.Untracked) + len(wts.Deleted)
		fmt.Printf("Diff against HEAD (%d change(s)):\n\n", totalChanges)

		for _, f := range wts.Staged {
			fmt.Printf("  %s[staged]   %s%s\n", colorGreen, f.Path, colorReset)
		}

		for _, f := range wts.Modified {
			fmt.Printf("  %s[modified] %s%s\n", colorYellow, f.Path, colorReset)

			// Show content-level diff
			diffLines, err := repository.GetContentDiff(f.Path)
			if err == nil {
				printContentDiff(diffLines)
			}
		}

		for _, f := range wts.Deleted {
			fmt.Printf("  %s[deleted]  %s%s\n", colorRed, f.Path, colorReset)
		}

		for _, f := range wts.Untracked {
			fmt.Printf("  %s[new]      %s%s\n", colorCyan, f.Path, colorReset)
		}

		fmt.Println()
	},
}

// printContentDiff renders a line-by-line diff with color-coded output.
func printContentDiff(lines []repository.DiffLine) {
	hasChanges := false
	for _, line := range lines {
		if line.Op != repository.DiffEqual {
			hasChanges = true
			break
		}
	}

	if !hasChanges {
		return
	}

	fmt.Println()
	for _, line := range lines {
		switch line.Op {
		case repository.DiffDelete:
			fmt.Printf("    %s- %s%s\n", colorRed, line.Content, colorReset)
		case repository.DiffInsert:
			fmt.Printf("    %s+ %s%s\n", colorGreen, line.Content, colorReset)
		case repository.DiffEqual:
			fmt.Printf("      %s\n", line.Content)
		}
	}
	fmt.Println()
}
func init() {
	rootCmd.AddCommand(diffCmd)
}
