package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"time"
)

type Point struct {
	x int
	y int
}

func parseInput(input io.Reader) ([]Point, []int, []int) {
	scanner := bufio.NewScanner(input)
	var galaxies []Point
	var emptyX []int
	var emptyYChecking []bool
	for x := 0; scanner.Scan(); x++ {
		line := scanner.Text()
		isEmpty := true
		for y, symbol := range line {
			if len(emptyYChecking) <= y {
				emptyYChecking = append(emptyYChecking, true)
			}
			if symbol == '#' {
				galaxies = append(galaxies, Point{x: x, y: y})
				isEmpty = false
				emptyYChecking[y] = false
			}
		}
		if isEmpty {
			emptyX = append(emptyX, x)
		}
	}
	var emptyY []int
	for y, isEmpty := range emptyYChecking {
		if isEmpty {
			emptyY = append(emptyY, y)
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return galaxies, emptyX, emptyY
}

func getDistance(galaxyA, galaxyB Point, emptyX []int, emptyY []int, emptyFactor float64) int64 {
	xDistance := math.Abs(float64(galaxyB.x) - float64(galaxyA.x))
	for _, x := range emptyX {
		if (galaxyA.x < x && x < galaxyB.x) || (galaxyA.x > x && x > galaxyB.x) {
			xDistance += emptyFactor - 1
		}
	}
	yDistance := math.Abs(float64(galaxyB.y) - float64(galaxyA.y))
	for _, y := range emptyY {
		if (galaxyA.y < y && y < galaxyB.y) || (galaxyA.y > y && y > galaxyB.y) {
			yDistance += emptyFactor - 1
		}
	}
	return int64(xDistance + yDistance)
}

func getResult(input io.Reader, emptyFactor float64) int64 {
	galaxies, emptyX, emptyY := parseInput(input)
	var distances int64
	for i, galaxyA := range galaxies {
		for _, galaxyB := range galaxies[i:] {
			if galaxyA.x == galaxyB.x && galaxyA.y == galaxyB.y {
				continue
			}
			distances += getDistance(galaxyA, galaxyB, emptyX, emptyY, emptyFactor)
		}
	}
	return distances
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

	emptyFactor := float64(1000000)
	result := getResult(inputFile, emptyFactor)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
