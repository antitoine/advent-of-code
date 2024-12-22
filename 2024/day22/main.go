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

func parseLine(line string) int64 {
	val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
	if err != nil {
		log.Fatalf("Unable to parse line: %v", err)
	}
	return val
}

func parseInput(input io.Reader) []int64 {
	scanner := bufio.NewScanner(input)

	var secrets []int64
	for b := 0; scanner.Scan(); b++ {
		secrets = append(secrets, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return secrets
}

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

type SequenceOfChanges [4]int64

type NbBananasPerSequence struct {
	delta            SequenceOfChanges
	bananasPerBuyers []*int64
}

func (n *NbBananasPerSequence) bananas() int64 {
	var sum int64
	for _, nbBananas := range n.bananasPerBuyers {
		if nbBananas != nil {
			sum += *nbBananas
		}
	}
	return sum
}

func getNumbersOfBananasFromSecret(secret int64) int64 {
	return secret % 10
}

func getResult(input io.Reader) int64 {
	secrets := parseInput(input)

	sequences := make(map[SequenceOfChanges]*NbBananasPerSequence)
	for b := 0; b < len(secrets); b++ {
		previousSecret := secrets[b]
		previousNbBananas := getNumbersOfBananasFromSecret(previousSecret)
		var delta SequenceOfChanges
		for i := 0; i < 2000; i++ {
			newSecret := nextSecret(previousSecret)
			newNbBananas := getNumbersOfBananasFromSecret(newSecret)
			change := newNbBananas - previousNbBananas
			if i <= 3 {
				delta[i] = change
			} else {
				delta = SequenceOfChanges{delta[1], delta[2], delta[3], change}
			}
			if i >= 3 {
				if nbBananasForSeq, ok := sequences[delta]; ok {
					if nbBananasForSeq.bananasPerBuyers[b] == nil {
						nbBananasForSeq.bananasPerBuyers[b] = &newNbBananas
					}
				} else {
					newSeq := &NbBananasPerSequence{delta, make([]*int64, len(secrets))}
					newSeq.bananasPerBuyers[b] = &newNbBananas
					sequences[delta] = newSeq
				}
			}
			previousSecret = newSecret
			previousNbBananas = newNbBananas
		}
	}

	var maxNbBananas int64
	var associatedSequence SequenceOfChanges
	for _, nbBananasForSeq := range sequences {
		nbBananas := nbBananasForSeq.bananas()
		if nbBananas > maxNbBananas {
			maxNbBananas = nbBananas
			associatedSequence = nbBananasForSeq.delta
		}
	}

	fmt.Printf("Max number of bananas: %d\n", maxNbBananas)
	fmt.Printf("Associated sequence: %v\n", associatedSequence)

	return maxNbBananas
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
