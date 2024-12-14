package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Robot struct {
	position image.Point
	velocity image.Point
}

var lineRegex = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

func parseLine(line string) Robot {
	matches := lineRegex.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("Unable to parse line: %s", line)
	}

	x, errParsingX := strconv.Atoi(matches[1])
	if errParsingX != nil {
		log.Fatalf("Unable to parse X from line: %s", errParsingX)
	}

	y, errParsingY := strconv.Atoi(matches[2])
	if errParsingY != nil {
		log.Fatalf("Unable to parse Y from line: %s", errParsingY)
	}

	vx, errParsingVX := strconv.Atoi(matches[3])
	if errParsingVX != nil {
		log.Fatalf("Unable to parse VX from line: %s", errParsingVX)
	}

	vy, errParsingVY := strconv.Atoi(matches[4])
	if errParsingVY != nil {
		log.Fatalf("Unable to parse VY from line: %s", errParsingVY)
	}

	return Robot{
		position: image.Pt(x, y),
		velocity: image.Pt(vx, vy),
	}
}

func parseInput(input io.Reader) []Robot {
	scanner := bufio.NewScanner(input)

	var robots []Robot
	for scanner.Scan() {
		robots = append(robots, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return robots
}

func getResult(input io.Reader, sizeX, sizeY, nbSeconds int) int64 {
	robots := parseInput(input)
	quadrants := [4]image.Rectangle{
		image.Rect(0, 0, sizeX/2, sizeY/2),
		image.Rect((sizeX/2)+1, 0, sizeX, sizeY/2),
		image.Rect(0, (sizeY/2)+1, sizeX/2, sizeY),
		image.Rect((sizeX/2)+1, (sizeY/2)+1, sizeX, sizeY),
	}
	nbRobotsInQuadrants := [4]int64{0, 0, 0, 0}
	for _, robot := range robots {
		positionAfterSeconds := robot.position.Add(robot.velocity.Mul(nbSeconds)).Mod(image.Rect(0, 0, sizeX, sizeY))
		for i, quadrant := range quadrants {
			if positionAfterSeconds.In(quadrant) {
				nbRobotsInQuadrants[i]++
				break
			}
		}
	}
	result := int64(1)
	for _, nbRobots := range nbRobotsInQuadrants {
		result *= nbRobots
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

	result := getResult(inputFile, 101, 103, 100)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
