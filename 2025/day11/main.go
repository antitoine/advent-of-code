package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func parseInput(input io.Reader) map[string][]string {
	graph := make(map[string][]string)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		device := parts[0]
		outputs := strings.Split(parts[1], " ")
		graph[device] = outputs
	}
	return graph
}

type memoKey struct {
	node       string
	visitedDac bool
	visitedFft bool
}

func countPaths(graph map[string][]string, memo map[memoKey]int64, current, target string, visitedDac, visitedFft bool) int64 {
	if current == "dac" {
		visitedDac = true
	}
	if current == "fft" {
		visitedFft = true
	}

	key := memoKey{current, visitedDac, visitedFft}
	if cached, exists := memo[key]; exists {
		return cached
	}

	if current == target {
		if visitedDac && visitedFft {
			return 1
		}
		return 0
	}

	outputs, exists := graph[current]
	if !exists {
		return 0
	}

	var count int64
	for _, next := range outputs {
		count += countPaths(graph, memo, next, target, visitedDac, visitedFft)
	}

	memo[key] = count
	return count
}

func getResult(input io.Reader) int64 {
	graph := parseInput(input)
	memo := make(map[memoKey]int64)
	return countPaths(graph, memo, "svr", "out", false, false)
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
