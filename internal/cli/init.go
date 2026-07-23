package cli

import (
	"fmt"
	"os"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an Evolution repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repository.Init(); err != nil {
			fmt.Println(err)
			return
		}

		// fmt.Println("Repository initialized!")
		cwd, _ := os.Getwd()
		fmt.Printf("Evolution repository initialized at:\n%s\n\n", cwd)
		fmt.Println("Try:")
		fmt.Println(" evo --help")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
