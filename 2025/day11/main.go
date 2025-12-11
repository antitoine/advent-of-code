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

func countPaths(graph map[string][]string, current, target string) int64 {
	if current == target {
		return 1
	}
	outputs, exists := graph[current]
	if !exists {
		return 0
	}
	var count int64
	for _, next := range outputs {
		count += countPaths(graph, next, target)
	}
	return count
}

func getResult(input io.Reader) int64 {
	graph := parseInput(input)
	return countPaths(graph, "you", "out")
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
