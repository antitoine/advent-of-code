package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func main() {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	var finalSum int64
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		textLine := scanner.Text()
		firstDigitIdx := strings.IndexFunc(textLine, isDigit)
		if firstDigitIdx < 0 {
			log.Fatalf("Unable to find the first digit in the line: %v", textLine)
		}
		LastDigitIdx := strings.LastIndexFunc(textLine, isDigit)
		if LastDigitIdx < 0 {
			log.Fatalf("Unable to find the last digit in the line: %v", textLine)
		}
		textLineRunes := []rune(textLine)
		lineNumber, errParsing := strconv.ParseInt(string(textLineRunes[firstDigitIdx])+string(textLineRunes[LastDigitIdx]), 10, 64)
		if errParsing != nil {
			log.Fatalf("Unable to parse the line number (%v with %v) in line text '%s': %v", textLine[firstDigitIdx], textLine[LastDigitIdx], textLine, errParsing)
		}
		finalSum += lineNumber
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	log.Printf("Final sum: %d", finalSum)
}
