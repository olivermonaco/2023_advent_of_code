// https://adventofcode.com/2023/day/7
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/071223_cards/cards"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cards/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	result := cards.CalculatePartOne(ctx, data)
	fmt.Printf("result is %d", result)
}
