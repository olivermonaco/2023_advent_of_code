package hot_springs

import (
	"cmp"
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	brokenSpring   = '#'
	possibleSpring = '?'
	nonSpring      = '.'
)

var resultsCache [][]int

type charsBrokenSpans struct {
	chars       []rune
	brokenSpans [][2]int
}

func (cBs charsBrokenSpans) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("chars", string(cBs.chars)).
		Interface("broken_spans", cBs.brokenSpans)
}

// -------- //

type separatedString struct {
	chars []rune
	keys  []int
}

func (sepStr separatedString) strings() []string {
	return append(
		[]string{
			(string(sepStr.chars))},
		kit.Map(sepStr.keys, func(k int) string { return strconv.Itoa(k) })...,
	)
}

func (sepString separatedString) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("runes", string(sepString.chars)).
		Ints("keys", sepString.keys)
}

type separatedStrings []separatedString

func (sepStrings separatedStrings) MarshalZerologArray(a *zerolog.Array) {
	for _, sepString := range sepStrings {
		a.Object(sepString)
	}
}

// -------- //

type stringGroup struct {
	keys []int
	charsBrokenSpans
}

func (sG stringGroup) MarshalZerologObject(e *zerolog.Event) {
	e.
		Ints("keys", sG.keys).
		Object("chars_broken_spans", sG.charsBrokenSpans)
}

func (sG stringGroup) copy() stringGroup {
	cB := charsBrokenSpans{
		chars:       kit.Map(sG.chars, func(c rune) rune { return c }),
		brokenSpans: kit.Map(sG.brokenSpans, func(bkSp [2]int) [2]int { return [2]int{bkSp[0], bkSp[1]} }),
	}
	new := stringGroup{
		keys:             kit.Map(sG.keys, func(k int) int { return k }),
		charsBrokenSpans: cB,
	}
	return new
}

func (sG stringGroup) initRange() separatedStringRefs {
	initialRange := make(separatedStringRefs, 0, len(sG.keys))

	var i, j int
	for _, key := range sG.keys {
		j = i + key - 1
		if i+key > len(sG.chars) {
			// i + key range is entirely greater than the sepString.s len, even at the outset
			return separatedStringRefs{}
		}
		initialRange = append(
			initialRange,
			separatedStringRef{
				start: i, end: j,
			},
		)
		i = j + 2
	}
	return initialRange
}

func (sG stringGroup) possibleCombos2(ctx context.Context) []refBuffGroups {
	// l := log.Ctx(ctx).With().Logger()
	if len(sG.chars) == 0 || len(sG.keys) == 0 {
		return []refBuffGroups{}
	}

	initialRange := sG.initRange()
	sortedBroken := initialRange.catchupBroken(sG.charsBrokenSpans)
	if sortedBroken == nil {
		return []refBuffGroups{}
	}
	initialRefBuffs := sortedBroken.toRefBuffs(len(sG.chars) - 1)
	initRefBuffGroups := initialRefBuffs.toGroups()

	combos := initRefBuffGroups.shiftGroups(sG.charsBrokenSpans, len(initRefBuffGroups)-1, len(sG.charsBrokenSpans.chars))
	minLen := kit.Reduce(combos, len(combos[0]), func(i int, rBGs refBuffGroups) int {
		if len(rBGs) < i {
			return len(rBGs)
		}
		return i
	})
	combos = kit.Filter(combos, func(a refBuffGroups) bool {
		return len(a) == minLen
	})
	slices.SortFunc(combos, func(a, b refBuffGroups) int {
		return a.compare(b)
	})
	combos = slices.CompactFunc(combos, func(a, b refBuffGroups) bool {
		return a.equivalent(b)
	})
	// l.Info().Interface("combos", combos).Send()
	return combos
}

// -------- //

type stringGroups []stringGroup

func (sGs stringGroups) MarshalZerologArray(a *zerolog.Array) {
	for _, sG := range sGs {
		a.Object(sG)
	}
}

func (sGs stringGroups) copy() stringGroups {
	return kit.Map(sGs, func(sG stringGroup) stringGroup { return sG.copy() })
}

func (sGs stringGroups) possibleCombos2(ctx context.Context) int {
	// l := log.Ctx(ctx).With().Logger()
	allcombos := make([]int, 0, len(sGs))
	// for i, sG := range sGs {
	for _, sG := range sGs {
		if brokenLen(sG.brokenSpans) > kit.Sum(sG.keys) {
			return 0
		}
		// l.Info().Int("idx", i).Msg("starting idx")
		mKL := minKeysLen(sG.keys)
		if mKL > len(sG.chars) {
			return 0
		}
		combos := sG.possibleCombos2(ctx)
		if len(combos) == 0 && len(sG.keys) > 0 {
			return 0
		}
		var groupTotals []int
		for _, rBGs := range combos {
			rGBsTotals := rBGs.calcTotals()
			if rGBsTotal := kit.Sum(rGBsTotals); rGBsTotal > 0 {
				groupTotals = append(groupTotals, kit.Sum(rGBsTotals))
			}
		}
		if groupTotal := kit.Sum(groupTotals); groupTotal > 0 {
			allcombos = append(allcombos, kit.Sum(groupTotals))
		}
	}

	return kit.Mult(allcombos)
}

