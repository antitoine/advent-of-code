package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func parseLine(line string, y int) ([]image.Point, []int) {
	var startingPoints []image.Point
	lineNumbers := make([]int, len(line))
	for i, numStr := range line {
		num, errParsing := strconv.Atoi(string(numStr))
		if errParsing != nil {
			log.Fatalf("Unable to parse the number '%s': %v", string(numStr), errParsing)
		}
		if num == 0 {
			startingPoints = append(startingPoints, image.Point{X: i, Y: y})
		}
		lineNumbers[i] = num
	}
	return startingPoints, lineNumbers
}

func parseInput(input io.Reader) ([]image.Point, [][]int) {
	scanner := bufio.NewScanner(input)

	var startingPoints []image.Point
	var grid [][]int
	for y := 0; scanner.Scan(); y++ {
		newStartingPoints, row := parseLine(scanner.Text(), y)
		grid = append(grid, row)
		startingPoints = append(startingPoints, newStartingPoints...)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return startingPoints, grid
}

var directions = []image.Point{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}}

func countTrailheadScore(startingPoint image.Point, grid [][]int) int64 {
	if grid[startingPoint.Y][startingPoint.X] == 9 {
		return 1
	}
	var score int64
	for _, direction := range directions {
		neighbour := startingPoint.Add(direction)
		if neighbour.In(image.Rect(0, 0, len(grid[0]), len(grid))) && grid[neighbour.Y][neighbour.X] == grid[startingPoint.Y][startingPoint.X]+1 {
			score += countTrailheadScore(neighbour, grid)
		}
	}
	return score
}

func getResult(input io.Reader) int64 {
	startingPoints, grid := parseInput(input)
	var sumScore int64
	for _, startingPoint := range startingPoints {
		score := countTrailheadScore(startingPoint, grid)
		sumScore += score
	}
	return sumScore
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
