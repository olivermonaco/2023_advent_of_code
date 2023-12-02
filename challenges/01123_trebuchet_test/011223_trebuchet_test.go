package trebuchet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrebuchet(t *testing.T) {
	tests := []struct {
		input     []string
		expOutput int
	}{
		{
			input: []string{"1abc2",
				"pqr3stu8vwx",
				"a1b2c3d4e5f",
				"treb7uchet"},
			expOutput: 77,
		},
	}

	for _, tt := range tests {
		actual := Trebuchet(tt.input)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
