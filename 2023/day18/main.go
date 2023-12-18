package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Direction int

const (
	Up    Direction = iota
	Down  Direction = iota
	Left  Direction = iota
	Right Direction = iota
)

type Position struct {
	i, j int64
}

func (p Position) Move(direction Direction) Position {
	switch direction {
	case Up:
		return Position{p.i - 1, p.j}
	case Down:
		return Position{p.i + 1, p.j}
	case Left:
		return Position{p.i, p.j - 1}
	case Right:
		return Position{p.i, p.j + 1}
	}
	log.Fatalf("Unknown direction: %d", direction)
	return p
}

var lineRegex = regexp.MustCompile(`^([UDLR]) (\d+) \(#([0-9a-f]{5})([0-9a-f])\)$`)

func parseLine(line string) (Direction, int64) {
	matches := lineRegex.FindStringSubmatch(line)
	if matches == nil || len(matches) != 5 {
		log.Fatalf("Unable to parse line: %s", line)
	}

	distance, errParsingDistance := strconv.ParseInt(matches[3], 16, 64)
	if errParsingDistance != nil {
		log.Fatalf("Unable to parse distance: %s", matches[2])
	}

	// 0 means R, 1 means D, 2 means L, and 3 means U.
	directionStr := matches[4]
	var direction Direction
	switch directionStr {
	case "0":
		direction = Right
	case "1":
		direction = Down
	case "2":
		direction = Left
	case "3":
		direction = Up
	}

	return direction, distance
}

func parseInput(input io.Reader) ([]Position, int64) {
	scanner := bufio.NewScanner(input)

	var vertices []Position
	currentPosition := Position{0, 0}
	vertices = append(vertices, currentPosition)
	var boundaryPoints int64
	for scanner.Scan() {
		direction, distance := parseLine(scanner.Text())
		boundaryPoints += distance
		for i := int64(0); i < distance; i++ {
			currentPosition = currentPosition.Move(direction)
		}
		vertices = append(vertices, currentPosition)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return vertices, boundaryPoints
}

func getResult(input io.Reader) int64 {
	vertices, boundaryPoints := parseInput(input)

	// Compute area using the shoelace formula
	// https://en.wikipedia.org/wiki/Shoelace_formula
	var area int64
	for i := int64(0); i < int64(len(vertices))-1; i++ {
		pointA := vertices[i]
		pointB := vertices[i+1]
		area += (pointA.j * pointB.i) - (pointA.i * pointB.j)
	}
	area /= 2

	// Computed points
	// https://en.wikipedia.org/wiki/Pick%27s_theorem
	return area - (boundaryPoints / 2) + 1 + boundaryPoints
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
