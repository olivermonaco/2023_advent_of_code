package hot_springs

import (
	"context"
	"os"
	"testing"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}

func TestHotSprings_PartOne(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     21,
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

// TODO: either uncomment and fix or delete
// func TestCalcNonBrokenTotals(t *testing.T) {
// 	tests := []struct {
// 		in  refBuffs
// 		exp int
// 	}{
// 		{
// 			// ????????????? 4,1,2,1
// 			in: refBuffs{
// 				refBuff{
// 					separatedStringRef: separatedStringRef{
// 						start: 0,
// 						end:   3,
// 					},
// 				},
// 				refBuff{
// 					separatedStringRef: separatedStringRef{
// 						start: 5,
// 						end:   5,
// 					},
// 				},
// 				refBuff{
// 					separatedStringRef: separatedStringRef{
// 						start: 7,
// 						end:   8,
// 					},
// 				},
// 				refBuff{
// 					separatedStringRef: separatedStringRef{
// 						start: 10,
// 						end:   10,
// 					},
// 					rBuff: 3,
// 				},
// 			},
// 			exp: 15,
// 		},
// 	}

// 	for _, tt := range tests {
// 		actual := tt.in.calcNonBrokenTotals()
// 		// actualSum := kit.Map(actual, func(i int) int { return sumConsecNums(i, 0) })
// 		actualSum := kit.Sum(actual)
// 		assert.Equalf(t, tt.exp, actualSum,
// 			"inequal expected:%d\nand actual:\n%d", tt.exp, kit.Sum(actual),
// 		)
// 	}
// }

func TestCalcTotals(t *testing.T) {
	tests := []struct {
		in  refBuffGroups
		exp int
	}{
		{
			// ????????##?????? 1,1,4,1,1
			in: refBuffGroups{
				refBuffGroup{
					refBuffs: refBuffs{
						refBuff{
							separatedStringRef: separatedStringRef{
								start: 0,
								end:   0,
							},
						},
						refBuff{
							separatedStringRef: separatedStringRef{
								start: 2,
								end:   2,
							},
							rBuff: 2,
						},
					},
				},
				refBuffGroup{
					refBuffs: refBuffs{
						refBuff{
							separatedStringRef: separatedStringRef{
								start:       6,
								end:         9,
								brokenSpans: [][2]int{{8, 9}},
							},
						},
					},
					brokenSpans: [][2]int{{8, 9}},
				},
				refBuffGroup{
					refBuffs: refBuffs{
						refBuff{
							separatedStringRef: separatedStringRef{
								start: 11,
								end:   11,
							},
						},
						refBuff{
							separatedStringRef: separatedStringRef{
								start: 13,
								end:   13,
							},
							rBuff: 2,
						},
					},
				},
			},
			exp: 81,
		},
	}

	for _, tt := range tests {
		actual := tt.in.calcTotals()
		actualSum := kit.Sum(actual)
		assert.Equalf(t, tt.exp, actualSum,
			"inequal expected:%d\nand actual:\n%d", tt.exp, kit.Sum(actual),
		)
	}
}
