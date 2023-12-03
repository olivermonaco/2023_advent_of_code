package cubes

import "github.com/rs/zerolog/log"

var (
	Red   = Color("red")
	Green = Color("green")
	Blue  = Color("blue")
)

type Color string

type Cube interface {
	ShowColor() Color
}

type RedCube struct{}

func (RC RedCube) ShowColor() Color {
	return Red
}

type BlueCube struct{}

func (BC BlueCube) ShowColor() Color {
	return Blue
}

type GreenCube struct{}

func (GC GreenCube) ShowColor() Color {
	return Green
}

type Turn struct {
	RedCubes   []RedCube
	BlueCubes  []BlueCube
	GreenCubes []GreenCube
}

type Game struct {
	Turns []Turn
}

func Calculate(input []string, partFunc func(string) int) int {
	var result int
	for _, s := range input {
		intResult := partFunc(s)
		result += intResult
		log.Info().
			Str("processed_str", s).
			Int("int_from_runes", intResult).
			Int("updated_result", result).Msg("")
	}
	return result
}
