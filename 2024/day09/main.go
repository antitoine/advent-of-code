package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

func parseInput(input io.Reader) []int {
	scanner := bufio.NewScanner(input)

	if !scanner.Scan() {
		log.Fatalf("Unable to scan the input file correctly")
	}

	line := scanner.Text()
	result := make([]int, 0)
	isBlocks := true
	blockNumber := 0
	for _, numStr := range line {
		num, errParsing := strconv.Atoi(string(numStr))
		if errParsing != nil {
			log.Fatalf("Unable to parse the number '%s': %v", string(numStr), errParsing)
		}
		for i := 0; i < num; i++ {
			if isBlocks {
				result = append(result, blockNumber)
			} else {
				result = append(result, -1)
			}
		}
		if isBlocks {
			blockNumber++
		}
		isBlocks = !isBlocks
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return result
}

func countBlockSize(disk []int, j int) int {
	blockSize := 0
	for k := j; k >= 0 && disk[k] == disk[j]; k-- {
		blockSize++
	}
	return blockSize
}

func sortFreeSpaces(freeSpaces [][]int) {
	sort.Slice(freeSpaces, func(i, j int) bool {
		if freeSpaces[i][0] == freeSpaces[j][0] {
			return freeSpaces[i][1] < freeSpaces[j][1]
		}
		return freeSpaces[i][0] < freeSpaces[j][0]
	})
}

func getResult(input io.Reader) int64 {
	initialState := parseInput(input)

	var freeSpaces [][]int
	for i := 0; i < len(initialState); i++ {
		if initialState[i] != -1 {
			continue
		}
		startIdx := i
		for i < len(initialState) && initialState[i] == -1 {
			i++
		}
		freeSpaces = append(freeSpaces, []int{i - startIdx, startIdx})
	}
	sortFreeSpaces(freeSpaces)

	j := len(initialState) - 1
	for j >= 0 {
		if initialState[j] == -1 {
			j--
			continue
		}
		blockSize := countBlockSize(initialState, j)
		freeSpaceIdx := sort.Search(len(freeSpaces), func(f int) bool {
			return freeSpaces[f][0] >= blockSize
		})
		if freeSpaceIdx == len(freeSpaces) || freeSpaces[freeSpaceIdx][1] > j {
			j -= blockSize
		} else {
			startIdx := freeSpaces[freeSpaceIdx][1]
			for k := 0; k < blockSize; k++ {
				initialState[startIdx+k] = initialState[j]
				initialState[j] = -1
				j--
				freeSpaces[freeSpaceIdx][0]--
				freeSpaces[freeSpaceIdx][1]++
			}
			if freeSpaces[freeSpaceIdx][0] == 0 {
				freeSpaces = append(freeSpaces[:freeSpaceIdx], freeSpaces[freeSpaceIdx+1:]...)
			}
			sortFreeSpaces(freeSpaces)
		}
	}

	var checksum int64
	for i := 0; i < len(initialState); i++ {
		if initialState[i] != -1 {
			checksum += int64(initialState[i]) * int64(i)
		}
	}

	return checksum
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
