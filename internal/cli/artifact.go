package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage and inspect typed AI artifacts",
	Long:  "Commands to register, list, inspect, and perform semantic diffs on typed AI artifacts (prompts, tools, memory, etc.).",
}

var artifactAddCmd = &cobra.Command{
	Use:   "add <type> <name> <path>",
	Short: "Register a typed AI artifact into evolution.manifest.json",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		artType, name, path := args[0], args[1], args[2]

		if err := repository.RegisterArtifactInManifest(artType, name, path); err != nil {
			fmt.Printf("Error registering artifact: %v\n", err)
			return
		}

		fmt.Printf("Registered %s artifact '%s' (path: %s) in evolution.manifest.json\n", artType, name, path)
	},
}

var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all AI artifacts attached to the current HEAD commit",
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.Exists() {
			fmt.Println("Error: Not an Evolution repository (run 'evo init' first)")
			return
		}

		artifacts, err := repository.GetHeadCommitArtifacts()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(artifacts) == 0 {
			fmt.Println("No artifacts attached to HEAD commit.")
			return
		}

		fmt.Println("Attached Artifacts in HEAD Commit:")
		for category, list := range artifacts {
			fmt.Printf("\n  %s[%s]%s (%d)\n", colorCyan, category, colorReset, len(list))
			for _, a := range list {
				hashStr := a.GetHash()
				if len(hashStr) > 8 {
					hashStr = hashStr[:8]
				} else if hashStr == "" {
					hashStr = "unhashed"
				}
				fmt.Printf("    %-20s %-10s path: %s\n", a.GetName(), fmt.Sprintf("(%s)", hashStr), a.GetPath())
			}
		}
		fmt.Println()
	},
}

var artifactShowCmd = &cobra.Command{
	Use:   "show <type> <hash>",
	Short: "Display content of a stored artifact",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		artType, hash := args[0], args[1]

		artifact, err := repository.LoadArtifact(artType, hash)
		if err != nil {
			fmt.Printf("Error loading artifact: %v\n", err)
			return
		}

		data, err := json.MarshalIndent(artifact, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting artifact: %v\n", err)
			return
		}

		fmt.Println(string(data))
	},
}

var artifactDiffCmd = &cobra.Command{
	Use:   "diff <commit1> <commit2>",
	Short: "Perform semantic diff between AI artifacts across two commits",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		commitID1, commitID2 := args[0], args[1]

		changes, err := repository.CompareCommitArtifacts(commitID1, commitID2)
		if err != nil {
			fmt.Printf("Error performing semantic diff: %v\n", err)
			return
		}

		if len(changes) == 0 {
			fmt.Println("No artifact changes detected between commits.")
			return
		}

		fmt.Printf("Semantic Artifact Diff (%s..%s):\n\n", commitID1[:8], commitID2[:8])
		for _, c := range changes {
			var color string
			switch c.Action {
			case "added":
				color = colorGreen
			case "modified":
				color = colorYellow
			case "removed":
				color = colorRed
			default:
				color = colorReset
			}

			fmt.Printf("  %s[%-8s]%s  %-15s %s: %s\n", color, c.Action, colorReset, c.Category, c.ArtifactName, c.Details)
		}
		fmt.Println()
	},
}

func init() {
	artifactCmd.AddCommand(artifactAddCmd)
	artifactCmd.AddCommand(artifactListCmd)
	artifactCmd.AddCommand(artifactShowCmd)
	artifactCmd.AddCommand(artifactDiffCmd)

	rootCmd.AddCommand(artifactCmd)
}
