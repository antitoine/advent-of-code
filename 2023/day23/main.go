package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

type Direction int

const (
	North Direction = iota
	South Direction = iota
	East  Direction = iota
	West  Direction = iota
)

type Position struct {
	i int
	j int
}

func (p Position) Move(direction Direction) Position {
	switch direction {
	case North:
		return Position{p.i - 1, p.j}
	case South:
		return Position{p.i + 1, p.j}
	case East:
		return Position{p.i, p.j + 1}
	case West:
		return Position{p.i, p.j - 1}
	}
	log.Fatalf("Unknown direction: %d", direction)
	return p
}

type Grid [][]rune

func (g Grid) isValidMove(position Position) bool {
	return position.i >= 0 &&
		position.i < len(g) &&
		position.j >= 0 &&
		position.j < len(g[position.i]) &&
		g[position.i][position.j] != '#'
}

func (g Grid) isSlopeWithDirection(position Position) (bool, Direction) {
	switch g[position.i][position.j] {
	case 'v':
		return true, South
	case '>':
		return true, East
	case '<':
		return true, West
	case '^':
		return true, North
	}
	return false, North
}

func (g Grid) String(positions ...Position) string {
	var result string
	for i, row := range g {
		for j, cell := range row {
			if slices.Contains(positions, Position{i, j}) {
				result += "O"
			} else {
				result += fmt.Sprintf("%c", cell)
			}
		}
		result += "\n"
	}
	return result
}

type Link struct {
	to   *Node
	cost int64
}

type Node struct {
	position  Position
	neighbors []Link
}

type Graph map[Position]*Node

func (g Grid) Graph() Graph {
	graph := make(Graph)
	for i, row := range g {
		for j, cell := range row {
			position := Position{i, j}
			if cell != '#' {
				graph[position] = &Node{
					position:  position,
					neighbors: make([]Link, 0),
				}
			}
		}
	}

	for _, node := range graph {
		for _, direction := range []Direction{North, South, East, West} {
			if newNeighbor := node.position.Move(direction); g.isValidMove(newNeighbor) {
				node.neighbors = append(node.neighbors, Link{
					to:   graph[newNeighbor],
					cost: 1,
				})
			}
		}
	}

	for {
		hasCompressed := false
		for _, node := range graph {
			if len(node.neighbors) == 2 {
				hasCompressed = true
				neighbor1 := node.neighbors[0]
				neighbor2 := node.neighbors[1]

				var neighbor1Links []Link
				for _, neighbor1Link := range graph[neighbor1.to.position].neighbors {
					if neighbor1Link.to != node {
						neighbor1Links = append(neighbor1Links, neighbor1Link)
					}
				}
				neighbor1Links = append(neighbor1Links, Link{
					to:   neighbor2.to,
					cost: neighbor1.cost + neighbor2.cost,
				})
				graph[neighbor1.to.position].neighbors = neighbor1Links

				var neighbor2Links []Link
				for _, neighbor2Link := range graph[neighbor2.to.position].neighbors {
					if neighbor2Link.to != node {
						neighbor2Links = append(neighbor2Links, neighbor2Link)
					}
				}
				neighbor2Links = append(neighbor2Links, Link{
					to:   neighbor1.to,
					cost: neighbor1.cost + neighbor2.cost,
				})
				graph[neighbor2.to.position].neighbors = neighbor2Links

				delete(graph, node.position)
				break
			}
		}
		if !hasCompressed {
			break
		}
	}

	return graph
}

func (g Graph) GetHighestPath(start, end Position, initialPathLength int64, visited []Position) int64 {
	if start == end {
		return initialPathLength
	}
	maxLength := int64(-1)
	for _, neighbor := range g[start].neighbors {
		if slices.Contains(visited, neighbor.to.position) {
			continue
		}
		pathLength := g.GetHighestPath(neighbor.to.position, end, initialPathLength+neighbor.cost, append(visited, neighbor.to.position))
		if pathLength > maxLength {
			maxLength = pathLength
		}
	}
	return maxLength
}

type Step struct {
	previousPositions []Position
	visited           map[Position]bool
	position          Position
}

func StartStep(position Position) *Step {
	newStep := &Step{
		previousPositions: make([]Position, 0),
		visited:           make(map[Position]bool),
		position:          position,
	}
	newStep.visited[position] = true
	return newStep
}

func (s *Step) Move(grid Grid, direction Direction) (*Step, bool) {
	newPosition := s.position.Move(direction)

	if !grid.isValidMove(newPosition) {
		return nil, false
	}
	if s.visited[newPosition] {
		return nil, false
	}

	newStep := &Step{
		previousPositions: make([]Position, len(s.previousPositions), len(s.previousPositions)+1),
		visited:           make(map[Position]bool),
		position:          newPosition,
	}
	for i, previousPosition := range s.previousPositions {
		newStep.previousPositions[i] = previousPosition
		newStep.visited[previousPosition] = true
	}
	newStep.previousPositions = append(newStep.previousPositions, s.position)
	newStep.visited[s.position] = true
	newStep.visited[newPosition] = true

	isSlope, slopeDirection := grid.isSlopeWithDirection(newPosition)
	if isSlope {
		return newStep.Move(grid, slopeDirection)
	}

	return newStep, true
}

func (s *Step) PathLength() int {
	return len(s.previousPositions)
}

func (s *Step) GetPath() []Position {
	return append(s.previousPositions, s.position)
}

func (g Grid) GetHighestPath(start, end Position) *Step {
	allowedMoves := []Direction{North, South, East, West}

	var stepQueue []*Step
	stepQueue = append(stepQueue, StartStep(start))

	var highestPath *Step

	for len(stepQueue) > 0 {
		step := stepQueue[0]
		stepQueue = stepQueue[1:]
		if step.position == end {
			if highestPath == nil || step.PathLength() > highestPath.PathLength() {
				highestPath = step
			}
			continue
		}

		for _, move := range allowedMoves {
			if newStep, allowed := step.Move(g, move); allowed {
				stepQueue = append(stepQueue, newStep)
			}
		}
	}

	return highestPath
}

func parseInput(input io.Reader) Grid {
	scanner := bufio.NewScanner(input)

	var grid Grid
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid
}

func getResultPart1(input io.Reader) int {
	grid := parseInput(input)
	start := Position{0, 1}
	end := Position{len(grid) - 1, len(grid[len(grid)-1]) - 2}
	return grid.GetHighestPath(start, end).PathLength()
}

func getResultPart2(input io.Reader) int64 {
	grid := parseInput(input)
	graph := grid.Graph()
	start := Position{0, 1}
	end := Position{len(grid) - 1, len(grid[len(grid)-1]) - 2}
	return graph.GetHighestPath(start, end, 0, []Position{start})
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
