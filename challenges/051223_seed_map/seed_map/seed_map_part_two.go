package seed_map

import (
	"cmp"
	"context"
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// seedsBounds represents a lower / upper bound of a contiguous seed value
// inclusive of start / end values
type seedsBounds struct {
	start, end int
}

func (sBs seedsBounds) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int("seeds_start", sBs.start).
		Int("seeds_end", sBs.end)
}

type seedsBondsSl []seedsBounds

func (sBsSl seedsBondsSl) MarshalZerologArray(a *zerolog.Array) {
	for _, sBs := range sBsSl {
		a.Object(sBs)
	}
}

// Intersection only cares if the receiver seedsBounds intersects with the other seedsBounds
func (sBs seedsBounds) Intersection(other seedsBounds) *seedsBounds {
	if sBs.end < other.start || sBs.start > other.end {
		return nil
	}
	return &seedsBounds{
		start: max(sBs.start, other.start),
		end:   min(sBs.end, other.end),
	}
}

// Pseudo code for range calculation
// for each seed range:
// - get the lower / upper bounds of the seed range (slice of 1)
// - the translators now have a source value lower / upper bound range for each mapping
// 		- each mapping has a + / - index amount based on the dest relationship
// 		- send the lowest / highest of each tuple to the next translator
// - pass these into the translator, and have the translator divide those up
//   into a slice of lower / upper bounds. Each of these will be fed into the next translator,
//   until the final translator will spit out a slice of lower / upper bounds

func CalculatePartTwo_CheckBounds(ctx context.Context, input []string) int {
	l := log.Ctx(ctx).With().Logger()

	seedsStartsRanges, lines := extractNumsCutLines(input)
	seedsBoundsSourceSl := buildSeedsBoundsSlice(seedsStartsRanges)

	translatorMap := buildTranslatorsBounds(lines)
	var minOverallLocation *int
	for _, sBsSource := range seedsBoundsSourceSl {
		locationBounds := getLocationBounds(ctx, []seedsBounds{sBsSource}, translatorMap)
		l.Info().Int("start_num", sBsSource.start).Msg("start Num")
		l.Info().
			Int("loc_bound", locationBounds[0].start).
			Msg("loc bound")

		if minOverallLocation == nil || *minOverallLocation > locationBounds[0].start {
			minOverallLocation = &locationBounds[0].start
			if kit.Deref(minOverallLocation) == 0 {
				fmt.Println("why zero tho")
			}
			l.Info().Int("min_loc", *minOverallLocation).Msg("min loc overall lowered")
		}
	}
	return *minOverallLocation
}

func buildTranslatorsBounds(lines []string) map[string]translatorBounds {
	var sourceName, destName *string

	translatorMap := make(map[string]translatorBounds, 7) // number of gardening translations
	var t translatorBounds
	for idx, line := range lines {
		if len(line) == 0 {
			sourceName, destName = nil, nil
			continue
		}
		if sourceName == nil || destName == nil {
			sourceName, destName = sourceDestNames(line)
			if sourceName == nil || destName == nil {
				panic(line)
			}
			t = translatorBounds{
				sourceName: *sourceName,
				destName:   *destName,
			}
			continue
		}
		rangeStrValues := strings.Fields(line)
		rangeValues := make([]int, 0, len(rangeStrValues))
		for _, rangeStr := range rangeStrValues {
			val, err := strconv.Atoi(rangeStr)
			if err != nil {
				panic(err)
			}
			rangeValues = append(rangeValues, val)
		}

		t.mappingsBounds = append(
			t.mappingsBounds,
			mappingBounds{
				seedsBounds: seedsBounds{
					start: rangeValues[1],
					end:   rangeValues[1] + rangeValues[2] - 1,
				},
				offset: rangeValues[0] - rangeValues[1],
			},
		)
		// current line could be the end of the file,
		// or next line could be an empty str
		if idx+1 == len(lines) || len(lines[idx+1]) == 0 {
			// sort by source start
			slices.SortStableFunc(
				t.mappingsBounds,
				func(a mappingBounds, b mappingBounds) int {
					return cmp.Compare(a.seedsBounds.start, b.seedsBounds.start)
				},
			)
			translatorMap[t.sourceName] = t
		}
	}

	return translatorMap
}

type mappingBounds struct {
	seedsBounds
	offset int
}

func (mBs mappingBounds) mappedValue(sourceVal int) *int {
	if sourceVal >= mBs.seedsBounds.start && sourceVal <= mBs.seedsBounds.end {
		ret := sourceVal + mBs.offset
		// if ret == 0 {
		// 	fmt.Println("\nzero found")
		// }
		// if ret == 34147766 {
		// 	fmt.Println("\n 34147766 found")
		// }
		return kit.Ptr(ret)
	}
	return nil
}

func (mBs mappingBounds) toDest(sourceVal int) int {
	if destVal := mBs.mappedValue(sourceVal); destVal != nil {
		return *destVal
	}
	return sourceVal
}

type mappingsBounds []mappingBounds

