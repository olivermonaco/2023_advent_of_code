package kit

import (
	"bufio"
	"context"
	"os"

	"github.com/rs/zerolog/log"
)

func ReadFileConstructLines(ctx context.Context, filename string) []string {
	l := log.Ctx(ctx).With().Caller().Logger()

	l.Info().Msgf("filename is %s", filename)
	file, err := os.Open(filename)
	if err != nil {

		l.Err(err).Msg("Error opening file")
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
		l.Err(err).Msg("Error reading file")
		return nil
	}
	return lines
}
