package cubes

import (
	"context"
	"fmt"
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
				Count: 14,
			},
			Color: Blue{},
		},
		Green{}: {
			Cubes: Cubes{
				Count: 13,
			},
			Color: Green{},
		},
	},
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

func (cc ColoredCubes[Color]) MarshalZerologObject(e *zerolog.Event) {
	if cc.Count > 0 {
		e.Str("cube_color", cc.Color.ShowColorStr()).
			Int("cubes_count", cc.Count)
	}
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
	Index int
	Cubes map[Color]ColoredCubes[Color]
}

type Turns []Turn

func (t *Turn) createColorMaps() {
	t.Cubes = make(map[Color]ColoredCubes[Color])
	for _, color := range []Color{Red{}, Blue{}, Green{}} {
		t.Cubes[color] = ColoredCubes[Color]{Color: color}
	}
}

func (t Turn) AddToTurn(cc ColoredCubes[Color]) {
	t.Cubes[cc.Color] = cc
}

func (t Turn) CalculateCountPowers() int {
	var result int
	for _, cc := range t.Cubes {
		if result == 0 {
			result = cc.Count
			continue
		}
		result *= cc.Count
	}
	return result
}

func (t Turn) MarshalZerologObject(e *zerolog.Event) {
	e.Int("turn_idx", t.Index)

	for _, cc := range t.Cubes {
		cc.MarshalZerologObject(e)
	}
}

func (tt Turns) MarshalZerologArray(a *zerolog.Array) {
	for _, t := range tt {
		a.Object(t)
	}
}

type Game struct {
	ID    int
	Turns []Turn
}

func (g Game) MarshalZerologObject(e *zerolog.Event) {
	e.Int("game_id", g.ID)
	turns := Turns(g.Turns)
	e.Array("turns", turns)
}

func (g Game) ImpossibleTurns(compareTurn Turn) Turns {
	var impossibleTurns Turns
	for _, turn := range g.Turns {
		for color, cubes := range turn.Cubes {
			if impossibleCubes := compareTurn.Cubes[color].CompareCube(cubes); impossibleCubes != nil {
				impossibleTurns = append(impossibleTurns, turn)
			}
		}
	}
	return impossibleTurns
}

func (g Game) MinCubesColor() Turn {
	var minCubesTurn Turn
	minCubesTurn.createColorMaps()

	for _, turn := range g.Turns {
		for color, cc := range turn.Cubes {
			if minCubesTurn.Cubes[color].Count < cc.Count {
				minCubesTurn.Cubes[color] = cc
			}
		}
	}
	return minCubesTurn
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
	for idx, turn := range turnsStr {
		var t Turn
		t.createColorMaps()
		t.Index = idx
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

func CalculatePartOne(ctx context.Context, input []string, compareTurn Turn) int {
	l := log.Ctx(ctx).With().Logger()

	var result int
	for _, s := range input {
		log.Info().Str("input_string", s).Msg("")
		game := toGame(s)
		if impossibleTurns := game.ImpossibleTurns(compareTurn); len(impossibleTurns) > 0 {
			l.Error().
				Str("processed_str", s).
				Int("game_id", game.ID).
				Array("impossible_turns", impossibleTurns).
				Int("result", result).
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
			Object("game", game).
			Int("updated_result", result).Msg("valid game")
		fmt.Println(result)
	}
	return result
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	ctx = log.Logger.WithContext(ctx)
	l := log.Ctx(ctx).With().Caller().Logger()

	var result int
	for _, s := range input {
		log.Info().Str("input_string", s).Msg("")
		game := toGame(s)
		minCubesForGame := game.MinCubesColor()
		powerResult := minCubesForGame.CalculateCountPowers()
		result += powerResult
		l.Info().
			Object("min_cubes_for_game", minCubesForGame).
			Int("power_result", powerResult).
			Msg("")
		fmt.Println(result)
	}
	return result
}
