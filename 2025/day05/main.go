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

type Range struct {
	start, end int64
}

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)

	var ranges []Range
	var ingredientIDs []int64
	parsingRanges := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			parsingRanges = false
			continue
		}

		if parsingRanges {
			parts := strings.Split(line, "-")
			start, _ := strconv.ParseInt(parts[0], 10, 64)
			end, _ := strconv.ParseInt(parts[1], 10, 64)
			ranges = append(ranges, Range{start: start, end: end})
		} else {
			id, _ := strconv.ParseInt(line, 10, 64)
			ingredientIDs = append(ingredientIDs, id)
		}
	}

	var freshCount int64
	for _, id := range ingredientIDs {
		for _, r := range ranges {
			if id >= r.start && id <= r.end {
				freshCount++
				break
			}
		}
	}

	return freshCount
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
