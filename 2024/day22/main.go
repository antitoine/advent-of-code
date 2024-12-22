package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
    Calculate the result of multiplying the secret number by 64. Then, mix this result into the secret number. Finally, prune the secret number.
    Calculate the result of dividing the secret number by 32. Round the result down to the nearest integer. Then, mix this result into the secret number. Finally, prune the secret number.
    Calculate the result of multiplying the secret number by 2048. Then, mix this result into the secret number. Finally, prune the secret number.

Each step of the above process involves mixing and pruning:

    To mix a value into the secret number, calculate the bitwise XOR of the given value and the secret number. Then, the secret number becomes the result of that operation. (If the secret number is 42 and you were to mix 15 into the secret number, the secret number would become 37.)
    To prune the secret number, calculate the value of the secret number modulo 16777216. Then, the secret number becomes the result of that operation. (If the secret number is 100000000 and you were to prune the secret number, the secret number would become 16113920.)
*/

func mixSecret(secret, value int64) int64 {
	return secret ^ value
}

func pruneSecret(secret int64) int64 {
	return secret % 16777216
}

func nextSecret(secret int64) int64 {
	step1 := pruneSecret(mixSecret(secret, secret*64))
	step2 := pruneSecret(mixSecret(step1, step1/32))
	step3 := pruneSecret(mixSecret(step2, step2*2048))
	return step3
}

func parseLine(line string) int64 {
	val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
	if err != nil {
		log.Fatalf("Unable to parse line: %v", err)
	}
	return val
}

func parseInput(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)

	var sumSecrets int64
	for scanner.Scan() {
		initialSecret := parseLine(scanner.Text())
		secret := initialSecret
		for i := 0; i < 2000; i++ {
			secret = nextSecret(secret)
		}
		sumSecrets += secret
		fmt.Printf("%d: %d\n", initialSecret, secret)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return sumSecrets
}

func getResult(input io.Reader) int64 {
	return parseInput(input)
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
