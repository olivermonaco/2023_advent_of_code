// https://adventofcode.com/2023/day/5
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/051223_seed_map/seed_map"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}

func main() {
	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("seed_map/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := seed_map.CalculatePartOne(ctx, data)
	// result := seed_map.CalculatePartTwo_CheckBounds(ctx, data)
	// result := seed_map.CalculatePartTwo_BruteForceConcurrent(ctx, data)
	result := seed_map.CalculatePartTwo_ReverseLocationConcurrent(ctx, data)
	fmt.Printf("result is %d", result)
}
