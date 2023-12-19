package maps

import (
	"context"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/stretchr/testify/assert"
)

func TestMaps_PartOne(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/part_one/example_input.txt",
			expOutput:     2,
		},
		{
			inputFilename: "test_files/part_one/example_input2.txt",
			expOutput:     6,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartOne(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestMaps_PartTwo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/part_two/example_input.txt",
			expOutput:     6,
		},
		{
			inputFilename: "test_files/part_two/om_ex1.txt",
			expOutput:     5,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartTwo(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
