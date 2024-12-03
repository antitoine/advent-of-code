package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

var operations = regexp.MustCompile(`mul\((\d+),(\d+)\)|don't\(\)|do\(\)`)
var enabled = true

func parseLine(line string) int64 {
	matches := operations.FindAllStringSubmatch(line, -1)
	if matches == nil {
		return 0
	}

	var result int64
	for _, match := range matches {
		if len(match) <= 0 {
			log.Fatalf("Unexpected number of matches: %v", match)
		}
		if match[0] == "do()" {
			enabled = true
			continue
		}
		if match[0] == "don't()" {
			enabled = false
			continue
		}
		if len(match) != 3 {
			log.Fatalf("Unexpected number of matches ofr a mult operation: %v", match)
		}
		if !enabled {
			continue
		}
		firstStr, secondStr := match[1], match[2]
		first, errFirst := strconv.ParseInt(firstStr, 10, 64)
		if errFirst != nil {
			log.Fatalf("Unable to parse first number: %v", errFirst)
		}
		second, errSecond := strconv.ParseInt(secondStr, 10, 64)
		if errSecond != nil {
			log.Fatalf("Unable to parse second number: %v", errSecond)
		}
		result += first * second
	}

	return result
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
