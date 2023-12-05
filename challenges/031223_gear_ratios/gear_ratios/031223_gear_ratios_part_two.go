package gear_ratios

import (
	"context"
	"strconv"
	"unicode"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

const GEAR_CHAR = '*'

// shoulda been finding the non digit, non . chars all along...
func CalculatePartTwo(ctx context.Context, input []string) int {
	ctx = log.Logger.WithContext(ctx)

	var result int
	for i, s := range input {
		var neighborLines []string
		if i != 0 {
			neighborLines = append(neighborLines, input[i-1])
		}
		if i != len(input)-1 {
			neighborLines = append(neighborLines, input[i+1])
		}
		result += calculatePartTwoLineSum(ctx, s, neighborLines)
	}
	return result
}

func neighborDigit(char rune) *rune {
	if unicode.IsDigit(char) {
		return &char
	}
	return nil
}

func getCompleteNumstr(runesStr []rune, idxStart int) string {
	var num []rune
	for i := idxStart; i <= len(runesStr)-1; i += 1 {
		if !unicode.IsDigit(runesStr[i]) {
			break
		}
		num = append(num, runesStr[i])
	}

	for i := idxStart - 1; i >= 0; i -= 1 {
		if !unicode.IsDigit(runesStr[i]) {
			break
		}
		num = append([]rune{runesStr[i]}, num...)
	}
	return string(num)
}

func getNeighborIdxs(
	runesNeighborLine []rune,
	lineIdxToPositionIdx []linePositionIdx,
	charIdx int,
) []linePositionIdx {
	lPIdxs := make([]linePositionIdx, 0, 3)

	if charIdx != 0 {
		prevDigit := neighborDigit(runesNeighborLine[charIdx-1])
		if prevDigit != nil {
			lPIdxs = append(
				lPIdxs,
				linePositionIdx{
					line:        &runesNeighborLine,
					positionIdx: -1,
				},
			)
		}
	}

	if charIdx != len(runesNeighborLine)-1 {
		next := neighborDigit(runesNeighborLine[charIdx+1])
		if next != nil {
			lPIdxs = append(
				lPIdxs,
				linePositionIdx{
					line:        &runesNeighborLine,
					positionIdx: 1,
				},
			)
		}
	}

	same := neighborDigit(runesNeighborLine[charIdx])
	if same != nil {
		lPIdxs = append(
			lPIdxs,
			linePositionIdx{
				line:        &runesNeighborLine,
				positionIdx: 0,
			},
		)
	}
	if len(lPIdxs) <= 1 {
		return lPIdxs
	}
	if lPIdxs[len(lPIdxs)-1].positionIdx-lPIdxs[0].positionIdx > 1 {
		if len(lPIdxs) == 3 {
			// all three positions are in a row, only one number to count
			// start from the middle when searching
			return []linePositionIdx{lPIdxs[1]}
		}
		// prev and next idxs have digits, but not the middle. counts as 2
		return []linePositionIdx{lPIdxs[0], lPIdxs[len(lPIdxs)-1]}
	}
	// two in a row, so just need to return one of them
	return []linePositionIdx{lPIdxs[0]}
}

type linePositionIdx struct {
	line        *[]rune
	positionIdx int
}

func calculatePartTwoLineSum(
	ctx context.Context,
	currentLine string,
	neighborLines []string,
) int {
	// l := log.Ctx(ctx).With().Logger()
	runesInStr := []rune(currentLine)

	var lineTotal int

	for idx, char := range runesInStr {
		if char != GEAR_CHAR {
			continue
		}
		lPIdxs := make([]linePositionIdx, 0, 8)

		if idx != 0 {
			prev := neighborDigit(runesInStr[idx-1])
			if prev != nil {
				lPIdxs = append(
					lPIdxs,
					linePositionIdx{
						line:        kit.Ptr([]rune(currentLine)),
						positionIdx: -1,
					},
				)
			}
		}

		if idx != len(runesInStr)-1 {
			next := neighborDigit(runesInStr[idx+1])
			if next != nil {
				lPIdxs = append(
					lPIdxs,
					linePositionIdx{
						line:        kit.Ptr([]rune(currentLine)),
						positionIdx: 1,
					},
				)
			}
		}

		for _, neighborLine := range neighborLines {
			neighborLineRunes := []rune(neighborLine)
			lPIdxs = append(
				lPIdxs,
				getNeighborIdxs(
					neighborLineRunes,
					lPIdxs,
					idx,
				)...,
			)
		}
		if len(lPIdxs) > 2 {
			// too many adjacent digits
			continue
		}
		nums := make([]int, 0, 2)
		for _, lineIdxPositionIdx := range lPIdxs {
			numStr := getCompleteNumstr(
				*lineIdxPositionIdx.line,
				idx+lineIdxPositionIdx.positionIdx,
			)
			if numStr == "" {
				continue
			}
			num, err := strconv.Atoi(numStr)
			if err != nil {
				panic(numStr)
			}
			nums = append(nums, num)
		}
		if len(nums) != 2 {
			continue
		}
		gearTotal := 1
		for _, num := range nums {
			gearTotal *= num
		}
		lineTotal += gearTotal
	}
	return lineTotal
}
