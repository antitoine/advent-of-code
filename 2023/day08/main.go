package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

type goRight = bool

func parseDirections(directionStr string) []goRight {
	directions := make([]goRight, len(directionStr))
	for i, direction := range directionStr {
		directions[i] = direction == 'R'
	}
	return directions
}

var nodeRegex = regexp.MustCompile(`^([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)$`)

func parseNodes(nodesStr string) (string, string, string) {
	matches := nodeRegex.FindStringSubmatch(nodesStr)
	if len(matches) != 4 {
		log.Fatalf("Invalid node string: %s", nodesStr)
	}
	return matches[1], matches[2], matches[3]
}

type Node struct {
	name  string
	left  *Node
	right *Node
}

func parseInput(input io.Reader) ([]goRight, map[string]*Node) {
	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		log.Fatalf("Unable to read first line")
	}
	directions := parseDirections(scanner.Text())
	if !scanner.Scan() {
		log.Fatalf("Unable to read second line")
	}
	nodes := make(map[string]*Node)
	for scanner.Scan() {
		name, leftName, rightName := parseNodes(scanner.Text())
		currentNode, currentNodeFound := nodes[name]
		if !currentNodeFound {
			currentNode = &Node{name: name}
			nodes[name] = currentNode
		}
		leftNode, leftNodeFound := nodes[leftName]
		if !leftNodeFound {
			leftNode = &Node{name: leftName}
			nodes[leftName] = leftNode
		}
		rightNode, rightNodeFound := nodes[rightName]
		if !rightNodeFound {
			rightNode = &Node{name: rightName}
			nodes[rightName] = rightNode
		}
		currentNode.left = nodes[leftName]
		currentNode.right = nodes[rightName]
	}
	return directions, nodes
}

func getResult(input io.Reader) int64 {
	directions, nodes := parseInput(input)
	currentNode, startingNodeFound := nodes["AAA"]
	if !startingNodeFound {
		log.Fatalf("Unable to find starting node")
	}
	var step int64
	for step = 0; currentNode.name != "ZZZ"; step++ {
		if directions[step%int64(len(directions))] {
			currentNode = currentNode.right
		} else {
			currentNode = currentNode.left
		}
	}
	return step
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
