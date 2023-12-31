package hot_springs

import (
	"context"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
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
	// rows := make([]row, 0, len(input))
	for _, line := range input {
		r := parseLine(line)

		rowCombos := r.calcSpringLocCombos()
		total += rowCombos
	}
	return total
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func parseLine(line string) row {
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

	return row{
		sGs:             []stringGroup{separateConsecutiveStrings(rowInfo[0])},
		consecutiveKeys: consecNums,
	}
}

// func matchRangeValidItems(validSpringsRunes []rune, consecNums []int, consecNumsIdx int) []int {
// 	var (
// 		runeIdxStart, runeIdxEnd int
// 		matchingconsecNums       []int
// 	)

// 	for {
// 		// catch window up with possible range
// 		if consecNumsIdx == len(consecNums) || runeIdxEnd > len(validSpringsRunes) {
// 			break
// 		}
// 		if runeIdxEnd-runeIdxStart+1 < consecNums[consecNumsIdx] {
// 			runeIdxEnd++
// 			continue
// 		}
// 		runeIdxEnd = shiftEndForKnownBrokenSpring(
// 			runeIdxStart,
// 			runeIdxEnd,
// 			validSpringsRunes,
// 		)

// 		matchingconsecNums = append(
// 			matchingconsecNums,
// 			consecNums[consecNumsIdx],
// 		)
// 		if string(validSpringsRunes[runeIdxEnd]) == possibleSpring &&
// 			runeIdxEnd+2 < len(validSpringsRunes) {
// 			runeIdxEnd++
// 		}
// 		runeIdxStart = runeIdxEnd + 1
// 		runeIdxEnd = runeIdxStart
// 		consecNumsIdx++
// 	}
// 	return matchingconsecNums
// }

// func shiftEndForKnownBrokenSpring(
// 	runeIdxStart, runeIdxEnd int,
// 	validSpringsRunes []rune,
// ) int {
// 	for {
// 		if runeIdxEnd < len(validSpringsRunes)-1 {
// 			if string(validSpringsRunes[runeIdxEnd+1]) == brokenSpring {
// 				runeIdxEnd++
// 				continue
// 			}
// 		}
// 		break
// 	}
// 	return runeIdxEnd
// }

func separateConsecutiveStrings(s string) []separatedString {
	var (
		fullRow              []separatedString // eg. [['#','#','#','?','?','#','#'], ['?', '#', '#', '?']]
		sG                   separatedString   // eg.  ['#','#','#','?','?','#','#']
		curConsecutiveString []rune
		i                    int
	)
	for {
		if i == len(s) {
			break
		}
		if []rune(s)[i] == nonSpring {
			if len(curConsecutiveString) > 0 {
				fullRow = append(
					fullRow,
					separatedString{s: curConsecutiveString},
				)
				fullRow = append(fullRow, sG)
				curConsecutiveString = []rune{}
			}
			i++
			continue
		}
		curConsecutiveString = append(curConsecutiveString, []rune(s)[i])
		i++
	}

	return fullRow
}

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
	var keysTotal int
	for {
		if span > keysTotal || curKeyIdx == len(consecutiveKeys) {
			break
		}
		keysTotal += consecutiveKeys[curKeyIdx]
		curKeyIdx++
	}
	return curKeyIdx - 1
}

// invariant: after constructing the stringGroups,
// if a sepString.s has a brokenSpring,
// the end of the sepString.s will always have the (consecutive) brokenSpring(s)
// EXCEPT for the totaling done in stringGroups.calcRecursiveTotal()
func (r row) calcSpringLocCombos() int {
	var (
		sGIdx, curKeyIdx int
	)
	for {
		if sGIdx == len(r.sGs) {
			break
		}
		sG := r.sGs[sGIdx]
		var (
			sepStringsIdx  int
			newStringGroup stringGroup
		)
		for {
			if sepStringsIdx == len(sG) {
				r.sGs[sGIdx] = newStringGroup
				sGIdx++
				break
			}
			sepString := sG[sepStringsIdx]
			consecutiveBrokenStringsLeft := remainingConsecutiveSprings(sepString.s)

			var i, j int
			tempSepString := separatedString{}
			for {
				if j >= len(sepString.s) {
					if i < j && curKeyIdx < len(r.consecutiveKeys) {
						newKeyIdx := catchKeysUp(j-i-1, curKeyIdx, r.consecutiveKeys[curKeyIdx:])
						newKeyIdx = max(newKeyIdx, 0)
						tempSepString.validConsecKeys = r.consecutiveKeys[curKeyIdx:][newKeyIdx:]
						tempSepString.s = append(tempSepString.s, sepString.s[i:j]...)
						newStringGroup = append(newStringGroup, tempSepString)
					}
					break
				}
				keysLeft := len(r.consecutiveKeys) - curKeyIdx - 1
				if keysLeft <= len(consecutiveBrokenStringsLeft) && len(consecutiveBrokenStringsLeft) > 0 {
					var bSIdx int
					for {
						if bSIdx == len(consecutiveBrokenStringsLeft) {
							bSIdx--
							break
						}
						if j+1 <= consecutiveBrokenStringsLeft[bSIdx][0] {
							bSIdx = bSIdx - 1
							break
						}
						bSIdx++
					}
					bSIdx = max(bSIdx, 0)
					j = consecutiveBrokenStringsLeft[bSIdx][1]
					for {
						if j-i+1 >= r.consecutiveKeys[curKeyIdx] {
							break
						}
						j++
					}

					tempSepString.validConsecKeys = append(
						tempSepString.validConsecKeys,
						r.consecutiveKeys[curKeyIdx],
					)
					tempSepString.s = sepString.s[i : j+1]
					newStringGroup = append(newStringGroup, tempSepString)
					j += 2
					i = j
					consecutiveBrokenStringsLeft = consecutiveBrokenStringsLeft[bSIdx+1:]
					curKeyIdx++
					tempSepString = separatedString{}
					continue
				}
				if j+2 < len(sepString.s) && sepString.s[j+2] == brokenSpring {
					newKeyIdx := catchKeysUp(j-i+1, curKeyIdx, r.consecutiveKeys[curKeyIdx:])
					tempSepString.validConsecKeys = r.consecutiveKeys[curKeyIdx : newKeyIdx+1]
					tempSepString.s = sepString.s[i:j]

					newStringGroup = append(newStringGroup, tempSepString)
					tempSepString = separatedString{}
					curKeyIdx = newKeyIdx + 1
					j += 2
					i = j
					continue
				}
				j++
			}
			sepStringsIdx++
		}
	}
	total := r.sGs.calcRecursiveTotal()
	return total
}

