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
	"github.com/rs/zerolog/log"
)

// this challenge was interesting, so I looked up some reddit threads on it after, saw this approach, and decided to try it out
func CalculatePartTwo_ReverseLocationConcurrent(ctx context.Context, input []string) int {
	l := log.Ctx(ctx).With().Logger()

	seedsStartsRanges, lines := extractNumsCutLines(input)
	seedsBoundsSourceSl := buildSeedsBoundsSlice(seedsStartsRanges)

	translatorDestMap := buildTranslatorsBoundsDestKey(lines)
	maxGRs := runtime.NumCPU()
	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan struct{}, maxGRs)

	l.Info().Int("max_goroutines", maxGRs).Msg("using max goroutines")

	i := 0
	var validLocations []int
	for {
		wg.Add(1)
		ch <- struct{}{}
		go func(i int) {
			idx := i
			defer func() { <-ch }()
			defer wg.Done()
			potentialSeed := calcSeedValue(idx, translatorDestMap)
			validSeed := findValidSeed(seedsBoundsSourceSl, potentialSeed)
			if validSeed != nil {
				l.Info().
					Int("valid_seed", *validSeed).
					Msg("found valid seed")
				mu.Lock()
				{
					validLocations = append(validLocations, idx)
				}
				mu.Unlock()
				return
			}
			if i%5 == 0 {
				l.Info().
					Int("tried_loc", i).
					Msg("tried location")
			}
		}(i)
		if len(validLocations) > 0 {
			break
		}
		i++
	}
	wg.Wait()
	close(ch)
	minLocation := validLocations[0]
	for _, location := range validLocations {
		if location < minLocation {
			minLocation = location
		}
	}
	return minLocation
}

// func CalculatePartTwo_ReverseLocation(ctx context.Context, input []string) int {
// 	l := log.Ctx(ctx).With().Logger()

// 	seedsStartsRanges, lines := extractNumsCutLines(input)
// 	seedsBoundsSourceSl := buildSeedsBoundsSlice(seedsStartsRanges)

// 	translatorDestMap := buildTranslatorsBoundsDestKey(lines)
// 	i := 0
// 	for {
// 		potentialLoc := calcSeedValue(i, translatorDestMap)
// 		validSeed := findValidSeed(seedsBoundsSourceSl, potentialLoc)
// 		if validSeed != nil {
// 			l.Info().
// 				Int("valid_seed", *validSeed).
// 				Msg("found valid seed")
// 			return i
// 		}
// 		i++
// 		if i%5 == 0 {
// 			l.Info().
// 				Int("tried_loc", i).
// 				Msg("tried location")
// 		}
// 	}
// }

func calcSeedValue(potentialLoc int, translatorDestMap map[string]translatorBounds) int {
	destName := "location"
	for {
		translator, ok := translatorDestMap[destName]
		if !ok {
			panic(fmt.Sprintf("destName %s, i is %d", destName, potentialLoc))
		}
		potentialLoc, destName = translator.toSource(potentialLoc)
		if destName == "seed" {
			break
		}
	}
	return potentialLoc
}

func findValidSeed(seedsBoundsSourceSl []seedsBounds, i int) *int {
	var validSeed *int
	for _, seed := range seedsBoundsSourceSl {
		if seed.start <= i && seed.end >= i {
			return kit.Ptr(i)
		}
	}
	return validSeed
}

func (t translatorBounds) toSource(i int) (int, string) {
	for _, mPsBs := range t.mappingsBounds {
		if mPsBs.start <= i && i <= mPsBs.end {
			i += mPsBs.offset
			return i, t.sourceName
		}
	}
	return i, t.sourceName
}

func buildTranslatorsBoundsDestKey(lines []string) map[string]translatorBounds {
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
					start: rangeValues[0],
					end:   rangeValues[0] + rangeValues[2] - 1,
				},
				offset: rangeValues[1] - rangeValues[0],
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
			translatorMap[t.destName] = t
		}
	}

	return translatorMap
}
