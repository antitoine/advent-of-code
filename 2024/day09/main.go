package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func parseInput(input io.Reader) []int {
	scanner := bufio.NewScanner(input)

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the input file correctly")
	}

	line := scanner.Text()
	result := make([]int, 0)
	isBlocks := true
	blockNumber := 0
	for _, numStr := range line {
		num, errParsing := strconv.Atoi(string(numStr))
		if errParsing != nil {
			log.Fatalf("Unable to parse the number '%s': %v", string(numStr), errParsing)
		}
		for i := 0; i < num; i++ {
			if isBlocks {
				result = append(result, blockNumber)
			} else {
				result = append(result, -1)
			}
		}
		if isBlocks {
			blockNumber++
		}
		isBlocks = !isBlocks
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return result
}

func getResult(input io.Reader) int64 {
	initialState := parseInput(input)
	var checksum int64
	i := 0
	for i < len(initialState) && initialState[i] != -1 {
		checksum += int64(initialState[i]) * int64(i)
		i++
	}
	for j := len(initialState) - 1; j > i; j-- {
		if initialState[j] == -1 {
			continue
		}
		initialState[i] = initialState[j]
		initialState[j] = -1
		for i < len(initialState) && initialState[i] != -1 {
			checksum += int64(initialState[i]) * int64(i)
			i++
		}
	}
	return checksum
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
