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

func parseRule(line string) (int, int) {
	parts := strings.Split(line, "|")
	if len(parts) != 2 {
		log.Fatalf("Invalid rule: %s", line)
	}
	before, errBefore := strconv.Atoi(parts[0])
	if errBefore != nil {
		log.Fatalf("Invalid rule for before number: %s", line)
	}
	after, errAfter := strconv.Atoi(parts[1])
	if errAfter != nil {
		log.Fatalf("Invalid rule for after number: %s", line)
	}
	return before, after
}

func parseUpdate(line string) []int {
	var updates []int
	parts := strings.Split(line, ",")
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			log.Fatalf("Invalid update number: %s", line)
		}
		updates = append(updates, num)
	}
	return updates
}

func parseInput(input io.Reader) (map[int][]int, [][]int) {
	scanner := bufio.NewScanner(input)

	rules := make(map[int][]int)
	var updates [][]int
	switchToUpdates := false
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			switchToUpdates = true
			continue
		}
		if !switchToUpdates {
			before, after := parseRule(text)
			rules[before] = append(rules[before], after)
		} else {
			updates = append(updates, parseUpdate(text))
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return rules, updates
}

func sortUpdate(rules map[int][]int, update []int) bool {
	valid := true
	process := make(map[int]int)
	for i, nb := range update {
		process[nb] = i
		if rule, found := rules[nb]; found {
			for _, after := range rule {
				if idx, foundAfter := process[after]; foundAfter {
					valid = false
					for j := i; j > idx; j-- {
						process[update[j-1]] = j
						update[j] = update[j-1]
					}
					update[idx] = nb
					process[nb] = idx
					break
				}
			}
		}
	}
	return valid
}

func getResult(input io.Reader) int64 {
	rules, updates := parseInput(input)

	var result int64
	for _, update := range updates {
		isAlreadySorted := sortUpdate(rules, update)
		isValid := isAlreadySorted
		for !isValid {
			isValid = sortUpdate(rules, update)
		}
		if !isAlreadySorted {
			middlePage := update[len(update)/2]
			result += int64(middlePage)
		}
	}

	return result
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
