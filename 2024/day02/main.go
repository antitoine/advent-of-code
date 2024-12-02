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

func parseLine(line string) []int {
	parts := strings.Split(line, " ")
	numbers := make([]int, len(parts))
	for i, part := range parts {
		number, errNumber := strconv.Atoi(part)
		if errNumber != nil {
			log.Fatalf("Unable to parse number: %v", errNumber)
		}
		numbers[i] = number
	}
	return numbers
}

func parseInput(input io.Reader) int {
	scanner := bufio.NewScanner(input)

	var safeReports int
	for scanner.Scan() {
		numbers := parseLine(scanner.Text())
		if len(numbers) <= 1 {
			safeReports++
			continue
		}
		ascending := numbers[0] < numbers[1]
		isSafe := true
		for i := 0; i < len(numbers)-1; i++ {
			if numbers[i] == numbers[i+1] {
				isSafe = false
				break
			} else if ascending && (numbers[i] > numbers[i+1] || numbers[i+1]-numbers[i] > 3) {
				isSafe = false
				break
			} else if !ascending && (numbers[i] < numbers[i+1] || numbers[i]-numbers[i+1] > 3) {
				isSafe = false
				break
			}
		}
		if isSafe {
			safeReports++
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return safeReports
}

func getResult(input io.Reader) int {
	return parseInput(input)
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
