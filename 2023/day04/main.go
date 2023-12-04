package main

import (
	"bufio"
	"io"
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

func getSumOfWinningCards(input io.Reader) int64 {
	var finalSum int64
	var cardsFactor []int64
	scanner := bufio.NewScanner(input)
	for cardIdx := 0; scanner.Scan(); cardIdx++ {
		line := scanner.Text()

		results := cardLineRegex.FindStringSubmatch(line)
		if len(results) != 4 {
			log.Fatalf("Unable to parse card line '%s', getting results length of %d", line, len(results))
		}
		winningNumbers := numbersStrToInts(results[2])
		gettingNumbers := numbersStrToInts(results[3])

		winningIndex := 0
		gettingIndex := 0
		match := 0
		for {
			if winningIndex >= len(winningNumbers) || gettingIndex >= len(gettingNumbers) {
				break
			} else if winningNumbers[winningIndex] == gettingNumbers[gettingIndex] {
				match++
				gettingIndex++
			} else if winningNumbers[winningIndex] < gettingNumbers[gettingIndex] {
				winningIndex++
			} else if winningNumbers[winningIndex] > gettingNumbers[gettingIndex] {
				gettingIndex++
			} else {
				log.Fatalf("Unkown case: winningIndex=%d, gettingIndex=%d", winningIndex, gettingIndex)
			}
		}
		if cardIdx >= len(cardsFactor) {
			cardsFactor = append(cardsFactor, 1)
		} else {
			cardsFactor[cardIdx]++
		}
		for i := 1; i <= match; i++ {
			if cardIdx+i >= len(cardsFactor) {
				cardsFactor = append(cardsFactor, 1*cardsFactor[cardIdx])
			} else {
				cardsFactor[cardIdx+i] += 1 * cardsFactor[cardIdx]
			}
		}
		finalSum += cardsFactor[cardIdx]
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
