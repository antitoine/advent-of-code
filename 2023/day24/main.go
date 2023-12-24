package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Coordinates struct {
	x int
	y int
	z int
}

func parseCoordinates(coordinates string) Coordinates {
	parts := strings.Split(coordinates, ",")
	if len(parts) != 3 {
		log.Fatalf("Unable to parse coordinates: %s", coordinates)
	}
	xStr := strings.TrimSpace(parts[0])
	x, errParsingX := strconv.Atoi(xStr)
	if errParsingX != nil {
		log.Fatalf("Unable to parse x coordinate: %s", xStr)
	}
	yStr := strings.TrimSpace(parts[1])
	y, errParsingY := strconv.Atoi(yStr)
	if errParsingY != nil {
		log.Fatalf("Unable to parse y coordinate: %s", yStr)
	}
	zStr := strings.TrimSpace(parts[2])
	z, errParsingZ := strconv.Atoi(zStr)
	if errParsingZ != nil {
		log.Fatalf("Unable to parse z coordinate: %s", zStr)
	}
	return Coordinates{
		x: x,
		y: y,
		z: z,
	}
}

type Zone struct {
	min Coordinates
	max Coordinates
}

type Hailstone struct {
	position Coordinates
	velocity Coordinates
}

func (h Hailstone) SlopeXY() float64 {
	return float64(h.velocity.y) / float64(h.velocity.x)
}

func (h Hailstone) IsPossibleXY(x float64, y float64) bool {
	if h.velocity.x < 0 && float64(h.position.x) < x {
		return false
	}
	if h.velocity.x > 0 && float64(h.position.x) > x {
		return false
	}
	if h.velocity.y < 0 && float64(h.position.y) < y {
		return false
	}
	if h.velocity.y > 0 && float64(h.position.y) > y {
		return false
	}
	return true
}

func (h Hailstone) IsXYCollidingWith(other Hailstone, in Zone) bool {
	// Initial point and vector to line equation y = a x + b:
	// x = x0 + vx t
	// t = (x - x0) / vx
	// y = y0 + vy t
	// t = (y - y0) / vy
	// (x - x0) / vx = (y - y0) / vy
	// (x - x0) vy = (y - y0) vx
	// x vy - x0 vy = y vx - y0 vx
	// x vy - y vx = x0 vy - y0 vx
	// y vx - x vy = y0 vx - x0 vy
	// y vx = x vy + y0 vx - x0 vy
	// y = (x vy + y0 vx - x0 vy) / vx
	// a = vy / vx
	// b = (y0 vx - x0 vy) / vx

	a1 := h.SlopeXY()
	a2 := other.SlopeXY()
	if a1 == a2 {
		return false // parallel lines
	}

	b1 := ((float64(h.position.y) * float64(h.velocity.x)) - (float64(h.position.x) * float64(h.velocity.y))) / float64(h.velocity.x)
	b2 := ((float64(other.position.y) * float64(other.velocity.x)) - (float64(other.position.x) * float64(other.velocity.y))) / float64(other.velocity.x)

	// Common X of two lines:
	// y = a1 x + b1
	// y = a2 x + b2
	// a1 x + b1 = a2 x + b2
	// a1 x - a2 x = b2 - b1
	// x (a1 - a2) = b2 - b1
	// x = (b2 - b1) / (a1 - a2)

	commonX := (b2 - b1) / (a1 - a2)
	if commonX < float64(in.min.x) || commonX > float64(in.max.x) {
		return false
	}

	// Equation x based:
	// y = a x + b
	// a x = y - b
	// x = (y - b) / a

	// Common Y of two lines:
	// x = (y - b1) / a1
	// x = (y - b2) / a2
	// (y - b1) / a1 = (y - b2) / a2
	// (y - b1) a2 = (y - b2) a1
	// y a2 - b1 a2 = y a1 - b2 a1
	// y a2 - y a1 = b1 a2 - b2 a1
	// y (a2 - a1) = b1 a2 - b2 a1
	// y = (b1 a2 - b2 a1) / (a2 - a1)

	commonY := (b1*a2 - b2*a1) / (a2 - a1)
	if commonY < float64(in.min.y) || commonY > float64(in.max.y) {
		return false
	}

	// Check if the common point is in both lines:
	return h.IsPossibleXY(commonX, commonY) && other.IsPossibleXY(commonX, commonY)
}

func parseLine(line string) Hailstone {
	coordinates := strings.Split(line, "@")
	if len(coordinates) != 2 {
		log.Fatalf("Unable to parse line: %s", line)
	}
	return Hailstone{
		position: parseCoordinates(coordinates[0]),
		velocity: parseCoordinates(coordinates[1]),
	}
}

func parseInput(input io.Reader) []Hailstone {
	scanner := bufio.NewScanner(input)

	var hailstones []Hailstone
	for scanner.Scan() {
		hailstones = append(hailstones, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return hailstones
}

func getResult(input io.Reader, testZone Zone) int64 {
	hailstones := parseInput(input)

	var count int64
	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			if hailstones[i].IsXYCollidingWith(hailstones[j], testZone) {
				count++
			}
		}
	}

	return count
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

	testZone := Zone{
		min: Coordinates{
			x: 200000000000000,
			y: 200000000000000,
		},
		max: Coordinates{
			x: 400000000000000,
			y: 400000000000000,
		},
	}
	result := getResult(inputFile, testZone)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
