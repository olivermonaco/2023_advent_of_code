// https://adventofcode.com/2023/day/3
package main

import (
	"context"
	"fmt"

	"github.com/olivermonaco/2023_advent_of_code/challenges/031223_gear_ratios/gear_ratios"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func main() {
	ctx := context.Background()
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("gear_ratios/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := gear_ratios.CalculatePartOne(ctx, data)
	result := gear_ratios.CalculatePartTwo(ctx, data)
	fmt.Printf("result is %d", result)
}
