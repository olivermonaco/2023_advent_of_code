package gear_ratios

import (
	"context"
	"unicode"

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
	}
}

func neighborDigit(char rune) *rune {
	if unicode.IsDigit(char) {
		return &char
	}
	return nil
}

func getNeighborIdxs(
	runesNeighborLine []rune,
	lineIdxToPositionIdx []linePositionIdx,
	lineIdx, charIdx int,
) []linePositionIdx {
	var lPIdxs []linePositionIdx

	if charIdx != 0 {
		prevDigit := neighborDigit(runesNeighborLine[charIdx-1])
		if prevDigit != nil {
			lPIdxs = append(
				lPIdxs,
				linePositionIdx{
					lineIdx:     lineIdx,
					positionIdx: -1,
				},
			)
		}
	}

	if charIdx != len(runesNeighborLine)-1 {
		next := neighborDigit(runesNeighborLine[charIdx-1])
		if next != nil {
			lPIdxs = append(
				lPIdxs,
				linePositionIdx{
					lineIdx:     lineIdx,
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
				lineIdx:     lineIdx,
				positionIdx: 0,
			},
		)
	}
	if len(lPIdxs) == 1 {
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
	return []linePositionIdx{}
}

type linePositionIdx struct {
	lineIdx     int
	positionIdx int
}

func calculatePartTwoLineSum(
	ctx context.Context,
	currentLine string,
	neighborLines []string,
) int {
	runesInStr := []rune(currentLine)
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
						lineIdx:     -1,
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
						lineIdx:     -1,
						positionIdx: 1,
					},
				)
			}
		}

		for neighborLineIdx, neighborLine := range neighborLines {
			neighborLineRunes := []rune(neighborLine)
			lPIdxs = append(
				lPIdxs,
				getNeighborIdxs(
					neighborLineRunes,
					lPIdxs,
					neighborLineIdx,
					idx,
				)...,
			)
		}
		if len(lPIdxs) > 2 {
			// too many adjacent digits
			continue
		}
		// TODO: Left off here, navigate out for each digit
	}
}
