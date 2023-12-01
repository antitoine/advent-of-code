package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var digitAsLetters = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func getFirstDigitFromString(line string) string {
	for i, r := range []rune(line) {
		if isDigit(r) {
			return string(r)
		}
		for letters, digit := range digitAsLetters {
			if strings.HasPrefix(line[i:], letters) {
				return digit
			}
		}
	}
	log.Fatalf("Unable to get first digit from provided line '%s'", line)
	return "0"
}

func getLastDigitFromString(line string) string {
	for i := len(line) - 1; i >= 0; i-- {
		r := rune(line[i])
		if isDigit(r) {
			return string(r)
		}
		for letters, digit := range digitAsLetters {
			if strings.HasSuffix(line[:i+1], letters) {
				return digit
			}
		}
	}
	log.Fatalf("Unable to get last digit from provided line '%s'", line)
	return "0"
}

func getDigitsFromString(line string) int64 {
	firstDigit := getFirstDigitFromString(line)
	lastDigit := getLastDigitFromString(line)
	number, errParsing := strconv.ParseInt(firstDigit+lastDigit, 10, 64)
	if errParsing != nil {
		log.Fatalf("Unable to parse digits provided '%s' & '%s': %v", firstDigit, lastDigit, errParsing)
	}
	return number
}

func getSumOfDigits(file *os.File) int64 {
	var finalSum int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		finalSum += getDigitsFromString(scanner.Text())
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

	finalSum := getSumOfDigits(inputFile)

	log.Printf("Final sum: %d", finalSum)
}
