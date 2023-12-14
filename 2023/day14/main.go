package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

type Place string

const (
	Empty          Place = "."
	CubeShapedRock Place = "#"
	RoundedRock    Place = "O"
)

type Platform [][]Place

func (p Platform) String() string {
	var result string
	for _, row := range p {
		for _, place := range row {
			result += string(place)
		}
		result += "\n"
	}
	return result
}

func parseInput(input io.Reader) Platform {
	scanner := bufio.NewScanner(input)
	var platform Platform
	for scanner.Scan() {
		line := scanner.Text()
		var row []Place
		for _, char := range line {
			row = append(row, Place(char))
		}
		platform = append(platform, row)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return platform
}

func tiltingTheLever(platform Platform) Platform {
	for columnIdx := 0; columnIdx < len(platform[0]); columnIdx++ {
		var nextBlockingRow int
		for rowIdx := 0; rowIdx < len(platform); rowIdx++ {
			switch platform[rowIdx][columnIdx] {
			case CubeShapedRock:
				nextBlockingRow = rowIdx + 1
			case RoundedRock:
				platform[rowIdx][columnIdx] = Empty
				platform[nextBlockingRow][columnIdx] = RoundedRock
				nextBlockingRow++
			}
		}
	}
	return platform
}

func computeLoad(platform Platform) int {
	var load int
	for rowIdx := 0; rowIdx < len(platform); rowIdx++ {
		for _, place := range platform[rowIdx] {
			if place == RoundedRock {
				load += len(platform) - rowIdx
			}
		}
	}
	return load
}

func getResult(input io.Reader) int {
	platform := parseInput(input)
	//log.Printf("Initial platform:\n%s", platform)
	tiltedPlatform := tiltingTheLever(platform)
	//log.Printf("Tilted platform:\n%s", tiltedPlatform)
	return computeLoad(tiltedPlatform)
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
