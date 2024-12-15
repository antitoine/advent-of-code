package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"time"
)

type State struct {
	matrix [][]rune
	space  image.Rectangle
	robot  image.Point
}

func (s *State) moveElt(instruction image.Point, elt image.Point) bool {
	newPosition := elt.Add(instruction)
	if !newPosition.In(s.space) {
		return false
	}
	if s.matrix[newPosition.Y][newPosition.X] == '#' {
		return false
	}
	if s.matrix[newPosition.Y][newPosition.X] == 'O' && !s.moveElt(instruction, newPosition) {
		return false
	}
	s.matrix[newPosition.Y][newPosition.X], s.matrix[elt.Y][elt.X] = s.matrix[elt.Y][elt.X], s.matrix[newPosition.Y][newPosition.X]
	return true
}

func (s *State) Move(instruction image.Point) {
	if s.moveElt(instruction, s.robot) {
		s.robot = s.robot.Add(instruction)
	}
}

func (s *State) GetBoxPositions() []image.Point {
	var boxes []image.Point
	for y, row := range s.matrix {
		for x, char := range row {
			if char == 'O' {
				boxes = append(boxes, image.Pt(x, y))
			}
		}
	}
	return boxes
}

func (s *State) Print() {
	for _, row := range s.matrix {
		fmt.Println(string(row))
	}
}

func parseMap(lines []string) *State {
	var robot image.Point
	var matrix [][]rune

	for y, line := range lines {
		var row []rune
		for x, char := range line {
			row = append(row, char)
			if char == '@' {
				robot = image.Pt(x, y)
			}
		}
		matrix = append(matrix, row)
	}

	return &State{
		matrix: matrix,
		space:  image.Rect(0, 0, len(matrix[0]), len(matrix)),
		robot:  robot,
	}
}

var directions = map[rune]image.Point{
	'^': image.Pt(0, -1),
	'v': image.Pt(0, 1),
	'<': image.Pt(-1, 0),
	'>': image.Pt(1, 0),
}

func parseInstructions(line string) []image.Point {
	var instructions []image.Point
	for _, char := range line {
		if direction, ok := directions[char]; ok {
			instructions = append(instructions, direction)
		} else {
			log.Fatalf("Invalid direction: %c", char)
		}
	}
	return instructions
}

func parseInput(input io.Reader) (*State, []image.Point) {
	scanner := bufio.NewScanner(input)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	state := parseMap(lines)

	var instructions []image.Point
	for scanner.Scan() {
		instructions = append(instructions, parseInstructions(scanner.Text())...)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return state, instructions
}

func getResult(input io.Reader) int64 {
	state, instructions := parseInput(input)
	for _, instruction := range instructions {
		state.Move(instruction)
	}
	boxes := state.GetBoxPositions()
	var result int64
	for _, box := range boxes {
		result += int64(box.X + box.Y*100)
	}
	return result
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
