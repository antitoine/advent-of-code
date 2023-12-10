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
	isLoopNode     bool
}

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

func parseInput(input io.Reader) ([][]*Node, *Node) {
	scanner := bufio.NewScanner(input)
	var nodes []*Node
	var graph [][]*Node
	var startingNode *Node
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		lineNodes := make([]*Node, len(line))
		for j, symbol := range line {
			node := &Node{
				i:      i,
				j:      j,
				symbol: symbol,
			}
			lineNodes[j] = node
			nodes = append(nodes, node)
		}
		graph = append(graph, lineNodes)
	}
	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	for _, node := range nodes {
		switch node.symbol {
		case '.':
			continue
		case '|':
			if node.i > 0 && node.i < len(graph)-1 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i-1][node.j], graph[node.i+1][node.j])
			} else {
				node.symbol = '.'
				continue
			}
		case '-':
			if node.j > 0 && node.j < len(graph[node.i])-1 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i][node.j-1], graph[node.i][node.j+1])
			} else {
				node.symbol = '.'
				continue
			}
		case 'L':
			if node.i > 0 && node.j < len(graph[node.i])-1 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i-1][node.j], graph[node.i][node.j+1])
			} else {
				node.symbol = '.'
				continue
			}
		case 'J':
			if node.i > 0 && node.j > 0 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i-1][node.j], graph[node.i][node.j-1])
			} else {
				node.symbol = '.'
				continue
			}
		case '7':
			if node.i < len(graph)-1 && node.j > 0 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i+1][node.j], graph[node.i][node.j-1])
			} else {
				node.symbol = '.'
				continue
			}
		case 'F':
			if node.i < len(graph)-1 && node.j < len(graph[node.i])-1 {
				node.connectedNodes = append(node.connectedNodes, graph[node.i+1][node.j], graph[node.i][node.j+1])
			} else {
				node.symbol = '.'
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
		fromNorth := false
		fromSouth := false
		fromEast := false
		fromWest := false
		if startingNode.i > 0 && graph[startingNode.i-1][startingNode.j].isConnectedToNode(startingNode) {
			startingNode.connectedNodes = append(startingNode.connectedNodes, graph[startingNode.i-1][startingNode.j])
			fromNorth = true
		}
		if startingNode.i < len(graph)-1 && graph[startingNode.i+1][startingNode.j].isConnectedToNode(startingNode) {
			startingNode.connectedNodes = append(startingNode.connectedNodes, graph[startingNode.i+1][startingNode.j])
			fromSouth = true
		}
		if startingNode.j > 0 && graph[startingNode.i][startingNode.j-1].isConnectedToNode(startingNode) {
			startingNode.connectedNodes = append(startingNode.connectedNodes, graph[startingNode.i][startingNode.j-1])
			fromWest = true
		}
		if startingNode.j < len(graph[startingNode.i])-1 && graph[startingNode.i][startingNode.j+1].isConnectedToNode(startingNode) {
			startingNode.connectedNodes = append(startingNode.connectedNodes, graph[startingNode.i][startingNode.j+1])
			fromEast = true
		}
		if len(startingNode.connectedNodes) != 2 {
			log.Fatalf("Starting node has %d connected nodes", len(startingNode.connectedNodes))
		}

		if fromNorth && fromSouth && !fromEast && !fromWest {
			startingNode.symbol = '|'
		} else if !fromNorth && !fromSouth && fromEast && fromWest {
			startingNode.symbol = '-'
		} else if fromNorth && !fromSouth && fromEast && !fromWest {
			startingNode.symbol = 'L'
		} else if fromNorth && !fromSouth && !fromEast && fromWest {
			startingNode.symbol = 'J'
		} else if !fromNorth && fromSouth && !fromEast && fromWest {
			startingNode.symbol = '7'
		} else if !fromNorth && fromSouth && fromEast && !fromWest {
			startingNode.symbol = 'F'
		} else {
			log.Fatalf("Unable to determine starting node symbol: %#v", startingNode)
		}
	}

	return graph, startingNode
}

func printGraph(graph [][]*Node) {
	for _, line := range graph {
		lineStr := ""
		for _, node := range line {
			lineStr += string(node.symbol)
		}
		log.Printf("%s", lineStr)
	}
}

func navigate(startingNode *Node) int64 {
	prevNode := startingNode
	currentNode := startingNode.connectedNodes[0]
	currentNode.isLoopNode = true
	var step int64
	for step = 1; currentNode != startingNode; step++ {
		var c int
		for c = 0; c < len(currentNode.connectedNodes) && currentNode.connectedNodes[c] == prevNode; c++ {
		}
		if c == len(currentNode.connectedNodes) {
			log.Fatalf("Unable to find next node from %v", currentNode)
		}
		prevNode = currentNode
		currentNode = currentNode.connectedNodes[c]
		currentNode.isLoopNode = true
	}
	return step
}

func getResultPart1(input io.Reader) int64 {
	graph, startingNode := parseInput(input)
	printGraph(graph)
	log.Printf("Starting node at i=%d j=%d", startingNode.i, startingNode.j)
	step := navigate(startingNode)
	if step%2 != 0 {
		log.Fatalf("Step is not even: %d", step)
	}
	return step / 2
}

func getResultPart2(input io.Reader) int64 {
	graph, startingNode := parseInput(input)
	printGraph(graph)
	log.Printf("Starting node at i=%d j=%d", startingNode.i, startingNode.j)
	navigate(startingNode)
	var nodesInsideLoop int64
	for i, line := range graph {
		if i == 0 || i == len(graph)-1 {
			continue
		}
		for j, node := range line {
			if j == 0 || j == len(line)-1 || node.isLoopNode {
				continue
			}
			var loopIntersection int64
			for expJ := j - 1; expJ >= 0; expJ-- {
				if graph[i][expJ].isLoopNode {
					if graph[i][expJ].symbol == '|' {
						loopIntersection++
					} else if graph[i][expJ].symbol == 'J' {
						expJ--
						for expJ >= 0 && graph[i][expJ].symbol == '-' {
							expJ--
						}
						if expJ < 0 {
							log.Fatalf("Unable to find node before J")
						}
						if graph[i][expJ].symbol == 'F' {
							loopIntersection++
						}
					} else if graph[i][expJ].symbol == '7' {
						expJ--
						for expJ >= 0 && graph[i][expJ].symbol == '-' {
							expJ--
						}
						if expJ < 0 {
							log.Fatalf("Unable to find node before J")
						}
						if graph[i][expJ].symbol == 'L' {
							loopIntersection++
						}
					}
				}
			}
			if loopIntersection%2 == 1 {
				nodesInsideLoop++
			}
		}
	}
	return nodesInsideLoop
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

	result := getResultPart2(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
