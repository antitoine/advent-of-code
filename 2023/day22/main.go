package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Position struct {
	x int
	y int
	z int
}

type Brick struct {
	id   rune
	from Position
	to   Position
}

func (b *Brick) CountCubes() int {
	return (b.to.x - b.from.x + 1) * (b.to.y - b.from.y + 1) * (b.to.z - b.from.z + 1)
}

func (b *Brick) Copy() *Brick {
	return &Brick{
		id:   b.id,
		from: b.from,
		to:   b.to,
	}
}

type Plan struct {
	maxX      int
	maxY      int
	maxZ      int
	positions map[int]map[int]map[int]*Brick
	bricks    Bricks
}

func NewPlan() *Plan {
	return &Plan{
		maxX:      0,
		maxY:      0,
		maxZ:      0,
		positions: make(map[int]map[int]map[int]*Brick),
		bricks:    make(Bricks, 0),
	}
}

func (p *Plan) CopyWithout(removedBrick *Brick) *Plan {
	newPlan := NewPlan()
	for _, brick := range p.bricks {
		if brick == removedBrick {
			continue
		}
		newPlan.Add(brick.Copy())
	}
	return newPlan
}

func (p *Plan) Add(brick *Brick) {
	p.bricks = append(p.bricks, brick)
	if brick.to.x > p.maxX {
		p.maxX = brick.to.x
	}
	if brick.to.y > p.maxY {
		p.maxY = brick.to.y
	}
	if brick.to.z > p.maxZ {
		p.maxZ = brick.to.z
	}
	for x := brick.from.x; x <= brick.to.x; x++ {
		if p.positions[x] == nil {
			p.positions[x] = make(map[int]map[int]*Brick)
		}
		for y := brick.from.y; y <= brick.to.y; y++ {
			if p.positions[x][y] == nil {
				p.positions[x][y] = make(map[int]*Brick)
			}
			for z := brick.from.z; z <= brick.to.z; z++ {
				if p.positions[x][y][z] != nil {
					log.Fatalf("Brick with ID %s cannot be added to the plan, another brick is already there at position  %d,%d,%d", string(brick.id), x, y, z)
				}
				p.positions[x][y][z] = brick
			}
		}
	}
}

func (p *Plan) GetBrickAt(position Position) *Brick {
	if p.positions == nil {
		return nil
	}
	if p.positions[position.x] == nil {
		return nil
	}
	if p.positions[position.x][position.y] == nil {
		return nil
	}
	return p.positions[position.x][position.y][position.z]
}

func (p *Plan) PlanXZ() string {
	var result strings.Builder
	for z := p.maxZ; z > 0; z-- {
		for x := 0; x <= p.maxX; x++ {
			var firsBrickVisible *Brick
			for y := 0; y <= p.maxY && firsBrickVisible == nil; y++ {
				if brick := p.GetBrickAt(Position{x: x, y: y, z: z}); brick != nil {
					firsBrickVisible = brick
				}
			}
			if firsBrickVisible == nil {
				result.WriteRune('.')
			} else {
				result.WriteRune(firsBrickVisible.id)
			}
		}
		result.WriteString("\n")
	}
	for x := 0; x <= p.maxX; x++ {
		result.WriteString("-")
	}
	return result.String()
}

func (p *Plan) PlanYZ() string {
	var result strings.Builder
	for z := p.maxZ; z > 0; z-- {
		for y := 0; y <= p.maxY; y++ {
			var firsBrickVisible *Brick
			for x := 0; x <= p.maxX && firsBrickVisible == nil; x++ {
				if brick := p.GetBrickAt(Position{x: x, y: y, z: z}); brick != nil {
					firsBrickVisible = brick
				}
			}
			if firsBrickVisible == nil {
				result.WriteRune('.')
			} else {
				result.WriteRune(firsBrickVisible.id)
			}
		}
		result.WriteString("\n")
	}
	for y := 0; y <= p.maxY; y++ {
		result.WriteString("-")
	}
	return result.String()
}

func (p *Plan) Bricks() Bricks {
	return p.bricks
}

type Bricks []*Brick

func (b Bricks) SortOnZ() Bricks {
	newBricks := make(Bricks, len(b))
	copy(newBricks, b)

	// Sorts brick on Z axis
	slices.SortStableFunc(newBricks, func(a, b *Brick) int {
		if a.from.z < b.from.z {
			return -1
		}
		if a.from.z > b.from.z {
			return 1
		}
		return 0
	})

	return newBricks
}

// Stabilize make all bricks fallen at the lowest position (on Z) regarding other bricks
func (p *Plan) Stabilize() (*Plan, int) {
	lowestZ := make([][]int, p.maxX+1)
	for x := 0; x <= p.maxX; x++ {
		lowestZ[x] = make([]int, p.maxY+1)
		for y := 0; y <= p.maxY; y++ {
			lowestZ[x][y] = 1
		}
	}

	countMovedBlocks := 0
	newPlan := NewPlan()
	for _, brick := range p.bricks.SortOnZ() {
		newBrick := brick.Copy()
		minZ := 1
		for x := brick.from.x; x <= brick.to.x; x++ {
			for y := brick.from.y; y <= brick.to.y; y++ {
				if lowestZ[x][y] > minZ {
					minZ = lowestZ[x][y]
				}
			}
		}
		diff := brick.from.z - minZ
		newBrick.from.z -= diff
		newBrick.to.z -= diff
		for x := brick.from.x; x <= brick.to.x; x++ {
			for y := brick.from.y; y <= brick.to.y; y++ {
				lowestZ[x][y] = newBrick.to.z + 1
			}
		}
		if diff > 0 {
			countMovedBlocks++
		}
		newPlan.Add(newBrick)
	}

	return newPlan, countMovedBlocks
}

