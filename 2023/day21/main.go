package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Grid [][]rune

type Position struct {
	i int64
	j int64
}

func parseInput(input io.Reader) (Grid, Position) {
	scanner := bufio.NewScanner(input)

	var grid Grid
	var start Position
	for i := int64(0); scanner.Scan(); i++ {
		line := scanner.Text()
		if idx := strings.Index(line, "S"); idx >= 0 {
			start = Position{i, int64(idx)}
			line = strings.Replace(line, "S", ".", 1)
		}
		grid = append(grid, []rune(line))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid, start
}

func (g Grid) isValidMove(position Position) bool {
	rows, cols := int64(len(g)), int64(len(g[0]))
	return position.i >= 0 && position.i < rows && position.j >= 0 && position.j < cols && g[position.i][position.j] == '.'
}

func (g Grid) countReachablePositions(start Position, moves uint64) int64 {
	directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	exploration := []Position{start}

	// BFS
	movesLeft := moves
	for movesLeft > 0 {
		nextUniqueExploration := make(map[Position]bool)
		for i := 0; i < len(exploration); i++ {
			current := exploration[i]
			for _, dir := range directions {
				next := Position{current.i + dir.i, current.j + dir.j}
				if g.isValidMove(next) {
					nextUniqueExploration[next] = true
				}
			}
		}
		nextExploration := make([]Position, 0, len(nextUniqueExploration))
		for next := range nextUniqueExploration {
			nextExploration = append(nextExploration, next)
		}
		exploration = nextExploration
		movesLeft--
	}

	return int64(len(exploration))
}

func getResultPart1(input io.Reader, moves uint64) int64 {
	grid, start := parseInput(input)
	return grid.countReachablePositions(start, moves)
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

	result := getResultPart1(inputFile, 64)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
