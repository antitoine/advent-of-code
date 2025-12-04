package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func getResult(input io.Reader) int64 {
	var grid [][]byte
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])

	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	var totalRemoved int64
	for {
		var toRemove [][2]int
		for r := range rows {
			for c := range cols {
				if grid[r][c] != '@' {
					continue
				}

				adjacentRolls := 0
				for _, d := range directions {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
						adjacentRolls++
					}
				}

				if adjacentRolls < 4 {
					toRemove = append(toRemove, [2]int{r, c})
				}
			}
		}

		if len(toRemove) == 0 {
			break
		}

		for _, pos := range toRemove {
			grid[pos[0]][pos[1]] = '.'
		}
		totalRemoved += int64(len(toRemove))
	}

	return totalRemoved
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
