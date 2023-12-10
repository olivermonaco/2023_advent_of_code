package seed_map

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

type mapping struct {
	sourceStart, destStart, rangeVal int
}

func (m mapping) mappedValue(sourceVal int) *int {
	if sourceVal >= m.sourceStart && sourceVal <= m.sourceStart+m.rangeVal-1 {
		return kit.Ptr(m.destStart + (sourceVal - m.sourceStart))
	}
	return nil
}

type mappings []mapping

func (mps mappings) toDest(sourceVal int) int {
	for _, mapping := range mps {
		if destVal := mapping.mappedValue(sourceVal); destVal != nil {
			return *destVal
		}
	}
	return sourceVal
}

type translator struct {
	sourceName, destName string
	mappings
}

func (t translator) translate(sourceVal int, sourceName string) (int, string) {
	if sourceName != t.sourceName {
		panic(
			fmt.Errorf(
				"translator recieved source name %s, but expected source name %s",
				sourceName, t.sourceName,
			),
		)
	}
	return t.mappings.toDest(sourceVal), t.destName
}

func sourceDestNames(line string) (*string, *string) {
	sourceName, destNamePlus, found := strings.Cut(line, "-to-")
	if found {
		destName := strings.Fields(destNamePlus)[0]
		return &sourceName, &destName
	}
	return nil, nil
}

func buildTranslators(lines []string) map[string]translator {
	var sourceName, destName *string

	translatorMap := make(map[string]translator, 7)
	var t translator
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
			t = translator{
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

		t.mappings = append(
			t.mappings,
			mapping{
				destStart:   rangeValues[0],
				sourceStart: rangeValues[1],
				rangeVal:    rangeValues[2],
			},
		)
		if idx+1 == len(lines) || len(lines[idx+1]) == 0 {
			translatorMap[t.sourceName] = t
		}
	}
	return translatorMap
}

func extractNumsCutLines(lines []string) ([]int, []string) {
	for idx, line := range lines {
		seedsStr, found := strings.CutPrefix(line, "seeds:")
		if found {
			// cut the input lines at the index
			lines = lines[idx+1:]

			seedsSlice := strings.Fields(seedsStr)
			var seeds []int
			for _, seedStr := range seedsSlice {
				seed, err := strconv.Atoi(seedStr)
				if err != nil {
					panic(seed)
				}
				seeds = append(seeds, seed)
			}
			return seeds, lines
		}
	}
	return nil, nil
}

func getLocation(seed int, translatorMap map[string]translator) int {
	sourceName := "seed"
	sourceVal := seed
	for {
		translator, ok := translatorMap[sourceName]
		if !ok {
			panic(sourceName)
		}
		sourceVal, sourceName = translator.translate(sourceVal, sourceName)
		if sourceName == "location" {
			return sourceVal
		}
	}
}

func CalculatePartOne(ctx context.Context, input []string) int {
	l := log.Ctx(ctx).With().Logger()

	seeds, lines := extractNumsCutLines(input)
	translatorMap := buildTranslators(lines)
	var minimumLocation *int
	for _, seed := range seeds {
		seedLocation := getLocation(seed, translatorMap)
		l.Info().Int("seed_location", seedLocation).Msg("")
		if minimumLocation == nil || seedLocation < *minimumLocation {
			minimumLocation = &seedLocation
		}
	}
	return *minimumLocation
}
