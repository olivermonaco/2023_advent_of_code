package trebuchet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrebuchet(t *testing.T) {
	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     142,
		},
		{
			inputFilename: "test_files/om_ex1.txt",
			expOutput:     124,
		},
	}

	for _, tt := range tests {
		actual := Trebuchet(tt.inputFilename)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
