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

func getTestValueIfValid(testNum int64, values []int64, currentResult int64) int64 {
	if currentResult > testNum {
		return 0
	}
	if len(values) == 0 {
		if testNum == currentResult {
			return testNum
		}
		return 0
	}
	if currentResult == 0 {
		return getTestValueIfValid(testNum, values[1:], values[0])
	}
	addition := getTestValueIfValid(testNum, values[1:], currentResult+values[0])
	if addition > 0 {
		return addition
	}
	multiplication := getTestValueIfValid(testNum, values[1:], currentResult*values[0])
	if multiplication > 0 {
		return multiplication
	}
	concatSum, errConcat := strconv.ParseInt(fmt.Sprintf("%d%d", currentResult, values[0]), 10, 64)
	if errConcat != nil {
		log.Fatalf("Unable to concatenate %d and %d: %v", currentResult, values[0], errConcat)
	}
	concatenation := getTestValueIfValid(testNum, values[1:], concatSum)
	if concatenation > 0 {
		return concatenation
	}
	return 0
}

func parseLine(line string) int64 {
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		log.Fatalf("Invalid line: %s", line)
	}
	testNum, errTestNum := strconv.ParseInt(parts[0], 10, 64)
	if errTestNum != nil {
		log.Fatalf("Invalid test number '%s': %v", parts[0], errTestNum)
	}
	valuesStr := strings.Split(parts[1], " ")
	var values []int64
	for _, valueStr := range valuesStr {
		value, errValue := strconv.ParseInt(valueStr, 10, 64)
		if errValue != nil {
			log.Fatalf("Invalid value '%s': %v", valueStr, errValue)
		}
		values = append(values, value)
	}
	return getTestValueIfValid(testNum, values, 0)
}

func parseInput(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)

	var result int64
	for scanner.Scan() {
		result += parseLine(scanner.Text())
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return result
}

func getResult(input io.Reader) int64 {
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
