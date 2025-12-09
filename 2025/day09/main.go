package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

func parseInput(input io.Reader) []Point {
	var points []Point
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		points = append(points, Point{x, y})
	}
	return points
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getResult(input io.Reader) int64 {
	points := parseInput(input)

	maxArea := int64(0)
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			width := abs(points[j].x-points[i].x) + 1
			height := abs(points[j].y-points[i].y) + 1
			area := int64(width) * int64(height)
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
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