func (mpsBs mappingsBounds) toDestBounds(
	sourceVals []seedsBounds,
	sourceIdx, compareIdx int,
) ([]seedsBounds, int) {
	if sourceIdx == len(sourceVals) {
		return []seedsBounds{}, compareIdx
	}

	// stay on the last idx if we're in this function,
	// because it means sourceIdx hasn't reached the end yet
	compareIdx = min(compareIdx, len(mpsBs)-1)

	compareBounds := mpsBs[compareIdx]
	sourceBounds := sourceVals[sourceIdx]
	// source entirely greater than compare
	if sourceBounds.start > compareBounds.end {
		return []seedsBounds{}, compareIdx + 1
	}

	// compare entirely greater than source (don't increment)
	if compareBounds.start > sourceBounds.end {
		return []seedsBounds{}, compareIdx
	}
	// there's going to be an intersection if we get here
	var destSeedsBounds []seedsBounds
	intersection := sourceBounds.Intersection(compareBounds.seedsBounds)
	if intersection == nil {
		panic([]seedsBounds{sourceBounds, compareBounds.seedsBounds})
	}

	// check if there's a lower bound below the intersection
	if sourceBounds.start < compareBounds.start {
		destSeedsBounds = append(
			destSeedsBounds,
			seedsBounds{
				start: sourceBounds.start,
				end:   compareBounds.start - 1,
			},
		)
	}

	// add the intersection
	destSeedsBounds = append(
		destSeedsBounds,
		seedsBounds{
			start: compareBounds.toDest(intersection.start),
			end:   compareBounds.toDest(intersection.end),
		},
	)

	// if there's no upper bound above the intersection,
	// return and increment the compare
	if sourceBounds.end < compareBounds.end {
		return destSeedsBounds, compareIdx + 1
	}

	if compareIdx == len(mpsBs)-1 {
		// if it's the last in the mappings, just return what we have plus
		// the source bound above the last compare
		destSeedsBounds = append(
			destSeedsBounds,
			seedsBounds{
				start: compareBounds.end + 1,
				end:   sourceBounds.end,
			},
		)
		return destSeedsBounds, compareIdx + 1
	}
	// below here is a recursion, based on the scenario:
	// - the sourceBound end is greater than the compareBound end
	// - therefore, there might be another compareBound we need to compare beyond the current
	// - so, construct a a new slice of sourceVals based on 1 greater than the current compareBound end,
	// 	 and increment the compareIdx to look at the next compare
	nextSourceVals := append(
		[]seedsBounds{
			{
				start: compareBounds.end + 1,
				end:   sourceBounds.end,
			},
		},
		sourceVals[sourceIdx+1:]...,
	)
	nextDestBounds, nextCompareIdx := mpsBs.toDestBounds(
		nextSourceVals,
		0,
		compareIdx+1,
	)
	destSeedsBounds = append(
		destSeedsBounds,
		nextDestBounds...,
	)
	compareIdx = nextCompareIdx
	return destSeedsBounds, compareIdx
}

type translatorBounds struct {
	sourceName, destName string
	mappingsBounds
}

// coalesceSeedsBonds requires vals to be sorted by start,
// and have no duplicate values (eg. [(1, 5), (1, 5)] is invalid)
// coalesceSeedsBonds compacts the vals input so overlapping ranges coalesce
//
// eg. [(0, 10), (2, 12), (14, 20)] -> [(0, 12), (14, 20)]
func coalesceSeedsBonds(vals []seedsBounds) []seedsBounds {
	retSeedsBonds := make([]seedsBounds, 0, len(vals))
	if len(vals) == 0 {
		return vals
	}

	curSeedBound := vals[0]

	for i, sBs := range vals {

		if sBs.start <= curSeedBound.end+1 {
			curSeedBound = seedsBounds{
				start: curSeedBound.start,
				end:   sBs.end,
			}
		}
		if sBs.start > curSeedBound.end+1 {
			retSeedsBonds = append(retSeedsBonds, curSeedBound)
			curSeedBound = sBs
		}
		if i == len(vals)-1 {
			retSeedsBonds = append(retSeedsBonds, curSeedBound)
		}
	}
	return retSeedsBonds
}

func (t translatorBounds) translateBounds(sourceVals []seedsBounds, sourceName string) ([]seedsBounds, string) {
	if sourceName != t.sourceName {
		panic(
			fmt.Errorf(
				"translator recieved source name %s, but expected source name %s",
				sourceName, t.sourceName,
			),
		)
	}
	var destBounds []seedsBounds
	var sourceIdx int
	var compareIdx int
	for {
		nextSeedsBounds, newCompareIdx := t.mappingsBounds.toDestBounds(
			sourceVals,
			sourceIdx,
			compareIdx,
		)

		destBounds = append(destBounds, nextSeedsBounds...)

		if newCompareIdx > compareIdx {
			compareIdx = newCompareIdx
			continue
		}

		if len(destBounds) == 0 {
			// source vals were entirely above all of the compare values
			destBounds = append(destBounds, sourceVals[sourceIdx])
		}

		slices.SortStableFunc(
			destBounds,
			func(a, b seedsBounds) int {
				return cmp.Compare(a.start, b.start)
			},
		)
		destBounds = slices.CompactFunc(
			destBounds,
			func(a, b seedsBounds) bool {
				if a.start != b.start {
					return false
				}
				if a.end != b.end {
					return false
				}
				return true
			})
		destBounds = coalesceSeedsBonds(destBounds)
		// fmt.Printf("\nafter coalescing: %v\n", destBounds)

		sourceIdx++
		if sourceIdx >= len(sourceVals) {
			break
		}
	}
	return destBounds, t.destName
}

