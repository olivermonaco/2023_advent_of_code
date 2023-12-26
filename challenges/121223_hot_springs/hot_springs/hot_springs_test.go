package hot_springs

import (
	"context"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/stretchr/testify/assert"
)

func TestHotSprings_PartOne(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     0,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartOne(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestHotSprings_PartTwo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     0,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartTwo(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestCalcSpringLocCombos(t *testing.T) {
	tests := []struct {
		inStr    string
		inKeys   []int
		expected int
	}{
		{
			inStr:    "??????",
			inKeys:   []int{1, 1, 1},
			expected: 4,
		},
		{
			inStr:    "???????????",
			inKeys:   []int{1, 1, 2, 1},
			expected: 35,
		},
		{
			inStr:    "??????",
			inKeys:   []int{2, 1},
			expected: 6,
		},
		{
			inStr:    "??#???#",
			inKeys:   []int{4, 1},
			expected: 2,
		},
		{
			inStr:    "????#????#",
			inKeys:   []int{4, 1},
			expected: 4,
		},
		{
			inStr:    "??##??#???",
			inKeys:   []int{4, 1},
			expected: 2,
		},
		{
			inStr:    "?##????",
			inKeys:   []int{3, 1},
			expected: 5,
		},
	}

	for _, tt := range tests {
		actual := calcSpringLocCombos(tt.inStr, tt.inKeys)
		assert.Equalf(t, tt.expected, actual, "actual (%v) and expected (%v) inequal")
	}
}
