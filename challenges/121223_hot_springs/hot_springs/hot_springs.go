package hot_springs

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

func CalculatePartOne(ctx context.Context, input []string) int {
	var total int

	logRows := make([][]string, len(input)+1)
	logRows[0] = []string{"line_value", "num_combos"}

	for i, line := range input {
		l := log.Ctx(ctx).With().
			Int("line num", i).
			Str("line_val", line).
			Logger()

		l.Info().Msg("starting")
		sGsKeys := parseLine(ctx, line)

		rowCombos := sGsKeys.calcTotal1(ctx)
		l.Info().Int("num_combos", rowCombos).Send()
		logRows[i+1] = []string{strings.ReplaceAll(line, ",", "_"), fmt.Sprintf("%d", rowCombos)}

		total += rowCombos
	}
	writeToCSV(logRows, openCSVFile("results"))
	return total
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	var total int

	logRows := make([][]string, len(input)+1)
	logRows[0] = []string{"line_value", "num_combos"}

	for i, line := range input {
		l := log.Ctx(ctx).With().
			Int("line num", i).
			Str("line_val", line).
			Logger()

		l.Info().Msg("starting")
		sGsKeys := parseLine2(ctx, line)

		rowCombos := sGsKeys.calcTotal2(ctx)
		l.Info().Int("num_combos", rowCombos).Send()
		logRows[i+1] = []string{strings.ReplaceAll(line, ",", "_"), fmt.Sprintf("%d", rowCombos)}

		total += rowCombos
	}
	writeToCSV(logRows, openCSVFile("results"))
	return total
}

// generates numbers for the cache
func generateNum(numRBs, greatestNum int) {
	var result int
	if numRBs == 0 && len(resultsCache) == 0 {
		resultsCache = append(resultsCache, []int{})
	}
	if numRBs > len(resultsCache)-1 {
		generateNum(numRBs-1, greatestNum)
		resultsCache = append(resultsCache, []int{})
	}

	if greatestNum == 1 && len(resultsCache[numRBs]) == 0 {
		resultsCache[numRBs] = append(resultsCache[numRBs], []int{1}...)
	}

	if greatestNum-1 > len(resultsCache[numRBs])-1 {
		generateNum(numRBs, greatestNum-1)
	}

	lastResult := resultsCache[numRBs][len(resultsCache[numRBs])-1]
	curRBLen := len(resultsCache[numRBs]) - 1
	prevRBNum := 1
	if numRBs > 0 {
		if curRBLen+1 > len(resultsCache[numRBs-1])-1 {
			generateNum(numRBs-1, curRBLen+1)
		}
		prevRBNum = resultsCache[numRBs-1][curRBLen+1]
	}
	result = lastResult + prevRBNum
	resultsCache[numRBs] = append(resultsCache[numRBs], result)
}

func calcNums(numRBs, greatestNum int) int {
	if greatestNum == 0 {
		return 1
	}
	if numRBs <= len(resultsCache)-1 {
		existingCombos := resultsCache[numRBs]
		if greatestNum <= len(existingCombos)-1 {
			return existingCombos[greatestNum]
		}
	}
	generateNum(numRBs, greatestNum)

	return resultsCache[numRBs][greatestNum]
}

func bkSpansEqFunc(a, b [2]int) bool {
	if a[0] != b[0] {
		return false
	}
	if a[1] != b[1] {
		return false
	}
	return true
}

func brokenLen(brokenSpans [][2]int) int {
	var l int
	for _, span := range brokenSpans {
		l += span[1] - span[0] + 1
	}
	return l
}

func brokenSpansInRange(start, end int, allBrokenSpans [][2]int) [][2]int {
	brokenInRange := make([][2]int, 0, len(allBrokenSpans))
	for _, span := range allBrokenSpans {
		if span[0] > end {
			break
		}
		if start <= span[0] && end >= span[1] {
			brokenInRange = append(brokenInRange, span)
		}
	}
	return brokenInRange
}

func contains(refBGs []refBuffGroups, compare refBuffGroups, eqFunc func(a, b refBuffGroup) bool) bool {
	return slices.ContainsFunc(
		refBGs,
		func(rBGs refBuffGroups) bool {
			if len(rBGs) != len(compare) {
				return false
			}
			for i, rBG := range rBGs {
				if len(rBG.brokenSpans) != len(compare[i].brokenSpans) {
					return false
				}
				if !eqFunc(rBG, compare[i]) {
					return false
				}
			}
			return true
		},
	)
}

func fitAvailable(rB refBuff, brokenSpans [][2]int, n int) int {
	if rB.start+n > brokenSpans[0][0] {
		if rB.start+1 > brokenSpans[0][0] {
			return 0
		}
		return brokenSpans[0][0] - rB.start
	}
	return n
}

