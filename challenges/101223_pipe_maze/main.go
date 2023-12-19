// https://adventofcode.com/2023/day/10
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/101223_pipe_maze/pipe_maze"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("pipe_maze/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	result := pipe_maze.CalculatePartOne(ctx, data)
	fmt.Printf("result is %d", result)
}
