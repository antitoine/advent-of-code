package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func transpose(grid [][]int64) [][]int64 {
	result := make([][]int64, len(grid[0]))
	for i := range result {
		result[i] = make([]int64, len(grid))
	}

	for i, row := range grid {
		for j, value := range row {
			result[j][i] = value
		}
	}

	return result
}

func isEqual(line1 []int64, line2 []int64) bool {
	for i := range line1 {
		if line1[i] != line2[i] {
			return false
		}
	}
	return true
}

func getReflectionScoreForLines(grid [][]int64) int64 {
	for i, j := 0, 1; j < len(grid); i, j = i+1, j+1 {
		if isEqual(grid[i], grid[j]) {
			found := true
			for k := 0; i-k >= 0 && j+k < len(grid); k++ {
				if !isEqual(grid[i-k], grid[j+k]) {
					found = false
					break
				}
			}
			if found {
				return int64(i) + 1
			}
		}
	}
	return 0
}

func getReflectionScore(grid [][]int64) int64 {
	result := getReflectionScoreForLines(grid)
	if result > 0 {
		return result * 100
	}
	result = getReflectionScoreForLines(transpose(grid))
	if result > 0 {
		return result
	}
	log.Fatalf("Unable to find a reflection score for the grid: %v", grid)
	return 0
}

func parseInput(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	reflectionScore := int64(0)
	var grid [][]int64
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(grid) > 0 {
				reflectionScore += getReflectionScore(grid)
				grid = [][]int64{}
			}
			continue
		}

		var row []int64
		for _, char := range line {
			if char == '#' {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		grid = append(grid, row)
	}

	if len(grid) > 0 {
		reflectionScore += getReflectionScore(grid)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return reflectionScore
}

func getResult(input io.Reader) int64 {
	return parseInput(input)
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