func (sG stringGroup) possibleCombos(ctx context.Context) separatedStrings {
	if len(sG.chars) == 0 || len(sG.keys) == 0 {
		return separatedStrings{}
	}

	initialRange := sG.initRange()

	combos := initialRange.allCombos(ctx, sG.charsBrokenSpans, len(initialRange)-1, len(sG.chars))
	l := log.Ctx(ctx).With().
		Ints("found_keys", sG.keys).
		Str("string_group", string(sG.chars)).
		Logger()
	if len(combos) == 0 {
		l.Info().Msg("no combos found for string group")
	} else {
		logCombos := make([][]string, 0, (len(combos)*2)+2)
		logCombos = append(logCombos, []string{"string", "keys"})
		logCombos = append(logCombos, []string{
			string(sG.chars), strings.Join(
				kit.Map(
					sG.keys,
					func(k int) string { return strconv.Itoa(k) },
				),
				"_"),
		},
		)
		for i, combo := range combos {
			logCombos = append(logCombos, combo.strings())
			if i%len(sG.keys) == 0 {
				logCombos = append(logCombos, []string{"", ""})
			}
		}
		writeToCSV(logCombos, openCSVFile("combos"))
		l.Info().Array("combos", combos).Msg("found combo for string group")
	}

	return combos
}

func (sGs stringGroups) possibleCombos(ctx context.Context) int {
	l := log.Ctx(ctx).With().Logger()
	allcombos := make([]int, 0, len(sGs))
	for _, sG := range sGs {
		if brokenLen(sG.brokenSpans) > kit.Sum(sG.keys) {
			return 0
		}
		mKL := minKeysLen(sG.keys)
		if mKL > len(sG.chars) {
			return 0
		}

		combos := len(sG.possibleCombos(ctx))
		if combos == 0 && len(sG.keys) > 0 {
			return 0
		}
		if combos > 0 {
			mod := combos % len(sG.keys)
			if mod != 0 {
				l.Info().Array("string_groups", sGs).
					Msgf("found %d combos, but invalid",
						kit.Sum(allcombos),
					)
				panic(mod)
			}
			allcombos = append(allcombos, combos/len(sG.keys))
		}
	}

	return kit.Mult(allcombos)
}

// -------- //

type stringGroupsKeys struct {
	keys []int
	sGs  stringGroups
	// would have preferred to create a type range or span struct and have this be a slice of that in the end
	brokenSprings [][2]int
}

func (sGsKs stringGroupsKeys) calcTotal1(ctx context.Context) int {
	var sumCombos []int
	sGsKs.sGs[0].keys = sGsKs.keys
	shiftKeysTotal := sGsKs.shiftKeys(ctx, sGsKs.sGs, 0)

	sumCombos = append(sumCombos, shiftKeysTotal)
	return kit.Sum(sumCombos)
}

func (sGsKs stringGroupsKeys) calcTotal2(ctx context.Context) int {
	var sumCombos []int
	sGsKs.sGs[0].keys = sGsKs.keys
	shiftKeysTotal := sGsKs.shiftKeys2(ctx, sGsKs.sGs, 0)

	sumCombos = append(sumCombos, shiftKeysTotal)

	return kit.Sum(sumCombos)
}

func (sGsKs stringGroupsKeys) shiftKeys(
	ctx context.Context,
	sGs stringGroups,
	curSGIdx int,
) int {
	var total int
	if curSGIdx == len(sGs) {
		return total
	}

	curKeys := sGs[curSGIdx].keys
	// iterate to try and pass all keys to the next sG
	// if it's invalid there will be left over at the end of the slice
	for i := 0; i <= len(curKeys); i++ {
		tempSGs := sGs.copy()
		if curSGIdx == len(tempSGs)-1 {
			total += tempSGs.possibleCombos(ctx)
			break
		}

		tempSGs[curSGIdx].keys = curKeys[:i]
		// pass the remaining keys to the next
		tempSGs[curSGIdx+1].keys = curKeys[i:]
		sumKeys := kit.Sum(tempSGs[curSGIdx].keys)
		brokenSpansLen := brokenLen(tempSGs[curSGIdx].brokenSpans)
		if brokenSpansLen > sumKeys || sumKeys > len(tempSGs[curSGIdx].chars) {
			continue
		}

		nextTotal := sGsKs.shiftKeys(
			ctx,
			tempSGs.copy(),
			curSGIdx+1,
		)
		if nextTotal == 0 {
			// all totals so far are invalid if the next one is invalid
			continue
		}
		total += nextTotal

	}
	return total
}

