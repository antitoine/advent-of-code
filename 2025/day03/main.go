package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func maxJoltage(bank string) int {
	maxVal := 0
	for i := 0; i < len(bank)-1; i++ {
		d1 := int(bank[i] - '0')
		for j := i + 1; j < len(bank); j++ {
			d2 := int(bank[j] - '0')
			joltage := d1*10 + d2
			if joltage > maxVal {
				maxVal = joltage
			}
		}
	}
	return maxVal
}

func getResult(input io.Reader) int {
	scanner := bufio.NewScanner(input)
	total := 0
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
