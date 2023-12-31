package hot_springs

import (
	"context"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

const (
	brokenSpring   = "#"
	possibleSpring = "?"
	nonSpring      = "."
)

type stringGroups []stringGroup

type stringGroup []separatedString

type separatedString struct {
	s               []rune
	validConsecKeys []int
}

type row struct {
	stringGroups
	consecutiveKeys []int
}

// type stringAndKeys struct {
// 	s    string
// 	keys []int
// }

// type rangeAndLimit struct {
// 	low, high int
// 	limit     *int
// }

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
		stringGroups:    separateConsecutiveStrings(rowInfo[0]),
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

func separateConsecutiveStrings(s string) []stringGroup {
	var (
		fullRow   []stringGroup // eg. [["###", "??", "##"], ["?", "##", "?"]]
		sG        stringGroup   // eg. ["###", "??", "##"]
		curString strings.Builder
	)
	var i, j int
	for {
		if j == len(s) {
			break
		}
		switch {
		case string([]rune(s)[j]) == nonSpring:
			if curString.Len() > 0 {
				sG = append(
					sG,
					separatedString{
						s: []rune(curString.String()),
					},
				)
				fullRow = append(fullRow, sG)
				curString.Reset()
				sG = stringGroup{}
			}
			j++
			i = j
		case []rune(s)[i] == []rune(s)[j]:
			curString.WriteRune([]rune(s)[j])
			j++
		default:
			if curString.Len() > 0 {
				sG = append(
					sG,
					separatedString{
						s: []rune(curString.String()),
					},
				)
				curString.Reset()
			}
			i = j
		}
	}
	if curString.Len() > 0 {
		sG = append(
			sG,
			separatedString{
				s: []rune(curString.String()),
			},
		)
		fullRow = append(fullRow, sG)
	}
	return fullRow
}

func (sepString separatedString) knownBrokenSprings() []int {
	var knownBrokenSprings []int
	for i := 0; i < len(sepString.s); i++ {
		if string(sepString.s[i]) == brokenSpring {
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

func remainingConsecutiveSprings(sepStrings []separatedString) [][2]int {
	var (
		i, j, k              int
		lastSepString        = sepStrings[len(sepStrings)-1]
		remainingConsecutive [][2]int
	)

	for {
		if string(sepStrings[k].s[j]) == brokenSpring {
			j++
			continue
		}
		if i > j {
			remainingConsecutive = append(remainingConsecutive, [2]int{i, j - 1})
		}

		j++
		i = j
		if k == len(sepStrings)-1 && j == len(lastSepString.s) {
			break
		}
		if j == len(sepStrings[k].s) {
			j = 0
			i = j
			k++
		}

	}

	return remainingConsecutive
}

// func matchKeysToSepString(
// 	sepString separatedString,
// 	consecutiveKeys []int,
// 	curKeyIdx int,
// ) (string, []int) {
// 	var i, j int
// 	for {
// 		if j >=
// 		if j - i + 1 > consecutiveKeys[curKeyIdx]
// 	}
// }

// invariant: after constructing the stringGroups,
// if a sepString.s has a brokenSpring,
// the end of the sepString.s will always have the (consecutive) brokenSpring(s)
// EXCEPT for the totaling done in stringGroups.calcRecursiveTotal()
func (r row) calcSpringLocCombos() int {
	var (
		stringGroupIdx              int
		curKeyIdx, brokenSpringsIdx int
	)

	for {
		if stringGroupIdx >= len(r.stringGroups) {
			break
		}

		sG := r.stringGroups[stringGroupIdx]

		var sepStringIdx int
		for {
			if curKeyIdx == len(r.consecutiveKeys) || sepStringIdx == len(sG) {
				break
			}

			if len(sG[sepStringIdx].knownBrokenSprings()) > 0 {
				if len(sG[sepStringIdx].s) == r.consecutiveKeys[curKeyIdx] {
					// exact match of current keys, eg.
					// "??###" 5
					sG[sepStringIdx].validConsecKeys = append(
						sG[sepStringIdx].validConsecKeys,
						r.consecutiveKeys[curKeyIdx],
					)
					curKeyIdx++
					sepStringIdx++
					brokenSpringsIdx++
					continue
				}
				for rIdx := 0; rIdx < len(sG[sepStringIdx].s); rIdx++ {
					// add any extra ??'s at the beginning
					// of the current sepString to the previous sepString,
					// if the previous sepString has no brokenSprings
					// eg.
					// "???####" 5 means ? goes to the previous sepString
					numRunes := len(sG[sepStringIdx].s) - rIdx
					if numRunes == r.consecutiveKeys[curKeyIdx] {
						if sepStringIdx-1 > -1 {
							if len(sG[sepStringIdx-1].knownBrokenSprings()) == 0 {
								sG[sepStringIdx-1].s = append(
									sG[sepStringIdx-1].s,
									sG[sepStringIdx].s[:len(sG[sepStringIdx].s)-rIdx]...,
								)
							}
						}
						sG[sepStringIdx].validConsecKeys = append(
							sG[sepStringIdx].validConsecKeys,
							r.consecutiveKeys[curKeyIdx],
						)
						break
					}
				}
				if sepStringIdx+1 < len(sG) {
					// cut off the next ? for the barrier
					sG[sepStringIdx+1].s = sG[sepStringIdx+1].s[1:]
				}
				sepStringIdx++
				brokenSpringsIdx++
				curKeyIdx++
				continue
			}
			if len(sG[sepStringIdx].knownBrokenSprings()) == 0 {
				var runeIdx int
				keySum := r.consecutiveKeys[curKeyIdx]

				for {
					keysLeft := len(r.consecutiveKeys) - curKeyIdx - 1
					consecutiveBrokenStringsLeft := remainingConsecutiveSprings(sG[sepStringIdx:])
					if keysLeft <= len(consecutiveBrokenStringsLeft) {

					}

					if runeIdx >= len(sG[sepStringIdx].s) {
						if sepStringIdx+1 < len(sG) {
							// pass any extra ?? to the next one
							// eg. "?????" 3 > pass ? to the next
							// if the extra ?s are not used, they'll get passed back
							var passForwardIdx int
							if len(sG[sepStringIdx].validConsecKeys) > 0 {
								passForwardIdx = kit.Sum(sG[sepStringIdx].validConsecKeys) +
									len(sG[sepStringIdx].validConsecKeys) - 1
								if kBS := sG[sepStringIdx].knownBrokenSprings(); len(kBS) > 0 {
									passForwardIdx += kBS[0]
								}
							}
							sG[sepStringIdx+1].s = append(
								sG[sepStringIdx].s[passForwardIdx:],
								sG[sepStringIdx+1].s...,
							)
							sG[sepStringIdx].s = sG[sepStringIdx].s[:passForwardIdx]

							if len(sG[sepStringIdx].s) == 0 {
								sG = append(sG[:sepStringIdx], sG[sepStringIdx+1:]...)
								runeIdx = 0
								continue
							}
							if sepStringIdx > 0 {
								// cut off the extra ?
								sG[sepStringIdx+1].s = sG[sepStringIdx+1].s[1:]
							}
						}
						sepStringIdx++
						runeIdx = 0
						break
					}
					if runeIdx >= keySum {
						sG[sepStringIdx].validConsecKeys = append(
							sG[sepStringIdx].validConsecKeys,
							r.consecutiveKeys[curKeyIdx],
						)
						curKeyIdx++
						if curKeyIdx == len(r.consecutiveKeys) {
							break
						}
						keySum = r.consecutiveKeys[curKeyIdx] +
							len(r.consecutiveKeys[:curKeyIdx]) - 1
						runeIdx++
					}
					runeIdx++
				}
			}

		}
		r.stringGroups[stringGroupIdx] = sG
		stringGroupIdx++
	}

	total := r.stringGroups.calcRecursiveTotal()
	return total
}

// func upperLimitKnownBrokenSpring(
// 	rAL *rangeAndLimit,
// 	knownBrokenSpringsIdx int,
// 	knownBrokenSprings [][2]int,
// ) {
// 	if len(knownBrokenSprings) == 0 {
// 		return
// 	}
// 	if knownBrokenSprings[knownBrokenSpringsIdx][0]-2 < rAL.high+(rAL.high-rAL.low) {
// 		// min of next known broken spring minus 2,
// 		// or lowest known broken spring + # of consecutive springs

// 		limit := knownBrokenSprings[knownBrokenSpringsIdx][0] + rAL.high - rAL.low
// 		if knownBrokenSpringsIdx+1 <= len(knownBrokenSprings)-1 {
// 			limit = min(
// 				limit,
// 				knownBrokenSprings[knownBrokenSpringsIdx+1][0]-2,
// 			)
// 		}
// 		rAL.limit = &limit
// 	}

// }

// func splitPossibleAndBrokenSprings(s string) [][2]int {
// 	var consecutiveBroken [][2]int
// 	i := 0
// 	j := 0
// 	for {
// 		if j == len(s) {
// 			if string([]rune(s)[j-1]) == brokenSpring {
// 				consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
// 			}
// 			break
// 		}
// 		if string([]rune(s)[j]) == brokenSpring {
// 			j++
// 			continue
// 		}
// 		if j > i {
// 			consecutiveBroken = append(consecutiveBroken, [2]int{i, j - 1})
// 		}
// 		j++
// 		i = j
// 	}
// 	return consecutiveBroken
// }

// func calcRecursiveStringGroupTotal(sepStrings []separatedString, curIdx, prevOffset, maxNum int) int {

// 	var curOffset int

// 	for _, sepString := range sepStrings {
// 		var (
// 			keyIdx int
// 		)

// 	}

// 	if curIdx > 0 {
// 		prevHigh = sepStrings[curIdx-1].high
// 		for {
// 			lowerLimit := sepStrings[curIdx].low - 2 + curOffset
// 			if prevHigh+prevOffset > lowerLimit {
// 				curOffset++
// 				continue
// 			}
// 			break
// 		}
// 	}

// 	var total int
// 	if curIdx > len(sepStrings)-1 {
// 		return 0
// 	}
// 	if sepStrings[curIdx].limit != nil {
// 		curMax = *sepStrings[curIdx].limit + 1
// 	}

// 	if curIdx == len(sepStrings)-1 {
// 		for {
// 			lastStartOffset := sepStrings[curIdx].s[curOffset]
// 			lenOfRange := sepStrings[curIdx].high - sepStrings[curIdx].low + 1
// 			compareTotal := lastStartOffset + lenOfRange
// 			if compareTotal > curMax {
// 				break
// 			}
// 			total++
// 			curOffset++
// 		}
// 		return total
// 	}

// 	for {
// 		subTotal := calcRecursiveStringGroupTotal(sepStrings, curIdx+1, curOffset, maxNum)
// 		if subTotal == 0 {
// 			break
// 		}
// 		total += subTotal
// 		curOffset++
// 		if sepStrings[curIdx].limit != nil {
// 			if sepStrings[curIdx].high+curOffset > *sepStrings[curIdx].limit {
// 				break
// 			}
// 		}

// 	}
// 	return total
// }

func (sG *stringGroup) catchUp(stopAtIdx int, other stringGroup) {
	var tempSG stringGroup
	for i := 0; i < stopAtIdx; i++ {
		tempSG = append(tempSG, other[i])
	}
	*sG = tempSG
}

func (sGs stringGroups) calcRecursiveTotal() int {
	total := 1
	for _, sG := range sGs {
		combosToSum := []int{sG.calcRecursiveTotal()}
		combosToSum = append(
			combosToSum,
			sG.shiftAndReturnCombos()...,
		)
		total *= kit.Sum(combosToSum)
	}
	return total
}

func (sG stringGroup) calcRecursiveTotal() int {
	var totals []int
	for _, sepString := range sG {
		if len(sepString.knownBrokenSprings()) > 0 {
			totals = append(totals, 1)
			continue
		}
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
			// // catch the tempSG up with i - 1
			// if idxsBrokenSprings[i]-1 >= 0 {
			// 	tempSG.catchUp(i-1, sG)
			// }
			// eg. a string of 10 ?'s ("??????????")
			// and consecKeys of 2, 3, 1 would be 8
			// so there'd be two runes left over provide
			nextBrokenSprings := tempSG[curIdx+1].knownBrokenSprings()
			numFreeChars := len(tempSG[curIdx+1].s)
			if len(nextBrokenSprings) > 0 {
				numFreeChars = nextBrokenSprings[0]
			}
			for j := 0; j < numFreeChars; j++ {
				// as long as there's extra ?? in the next sepString,
				// add them to the end of the current sepString (containing brokenSprings)
				// and, if there's a previous sepString, add it to the prior sepString

				if string(tempSG[curIdx].s[0]) == brokenSpring ||
					string(tempSG[curIdx+1].s[0]) == brokenSpring {
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

// func calcRecursiveRangeTotal(ranges []rangeAndLimit, curIdx, prevOffset, maxNum int) int {

// 	var prevHigh, curOffset int

// 	curMax := maxNum

// 	if curIdx > 0 {
// 		prevHigh = ranges[curIdx-1].high
// 		for {
// 			lowerLimit := ranges[curIdx].low - 2 + curOffset
// 			if prevHigh+prevOffset > lowerLimit {
// 				curOffset++
// 				continue
// 			}
// 			break
// 		}
// 	}

// 	var total int
// 	if curIdx > len(ranges)-1 {
// 		return 0
// 	}
// 	if ranges[curIdx].limit != nil {
// 		curMax = *ranges[curIdx].limit + 1
// 	}

// 	if curIdx == len(ranges)-1 {
// 		for {
// 			lastStartOffset := ranges[curIdx].low + curOffset
// 			lenOfRange := ranges[curIdx].high - ranges[curIdx].low + 1
// 			compareTotal := lastStartOffset + lenOfRange
// 			if compareTotal > curMax {
// 				break
// 			}
// 			total++
// 			curOffset++
// 		}
// 		return total
// 	}

// 	for {
// 		subTotal := calcRecursiveRangeTotal(ranges, curIdx+1, curOffset, maxNum)
// 		if subTotal == 0 {
// 			break
// 		}
// 		total += subTotal
// 		curOffset++
// 		if ranges[curIdx].limit != nil {
// 			if ranges[curIdx].high+curOffset > *ranges[curIdx].limit {
// 				break
// 			}
// 		}

// 	}
// 	return total
// }
