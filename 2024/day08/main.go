package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"time"
)

func parseInput(input io.Reader) (map[rune][]image.Point, int, int) {
	scanner := bufio.NewScanner(input)

	antennas := make(map[rune][]image.Point)
	y := 0
	x := 0
	for scanner.Scan() {
		line := scanner.Text()
		x = len(line)
		for x, char := range line {
			if char == '.' {
				continue
			}
			antennas[char] = append(antennas[char], image.Point{X: x, Y: y})
		}
		y++
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return antennas, x, y
}

func isValidPosition(pos image.Point, xSize, ySize int) bool {
	return pos.X >= 0 && pos.X < xSize && pos.Y >= 0 && pos.Y < ySize
}

func getResult(input io.Reader) int {
	antennas, xSize, ySize := parseInput(input)
	antinodes := map[image.Point]struct{}{}
	for _, positions := range antennas {
		for _, antenna1 := range positions {
			for _, antenna2 := range positions {
				if antenna1 == antenna2 {
					continue
				}

				for d := antenna2.Sub(antenna1); isValidPosition(antenna2, xSize, ySize); antenna2 = antenna2.Add(d) {
					antinodes[antenna2] = struct{}{}
				}
			}
		}

	}
	return len(antinodes)
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
