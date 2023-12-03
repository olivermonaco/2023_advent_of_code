package trebuchet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrebuchet_PartOne(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/part_one/example_input.txt",
			expOutput:     142,
		},
		{
			inputFilename: "test_files/part_one/om_ex1.txt",
			expOutput:     124,
		},
	}

	for _, tt := range tests {
		data := ReadFileConstructLines(tt.inputFilename)
		actual := Calculate(ctx, data, IntFromStrPartOne)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestTrebuchet_PartTwo(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/part_two/example_input.txt",
			expOutput:     281,
		},
		{
			inputFilename: "test_files/part_two/om_ex1.txt",
			expOutput:     231,
		},
	}

	for _, tt := range tests {
		data := ReadFileConstructLines(tt.inputFilename)
		actual := Calculate(ctx, data, IntFromStrPartTwo)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
