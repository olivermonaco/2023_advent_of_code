package cosmic

import (
	"context"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

const galaxyChar = '#'

type galaxy struct {
	id                 int
	xCoord, yCoord     int
	galaxyIDToDistance map[int]int
}

type space struct {
	emptyColumns, emptyRows map[int]struct{}
	galaxies                []galaxy
}

func CalculatePartOne(ctx context.Context, input []string) int {
	s := createSpace(input)
	var total int
	for i, g := range s.galaxies {
		if i == len(s.galaxies)-1 {
			break
		}
		g.compareOtherGalaxies(s.galaxies[i+1:], s.emptyColumns, s.emptyRows)
		total += g.countTotal()
	}

	return total
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func createSpace(lines []string) space {
	var s space
	s.emptyRows = make(map[int]struct{})
	s.emptyColumns = make(map[int]struct{})
	for rowIdx, line := range lines {
		// add to map of columns with
		s.emptyRows[rowIdx] = struct{}{}
		for columnIdx, r := range line {
			if r == galaxyChar {
				g := galaxy{
					xCoord: columnIdx,
					yCoord: rowIdx,
				}
				if len(s.galaxies) > 0 {
					g.id = s.galaxies[len(s.galaxies)-1].id + 1
				}
				s.galaxies = append(s.galaxies, g)
			}
		}
	}
	// create map of columns with no galaxies
	for colNum := range lines[0] {
		s.emptyColumns[colNum] = struct{}{}
	}

	s.deleteRowsCols()

	return s
}

func (s space) deleteRowsCols() {
	for _, g := range s.galaxies {
		delete(s.emptyColumns, g.xCoord)
		delete(s.emptyRows, g.yCoord)
	}
}

func (g *galaxy) compareOtherGalaxies(
	otherGalaxies []galaxy,
	emptyColumns, emptyRows map[int]struct{},
) {
	g.galaxyIDToDistance = make(map[int]int, len(otherGalaxies))
	for _, other := range otherGalaxies {
		diffX, diffY := g.diff(other)
		for colIdx := range emptyColumns {
			if kit.IsBetween(colIdx, g.xCoord, other.xCoord) {
				diffX++
			}
		}
		for rowIdx := range emptyRows {
			if kit.IsBetween(rowIdx, g.yCoord, other.yCoord) {
				diffY++
			}
		}
		g.galaxyIDToDistance[other.id] = diffX + diffY
	}
}

func (g *galaxy) countTotal() int {
	var total int
	for _, count := range g.galaxyIDToDistance {
		total += count
	}
	return total
}

func (g galaxy) diff(other galaxy) (int, int) {
	return kit.Abs(g.xCoord - other.xCoord), kit.Abs(g.yCoord - other.yCoord)
}

// func (s space) countSpace(x, y int) int {
// 	spaceCount := 1
// 	if _, ok := s.emptyRows[x]; ok {
// 		spaceCount++
// 	}
// 	if _, ok := s.emptyColumns[y]; ok {
// 		spaceCount++
// 	}
// 	return spaceCount
// }
