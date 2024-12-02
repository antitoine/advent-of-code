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

func isSafeLevels(ascending bool, a, b int) bool {
	if a == b {
		return false
	} else if ascending && (a > b || b-a > 3) {
		return false
	} else if !ascending && (a < b || a-b > 3) {
		return false
	}
	return true
}

func isSafeReport(numbers []int) bool {
	ascending := numbers[0] < numbers[1]
	for i := 0; i < len(numbers)-1; i++ {
		if !isSafeLevels(ascending, numbers[i], numbers[i+1]) {
			return false
		}
	}
	return true
}

func isSafeReportWithOneError(numbers []int) bool {
	aLessThanBCounter := 0
	for i := 0; i < 4; i++ {
		if numbers[i] < numbers[i+1] {
			aLessThanBCounter++
		} else {
			aLessThanBCounter--
		}
	}
	ascending := true
	if aLessThanBCounter < 0 {
		ascending = false
	}

	if !isSafeLevels(ascending, numbers[0], numbers[1]) {
		newSliceA := append([]int{}, numbers[1:]...)
		newSliceB := append([]int{numbers[0]}, numbers[2:]...)
		if !isSafeReport(newSliceB) && !isSafeReport(newSliceA) {
			return false
		}

	}
	for i := 1; i < len(numbers)-1; i++ {
		if !isSafeLevels(ascending, numbers[i], numbers[i+1]) {
			leftA := append([]int{}, numbers[:i]...)
			newSliceA := append(leftA, numbers[i+1:]...)
			leftB := append([]int{}, numbers[:i+1]...)
			newSliceB := append(leftB, numbers[i+2:]...)
			if !isSafeReport(newSliceB) && !isSafeReport(newSliceA) {
				return false
			}

		}
	}
	return true
}

func parseInput(input io.Reader) int {
	scanner := bufio.NewScanner(input)

	var safeReports int
	for scanner.Scan() {
		numbers := parseLine(scanner.Text())
		if isSafeReportWithOneError(numbers) {
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
