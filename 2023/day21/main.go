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

func modGrid(idx, length int64) int64 {
	m := idx % length
	if m < 0 {
		m += length
	}
	return m
}

func (p Position) modulo(g Grid) Position {
	return Position{
		i: modGrid(p.i, int64(len(g))),
		j: modGrid(p.j, int64(len(g[0]))),
	}
}

func (g Grid) countReachablePositions(start Position, moves int, infiniteGrid bool) int64 {
	directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	exploration := []Position{start}
	prevExplorationSize := 0

	var values []int
	var delta1 []int
	var delta2 []int

	maxCompute := 1000

	// BFS
	movesLeft := moves
	for i := 0; i < maxCompute && movesLeft > 0; i++ {
		nextUniqueExploration := make(map[Position]bool)
		for i := 0; i < len(exploration); i++ {
			current := exploration[i]
			for _, dir := range directions {
				next := Position{current.i + dir.i, current.j + dir.j}
				var isValid bool
				if infiniteGrid {
					isValid = g.isValidMove(next.modulo(g))
				} else {
					isValid = g.isValidMove(next)
				}
				if isValid {
					nextUniqueExploration[next] = true
				}
			}
		}
		nextExploration := make([]Position, 0, len(nextUniqueExploration))
		for next := range nextUniqueExploration {
			nextExploration = append(nextExploration, next)
		}
		newExplorationSize := len(nextExploration)
		prevExplorationSize = len(exploration)
		exploration = nextExploration
		movesLeft--

		values = append(values, newExplorationSize)
		delta1 = append(delta1, newExplorationSize-prevExplorationSize)
		if len(delta1) > len(g) {
			delta2 = append(delta2, delta1[len(delta1)-1]-delta1[len(delta1)-1-len(g)])
		} else {
			delta2 = append(delta2, 0)
		}
	}

	if movesLeft <= 0 {
		return int64(len(exploration))
	}

	restartAt := maxCompute - len(g)

	values = values[:restartAt]
	delta1 = delta1[:restartAt]
	delta2 = delta2[:restartAt]

	result := values[len(values)-1]
	for i := restartAt + 1; i < moves+1; i++ {
		delta2ForI := delta2[i-len(g)-1]
		delta1ForI := delta1[i-len(g)-1] + delta2ForI
		result += delta1ForI

		delta1 = append(delta1, delta1ForI)
		delta2 = append(delta2, delta2ForI)
	}

	return int64(result)
}

func (g Grid) String() string {
	var sb strings.Builder
	for _, line := range g {
		sb.WriteString(string(line))
		sb.WriteString("\n")
	}
	return sb.String()
}

func getResultPart1(input io.Reader, moves int) int64 {
	grid, start := parseInput(input)
	return grid.countReachablePositions(start, moves, false)
}

func getResultPart2(input io.Reader, moves int) int64 {
	grid, start := parseInput(input)
	return grid.countReachablePositions(start, moves, true)
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

	result := getResultPart2(inputFile, 26501365)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
