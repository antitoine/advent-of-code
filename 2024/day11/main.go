package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseLine(line string) []int64 {
	parts := strings.Split(line, " ")
	result := make([]int64, len(parts))
	for i, part := range parts {
		num, errParsing := strconv.ParseInt(part, 10, 64)
		if errParsing != nil {
			log.Fatalf("Unable to parse the number '%s': %v", part, errParsing)
		}
		result[i] = num
	}
	return result
}

func parseInput(input io.Reader) []int64 {
	scanner := bufio.NewScanner(input)

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the input file correctly")
	}

	stones := parseLine(scanner.Text())

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return stones
}

func splitDigits(stone int64) (int64, int64, bool) {
	stoneStr := fmt.Sprintf("%d", stone)
	if len(stoneStr)%2 == 1 {
		return 0, 0, false
	}
	leftStr := stoneStr[:len(stoneStr)/2]
	leftNum, errParsingLeft := strconv.ParseInt(leftStr, 10, 64)
	if errParsingLeft != nil {
		log.Fatalf("Unable to parse the number '%s': %v", leftStr, errParsingLeft)
	}
	rightStr := stoneStr[len(stoneStr)/2:]
	rightNum, errParsingRight := strconv.ParseInt(rightStr, 10, 64)
	if errParsingRight != nil {
		log.Fatalf("Unable to parse the number '%s': %v", rightStr, errParsingRight)
	}
	return leftNum, rightNum, true
}

func getResult(input io.Reader) int {
	state := parseInput(input)

	for b := 0; b < 25; b++ {
		var newState []int64
		for _, stone := range state {
			if stone == 0 {
				newState = append(newState, 1)
			} else if left, right, isEvenDigits := splitDigits(stone); isEvenDigits {
				newState = append(newState, left, right)
			} else {
				newState = append(newState, stone*2024)
			}
		}
		state = newState
	}

	return len(state)
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
