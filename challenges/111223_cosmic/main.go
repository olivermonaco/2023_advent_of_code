    // https://adventofcode.com/2023/day/11
package main

import (
	"context"
	"fmt"
	"os"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/olivermonaco/2023_advent_of_code/challenges/111223_cosmic/cosmic"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cosmic/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	result := cosmic.CalculatePartOne(ctx, data)
	fmt.Printf("result is %d", result)
}
