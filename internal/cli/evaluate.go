package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Urvish0/evolution/internal/repository"
	"github.com/spf13/cobra"
)

var evaluateCompareFlag []string

var evaluateCmd = &cobra.Command{
	Use:   "evaluate [target1] [target2]",
	Short: "Run AI quality evaluations or compare evaluation scores across commits",
	Long:  "Evaluates latency, token cost, safety guardrails, and correctness for executions, commits, or compares two commits side-by-side.",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Case 1: --compare flag with 2 targets OR 2 positional arguments
		var target1, target2 string
		if len(evaluateCompareFlag) == 2 {
			target1, target2 = evaluateCompareFlag[0], evaluateCompareFlag[1]
		} else if len(args) == 2 {
			target1, target2 = args[0], args[1]
		}

		if target1 != "" && target2 != "" {
			comp, err := repository.CompareCommitEvaluations(target1, target2)
			if err != nil {
				fmt.Printf("Error comparing commit evaluations: %v\n", err)
				return
			}

			fmt.Printf("%s=== Cross-Commit Quality Score Comparison ===%s\n", colorCyan, colorReset)
			fmt.Printf("Target 1: %s%s%s - %s (%d runs)\n",
				colorYellow, comp.Report1.CommitID[:8], colorReset, comp.Report1.CommitMsg, comp.Report1.ExecutionCount)
			fmt.Printf("Target 2: %s%s%s - %s (%d runs)\n\n",
				colorYellow, comp.Report2.CommitID[:8], colorReset, comp.Report2.CommitMsg, comp.Report2.ExecutionCount)

			fmt.Printf("Overall Quality Score: %.2f vs %.2f  %s\n\n",
				comp.Report1.OverallScore, comp.Report2.OverallScore, formatTrendIndicator(comp.OverallScoreDelta))

			fmt.Printf("%sEvaluator Score Breakdown:%s\n", colorCyan, colorReset)
			for name, delta := range comp.EvaluatorDeltas {
				m1 := comp.Report1.EvaluatorMeans[name]
				m2 := comp.Report2.EvaluatorMeans[name]
				fmt.Printf("  %-14s %.2f vs %.2f  %s\n", name, m1, m2, formatTrendIndicator(delta))
			}
			return
		}

		// Case 2: 1 positional argument (execution ID OR commit ID / branch)
		if len(args) == 1 {
			target := args[0]

			// First, try as execution ID
			state, exec, err := repository.ReplayExecution(target)
			if err == nil {
				res, err := repository.EvaluateExecution(exec, state, nil)
				if err == nil {
					fmt.Printf("%s=== AI Execution Evaluation Report ===%s\n", colorCyan, colorReset)
					fmt.Printf("Evaluation ID: %s%s%s\n", colorYellow, res.ID[:8], colorReset)
					fmt.Printf("Execution ID:  %s%s%s\n", colorYellow, res.ExecutionID[:8], colorReset)
					fmt.Printf("Commit ID:     %s%s%s\n", colorYellow, res.CommitID[:8], colorReset)
					fmt.Printf("Overall Score: %s%.2f / 1.00%s\n\n", colorGreen, res.OverallScore, colorReset)

					fmt.Printf("%sEvaluator Breakdown:%s\n", colorCyan, colorReset)
					for name, score := range res.Scores {
						scoreColor := colorGreen
						if score.Score < 0.7 {
							scoreColor = colorYellow
						}
						if score.Score < 0.5 {
							scoreColor = colorRed
						}
						fmt.Printf("  %-14s %s%.2f%s (%s)\n", name, scoreColor, score.Score, colorReset, score.Details)
					}
					return
				}
			}

			// Fallback: evaluate commit / branch
			report, err := repository.EvaluateCommit(target)
			if err != nil {
				fmt.Printf("Error evaluating target %s: %v\n", target, err)
				return
			}

			fmt.Printf("%s=== Commit Intelligence Quality Report ===%s\n", colorCyan, colorReset)
			fmt.Printf("Commit ID:     %s%s%s\n", colorYellow, report.CommitID[:8], colorReset)
			fmt.Printf("Message:       %s\n", report.CommitMsg)
			fmt.Printf("Executions:    %d runs evaluated\n", report.ExecutionCount)
			fmt.Printf("Overall Score: %s%.2f / 1.00%s\n\n", colorGreen, report.OverallScore, colorReset)

			fmt.Printf("%sMean Evaluator Scores:%s\n", colorCyan, colorReset)
			for name, score := range report.EvaluatorMeans {
				scoreColor := colorGreen
				if score < 0.7 {
					scoreColor = colorYellow
				}
				if score < 0.5 {
					scoreColor = colorRed
				}
				fmt.Printf("  %-14s %s%.2f%s\n", name, scoreColor, score, colorReset)
			}
			return
		}

		// Case 3: 0 arguments provided -> evaluate active branch HEAD commit
		currentBranch, err := repository.GetCurrentBranchName()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			return
		}

		report, err := repository.EvaluateCommit(currentBranch)
		if err != nil {
			fmt.Printf("Error evaluating active branch: %v\n", err)
			return
		}

		fmt.Printf("%s=== Commit Intelligence Quality Report ===%s\n", colorCyan, colorReset)
		fmt.Printf("Commit ID:     %s%s%s\n", colorYellow, report.CommitID[:8], colorReset)
		fmt.Printf("Message:       %s\n", report.CommitMsg)
		fmt.Printf("Executions:    %d runs evaluated\n", report.ExecutionCount)
		fmt.Printf("Overall Score: %s%.2f / 1.00%s\n\n", colorGreen, report.OverallScore, colorReset)

		fmt.Printf("%sMean Evaluator Scores:%s\n", colorCyan, colorReset)
		for name, score := range report.EvaluatorMeans {
			scoreColor := colorGreen
			if score < 0.7 {
				scoreColor = colorYellow
			}
			if score < 0.5 {
				scoreColor = colorRed
			}
			fmt.Printf("  %-14s %s%.2f%s\n", name, scoreColor, score, colorReset)
		}
	},
}

