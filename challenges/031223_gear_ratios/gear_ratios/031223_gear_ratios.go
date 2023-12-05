package gear_ratios

import (
	"context"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

const INVALID_CHAR = '.'

func CalculatePartOne(ctx context.Context, input []string) int {
	ctx = log.Logger.WithContext(ctx)
	// l := log.Ctx(ctx).With().Caller().Logger()

	// assume line 1 is the num characters for rest of lines

	var result int
	for i, s := range input {
		// l.Info().Str("input_string", s).Msg("")
		var neighborLines []string
		if i != 0 {
			neighborLines = append(neighborLines, input[i-1])
		}
		if i != len(input)-1 {
			neighborLines = append(neighborLines, input[i+1])
		}

		lineSum := calculateLineSum(ctx, s, neighborLines)
		result += lineSum
		// log.Info().
		// 	Int("updated_result", result).Msg("valid game")
		// fmt.Println(result)
	}
	return result
}

// non nil from this means the neighbor character is valid
func validNeighborChar(char rune) *rune {
	fmt.Println(string(char))
	if char != INVALID_CHAR && !unicode.IsDigit(char) {
		return &char
	}
	return nil
}

func validCharFromNeighborLines(otherLine string, charRangeStart, charRangeEnd int) *rune {

	runesInStr := []rune(otherLine)

	if charRangeStart > 0 {
		validRune := validNeighborChar(runesInStr[charRangeStart-1])
		if validRune != nil {
			return validRune
		}
	}

	for i := charRangeStart; i <= charRangeEnd+1; i += 1 {
		if i == len(runesInStr) {
			break
		}
		validRune := validNeighborChar(runesInStr[i])
		if validRune != nil {
			return validRune
		}
	}

	return nil
}

func convertNumStrAddSum(numStr string) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic(err)
	}
	return num
}

// func processNumber

// two pointer / sliding window to grab numbers in a row, then compare neighbors
// below for iterating through strings by rune
// https://gobyexample.com/strings-and-runes#:~:text=To%20count%20how%20many%20runes,this%20count%20may%20be%20surprising.
func calculateLineSum(
	ctx context.Context,
	currentLine string,
	neighborLines []string,
) int {
	l := log.Ctx(ctx).With().Logger()

	runesInStr := []rune(currentLine)

	// not every character has the same number of bytes,
	// and because we're ranging thru them by bytes for the current line,
	// we need to track for the lines up and down from the current line what character we're on
	var currentCharIdx int
	var numStrStartIdx *int
	var numStrEndIdx *int

	var prevChar rune
	var totalForLine int

	for i, w := 0, 0; i < len(currentLine); i += w {
		currentCharIdx += 1

		currentChar, width := utf8.DecodeRuneInString(currentLine[i:])
		w = width
		// currentCharStr := string(currentChar)
		// fmt.Println(currentCharStr + " current character")

		if currentChar == INVALID_CHAR || !unicode.IsDigit(currentChar) {
			prevChar = currentChar
			continue
		}
		if numStrStartIdx == nil {
			// need the increment at the top of the loop to guarantee it happens,
			// but the real idx is 1 less than currenCharIdx
			numStrStartIdx = kit.Ptr(currentCharIdx - 1)
		}
		numStrEndIdx = kit.Ptr(currentCharIdx - 1)

		if i != len(currentLine) {
			nextChar, _ := utf8.DecodeRuneInString(currentLine[i+width:])
			if validNextChar := validNeighborChar(nextChar); validNextChar != nil {
				totalForLine += convertNumStrAddSum(
					string(runesInStr[*numStrStartIdx : *numStrEndIdx+1]),
				)
				l.Info().
					Str("valid_num_for_line", string(runesInStr[*numStrStartIdx:*numStrEndIdx+1])).
					Msg("found valid number")
				numStrStartIdx = nil
				numStrEndIdx = nil
				continue
			}
			if unicode.IsDigit(nextChar) {
				continue
			}
		}

		if i != 0 {
			if validPrevChar := validNeighborChar(prevChar); validPrevChar != nil {
				totalForLine += convertNumStrAddSum(
					string(runesInStr[*numStrStartIdx : *numStrEndIdx+1]),
				)
				l.Info().
					Str("valid_num_for_line", string(runesInStr[*numStrStartIdx:*numStrEndIdx+1])).
					Msg("found valid number")
				numStrStartIdx = nil
				continue
			}
		}
		if numStrStartIdx == nil {
			continue
		}

		for _, neighborLine := range neighborLines {
			validCharFromLine := validCharFromNeighborLines(neighborLine, *numStrStartIdx, *numStrEndIdx)
			if validCharFromLine != nil {
				totalForLine += convertNumStrAddSum(
					string(runesInStr[*numStrStartIdx : *numStrEndIdx+1]),
				)
				l.Info().
					Str("valid_num_for_line", string(runesInStr[*numStrStartIdx:*numStrEndIdx+1])).
					Msg("found valid number")
				numStrStartIdx = nil
				break
			}
		}
		if numStrStartIdx != nil && numStrEndIdx != nil {
			l.Warn().Str("invalid_number", string(runesInStr[*numStrStartIdx:*numStrEndIdx+1])).Msgf("invalid number")
			if string(runesInStr[*numStrStartIdx:*numStrEndIdx+1]) == "691" {
				fmt.Println("panic is starting")
			}
		}
		// no valid neighbor characters found, reset
		numStrStartIdx = nil
	}
	return totalForLine
}
