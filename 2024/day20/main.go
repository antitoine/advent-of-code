package main

import (
	"bufio"
	"container/heap"
	"fmt"
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

var ZeroPoint = image.Pt(-1, -1)

type Cheats struct {
	first  image.Point
	second image.Point
}

type Path struct {
	pos         image.Point
	steps       int
	chests      Cheats
	previousPos map[image.Point]int
}

func (p Path) Next(pos image.Point, withCheat bool) Path {
	newCheats := p.chests
	if withCheat {
		if newCheats.first == ZeroPoint {
			newCheats.first = pos
		} else if newCheats.second == ZeroPoint {
			newCheats.second = pos
		} else {
			log.Fatalf("More than two cheats found")
		}
	}
	return Path{
		pos:    pos,
		steps:  p.steps + 1,
		chests: newCheats,
	}
}

func (p Path) SaveAndNext(pos image.Point, withCheat bool) Path {
	next := p.Next(pos, withCheat)
	newPreviousPos := make(map[image.Point]int)
	for k, v := range p.previousPos {
		newPreviousPos[k] = v
	}
	newPreviousPos[p.pos] = p.steps
	next.previousPos = newPreviousPos
	return next
}

func InitPath(pos image.Point) Path {
	return Path{
		pos:    pos,
		steps:  0,
		chests: Cheats{ZeroPoint, ZeroPoint},
	}
}

type SmallestPathHeap []Path

func (h SmallestPathHeap) Len() int           { return len(h) }
func (h SmallestPathHeap) Less(i, j int) bool { return h[i].steps < h[j].steps }
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

func getSmallestPathWithoutCheats(track Track) *Path {
	paths := &SmallestPathHeap{}
	heap.Push(paths, InitPath(track.start))
	visited := make(map[image.Point]struct{})

	for paths.Len() > 0 {
		path := heap.Pop(paths).(Path)

		if path.pos == track.end {
			return &path
		}

		if _, ok := visited[path.pos]; ok {
			continue
		}
		visited[path.pos] = struct{}{}

		for _, dirDelta := range directions {
			nextPos := path.pos.Add(dirDelta)
			if !nextPos.In(track.space) {
				continue
			}
			if track.track[nextPos.Y][nextPos.X] != '#' {
				heap.Push(paths, path.SaveAndNext(nextPos, false))
			}
		}
	}

	return nil
}

type VisitedKey struct {
	pos    image.Point
	cheats Cheats
}

func getAllPathsWithCheats(track Track, withoutCheatPath Path) []Path {
	paths := &SmallestPathHeap{}
	heap.Push(paths, InitPath(track.start))
	var possiblesPaths []Path
	visited := make(map[VisitedKey]struct{})

	for paths.Len() > 0 {
		path := heap.Pop(paths).(Path)

		if path.steps >= withoutCheatPath.steps {
			continue
		}

		if track.track[path.pos.Y][path.pos.X] == '#' && path.chests.second != ZeroPoint {
			continue
		}

		if path.pos == track.end {
			possiblesPaths = append(possiblesPaths, path)
			continue
		}

		if path.chests.second != ZeroPoint {
			if prevPathSteps, ok := withoutCheatPath.previousPos[path.pos]; ok && path.steps >= prevPathSteps {
				continue
			}
		}

		if _, ok := visited[VisitedKey{path.pos, path.chests}]; ok {
			continue
		}
		visited[VisitedKey{path.pos, path.chests}] = struct{}{}

		for _, dirDelta := range directions {
			nextPos := path.pos.Add(dirDelta)
			if !nextPos.In(track.space) {
				continue
			}
			if path.chests.first != ZeroPoint && path.chests.second == ZeroPoint {
				heap.Push(paths, path.Next(nextPos, true))
			} else if track.track[nextPos.Y][nextPos.X] == '#' && path.chests.first == ZeroPoint {
				heap.Push(paths, path.Next(nextPos, true))
			} else if track.track[nextPos.Y][nextPos.X] != '#' {
				heap.Push(paths, path.Next(nextPos, false))
			}
		}
	}

	return possiblesPaths
}

func getResult(input io.Reader, nbLeastSavingSteps int) int64 {
	track := parseInput(input)

	smallestPathWithoutCheats := getSmallestPathWithoutCheats(track)
	if smallestPathWithoutCheats == nil {
		log.Fatalf("No path found")
	}
	fmt.Printf("Smallest path without cheats: %d\n", smallestPathWithoutCheats.steps)

	pathsWithCheats := getAllPathsWithCheats(track, *smallestPathWithoutCheats)

	var nbCheatsSavingSteps int64
	for _, path := range pathsWithCheats {
		savingSteps := smallestPathWithoutCheats.steps - path.steps
		if savingSteps >= nbLeastSavingSteps {
			nbCheatsSavingSteps++
		}
	}

	// Print the total number of cheats grouped by the amount of time they save
	//groups := make(map[int]map[Cheats]struct{})
	//for _, path := range pathsWithCheats {
	//	savingSteps := smallestPathWithoutCheats.steps - path.steps
	//	if _, ok := groups[savingSteps]; !ok {
	//		groups[savingSteps] = make(map[Cheats]struct{})
	//	}
	//	groups[savingSteps][path.chests] = struct{}{}
	//}
	//
	//for saving := 0; saving <= smallestPathWithoutCheats.steps; saving++ {
	//	if cheats, ok := groups[saving]; ok {
	//		log.Printf("There are %d cheats that save %d picoseconds", len(cheats), saving)
	//	}
	//}

	return nbCheatsSavingSteps
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

	result := getResult(inputFile, 100)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
