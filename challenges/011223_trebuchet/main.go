// https://adventofcode.com/2023/day/1
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olivermonaco/2023_advent_of_code/challenges/011223_trebuchet/trebuchet"
	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()
	ctx = log.Logger.WithContext(ctx)
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("trebuchet/test_files/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)
	result := trebuchet.Calculate(data, trebuchet.IntFromStrPartTwo)
	fmt.Printf("result is %d", result)
}