func (sGs stringGroups) calcRecursiveTotal() int {
	var combosToSum []int
	for _, sG := range sGs {
		combosToSum = append(
			combosToSum,
			sG.calcRecursiveTotal(),
		)
		combosToSum = append(
			combosToSum,
			sG.shiftAndReturnCombos()...,
		)
	}
	total := kit.Sum(combosToSum)

	return total
}

func (sG stringGroup) calcRecursiveTotal() int {
	total := 1
	for _, sepString := range sG {
		total *= sepString.calcRecursiveTotal()
	}

	return total
}

func (sepString separatedString) calcRecursiveTotal() int {
	totals := []int{1}
	if len(sepString.knownBrokenSprings()) == 0 {
		totals = append(totals,
			calcRecursiveSepStringTotal(
				sepString.s, 0, 0, sepString.validConsecKeys,
			),
		)
	}

	return kit.Mult(totals)
}

func calcRecursiveSepStringTotal(runes []rune, spansIdx, offset int, spans []int) int {
	var total int
	prevSpansSum := kit.Sum(spans[:spansIdx])
	prevSpansSum += len(spans[:spansIdx]) // add the buffer characters
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
		subtotal := calcRecursiveSepStringTotal(runes, spansIdx+1, offset+i, spans)
		if subtotal == 0 {
			break
		}
		total += subtotal
		i++
	}
	return total
}

// shiftAndReturnCombos shifts the sepStrings in the string group rune by rune,
// to provide the number of possible combos per stringGroup possibility
// eg. [("??", 1), ("??##", 4), ("????", 2)] having possibilities of:
//
//	[("???", 1), ("?##?", 4), ("???", 2)] -> (3*1*2) = 6 possibilities
//	[("????", 1), ("##??", 4), ("??", 2)] -> (4*1*1) = 4 possibilities
//
// returns []int{6, 4}
func (sG stringGroup) shiftAndReturnCombos() []int {
	var combos []int

	idxsBrokenSprings := sG.knownBrokenSepStrings()
	for i := 0; i < len(idxsBrokenSprings); i++ {
		// copy the stringGroup
		tempSG := stringGroup(kit.Map(sG, func(sepStr separatedString) separatedString { return sepStr }))
		if idxsBrokenSprings[i] < len(sG)-1 {
			curIdx := idxsBrokenSprings[i]
			nextBrokenSprings := tempSG[curIdx+1].knownBrokenSprings()
			numFreeChars := len(tempSG[curIdx+1].s)
			if len(nextBrokenSprings) > 0 {
				numFreeChars = nextBrokenSprings[0]
			}
			numFreeChars = min(numFreeChars, kit.Sum(tempSG[curIdx].validConsecKeys)-1)
			for j := 0; j < numFreeChars; j++ {
				// as long as there's extra ?? in the next sepString,
				// add them to the end of the current sepString (containing brokenSprings)
				// and, if there's a previous sepString, add it to the prior sepString

				if tempSG[curIdx].s[0] == brokenSpring ||
					tempSG[curIdx+1].s[0] == brokenSpring {
					// stop at the first brokenSpring in the current sepString or next sepString
					break
				}
				if idxsBrokenSprings[i] > 0 {
					tempSG[i-1].s = append(tempSG[i-1].s, tempSG[i-1].s[:j+1]...)
				}
				// take up to j from the next sepString
				tempSG[curIdx].s = append(tempSG[curIdx].s, tempSG[curIdx+1].s[:1]...)
				tempSG[curIdx+1].s = tempSG[curIdx+1].s[1:]
				// remove up to j from current sepString
				tempSG[curIdx].s = tempSG[curIdx].s[1:]

				combos = append(combos, tempSG.calcRecursiveTotal())
			}
		}
	}
	return combos
}
