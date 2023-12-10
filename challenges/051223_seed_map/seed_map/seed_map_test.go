package seed_map

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

func TestSeedMap_PartOne(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     35,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartOne(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestSeedMap_PartTwoBruteForceConcurrent(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     46,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartTwo_BruteForceConcurrent(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestSeedMap_PartTwo_CheckBounds(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     46,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartTwo_CheckBounds(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}

func TestSeedMap_PartTwo_ReverseLocation(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	tests := []struct {
		inputFilename string
		expOutput     int
	}{
		{
			inputFilename: "test_files/example_input.txt",
			expOutput:     46,
		},
		{
			inputFilename: "test_files/om_ex1.txt",
			expOutput:     0,
		},
	}

	for _, tt := range tests {
		data := kit.ReadFileConstructLines(ctx, tt.inputFilename)
		actual := CalculatePartTwo_ReverseLocationConcurrent(ctx, data)
		assert.Equalf(t, tt.expOutput, actual, "inequal expected:%d\nand actual:\n%d")
	}
}
