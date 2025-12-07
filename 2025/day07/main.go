package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	var grid []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			grid = append(grid, line)
		}
	}

	if len(grid) == 0 {
		return 0
	}

	// Find starting position 'S'
	startCol := -1
	startRow := -1
	for row, line := range grid {
		for col, ch := range line {
			if ch == 'S' {
				startCol = col
				startRow = row
				break
			}
		}
		if startCol != -1 {
			break
		}
	}

	if startCol == -1 {
		return 0
	}

	width := len(grid[0])

	// Track active beam columns using a map as a set
	activeBeams := make(map[int]bool)
	activeBeams[startCol] = true

	splitCount := int64(0)

	// Process rows starting from the row after S
	for row := startRow + 1; row < len(grid); row++ {
		if len(activeBeams) == 0 {
			break
		}

		newBeams := make(map[int]bool)
		line := grid[row]

		for col := range activeBeams {
			if col < 0 || col >= width {
				continue
			}

			ch := line[col]
			if ch == '^' {
				splitCount++
				leftCol := col - 1
				rightCol := col + 1
				if leftCol >= 0 {
					newBeams[leftCol] = true
				}
				if rightCol < width {
					newBeams[rightCol] = true
				}
			} else {
				newBeams[col] = true
			}
		}

		activeBeams = newBeams
	}

	return splitCount
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
