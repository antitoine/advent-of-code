package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func parseLine(line string) (int, int) {
	parts := strings.Split(line, "   ")
	if len(parts) != 2 {
		log.Fatalf("Unable to parse line: %s", line)
	}
	firstNumber, errFirstNumber := strconv.Atoi(parts[0])
	if errFirstNumber != nil {
		log.Fatalf("Unable to parse first number: %v", errFirstNumber)
	}
	secondNumber, errSecondNumber := strconv.Atoi(parts[1])
	if errSecondNumber != nil {
		log.Fatalf("Unable to parse second number: %v", errSecondNumber)
	}
	return firstNumber, secondNumber
}

func parseInput(input io.Reader) ([]int, []int) {
	scanner := bufio.NewScanner(input)

	var firstList []int
	var secondList []int
	for scanner.Scan() {
		firstNumber, secondNumber := parseLine(scanner.Text())
		firstList = append(firstList, firstNumber)
		secondList = append(secondList, secondNumber)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	if len(firstList) != len(secondList) {
		log.Fatalf("First and second list have different lengths: %d and %d", len(firstList), len(secondList))
	}

	return firstList, secondList
}

func getResult(input io.Reader) int {
	firstList, secondList := parseInput(input)
	slices.Sort(firstList)
	slices.Sort(secondList)

	var result int
	for i := 0; i < len(firstList); i++ {
		if firstList[i] < secondList[i] {
			result += secondList[i] - firstList[i]
		} else {
			result += firstList[i] - secondList[i]
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
