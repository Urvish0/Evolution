package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	manifestNameFlag string
	manifestDescFlag string
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Manage the Intelligence Manifest (evolution.manifest.json)",
	Long:  "Commands to initialize, validate, and view the Intelligence Manifest describing AI system state.",
}

var manifestInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter evolution.manifest.json file",
	Run: func(cmd *cobra.Command, args []string) {
		if repository.ManifestExists() {
			fmt.Printf("Error: %s already exists in the current directory\n", repository.ManifestFileName)
			return
		}

		if err := repository.InitManifest(manifestNameFlag, manifestDescFlag); err != nil {
			fmt.Printf("Error creating manifest: %v\n", err)
			return
		}

		fmt.Printf("Created %s conforming to Spec v0.1\n", repository.ManifestFileName)
	},
}

var manifestValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate evolution.manifest.json against Specification v0.1",
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.ManifestExists() {
			fmt.Printf("Error: %s not found in current directory\n", repository.ManifestFileName)
			return
		}

		m, err := repository.LoadManifest(repository.ManifestFileName)
		if err != nil {
			fmt.Printf("Error loading manifest: %v\n", err)
			return
		}

		if err := repository.ValidateManifest(m); err != nil {
			fmt.Printf("❌ Validation Failed: %v\n", err)
			return
		}

		fmt.Printf("✅ %s is valid! (Spec v%s, Name: %s)\n", repository.ManifestFileName, m.Version, m.Name)
	},
}

var manifestShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the active Intelligence Manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if !repository.ManifestExists() {
			fmt.Printf("Error: %s not found in current directory\n", repository.ManifestFileName)
			return
		}

		m, err := repository.LoadManifest(repository.ManifestFileName)
		if err != nil {
			fmt.Printf("Error loading manifest: %v\n", err)
			return
		}

		data, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting manifest: %v\n", err)
			return
		}

		fmt.Println(string(data))
	},
}

func init() {
	manifestInitCmd.Flags().StringVarP(&manifestNameFlag, "name", "n", "", "Name of the intelligence system")
	manifestInitCmd.Flags().StringVarP(&manifestDescFlag, "description", "d", "", "Description of the intelligence system")

	manifestCmd.AddCommand(manifestInitCmd)
	manifestCmd.AddCommand(manifestValidateCmd)
	manifestCmd.AddCommand(manifestShowCmd)

	rootCmd.AddCommand(manifestCmd)
}
