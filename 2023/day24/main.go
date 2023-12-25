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

func (c Coordinates) Diff(other Coordinates) Coordinates {
	return Coordinates{
		x: c.x - other.x,
		y: c.y - other.y,
		z: c.z - other.z,
	}
}

func (c Coordinates) Div(t int) Coordinates {
	return Coordinates{
		x: c.x / t,
		y: c.y / t,
		z: c.z / t,
	}
}

type Zone struct {
	min Coordinates
	max Coordinates
}

type Trajectory struct {
	position Coordinates
	velocity Coordinates
}

func (t Trajectory) PositionAt(time int) Coordinates {
	return Coordinates{
		x: t.position.x + (t.velocity.x * time),
		y: t.position.y + (t.velocity.y * time),
		z: t.position.z + (t.velocity.z * time),
	}
}

func (t Trajectory) SlopeXY() float64 {
	return float64(t.velocity.y) / float64(t.velocity.x)
}

func (t Trajectory) IsPossibleXY(x float64, y float64) bool {
	if t.velocity.x < 0 && float64(t.position.x) < x {
		return false
	}
	if t.velocity.x > 0 && float64(t.position.x) > x {
		return false
	}
	if t.velocity.y < 0 && float64(t.position.y) < y {
		return false
	}
	if t.velocity.y > 0 && float64(t.position.y) > y {
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

	b1 := ((float64(t.position.y) * float64(t.velocity.x)) - (float64(t.position.x) * float64(t.velocity.y))) / float64(t.velocity.x)
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
	return t.IsPossibleXY(commonX, commonY) && other.IsPossibleXY(commonX, commonY)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (t Trajectory) GetMinimumTime(other Trajectory) int {
	if t.position.x == other.position.x && t.position.y == other.position.y && t.position.z == other.position.z {
		return 0
	}

	if t.velocity.x >= 0 && other.velocity.x <= 0 && t.position.x > other.position.x {
		return -1
	}
	if t.velocity.x == 0 && other.velocity.x == 0 && t.position.x != other.position.x {
		return -1
	}
	if t.velocity.x <= 0 && other.velocity.x >= 0 && t.position.x < other.position.x {
		return -1
	}
	if t.velocity.y >= 0 && other.velocity.y <= 0 && t.position.y > other.position.y {
		return -1
	}
	if t.velocity.y == 0 && other.velocity.y == 0 && t.position.y != other.position.y {
		return -1
	}
	if t.velocity.y <= 0 && other.velocity.y >= 0 && t.position.y < other.position.y {
		return -1
	}
	if t.velocity.z >= 0 && other.velocity.z <= 0 && t.position.z > other.position.z {
		return -1
	}
	if t.velocity.z == 0 && other.velocity.z == 0 && t.position.z != other.position.z {
		return -1
	}
	if t.velocity.z <= 0 && other.velocity.z >= 0 && t.position.z < other.position.z {
		return -1
	}

	minimumTime := -1
	if t.position.x < other.position.x {
		timeNeeded := (other.position.x - t.position.x) / max(1, abs(other.velocity.x-t.velocity.x))
		minimumTime = timeNeeded
	} else if t.position.x > other.position.x {
		timeNeeded := (t.position.x - other.position.x) / max(1, abs(t.velocity.x-other.velocity.x))
		if minimumTime == -1 || timeNeeded < minimumTime {
			minimumTime = timeNeeded
		}
	}
	if t.position.y < other.position.y {
		timeNeeded := (other.position.y - t.position.y) / max(1, abs(other.velocity.y-t.velocity.y))
		if minimumTime == -1 || timeNeeded < minimumTime {
			minimumTime = timeNeeded
		}
	} else if t.position.y > other.position.y {
		timeNeeded := (t.position.y - other.position.y) / max(1, abs(t.velocity.y-other.velocity.y))
		if minimumTime == -1 || timeNeeded < minimumTime {
			minimumTime = timeNeeded
		}
	}
	if t.position.z < other.position.z {
		timeNeeded := (other.position.z - t.position.z) / max(1, abs(other.velocity.z-t.velocity.z))
		if minimumTime == -1 || timeNeeded < minimumTime {
			minimumTime = timeNeeded
		}
	} else if t.position.z > other.position.z {
		timeNeeded := (t.position.z - other.position.z) / max(1, abs(t.velocity.z-other.velocity.z))
		if minimumTime == -1 || timeNeeded < minimumTime {
			minimumTime = timeNeeded
		}
	}

	return minimumTime
}

func (t Trajectory) CollidingWith(other Trajectory) (bool, int) {
	collidingAtT := 0
	currentA := Trajectory{t.position, t.velocity}
	currentB := Trajectory{other.position, other.velocity}
	minimumTime := currentA.GetMinimumTime(currentB)
	for ; minimumTime >= 1; minimumTime = currentA.GetMinimumTime(currentB) {
		collidingAtT += minimumTime
		currentA = Trajectory{currentA.PositionAt(collidingAtT), currentA.velocity}
		currentB = Trajectory{currentB.PositionAt(collidingAtT), currentB.velocity}
	}

	if currentA.position.x == currentB.position.x && currentA.position.y == currentB.position.y && currentA.position.z == currentB.position.z {
		return true, collidingAtT
	}

	return false, 0
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

func GetResultPart2(input io.Reader) int {
	hailstones := parseInput(input)

	minTime := 1
	maxTime := 100
	progress := 0
	totalProgress := len(hailstones) * len(hailstones) * (maxTime - minTime) * (maxTime - minTime - 1) / 2
	for t1 := minTime; t1 < maxTime; t1++ {
		for t2 := t1 + 1; t2 <= maxTime; t2++ {
			for i := 0; i < len(hailstones); i++ {
				for j := 0; j < len(hailstones); j++ {
					progress++
					if progress%100000 == 0 {
						log.Printf("Progress: %d%% (t1=%d / t2=%d)", (progress*100)/totalProgress, t1, t2)
					}

					if i == j {
						continue
					}
					positionAtT1 := hailstones[i].PositionAt(t1)
					positionAtT2 := hailstones[j].PositionAt(t2)

					velocity := positionAtT2.Diff(positionAtT1).Div(t2 - t1)
					rock := Trajectory{
						position: positionAtT1.Diff(velocity),
						velocity: velocity,
					}
					allColliding := true
					for k := 0; k < len(hailstones); k++ {
						if k == i || k == j {
							continue
						}
						colliding, _ := rock.CollidingWith(hailstones[k])
						if !colliding {
							allColliding = false
							break
						}
					}
					if allColliding {
						return rock.position.x + rock.position.y + rock.position.z
					}
				}
			}
		}
	}

	log.Fatalf("Unable to find a solution")

	return 0
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

	result := GetResultPart2(inputFile) // 1025019997186820

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
