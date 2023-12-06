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
		var lineResult int
		if numWinners != nil {
			lineResult = max(int(math.Pow(float64(2), float64(*numWinners-1))), 1)
		}
		result += lineResult
	}
	return result
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	// ctx = log.Logger.WithContext(ctx)
	// l := log.Ctx(ctx).With().Logger()

	_, _, result := winnersNum(ctx, 0, input)
	return result
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

func winnersNum(ctx context.Context, curIdx int, lines []string) (int, int, int) {
	l := log.Ctx(ctx).With().Logger()

	if curIdx == len(lines) {
		return 0, 0, 0
	}

	haveNums, winningNums := separateNums(lines[curIdx])
	var numWinners *int
	for k := range haveNums {
		if _, ok := winningNums[k]; ok {
			numWinners = kit.Ptr(kit.Deref(numWinners) + 1)
			// l.Info().Int("winning_num", k).Msg("")
		}
	}
	var lineResult int
	if numWinners != nil {
		lineResult = max(int(math.Pow(float64(2), float64(*numWinners-1))), 1)
	}

	// TODO: left off here, log some stuff
	numCopies := kit.Deref(numWinners)
	resultWithCopies := lineResult
	for i := 0; i <= numCopies; i += 1 {
		nthLineResult, _, _ := winnersNum(ctx, curIdx+i, lines)
		resultWithCopies += nthLineResult
	}
	_, _, nextTotal := winnersNum(ctx, curIdx+1, lines)
	return lineResult, resultWithCopies, resultWithCopies + nextTotal
}
