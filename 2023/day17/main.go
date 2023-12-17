package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Grid [][]uint8

func (g Grid) IsAllowed(position Position) bool {
	return position.i >= 0 && position.i < len(g) && position.j >= 0 && position.j < len(g[0])
}

func (g Grid) GetHeatLoss(position Position) uint8 {
	return g[position.i][position.j]
}

func (g Grid) String() string {
	var result string
	for _, row := range g {
		for _, heatLoss := range row {
			result += fmt.Sprintf("%d", heatLoss)
		}
		result += "\n"
	}
	return result
}

type Direction int

const (
	North Direction = iota
	South Direction = iota
	East  Direction = iota
	West  Direction = iota
)

func (d Direction) TurnLeft() Direction {
	switch d {
	case North:
		return West
	case South:
		return East
	case East:
		return North
	case West:
		return South
	}
	log.Fatalf("Unknown direction: %d", d)
	return d
}

func (d Direction) TurnRight() Direction {
	switch d {
	case North:
		return East
	case South:
		return West
	case East:
		return South
	case West:
		return North
	}
	log.Fatalf("Unknown direction: %d", d)
	return d
}

type Position struct {
	i, j int
}

func (p Position) Move(direction Direction) Position {
	switch direction {
	case North:
		return Position{p.i - 1, p.j}
	case South:
		return Position{p.i + 1, p.j}
	case East:
		return Position{p.i, p.j + 1}
	case West:
		return Position{p.i, p.j - 1}
	}
	log.Fatalf("Unknown direction: %d", direction)
	return p
}

type Step struct {
	position              Position
	heatLoss              uint64
	direction             Direction
	currentStraightLength int
	queueIndex            int
}

func (s *Step) Move(grid Grid, direction Direction, resetStraightLength bool) (*Step, bool) {
	newPosition := s.position.Move(direction)
	if grid.IsAllowed(newPosition) {
		step := &Step{
			position:  newPosition,
			heatLoss:  s.heatLoss + uint64(grid.GetHeatLoss(newPosition)),
			direction: direction,
		}
		if resetStraightLength {
			step.currentStraightLength = 1
		} else {
			step.currentStraightLength = s.currentStraightLength + 1
		}
		return step, true
	}
	return nil, false
}

type StepQueue []*Step

func (pq StepQueue) Len() int           { return len(pq) }
func (pq StepQueue) Less(i, j int) bool { return pq[i].heatLoss < pq[j].heatLoss }
func (pq StepQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *StepQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Step)
	item.queueIndex = n
	*pq = append(*pq, item)
}

func (pq *StepQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.queueIndex = -1
	*pq = old[0 : n-1]
	return item
}

type Visit struct {
	position              Position
	direction             Direction
	currentStraightLength int
}

type VisitedSet map[Visit]bool

func (v VisitedSet) IsVisited(step *Step) bool {
	visit := Visit{
		position:              step.position,
		direction:             step.direction,
		currentStraightLength: step.currentStraightLength,
	}
	_, exists := v[visit]
	v[visit] = true
	return exists
}

const minStraightDistance = 4
const maxStraightDistance = 10

func getMinimumHeatLoss(grid Grid) uint64 {
	stepQueue := StepQueue{}
	heap.Init(&stepQueue)
	heap.Push(&stepQueue, &Step{position: Position{0, 0}, direction: East, currentStraightLength: 0})
	heap.Push(&stepQueue, &Step{position: Position{0, 0}, direction: South, currentStraightLength: 0})

	visits := make(VisitedSet)

	for stepQueue.Len() > 0 {
		step := heap.Pop(&stepQueue).(*Step)

		if visits.IsVisited(step) {
			continue
		}

		if step.position.i == len(grid)-1 && step.position.j == len(grid[0])-1 {
			if step.currentStraightLength >= minStraightDistance {
				return step.heatLoss
			}
			continue
		}

		// If we could move in the same direction, do it
		if step.currentStraightLength < maxStraightDistance {
			if newStep, allowed := step.Move(grid, step.direction, false); allowed {
				heap.Push(&stepQueue, newStep)
			}
		}

		if step.currentStraightLength >= minStraightDistance {
			// Try turning left
			if newStep, allowed := step.Move(grid, step.direction.TurnLeft(), true); allowed {
				heap.Push(&stepQueue, newStep)
			}

			// Try turning right
			if newStep, allowed := step.Move(grid, step.direction.TurnRight(), true); allowed {
				heap.Push(&stepQueue, newStep)
			}
		}
	}

	log.Fatalf("Unable to find a path")

	return 0
}

func parseInput(input io.Reader) Grid {
	scanner := bufio.NewScanner(input)

	var grid Grid
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]uint8, len(line))
		for column, char := range line {
			heatLoss, errParsing := strconv.ParseUint(string(char), 10, 8)
			if errParsing != nil {
				log.Fatalf("Unable to parse uint8 from character '%s': %v", string(char), errParsing)
			}
			row[column] = uint8(heatLoss)
		}
		grid = append(grid, row)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return grid
}

func getResult(input io.Reader) uint64 {
	grid := parseInput(input)
	//log.Printf("Grid:\n%s", grid)
	return getMinimumHeatLoss(grid)
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
