package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

var numbersRegex = regexp.MustCompile(`\s*(-?\d*)\s*`)

func numbersStrToInts(numbersStr string) []int64 {
	var numbers []int64
	for _, numberStr := range numbersRegex.FindAllStringSubmatch(numbersStr, -1) {
		if len(numberStr) != 2 {
			log.Fatalf("Unable to parse number '%s', getting results length of %d", numbersStr, len(numberStr))
		}
		numberInt, errParsingNumber := strconv.ParseInt(numberStr[1], 10, 64)
		if errParsingNumber != nil {
			log.Fatalf("Unable to parse number '%s': %v", numberStr[1], errParsingNumber)
		}
		numbers = append(numbers, numberInt)
	}
	return numbers
}

func parseInput(input io.Reader) [][]int64 {
	scanner := bufio.NewScanner(input)
	var linesNumbers [][]int64
	for scanner.Scan() {
		linesNumbers = append(linesNumbers, numbersStrToInts(scanner.Text()))
	}
	return linesNumbers
}

func derivative(numbers []int64) ([]int64, bool) {
	if len(numbers) == 0 {
		log.Fatalf("Unable to calculate derivative of empty numbers")
		return []int64{}, false
	}
	onlyZeros := true
	derivativeNumbers := make([]int64, len(numbers)-1)
	for i := 1; i < len(numbers); i++ {
		derivativeNumbers[i-1] = numbers[i] - numbers[i-1]
		if derivativeNumbers[i-1] != 0 {
			onlyZeros = false
		}
	}
	return derivativeNumbers, onlyZeros
}

func extrapolateForward(numbers []int64) int64 {
	if len(numbers) == 0 {
		log.Fatalf("Unable to extrapolate empty numbers")
		return 0
	}
	derivativeNumbers, onlyZeros := derivative(numbers)
	if onlyZeros {
		return numbers[len(numbers)-1]
	} else {
		return numbers[len(numbers)-1] + extrapolateForward(derivativeNumbers)
	}
}

func extrapolateBackward(numbers []int64) int64 {
	if len(numbers) == 0 {
		log.Fatalf("Unable to extrapolate empty numbers")
		return 0
	}
	derivativeNumbers, onlyZeros := derivative(numbers)
	if onlyZeros {
		return numbers[0]
	} else {
		return numbers[0] - extrapolateBackward(derivativeNumbers)
	}
}

func getResultPart1(input io.Reader) int64 {
	var result int64
	for _, numbers := range parseInput(input) {
		result += extrapolateForward(numbers)
	}
	return result
}

func getResultPart2(input io.Reader) int64 {
	var result int64
	for _, numbers := range parseInput(input) {
		result += extrapolateBackward(numbers)
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

	result := getResultPart2(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
