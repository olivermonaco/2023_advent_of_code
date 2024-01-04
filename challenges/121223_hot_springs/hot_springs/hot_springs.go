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
	// rows := make([]row, 0, len(input))
	for i, line := range input {
		l := log.Ctx(ctx).With().Int("line num", i).Str("line_val", line).Logger()
		ctx = l.WithContext(ctx)
		// rows = append(rows, parseLine(ctx, line))
		// l := log.Ctx(ctx).With().Logger()
		l.Info().Msg("tests")

		r := parseLine(ctx, line)

		// TODO: ?###???????? 3,2,1 gives 1, should give more left off here
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

	strGroups, knownConsecBroken := createStringGroups(rowInfo[0])
	strGroups = strGroups.addKeys(consecNums, knownConsecBroken)

	return row{
		sGs:             strGroups,
		consecutiveKeys: consecNums,
	}
}

func createStringGroups(s string) (stringGroups, [][2]int) {
	separated := strings.Split(s, ".")

	sGs := make(stringGroups, 0, len(separated))
	knownConsecBrokenSprings := make([][2]int, 0, len(separated))

	for _, sG := range separated {
		if len(sG) > 0 {
			sepString := separatedString{s: []rune(sG)}
			knownConsecBrokenSprings = append(
				knownConsecBrokenSprings,
				sepString.knownBrokenSprings()...,
			)
			sGs = append(sGs,
				stringGroup{
					sepString,
				},
			)
		}
	}
	return sGs, knownConsecBrokenSprings
}

func (sGs stringGroups) addKeys(keys []int, knownConsecBroken [][2]int) stringGroups {
	// todo: do this with pointers
	var keysIdx int
	knownConsecBrokenIdx := max(len(knownConsecBroken)-1, 0)

	newSGs := make(stringGroups, 0, len(sGs))
	for _, sG := range sGs {
		sG, keysIdx, knownConsecBrokenIdx = sG.addKeys(
			keys, knownConsecBroken,
			keysIdx, knownConsecBrokenIdx,
		)
		if len(sG) == 0 {
			continue
		}
		newSGs = append(newSGs, sG)
	}
	return newSGs
}

