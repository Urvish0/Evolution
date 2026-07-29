package repository

import (
	"strings"
)

// DiffOp represents the type of change for a single line.
type DiffOp int

const (
	DiffEqual  DiffOp = iota // Line is unchanged
	DiffDelete               // Line was removed
	DiffInsert               // Line was added
)

// DiffLine represents a single line in a unified diff output.
type DiffLine struct {
	Op      DiffOp
	Content string
}

// ComputeLineDiff computes a line-by-line diff between oldText and newText
// using the Longest Common Subsequence (LCS) algorithm.
// Returns a slice of DiffLine entries representing the minimal edit script.
func ComputeLineDiff(oldText, newText string) []DiffLine {
	oldLines := splitLines(oldText)
	newLines := splitLines(newText)

	// Build LCS table
	lcs := buildLCSTable(oldLines, newLines)

	// Backtrack through LCS table to produce the edit script
	return backtrackDiff(oldLines, newLines, lcs)
}

// splitLines splits text into lines, handling empty input gracefully.
func splitLines(text string) []string {
	if text == "" {
		return []string{}
	}
	// Remove trailing newline to avoid a phantom empty line at the end
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return []string{}
	}
	return strings.Split(text, "\n")
}

// buildLCSTable constructs the dynamic programming table for LCS.
// lcs[i][j] = length of LCS of oldLines[:i] and newLines[:j].
func buildLCSTable(oldLines, newLines []string) [][]int {
	m := len(oldLines)
	n := len(newLines)

	// Allocate (m+1) x (n+1) table initialized to 0
	lcs := make([][]int, m+1)
	for i := range lcs {
		lcs[i] = make([]int, n+1)
	}

	// Fill bottom-up
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else {
				if lcs[i-1][j] >= lcs[i][j-1] {
					lcs[i][j] = lcs[i-1][j]
				} else {
					lcs[i][j] = lcs[i][j-1]
				}
			}
		}
	}

	return lcs
}

// backtrackDiff walks the LCS table backwards to produce the diff output.
func backtrackDiff(oldLines, newLines []string, lcs [][]int) []DiffLine {
	var result []DiffLine

	i := len(oldLines)
	j := len(newLines)

	// Backtrack from bottom-right corner
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			// Lines match — part of LCS (unchanged)
			result = append(result, DiffLine{Op: DiffEqual, Content: oldLines[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			// Line was added in new version
			result = append(result, DiffLine{Op: DiffInsert, Content: newLines[j-1]})
			j--
		} else if i > 0 {
			// Line was removed from old version
			result = append(result, DiffLine{Op: DiffDelete, Content: oldLines[i-1]})
			i--
		}
	}

	// Reverse because we built it backwards
	reverse(result)
	return result
}

// reverse reverses a slice of DiffLine in place.
func reverse(s []DiffLine) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
