// https://adventofcode.com/2023/day/12
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/olivermonaco/2023_advent_of_code/challenges/121223_hot_springs/hot_springs"
	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func initFile() *os.File {
	t := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("log_%s.txt", t.Format("01-06-2006_3:04pm")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return file
}
func initFileWriter(file *os.File) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{Out: file, NoColor: true}
}

func initLogger() *os.File {

	file := initFile()
	fileWriter := initFileWriter(file)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	multi := zerolog.MultiLevelWriter(
		consoleWriter,
		fileWriter,
	)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

	return file
}

func main() {
	file := initLogger()
	defer file.Close()
	ctx := log.Logger.WithContext(context.Background())
	// Get the absolute path of the current file
	// filename := "example_input12.txt"
	// filename := "example_input.txt"
	filename := "puzzle_input.txt"
	// relFilepath := fmt.Sprintf("hot_springs/test_files/%s", filename)
	relFilepath := fmt.Sprintf("hot_springs/%s", filename)
	data := kit.ReadFileConstructLines(ctx, relFilepath)

	// result := hot_springs.CalculatePartOne(ctx, data)
	result := hot_springs.CalculatePartTwo(ctx, data)
	log.Logger.Info().Int("result", result).Send()
	fmt.Printf("result is %d", result)
}
