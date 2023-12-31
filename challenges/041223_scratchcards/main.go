// https://adventofcode.com/2023/day/4
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olivermonaco/2023_advent_of_code/challenges/041223_scratchcards/scratchcards"
	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("scratchcards/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := gear_ratios.CalculatePartOne(ctx, data)
	result := scratchcards.CalculatePartTwo(ctx, data)
	fmt.Printf("result is %d", result)
}
