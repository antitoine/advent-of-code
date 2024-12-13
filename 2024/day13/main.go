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

const prizeCorrection = 10000000000000

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

	return image.Pt(x+prizeCorrection, y+prizeCorrection)
}

type Game struct {
	ButtonA image.Point
	ButtonB image.Point
	Prize   image.Point
}

const costButtonA = 3
const costButtonB = 1

func (g Game) Solve() int {
	det := g.ButtonA.X*g.ButtonB.Y - g.ButtonA.Y*g.ButtonB.X
	if det == 0 {
		return 0
	}
	xNumerator := g.Prize.X*g.ButtonB.Y - g.Prize.Y*g.ButtonB.X
	yNumerator := g.ButtonA.X*g.Prize.Y - g.ButtonA.Y*g.Prize.X

	if xNumerator%det != 0 || yNumerator%det != 0 {
		return 0
	}

	aParticular := xNumerator / det
	bParticular := yNumerator / det

	return costButtonA*aParticular + costButtonB*bParticular
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