func (sGsKs stringGroupsKeys) shiftKeys2(
	ctx context.Context,
	sGs stringGroups,
	curSGIdx int,
) int {
	var total int
	if curSGIdx == len(sGs) {
		return total
	}

	curKeys := sGs[curSGIdx].keys
	// iterate to try and pass all keys to the next sG
	// if it's invalid there will be left over at the end of the slice
	for i := 0; i <= len(curKeys); i++ {
		tempSGs := sGs.copy()
		if curSGIdx == len(tempSGs)-1 {
			mKL := minKeysLen(tempSGs[curSGIdx].keys)
			brokenSpansLen := brokenLen(tempSGs[curSGIdx].brokenSpans)
			if brokenSpansLen > mKL || mKL > len(tempSGs[curSGIdx].chars) {
				continue
			}
			total += tempSGs.possibleCombos2(ctx)
			break
		}

		tempSGs[curSGIdx].keys = curKeys[:i]
		// pass the remaining keys to the next
		tempSGs[curSGIdx+1].keys = curKeys[i:]

		mKL := minKeysLen(tempSGs[curSGIdx].keys)
		brokenSpansLen := brokenLen(tempSGs[curSGIdx].brokenSpans)
		if brokenSpansLen > mKL || mKL > len(tempSGs[curSGIdx].chars) {
			continue
		}

		nextTotal := sGsKs.shiftKeys2(
			ctx,
			tempSGs.copy(),
			curSGIdx+1,
		)
		if nextTotal == 0 {
			// all totals so far are invalid if the next one is invalid
			continue
		}
		total += nextTotal

	}
	return total
}

// -------- //

type separatedStringRef struct {
	start, end  int
	brokenSpans [][2]int
}

func (sepStrRef separatedStringRef) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int("end", sepStrRef.end).
		Int("start", sepStrRef.start).
		Interface("broken_spans", sepStrRef.brokenSpans)
}

func (sepStrRef separatedStringRef) copy() separatedStringRef {
	return separatedStringRef{
		start: sepStrRef.start,
		end:   sepStrRef.end,
		brokenSpans: kit.Map(
			sepStrRef.brokenSpans,
			func(broken [2]int) [2]int { return [2]int{broken[0], broken[1]} },
		),
	}
}

func (sepStrRef *separatedStringRef) addBrokenSpans(allBrokenSpans [][2]int) {
	sepStrRef.brokenSpans = brokenSpansInRange(sepStrRef.start, sepStrRef.end, allBrokenSpans)
}

// ------ //

type separatedStringRefs []separatedStringRef

func (sepStrRefs separatedStringRefs) copy() separatedStringRefs {
	return kit.Map(sepStrRefs, func(sepStrRef separatedStringRef) separatedStringRef { return sepStrRef.copy() })
}

func (sepStrRefs separatedStringRefs) MarshalZerologArray(a *zerolog.Array) {
	for _, sepStrRef := range sepStrRefs {
		a.Object(sepStrRef)
	}
}

func (sepStrRefs separatedStringRefs) validCombo(ctx context.Context, cBs charsBrokenSpans) separatedStrings {
	l := log.Ctx(ctx).With().Logger()
	if len(sepStrRefs) == 0 || len(cBs.chars) == 0 {
		return separatedStrings{}
	}

	if len(sepStrRefs.remainingBrokenSpans(cBs.brokenSpans)) > 0 {
		// unencompassed broken springs
		return separatedStrings{}
	}

	validCombo := make(separatedStrings, 0, len(sepStrRefs))

	for i, sepStrRef := range sepStrRefs {
		cutoff := len(cBs.chars)
		if i < len(sepStrRefs)-1 {
			cutoff = sepStrRefs[i+1].start - 1
		}

		if sepStrRef.end >= cutoff {
			// check end is too close to the nextMax
			// calling funcs should account for this already but just in case
			return separatedStrings{}
		}
		validSepStr := sepStrRef.toSeparatedString(cBs.chars)

		validCombo = append(validCombo, validSepStr)
	}
	l.Info().Array("sep_str_refs", sepStrRefs).Msg("valid_combo")
	return validCombo
}

func (sepStrRefs separatedStringRefs) bumpForward(refIdx, numSpaces int) separatedStringRefs {
	new := make(separatedStringRefs, 0, len(sepStrRefs))
	for _, sepStrRef := range sepStrRefs[refIdx:] {
		new = append(
			new,
			separatedStringRef{
				start: sepStrRef.start + numSpaces,
				end:   sepStrRef.end + numSpaces,
			},
		)
	}
	return new
}