type seedsRanges struct {
	start, seedRange int
}

func (sRs seedsRanges) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int("seed_start_idx", sRs.start).
		Int("seed_range_from_start", sRs.seedRange)
}
func buildSeedsBoundsSlice(startsRanges []int) []seedsBounds {
	seedsBoundsSourceSl := make([]seedsBounds, 0, (len(startsRanges)/2)+1)

	for i, startOrRange := range startsRanges {
		if i%2 == 1 {
			seedsBoundsSourceSl = append(
				seedsBoundsSourceSl,
				seedsBounds{
					start: startsRanges[i-1],
					end:   startsRanges[i-1] + startOrRange - 1,
				},
			)
		}
	}
	slices.SortFunc(seedsBoundsSourceSl,
		func(a, b seedsBounds) int { return cmp.Compare(a.start, b.start) },
	)
	return seedsBoundsSourceSl
}

func getLocationBounds(
	ctx context.Context,
	sourceBounds []seedsBounds,
	translatorMap map[string]translatorBounds,
) []seedsBounds {
	l := log.Ctx(ctx).With().Logger()

	sourceName := "seed"
	curSourceBoundsRanges := sourceBounds
	for {
		translator, ok := translatorMap[sourceName]
		if !ok {
			panic(sourceName)
		}
		curSourceBoundsRanges, sourceName = translator.translateBounds(curSourceBoundsRanges, sourceName)
		l.Info().
			Str("new_source_name", sourceName).
			Int("len_curboundranges", len(curSourceBoundsRanges)).
			Msg("")
		// l.Info().
		// 	Array("source_bounds_ranges", seedsBondsSl(curSourceBoundsRanges)).
		// 	Str("new_source_name", sourceName).
		// 	Msg("new ranges calculated")
		if sourceName == "location" {
			return curSourceBoundsRanges
		}
	}
}

// this isn't right...
func CalculatePartTwo_BruteForceConcurrent(ctx context.Context, input []string) int {
	l := log.Ctx(ctx).With().Logger()

	seedsStartsRanges, lines := extractNumsCutLines(input)
	seedsRangesSlice := buildSeedsRangesSlice(seedsStartsRanges)

	translatorMap := buildTranslators(lines)
	var minOverallLocation *int
	for _, sRs := range seedsRangesSlice {
		minLocForRange := minLocationForRange(ctx, sRs, translatorMap)
		if minOverallLocation == nil || minLocForRange < *minOverallLocation {
			minOverallLocation = &minLocForRange
			l.Info().Int("min_loc", minLocForRange).Msg("min loc overall lowered")
		}
	}
	return *minOverallLocation
}

func buildSeedsRangesSlice(startsRanges []int) []seedsRanges {
	seedsRangesSourceSl := make([]seedsRanges, 0, (len(startsRanges)/2)+1)

	for i, startRange := range startsRanges {
		if i%2 == 1 {
			seedsRangesSourceSl = append(
				seedsRangesSourceSl,
				seedsRanges{
					start:     startsRanges[i-1],
					seedRange: startRange,
				},
			)
		}
	}
	return seedsRangesSourceSl
}

func minLocationForRange(ctx context.Context, sRs seedsRanges, translatorMap map[string]translator) int {
	l := log.Ctx(ctx).With().Logger()

	maxGRs := runtime.NumCPU()
	var wg sync.WaitGroup

	var mu sync.Mutex
	ch := make(chan struct{}, maxGRs)
	l.Info().Int("max_goroutines", maxGRs).Msg("using max goroutines")

	var minSeedLoc *int

	// l.Info().Object("seed", sRs).Msg("Processing seed")

	for seedIdx := 0; seedIdx <= sRs.seedRange-1; seedIdx++ {
		wg.Add(1)
		seed := sRs.start + seedIdx
		ch <- struct{}{}
		go func(seed int) {
			defer func() { <-ch }()
			defer wg.Done()
			// below blocks goroutines until the goroutine receives work to do
			seedLoc := getLocation(seed, translatorMap)
			mu.Lock()
			defer mu.Unlock()
			if minSeedLoc == nil || seedLoc < *minSeedLoc {
				minSeedLoc = &seedLoc
				l.Info().Int("seed_location", seedLoc).Msg("new min seed location")
			}

		}(seed)
	}
	wg.Wait()
	close(ch)

	return *minSeedLoc
}
