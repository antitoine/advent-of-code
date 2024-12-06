package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

type Position struct {
	x, y int
}

type Guard struct {
	position  Position
	direction rune
}

func (g *Guard) Next(grid [][]rune) {
	switch g.direction {
	case '^':
		g.position.y--
	case 'v':
		g.position.y++
	case '<':
		g.position.x--
	case '>':
		g.position.x++
	}
	if g.position.x < 0 || g.position.x >= len(grid[0]) || g.position.y < 0 || g.position.y >= len(grid) {
		return
	}
	if grid[g.position.y][g.position.x] == '#' {
		switch g.direction {
		case '^':
			g.position.y++
			g.position.x++
			g.direction = '>'
		case 'v':
			g.position.y--
			g.position.x--
			g.direction = '<'
		case '<':
			g.position.x++
			g.position.y--
			g.direction = '^'
		case '>':
			g.position.x--
			g.position.y++
			g.direction = 'v'
		}
	}
}

var directions = []rune{'^', '>', 'v', '<'}

func parseInput(input io.Reader) ([][]rune, Guard) {
	scanner := bufio.NewScanner(input)

	var grid [][]rune
	var guard Guard
	y := 0
	for scanner.Scan() {
		line := []rune(scanner.Text())
		for x, char := range line {
			if slices.Contains(directions, char) {
				guard = Guard{
					position:  Position{x, y},
					direction: char,
				}
				line[x] = '.'
				break
			}
		}
		grid = append(grid, line)
		y++
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid, guard
}

func getResult(input io.Reader) int {
	grid, guard := parseInput(input)

	visited := make(map[Position]struct{})
	previousState := make(map[Guard]struct{})
	for guard.position.x >= 0 && guard.position.x < len(grid[0]) && guard.position.y >= 0 && guard.position.y < len(grid) {
		visited[guard.position] = struct{}{}
		if _, ok := previousState[guard]; ok {
			break
		}
		previousState[guard] = struct{}{}
		guard.Next(grid)
	}

	return len(visited)
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
