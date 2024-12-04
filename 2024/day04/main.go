package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func parseInput(input io.Reader) []string {
	scanner := bufio.NewScanner(input)

	var grid []string
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid
}

func isValidMatrix(grid []string, i, j int) bool {
	if i+2 >= len(grid) || j+2 >= len(grid[i]) {
		return false
	}
	if grid[i][j] != 'M' && grid[i][j] != 'S' {
		return false
	}
	if grid[i+1][j+1] != 'A' {
		return false
	}
	if grid[i][j] == 'M' && grid[i][j+2] == 'S' && grid[i+2][j] == 'M' && grid[i+2][j+2] == 'S' {
		return true
	} else if grid[i][j] == 'S' && grid[i][j+2] == 'M' && grid[i+2][j] == 'S' && grid[i+2][j+2] == 'M' {
		return true
	} else if grid[i][j] == 'M' && grid[i][j+2] == 'M' && grid[i+2][j] == 'S' && grid[i+2][j+2] == 'S' {
		return true
	} else if grid[i][j] == 'S' && grid[i][j+2] == 'S' && grid[i+2][j] == 'M' && grid[i+2][j+2] == 'M' {
		return true
	}
	return false
}

func getResult(input io.Reader) int {
	grid := parseInput(input)
	var cnt int
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if isValidMatrix(grid, i, j) {
				cnt++
			}
		}
	}
	return cnt
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
