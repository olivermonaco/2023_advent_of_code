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

		l.Info().Send()
		r := parseLine(ctx, line)

		rowCombos := r.sGs.calcRecursiveTotal(ctx)
		l.Info().Int("num_combos", rowCombos).Send()
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
	strGroups = strGroups.addKeysAndSeparate(ctx, consecNums, knownConsecBroken)

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

func (sGs stringGroups) addKeysAndSeparate(
	ctx context.Context,
	keys []int, knownConsecBroken [][2]int,
) stringGroups {
	// todo: do this with pointers
	var (
		keysIdx, knownConsecBrokenIdx int
	)
	newSGs := make(stringGroups, 0, len(sGs))
	for _, sG := range sGs {
		sG, keysIdx, knownConsecBrokenIdx = sG.addKeysAndSeparate(
			ctx,
			keys, knownConsecBroken,
			keysIdx, knownConsecBrokenIdx,
		)
		if len(sG) == 0 {
			continue
		}
		newSGs = append(newSGs, sG)
		if keysIdx == len(keys) {
			break
		}
	}
	return newSGs
}

func (sG stringGroup) addKeysAndSeparate(
	ctx context.Context,
	keys []int, knownConsecBroken [][2]int,
	curKeyIdx, knownConsecBrokenIdx int,
) (stringGroup, int, int) {
	var newSG stringGroup

	for _, sepString := range sG {
		// we're making a new string group from each separated string in the string group
		// by cutting and expanding the separated string into multiple
		var stringGroupFromSepString stringGroup
		stringGroupFromSepString, curKeyIdx, knownConsecBrokenIdx = sepString.addKeysAndSeparate(
			keys,
			knownConsecBroken,
			curKeyIdx,
			knownConsecBrokenIdx,
		)
		newSG = append(newSG, stringGroupFromSepString...)
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

func (sepString separatedString) addKeysAndSeparate(
	keys []int, knownConsecBroken [][2]int,
	curKeyIdx, knownConsecBrokenIdx int,
) (stringGroup, int, int) {
	var (
		i, curBrokenSpringsIdx int
		newStringGroup         stringGroup
	)

	// TODO: account for no matches for keys, but still want to appends

	curBrokenSprings := findSpansOfCharacters(sepString.s, 0, brokenSpring)

	j := len(sepString.s) // invalid for ranges, but logic catches

	for {

		if curKeyIdx < len(keys) {
			j = i + keys[curKeyIdx] - 1
		}

		newKeyIdx := curKeyIdx
		if curBrokenSpringsIdx > len(curBrokenSprings)-1 {
			// no broken springs left, get as many runes / keys as you can
			j = len(sepString.s) - 1
			newKeyIdx = catchKeysUp(j-i+1, curKeyIdx, keys)
		}

		if i > len(sepString.s) || j-i > len(sepString.s)-1 {
			// expected range or current idx is entirely greater than the remaining sepString.s len
			if i < len(sepString.s)-1 && len(newStringGroup) > 0 {
				// but there could be remaining ?s left to pass to the previous string group
				newSepString := separatedString{s: sepString.s[i : len(sepString.s)-1]}
				newStringGroup = newSepString.addPassBack(newStringGroup, len(newStringGroup)-1)
			}
			break
		}

		lowerLeft := i
		catchUpWithBrokenSprings := catchUpWithBrokenSprings(
			keys,
			knownConsecBroken,
			curKeyIdx,
			curBrokenSpringsIdx,
		)
		brokenSpringInRange := brokenSpringsExist(
			sepString.s,
			keys,
			i,
			curKeyIdx,
		)

		if catchUpWithBrokenSprings || brokenSpringInRange {
			lowerLeft, j = sepString.leftMostSpan(
				i,
				curKeyIdx,
				curBrokenSpringsIdx,
				keys,
				knownConsecBroken[knownConsecBrokenIdx:],
			)
		}

		newSepString := separatedString{
			s:               sepString.s[i : j+1],
			validConsecKeys: keys[curKeyIdx : newKeyIdx+1],
		}

		// lowerRight, _ := sepString.rightMostSpan(i, curKeyIdx, keys)
		// if j < len(sepString.s)-1 && lowerRight+keys[curKeyIdx] >= len(sepString.s) {
		// 	// TODO: this isn't right I think
		// 	// we haven't hit the right most limit,
		// 	// and there's room to add possible springs to the
		// 	newSepString.s = append(
		// 		newSepString.s,
		// 		sepString.s[j+1:len(sepString.s)]...,
		// 	)
		// }
		if len(newStringGroup) > 0 && i < lowerLeft {
			// add any extras before the lower limit to the previous
			newStringGroup[len(newStringGroup)-1].s = append(
				newStringGroup[len(newStringGroup)-1].s,
				newSepString.s[:lowerLeft-i]...,
			)
			// cut off those characters from the other
			newSepString.s = newSepString.s[lowerLeft-i:]
		}

		newStringGroup = append(newStringGroup, newSepString)

		if brokenSpringSpan := newSepString.knownBrokenSprings(); len(brokenSpringSpan) > 0 {
			curBrokenSpringsIdx += len(brokenSpringSpan)
		}
		i = j + 2
		curKeyIdx++

		if curKeyIdx == len(keys) {
			if i < len(sepString.s)-1 && len(newStringGroup) > 0 {
				newSepString := separatedString{s: sepString.s[i : len(sepString.s)-1]}
				newStringGroup = newSepString.addPassBack(newStringGroup, len(newStringGroup)-1)
			}
			break
		}
	}
	knownConsecBrokenIdx += len(curBrokenSprings)
	return newStringGroup, curKeyIdx, knownConsecBrokenIdx
}

func (sepString *separatedString) addPassBack(sG stringGroup, sGIdx int) stringGroup {
	if sGIdx < 0 {
		return sG
	}
	if sGIdx == 0 && sG[sGIdx].s[len(sG[sGIdx].s)-1] == brokenSpring {
		// if it's already left aligned, and it's the first one
		// there's no room to add to it
		return sG
	}

	brokenSpringsSpans := sepString.knownBrokenSprings()

	upperBound := len(sepString.s)
	if len(brokenSpringsSpans) > 0 {
		_, upperBound = sepString.rightMostSpan(0, 0, sepString.validConsecKeys)
		upperBound++ // upper bound of a slice, but rightMostSpan is inclusive of the upper bound
	}
	reqRunes := minReqRunes(sepString.validConsecKeys)
	numPassBack := upperBound - reqRunes
	addPossible := make([]rune, 0, numPassBack)
	for i := 0; i < numPassBack; i++ {
		addPossible = append(addPossible, possibleSpring)
		sepString.s = sepString.s[1:]
	}
	// add any extra possible springs
	sG[sGIdx].s = append(sG[sGIdx].s, addPossible...)

	if len(sepString.s) > 0 {
		// add to string group if there are remaining possible springs
		sG = append(sG, *sepString)
	}

	sG = sG[sGIdx].addPassBack(sG, sGIdx-1)
	return sG
}

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

func (sGs stringGroups) calcRecursiveTotal(ctx context.Context) int {
	var combosToMult []int
	for _, sG := range sGs {
		sGTotal := sG.calcRecursiveTotal(ctx)
		shiftedTotals := sG.shiftAndReturnCombos(ctx)
		combosToMult = append(combosToMult, kit.Sum(append(shiftedTotals, sGTotal)))
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
	if len(sepString.validConsecKeys) == 0 {
		return 0
	}
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

// minReqRunes calculates the minimum number of runes required to hold a slice of keys
func minReqRunes(keys []int) int {
	return max(
		kit.Sum(keys)+len(keys)-1,
		0,
	)
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
func (sepString separatedString) leftMostSpan(
	lowerLimit, curKeyIdx, brokenSpringsIdx int,
	keys []int,
	knownConsecBroken [][2]int,
) (int, int) {
	if len(knownConsecBroken)-(brokenSpringsIdx) == 0 {
		// order doesn't matter, send it
		return lowerLimit, len(sepString.s) - 1
	}
	i := knownConsecBroken[brokenSpringsIdx][0]
	j := knownConsecBroken[brokenSpringsIdx][1]
	if j-i+1 < keys[curKeyIdx] {
		for {
			// decrease i until it reaches the beginning of sepString.s,
			// or until it fulfills the curKey
			if i == lowerLimit || j-i+1 == keys[curKeyIdx] {
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
		// make sure the next rune is not a broken spring, otherwise shift entire thing
		// could be that we need to shift the entire span
		// eg. ?#??#? 4
		for {
			if j == len(sepString.s)-1 {
				break
			}
			if sepString.s[j+1] == brokenSpring {
				i++
				j++
				continue
			}
			break
		}
	}

	return i, j
}

// leftMostSpan gets the right most span possible for the current key idx considering the rune idx passed in
func (sepString separatedString) rightMostSpan(runeIdx, curKeyIdx int, keys []int) (int, int) {
	brokenSpringSpans := findSpansOfCharacters(sepString.s, runeIdx, brokenSpring)
	if len(brokenSpringSpans) == 0 || curKeyIdx == len(keys) {
		// order doesn't matter, send it
		return runeIdx, len(sepString.s) - 1
	}

	if len(brokenSpringSpans) == 0 {
		return 0, len(sepString.s) - 1
	}
	i := brokenSpringSpans[0][0]
	j := brokenSpringSpans[0][1]
	if j-i+1 < keys[curKeyIdx] {
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
		for {
			// decrease i until it reaches the beginning of sepString.s,
			// or until it fulfills the curKey
			if i == -1 || j-i+1 == keys[curKeyIdx] {
				break
			}
			i--
		}
	}

	return i, j
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
		curIdx := idxsBrokenSprings[i]
		brokenSpringSpan := tempSG[curIdx].knownBrokenSprings()
		if kit.Sum(tempSG[curIdx].validConsecKeys) == brokenSpringSpan[0][1]-brokenSpringSpan[0][0]+1 {
			// can't advance the broken springs, because all keys match the number of broken springs
			// at this point, everything should be right aligned
			// eg. ??### 3
			continue
		}
		// numFreeRunes := len(tempSG[curIdx].s) - minReqRunes(tempSG[curIdx].validConsecKeys)

		_, secondHighest := tempSG[curIdx].leftMostSpan(0, 0, 0, tempSG[curIdx].validConsecKeys, brokenSpringSpan)
		_, highest := tempSG[curIdx].rightMostSpan(0, 0, tempSG[curIdx].validConsecKeys)
		for j := 0; j < highest-secondHighest; j++ {
			if curIdx > 0 {
				tempSG[curIdx-1].s = append(tempSG[curIdx-1].s, tempSG[curIdx].s[0])
			}
			tempSG[curIdx].s = tempSG[curIdx].s[:len(tempSG[curIdx].s)-1]
			combos = append(combos, tempSG.calcRecursiveTotal(ctx))
		}

		var addtlFreeRunes int
		if curIdx < len(idxsBrokenSprings)-1 {
			addtlFreeRunes = findLeftMostRemaining(tempSG[curIdx+1])
		}

		for j := 0; j < addtlFreeRunes; j++ {

			if tempSG[curIdx].s[0] == brokenSpring ||
				tempSG[curIdx+1].s[0] == brokenSpring {
				// stop at the first brokenSpring in the current sepString or next sepString
				break
			}
			if idxsBrokenSprings[i] > 0 {
				// pass back any ??'s to the previous
				tempSG[i-1].s = append(tempSG[curIdx-1].s, tempSG[curIdx-1].s[:j+1]...)
			}
			// take one from the next sepString
			tempSG[curIdx].s = append(tempSG[curIdx].s, tempSG[curIdx+1].s[:1]...)
			tempSG[curIdx+1].s = tempSG[curIdx+1].s[1:]
			// remove one from current sepString
			tempSG[curIdx].s = tempSG[curIdx].s[1:]

			combos = append(combos, tempSG.calcRecursiveTotal(ctx))
		}
	}
	return combos
}

// catchUpWithBrokenSprings determines if there are more consecutive broken springs in a row
// than remaining keys
func catchUpWithBrokenSprings(keys []int, knownConsecBroken [][2]int, keysIdx, brokenSpringsIdx int) bool {
	consecBrokenLeft := (len(knownConsecBroken) - (brokenSpringsIdx + brokenSpringsIdx))

	return consecBrokenLeft >= remainingKeys(keys, keysIdx)
}

// brokenSpringsExist tells if there's a broken spring in the given range.
// The given range includes the current rune index to the current rune index plus the current key.
// This includes the buffer index (the index after the last rune in this range)
// eg. ???# 3 would return true (left to right of course)
func brokenSpringsExist(runes []rune, keys []int, runeIdx, keyIdx int) bool {
	return len(findSpansOfCharacters(runes[runeIdx:runeIdx+keys[keyIdx]], 0, brokenSpring)) > 0
}

// just to make it easier to read
func remainingKeys(keys []int, curIdx int) int {
	return len(keys) - curIdx
}
