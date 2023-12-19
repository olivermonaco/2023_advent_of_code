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
    1:[]rune{'-', '7', 'J', 'S'},
    2:[]rune{'|', 'J', 'L', 'S'},
    3:[]rune{'-', 'F', 'L', 'S'},
}


type canProgress interface {
	addToOrigin(o origin)
	next(lines []string) canProgress
}

type pipe struct {
	lastPipeRelation int
	pipeType       rune
}

type origin struct {
	pipe
	pipesSeen []pipe
}

func calcLastPipe(nextPipeRelation int) int {
    return (nextPipeRelation + len(validCoordToPipeTypes) / 2) % len(validCoordToPipeTypes)
}

func (o origin) validNext(r rune, relation int) *pipe {
    var validPipe *pipe
        for coord, pipeType := range validCoordToPipeTypes[relation] {
            if r == pipeType {
                validPipe = kit.Ptr(pipe{
                    lastPipe: calcLastPipe(0),
                    pipeType: r,
                })
            }
        }
        return validPipe
}

func (o origin) next(lines []string) pipe {
	var prevLine, nextLine *string
	if o.yCoord > 0 {
		prevLine = kit.Ptr(lines[o.yCoord-1])
	}
	if o.yCoord < len(lines)-1 {
		prevLine = kit.Ptr(lines[o.yCoord+1])
	}
	var nP pipe
	curLine := lines[o.yCoord]
	if prevLine != nil {
        if []rune(curLine)[o.xCoord]
	}

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
