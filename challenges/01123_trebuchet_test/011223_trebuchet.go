package trebuchet

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(os.Stdout).With().Caller().Logger()
}

func Trebuchet(input []string) int {
	var result int
	for _, s := range input {
		result += processString(s)
	}
	return result
}

func processString(s string) int {
	firstRune := numRuneInString(s)
	lastRune := numRuneInString(reverseString([]rune(s)))
	int, err := strconv.Atoi(string([]rune{*firstRune, *lastRune}))
	if err != nil {
		log.Err(err).Msg("error")
		return 0
	}
	return int
}

func reverseString(s []rune) string {
	var reversed []rune

	for i := len(s) - 1; i > -1; i -= 1 {
		reversed = append(reversed, s[i])
	}
	return string(reversed)

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