func (sepStrRefs separatedStringRefs) catchupBroken(cBs charsBrokenSpans) *separatedStringRefs {
	if len(cBs.brokenSpans) == 0 {
		return &sepStrRefs
	}
	remaining := sepStrRefs.remainingBrokenSpans(cBs.brokenSpans)
	if len(remaining) == 0 {
		if sepStrRefs[len(sepStrRefs)-1].end > len(cBs.chars)-1 {
			return nil
		}
		for i := 0; i < len(sepStrRefs); i++ {
			sepStrRefs[i].addBrokenSpans(cBs.brokenSpans)
		}
		return &sepStrRefs
	}
	least := remaining[0]
	i := nextLowest(sepStrRefs, remaining[0])
	if i == -1 {
		// invalid case, lowest sepStrRef may not be big enough for lowest broken span
		return nil
	}
	c := sepStrRefs.copy()
	if i < len(c)-1 && c[i+1].start-1 <= least[1] {
		numSpaces := least[1] + 2 - c[i+1].start
		c = append(c[:i+1], c.bumpForward(i+1, numSpaces)...)
	}
	diff := c[i].end - c[i].start + 1
	c[i].end = least[1]
	c[i].start = c[i].end - diff + 1
	return c.catchupBroken(cBs)
}

func (sepStrRefs separatedStringRefs) allCombos(ctx context.Context, cBSs charsBrokenSpans, refIdx, cutoff int) separatedStrings {
	validCombos := make(separatedStrings, 0, len(sepStrRefs))
	if refIdx == -1 {
		return sepStrRefs.validCombo(ctx, cBSs)
	}
	c := sepStrRefs.copy()
	validCombos = append(validCombos, c.allCombos(ctx, cBSs, refIdx-1, c[refIdx].start-1)...)

	if c[refIdx].end+1 >= cutoff {
		return validCombos
	}
	c[refIdx].start++
	c[refIdx].end++
	validCombos = append(validCombos, c.allCombos(ctx, cBSs, refIdx, cutoff)...)
	return validCombos
}

func (sepStrRef separatedStringRef) toSeparatedString(runes []rune) separatedString {
	return separatedString{
		chars: runes[sepStrRef.start : sepStrRef.end+1],
		keys:  []int{sepStrRef.end + 1 - sepStrRef.start},
	}
}

func (rBG *refBuffGroup) removeLastRefBuff() {
	// remove the refBuff and brokenSpan we're moving to its own group from the current
	if len(rBG.refBuffs[len(rBG.refBuffs)-1].brokenSpans) > 0 {
		numToRemove := len(rBG.refBuffs[len(rBG.refBuffs)-1].brokenSpans)
		numSpans := len(rBG.brokenSpans)
		rBG.brokenSpans = rBG.brokenSpans[:numSpans-numToRemove]
	}
	rBG.refBuffs = rBG.refBuffs[:len(rBG.refBuffs)-1]
}

func (sepStrRefs separatedStringRefs) remainingBrokenSpans(brokenSpringSpans [][2]int) [][2]int {
	var sepStrRefIdx, brokenIdx int

	remBroken := make([][2]int, 0, len(brokenSpringSpans))

	for {
		if brokenIdx == len(brokenSpringSpans) {
			break
		}
		if sepStrRefIdx == len(sepStrRefs) {
			remBroken = append(remBroken, brokenSpringSpans[brokenIdx:]...)
			break
		}
		if sepStrRefs[sepStrRefIdx].start > brokenSpringSpans[brokenIdx][0] {
			remBroken = append(remBroken, brokenSpringSpans[brokenIdx])
			brokenIdx++
			continue
		}
		if sepStrRefs[sepStrRefIdx].end < brokenSpringSpans[brokenIdx][1] {
			sepStrRefIdx++
			continue
		}
		brokenIdx++
	}
	return remBroken
}

func (sepStrRefs separatedStringRefs) toRefBuffs(maxIdx int) refBuffs {
	refBuffs := make(refBuffs, len(sepStrRefs))
	for i := len(sepStrRefs) - 1; i > -1; i-- {
		refBuff := refBuff{
			separatedStringRef: sepStrRefs[i],
		}
		if sepStrRefs[i].end < maxIdx {
			refBuff.rBuff = maxIdx - sepStrRefs[i].end
		}
		maxIdx = refBuff.start - 2
		refBuffs[i] = refBuff
	}
	return refBuffs
}

// -------- //

type refBuff struct {
	separatedStringRef
	lBuff, rBuff int
}

func (rB refBuff) compare(other refBuff) int {
	return cmp.Compare(rB.start, other.start)
}

func (rB refBuff) eq(other refBuff) bool {
	if rB.start != other.start {
		return false
	}
	if rB.end != other.end {
		return false
	}
	return slices.EqualFunc(rB.brokenSpans, other.brokenSpans, bkSpansEqFunc)
}

