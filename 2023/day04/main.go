package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var cardLineRegex = regexp.MustCompile(`^Card\s*(\d*): ([\s\d]*) \| ([\s\d]*)$`)
var numbersRegex = regexp.MustCompile(`\s*(\d*)\s*`)

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
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
	return numbers
}

func getSumOfWinningCards(file *os.File) int64 {
	var finalSum int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		results := cardLineRegex.FindStringSubmatch(line)
		if len(results) != 4 {
			log.Fatalf("Unable to parse card line '%s', getting results length of %d", line, len(results))
		}
		winningNumbers := numbersStrToInts(results[2])
		gettingNumbers := numbersStrToInts(results[3])

		winningIndex := 0
		gettingIndex := 0
		var points int64
		for {
			if winningIndex >= len(winningNumbers) || gettingIndex >= len(gettingNumbers) {
				break
			} else if winningNumbers[winningIndex] == gettingNumbers[gettingIndex] {
				if points == 0 {
					points = 1
				} else {
					points *= 2
				}
				gettingIndex++
			} else if winningNumbers[winningIndex] < gettingNumbers[gettingIndex] {
				winningIndex++
			} else if winningNumbers[winningIndex] > gettingNumbers[gettingIndex] {
				gettingIndex++
			} else {
				log.Fatalf("Unkown case: winningIndex=%d, gettingIndex=%d", winningIndex, gettingIndex)
			}
		}
		finalSum += points
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return finalSum
}

func main() {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	finalSum := getSumOfWinningCards(inputFile)

	log.Printf("Final sum: %d", finalSum)
}
