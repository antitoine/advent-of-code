package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseLine(line string) image.Point {
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		log.Fatalf("Invalid line: %s", line)
	}
	x, errX := strconv.Atoi(parts[0])
	if errX != nil {
		log.Fatalf("Invalid x: %s", parts[0])
	}
	y, errY := strconv.Atoi(parts[1])
	if errY != nil {
		log.Fatalf("Invalid y: %s", parts[1])
	}
	return image.Point{X: x, Y: y}
}

func parseInput(input io.Reader) []image.Point {
	scanner := bufio.NewScanner(input)

	var corruptedBytes []image.Point
	for scanner.Scan() {
		corruptedBytes = append(corruptedBytes, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return corruptedBytes
}

func printMap(space image.Rectangle, corruptedBytes []image.Point) {
	corruptedBytesMap := make(map[image.Point]struct{})
	for c := 0; c < len(corruptedBytes); c++ {
		corruptedBytesMap[corruptedBytes[c]] = struct{}{}
	}

	for y := 0; y < space.Dy(); y++ {
		for x := 0; x < space.Dx(); x++ {
			if _, ok := corruptedBytesMap[image.Point{X: x, Y: y}]; ok {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

type Path struct {
	Length int64
	Pos    image.Point
}

type PathHeap []Path

func (h PathHeap) Len() int           { return len(h) }
func (h PathHeap) Less(i, j int) bool { return h[i].Length < h[j].Length }
func (h PathHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PathHeap) Push(x any) {
	*h = append(*h, x.(Path))
}

func (h *PathHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var directions = []image.Point{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}}

func shortestPath(space image.Rectangle, corruptedBytes []image.Point) int64 {
	visited := make(map[image.Point]struct{})
	visited[space.Min] = struct{}{}

	paths := &PathHeap{Path{Length: 0, Pos: space.Min}}
	heap.Init(paths)

	corruptedBytesMap := make(map[image.Point]struct{})
	for _, corruptedByte := range corruptedBytes {
		corruptedBytesMap[corruptedByte] = struct{}{}
	}

	dest := space.Max.Sub(image.Point{X: 1, Y: 1})
	for paths.Len() > 0 {
		path := heap.Pop(paths).(Path)
		if path.Pos == dest {
			return path.Length
		}
		for _, dir := range directions {
			nextPos := path.Pos.Add(dir)
			if !nextPos.In(space) {
				continue
			}
			if _, ok := corruptedBytesMap[nextPos]; ok {
				continue
			}
			if _, ok := visited[nextPos]; ok {
				continue
			}
			visited[nextPos] = struct{}{}
			heap.Push(paths, Path{Length: path.Length + 1, Pos: nextPos})
		}
	}

	return -1
}

func getResult(input io.Reader, space image.Rectangle, nbCorruptedBytes int) int64 {
	corruptedBytes := parseInput(input)

	printMap(space, corruptedBytes[:nbCorruptedBytes])

	return shortestPath(space, corruptedBytes[:nbCorruptedBytes])
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

	result := getResult(inputFile, image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(71, 71)}, 1024)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
