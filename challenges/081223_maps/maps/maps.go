package maps

import (
	"cmp"
	"context"
	"slices"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
	"github.com/rs/zerolog/log"
)

type node struct {
	data        string
	right, left *node
}

type nodeIterWrapper struct {
	currentNode    *node
	stepsTaken     int
	instructionIdx int
}

func (nIW *nodeIterWrapper) findNextEnd(
	instructions []instruction,
	endDataLogic func(node) bool,
) {
	currentNode := *nIW.currentNode
	for {
		nIW.instructionIdx = nIW.stepsTaken % len(instructions)
		currentInstruction := instructions[nIW.instructionIdx]
		currentNode = currentNode.next(currentInstruction.goLeft)
		nIW.stepsTaken++
		if endDataLogic(currentNode) {
			nIW.currentNode = &currentNode
			break
		}
	}
}

func (nIW *nodeIterWrapper) findNextEqualOrAbove(
	instructions []instruction,
	highestCurrentSteps int,
	endDataLogic func(node) bool,
) {
	if nIW.stepsTaken == highestCurrentSteps {
		return
	}
	for {
		nIW.findNextEnd(instructions, endDataLogic)
		if nIW.stepsTaken >= highestCurrentSteps {
			break
		}
	}
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

func CalculatePartOne(ctx context.Context, input []string) int {
	const (
		startData = "AAA"
		endData   = "ZZZ"
	)

	instructions := getInstructions(input[0])

	initNodeMap := make(map[string]*node)

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
	numStepsToEnd, _ := findNumStepsToEnd(
		startNode,
		instructions,
		0,
		func(n node) bool { return n.data == endData },
	)
	return numStepsToEnd
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	const (
		startData = 'A'
		endData   = 'Z'
	)

	l := log.Ctx(ctx).With().Logger()

	instructions := getInstructions(input[0])

	initNodeMap := make(map[string]*node)

	for _, line := range input[1:] {
		if line == "" {
			continue
		}
		data, leftData, rightData := parseNodeData(line)
		createNodeAddToMap(data, leftData, rightData, initNodeMap)
	}
	startNodes := findStartNodes(startData, initNodeMap)

	nodeIterWrappers := make([]nodeIterWrapper, 0, len(startNodes))
	for _, n := range startNodes {
		nodeIterWrappers = append(
			nodeIterWrappers,
			nodeIterWrapper{
				currentNode: n,
			},
		)
	}

	endDataLogic := func(n node) bool { return []rune(n.data)[2] == endData }
	for i, nIW := range nodeIterWrappers {
		nIW.findNextEnd(instructions, endDataLogic)
		nodeIterWrappers[i] = nIW
	}

	for {
		slices.SortFunc(nodeIterWrappers, func(a, b nodeIterWrapper) int {
			return cmp.Compare(a.stepsTaken, b.stepsTaken)
		})
		highestCurrent := nodeIterWrappers[len(nodeIterWrappers)-1]
		for i, nIW := range nodeIterWrappers[:len(nodeIterWrappers)-1] {
			nIW.findNextEqualOrAbove(
				instructions,
				highestCurrent.stepsTaken,
				endDataLogic,
			)
			nodeIterWrappers[i] = nIW
		}
		uniqueNodeIterWrappers := slices.CompactFunc(
			nodeIterWrappers,
			func(a, b nodeIterWrapper) bool {
				return a.stepsTaken == b.stepsTaken
			},
		)
		// logging
		uniqueSteps := make([]int, 0, len(uniqueNodeIterWrappers))
		for _, uNIW := range uniqueNodeIterWrappers {
			uniqueSteps = append(uniqueSteps, uNIW.stepsTaken)
		}
		slices.Sort(uniqueSteps)
		r := uniqueSteps[len(uniqueSteps)-1] - uniqueSteps[0]

		if len(uniqueNodeIterWrappers) == 1 {
			l.Info().
				Ints("unique steps taken", uniqueSteps).
				Int("range", r).
				Send()
			break
		}

		l.Info().
			Ints("unique steps taken", uniqueSteps).
			Int("range", r).
			Send()
	}
	return nodeIterWrappers[0].stepsTaken
}

func findStartNodes(startRune rune, initNodeMap map[string]*node) []*node {
	var startingNodes []*node
	for k, v := range initNodeMap {
		if []rune(k)[2] == startRune {
			startingNodes = append(startingNodes, v)
		}
	}
	return startingNodes
}

func findNumStepsToEnd(
	startNode *node,
	instructions []instruction,
	instructionIdx int,
	endDataLogic func(n node) bool,
) (int, *node) {
	currentNode := *startNode
	var stepCounter int
	for {
		idx := (stepCounter + instructionIdx) % len(instructions)
		currentInstruction := instructions[idx]
		currentNode = currentNode.next(currentInstruction.goLeft)
		stepCounter++
		if endDataLogic(currentNode) {
			break
		}
	}
	return stepCounter, &currentNode
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
	leftNodeData, rightNodeData, found := strings.Cut(nodes, ",")
	if !found {
		panic(line)
	}
	leftNodeData = strings.Trim(strings.TrimSpace(leftNodeData), "(")
	rightNodeData = strings.Trim(strings.TrimSpace(rightNodeData), ")")

	return strings.TrimSpace(nodeData), leftNodeData, rightNodeData
}

func createNodeAddToMap(
	data, leftData, rightData string,
	nodeDataToNode map[string]*node,
) {

	n := kit.Ptr(node{data: data})
	if existing, ok := nodeDataToNode[data]; ok {
		n = existing
	}

	leftNode := kit.Ptr(node{data: leftData})

	if leftData == data {
		leftNode = n
	}
	existingLeftNode, ok := nodeDataToNode[leftData]
	if !ok {
		nodeDataToNode[leftData] = leftNode
	} else {
		leftNode = existingLeftNode
	}
	n.left = leftNode

	rightNode := kit.Ptr(node{data: rightData})
	if rightData == data {
		rightNode = n
	}
	existingRightNode, ok := nodeDataToNode[rightData]
	if !ok {
		nodeDataToNode[rightData] = rightNode
	} else {
		rightNode = existingRightNode
	}
	n.right = rightNode

	nodeDataToNode[data] = n
}
