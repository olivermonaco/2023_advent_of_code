// https://adventofcode.com/2023/day/1
package main

import (
	"context"
	"fmt"

	"github.com/olivermonaco/2023_advent_of_code/challenges/011223_trebuchet/trebuchet"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("trebuchet/test_files/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)
	result := trebuchet.Calculate(ctx, data, trebuchet.IntFromStrPartTwo)
	fmt.Printf("result is %d", result)
}
