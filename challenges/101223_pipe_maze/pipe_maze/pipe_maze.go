package pipe_maze

import (
	"context"
	"math"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

const (
	startRune = "S"
)

var (
	validCoordToPipeTypes = map[int][]string{
		0: {"|", "7", "F", "S"},
		1: {"-", "7", "J", "S"},
		2: {"|", "J", "L", "S"},
		3: {"-", "F", "L", "S"},
	}
	pipeTypeToCoords = map[string]map[int]struct{}{
		"|": {0: {}, 2: {}},
		"J": {0: {}, 3: {}},
		"L": {0: {}, 1: {}},
		"-": {1: {}, 3: {}},
		"F": {1: {}, 2: {}},
		"7": {2: {}, 3: {}},
	}
)

type pipe struct {
	xCoord, yCoord   int
	lastPipeRelation int
	pipeType         string
}

type origin struct {
	pipe
	pipesSeen []pipe
}

func (o origin) next(lines []string, potentialPipes map[int]string) pipe {
	var nextPipe pipe

	for relation, pipeType := range potentialPipes {
		if validPipe := o.validNext(pipeType, relation); validPipe != nil {
			nextPipe = *validPipe
			nextPipe.xCoord = o.xCoord + (nextPipe.lastPipeRelation-2)%2
			nextPipe.yCoord = o.yCoord - (nextPipe.lastPipeRelation-1)%2
			break
		}
	}
	return nextPipe
}

func getSurroundingTiles(lines []string, xCoord, yCoord int) map[int]string {
	surroundingTiles := make(map[int]string, 4)

	if yCoord > 0 {
		// get surrounding tile from line above current
		surroundingTiles[0] = string([]rune(lines[yCoord-1])[xCoord])
	}
	if yCoord < len(lines)-1 {
		// get surrounding tile from line below current
		surroundingTiles[2] = string([]rune(lines[yCoord+1])[xCoord])
	}
	if xCoord < len(lines)-1 {
		surroundingTiles[1] = string([]rune(lines[yCoord])[xCoord+1])
	}
	if xCoord > 0 {
		surroundingTiles[3] = string([]rune(lines[yCoord])[xCoord-1])
	}
	return surroundingTiles
}

func calcPrevPipeLoc(nextPipeRelation int) int {
	return (nextPipeRelation + len(validCoordToPipeTypes)/2) % len(validCoordToPipeTypes)
}

func (p pipe) validNext(r string, relation int) *pipe {
	var validPipe *pipe
	for _, pipeType := range validCoordToPipeTypes[relation] {
		if string(r) == pipeType {
			validPipe = kit.Ptr(
				pipe{
					lastPipeRelation: calcPrevPipeLoc(relation),
					pipeType:         string(r),
				},
			)
		}
	}
	return validPipe
}

func (p pipe) next(lines []string, surroundingTiles map[int]string) pipe {
	nextPipeOptions, ok := pipeTypeToCoords[p.pipeType]
	if !ok {
		panic(p.pipeType)
	}

	var nextPipeLoc *int
	for pipeOption := range nextPipeOptions {
		if pipeOption != p.lastPipeRelation {
			nextPipeLoc = kit.Ptr(pipeOption)
			break
		}
	}

	nextPipeType, ok := surroundingTiles[*nextPipeLoc] // as per usual, accept a panic given it's a script
	if !ok {
		panic(nextPipeLoc)
	}
	var nextPipe pipe
	nextPipe.pipeType = nextPipeType
	nextPipe.lastPipeRelation = calcPrevPipeLoc(*nextPipeLoc)
	nextPipe.xCoord = p.xCoord + (nextPipe.lastPipeRelation-2)%2
	nextPipe.yCoord = p.yCoord - (nextPipe.lastPipeRelation-1)%2

	return nextPipe
}

func findOrigin(lines []string) origin {
	var o origin
	for lineIdx, line := range lines {
		for runeIdx, r := range []rune(line) {
			if string(r) == startRune {
				o = origin{
					pipe: pipe{
						pipeType: startRune,
						xCoord:   runeIdx,
						yCoord:   lineIdx,
					},
				}
			}
		}
	}
	return o
}

func CalculatePartOne(ctx context.Context, input []string) int {
	o := findOrigin(input)

	potentialPipes := getSurroundingTiles(input, o.xCoord, o.yCoord)
	current := o.next(input, potentialPipes)
	o.pipesSeen = append(o.pipesSeen, current)
	if current.pipeType == o.pipeType {
		return 0
	}

	for {
		surroundingTiles := getSurroundingTiles(input, current.xCoord, current.yCoord)
		current = current.next(input, surroundingTiles)
		if current.pipeType == o.pipeType {
			break
		}
		o.pipesSeen = append(o.pipesSeen, current)
	}

	farthest := float64(len(o.pipesSeen)) / float64(2)
	return int(math.Ceil(farthest))
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}
