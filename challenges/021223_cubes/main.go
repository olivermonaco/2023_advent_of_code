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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()
	ctx = log.Logger.WithContext(ctx)
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cube/test_files/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	compareTurn := cubes.Turn{
		Cubes: map[cubes.Color]cubes.ColoredCubes[cubes.Color]{
			cubes.Red{}: {
				Cubes: cubes.Cubes{
					Count: 12,
				},
			},
			cubes.Blue{}: {
				Cubes: cubes.Cubes{
					Count: 13,
				},
			},
			cubes.Green{}: {
				Cubes: cubes.Cubes{
					Count: 14,
				},
			},
		},
	}
	result := cubes.Calculate(ctx, data, compareTurn)
	fmt.Printf("result is %d", result)
}