func (rB *refBuff) leftAlign(lCutoff int) {
	diff := rB.lBuff
	if len(rB.brokenSpans) > 0 {
		diff = min(rB.start-lCutoff-1, rB.end-rB.brokenSpans[len(rB.brokenSpans)-1][1])
		rB.start -= diff
		rB.end -= diff
		rB.rBuff += diff
		rB.lBuff = 0
		return
	}
	rB.start -= diff
	rB.end -= diff
	rB.rBuff += diff
	rB.lBuff -= diff

	above := (rB.start - rB.lBuff)
	above = above - lCutoff - 1
	if above > 0 {
		rB.start -= above
		rB.end -= above
		rB.rBuff += above
	}
}

func (rB refBuff) copy() refBuff {
	return refBuff{
		separatedStringRef: rB.separatedStringRef.copy(),
		rBuff:              rB.rBuff,
		lBuff:              rB.lBuff,
	}
}

type createRefBuffInput struct {
	start, end, lBuff, rBuff int
	allBrokenSpans           [][2]int
}

func createRefBuff(input createRefBuffInput) refBuff {
	return refBuff{
		separatedStringRef: separatedStringRef{start: input.start, end: input.end, brokenSpans: input.allBrokenSpans},
		lBuff:              input.lBuff,
		rBuff:              input.rBuff,
	}
}

func (rB refBuff) shift(chars []rune, lowerCutoff, upperCutoff int) *refBuff {
	if rB.end >= upperCutoff {
		shiftBack := rB.end - upperCutoff + 1
		if rB.start-shiftBack < lowerCutoff {
			return nil
		}
		rB.start = rB.start - shiftBack
		rB.end = rB.end - shiftBack
	}
	if rB.start-1 > 0 && chars[rB.start-1] == brokenSpring {
		return nil
	}
	return &rB
}

type refBuffs []refBuff

func (rBs refBuffs) sepStrRefs() separatedStringRefs {
	sepStrRefs := make(separatedStringRefs, 0, len(rBs))
	for _, rB := range rBs {
		sepStrRefs = append(sepStrRefs, rB.separatedStringRef)
	}
	return sepStrRefs
}

func (rBs refBuffs) toGroups() refBuffGroups {
	var curRBG refBuffGroup
	rBGs := make(refBuffGroups, 0, len(rBs))
	for i, rB := range rBs {
		if len(rB.brokenSpans) > 0 {
			if i > 0 && len(curRBG.refBuffs) > 0 {
				rBGs = append(rBGs, curRBG)
			}
			rBGs = append(
				rBGs,
				refBuffGroup{
					refBuffs:    refBuffs{rB},
					brokenSpans: rB.brokenSpans,
				},
			)
			curRBG = refBuffGroup{}
			continue
		}
		curRBG.refBuffs = append(curRBG.refBuffs, rB)
		if i == len(rBs)-1 {
			rBGs = append(rBGs, curRBG)
		}
	}
	return rBGs
}

// assumes that the rBuff is always the last indexed refBuff
func (rBs refBuffs) calcNonBrokenTotals() int {
	return calcNums(len(rBs)-1, rBs[len(rBs)-1].rBuff)
}

type refBuffGroup struct {
	refBuffs
	brokenSpans [][2]int
}

func (rBG refBuffGroup) add(n int) *refBuffGroup {
	if len(rBG.brokenSpans) > 0 {
		maxN := fitAvailable(rBG.refBuffs[0], rBG.brokenSpans, n)
		if maxN < n {
			return nil
		}
		cRBG := rBG.copy()

		// in this case we need to add to start / end,
		// because the n is coming from a non broken group
		cRBG.refBuffs[0].start += maxN
		cRBG.refBuffs[0].end += maxN
		cRBG.refBuffs[0].lBuff += maxN
		return &cRBG
	}
	cRBG := rBG.copy()
	cRBG.refBuffs[len(cRBG.refBuffs)-1].rBuff += n
	cRBG.refBuffs[0].lBuff--
	cRBG.refBuffs[0].start++
	cRBG.refBuffs[0].end++
	cRBG.refBuffs[0].rBuff++
	return &cRBG
}

func (rBG refBuffGroup) addMax(n, passedLBuff int) *refBuffGroup {
	if len(rBG.brokenSpans) > 0 {
		maxN := fitAvailable(rBG.refBuffs[0], rBG.brokenSpans, n)
		if maxN == 0 {
			return nil
		}
		cRBG := rBG.copy()
		cRBG.refBuffs[0].rBuff += maxN
		return &cRBG
	}
	cRBG := rBG.copy()
	cRBG.refBuffs[len(cRBG.refBuffs)-1].rBuff += passedLBuff
	return &cRBG
}

