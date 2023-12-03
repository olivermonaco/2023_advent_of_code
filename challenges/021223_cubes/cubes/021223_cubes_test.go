package cubes

import (
	"context"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/stretchr/testify/assert"
)

func TestCubes_PartOne(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		compareTurn   Turn
		expOutput     int
	}{
		{
			inputFilename: "test_files/part_one/example_input.txt",
			compareTurn:   CompareTurn,
			expOutput:     8,
		},
		{
			inputFilename: "test_files/part_one/om_ex1.txt",
			compareTurn:   CompareTurn,
			expOutput:     79,
		},
		{
			inputFilename: "test_files/part_one/om_ex2.txt",
			compareTurn:   CompareTurn,
			expOutput:     92,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := Calculate(ctx, data, tt.compareTurn)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
