package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Shape struct {
	index     int
	grid      [][]bool
	cellCount int
}

type Region struct {
	width  int
	height int
	counts []int
}

// Piece represents a specific shape variant with precomputed cells
type Piece struct {
	shapeIdx int
	cells    [][2]int // relative positions of filled cells
	rows     int
	cols     int
}

func parseGrid(lines []string) [][]bool {
	grid := make([][]bool, len(lines))
	for i, line := range lines {
		grid[i] = make([]bool, len(line))
		for j, char := range line {
			grid[i][j] = (char == '#')
		}
	}
	return grid
}

func countCells(grid [][]bool) int {
	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				count++
			}
		}
	}
	return count
}

func rotate90(grid [][]bool) [][]bool {
	rows := len(grid)
	if rows == 0 {
		return grid
	}
	cols := len(grid[0])

	result := make([][]bool, cols)
	for i := range result {
		result[i] = make([]bool, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result[j][rows-1-i] = grid[i][j]
		}
	}

	return result
}

func flipHorizontal(grid [][]bool) [][]bool {
	rows := len(grid)
	if rows == 0 {
		return grid
	}
	cols := len(grid[0])

	result := make([][]bool, rows)
	for i := range result {
		result[i] = make([]bool, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = grid[i][cols-1-j]
		}
	}

	return result
}

func gridKey(grid [][]bool) string {
	var sb strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func gridToCells(grid [][]bool) [][2]int {
	var cells [][2]int
	for i, row := range grid {
		for j, cell := range row {
			if cell {
				cells = append(cells, [2]int{i, j})
			}
		}
	}
	return cells
}

func generatePieces(shapeIdx int, grid [][]bool) []Piece {
	seen := make(map[string]bool)
	var pieces []Piece

	var addVariant func([][]bool)
	addVariant = func(g [][]bool) {
		key := gridKey(g)
		if seen[key] {
			return
		}
		seen[key] = true

		pieces = append(pieces, Piece{
			shapeIdx: shapeIdx,
			cells:    gridToCells(g),
			rows:     len(g),
			cols:     len(g[0]),
		})

		rotated := rotate90(g)
		addVariant(rotated)
		addVariant(rotate90(rotated))
		addVariant(rotate90(rotate90(rotated)))

		flipped := flipHorizontal(g)
		addVariant(flipped)
		addVariant(rotate90(flipped))
		addVariant(rotate90(rotate90(flipped)))
		addVariant(rotate90(rotate90(rotate90(flipped))))
	}

	addVariant(grid)
	return pieces
}

func parseInput(input io.Reader) ([]Shape, []Region, error) {
	scanner := bufio.NewScanner(input)
	var shapes []Shape
	var regions []Region

	var currentShape *Shape
	var currentGrid []string
	var inShapeSection = true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			if currentShape != nil && len(currentGrid) > 0 {
				currentShape.grid = parseGrid(currentGrid)
				currentShape.cellCount = countCells(currentShape.grid)
				shapes = append(shapes, *currentShape)
				currentShape = nil
				currentGrid = nil
			}
			continue
		}

		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			inShapeSection = false
			if currentShape != nil && len(currentGrid) > 0 {
				currentShape.grid = parseGrid(currentGrid)
				currentShape.cellCount = countCells(currentShape.grid)
				shapes = append(shapes, *currentShape)
				currentShape = nil
				currentGrid = nil
			}

			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}

			sizePart := strings.TrimSuffix(parts[0], ":")
			sizeParts := strings.Split(sizePart, "x")
			if len(sizeParts) != 2 {
				continue
			}

			width, err := strconv.Atoi(sizeParts[0])
			if err != nil {
				continue
			}

			height, err := strconv.Atoi(sizeParts[1])
			if err != nil {
				continue
			}

			counts := make([]int, len(parts)-1)
			for i := 1; i < len(parts); i++ {
				count, err := strconv.Atoi(parts[i])
				if err != nil {
					return nil, nil, fmt.Errorf("invalid count: %v", err)
				}
				counts[i-1] = count
			}

			regions = append(regions, Region{
				width:  width,
				height: height,
				counts: counts,
			})
			continue
		}

		if inShapeSection {
			if strings.HasSuffix(line, ":") {
				if currentShape != nil && len(currentGrid) > 0 {
					currentShape.grid = parseGrid(currentGrid)
					currentShape.cellCount = countCells(currentShape.grid)
					shapes = append(shapes, *currentShape)
				}
				indexStr := strings.TrimSuffix(line, ":")
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid shape index: %v", err)
				}
				currentShape = &Shape{index: index}
				currentGrid = nil
			} else if currentShape != nil {
				currentGrid = append(currentGrid, line)
			}
		}
	}

	if currentShape != nil && len(currentGrid) > 0 {
		currentShape.grid = parseGrid(currentGrid)
		currentShape.cellCount = countCells(currentShape.grid)
		shapes = append(shapes, *currentShape)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return shapes, regions, nil
}