func (rBG refBuffGroup) compare(other refBuffGroup) int {
	return slices.CompareFunc(rBG.refBuffs, other.refBuffs, func(a, b refBuff) int {
		return a.compare(b)
	})
}

func (rBG refBuffGroup) eq(other refBuffGroup) bool {
	if !slices.EqualFunc(rBG.brokenSpans, other.brokenSpans, bkSpansEqFunc) {
		return false
	}
	return slices.EqualFunc(rBG.refBuffs, other.refBuffs, func(a, b refBuff) bool { return a.eq(b) })
}

func (rBG refBuffGroup) copy() refBuffGroup {
	return refBuffGroup{
		refBuffs:    kit.Map(rBG.refBuffs, func(rB refBuff) refBuff { return rB.copy() }),
		brokenSpans: kit.Map(rBG.brokenSpans, func(s [2]int) [2]int { return [2]int{s[0], s[1]} }),
	}
}

func (rBG refBuffGroup) calcTotal() int {
	if len(rBG.brokenSpans) == 0 {
		return rBG.calcNonBrokenTotals()
	}
	return 1
}

func (rBG *refBuffGroup) leftAlign(lCutoff int) {
	for i := range rBG.refBuffs {
		rBG.refBuffs[i].leftAlign(lCutoff)
		lCutoff = rBG.refBuffs[i].end + 1
	}
}

type refBuffGroups []refBuffGroup

func (rBGs refBuffGroups) copy() refBuffGroups {
	c := make(refBuffGroups, 0, len(rBGs))
	for _, rBG := range rBGs {
		c = append(c, rBG.copy())
	}
	return c
}

func (rBGs refBuffGroups) calcTotal2() int {
	configTotal := make([]int, 0, len(rBGs))
	for _, rBG := range rBGs {
		configTotal = append(configTotal, rBG.calcTotal())
	}
	return kit.Mult(configTotal)
}

func (rBGs refBuffGroups) compare(other refBuffGroups) int {
	if len(rBGs) != len(other) {
		panic(rBGs)
	}
	return slices.CompareFunc(rBGs, other, func(a, b refBuffGroup) int {
		return a.compare(b)
	})
}

// equivalent means broken springs are in the same groups here
func (rBGs refBuffGroups) equivalent(other refBuffGroups) bool {
	if len(rBGs) != len(other) {
		return false
	}
	for i, rBG := range rBGs {
		if len(rBG.brokenSpans) != len(other[i].brokenSpans) {
			return false
		}
	}
	return true
}

func (rBGs refBuffGroups) leftAlign() {
	lCutoff := -1
	for i := range rBGs {
		rBGs[i].leftAlign(lCutoff)
		if i < len(rBGs)-1 &&
			len(rBGs[i+1].refBuffs) > 0 &&
			len(rBGs[i+1].refBuffs[0].brokenSpans) > 0 {
			lCutoff = rBGs[i].refBuffs[len(rBGs[i].refBuffs)-1].end + 1
			continue

		}
		if i > 0 && len(rBGs[i].brokenSpans) > 0 {
			rBGs[i-1].refBuffs[len(rBGs[i-1].refBuffs)-1].rBuff = rBGs[i].refBuffs[0].start - 2 - rBGs[i-1].refBuffs[len(rBGs[i-1].refBuffs)-1].end
		}
		lCutoff = rBGs[i].refBuffs[len(rBGs[i].refBuffs)-1].end +
			rBGs[i].refBuffs[len(rBGs[i].refBuffs)-1].rBuff + 1
	}
}

