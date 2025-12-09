package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

func parseInput(input io.Reader) []Point {
	var points []Point
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		points = append(points, Point{x, y})
	}
	return points
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// isPointInsidePolygon uses ray casting to check if a point is inside the polygon
func isPointInsidePolygon(px, py float64, redTiles []Point) bool {
	n := len(redTiles)
	count := 0
	for i := 0; i < n; i++ {
		p1 := redTiles[i]
		p2 := redTiles[(i+1)%n]

		// Only count vertical edges (horizontal ray crosses vertical edges)
		if p1.x == p2.x {
			x := float64(p1.x)
			y1, y2 := float64(p1.y), float64(p2.y)
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			// Check if ray from (px, py) going right crosses this edge
			if px < x && py > y1 && py < y2 {
				count++
			}
		}
	}
	return count%2 == 1
}

func getResult(input io.Reader) int64 {
	redTiles := parseInput(input)

	// Coordinate compression: extract unique x and y coordinates
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range redTiles {
		xSet[p.x] = true
		ySet[p.y] = true
	}

	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)

	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)

	// Create index maps
	xIndex := make(map[int]int)
	for i, x := range xCoords {
		xIndex[x] = i
	}
	yIndex := make(map[int]int)
	for i, y := range yCoords {
		yIndex[y] = i
	}

	nx, ny := len(xCoords), len(yCoords)

	// For each cell in the compressed grid, determine if it's outside the polygon
	// Cell (i, j) represents the region (xCoords[i], xCoords[i+1]) x (yCoords[j], yCoords[j+1])
	outside := make([][]bool, nx-1)
	for i := 0; i < nx-1; i++ {
		outside[i] = make([]bool, ny-1)
		for j := 0; j < ny-1; j++ {
			// Pick center of cell
			cx := float64(xCoords[i]+xCoords[i+1]) / 2.0
			cy := float64(yCoords[j]+yCoords[j+1]) / 2.0
			outside[i][j] = !isPointInsidePolygon(cx, cy, redTiles)
		}
	}

	// Build 2D prefix sum of outside cells
	// prefix[i][j] = sum of outside cells in [0, i) x [0, j)
	prefix := make([][]int, nx)
	for i := 0; i <= nx-1; i++ {
		prefix[i] = make([]int, ny)
	}
	for i := 1; i < nx; i++ {
		for j := 1; j < ny; j++ {
			val := 0
			if outside[i-1][j-1] {
				val = 1
			}
			prefix[i][j] = prefix[i-1][j] + prefix[i][j-1] - prefix[i-1][j-1] + val
		}
	}

	// Helper to check if rectangle in compressed grid contains any outside cells
	hasOutsideCells := func(ix1, ix2, iy1, iy2 int) bool {
		if ix1 >= ix2 || iy1 >= iy2 {
			return false // No cells in degenerate rectangle
		}
		sum := prefix[ix2][iy2] - prefix[ix1][iy2] - prefix[ix2][iy1] + prefix[ix1][iy1]
		return sum > 0
	}

	// Helper to check if a horizontal segment is inside (for degenerate rectangles)
	isHorizontalSegmentInside := func(x1, x2, y int) bool {
		ix1, ix2 := xIndex[x1], xIndex[x2]
		if ix1 > ix2 {
			ix1, ix2 = ix2, ix1
		}
		// Check each sub-segment between consecutive x coordinates
		for i := ix1; i < ix2; i++ {
			cx := float64(xCoords[i]+xCoords[i+1]) / 2.0
			if !isPointInsidePolygon(cx, float64(y), redTiles) {
				return false
			}
		}
		return true
	}

	// Helper to check if a vertical segment is inside (for degenerate rectangles)
	isVerticalSegmentInside := func(x, y1, y2 int) bool {
		iy1, iy2 := yIndex[y1], yIndex[y2]
		if iy1 > iy2 {
			iy1, iy2 = iy2, iy1
		}
		// Check each sub-segment between consecutive y coordinates
		for j := iy1; j < iy2; j++ {
			cy := float64(yCoords[j]+yCoords[j+1]) / 2.0
			if !isPointInsidePolygon(float64(x), cy, redTiles) {
				return false
			}
		}
		return true
	}

	// For each pair of red tiles, check if all tiles in rectangle are valid
	maxArea := int64(0)
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			p1, p2 := redTiles[i], redTiles[j]
			x1, x2 := minInt(p1.x, p2.x), maxInt(p1.x, p2.x)
			y1, y2 := minInt(p1.y, p2.y), maxInt(p1.y, p2.y)

			ix1, ix2 := xIndex[x1], xIndex[x2]
			iy1, iy2 := yIndex[y1], yIndex[y2]

			valid := true
			if x1 == x2 {
				// Vertical line
				valid = isVerticalSegmentInside(x1, y1, y2)
			} else if y1 == y2 {
				// Horizontal line
				valid = isHorizontalSegmentInside(x1, x2, y1)
			} else {
				// Normal rectangle: check cells using prefix sum
				valid = !hasOutsideCells(ix1, ix2, iy1, iy2)
			}

			if valid {
				width := x2 - x1 + 1
				height := y2 - y1 + 1
				area := int64(width) * int64(height)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea
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
