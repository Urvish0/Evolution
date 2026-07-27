package cli

import (
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set repository or user configuration",
	Long:  "Manage settings such as user.name and user.email.",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value (e.g. user.name 'Urvish')",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := repository.SetUserSetting(key, value); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Set %s = %q\n", key, value)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all user configurations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := repository.LoadUserConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("user.name  = %s\n", cfg.Name)
		fmt.Printf("user.email = %s\n", cfg.Email)
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
