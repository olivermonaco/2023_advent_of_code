// https://adventofcode.com/2023/day/2
package main

import (
	"context"
	"fmt"

	"github.com/olivermonaco/2023_advent_of_code/challenges/021223_cubes/cubes"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cubes/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := cubes.CalculatePartOne(ctx, data, cubes.CompareTurn)
	result := cubes.CalculatePartTwo(ctx, data)
	fmt.Printf("result is %d", result)
}