func parseLine(ctx context.Context, line string) stringGroupsKeys {
	rowInfo := strings.Fields(line)
	if len(rowInfo) != 2 {
		panic(rowInfo)
	}
	consecNumsStr := strings.Split(rowInfo[1], ",")

	consecNums := make([]int, 0, len(consecNumsStr))
	for _, consec := range consecNumsStr {
		n, err := strconv.Atoi(string(consec))
		if err != nil {
			panic(err)
		}
		consecNums = append(consecNums, n)
	}

	strGroups := createStringGroupsPt1(rowInfo[0], consecNums)

	return strGroups
}

func parseLine2(ctx context.Context, line string) stringGroupsKeys {
	rowInfo := strings.Fields(line)
	if len(rowInfo) != 2 {
		panic(rowInfo)
	}
	consecNumsStr := strings.Split(rowInfo[1], ",")

	consecNums := make([]int, 0, len(consecNumsStr))
	for _, consec := range consecNumsStr {
		n, err := strconv.Atoi(string(consec))
		if err != nil {
			panic(err)
		}
		consecNums = append(consecNums, n)
	}

	strGroups := createStringGroupsPt2(rowInfo[0], consecNums)

	return strGroups
}

func createStringGroupsPt1(s string, consecBrokenSprings []int) stringGroupsKeys {
	separated := strings.Split(s, ".")

	sGsKs := stringGroupsKeys{
		sGs:  make([]stringGroup, 0, len(separated)),
		keys: consecBrokenSprings,
	}

	for _, sGStr := range separated {
		if len(sGStr) > 0 {
			sG := stringGroup{
				charsBrokenSpans: charsBrokenSpans{
					chars:       []rune(sGStr),
					brokenSpans: brokenSpans([]rune(sGStr)),
				},
			}

			sGsKs.sGs = append(sGsKs.sGs, sG)
			sGsKs.brokenSprings = append(sGsKs.brokenSprings, sG.brokenSpans...)
		}
	}
	return sGsKs
}

func dupeFiveX[T any](vals []T, addVal *T) []T {
	new := make([]T, 0, (len(vals)*5)+5)
	for i := 0; i < 5; i++ {
		new = append(new, vals...)
		if addVal != nil {
			new = append(new, *addVal)
		}
	}
	if addVal != nil {
		new = new[:len(new)-1]
	}
	return new
}
func createStringGroupsPt2(s string, consecBrokenSprings []int) stringGroupsKeys {
	chars := []rune(s)
	chars = dupeFiveX(chars, kit.Ptr(possibleSpring))
	s = string(chars)

	separated := strings.Split(s, ".")
	sGsKs := stringGroupsKeys{
		sGs:  make([]stringGroup, 0, len(separated)),
		keys: dupeFiveX(consecBrokenSprings, nil),
	}

	for _, sGStr := range separated {
		if len(sGStr) > 0 {
			sG := stringGroup{
				charsBrokenSpans: charsBrokenSpans{
					chars:       []rune(sGStr),
					brokenSpans: brokenSpans([]rune(sGStr)),
				},
			}

			sGsKs.sGs = append(sGsKs.sGs, sG)
			sGsKs.brokenSprings = append(sGsKs.brokenSprings, sG.brokenSpans...)
		}
	}
	return sGsKs
}

func brokenSpans(runes []rune) [][2]int {
	var (
		i, j         int
		consecBroken [][2]int
	)

	for {
		if j == len(runes) {
			if i < j {
				consecBroken = append(consecBroken, [2]int{i, j - 1})
			}
			break
		}
		if runes[j] == brokenSpring {
			j++
			continue
		}
		if i < j {
			consecBroken = append(consecBroken, [2]int{i, j - 1})
		}

		j++
		i = j
	}
	return consecBroken
}

func nextLowest(sepStrRefs separatedStringRefs, span [2]int) int {
	for i := len(sepStrRefs) - 1; i > -1; i-- {
		if sepStrRefs[i].end < span[1] {
			return i
		}
	}
	return -1
}

func minKeysLen(keys []int) int {
	return max(kit.Sum(keys)+(len(keys)-1), 0)
}

func contractSpansToDiff(spans [][2]int, diff int) [][2]int {
	i := 0
	for {
		if i == len(spans) {
			return [][2]int{}
		}
		if spans[len(spans)-1][1]-spans[i][0] < diff {
			break
		}
		i++
	}
	return spans[i:]
}

// ------ //

func openCSVFile(fPrepend string) *os.File {
	t := time.Now()
	file, err := os.OpenFile(
		fmt.Sprintf("%s_%s.csv", fPrepend, t.Format("01-06-2006_3:04pm")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}
	return file
}

func writeToCSV(logRows [][]string, f *os.File) {
	writer := csv.NewWriter(f)
	defer writer.Flush()

	err := writer.WriteAll(logRows)
	if err != nil {
		panic(err)
	}
}
