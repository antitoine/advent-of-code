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
	x float64
	y float64
	z float64
}

func (c Coordinates) Add(other Coordinates) Coordinates {
	return Coordinates{
		x: c.x + other.x,
		y: c.y + other.y,
		z: c.z + other.z,
	}
}

func (c Coordinates) Multiply(scalar float64) Coordinates {
	return Coordinates{
		x: c.x * scalar,
		y: c.y * scalar,
		z: c.z * scalar,
	}
}

func (c Coordinates) Subtract(other Coordinates) Coordinates {
	return Coordinates{
		x: c.x - other.x,
		y: c.y - other.y,
		z: c.z - other.z,
	}
}

func (c Coordinates) Divide(scalar float64) Coordinates {
	return Coordinates{
		x: c.x / scalar,
		y: c.y / scalar,
		z: c.z / scalar,
	}
}

func (c Coordinates) CrossProduct(other Coordinates) Coordinates {
	return Coordinates{
		x: c.y*other.z - c.z*other.y,
		y: c.z*other.x - c.x*other.z,
		z: c.x*other.y - c.y*other.x,
	}
}

func (c Coordinates) DotProduct(other Coordinates) float64 {
	return c.x*other.x + c.y*other.y + c.z*other.z
}

type Zone struct {
	min Coordinates
	max Coordinates
}

type Trajectory struct {
	position Coordinates
	velocity Coordinates
}

func (t Trajectory) SlopeXY() float64 {
	return t.velocity.y / t.velocity.x
}

func (t Trajectory) IsPossibleXY(x float64, y float64) bool {
	if t.velocity.x < 0 && t.position.x < x {
		return false
	}
	if t.velocity.x > 0 && t.position.x > x {
		return false
	}
	if t.velocity.y < 0 && t.position.y < y {
		return false
	}
	if t.velocity.y > 0 && t.position.y > y {
		return false
	}
	return true
}

func (t Trajectory) IsXYCollidingWith(other Trajectory, in Zone) bool {
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

	a1 := t.SlopeXY()
	a2 := other.SlopeXY()
	if a1 == a2 {
		return false // parallel lines
	}

	b1 := ((t.position.y * t.velocity.x) - (t.position.x * t.velocity.y)) / t.velocity.x
	b2 := ((other.position.y * other.velocity.x) - (other.position.x * other.velocity.y)) / other.velocity.x

	// Common X of two lines:
	// y = a1 x + b1
	// y = a2 x + b2
	// a1 x + b1 = a2 x + b2
	// a1 x - a2 x = b2 - b1
	// x (a1 - a2) = b2 - b1
	// x = (b2 - b1) / (a1 - a2)

	commonX := (b2 - b1) / (a1 - a2)
	if commonX < in.min.x || commonX > in.max.x {
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
	if commonY < in.min.y || commonY > in.max.y {
		return false
	}

	// Check if the common point is in both lines:
	return t.IsPossibleXY(commonX, commonY) && other.IsPossibleXY(commonX, commonY)
}

func parseCoordinates(coordinates string) Coordinates {
	parts := strings.Split(coordinates, ",")
	if len(parts) != 3 {
		log.Fatalf("Unable to parse coordinates: %s", coordinates)
	}
	xStr := strings.TrimSpace(parts[0])
	x, errParsingX := strconv.ParseFloat(xStr, 64)
	if errParsingX != nil {
		log.Fatalf("Unable to parse x coordinate: %s", xStr)
	}
	yStr := strings.TrimSpace(parts[1])
	y, errParsingY := strconv.ParseFloat(yStr, 64)
	if errParsingY != nil {
		log.Fatalf("Unable to parse y coordinate: %s", yStr)
	}
	zStr := strings.TrimSpace(parts[2])
	z, errParsingZ := strconv.ParseFloat(zStr, 64)
	if errParsingZ != nil {
		log.Fatalf("Unable to parse z coordinate: %s", zStr)
	}
	return Coordinates{
		x: x,
		y: y,
		z: z,
	}
}

func parseLine(line string) Trajectory {
	coordinates := strings.Split(line, "@")
	if len(coordinates) != 2 {
		log.Fatalf("Unable to parse line: %s", line)
	}
	return Trajectory{
		position: parseCoordinates(coordinates[0]),
		velocity: parseCoordinates(coordinates[1]),
	}
}

func parseInput(input io.Reader) []Trajectory {
	scanner := bufio.NewScanner(input)

	var hailstones []Trajectory
	for scanner.Scan() {
		hailstones = append(hailstones, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return hailstones
}

func GetResultPart1(input io.Reader, testZone Zone) int64 {
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

func planeLineIntersection(p0, n, p, v Coordinates) (Coordinates, float64) {
	denominator := v.DotProduct(n)
	numerator := p0.Subtract(p).DotProduct(n)
	t := numerator / denominator
	intersection := p.Add(v.Multiply(t))
	return intersection, t
}

func reduceHailstonesPositions(hailstones []Trajectory) []Trajectory {
	result := make([]Trajectory, len(hailstones))
	for i := 0; i < len(hailstones); i++ {
		result[i] = Trajectory{
			position: Coordinates{
				x: hailstones[i].position.x / 1e12,
				y: hailstones[i].position.y / 1e12,
				z: hailstones[i].position.z / 1e12,
			},
			velocity: hailstones[i].velocity,
		}
	}
	return result
}

// https://en.wikipedia.org/wiki/Line%E2%80%93plane_intersection
func GetResultPart2(input io.Reader) int64 {
	originalHailstones := parseInput(input)
	hailstones := reduceHailstonesPositions(originalHailstones)

	p0, v0 := hailstones[0].position, hailstones[0].velocity
	p1, v1 := hailstones[1].position.Subtract(p0), hailstones[1].velocity.Subtract(v0)
	p2, v2 := hailstones[2].position.Subtract(p0), hailstones[2].velocity.Subtract(v0)
	p3, v3 := hailstones[3].position.Subtract(p0), hailstones[3].velocity.Subtract(v0)

	rockPlaneN := p1.CrossProduct(p1.Add(v1))

	origin := Coordinates{x: 0, y: 0, z: 0}
	p02, t02 := planeLineIntersection(origin, rockPlaneN, p2, v2)
	p03, t03 := planeLineIntersection(origin, rockPlaneN, p3, v3)

	vRock := p02.Subtract(p03).Divide(t02 - t03)
	pRock := p02.Subtract(vRock.Multiply(t02))

	pRock = pRock.Add(p0)
	vRock = vRock.Add(v0)

	rock := Trajectory{position: pRock, velocity: vRock}

	return int64(rock.position.x*1e12 + rock.position.y*1e12 + rock.position.z*1e12)
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

	result := GetResultPart2(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
