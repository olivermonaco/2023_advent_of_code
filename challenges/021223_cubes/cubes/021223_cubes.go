package cubes

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var CompareTurn = Turn{
	Cubes: map[Color]ColoredCubes[Color]{
		Red{}: {
			Cubes: Cubes{
				Count: 12,
			},
			Color: Red{},
		},
		Blue{}: {
			Cubes: Cubes{
				Count: 13,
			},
			Color: Blue{},
		},
		Green{}: {
			Cubes: Cubes{
				Count: 14,
			},
			Color: Green{},
		},
	},
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type Color interface {
	ShowColorStr() string
}

type Red struct{}

func (r Red) ShowColorStr() string {
	return "red"
}

type Green struct{}

func (g Green) ShowColorStr() string {
	return "green"
}

type Blue struct{}

func (b Blue) ShowColorStr() string {
	return "blue"
}

type Cubes struct {
	Count int
}

type ColoredCubes[T Color] struct {
	Cubes
	Color T
}

func (cc ColoredCubes[T]) CompareCube(other ColoredCubes[T]) *ColoredCubes[T] {
	if cc.Count < other.Count {
		return &other
	}
	return nil
}

func coloredCubesFromColor[T Color](cubes Cubes, c T) ColoredCubes[T] {
	return ColoredCubes[T]{
		Cubes: cubes,
		Color: c,
	}
}

type Turn struct {
	Cubes map[Color]ColoredCubes[Color]
}

func (t *Turn) createColorMaps() {
	t.Cubes = make(map[Color]ColoredCubes[Color])
	for _, color := range []Color{Red{}, Blue{}, Green{}} {
		t.Cubes[color] = ColoredCubes[Color]{Color: color}
	}
}

func (t Turn) AddToTurn(cc ColoredCubes[Color]) {
	t.Cubes[cc.Color] = cc
}

type Game struct {
	ID    int
	Turns []Turn
}

func (G Game) ImpossibleTurns(compareTurn Turn) []Turn {
	var impossibleTurns []Turn
	for _, turn := range G.Turns {
		for color, cubes := range turn.Cubes {
			if impossibleCubes := compareTurn.Cubes[color].CompareCube(cubes); impossibleCubes != nil {
				impossibleTurns = append(impossibleTurns, turn)
			}
		}
	}
	return impossibleTurns
}

func extractGameIDTurns(s string) (int, string) {
	_, remaining, found := strings.Cut(s, "Game ")
	if !found {
		panic(s)
	}

	num, remaining, found := strings.Cut(remaining, ":")
	if !found {
		panic(s)
	}

	gameID, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}

	return gameID, remaining
}

func toColoredCube(s string) ColoredCubes[Color] {
	countColor := strings.Fields(s)
	if len(countColor) != 2 {
		panic(s)
	}
	count, err := strconv.Atoi(countColor[0])
	if err != nil {
		panic(err)
	}
	cubes := Cubes{Count: count}

	var c Color
	colorStr := countColor[1]
	switch colorStr {
	case "red":
		c = Red{}
	case "blue":
		c = Blue{}
	case "green":
		c = Green{}
	}
	cc := coloredCubesFromColor(cubes, c)
	return cc
}

func extractTurns(line string) []Turn {
	var turns []Turn

	turnsStr := strings.Split(line, ";")
	for _, turn := range turnsStr {
		var t Turn
		t.createColorMaps()
		cubes := strings.Split(turn, ",")
		for _, cube := range cubes {
			cc := toColoredCube(cube)
			t.AddToTurn(cc)
		}
		turns = append(turns, t)
	}
	return turns
}

func toGame(line string) Game {
	var g Game
	gameID, turnsStr := extractGameIDTurns(line)
	g.ID = gameID
	g.Turns = extractTurns(turnsStr)
	return g
}

func (compareCubes ColoredCubes[T]) LogInvalidCube(ctx context.Context, other ColoredCubes[T]) {
	l := log.Ctx(ctx).With().Caller().Logger()

	if impossibleCube := compareCubes.CompareCube(other); impossibleCube != nil {
		l.Error().
			Int(
				fmt.Sprintf("%s_cube_count", compareCubes.Color.ShowColorStr()),
				other.Cubes.Count,
			).Msg("impossible cube count")
	}
}

func Calculate(ctx context.Context, input []string, compareTurn Turn) int {
	ctx = log.Logger.WithContext(ctx)
	l := log.Ctx(ctx).With().Caller().Logger()

	var result int
	for _, s := range input {
		game := toGame(s)
		if impossibleTurns := game.ImpossibleTurns(compareTurn); len(impossibleTurns) > 0 {
			// just for logging purposes
			l.Error().
				Str("processed_str", s).
				Msgf("Game %d not possible", game.ID)
			for _, turn := range impossibleTurns {
				for color, cc := range turn.Cubes {
					compareCubes := compareTurn.Cubes[color]
					compareCubes.LogInvalidCube(ctx, cc)
				}
			}
			continue
		}
		result += game.ID
		log.Info().
			Int("game_id", game.ID).
			Int("updated_result", result).Msg("valid game")
	}
	return result
}