// PieceToPlace represents a piece that needs to be placed, with ordering info
type PieceToPlace struct {
	shapeIdx   int
	instanceID int // for ordering identical pieces
}

func canFit(region Region, shapes []Shape, allPieces [][]Piece) bool {
	// Quick area check
	totalCells := 0
	for shapeIdx, count := range region.counts {
		if shapeIdx < len(shapes) {
			totalCells += count * shapes[shapeIdx].cellCount
		}
	}
	if totalCells > region.width*region.height {
		return false
	}

	// Build list of pieces to place
	var piecesToPlace []PieceToPlace
	for shapeIdx, count := range region.counts {
		for i := 0; i < count; i++ {
			piecesToPlace = append(piecesToPlace, PieceToPlace{
				shapeIdx:   shapeIdx,
				instanceID: i,
			})
		}
	}

	if len(piecesToPlace) == 0 {
		return true
	}

	// Sort by number of variants (most constrained first)
	sort.Slice(piecesToPlace, func(i, j int) bool {
		pi, pj := piecesToPlace[i], piecesToPlace[j]
		vi := len(allPieces[pi.shapeIdx])
		vj := len(allPieces[pj.shapeIdx])
		if vi != vj {
			return vi < vj
		}
		if pi.shapeIdx != pj.shapeIdx {
			return pi.shapeIdx < pj.shapeIdx
		}
		return pi.instanceID < pj.instanceID
	})

	// Use a flat grid for speed
	grid := make([]bool, region.height*region.width)

	// Track last position for each shape to avoid duplicate orderings
	lastPos := make(map[int]int)

	return backtrack(piecesToPlace, grid, region.width, region.height, 0, allPieces, lastPos)
}

func backtrack(pieces []PieceToPlace, grid []bool, width, height, pieceIdx int, allPieces [][]Piece, lastPos map[int]int) bool {
	if pieceIdx >= len(pieces) {
		return true
	}

	piece := pieces[pieceIdx]
	variants := allPieces[piece.shapeIdx]

	// Get minimum position for this shape to avoid duplicate orderings
	minPos := 0
	if piece.instanceID > 0 {
		if lp, ok := lastPos[piece.shapeIdx]; ok {
			minPos = lp
		}
	}

	for _, variant := range variants {
		maxRow := height - variant.rows
		maxCol := width - variant.cols

		for row := 0; row <= maxRow; row++ {
			for col := 0; col <= maxCol; col++ {
				pos := row*width + col
				if pos < minPos {
					continue
				}

				if canPlace(variant, grid, width, row, col) {
					place(variant, grid, width, row, col)

					oldPos := lastPos[piece.shapeIdx]
					lastPos[piece.shapeIdx] = pos

					if backtrack(pieces, grid, width, height, pieceIdx+1, allPieces, lastPos) {
						return true
					}

					lastPos[piece.shapeIdx] = oldPos
					unplace(variant, grid, width, row, col)
				}
			}
		}
	}

	return false
}

func canPlace(piece Piece, grid []bool, width, row, col int) bool {
	for _, cell := range piece.cells {
		idx := (row+cell[0])*width + (col + cell[1])
		if grid[idx] {
			return false
		}
	}
	return true
}

func place(piece Piece, grid []bool, width, row, col int) {
	for _, cell := range piece.cells {
		idx := (row+cell[0])*width + (col + cell[1])
		grid[idx] = true
	}
}

func unplace(piece Piece, grid []bool, width, row, col int) {
	for _, cell := range piece.cells {
		idx := (row+cell[0])*width + (col + cell[1])
		grid[idx] = false
	}
}

func getResult(input io.Reader) int64 {
	shapes, regions, err := parseInput(input)
	if err != nil {
		log.Fatalf("Error parsing input: %v", err)
	}

	maxIndex := 0
	for _, shape := range shapes {
		if shape.index > maxIndex {
			maxIndex = shape.index
		}
	}

	shapeMap := make(map[int]Shape)
	for _, shape := range shapes {
		shapeMap[shape.index] = shape
	}

	allPieces := make([][]Piece, maxIndex+1)
	for i := 0; i <= maxIndex; i++ {
		if shape, ok := shapeMap[i]; ok {
			allPieces[i] = generatePieces(i, shape.grid)
		} else {
			allPieces[i] = []Piece{}
		}
	}

	count := 0
	for _, region := range regions {
		if canFit(region, shapes, allPieces) {
			count++
		}
	}

	return int64(count)
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
