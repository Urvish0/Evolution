package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var (
	recordInputFlag          string
	recordOutputFlag         string
	recordPromptTokensFlag   int
	recordCompTokensFlag     int
	recordDurationMsFlag     int64
	recordMetaPairs          map[string]string
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record an AI execution run bound to the current HEAD commit",
	Long:  "Records an AI system execution run including inputs, outputs, token consumption, latency, and custom metadata.",
	Run: func(cmd *cobra.Command, args []string) {
		if recordInputFlag == "" || recordOutputFlag == "" {
			fmt.Println("Error: Both --input and --output flags are required for recording.")
			return
		}

		tokens := repository.TokenUsage{
			PromptTokens:     recordPromptTokensFlag,
			CompletionTokens: recordCompTokensFlag,
		}

		exec, err := repository.RecordExecution(recordInputFlag, recordOutputFlag, recordDurationMsFlag, tokens, recordMetaPairs)
		if err != nil {
			fmt.Printf("Error recording execution: %v\n", err)
			return
		}

		fmt.Printf("Recorded Execution '%s' (Commit: %s)\n", exec.ID[:8], exec.CommitID[:8])
	},
}

var executionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Inspect recorded AI executions",
	Long:  "Commands to list and view full JSON details of recorded AI execution runs.",
}

var executionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded AI executions",
	Run: func(cmd *cobra.Command, args []string) {
		executions, err := repository.ListExecutions()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(executions) == 0 {
			fmt.Println("No recorded executions found.")
			return
		}

		fmt.Printf("%s%-10s %-10s %-10s %-12s %s%s\n",
			colorCyan, "ID", "COMMIT", "TOKENS", "DURATION", "TIMESTAMP", colorReset)

		for _, e := range executions {
			shortID := e.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}
			shortCommit := e.CommitID
			if len(shortCommit) > 8 {
				shortCommit = shortCommit[:8]
			}

			fmt.Printf("%s%-10s%s %s%-10s%s %-10d %-12s %s\n",
				colorYellow, shortID, colorReset,
				colorGray, shortCommit, colorReset,
				e.Tokens.TotalTokens,
				fmt.Sprintf("%dms", e.DurationMs),
				e.Timestamp)
		}
	},
}

var executionShowCmd = &cobra.Command{
	Use:   "show <execution_id>",
	Short: "Display full JSON details for a recorded execution",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execID := args[0]

		exec, err := repository.LoadExecution(execID)
		if err != nil {
			fmt.Printf("Error loading execution: %v\n", err)
			return
		}

		data, err := json.MarshalIndent(exec, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting execution: %v\n", err)
			return
		}

		fmt.Println(string(data))
	},
}

func init() {
	recordCmd.Flags().StringVarP(&recordInputFlag, "input", "i", "", "Execution input text / user query (required)")
	recordCmd.Flags().StringVarP(&recordOutputFlag, "output", "o", "", "Execution output text / AI response (required)")
	recordCmd.Flags().IntVar(&recordPromptTokensFlag, "tokens-prompt", 0, "Prompt tokens consumed")
	recordCmd.Flags().IntVar(&recordCompTokensFlag, "tokens-completion", 0, "Completion tokens consumed")
	recordCmd.Flags().Int64Var(&recordDurationMsFlag, "duration", 0, "Execution duration in milliseconds")
	recordCmd.Flags().StringToStringVar(&recordMetaPairs, "meta", map[string]string{}, "Custom metadata key=value pairs")

	executionCmd.AddCommand(executionListCmd)
	executionCmd.AddCommand(executionShowCmd)

	rootCmd.AddCommand(recordCmd)
	rootCmd.AddCommand(executionCmd)
}
