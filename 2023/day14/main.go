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

// rotate returns the same platform rotated 90 degrees clockwise
func rotate(platform Platform) Platform {
	var rotatedPlatform Platform
	for columnIdx := 0; columnIdx < len(platform[0]); columnIdx++ {
		var row []Place
		for rowIdx := len(platform) - 1; rowIdx >= 0; rowIdx-- {
			row = append(row, platform[rowIdx][columnIdx])
		}
		rotatedPlatform = append(rotatedPlatform, row)
	}
	return rotatedPlatform
}

var loadMemory = make(map[string]int)

func computeLoad(platform Platform) int {
	if load, ok := loadMemory[platform.String()]; ok {
		return load
	}
	var load int
	for rowIdx := 0; rowIdx < len(platform); rowIdx++ {
		for _, place := range platform[rowIdx] {
			if place == RoundedRock {
				load += len(platform) - rowIdx
			}
		}
	}
	loadMemory[platform.String()] = load
	return load
}

var cycleMemory = make(map[string]Platform)

func cycle(initPlatform Platform) Platform {
	if platform, ok := cycleMemory[initPlatform.String()]; ok {
		return platform
	}
	platform := tiltingTheLever(initPlatform)    // tilt north
	platform = tiltingTheLever(rotate(platform)) // tilt west
	platform = tiltingTheLever(rotate(platform)) // tilt south
	platform = tiltingTheLever(rotate(platform)) // tilt east
	platform = rotate(platform)                  // get initial position
	cycleMemory[initPlatform.String()] = platform
	return platform
}

var platformIdxMemory = make(map[string]int)

func getResult(input io.Reader) int {
	initPlatform := parseInput(input)

	platform := initPlatform
	var cycleStart int
	var cycleEnd int

	for i := 0; i < 1000000000; i++ {
		platform = cycle(platform)

		if existingIdx, found := platformIdxMemory[platform.String()]; found {
			cycleStart = existingIdx
			cycleEnd = i
			break
		}
		platformIdxMemory[platform.String()] = i
	}
	cycleLength := cycleEnd - cycleStart

	remainingIterations := 1000000000 - cycleEnd - 1
	remainingCycles := remainingIterations % cycleLength
	for i := 0; i < remainingCycles; i++ {
		platform = cycle(platform)
	}

	log.Printf("Final platform:\n%s", platform)
	return computeLoad(platform)
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
