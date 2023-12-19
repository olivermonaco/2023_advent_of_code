package pipe_maze

import (
	"context"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

const (
	startRune = 'S'
)

var validCoordToPipeTypes = map[int][]rune{
	0: []rune{'|', '7', 'F', 'S'},
	1: []rune{'-', '7', 'J', 'S'},
	2: []rune{'|', 'J', 'L', 'S'},
	3: []rune{'-', 'F', 'L', 'S'},
}

type canProgress interface {
	addToOrigin(o origin)
	next(lines []string) canProgress
}

type pipe struct {
	lastPipeRelation int
	pipeType         rune
}

type origin struct {
	pipe
	pipesSeen []pipe
}

func calcLastPipe(nextPipeRelation int) int {
	return (nextPipeRelation + len(validCoordToPipeTypes)/2) % len(validCoordToPipeTypes)
}

func (o origin) validNext(r rune, relation int) *pipe {
	var validPipe *pipe
	for _, pipeType := range validCoordToPipeTypes[relation] {
		if r == pipeType {
			validPipe = kit.Ptr(pipe{
				lastPipeRelation: calcLastPipe(relation),
				pipeType:         r,
			})
		}
	}
	return validPipe
}

func (o origin) next(potentialPipes map[int]rune) pipe {
	var nextPipe pipe
	for relation, pipeType := range potentialPipes {
		if o.validNext(pipeType, relation) != nil {
			nextPipe = *o.validNext(pipeType, relation)
		}
	}
	return nextPipe
}

func findOrigin(lines []string) origin {
	var o origin
	for lineIdx, line := range lines {
		for runeIdx, r := range []rune(line) {
			if r == startRune {
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
	return 0
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}
