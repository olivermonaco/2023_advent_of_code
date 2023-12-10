package scratchcards

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

func CalculatePartOne(ctx context.Context, input []string) int {
	ctx = log.Logger.WithContext(ctx)
	l := log.Ctx(ctx).With().Logger()

	var result int
	for _, s := range input {
		haveNums, winningNums := separateNums(s)
		var numWinners *int
		for k := range haveNums {
			if _, ok := winningNums[k]; ok {
				numWinners = kit.Ptr(kit.Deref(numWinners) + 1)
				l.Info().Int("winning_num", k).Msg("")
			}
		}
		if numWinners != nil {
			result += calcLineResult(*numWinners)
		}
	}
	return result
}

func CalculatePartTwo(ctx context.Context, lines []string) int {
	var total int
	for i := 0; i < len(lines); i++ {
		lineTotal := recursiveLineCount(ctx, i, lines)
		total += lineTotal
	}

	return total
}

func recursiveLineCount(ctx context.Context, curIdx int, lines []string) int {
	if curIdx >= len(lines) {
		return 0
	}
	resultCount := 1
	numWinners := resultForLine(lines, curIdx)
	for i := 1; i < numWinners+1; i++ {
		if curIdx+i >= len(lines) {
			break
		}
		resultCount += recursiveLineCount(ctx, curIdx+i, lines)
	}
	return resultCount
}

func toIntMap(strs []string) map[int]struct{} {
	ret := make(map[int]struct{}, len(strs))
	for _, s := range strs {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ret[i] = struct{}{}
	}
	return ret
}

func separateNums(s string) (map[int]struct{}, map[int]struct{}) {
	_, numbersStr, _ := strings.Cut(s, ":")
	haveNumsStr, winningNumsStr, _ := strings.Cut(numbersStr, "|")
	haveNumsSl := strings.Fields(haveNumsStr)
	winningNumsSl := strings.Fields(winningNumsStr)

	return toIntMap(haveNumsSl), toIntMap(winningNumsSl)

}

func winnersForLine(haveNums, winningNums map[int]struct{}) int {

	var numWinners int

	for k := range haveNums {
		if _, ok := winningNums[k]; ok {
			numWinners++
		}
	}
	return numWinners
}

func calcLineResult(numWinners int) int {
	return max(int(math.Pow(float64(2), float64(numWinners-1))), 1)
}

func resultForLine(lines []string, idx int) int {
	if idx >= len(lines) {
		return 0
	}

	haveNums, winningNums := separateNums(lines[idx])
	winners := winnersForLine(haveNums, winningNums)
	return winners
}
