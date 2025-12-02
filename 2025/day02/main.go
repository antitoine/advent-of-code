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

type idRange struct {
	start, end int
}

func parseInput(input io.Reader) ([]idRange, int) {
	var ranges []idRange
	maxId := 0
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				continue
			}
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])
			if end > maxId {
				maxId = end
			}
			ranges = append(ranges, idRange{start, end})
		}
	}
	return ranges, maxId
}

func generateInvalidNumbers(maxVal int) []int {
	invalidSet := make(map[int]struct{})

	// Generate all numbers that are a pattern repeated at least twice
	for patternLen := 1; ; patternLen++ {
		minPattern := 1
		if patternLen > 1 {
			minPattern = intPow(10, patternLen-1)
		}
		maxPattern := intPow(10, patternLen) - 1

		// The smallest invalid number for this pattern length (repeated twice)
		smallestInvalid := minPattern*intPow(10, patternLen) + minPattern
		if smallestInvalid > maxVal {
			break
		}

		for pattern := minPattern; pattern <= maxPattern; pattern++ {
			// Generate pattern repeated 2, 3, 4, ... times
			invalidNum := pattern
			multiplier := intPow(10, patternLen)
			for repeatCount := 2; ; repeatCount++ {
				invalidNum = invalidNum*multiplier + pattern
				if invalidNum > maxVal {
					break
				}
				invalidSet[invalidNum] = struct{}{}
			}
		}
	}

	invalid := make([]int, 0, len(invalidSet))
	for num := range invalidSet {
		invalid = append(invalid, num)
	}
	return invalid
}

func intPow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func getResult(input io.Reader) int {
	ranges, maxId := parseInput(input)

	invalidNumbers := generateInvalidNumbers(maxId)

	sum := 0
	for _, num := range invalidNumbers {
		for _, r := range ranges {
			if num >= r.start && num <= r.end {
				sum += num
				break
			}
		}
	}

	return sum
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