func (sG stringGroup) addKeys(
	keys []int, knownConsecBroken [][2]int,
	curKeyIdx, knownConsecBrokenIdx int,
) (stringGroup, int, int) {
	var (
		newSG        stringGroup
		newCurKeyIdx int
	)
	for _, sepString := range sG {
		// we're making a new string group from each separated string in the string group
		// by cutting and expanding the separated string into multiple
		var newStringGroup stringGroup
		newStringGroup, newCurKeyIdx, knownConsecBrokenIdx = sepString.addKeys2(
			keys,
			knownConsecBroken,
			curKeyIdx,
			knownConsecBrokenIdx,
		)
		if len(newStringGroup) == 0 {
			continue
		}
		curKeyIdx = newCurKeyIdx
		newSG = append(newSG, newStringGroup...)
	}
	return newSG, curKeyIdx, knownConsecBrokenIdx
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

// addKeys2 is still iterating from right to left
func (sepString separatedString) addKeys2(
	keys []int, knownConsecBroken [][2]int,
	curKeyIdx, knownConsecBrokenIdx int,
) (stringGroup, int, int) {
	var (
		i, curBrokenSpringsIdx int
		newStringGroup         stringGroup
	)
	j := len(sepString.s) - 1

	curBrokenSprings := findSpansOfCharacters(sepString.s, 0, brokenSpring)

	for {
		if len(knownConsecBroken)-knownConsecBrokenIdx >= len(keys)-curKeyIdx {
			// get the right most span for greatest current broken idx
			_, j = sepString.leftMostSpan(
				curBrokenSprings[curBrokenSpringsIdx][0],
				curKeyIdx,
				keys,
			)
		}
		newKeyIdx := catchKeysUp(j-i+1, curKeyIdx, keys)
		newSepString := separatedString{
			s:               sepString.s[i : j+1],
			validConsecKeys: keys[curKeyIdx : newKeyIdx+1],
		}

		newStringGroup = append(newStringGroup, newSepString)

		curBrokenSpringsIdx++

		curKeyIdx = newKeyIdx + 1
		i = j + 2

		if i >= len(sepString.s) {
			break
		}
	}
	knownConsecBrokenIdx += len(curBrokenSprings)
	return newStringGroup, curKeyIdx, knownConsecBrokenIdx
}

// func (sepString separatedString) addKeys(keys []int, curKeyIdx int) (stringGroup, int) {
// 	var (
// 		i, j             int
// 		seenBrokenSpring bool
// 	)

// 	consecBrokenSprings := sepString.knownBrokenSprings()

// 	var newStringGroup []separatedString

// 	for {
// 		if j >= len(sepString.s) || curKeyIdx == len(keys) {
// 			break
// 		}
// 		if sepString.s[j] == brokenSpring {
// 			seenBrokenSpring = true
// 		}

// 		if j-i+1 >= keys[curKeyIdx] &&
// 			remainingKeys(keys, curKeyIdx) <= len(consecBrokenSprings) {
// 			// need to catch the idx up with the next consecutive number of broken springs
// 			brokenSpringSpans := findSpansOfCharacters(sepString.s, j, brokenSpring)
// 			j = brokenSpringSpans[0][1]

// 			if j-i >= keys[curKeyIdx] {
// 				newKeyIdx := catchKeysUp(j-i, curKeyIdx, keys)
// 				newSepString := separatedString{
// 					s:               sepString.s[i : j+1],
// 					validConsecKeys: keys[curKeyIdx : newKeyIdx+1],
// 				}
// 				newStringGroup = passBackAppend(newSepString, newStringGroup)
// 				curKeyIdx = newKeyIdx + 1
// 				j += 2
// 				i = j
// 				seenBrokenSpring = false
// 				continue
// 			}
// 			// if not, keep going
// 		}

// 		if j-i+1 >= keys[curKeyIdx] &&
// 			remainingKeys(keys, curKeyIdx) > len(consecBrokenSprings) &&
// 			!seenBrokenSpring {
// 			// eg. ????## 2, 3
// 			newSepString := separatedString{
// 				s:               sepString.s[i : j+1],
// 				validConsecKeys: keys[curKeyIdx : curKeyIdx+1],
// 			}
// 			newStringGroup = appendOrCombine(newSepString, newStringGroup)

// 			curKeyIdx++
// 			j += 2
// 			i = j
// 			continue
// 		}

// 		if j-i+1 == keys[curKeyIdx] || j == len(sepString.s) {
// 			newSepString := separatedString{
// 				s:               sepString.s[i : j+1],
// 				validConsecKeys: []int{keys[curKeyIdx]},
// 			}
// 			newStringGroup = passBackAppend(newSepString, newStringGroup)

// 			curKeyIdx++
// 			j += 2
// 			i = j
// 			seenBrokenSpring = false
// 			continue
// 		}
// 		j++
// 	}

// 	return newStringGroup, curKeyIdx
// }

func passBackAppend(newSepString separatedString, newStringGroup stringGroup) stringGroup {
	possiblePassBack := findLeftMostRemaining(newSepString)
	if len(newStringGroup) > 0 {
		for i := 0; i < possiblePassBack; i++ {
			newStringGroup[len(newStringGroup)-1].s = append(
				newStringGroup[len(newStringGroup)-1].s,
				possibleSpring,
			)
		}
	}
	return append(newStringGroup, newSepString)
}

// appendOrCombine appends the new separated string to the slice of separated strings,
// or combines it with the previous sepString (adding an extra possible string for a buffer)
func appendOrCombine(newSepString separatedString, newStringGroup stringGroup) []separatedString {
	if len(newStringGroup) == 0 {
		return []separatedString{newSepString}
	}
	if len(newStringGroup[len(newStringGroup)-1].knownBrokenSprings()) > 0 {
		return append(newStringGroup, newSepString)
	}
	// reconstruct the sepStrings considering the last in the known broken
	// and add back in the buffer that was cut off by the idx jump
	newSepString.s = append(
		newStringGroup[len(newStringGroup)-1].s,
		newSepString.s...,
	)
	newSepString.validConsecKeys = append(
		newStringGroup[len(newStringGroup)-1].validConsecKeys,
		newSepString.validConsecKeys...,
	)
	// cut off the last one, because we're about to replace it
	newStringGroup = newStringGroup[:len(newStringGroup)-1]

	return append(newStringGroup, newSepString)
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

func (sepString separatedString) knownBrokenSprings() [][2]int {
	var (
		i, j         int
		consecBroken [][2]int
	)

	for {
		if j == len(sepString.s) {
			if i < j {
				consecBroken = append(consecBroken, [2]int{i, j - 1})
			}
			break
		}
		if sepString.s[j] == brokenSpring {
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

func (sG stringGroup) knownBrokenSepStrings() []int {
	var brokenSpringsIdxs []int
	for i, sepString := range sG {
		if len(sepString.knownBrokenSprings()) > 0 {
			brokenSpringsIdxs = append(brokenSpringsIdxs, i)
		}
	}
	return brokenSpringsIdxs
}

// iterate in reverse
func catchKeysUp(span, curKeyIdx int, consecutiveKeys []int) int {
	newKeyIdx := curKeyIdx
	for {
		if span < minReqRunes(consecutiveKeys[curKeyIdx:newKeyIdx+1]) {
			break
		}
		newKeyIdx++

		if newKeyIdx == len(consecutiveKeys) {
			break
		}
	}
	return newKeyIdx - 1
}

func (r row) calcSpringLocCombos(ctx context.Context) int {
	total := r.sGs.calcRecursiveTotal(ctx)
	return total
}

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

// func findLeftMostRemaining(sepString separatedString) int {
// 	brokenSpringSpans := findSpansOfCharacters(sepString.s, 0, brokenSpring)
// 	if len(brokenSpringSpans) == 0 {
// 		// order doesn't matter, math it
// 		sumNums := kit.Sum(sepString.validConsecKeys)
// 		lenNums := len(sepString.validConsecKeys) - 1
// 		return len(sepString.s) - sumNums - lenNums
// 	}
// 	var i, j, remainingLeftMost int
// 	curKeyIdx := len(sepString.validConsecKeys) - 1
// 	brokenSpringSpansIdx := len(brokenSpringSpans) - 1
// 	curRunes := sepString.s
// 	for {
// 		if brokenSpringSpansIdx < 0 || curKeyIdx < 0 {
// 			sumNums := kit.Sum(sepString.validConsecKeys[:curKeyIdx+1])
// 			lenNums := len(sepString.validConsecKeys[:curKeyIdx+1])
// 			remainingLeftMost = (j - i + 1) - sumNums - lenNums
// 			break
// 		}
// 		i = brokenSpringSpans[brokenSpringSpansIdx][0]
// 		j = brokenSpringSpans[brokenSpringSpansIdx][1]
// 		curSpan := j - i + 1
// 		curKey := sepString.validConsecKeys[curKeyIdx]
// 		if curSpan < curKey {
// 			// increase j until the end of the slice, or it fulfills the curKey
// 			for {
// 				if j == len(curRunes) || j-i == sepString.validConsecKeys[curKeyIdx] {
// 					break
// 				}
// 				j++
// 			}
// 			for {
// 				curSpan = j - i
// 				if i == 0 || curSpan == sepString.validConsecKeys[curKeyIdx] {
// 					break
// 				}
// 				i--
// 			}
// 		}
// 		j = max(i-2, -1) // 2 for buffer, and min -1 to account for no remainingLeftMost
// 		i = 0            // will get increased if there's any remaining brokenSprings
// 		curRunes = curRunes[:j+1]
// 		curKeyIdx--
// 		brokenSpringSpansIdx--
// 	}
// 	return remainingLeftMost
// }

// minReqRunes calculates the minimum number of runes required to hold a slice of keys
func minReqRunes(keys []int) int {
	return kit.Sum(keys) + len(keys) - 1
}

func findLeftMostRemaining(sepString separatedString) int {
	brokenSpringSpans := findSpansOfCharacters(sepString.s, 0, brokenSpring)
	if len(brokenSpringSpans) == 0 {
		// order doesn't matter, math it
		return len(sepString.s) - minReqRunes(sepString.validConsecKeys)
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
		j = max(i-2, -1) // 2 for buffer, and min -1 to account for no remainingLeftMost
		i = 0            // will get increased if there's any remaining brokenSprings
		curRunes = curRunes[:j+1]
		curKeyIdx--
		brokenSpringSpansIdx--
	}
	return remainingLeftMost
}

// leftMostSpan gets the right most span possible for the current key idx considering the rune idx passed in
func (sepString separatedString) leftMostSpan(runeIdx, curKeyIdx int, keys []int) (int, int) {
	brokenSpringSpans := findSpansOfCharacters(sepString.s, runeIdx, brokenSpring)
	if len(brokenSpringSpans) == 0 {
		// order doesn't matter, send it
		return runeIdx, len(sepString.s) - 1
	}
	var (
		i, j int
	)
	if len(brokenSpringSpans) == 0 {
		return 0, len(sepString.s) - 1
	}
	j = brokenSpringSpans[0][1]
	if j-i+1 < keys[curKeyIdx] {
		for {
			// decrease i until it reaches the beginning of sepString.s,
			// or until it fulfills the curKey
			if i == -1 || j-i+1 == keys[curKeyIdx] {
				break
			}
			i--
		}

		// increase j until the end of the slice,
		// or it fulfills the curKey
		for {
			if j == len(sepString.s) || j-i+1 == keys[curKeyIdx] {
				// off by 1 if j was len(sepString.s), so decrement by 1
				// there are two possible cases here, so avoid an extra if by using min
				j = min(j, len(sepString.s)-1)
				break
			}
			j++
		}
	}

	return i, j
}

// // findRightMostSpan gets the right most span possible for the current key idx considering the rune idx passed in
// func (sepString separatedString) findRightMostSpan(runeIdx, curKeyIdx int, keys []int) (int, int) {
// 	brokenSpringSpans := findSpansOfCharacters(sepString.s, runeIdx, brokenSpring)
// 	if len(brokenSpringSpans) == 0 {
// 		// order doesn't matter, send it
// 		return runeIdx, len(sepString.s) - 1
// 	}
// 	var (
// 		i, j int
// 	)
// 	brokenSpringSpansIdx := len(brokenSpringSpans) - 1
// 	j = brokenSpringSpans[brokenSpringSpansIdx][1]
// 	for {
// 		if brokenSpringSpansIdx < 0 {
// 			break
// 		}
// 		// at each iteration start i over at the beginning of the span
// 		i = brokenSpringSpans[brokenSpringSpansIdx][0]
// 		if j-i+1 < keys[curKeyIdx] {
// 			// increase j until the end of the slice,
// 			// or it fulfills the curKey
// 			for {
// 				if j == len(sepString.s) || j-i+1 == keys[curKeyIdx] {
// 					// off by 1 if j was len(sepString.s), so decrement by 1
// 					// there are two possible cases here, so avoid an extra if by using min
// 					j = min(j, len(sepString.s)-1)
// 					break
// 				}
// 				j++
// 			}

// 			for {
// 				// decrease i until it reaches the beginning of sepString.s,
// 				// or until it fulfills the curKey
// 				if i == -1 || j-i+1 == keys[curKeyIdx] {
// 					break
// 				}
// 				i--
// 			}
// 		}
// 		newKeyIdx := curKeyIdx
// 		// catch the curKeyIdx up with the span
// 		for {
// 			if minReqRunes(keys[newKeyIdx:curKeyIdx+1]) > j-i+1 {
// 				break
// 			}
// 			newKeyIdx--
// 		}
// 		curKeyIdx = newKeyIdx + 1
// 		// account for the buffer
// 		j = i - 1
// 		brokenSpringSpansIdx--
// 	}

// 	return i, i + keys[curKeyIdx]
// }

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
			numFreeRunes := findLeftMostRemaining(tempSG[curIdx+1])

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
