package main

import (
	"bufio"
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

func getGroupId(nodes ...string) string {
	sort.Strings(nodes)
	return strings.Join(nodes, ",")
}

func getAllGroups(nodes []string, size int) [][]string {
	if size == 0 || size > len(nodes) {
		return [][]string{}
	}

	var result [][]string
	var combination []string
	var helper func(int, int)

	helper = func(start, remaining int) {
		if remaining == 0 {
			group := make([]string, len(combination))
			copy(group, combination)
			result = append(result, group)
			return
		}

		for i := start; i <= len(nodes)-remaining; i++ {
			combination = append(combination, nodes[i])
			helper(i+1, remaining-1)
			combination = combination[:len(combination)-1] // backtrack
		}
	}

	helper(0, size)
	return result
}

func getGroups(connections Connections, nbNodes int) []string {
	groups := make(map[string]int)
	for node, neighbours := range connections {
		nodes := append(neighbours, node)
		for _, group := range getAllGroups(nodes, nbNodes) {
			groups[getGroupId(group...)]++
		}
	}

	var clusters []string
	for groupId, count := range groups {
		if count >= nbNodes {
			clusters = append(clusters, groupId)
		}
	}

	return clusters
}

func getResult(input io.Reader) string {
	connections := parseInput(input)

	groupSize := 3
	groups := getGroups(connections, groupSize)
	for len(groups) > 1 {
		groupSize++
		groups = getGroups(connections, groupSize)
	}

	if len(groups) != 1 {
		log.Fatalf("Expected 1 group, got %d", len(groups))
	}

	return groups[0]
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

	log.Printf("Final result: %s", result)
	log.Printf("Execution took %s", time.Since(start))
}
