package hot_springs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

const (
	brokenSpring   = '#'
	possibleSpring = '?'
	nonSpring      = '.'
)

type stringGroups []stringGroup

type stringGroup []separatedString

type separatedString struct {
	s               []rune
	validConsecKeys []int
}

type row struct {
	sGs             stringGroups
	consecutiveKeys []int
}

func CalculatePartOne(ctx context.Context, input []string) int {
	var total int
	rows := make([]row, 0, len(input))
	for i, line := range input {
		l := log.Ctx(ctx).With().Int("line num", i).Str("line_val", line).Logger()
		ctx = l.WithContext(ctx)
		rows = append(rows, parseLine(ctx, line))
		// l := log.Ctx(ctx).With().Logger()
		l.Info().Msg("tests")

		r := parseLine(ctx, line)

		rowCombos := r.calcSpringLocCombos(ctx)
		total += rowCombos
	}
	return total
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func parseLine(ctx context.Context, line string) row {
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

	strGroups := createStringGroups(rowInfo[0])
	strGroups = strGroups.addKeys(consecNums)

	return row{
		sGs:             strGroups,
		consecutiveKeys: consecNums,
	}
}

func createStringGroups(s string) stringGroups {
	separated := strings.Split(s, ".")

	sGs := make(stringGroups, 0, len(separated))

	for _, sG := range separated {
		if len(sG) > 0 {
			sGs = append(sGs,
				stringGroup{
					separatedString{s: []rune(sG)},
				},
			)
		}
	}
	return sGs
}

func (sGs stringGroups) addKeys(keys []int) stringGroups {
	// todo: do this with pointers
	var keysIdx int

	for i, sG := range sGs {

		sG, keysIdx = sG.addKeys(keys, keysIdx)
		if len(sG) == 0 {
			continue
		}
		sGs[i] = sG
	}
	return sGs
}

func (sG stringGroup) addKeys(keys []int, curKeyIdx int) (stringGroup, int) {
	var newSG stringGroup
	for _, sepString := range sG {
		newSepStrings, newCurKeyIdx := sepString.addKeys(keys, curKeyIdx)
		newSG = append(newSG, newSepStrings...)
		curKeyIdx = newCurKeyIdx
	}
	return newSG, curKeyIdx
}

// no len means no character spans found
func findSpansOfCharacters(s []rune, curIdx int, findRune rune) [][2]int {
	var (
		start *int
		spans [][2]int
	)

	j := curIdx
	for {
		if j == len(s) {
			if start == nil {
				break
			}
			spans = append(spans, [2]int{*start, j - 1})
			break
		}
		if s[j] != findRune && start != nil {
			spans = append(spans, [2]int{*start, j - 1})
			start = nil
		}
		if s[j] == findRune && start == nil {
			tempJ := j
			start = &tempJ
		}
		j++
	}
	return spans
}

// just to make it easier to read
func remainingKeys(keys []int, curIdx int) int {
	return len(keys) - curIdx
}

func (sepString separatedString) addKeys(keys []int, curKeyIdx int) (stringGroup, int) {
	var (
		i, j             int
		seenBrokenSpring bool
	)

	consecBrokenSprings := remainingConsecutiveSprings(sepString.s)

	var newSepStrings []separatedString

	for {
		if j >= len(sepString.s) || curKeyIdx == len(keys) {
			break
		}
		if sepString.s[j] == brokenSpring {
			seenBrokenSpring = true
		}
		rKeys := remainingKeys(keys, curKeyIdx)

		if j-i+1 >= keys[curKeyIdx] && rKeys <= len(consecBrokenSprings) {
			// need to catch the idx up with the next consecutive number of broken springs
			brokenSpringSpans := findSpansOfCharacters(sepString.s, j, brokenSpring)
			j = brokenSpringSpans[0][1]
			if j-i >= keys[curKeyIdx] {
				newKeyIdx := catchKeysUp(j-i, curKeyIdx, keys)
				newSepStrings = append(newSepStrings,
					separatedString{
						s:               sepString.s[i : j+1],
						validConsecKeys: keys[curKeyIdx : newKeyIdx+1],
					},
				)
				curKeyIdx = newKeyIdx + 1
				j += 2
				i = j
				seenBrokenSpring = false
				continue
			}
			// if not, keep going
		}

		if j-i+1 >= keys[curKeyIdx] &&
			remainingKeys(keys, curKeyIdx) > len(consecBrokenSprings) &&
			!seenBrokenSpring {
			// eg. ????## 2, 3
			// find the end of the current ?? run
			possibleSpringSpans := findSpansOfCharacters(sepString.s, j, possibleSpring)
			j = possibleSpringSpans[0][1]
			// find the most number of keys that can fit in it
			newKeyIdx := catchKeysUp(j-i+1, curKeyIdx, keys)
			newSepStrings = append(newSepStrings,
				separatedString{
					s:               sepString.s[i : j+1],
					validConsecKeys: keys[curKeyIdx : newKeyIdx+1],
				},
			)
			curKeyIdx = newKeyIdx + 1
			j += 2
			i = j
			continue
		}

		if j-i+1 == keys[curKeyIdx] {
			newSepStrings = append(newSepStrings,
				separatedString{
					s:               sepString.s[i : j+1],
					validConsecKeys: []int{keys[curKeyIdx]},
				},
			)
			curKeyIdx++
			j += 2
			i = j
			seenBrokenSpring = false
			continue
		}
		j++
	}

	return newSepStrings, curKeyIdx
}

// func separateConsecutiveStrings(s string) stringGroups {
// 	var (
// 		fullRow              stringGroups // eg. [['#','#','#','?','?','#','#'], ['?', '#', '#', '?']]
// 		curConsecutiveString []rune
// 		i                    int
// 	)
// 	for {
// 		if i == len(s) {
// 			if len(curConsecutiveString) > 0 {
// 				fullRow = append(
// 					fullRow,
// 					stringGroup{
// 						separatedString{s: curConsecutiveString},
// 					},
// 				)

// 			}
// 			break
// 		}
// 		if []rune(s)[i] == nonSpring {
// 			if len(curConsecutiveString) > 0 {
// 				fullRow = append(
// 					fullRow,
// 					stringGroup{
// 						separatedString{s: curConsecutiveString},
// 					},
// 				)
// 				curConsecutiveString = []rune{}
// 			}
// 			i++
// 			continue
// 		}
// 		curConsecutiveString = append(curConsecutiveString, []rune(s)[i])
// 		i++
// 	}

// 	return fullRow
// }

func (sepString separatedString) knownBrokenSprings() []int {
	var knownBrokenSprings []int
	for i := 0; i < len(sepString.s); i++ {
		if sepString.s[i] == brokenSpring {
			knownBrokenSprings = append(knownBrokenSprings, i)
		}
	}
	return knownBrokenSprings
}

func (sG stringGroup) knownBrokenSepStrings() []int {
	var brokenSpringsIdxs []int
	for i, sepString := range sG {
		if len(sepString.knownBrokenSprings()) > 0 {
			brokenSpringsIdxs = append(brokenSpringsIdxs, i)
		}
	}
	return brokenSpringsIdxs
}

func remainingConsecutiveSprings(s []rune) [][2]int {
	var (
		i, j                 int
		remainingConsecutive [][2]int
	)

	for {
		if j == len(s) {
			if i < j {
				remainingConsecutive = append(remainingConsecutive, [2]int{i, j - 1})
			}
			break
		}
		if s[j] == brokenSpring {
			j++
			continue
		}
		if i < j {
			remainingConsecutive = append(remainingConsecutive, [2]int{i, j - 1})
		}

		j++
		i = j
	}

	return remainingConsecutive
}

func catchKeysUp(span, curKeyIdx int, consecutiveKeys []int) int {
	keysTotal := consecutiveKeys[curKeyIdx]
	for {
		if span < keysTotal {
			break
		}
		curKeyIdx++

		if curKeyIdx == len(consecutiveKeys) {
			break
		}
		keysTotal += consecutiveKeys[curKeyIdx] + 1
	}
	return curKeyIdx - 1
}

func (r row) calcSpringLocCombos(ctx context.Context) int {
	total := r.sGs.calcRecursiveTotal(ctx)
	return total
}

// // invariant: after constructing the stringGroups,
// // if a sepString.s has a brokenSpring,
// // the end of the sepString.s will always have the (consecutive) brokenSpring(s)
// // EXCEPT for the totaling done in stringGroups.calcRecursiveTotal()
// func (r row) calcSpringLocCombos() int {
// 	var (
// 		sGIdx, curKeyIdx int
// 	)
// 	for {
// 		if sGIdx == len(r.sGs) {
// 			break
// 		}
// 		sG := r.sGs[sGIdx]
// 		var (
// 			sepStringsIdx  int
// 			newStringGroup stringGroup
// 		)
// 		for {
// 			if sepStringsIdx == len(sG) {
// 				r.sGs[sGIdx] = newStringGroup
// 				sGIdx++
// 				break
// 			}
// 			sepString := sG[sepStringsIdx]
// 			consecutiveBrokenStringsLeft := remainingConsecutiveSprings(sepString.s)

// 			var i, j int
// 			tempSepString := separatedString{}
// 			for {
// 				if j >= len(sepString.s) {
// 					if i < j && curKeyIdx < len(r.consecutiveKeys) {
// 						newKeyIdx := catchKeysUp(j-i, curKeyIdx, r.consecutiveKeys)
// 						tempSepString.validConsecKeys = r.consecutiveKeys[curKeyIdx : newKeyIdx+1]
// 						tempSepString.s = append(tempSepString.s, sepString.s[i:j]...)
// 						newStringGroup = append(newStringGroup, tempSepString)
// 						curKeyIdx = newKeyIdx + 1
// 					}
// 					break
// 				}
// 				keysLeft := len(r.consecutiveKeys) - curKeyIdx
// 				if keysLeft <= len(consecutiveBrokenStringsLeft) &&
// 					len(consecutiveBrokenStringsLeft) > 0 {
// 					var bSIdx int
// 					for {
// 						if bSIdx == len(consecutiveBrokenStringsLeft) {
// 							bSIdx--
// 							break
// 						}
// 						if j+1 <= consecutiveBrokenStringsLeft[bSIdx][0] {
// 							bSIdx = bSIdx - 1
// 							break
// 						}
// 						bSIdx++
// 					}
// 					bSIdx = max(bSIdx, 0)
// 					j = consecutiveBrokenStringsLeft[bSIdx][1]
// 					for {
// 						if j-i+1 >= r.consecutiveKeys[curKeyIdx] {
// 							break
// 						}
// 						j++
// 					}

// 					tempSepString.validConsecKeys = append(
// 						tempSepString.validConsecKeys,
// 						r.consecutiveKeys[curKeyIdx],
// 					)
// 					tempSepString.s = sepString.s[i : j+1]
// 					newStringGroup = append(newStringGroup, tempSepString)
// 					j += 2
// 					i = j
// 					consecutiveBrokenStringsLeft = consecutiveBrokenStringsLeft[bSIdx+1:]
// 					curKeyIdx++
// 					tempSepString = separatedString{}
// 					continue
// 				}
// 				if j+2 < len(sepString.s) &&
// 					sepString.s[j+1] != brokenSpring &&
// 					sepString.s[j+2] == brokenSpring {
// 					newKeyIdx := catchKeysUp(j-i, curKeyIdx, r.consecutiveKeys)
// 					tempSepString.validConsecKeys = r.consecutiveKeys[curKeyIdx : newKeyIdx+1]
// 					tempSepString.s = sepString.s[i:j]

// 					newStringGroup = append(newStringGroup, tempSepString)
// 					tempSepString = separatedString{}
// 					curKeyIdx = newKeyIdx + 1
// 					j += 2
// 					i = j
// 					continue
// 				}
// 				j++
// 			}
// 			sepStringsIdx++
// 		}
// 	}
// 	total := r.sGs.calcRecursiveTotal()
// 	return total
// }

func (sGs stringGroups) calcRecursiveTotal(ctx context.Context) int {
	var combosToMult []int
	for _, sG := range sGs {
		combosToMult = append(
			combosToMult,
			sG.calcRecursiveTotal(ctx),
		)
		combosToMult = append(
			combosToMult,
			sG.shiftAndReturnCombos(ctx)...,
		)
	}
	total := kit.Mult(combosToMult)

	return total
}

func (sG stringGroup) calcRecursiveTotal(ctx context.Context) int {
	total := 1
	for _, sepString := range sG {
		total *= sepString.calcRecursiveTotal(ctx)
	}

	return total
}

func (sepString separatedString) calcRecursiveTotal(ctx context.Context) int {
	totals := []int{1}
	if len(sepString.knownBrokenSprings()) == 0 {
		totals = append(totals,
			calcRecursiveSepStringTotal(
				ctx,
				sepString.s, 0, 0, sepString.validConsecKeys,
			),
		)
	}

	return kit.Mult(totals)
}

func calcRecursiveSepStringTotal(ctx context.Context, runes []rune, spansIdx, offset int, spans []int) int {
	var total int
	prevSpansSum := kit.Sum(spans[:spansIdx])
	prevSpansSum += len(spans[:spansIdx]) // add the buffer characters
	if spansIdx+1 > len(spans) {
		fmt.Println("test")
	}
	curSpan := spans[spansIdx]

	if spansIdx == len(spans)-1 {
		var combos int
		curEnd := prevSpansSum + curSpan + offset
		for i := 0; i+curEnd < len(runes)+1; i++ {
			combos++
		}
		return combos
	}
	i := 0
	for {
		subtotal := calcRecursiveSepStringTotal(ctx, runes, spansIdx+1, offset+i, spans)
		if subtotal == 0 {
			break
		}
		total += subtotal
		i++
	}
	return total
}

// start at the right most brokenSpring end, and work backward fulfilling the keys
func findRemainingLeftMostRunes(sepString separatedString) int {
	brokenSpringSpans := findSpansOfCharacters(sepString.s, 0, brokenSpring)
	if len(brokenSpringSpans) == 0 {
		// order doesn't matter, math it
		sumNums := kit.Sum(sepString.validConsecKeys)
		lenNums := len(sepString.validConsecKeys) - 1
		return len(sepString.s) - sumNums - lenNums
	}
	var i, j, remainingLeftMost int
	curKeyIdx := len(sepString.validConsecKeys) - 1
	brokenSpringSpansIdx := len(brokenSpringSpans) - 1
	curRunes := sepString.s
	for {
		if brokenSpringSpansIdx < 0 || curKeyIdx < 0 {
			sumNums := kit.Sum(sepString.validConsecKeys[:curKeyIdx+1])
			lenNums := len(sepString.validConsecKeys[:curKeyIdx+1])
			remainingLeftMost = (j - i + 1) - sumNums - lenNums
			break
		}
		i = brokenSpringSpans[brokenSpringSpansIdx][0]
		j = brokenSpringSpans[brokenSpringSpansIdx][1]
		curSpan := j - i + 1
		curKey := sepString.validConsecKeys[curKeyIdx]
		if curSpan < curKey {
			// increase j until the end of the slice, or it fulfills the curKey
			for {
				if j == len(curRunes) || j-i == sepString.validConsecKeys[curKeyIdx] {
					break
				}
				j++
			}
			for {
				curSpan = j - i
				if i == 0 || curSpan == sepString.validConsecKeys[curKeyIdx] {
					break
				}
				i--
			}
		}
		j = i - 2 // 2 for buffer
		i = 0     // will get increased if there's any remaining brokenSprings
		curRunes = curRunes[:j+1]
		curKeyIdx--
		brokenSpringSpansIdx--
	}
	return remainingLeftMost
}

// shiftAndReturnCombos shifts the sepStrings in the string group rune by rune,
// to provide the number of possible combos per stringGroup possibility
// eg. [("??", 1), ("??##", 4), ("????", 2)] having possibilities of:
//
//	[("???", 1), ("?##?", 4), ("???", 2)] -> (3*1*2) = 6 possibilities
//	[("????", 1), ("##??", 4), ("??", 2)] -> (4*1*1) = 4 possibilities
//
// returns []int{6, 4}
func (sG stringGroup) shiftAndReturnCombos(ctx context.Context) []int {
	var combos []int

	idxsBrokenSprings := sG.knownBrokenSepStrings()
	for i := 0; i < len(idxsBrokenSprings); i++ {
		// copy the stringGroup
		tempSG := stringGroup(kit.Map(sG, func(sepStr separatedString) separatedString { return sepStr }))
		if idxsBrokenSprings[i] < len(sG)-1 {
			curIdx := idxsBrokenSprings[i]
			if kit.Sum(tempSG[curIdx].validConsecKeys) == len(tempSG[curIdx].knownBrokenSprings()) {
				// can't advance the broken springs, because all keys match the number of broken springs
				// eg. ??### 3
				break
			}
			numFreeRunes := findRemainingLeftMostRunes(tempSG[curIdx+1])

			for j := 0; j < numFreeRunes+1; j++ {

				if tempSG[curIdx].s[0] == brokenSpring ||
					tempSG[curIdx+1].s[0] == brokenSpring {
					// stop at the first brokenSpring in the current sepString or next sepString
					break
				}
				if idxsBrokenSprings[i] > 0 {
					// pass back any ??'s to the previous
					tempSG[i-1].s = append(tempSG[curIdx-1].s, tempSG[curIdx-1].s[:j+1]...)
				}
				// take up to j from the next sepString
				tempSG[curIdx].s = append(tempSG[curIdx].s, tempSG[curIdx+1].s[:1]...)
				tempSG[curIdx+1].s = tempSG[curIdx+1].s[1:]
				// remove up to j from current sepString
				tempSG[curIdx].s = tempSG[curIdx].s[1:]

				combos = append(combos, tempSG.calcRecursiveTotal(ctx))
			}
		}
	}
	return combos
}

func matchingExpected(s, matchingStr string) bool {
	var matchingIdxs int
	for idx, r := range s {
		if idx > 7 {
			continue
		}
		if r == []rune(matchingStr)[idx] {
			matchingIdxs++
		}
	}
	return matchingIdxs == len(matchingStr)
}
