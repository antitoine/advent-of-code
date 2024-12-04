package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
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

func getVerticalWord(grid []string, i, j int) string {
	var word string
	for k := 0; k < 4 && i < len(grid); k++ {
		word += string(grid[i][j])
		i++
	}
	return word
}

func getDiagonalRightWord(grid []string, i, j int) string {
	var word string
	for k := 0; k < 4 && i < len(grid) && j < len(grid[i]); k++ {
		word += string(grid[i][j])
		i++
		j++
	}
	return word
}

func getDiagonalLeftWord(grid []string, i, j int) string {
	var word string
	for k := 0; k < 4 && i < len(grid) && j >= 0; k++ {
		word += string(grid[i][j])
		i++
		j--
	}
	return word
}

func cntWords(grid []string, i, j int) int {
	if grid[i][j] != 'X' && grid[i][j] != 'S' {
		return 0
	}
	var cnt int
	verticalWord := getVerticalWord(grid, i, j)
	diagonalRightWord := getDiagonalRightWord(grid, i, j)
	diagonalLeftWord := getDiagonalLeftWord(grid, i, j)
	if grid[i][j] == 'X' {
		if strings.HasPrefix(grid[i][j:], "XMAS") {
			cnt++
		}
		if strings.HasPrefix(verticalWord, "XMAS") {
			cnt++
		}
		if strings.HasPrefix(diagonalRightWord, "XMAS") {
			cnt++
		}
		if strings.HasPrefix(diagonalLeftWord, "XMAS") {
			cnt++
		}
	} else {
		if strings.HasPrefix(grid[i][j:], "SAMX") {
			cnt++
		}
		if strings.HasPrefix(verticalWord, "SAMX") {
			cnt++
		}
		if strings.HasPrefix(diagonalRightWord, "SAMX") {
			cnt++
		}
		if strings.HasPrefix(diagonalLeftWord, "SAMX") {
			cnt++
		}
	}
	return cnt
}

func getResult(input io.Reader) int {
	grid := parseInput(input)
	var cnt int
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			cnt += cntWords(grid, i, j)
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
