package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Cell struct {
	symbol    rune
	visited   []Direction
	energized bool
}

func (c *Cell) IsVisited(direction Direction) bool {
	for _, visitedDirection := range c.visited {
		if visitedDirection == direction {
			return true
		}
	}
	return false
}

func (c *Cell) Visit(direction Direction) {
	c.visited = append(c.visited, direction)
	c.energized = true
}

func (c *Cell) IsEnergized() bool {
	return c.energized
}

func (c *Cell) NextDirection(direction Direction) []Direction {
	switch c.symbol {
	case '.':
		return []Direction{direction}
	case '|':
		if direction == up || direction == down {
			return []Direction{direction}
		} else {
			return []Direction{up, down}
		}
	case '-':
		if direction == left || direction == right {
			return []Direction{direction}
		} else {
			return []Direction{left, right}
		}
	case '/':
		switch direction {
		case up:
			return []Direction{right}
		case down:
			return []Direction{left}
		case left:
			return []Direction{down}
		case right:
			return []Direction{up}
		}
	case '\\':
		switch direction {
		case up:
			return []Direction{left}
		case down:
			return []Direction{right}
		case left:
			return []Direction{up}
		case right:
			return []Direction{down}
		}
	}
	log.Fatalf("Unknown symbol: %v", c.symbol)
	return nil
}

type Graph [][]*Cell

func (g Graph) GetCell(position Position) *Cell {
	return g[position.row][position.column]
}

func (g Graph) SymbolGraph() string {
	var result string
	for _, row := range g {
		for _, cell := range row {
			result += string(cell.symbol)
		}
		result += "\n"
	}
	return result
}

func (g Graph) EnergizedGraph() string {
	var result string
	for _, row := range g {
		for _, cell := range row {
			if cell.energized {
				result += "#"
			} else {
				result += "."
			}
		}
		result += "\n"
	}
	return result
}

type Position struct {
	row    int
	column int
}

func (p Position) Next(direction Direction) Position {
	switch direction {
	case up:
		return Position{p.row - 1, p.column}
	case down:
		return Position{p.row + 1, p.column}
	case left:
		return Position{p.row, p.column - 1}
	case right:
		return Position{p.row, p.column + 1}
	}
	panic(fmt.Sprintf("Unknown direction: %v", direction))
}

type Direction int

const (
	up    Direction = iota
	down  Direction = iota
	left  Direction = iota
	right Direction = iota
)

func parseInput(input io.Reader) Graph {
	scanner := bufio.NewScanner(input)

	var graph Graph
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]*Cell, len(line))
		for column, char := range line {
			row[column] = &Cell{symbol: char}
		}
		graph = append(graph, row)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return graph
}

func getCountOfEnergizedCells(graph Graph, position Position, direction Direction) int64 {
	if position.row < 0 || position.row >= len(graph) {
		return 0
	}
	if position.column < 0 || position.column >= len(graph[position.row]) {
		return 0
	}
	cell := graph.GetCell(position)
	if cell.IsVisited(direction) {
		return 0
	}
	var result int64
	if !cell.IsEnergized() {
		result++
	}
	cell.Visit(direction)
	for _, nextDirection := range cell.NextDirection(direction) {
		result += getCountOfEnergizedCells(graph, position.Next(nextDirection), nextDirection)
	}
	return result
}

func getResult(input io.Reader) int64 {
	graph := parseInput(input)
	//log.Printf("Initial graph:\n%s", graph.SymbolGraph())
	energizedCells := getCountOfEnergizedCells(graph, Position{0, 0}, right)
	//log.Printf("Energized graph:\n%s", graph.EnergizedGraph())
	return energizedCells
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