func formatTrendIndicator(delta float64) string {
	if delta > 0.001 {
		return fmt.Sprintf("%s▲ +%.2f%s", colorGreen, delta, colorReset)
	} else if delta < -0.001 {
		return fmt.Sprintf("%s▼ %.2f%s", colorRed, delta, colorReset)
	}
	return fmt.Sprintf("%s▶  0.00%s", colorGray, colorReset)
}

var evaluateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded evaluation results",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := repository.ListEvaluationResults()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(list) == 0 {
			fmt.Println("No recorded evaluations found.")
			return
		}

		fmt.Printf("%s%-10s %-10s %-10s %-14s %s%s\n",
			colorCyan, "EVAL_ID", "EXEC_ID", "COMMIT", "OVERALL SCORE", "TIMESTAMP", colorReset)

		for _, r := range list {
			shortID := r.ID[:8]
			shortExec := r.ExecutionID[:8]
			shortCommit := r.CommitID[:8]

			scoreColor := colorGreen
			if r.OverallScore < 0.7 {
				scoreColor = colorYellow
			}
			if r.OverallScore < 0.5 {
				scoreColor = colorRed
			}

			fmt.Printf("%s%-10s%s %s%-10s%s %s%-10s%s %s%-14.2f%s %s\n",
				colorYellow, shortID, colorReset,
				colorGray, shortExec, colorReset,
				colorGray, shortCommit, colorReset,
				scoreColor, r.OverallScore, colorReset,
				r.Timestamp)
		}
	},
}

var evaluateShowCmd = &cobra.Command{
	Use:   "show <evaluation_id>",
	Short: "Display full JSON report for an evaluation result",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		res, err := repository.LoadEvaluationResult(id)
		if err != nil {
			fmt.Printf("Error loading evaluation: %v\n", err)
			return
		}

		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}

		fmt.Println(string(data))
	},
}

func init() {
	evaluateCmd.Flags().StringSliceVar(&evaluateCompareFlag, "compare", []string{}, "Compare evaluation scores across two commits/branches (e.g. --compare rev1,rev2)")
	evaluateCmd.AddCommand(evaluateListCmd)
	evaluateCmd.AddCommand(evaluateShowCmd)
	rootCmd.AddCommand(evaluateCmd)
}
