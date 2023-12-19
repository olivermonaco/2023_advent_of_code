package environment_report

import (
	"context"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func CalculatePartOne(ctx context.Context, input []string) int {
	var result int
	for _, line := range input {
		history := parseLine(line)
		histResults := recursiveCreateDiffs(history)
		var lineTotal int
		histResults = append(histResults, history[len(history)-1])
		for _, histResult := range histResults {
			lineTotal += histResult
		}
		result += lineTotal
	}

	return result
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func parseLine(line string) []int {
	return kit.Map(strings.Fields(line), func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return n
	})
}

func recursiveCreateDiffs(numberLine []int) []int {
	diffs := makeLineDiffs(numberLine)

	if len(diffs) == 0 || diffs[len(diffs)-1] == 0 {
		return []int{}
	}
	lastVal := []int{diffs[len(diffs)-1]}
	return append(lastVal, recursiveCreateDiffs(diffs)...)

}

func makeLineDiffs(numberLine []int) []int {
	diffs := make([]int, 0, len(numberLine)-1)
	for i := 1; i < len(numberLine); i++ {
		addDiff := numberLine[i] - numberLine[i-1]
		diffs = append(diffs, addDiff)
	}
	return diffs
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
