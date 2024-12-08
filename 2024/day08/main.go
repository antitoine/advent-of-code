package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

type Position struct {
	x, y int
}

type Frequency rune

func parseInput(input io.Reader) (map[Frequency][]Position, int, int) {
	scanner := bufio.NewScanner(input)

	antennas := make(map[Frequency][]Position)
	y := 0
	x := 0
	for scanner.Scan() {
		line := scanner.Text()
		x = len(line)
		for x, char := range line {
			if char == '.' {
				continue
			}
			antennas[Frequency(char)] = append(antennas[Frequency(char)], Position{x, y})
		}
		y++
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return antennas, x, y
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getResult(input io.Reader) int {
	antennas, xSize, ySize := parseInput(input)
	antinodes := map[Position]struct{}{}
	for _, positions := range antennas {
		for a1 := 0; a1 < len(positions)-1; a1++ {
			antenna1 := positions[a1]
			for a2 := a1 + 1; a2 < len(positions); a2++ {
				antenna2 := positions[a2]
				firstAntinode := Position{
					x: 2*antenna2.x - antenna1.x,
					y: 2*antenna2.y - antenna1.y,
				}
				secondAntinode := Position{
					x: 2*antenna1.x - antenna2.x,
					y: 2*antenna1.y - antenna2.y,
				}
				if firstAntinode.x >= 0 && firstAntinode.x < xSize && firstAntinode.y >= 0 && firstAntinode.y < ySize {
					antinodes[firstAntinode] = struct{}{}
				}
				if secondAntinode.x >= 0 && secondAntinode.x < xSize && secondAntinode.y >= 0 && secondAntinode.y < ySize {
					antinodes[secondAntinode] = struct{}{}
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
