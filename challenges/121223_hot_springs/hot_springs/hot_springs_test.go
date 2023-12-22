package hot_springs

import (
	"context"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		expected [][2]int
	}{
		// {
		// 	inStr:    "???",
		// 	inKeys:   []int{1, 1},
		// 	expected: [][2]int{{0, 0}, {2, 2}},
		// },
		// {
		// 	inStr:    "?###????????",
		// 	inKeys:   []int{4, 1, 1},
		// 	expected: [][2]int{{0, 3}, {5, 5}, {7, 7}},
		// },
		{
			inStr:    "?##?#????#",
			inKeys:   []int{3, 1, 1},
			expected: [][2]int{{0, 3}, {5, 5}, {7, 7}},
		},
	}

	for _, tt := range tests {
		actuals := calcSpringLocCombos(tt.inStr, tt.inKeys)
		require.Lenf(
			t, actuals, len(tt.expected),
			"actual keys:\n%v\nand expected keys:\n%v\ninequal",
			actuals, tt.expected,
		)
		for i, actual := range actuals {
			assert.Equalf(t, tt.expected[i], actual, "actual (%v) and expected (%v) inequal")
		}
	}
}
