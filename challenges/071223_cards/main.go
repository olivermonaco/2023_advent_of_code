// https://adventofcode.com/2023/day/7
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/071223_cards/cards"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func main() {
	t := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("log_%s.txt", t.Format("01-06-2006_3:04pm")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening log file")
	}
	defer file.Close()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	fileLogger := zerolog.New(file).With().Timestamp().Caller().Logger()

	consoleLogger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	multi := zerolog.MultiLevelWriter(consoleLogger, fileLogger)
	log.Logger = zerolog.New(multi)

	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	filename := "puzzle_input.txt"
	relFilepath := fmt.Sprintf("cards/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	result := cards.CalculatePartOne(ctx, data)
	fmt.Printf("result is %d", result)
}
