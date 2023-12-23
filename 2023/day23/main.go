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
		previousPositions: make([]Position, 0, len(s.previousPositions)+1),
		visited:           make(map[Position]bool),
		position:          newPosition,
	}
	for _, previousPosition := range s.previousPositions {
		newStep.previousPositions = append(newStep.previousPositions, previousPosition)
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

type StepQueue []*Step

func (sq *StepQueue) Push(step *Step) {
	*sq = append(*sq, step)
}

func (sq *StepQueue) Pop() *Step {
	old := *sq
	n := len(old)
	step := old[n-1]
	old[n-1] = nil
	*sq = old[0 : n-1]
	return step
}

func (g Grid) getHighestPath(start, end Position) *Step {
	allowedMoves := []Direction{North, South, East, West}

	stepQueue := StepQueue{StartStep(start)}
	var highestPath *Step

	for len(stepQueue) > 0 {
		step := stepQueue.Pop()

		if step.position == end {
			if highestPath == nil || step.PathLength() > highestPath.PathLength() {
				highestPath = step
			}
			continue
		}

		for _, move := range allowedMoves {
			if newStep, allowed := step.Move(g, move); allowed {
				stepQueue.Push(newStep)
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

func getResult(input io.Reader) int {
	grid := parseInput(input)
	//log.Printf("Grid:\n%v", grid.String())
	start := Position{0, 1}
	end := Position{len(grid) - 1, len(grid[len(grid)-1]) - 2}
	step := grid.getHighestPath(start, end)
	//log.Printf("Grid with path:\n%v", grid.String(step.GetPath()...))
	return step.PathLength()
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
