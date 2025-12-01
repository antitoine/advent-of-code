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

const initialValue = 50

var lineRegex = regexp.MustCompile(`^([LR])(\d+)$`)

func parseLine(line string) int {
	matches := lineRegex.FindStringSubmatch(line)
	if matches == nil || len(matches) != 3 {
		log.Fatalf("Unable to parse line: %s", line)
	}
	number, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatalf("Unable to parse number: %v", err)
	}
	if matches[1] == "L" {
		return -number
	}
	return number
}

func computePassword(input io.Reader) int {
	scanner := bufio.NewScanner(input)
	nbOfZeroes := 0
	currentValue := initialValue
	for scanner.Scan() {
		rotation := parseLine(scanner.Text())
		nbOfZeroes += max(rotation, -rotation) / 100

		newValue := currentValue + rotation%100
		if (rotation > 0 && newValue >= 100) || (rotation < 0 && newValue <= 0 && currentValue > 0) {
			nbOfZeroes++
		}
		currentValue = (newValue + 100) % 100
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return nbOfZeroes
}

func getResult(input io.Reader) int {
	return computePassword(input)
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