func (rBGs refBuffGroups) shiftGroups(cBs charsBrokenSpans, rBGIdx, upperCutoff int) []refBuffGroups {
	if rBGIdx == -1 {
		remainingBroken := rBGs.remainingBrokenSpans(cBs.brokenSpans)

		if len(remainingBroken) > 0 {
			return []refBuffGroups{}
		}
		lBuff := rBGs[0].refBuffs[0].start
		if len(rBGs[0].refBuffs[0].brokenSpans) > 0 {
			endDiff := rBGs[0].refBuffs[0].end - rBGs[0].refBuffs[0].brokenSpans[len(rBGs[0].refBuffs[0].brokenSpans)-1][1]
			lBuff = min(lBuff, endDiff)
		}
		if slices.ContainsFunc(rBGs, func(rBG refBuffGroup) bool {
			return len(rBG.refBuffs) == 0
		}) {
			return []refBuffGroups{}
		}
		rBGs[0].refBuffs[0].lBuff = lBuff
		rBGs.leftAlign()
		validCombo := rBGs.copy()
		return []refBuffGroups{validCombo}

	}
	validCombos := make([]refBuffGroups, 0, len(rBGs))
	nextRBGs := rBGs.copy()

	validCombos = append(
		validCombos,
		nextRBGs.shiftGroups(cBs, rBGIdx-1,
			nextRBGs[rBGIdx].refBuffs[0].start-1)...,
	)

	curRBG := nextRBGs[rBGIdx]
	passRefBuff := curRBG.refBuffs[len(curRBG.refBuffs)-1]
	passRefBuff.lBuff = 0
	passRefBuff.rBuff = 0

	refBuffDiff := passRefBuff.end - passRefBuff.start + 1
	lowerCutoff := passRefBuff.start

	remBroken := brokenSpansInRange(passRefBuff.end+1, upperCutoff, cBs.brokenSpans)
	if len(remBroken) > 0 {
		brokenEncompassed := contractSpansToDiff(remBroken, refBuffDiff)
		// for each span in broken encompassed:
		// 	- right align for all spans encompassed
		// 	- right align for just the last one
		// 	- return each of these results
		if len(brokenEncompassed) == 0 {
			// can't move group forward
			return validCombos
		}
		for spanIdx := len(brokenEncompassed) - 1; spanIdx > -1; spanIdx-- {
			bkRBGs := nextRBGs.copy()
			bkIdx := rBGIdx
			if spanIdx > 0 {
				lowerCutoff = brokenEncompassed[spanIdx-1][1] + 1
			}
			bkPassBuffRef := createRefBuff(
				createRefBuffInput{
					start:          brokenEncompassed[spanIdx][0],
					end:            brokenEncompassed[spanIdx][0] + refBuffDiff - 1,
					allBrokenSpans: brokenEncompassed[spanIdx:],
				},
			)
			shiftedBuffRef := bkPassBuffRef.shift(cBs.chars, lowerCutoff, upperCutoff)
			if shiftedBuffRef == nil {
				// as soon as it's bigger than the cutoff ranges,
				// subsequent ones will also be bigger than cutoff ranges
				break
			}
			nextGroup := refBuffGroup{
				refBuffs:    refBuffs{*shiftedBuffRef},
				brokenSpans: brokenEncompassed[spanIdx:],
			}
			if bkIdx < len(bkRBGs)-1 {
				// add lBuff to next one if there's room to
				bkRBGs[bkIdx+1].refBuffs[0].lBuff += upperCutoff - shiftedBuffRef.end - 1
			}
			// remove the refBuff and brokenSpan we're moving to its own group from the current
			bkRBGs[bkIdx].removeLastRefBuff()

			// remove refBuffGroup if there are no longer refBuffs in it
			if len(bkRBGs[bkIdx].refBuffs) == 0 {
				bkRBGs = append(bkRBGs[:bkIdx], bkRBGs[bkIdx+1:]...)[:len(bkRBGs)-1]
				bkIdx--
			}

			bkRBGs = append(
				bkRBGs[:bkIdx+1],
				append(refBuffGroups{nextGroup}, bkRBGs[bkIdx+1:]...)...,
			)

			validCombos = append(validCombos, bkRBGs.shiftGroups(cBs, bkIdx, shiftedBuffRef.start-1)...)
		}
		return validCombos
	}

	nextRBGs[rBGIdx].removeLastRefBuff()

	passRefBuff = createRefBuff(
		createRefBuffInput{
			start: upperCutoff - 1 - (refBuffDiff - 1),
			end:   upperCutoff - 1,
		},
	)
	passRefBuff.addBrokenSpans(cBs.brokenSpans)

	if len(passRefBuff.brokenSpans) == 0 &&
		rBGIdx < len(nextRBGs)-1 && len(nextRBGs[rBGIdx+1].brokenSpans) == 0 {
		// if the current passRefBuff has no brokenSpans,
		// and the subsequent rBG also has no broken spans,
		// add it to the subsequent rBG
		nextRBGs[rBGIdx+1].refBuffs = append(refBuffs{passRefBuff}, nextRBGs[rBGIdx+1].refBuffs...)
		if len(nextRBGs[rBGIdx].refBuffs) == 0 {
			// remove if there's no more refBuffs
			nextRBGs = append(nextRBGs[:rBGIdx], nextRBGs[rBGIdx+1:]...)
			rBGIdx--
		}
	} else {
		// add as a distinct refBuffGroup
		// because the subsequent one either does have broken spans, or doesn't exist
		// if the rBG at the rBGIdx ran out of refBuffs, decrement the rBGIdx
		// otherwise, the rBGIdx is still valid,
		// because it has dropped the refBuff that we've passed to a new RBG

		tempRBGs := refBuffGroups{
			refBuffGroup{
				refBuffs:    refBuffs{passRefBuff},
				brokenSpans: passRefBuff.brokenSpans,
			},
		}
		if len(nextRBGs[rBGIdx].refBuffs) > 0 {
			tempRBGs = append(refBuffGroups{nextRBGs[rBGIdx]}, tempRBGs...)
		}

		newTemp := append(nextRBGs[:rBGIdx].copy(), tempRBGs...)
		newTemp = append(newTemp, nextRBGs[rBGIdx+1:]...)
		nextRBGs = newTemp
		if len(tempRBGs) == 1 {
			rBGIdx--
		}
	}
	validCombos = append(validCombos, nextRBGs.shiftGroups(cBs, rBGIdx, passRefBuff.start-passRefBuff.lBuff-1)...)

	return validCombos
}

