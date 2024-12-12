package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func parseLine(line string) {
	strings.Split(line, " ")
}

func parseInput(input io.Reader) [][]rune {
	scanner := bufio.NewScanner(input)

	var grid [][]rune
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid
}

func isValidPoint(grid [][]rune, point image.Point) bool {
	return point.In(image.Rect(0, 0, len(grid[0]), len(grid)))
}

var directions = []image.Point{
	image.Pt(0, 1),
	image.Pt(1, 0),
	image.Pt(0, -1),
	image.Pt(-1, 0),
}

type Side struct {
	direction image.Point
	point     image.Point
}

func getZoneArea(grid [][]rune, point image.Point, visited map[image.Point]struct{}, sides map[Side]struct{}) int64 {
	var area int64
	for _, direction := range directions {
		neighbour := point.Add(direction)
		if !isValidPoint(grid, neighbour) || grid[neighbour.Y][neighbour.X] != grid[point.Y][point.X] {
			sides[Side{
				direction: direction,
				point:     point,
			}] = struct{}{}
			continue
		}
		if _, ok := visited[neighbour]; ok {
			continue
		}
		visited[neighbour] = struct{}{}

		neighbourArea := getZoneArea(grid, neighbour, visited, sides)
		area += neighbourArea
	}

	return area + 1
}

func getAxisAndPosition(direction image.Point, point image.Point) (int, int) {
	if direction.X == 0 {
		return point.Y, point.X
	} else {
		return point.X, point.Y
	}
}

func countDistinctSides(sides map[Side]struct{}) int64 {
	var sidesCnt int64
	for _, direction := range directions {
		sameSidePoints := make(map[int][]int)
		for side := range sides {
			if side.direction == direction {
				axis, position := getAxisAndPosition(direction, side.point)
				sameSidePoints[axis] = append(sameSidePoints[axis], position)
			}
		}
		for _, positions := range sameSidePoints {
			sort.Ints(positions)
			sidesCnt++
			for i := 1; i < len(positions); i++ {
				if positions[i]-positions[i-1] > 1 {
					sidesCnt++
				}
			}
		}
	}
	return sidesCnt
}

func getResult(input io.Reader) int64 {
	grid := parseInput(input)
	visited := make(map[image.Point]struct{})
	var result int64
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			point := image.Pt(x, y)
			if _, ok := visited[point]; ok {
				continue
			}
			visited[point] = struct{}{}
			sides := make(map[Side]struct{})
			area := getZoneArea(grid, point, visited, sides)
			result += area * countDistinctSides(sides)
		}
	}
	return result
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
