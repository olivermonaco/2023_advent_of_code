// https://adventofcode.com/2023/day/2
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olivermonaco/2023_advent_of_code/challenges/021223_cubes/cubes"
	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}

func main() {
	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cubes/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := cubes.CalculatePartOne(ctx, data, cubes.CompareTurn)
	result := cubes.CalculatePartTwo(ctx, data)
	fmt.Printf("result is %d", result)
}
