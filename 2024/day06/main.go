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

func (g Guard) NextPosition() Position {
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
	return g.position
}

func (g Guard) Next(grid [][]rune) Guard {
	g.position = g.NextPosition()
	if g.position.x < 0 || g.position.x >= len(grid[0]) || g.position.y < 0 || g.position.y >= len(grid) {
		return g
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
	return g
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

	startingPosition := guard.position

	obstaclesForLoop := make(map[Position]struct{})

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			obstacle := Position{x, y}
			if grid[obstacle.y][obstacle.x] == '.' && startingPosition != obstacle {
				grid[obstacle.y][obstacle.x] = '#'
				previousState := make(map[Guard]struct{})
				previousState[guard] = struct{}{}
				guardWithObstacle := guard.Next(grid)
				for guardWithObstacle.position.y >= 0 && guardWithObstacle.position.y < len(grid) && guardWithObstacle.position.x >= 0 && guardWithObstacle.position.x < len(grid[guardWithObstacle.position.y]) {
					if _, ok := previousState[guardWithObstacle]; ok {
						obstaclesForLoop[obstacle] = struct{}{}
						break
					}
					previousState[guardWithObstacle] = struct{}{}
					guardWithObstacle = guardWithObstacle.Next(grid)
				}
				grid[obstacle.y][obstacle.x] = '.'
			}
		}
	}

	return len(obstaclesForLoop)
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

	// 1980 Wrong too high
	// 1928 Good
	// 1909 Wrong oo low
	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
