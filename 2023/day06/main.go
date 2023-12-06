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

func numberStrToInt(numberStr string) int64 {
	numberWithoutSpaces := strings.ReplaceAll(numberStr, " ", "")
	numberInt, errParsingNumber := strconv.ParseInt(strings.ReplaceAll(numberWithoutSpaces, " ", ""), 10, 64)
	if errParsingNumber != nil {
		log.Fatalf("Unable to parse number '%s': %v", numberWithoutSpaces, errParsingNumber)
	}
	return numberInt
}

func parseInput(input io.Reader) (int64, int64) {
	scanner := bufio.NewScanner(input)

	// Seeds
	if !scanner.Scan() {
		log.Fatalf("Unable to scan the first line of the input correctly: %v", scanner.Err())
	}
	timeStr := scanner.Text()
	timeStrSplit := strings.Split(timeStr, ":")
	if len(timeStrSplit) != 2 {
		log.Fatalf("Unable to parse the first line of the input correctly after the split: %s", timeStr)
	}
	t := numberStrToInt(timeStrSplit[1])

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the second line of the input correctly: %v", scanner.Err())
	}

	distanceStr := scanner.Text()
	distanceStrSplit := strings.Split(distanceStr, ":")
	if len(distanceStrSplit) != 2 {
		log.Fatalf("Unable to parse the second line of the input correctly after the split: %s", distanceStr)
	}
	distance := numberStrToInt(distanceStrSplit[1])

	return t, distance
}

func getDistanceTraveledForTime(timeHolding int64, totalTime int64) int64 {
	return (totalTime - timeHolding) * timeHolding
}

func getMinimumTimeHoldingForDistance(distance int64, totalTime int64) int64 {
	for i := int64(0); i <= totalTime; i++ {
		if getDistanceTraveledForTime(i, totalTime) > distance {
			return i
		}
	}
	return -1
}

func getMaximumTimeHoldingForDistance(distance int64, totalTime int64) int64 {
	for i := totalTime; i >= 0; i-- {
		if getDistanceTraveledForTime(i, totalTime) > distance {
			return i
		}
	}
	return -1
}

func getResult(input io.Reader) int64 {
	t, d := parseInput(input)
	maximumTimeHolding := getMaximumTimeHoldingForDistance(d, t)
	minimumTimeHolding := getMinimumTimeHoldingForDistance(d, t)
	possibilities := maximumTimeHolding - minimumTimeHolding + 1

	return possibilities
}

func main() {
	start := time.Now()
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
