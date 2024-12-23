package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Connections map[string][]string

func parseLine(line string) (string, string) {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		log.Fatalf("Expected 2 parts, got %d for line %s", len(parts), line)
	}
	return parts[0], parts[1]
}

func parseInput(input io.Reader) Connections {
	scanner := bufio.NewScanner(input)

	connections := make(Connections)
	for scanner.Scan() {
		nodeA, nodeB := parseLine(scanner.Text())
		connections[nodeA] = append(connections[nodeA], nodeB)
		connections[nodeB] = append(connections[nodeB], nodeA)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return connections
}

func getGroupId(nodes []string) string {
	sort.Strings(nodes)
	return fmt.Sprintf("%s-%s-%s", nodes[0], nodes[1], nodes[2])
}

func getResult(input io.Reader) int {
	connections := parseInput(input)

	groups := make(map[string]int)
	for node, neighbours := range connections {
		for n := 0; n < len(neighbours)-1; n++ {
			for m := n + 1; m < len(neighbours); m++ {
				if strings.HasPrefix(node, "t") || strings.HasPrefix(neighbours[n], "t") || strings.HasPrefix(neighbours[m], "t") {
					groupId := getGroupId([]string{node, neighbours[n], neighbours[m]})
					groups[groupId]++
				}
			}
		}
	}

	var clusters []string
	for groupId, count := range groups {
		if count >= 3 {
			clusters = append(clusters, groupId)
		}
	}

	return len(clusters)
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
