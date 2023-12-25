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
		i, j, checkKeysIdx int
		initialIdxs        [][2]int
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
		for {
			// catch knownBrokenSprings up with the current idx
			if len(knownBrokenSprings) == 0 {
				break
			}
			if knownBrokenSprings[0][1] > j {
				break
			}
			knownBrokenSprings = knownBrokenSprings[1:]
		}

		if len(knownBrokenSprings) > len(contiguousKeys)-checkKeysIdx-1 {
			i++
			j++
			continue
		}
		initialIdxs = append(initialIdxs, [2]int{i, j})
		j += 2
		i = j

		checkKeysIdx++
	}

	// TODO: replace this, or use it if the known / unknown values can use it
	// var latestSpringIdx, latestSpringLen int
	// if len(initialIdxs) > 0 {
	// 	latestSpringIdx = initialIdxs[len(initialIdxs)-1][1]
	// 	latestSpringLen = initialIdxs[len(initialIdxs)-1][1] -
	// 		initialIdxs[len(initialIdxs)-1][0] + 1
	// }
	// combos := calculateCombos(len(contiguousKeys), latestSpringIdx, latestSpringLen, len(s), -1)
	combos := calcRecursiveRangeTotal(
		initialIdxs,
		0,
		0,
		len(s),
	)
	return combos
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
func calcRecursiveRangeTotal(ranges [][2]int, curIdx, offset, maxNum int) int {
	var total int
	lastStart := ranges[len(ranges)-1][0]
	lastEnd := ranges[len(ranges)-1][1]

	if curIdx == len(ranges)-1 {
		curStart := lastStart
		curEnd := lastEnd
		for {
			lastStartOffset := curStart + offset
			lenOfRange := curEnd - curStart + 1
			compareTotal := lastStartOffset + lenOfRange
			if compareTotal > maxNum {
				break
			}
			total++
			offset++
		}
		return total
	}
	// curStart := ranges[curIdx][0]
	// curEnd := ranges[curIdx][1]
	for {
		subTotal := calcRecursiveRangeTotal(ranges, curIdx+1, offset, maxNum)
		total += subTotal
		if subTotal == 0 {
			// this replaces the below
			break
		}
		// lastStartOffset := lastStart + offset
		// lastStartRange := lastEnd - lastStart + 1
		// otherCompareTotal := lastStartOffset + lastStartRange
		// if otherCompareTotal > maxNum {
		// 	break
		// }
		offset++
	}
	return total
}

func calcCurrentRange(curRangeStart, curRangeEnd, limit int) int {
	var total int
	for {
		currentLen := curRangeEnd - curRangeStart + 1
		if curRangeStart+currentLen == limit {
			break
		}
		total++
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
