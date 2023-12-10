package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

const mapSize = 140

type Node struct {
	i              int
	j              int
	symbol         rune
	connectedNodes []*Node
}

/*
   | is a vertical pipe connecting north and south.
   - is a horizontal pipe connecting east and west.
   L is a 90-degree bend connecting north and east.
   J is a 90-degree bend connecting north and west.
   7 is a 90-degree bend connecting south and west.
   F is a 90-degree bend connecting south and east.
   . is ground; there is no pipe in this tile.
   S is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.
*/

func (n *Node) isConnectedToNode(otherNode *Node) bool {
	if n == nil || otherNode == nil {
		return false
	}
	for _, connectedNode := range n.connectedNodes {
		if otherNode == connectedNode {
			return true
		}
	}
	return false
}

func parseInput(input io.Reader) *Node {
	scanner := bufio.NewScanner(input)
	var draftGraph []*Node
	nodes := make([][]*Node, mapSize)
	var startingNode *Node
	for i := 0; scanner.Scan(); i++ {
		nodes[i] = make([]*Node, mapSize)
		for j, symbol := range scanner.Text() {
			node := &Node{
				i:              i,
				j:              j,
				symbol:         symbol,
				connectedNodes: make([]*Node, 2),
			}
			draftGraph = append(draftGraph, node)
			nodes[i][j] = node
		}
	}
	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	for _, node := range draftGraph {
		switch node.symbol {
		case '.':
			continue
		case '|':
			if node.i > 0 && node.i < mapSize-1 {
				node.connectedNodes[0] = nodes[node.i-1][node.j]
				node.connectedNodes[1] = nodes[node.i+1][node.j]
			} else {
				continue
			}
		case '-':
			if node.j > 0 && node.j < mapSize-1 {
				node.connectedNodes[0] = nodes[node.i][node.j-1]
				node.connectedNodes[1] = nodes[node.i][node.j+1]
			} else {
				continue
			}
		case 'L':
			if node.i > 0 && node.j < mapSize-1 {
				node.connectedNodes[0] = nodes[node.i-1][node.j]
				node.connectedNodes[1] = nodes[node.i][node.j+1]
			} else {
				continue
			}
		case 'J':
			if node.i > 0 && node.j > 0 {
				node.connectedNodes[0] = nodes[node.i-1][node.j]
				node.connectedNodes[1] = nodes[node.i][node.j-1]
			} else {
				continue
			}
		case '7':
			if node.i < mapSize-1 && node.j > 0 {
				node.connectedNodes[0] = nodes[node.i+1][node.j]
				node.connectedNodes[1] = nodes[node.i][node.j-1]
			} else {
				continue
			}
		case 'F':
			if node.i < mapSize-1 && node.j < mapSize-1 {
				node.connectedNodes[0] = nodes[node.i+1][node.j]
				node.connectedNodes[1] = nodes[node.i][node.j+1]
			} else {
				continue
			}
		case 'S':
			startingNode = node
		default:
			log.Fatalf("Unknown symbol: %v", string(node.symbol))
		}
	}

	if startingNode == nil {
		log.Fatalf("No starting node found")
	} else {
		connectedNodesIdx := 0
		if startingNode.i > 0 && nodes[startingNode.i-1][startingNode.j].isConnectedToNode(startingNode) {
			startingNode.connectedNodes[connectedNodesIdx] = nodes[startingNode.i-1][startingNode.j]
			connectedNodesIdx++
		}
		if startingNode.i < mapSize-1 && nodes[startingNode.i+1][startingNode.j].isConnectedToNode(startingNode) {
			startingNode.connectedNodes[connectedNodesIdx] = nodes[startingNode.i+1][startingNode.j]
			connectedNodesIdx++
		}
		if startingNode.j > 0 && nodes[startingNode.i][startingNode.j-1].isConnectedToNode(startingNode) {
			startingNode.connectedNodes[connectedNodesIdx] = nodes[startingNode.i][startingNode.j-1]
			connectedNodesIdx++
		}
		if startingNode.j < mapSize-1 && nodes[startingNode.i][startingNode.j+1].isConnectedToNode(startingNode) {
			startingNode.connectedNodes[connectedNodesIdx] = nodes[startingNode.i][startingNode.j+1]
			connectedNodesIdx++
		}
		if connectedNodesIdx != 2 {
			log.Fatalf("Starting node has %d connected nodes", connectedNodesIdx)
		}
	}

	return startingNode
}

func getResult(input io.Reader) int64 {
	startingNode := parseInput(input)
	prevNode := startingNode
	currentNode := startingNode.connectedNodes[0]
	var step int64
	for step = 1; currentNode != startingNode; step++ {
		if currentNode.connectedNodes[0] != prevNode {
			prevNode = currentNode
			currentNode = currentNode.connectedNodes[0]
		} else if currentNode.connectedNodes[1] != prevNode {
			prevNode = currentNode
			currentNode = currentNode.connectedNodes[1]
		} else {
			log.Fatalf("Unable to find next node from %v", currentNode)
		}
	}
	if step%2 != 0 {
		log.Fatalf("Step is not even: %d", step)
	}
	return step / 2
}

func loadFile() *os.File {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	return inputFile
}

func main() {
	start := time.Now()
	inputFile := loadFile()
	defer inputFile.Close()

	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
