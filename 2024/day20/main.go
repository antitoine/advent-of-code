package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

type Track struct {
	track [][]rune
	start image.Point
	end   image.Point
	space image.Rectangle
}

func parseInput(input io.Reader) Track {
	scanner := bufio.NewScanner(input)

	var track [][]rune
	var start, end image.Point
	for y := 0; scanner.Scan(); y++ {
		line := []rune(scanner.Text())
		if startIdx := slices.Index(line, 'S'); startIdx != -1 {
			start = image.Pt(startIdx, y)
			line[startIdx] = '.'
		}
		if endIdx := slices.Index(line, 'E'); endIdx != -1 {
			end = image.Pt(endIdx, y)
			line[endIdx] = '.'
		}
		track = append(track, line)
	}
	space := image.Rect(0, 0, len(track[0]), len(track))

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return Track{
		track: track,
		start: start,
		end:   end,
		space: space,
	}
}

var directions = []image.Point{
	image.Pt(0, -1),
	image.Pt(-1, 0),
	image.Pt(0, 1),
	image.Pt(1, 0),
}

func getNormalPath(track Track) []image.Point {
	steps := []image.Point{
		track.start,
	}

	previousPos := image.Pt(-1, -1)
	currentPos := track.start
	for currentPos != track.end {
		for _, direction := range directions {
			nextPos := currentPos.Add(direction)
			if !nextPos.In(track.space) || previousPos == nextPos || track.track[nextPos.Y][nextPos.X] == '#' {
				continue
			}
			previousPos = currentPos
			currentPos = nextPos
			steps = append(steps, nextPos)
			break
		}
	}

	return steps
}

func manhattanDistance(p1, p2 image.Point) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getAllPossiblePathWithCheat(normalPath []image.Point, nbCheat int, minReductionSteps int) int64 {
	positions := make([]image.Point, 0, len(normalPath))
	distances := make([]int, 0, len(normalPath))

	for i, pos := range normalPath {
		positions = append(positions, pos)
		distances = append(distances, i)
	}

	var result int64
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			distance := manhattanDistance(positions[i], positions[j])
			diff := distances[j] - distances[i] - distance
			if distance == nbCheat && diff >= minReductionSteps {
				result++
			}
		}
	}

	return result
}

func getResult(input io.Reader, maxCheats int, nbLeastSavingSteps int) int64 {
	track := parseInput(input)

	path := getNormalPath(track)

	return getAllPossiblePathWithCheat(path, maxCheats, nbLeastSavingSteps)
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

	result := getResult(inputFile, 2, 100)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
