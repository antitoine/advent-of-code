package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

type GoRight = bool

func parseDirections(directionStr string) []GoRight {
	directions := make([]GoRight, len(directionStr))
	for i, direction := range directionStr {
		directions[i] = direction == 'R'
	}
	return directions
}

var nodeRegex = regexp.MustCompile(`^([1-9A-Z]+) = \(([1-9A-Z]+), ([1-9A-Z]+)\)$`)

func parseNodes(nodesStr string) (string, string, string) {
	matches := nodeRegex.FindStringSubmatch(nodesStr)
	if len(matches) != 4 {
		log.Fatalf("Invalid node string: %s", nodesStr)
	}
	return matches[1], matches[2], matches[3]
}

type Node struct {
	prefixName string
	endingName string
	left       *Node
	right      *Node
}

func parseInput(input io.Reader) ([]GoRight, map[string]*Node, []string) {
	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		log.Fatalf("Unable to read first line")
	}
	directions := parseDirections(scanner.Text())
	if !scanner.Scan() {
		log.Fatalf("Unable to read second line")
	}
	nodes := make(map[string]*Node)
	var startingNodes []string
	for scanner.Scan() {
		name, leftName, rightName := parseNodes(scanner.Text())
		prefixName := name[0:2]
		endingName := name[2:]
		currentNode, currentNodeFound := nodes[name]
		if !currentNodeFound {
			currentNode = &Node{prefixName: prefixName, endingName: endingName}
			nodes[name] = currentNode
		}
		if endingName == "A" {
			startingNodes = append(startingNodes, name)
		}
		leftNode, leftNodeFound := nodes[leftName]
		if !leftNodeFound {
			prefixLeftNodeName := leftName[0:2]
			endingLeftNodeName := leftName[2:]
			leftNode = &Node{prefixName: prefixLeftNodeName, endingName: endingLeftNodeName}
			nodes[leftName] = leftNode
		}
		rightNode, rightNodeFound := nodes[rightName]
		if !rightNodeFound {
			prefixRightNodeName := rightName[0:2]
			endingRightNodeName := rightName[2:]
			rightNode = &Node{prefixName: prefixRightNodeName, endingName: endingRightNodeName}
			nodes[rightName] = rightNode
		}
		currentNode.left = nodes[leftName]
		currentNode.right = nodes[rightName]
	}
	return directions, nodes, startingNodes
}

type Track struct {
	currentNode   *Node
	requireNbStep *int64
}

func greatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func leastCommonMultiple(a, b int64) int64 {
	return a * b / greatestCommonDivisor(a, b)
}

func getResult(input io.Reader) int64 {
	directions, nodes, startingNodes := parseInput(input)
	tracksNotCompleted := make([]*Node, len(startingNodes))
	for i, startingNode := range startingNodes {
		if node, ok := nodes[startingNode]; !ok {
			log.Fatalf("Unable to find starting node %s", startingNode)
		} else {
			tracksNotCompleted[i] = node
		}
	}
	var minRequiredSteps []int64
	var step int64
	for step = 0; len(tracksNotCompleted) > 0; step++ {
		goRight := directions[step%int64(len(directions))]
		for i := 0; i < len(tracksNotCompleted); {
			if goRight {
				tracksNotCompleted[i] = tracksNotCompleted[i].right
			} else {
				tracksNotCompleted[i] = tracksNotCompleted[i].left
			}
			if tracksNotCompleted[i].endingName == "Z" {
				minRequiredSteps = append(minRequiredSteps, step+1)
				tracksNotCompleted = append(tracksNotCompleted[:i], tracksNotCompleted[i+1:]...)
				continue
			}
			i++
		}
	}
	minRequiredStep := int64(1)
	for _, steps := range minRequiredSteps {
		minRequiredStep = leastCommonMultiple(minRequiredStep, steps)
	}
	return minRequiredStep
}

func main() {
	start := time.Now()
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
