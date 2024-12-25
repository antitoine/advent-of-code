package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

/*
0,5,3,4,3
1,2,0,5,3

5,0,2,1,3
4,3,4,0,2
3,0,2,0,1
*/
func parseLock(lines []string) [5]int {
	var lock [5]int
	for i := 1; i < len(lines); i++ {
		for j, c := range lines[i] {
			if c == '#' {
				lock[j]++
			}
		}
	}
	return lock
}

func parseKey(lines []string) [5]int {
	var key [5]int
	for i := 0; i < len(lines)-1; i++ {
		for j, c := range lines[i] {
			if c == '#' {
				key[j]++
			}
		}
	}
	return key
}

func parseInput(input io.Reader) ([][5]int, [][5]int) {
	scanner := bufio.NewScanner(input)

	var locks [][5]int
	var keys [][5]int

	parseLines := func(lines []string) {
		if lines[0][0] == '#' {
			locks = append(locks, parseLock(lines))
		} else {
			keys = append(keys, parseKey(lines))
		}
	}

	var currentLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			currentLines = append(currentLines, line)
			continue
		}
		parseLines(currentLines)
		currentLines = nil
	}
	parseLines(currentLines)

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return locks, keys
}

func getResult(input io.Reader) int64 {
	locks, keys := parseInput(input)

	var nbFit int64
	for _, lock := range locks {
		for _, key := range keys {
			fit := true
			for i := 0; i < 5 && fit; i++ {
				if lock[i]+key[i] > 5 {
					fit = false
				}
			}
			if fit {
				nbFit++
			}
		}
	}

	return nbFit
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
