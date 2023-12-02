package main

import (
	"fmt"
	"os"

	"github.com/olivermonaco/2023_advent_of_code/challenges/01123_trebuchet/trebuchet"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("trebuchet/test_files/%s", filename)
	result := trebuchet.Trebuchet(relFilepath)
	fmt.Printf("result is %d", result)
}