func (rBGs refBuffGroups) remainingBrokenSpans(brokenSpans [][2]int) [][2]int {
	allSepStrRefs := make(separatedStringRefs, 0, len(rBGs)*2)
	for _, rBG := range rBGs {
		allSepStrRefs = append(allSepStrRefs, rBG.sepStrRefs()...)
	}
	return allSepStrRefs.remainingBrokenSpans(brokenSpans)
}

func (rBG refBuffGroup) subtract(n int) *refBuffGroup {
	if len(rBG.brokenSpans) > 0 {
		if rBG.refBuffs[0].lBuff == 0 {
			if rBG.refBuffs[0].rBuff-n < 0 || rBG.refBuffs[0].start+n > rBG.brokenSpans[0][0] {
				return nil
			}
			cRBG := rBG.copy()
			cRBG.refBuffs[0].start += n
			cRBG.refBuffs[0].end += n
			cRBG.refBuffs[0].rBuff -= n
			return &cRBG

		}
		if rBG.refBuffs[0].lBuff-n <= 0 {
			return nil
		}
		cRBG := rBG.copy()
		cRBG.refBuffs[0].lBuff -= n
		cRBG.refBuffs[0].rBuff = max(cRBG.refBuffs[0].rBuff-n, 0)
		if cRBG.refBuffs[0].rBuff == 0 {
			return &cRBG
		}
		cRBG.refBuffs[0].start += n
		cRBG.refBuffs[0].end += n
		return &cRBG
	}

	if rBG.refBuffs[len(rBG.refBuffs)-1].rBuff-n < 0 {
		return nil
	}
	cRBG := rBG.copy()
	cRBG.refBuffs = kit.Map(cRBG.refBuffs, func(rB refBuff) refBuff {
		rB.start += n
		rB.end += n
		return rB
	})
	cRBG.refBuffs[len(cRBG.refBuffs)-1].rBuff -= n

	return &cRBG
}

func (rBGs refBuffGroups) calcTotals() []int {
	allGroups := rBGs.shiftBrokenGroups(len(rBGs) - 1)
	if !contains(allGroups, rBGs, func(a, b refBuffGroup) bool { return a.eq(b) }) {
		allGroups = append([]refBuffGroups{rBGs}, allGroups...)
	}
	totals := kit.Map(allGroups, func(rB refBuffGroups) int { return rB.calcTotal2() })
	return totals
}

// shift buffers (and start / end for broken groups) across groups where possible
func (rBGs refBuffGroups) shiftBrokenGroups(rBIdx int) []refBuffGroups {
	allGroups := make([]refBuffGroups, 0, len(rBGs))
	if rBIdx == -1 {
		return []refBuffGroups{rBGs}
	}

	nextTotals := make([]refBuffGroups, 0, 4)
	if len(rBGs[rBIdx].brokenSpans) == 0 {
		if rBIdx == 0 {
			return rBGs.shiftBrokenGroups(rBIdx - 1)
		}
		for i := 0; i < rBGs[rBIdx].refBuffs[len(rBGs[rBIdx].refBuffs)-1].rBuff+1; i++ {
			nextShifted := rBGs.copy()
			subtracted := nextShifted[rBIdx].subtract(i)
			if subtracted == nil {
				break
			}
			nextShifted[rBIdx] = *subtracted
			added := rBGs[rBIdx-1].add(i)
			if added == nil {
				break
			}
			nextShifted[rBIdx-1] = *added
			nextTotals = append(nextTotals, nextShifted.shiftBrokenGroups(rBIdx-1)...)
		}
		allGroups = append(allGroups, nextTotals...)
		return allGroups
	}

	for i := 0; i < rBGs[rBIdx].refBuffs[0].rBuff+1; i++ {
		nextShifted := rBGs.copy()
		subtracted := nextShifted[rBIdx].subtract(i)
		if subtracted == nil {
			break
		}
		nextShifted[rBIdx] = *subtracted

		var added *refBuffGroup
		if rBIdx > 0 {
			added = rBGs[rBIdx-1].addMax(i, subtracted.refBuffs[0].lBuff)
		}
		if added != nil {
			nextShifted[rBIdx-1] = *added
		}
		nextTotals = append(nextTotals, nextShifted.shiftBrokenGroups(rBIdx-1)...)
	}
	allGroups = append(allGroups, nextTotals...)

	return allGroups
}
