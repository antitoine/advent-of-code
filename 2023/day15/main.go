package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func ascii(char rune) int64 {
	return int64(char)
}

func hash(input string) int64 {
	var result int64
	for _, char := range input {
		if char == '\n' {
			continue
		}
		result += ascii(char)
		result *= 17
		result = result % 256
	}
	return result
}

func parseInput(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)

	var result int64
	for scanner.Scan() {
		for _, step := range strings.Split(scanner.Text(), ",") {
			result += hash(step)
		}
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
