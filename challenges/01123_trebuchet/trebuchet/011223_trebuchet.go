package trebuchet

import (
	"bufio"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func CalculatePartOne(input []string) int {
	var result int
	for _, s := range input {
		intResult := intFromStr(s)
		result += intResult
		log.Info().
			Str("processed_str", s).
			Int("int_from_runes", intResult).
			Int("updated_result", result).Msg("")
	}
	return result
}

func intFromStr(s string) int {
	firstRune := numRuneInString(s)
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", s)
	}
	lastRune := numRuneInString(reverseString([]rune(s)))
	if firstRune == nil {
		log.Error().Msgf("couldn't find rune in string: %s", reverseString([]rune(s)))
	}
	runes := []rune{*firstRune, *lastRune}
	result, err := strconv.Atoi(string(runes))
	if err != nil {
		log.Err(err).Msg("error")
		return 0
	}
	return result
}

func reverseString(s []rune) string {
	var reversed []rune

	for i := len(s) - 1; i > -1; i -= 1 {
		reversed = append(reversed, s[i])
	}
	reversedS := string(reversed)
	return reversedS

}

func numRuneInString(s string) *rune {
	for _, r := range s {
		_, err := strconv.Atoi(string(r))
		if err != nil {
			continue
		}
		return &r
	}
	return nil
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
