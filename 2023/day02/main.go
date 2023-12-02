package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Bag struct {
	red   int64
	green int64
	blue  int64
}

var drawsLineRegex = regexp.MustCompile(`(\d*) (red|blue|green)`)

func parseGameDraws(drawsLine string) ([]Bag, Bag) {
	var draws []Bag
	var minimumBagForGame Bag
	for _, drawStr := range strings.Split(drawsLine, "; ") {
		cubesMatches := drawsLineRegex.FindAllStringSubmatch(drawStr, -1)
		var draw Bag
		for _, cubeMatch := range cubesMatches {
			if len(cubeMatch) != 3 {
				log.Fatalf("Unable to parse cube match '%s', getting results length of %d", cubeMatch, len(cubeMatch))
			}
			cubeCount, errParsingCount := strconv.ParseInt(cubeMatch[1], 10, 64)
			if errParsingCount != nil {
				log.Fatalf("Unable to parse cube count '%s': %v", cubeMatch[1], errParsingCount)
			}
			cubeColor := cubeMatch[2]
			if cubeColor == "red" {
				draw.red += cubeCount
			} else if cubeColor == "green" {
				draw.green += cubeCount
			} else if cubeColor == "blue" {
				draw.blue += cubeCount
			} else {
				log.Fatalf("Unkown cube color '%s'", cubeColor)
			}
		}
		draws = append(draws, draw)
		if minimumBagForGame.red < draw.red {
			minimumBagForGame.red = draw.red
		}
		if minimumBagForGame.green < draw.green {
			minimumBagForGame.green = draw.green
		}
		if minimumBagForGame.blue < draw.blue {
			minimumBagForGame.blue = draw.blue
		}
	}
	return draws, minimumBagForGame
}

var gameLineRegex = regexp.MustCompile(`^Game (\d*): (.*)$`)

func parseGameLine(line string) (int64, []Bag, Bag) {
	results := gameLineRegex.FindStringSubmatch(line)
	if len(results) != 3 {
		log.Fatalf("Unable to parse game line '%s', getting results length of %d", line, len(results))
	}
	gameId, errParsingId := strconv.ParseInt(results[1], 10, 64)
	if errParsingId != nil {
		log.Fatalf("Unable to parse game id '%s': %v", results[1], errParsingId)
	}
	draws, minimumBagForGame := parseGameDraws(results[2])
	return gameId, draws, minimumBagForGame
}

func getSumOfGamePowerCubes(file *os.File) int64 {
	var finalSum int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, _, minimumBagForGame := parseGameLine(line)
		finalSum += minimumBagForGame.red * minimumBagForGame.green * minimumBagForGame.blue
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return finalSum
}

func main() {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	finalSum := getSumOfGamePowerCubes(inputFile)

	log.Printf("Final sum: %d", finalSum)
}
