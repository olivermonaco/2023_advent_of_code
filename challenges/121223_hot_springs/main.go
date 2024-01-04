// https://adventofcode.com/2023/day/12
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/121223_hot_springs/hot_springs"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// TODO: Left off here, run debugger on this and debug panics / len issues
func main() {
	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	filename := "example_input.txt"
	relFilepath := fmt.Sprintf("hot_springs/test_files/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	result := hot_springs.CalculatePartOne(ctx, data)
	fmt.Printf("result is %d", result)
}
