package hot_springs

import (
	"context"
	"strconv"
	"strings"
)

const (
	brokenSpring   = "#"
	possibleSpring = "#"
	nonSpring      = "."
)

type row struct {
	contiguousKeys                 []int
	springs                        string
	possibleRanges, definiteRanges [][2]int
}

type rangeAndLimit struct {
	low, high int
	limit     *int
}

func CalculatePartOne(ctx context.Context, input []string) int {
	rows := make([]row, 0, len(input))
	for _, line := range input {
		rows = append(rows, parseLine(line))
	}
	return 0
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func parseLine(line string) row {
	rowInfo := strings.Fields(line)
	if len(rowInfo) != 2 {
		panic(rowInfo)
	}
	contigNums := strings.Split(rowInfo[1], ",")

	cMap := make([]int, 0, len(contigNums))
	for _, contig := range contigNums {
		n, err := strconv.Atoi(string(contig))
		if err != nil {
			panic(err)
		}
		cMap = append(cMap, n)
	}

	return row{
		contiguousKeys: cMap,
		springs:        rowInfo[0],
	}
}

func (r row) findCombinations() int {
	var (
		contiguousKeyIdx, i int
	)

	for j := 0; j < len([]rune(r.springs)); j++ {
		if string([]rune(r.springs)[j]) == nonSpring {
			if i == j {
				// two nonSprings in succession
				i++
				continue
			}
			strLen := j - i + 1
			var checkKeys []int
			for {
				if r.contiguousKeys[contiguousKeyIdx] > strLen {
					break
				}
				strLen -= r.contiguousKeys[contiguousKeyIdx]
				checkKeys = append(checkKeys, contiguousKeyIdx)
				contiguousKeyIdx++
			}

		}
	}
	return 0
}

func calcSpringLocCombos(s string, contiguousKeys []int) int {
	var (
		i, j, checkKeysIdx, brokenSpringsIdx int
		initialIdxs                          []rangeAndLimit
	)

	knownBrokenSprings := identifyNumConsecutiveBrokenSprings(s)

	for {
		if checkKeysIdx == len(contiguousKeys) {
			break
		}

		if j-i+1 < contiguousKeys[checkKeysIdx] {
			j++
			continue
		}
		if len(knownBrokenSprings) > 0 {
			for {
				// catch knownBrokenSprings up with the current idx
				if knownBrokenSprings[brokenSpringsIdx][1] < i {
					brokenSpringsIdx++
				}
				break
			}
		}

		lenSpringsMinusIdx := len(knownBrokenSprings) - brokenSpringsIdx
		lenKeysMinusIdx := len(contiguousKeys) - checkKeysIdx

		if lenSpringsMinusIdx >= lenKeysMinusIdx && j < knownBrokenSprings[brokenSpringsIdx][0] {
			i++
			j++
			continue
		}

		rAL := rangeAndLimit{low: i, high: j}
		if brokenSpringsIdx < len(knownBrokenSprings) {
			if knownBrokenSprings[brokenSpringsIdx][1] >= i {
				upperLimitKnownBrokenSpring(&rAL, brokenSpringsIdx, knownBrokenSprings)
			}
		}
		initialIdxs = append(initialIdxs, rAL)
		j += 2
		i = j

		checkKeysIdx++
	}

	combos := calcRecursiveRangeTotal(
		initialIdxs,
		0,
		0,
		len(s),
	)
	return combos
}

func upperLimitKnownBrokenSpring(
	rAL *rangeAndLimit,
	knownBrokenSpringsIdx int,
	knownBrokenSprings [][2]int,
) {
	if len(knownBrokenSprings) == 0 {
		return
	}
	if knownBrokenSprings[knownBrokenSpringsIdx][0]-2 < rAL.high+(rAL.high-rAL.low) {
		// min of next known broken spring minus 2,
		// or lowest known broken spring + # of contiguous springs

		limit := knownBrokenSprings[knownBrokenSpringsIdx][0] + rAL.high - rAL.low
		if knownBrokenSpringsIdx+1 <= len(knownBrokenSprings)-1 {
			limit = min(
				limit,
				knownBrokenSprings[knownBrokenSpringsIdx+1][0]-2,
			)
		}
		rAL.limit = &limit
	}

}

func identifyNumConsecutiveBrokenSprings(s string) [][2]int {
	var consecutiveBroken [][2]int
	i := 0
	j := 0
	for {
		if j == len(s) {
			if string([]rune(s)[j-1]) == brokenSpring {
				consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
			}
			break
		}
		if string([]rune(s)[j]) == brokenSpring {
			j++
			continue
		}
		if j > i {
			consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
		}
		j++
		i = j
	}
	return consecutiveBroken
}

func calculateCombos(numKeys, latestSpringsIdx, latestNumSprings, numSpaces, diffBetweenSpaces int) int {
	var total int
	for {
		for curStart := numSpaces - latestSpringsIdx; curStart > -1; curStart-- {
			total += calcSeries(curStart, curStart, diffBetweenSpaces)
		}
		if latestSpringsIdx+latestNumSprings == numSpaces {
			break
		}
		latestSpringsIdx++
	}

	return total
}

// func calcNumSpaces(ranges [][2]int, curIdx, offset int) int {
// 	for _, r := range ranges[curIdx:] {

// 	}
// }

// TODO: test this and make it better
func calcRecursiveRangeTotal(ranges []rangeAndLimit, curIdx, prevOffset, maxNum int) int {

	var prevHigh, curOffset int

	curMax := maxNum

	if curIdx > 0 {
		prevHigh = ranges[curIdx-1].high
		for {
			lowerLimit := ranges[curIdx].low - 2 + curOffset
			if prevHigh+prevOffset > lowerLimit {
				curOffset++
				continue
			}
			break
		}
	}

	var total int
	if ranges[curIdx].limit != nil {
		curMax = *ranges[curIdx].limit + 1
	}

	if curIdx == len(ranges)-1 {
		for {
			lastStartOffset := ranges[curIdx].low + curOffset
			lenOfRange := ranges[curIdx].high - ranges[curIdx].low + 1
			compareTotal := lastStartOffset + lenOfRange
			if compareTotal > curMax {
				break
			}
			total++
			curOffset++
		}
		return total
	}

	for {
		subTotal := calcRecursiveRangeTotal(ranges, curIdx+1, curOffset, maxNum)
		if subTotal == 0 {
			break
		}
		total += subTotal
		curOffset++
		if ranges[curIdx].limit != nil {
			if ranges[curIdx].high+curOffset > *ranges[curIdx].limit {
				break
			}
		}

	}
	return total
}

func calcSeries(numterms, firstTerm, diff int) int {
	total := (2 * firstTerm)
	total += ((numterms - 1) * diff)
	total *= numterms
	total /= 2
	return total
}
