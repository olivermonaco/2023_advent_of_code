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
	cMap := make([]int, 0, len(rowInfo[1]))
	for _, contig := range rowInfo[1] {
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

func calcSpringLocCombos(s string, contiguousKeys []int) [][2]int {
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
	return initialIdxs
}

func identifyNumConsecutiveBrokenSprings(s string) [][2]int {
	var consecutiveBroken [][2]int
	i := 0
	j := 0
	for {
		if j == len(s) && j > i {
			consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
			break
		}
		if string([]rune(s)[j]) != brokenSpring {
			if j > i {
				consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
			}
			j++
			i = j
			continue
		}
		if j >= len(s) {
			break
		}
		j++
	}
	return consecutiveBroken
}
