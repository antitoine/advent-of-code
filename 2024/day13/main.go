package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var buttonRegex = regexp.MustCompile(`Button \w: X\+(\d+), Y\+(\d+)`)

func parseButton(line string) image.Point {
	matches := buttonRegex.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("Unable to parse button from line: %s", line)
	}

	x, errParsingX := strconv.Atoi(matches[1])
	if errParsingX != nil {
		log.Fatalf("Unable to parse X from line: %s", errParsingX)
	}

	y, errParsingY := strconv.Atoi(matches[2])
	if errParsingY != nil {
		log.Fatalf("Unable to parse Y from line: %s", errParsingY)
	}

	return image.Pt(x, y)
}

var prizeRegex = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func parsePrize(line string) image.Point {
	matches := prizeRegex.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("Unable to parse prize from line: %s", line)
	}

	x, errParsingX := strconv.Atoi(matches[1])
	if errParsingX != nil {
		log.Fatalf("Unable to parse X from line: %s", errParsingX)
	}

	y, errParsingY := strconv.Atoi(matches[2])
	if errParsingY != nil {
		log.Fatalf("Unable to parse Y from line: %s", errParsingY)
	}

	return image.Pt(x, y)
}

type Game struct {
	ButtonA image.Point
	ButtonB image.Point
	Prize   image.Point
}

const maxButton = 100
const costButtonA = 3
const costButtonB = 1

func (g Game) Solve() int {
	minC := math.MaxInt

	for a := 0; a <= maxButton; a++ {
		for b := 0; b <= maxButton; b++ {
			if g.ButtonA.X*a+g.ButtonB.X*b == g.Prize.X && g.ButtonA.Y*a+g.ButtonB.Y*b == g.Prize.Y {
				c := costButtonA*a + costButtonB*b
				if c < minC {
					minC = c
				}
			}
		}
	}

	if minC < math.MaxInt {
		return minC
	}

	return 0
}

func parseInput(input io.Reader) []Game {
	scanner := bufio.NewScanner(input)

	var games []Game
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		buttonA := parseButton(scanner.Text())

		if !scanner.Scan() {
			log.Fatalf("Unable to scan the input file correctly for button B: %v", scanner.Err())
		}
		buttonB := parseButton(scanner.Text())

		if !scanner.Scan() {
			log.Fatalf("Unable to scan the input file correctly for button B: %v", scanner.Err())
		}
		prize := parsePrize(scanner.Text())

		games = append(games, Game{
			ButtonA: buttonA,
			ButtonB: buttonB,
			Prize:   prize,
		})
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return games
}

func getResult(input io.Reader) int64 {
	games := parseInput(input)

	var result int64
	for _, game := range games {
		result += int64(game.Solve())
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
