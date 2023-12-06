package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var numbersRegex = regexp.MustCompile(`\s*(\d*)\s*`)

func numbersStrToInts(numbersStr string) []int64 {
	var numbers []int64
	for _, numberStr := range numbersRegex.FindAllStringSubmatch(numbersStr, -1) {
		if len(numberStr) != 2 {
			log.Fatalf("Unable to parse number '%s', getting results length of %d", numbersStr, len(numberStr))
		}
		numberInt, errParsingNumber := strconv.ParseInt(numberStr[1], 10, 64)
		if errParsingNumber != nil {
			log.Fatalf("Unable to parse number '%s': %v", numberStr[1], errParsingNumber)
		}
		numbers = append(numbers, numberInt)
	}
	return numbers
}

func parseInput(input io.Reader) ([]int64, []int64) {
	scanner := bufio.NewScanner(input)

	// Seeds
	if !scanner.Scan() {
		log.Fatalf("Unable to scan the first line of the input correctly: %v", scanner.Err())
	}
	timesStr := scanner.Text()
	timesStrSplit := strings.Split(timesStr, ":")
	if len(timesStrSplit) != 2 {
		log.Fatalf("Unable to parse the first line of the input correctly after the split: %s", timesStr)
	}
	times := numbersStrToInts(timesStrSplit[1])

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the second line of the input correctly: %v", scanner.Err())
	}

	distancesStr := scanner.Text()
	distancesStrSplit := strings.Split(distancesStr, ":")
	if len(distancesStrSplit) != 2 {
		log.Fatalf("Unable to parse the second line of the input correctly after the split: %s", timesStr)
	}
	distances := numbersStrToInts(distancesStrSplit[1])

	return times, distances
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
	times, distances := parseInput(input)
	if len(times) != len(distances) {
		log.Fatalf("Unable to calculate the result, times and distances have different lengths: %d and %d", len(times), len(distances))
	}
	result := int64(1)
	for i := 0; i < len(times); i++ {
		maximumTimeHolding := getMaximumTimeHoldingForDistance(distances[i], times[i])
		minimumTimeHolding := getMinimumTimeHoldingForDistance(distances[i], times[i])
		possibilities := maximumTimeHolding - minimumTimeHolding + 1
		result *= possibilities
	}

	return result
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
