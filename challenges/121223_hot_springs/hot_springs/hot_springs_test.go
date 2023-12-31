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
			inputFilename: "test_files/example_input2.txt",
			expOutput:     21,
		},
		// {
		// 	inputFilename: "test_files/example_input2.txt",
		// 	expOutput:     21,
		// },
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

// TODO: this passes, left off here, need to test the actual flow now
func TestCalcSpringLocCombos(t *testing.T) {
	tests := []struct {
		inStr    string
		inRow    row
		inKeys   []int
		expected int
	}{
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("??????"),
						},
					},
				},
				consecutiveKeys: []int{1, 1, 1},
			},
			expected: 4,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("???????????"),
						},
					},
				},
				consecutiveKeys: []int{1, 1, 2, 1},
			},
			expected: 35,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("??????"),
						},
					},
				},
				consecutiveKeys: []int{2, 1},
			},
			expected: 6,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("??#???#"),
						},
					},
				},
				consecutiveKeys: []int{4, 1},
			},
			expected: 2,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("????#?????#"),
						},
					},
				},
				consecutiveKeys: []int{4, 1},
			},
			expected: 4,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("??##??#???"),
						},
					},
				},
				consecutiveKeys: []int{4, 1},
			},
			expected: 2,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("?##????"),
						},
					},
				},
				consecutiveKeys: []int{3, 1},
			},
			expected: 5,
		},
		{
			inRow: row{
				sGs: []stringGroup{
					{
						{
							s: []rune("?#?#?#?#?#?#?#?"),
						},
					},
				},
				consecutiveKeys: []int{1, 3, 1, 6},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {

		actual := tt.inRow.calcSpringLocCombos()
		assert.Equalf(t, tt.expected, actual, "actual (%v) and expected (%v) inequal")
	}
}
