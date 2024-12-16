package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"time"
)

type Board struct {
	cells [][]rune
	space image.Rectangle
	start image.Point
	end   image.Point
}

func (b *Board) Print() {
	for _, row := range b.cells {
		fmt.Println(string(row))
	}
}

func parseInput(input io.Reader) *Board {
	scanner := bufio.NewScanner(input)

	board := &Board{}

	for y := 0; scanner.Scan(); y++ {
		var row []rune
		for x, char := range scanner.Text() {
			row = append(row, char)
			if char == 'S' {
				board.start = image.Pt(x, y)
			} else if char == 'E' {
				board.end = image.Pt(x, y)
			}
		}
		board.cells = append(board.cells, row)
	}

	board.space = image.Rect(0, 0, len(board.cells[0]), len(board.cells))

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return board
}

type Direction rune

const (
	North Direction = 'N'
	West  Direction = 'W'
	South Direction = 'S'
	East  Direction = 'E'
)

var directions = map[Direction]image.Point{
	North: image.Pt(0, -1),
	West:  image.Pt(-1, 0),
	South: image.Pt(0, 1),
	East:  image.Pt(1, 0),
}

func nextDirection(dir Direction) Direction {
	switch dir {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	}
	return 0
}

type Path struct {
	score int64
	dir   Direction
	pos   image.Point
}

type SmallestPathHeap []Path

func (h SmallestPathHeap) Len() int           { return len(h) }
func (h SmallestPathHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h SmallestPathHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *SmallestPathHeap) Push(x any) {
	*h = append(*h, x.(Path))
}

func (h *SmallestPathHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func smallestPath(board *Board) int64 {
	paths := &SmallestPathHeap{}
	heap.Push(paths, Path{0, East, board.start})
	visited := make(map[image.Point]struct{})

	for paths.Len() > 0 {
		path := heap.Pop(paths).(Path)

		if path.pos == board.end {
			return path.score
		}

		nextPos := path.pos.Add(directions[path.dir])
		if nextPos.In(board.space) && board.cells[nextPos.Y][nextPos.X] != '#' {
			if _, ok := visited[nextPos]; !ok {
				visited[nextPos] = struct{}{}
				heap.Push(paths, Path{path.score + 1, path.dir, nextPos})
			}
		}

		nextDir := nextDirection(path.dir)
		nextPos = path.pos.Add(directions[nextDir])
		if nextPos.In(board.space) && board.cells[nextPos.Y][nextPos.X] != '#' {
			if _, ok := visited[nextPos]; !ok {
				visited[nextPos] = struct{}{}
				heap.Push(paths, Path{path.score + 1001, nextDir, nextPos})
			}
		}

		nextDir = nextDirection(nextDir)
		nextPos = path.pos.Add(directions[nextDir])
		if nextPos.In(board.space) && board.cells[nextPos.Y][nextPos.X] != '#' {
			if _, ok := visited[nextPos]; !ok {
				visited[nextPos] = struct{}{}
				heap.Push(paths, Path{path.score + 2001, nextDir, nextPos})
			}
		}

		nextDir = nextDirection(nextDir)
		nextPos = path.pos.Add(directions[nextDir])
		if nextPos.In(board.space) && board.cells[nextPos.Y][nextPos.X] != '#' {
			if _, ok := visited[nextPos]; !ok {
				visited[nextPos] = struct{}{}
				heap.Push(paths, Path{path.score + 1001, nextDir, nextPos})
			}
		}
	}

	return 0
}

func getResult(input io.Reader) int64 {
	board := parseInput(input)
	return smallestPath(board)
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
