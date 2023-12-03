package trebuchet

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

var strRepToInt = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func Calculate(input []string, partFunc func(string) int) int {
	var result int
	for _, s := range input {
		intResult := partFunc(s)
		result += intResult
		log.Info().
			Str("processed_str", s).
			Int("int_from_runes", intResult).
			Int("updated_result", result).Msg("")
	}
	return result
}

func IntFromStrPartOne(s string) int {
	firstRune, _ := numRuneInStr(s)
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", s)
		panic(s)
	}
	lastRune, _ := numRuneInStr(kit.ReverseString(s))
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", kit.ReverseString(s))
		panic(s)
	}
	runes := []rune{*firstRune, *lastRune}
	result, err := strconv.Atoi(string(runes))
	if err != nil {
		log.Err(err).Msg("error")
		return 0
	}
	return result
}

func IntFromStrPartTwo(s string) int {
	firstRune, firstRuneIdx := numRuneInStr(s)
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", s)
	}
	lastRune, lastRuneIdx := numRuneInStr(kit.ReverseString(s))
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", kit.ReverseString(s))
	}

	firstStrMatch, firstStrRepIdx := strRepInStr(s, func(s string) string { return s })
	if firstStrRepIdx >= 0 && firstStrRepIdx <= firstRuneIdx {
		firstRune = kit.Ptr(
			runeFromInt(firstStrMatch),
		)
	}

	lastStrMatch, lastStrRepIdx := strRepInStr(
		kit.ReverseString(s), func(s string) string { return kit.ReverseString(s) },
	)
	if lastStrRepIdx >= 0 && lastStrRepIdx <= lastRuneIdx {
		lastRune = kit.Ptr(
			runeFromInt(lastStrMatch),
		)
	}

	runes := []rune{*firstRune, *lastRune}
	result, err := strconv.Atoi(string(runes))
	if err != nil {
		log.Err(err).Msg("error")
		return 0
	}
	return result
}

func numRuneInStr(s string) (*rune, int) {
	for idx, r := range s {
		_, err := strconv.Atoi(string(r))
		if err != nil {
			continue
		}
		return &r, idx
	}
	return nil, 0
}

func strRepInStr(s string, compare_direction func(string) string) (int, int) {
	minIdx := len(s) // len of bytes, but that's fine, just need it to be bigger than the number of runes
	returnInt := -1
	for strRep, v := range strRepToInt {
		strRep = compare_direction(strRep)
		if idx := strings.Index(s, strRep); idx >= 0 {
			if idx < minIdx {
				minIdx = idx
				returnInt = v
			}
		}
	}
	if returnInt == -1 {
		return -1, -1
	}
	return returnInt, minIdx
}

// not really safe or checking anything, but moving fast
func runeFromInt(i int) rune {
	numRep := strconv.Itoa(i)
	return []rune(numRep)[0]
}

func ReadFileConstructLines(filename string) []string {
	log.Info().Msgf("filename is %s", filename)
	file, err := os.Open(filename)
	if err != nil {

		log.Err(err).Msg("Error opening file")
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Err(err).Msg("Error reading file")
		return nil
	}
	return lines
}
