package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func parseInventory(line string) map[rune][]string {
	parts := strings.Split(line, ", ")
	inventory := make(map[rune][]string, len(parts))
	for _, part := range parts {
		inventory[rune(part[0])] = append(inventory[rune(part[0])], part)
	}
	return inventory
}

func parseModel(line string) string {
	return line
}

func parseInput(input io.Reader) (map[rune][]string, []string) {
	scanner := bufio.NewScanner(input)

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the input file correctly")
	}
	inventory := parseInventory(scanner.Text())

	var models []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		models = append(models, parseModel(line))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return inventory, models
}

var cache = make(map[string]bool)

func isModelCompatible(model string, inventory map[rune][]string) bool {
	if model == "" {
		return true
	}
	if result, ok := cache[model]; ok {
		return result
	}
	firstColor := rune(model[0])
	available, ok := inventory[firstColor]
	if !ok {
		return false
	}
	for _, part := range available {
		if strings.HasPrefix(model, part) && isModelCompatible(model[len(part):], inventory) {
			cache[model] = true
			return true
		}
	}
	cache[model] = false
	return false
}

func getResult(input io.Reader) int64 {
	inventory, models := parseInput(input)

	var nbCompatibleModels int64
	for _, model := range models {
		if isModelCompatible(model, inventory) {
			nbCompatibleModels++
		}
	}

	return nbCompatibleModels
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
