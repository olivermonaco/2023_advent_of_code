package gear_ratios

import (
	"context"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/stretchr/testify/assert"
)

func TestGearRatios_PartOne(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		// {
		// 	inputFilename: "test_files/part_one/om_ex2.txt",
		// 	expOutput:     35,
		// },
		{
			inputFilename: "test_files/part_one/example_input.txt",
			expOutput:     4361,
		},
		// {
		// 	inputFilename: "puzzle_input.txt",
		// 	expOutput:     4361,
		// },
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartOne(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
