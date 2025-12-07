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

	// Track count of timelines at each column position
	timelines := make(map[int]int64)
	timelines[startCol] = 1

	// Process rows starting from the row after S
	for row := startRow + 1; row < len(grid); row++ {
		if len(timelines) == 0 {
			break
		}

		newTimelines := make(map[int]int64)
		line := grid[row]

		for col, count := range timelines {
			if col < 0 || col >= width {
				continue
			}

			ch := line[col]
			if ch == '^' {
				leftCol := col - 1
				rightCol := col + 1
				if leftCol >= 0 {
					newTimelines[leftCol] += count
				}
				if rightCol < width {
					newTimelines[rightCol] += count
				}
			} else {
				newTimelines[col] += count
			}
		}

		timelines = newTimelines
	}

	// Sum all timelines
	var total int64
	for _, count := range timelines {
		total += count
	}

	return total
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
