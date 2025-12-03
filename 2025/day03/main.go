package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

const digitsToSelect = 12

func maxJoltage(bank string) int64 {
	n := len(bank)
	result := int64(0)
	start := 0

	for i := 0; i < digitsToSelect; i++ {
		end := n - digitsToSelect + i

		maxIdx := start
		for j := start + 1; j <= end; j++ {
			if bank[j] > bank[maxIdx] {
				maxIdx = j
			}
		}

		result = result*10 + int64(bank[maxIdx]-'0')
		start = maxIdx + 1
	}

	return result
}

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	total := int64(0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		total += maxJoltage(line)
	}
	return total
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