func (p *Plan) ComputeBricksLinks() (map[*Brick]Bricks, map[*Brick]Bricks) {
	brickSupportsBricks := make(map[*Brick]Bricks)
	brickSupportedByBricks := make(map[*Brick]Bricks)
	for _, brick := range p.bricks {
		uniqueSupportsBricks := make(map[*Brick]bool)
		uniqueSupportedByBricks := make(map[*Brick]bool)
		for x := brick.from.x; x <= brick.to.x; x++ {
			for y := brick.from.y; y <= brick.to.y; y++ {
				if supportedBrick := p.GetBrickAt(Position{x: x, y: y, z: brick.to.z + 1}); supportedBrick != nil {
					uniqueSupportsBricks[supportedBrick] = true
				}
				if onBrick := p.GetBrickAt(Position{x: x, y: y, z: brick.from.z - 1}); onBrick != nil {
					uniqueSupportedByBricks[onBrick] = true
				}
			}
		}

		var supportsBricks Bricks
		for supportedBrick := range uniqueSupportsBricks {
			supportsBricks = append(supportsBricks, supportedBrick)
		}
		brickSupportsBricks[brick] = supportsBricks

		var supportedByBricks Bricks
		for onBrick := range uniqueSupportedByBricks {
			supportedByBricks = append(supportedByBricks, onBrick)
		}
		brickSupportedByBricks[brick] = supportedByBricks
	}
	return brickSupportsBricks, brickSupportedByBricks
}

func (p *Plan) CountPossibleBricksToDisintegrate() int {
	brickSupportsBricks, brickSupportedByBricks := p.ComputeBricksLinks()
	var count int
	for _, supportsBricks := range brickSupportsBricks {
		if len(supportsBricks) == 0 {
			count++
			continue
		}
		couldBeDisintegrated := true
		for _, supportedBrick := range supportsBricks {
			if brickOnBricks := brickSupportedByBricks[supportedBrick]; len(brickOnBricks) <= 1 {
				couldBeDisintegrated = false
				break
			}
		}
		if couldBeDisintegrated {
			count++
		}
	}

	return count
}

func (p *Plan) CountBricksFallIfDisintegrated() int {
	totalBricksFall := 0
	for _, brick := range p.bricks {
		newPlan := p.CopyWithout(brick)
		_, movedBricks := newPlan.Stabilize()
		totalBricksFall += movedBricks
	}
	return totalBricksFall
}

func parseCoordinates(coordinates string) Position {
	coords := strings.Split(coordinates, ",")
	if len(coords) != 3 {
		log.Fatalf("Invalid coordinates: %s", coordinates)
	}
	x, errParsingX := strconv.Atoi(coords[0])
	if errParsingX != nil {
		log.Fatalf("Invalid X coordinate: %s", coords[0])
	}
	y, errParsingY := strconv.Atoi(coords[1])
	if errParsingY != nil {
		log.Fatalf("Invalid Y coordinate: %s", coords[1])
	}
	z, errParsingZ := strconv.Atoi(coords[2])
	if errParsingZ != nil {
		log.Fatalf("Invalid Z coordinate: %s", coords[2])
	}
	return Position{
		x: x,
		y: y,
		z: z,
	}
}

func parseBrick(line string) *Brick {
	coords := strings.Split(line, "~")
	if len(coords) != 2 {
		log.Fatalf("Invalid brick line: %s", line)
	}
	from := parseCoordinates(coords[0])
	to := parseCoordinates(coords[1])
	if from.x > to.x || from.y > to.y || from.z > to.z {
		log.Fatalf("Invalid brick line: %s", line)
	}
	return &Brick{
		from: parseCoordinates(coords[0]),
		to:   parseCoordinates(coords[1]),
	}
}

const allowedBrickChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func parseInput(input io.Reader) *Plan {
	scanner := bufio.NewScanner(input)

	plan := NewPlan()
	brickId := 0
	for scanner.Scan() {
		brick := parseBrick(scanner.Text())
		brick.id = rune(allowedBrickChars[brickId])
		brickId++
		if brickId >= len(allowedBrickChars) {
			brickId = 0
		}
		plan.Add(brick)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return plan
}

func getResultPart1(input io.Reader) int {
	plan := parseInput(input)
	stabilizedPlan, _ := plan.Stabilize()
	return stabilizedPlan.CountPossibleBricksToDisintegrate()
}

func getResultPart2(input io.Reader) int {
	plan := parseInput(input)
	stabilizedPlan, _ := plan.Stabilize()
	return stabilizedPlan.CountBricksFallIfDisintegrated()
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

	result := getResultPart2(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
