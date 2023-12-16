package maps

import (
	"context"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

func CalculatePartOne(ctx context.Context, input []string) int {
	const (
		startData = "AAA"
		endData   = "ZZZ"
	)

	instructions := getInstructions(input[0])

	initNodeMap := make(map[string]node)

	for _, line := range input[1:] {
		if line == "" {
			continue
		}
		data, leftData, rightData := parseNodeData(line)
		createNodeAddToMap(data, leftData, rightData, initNodeMap)
	}
	startNode, ok := initNodeMap[startData]
	if !ok {
		panic(initNodeMap)
	}
	currentNode := startNode
	var stepCounter int
	for {
		currentNode = currentNode.next(instructions[stepCounter].goLeft)
		stepCounter++
		if currentNode.data == endData {
			break
		}
	}
	return stepCounter
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

type node struct {
	data        string
	right, left *node
}

type instruction struct {
	original string
	goLeft   bool
}

func (n node) next(goLeft bool) node {
	if !goLeft {
		return *n.right
	}
	return *n.left
}

func getInstructions(line string) []instruction {
	var leftRightMap = map[rune]bool{
		'R': false,
		'L': true,
	}

	var instructions []instruction
	for _, r := range line {
		b, ok := leftRightMap[r]
		if !ok {
			panic(r)
		}
		instructions = append(
			instructions,
			instruction{
				original: string(r),
				goLeft:   b,
			},
		)
	}
	return instructions
}

func parseNodeData(line string) (string, string, string) {
	nodeData, nodes, found := strings.Cut(line, "=")
	if !found {
		panic(line)
	}
	rightNodeData, leftNodeData, found := strings.Cut(nodes, ",")
	if !found {
		panic(line)
	}
	rightNodeData = strings.TrimSpace(strings.Trim(rightNodeData, "("))
	leftNodeData = strings.TrimSpace(strings.Trim(leftNodeData, ")"))

	return strings.TrimSpace(nodeData), leftNodeData, rightNodeData
}

func createNodeAddToMap(
	data, leftData, rightData string,
	nodeDataToNode map[string]node,
) {

	n := node{
		data: data,
	}
	existing, ok := nodeDataToNode[data]
	if ok {
		n = existing
	}

	rightNode, ok := nodeDataToNode[rightData]
	if !ok {
		rightNode = node{
			data: rightData,
		}
		nodeDataToNode[rightData] = rightNode
	}
	n.right = kit.Ptr(rightNode)

	leftNode, ok := nodeDataToNode[leftData]
	if !ok {
		leftNode = node{
			data: leftData,
		}
		nodeDataToNode[rightData] = leftNode
	}
	n.left = kit.Ptr(leftNode)
}
