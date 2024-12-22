package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Code rune

const (
	C0 Code = '0'
	C1 Code = '1'
	C2 Code = '2'
	C3 Code = '3'
	C4 Code = '4'
	C5 Code = '5'
	C6 Code = '6'
	C7 Code = '7'
	C8 Code = '8'
	C9 Code = '9'
	CA Code = 'A'
)

type NumericalMap map[Code]image.Point

var numericalMap = NumericalMap{
	CA: image.Pt(2, 0),
	C0: image.Pt(1, 0),
	C1: image.Pt(0, 1),
	C2: image.Pt(1, 1),
	C3: image.Pt(2, 1),
	C4: image.Pt(0, 2),
	C5: image.Pt(1, 2),
	C6: image.Pt(2, 2),
	C7: image.Pt(0, 3),
	C8: image.Pt(1, 3),
	C9: image.Pt(2, 3),
}

type Direction rune

const (
	DU Direction = '^'
	DD Direction = 'v'
	DL Direction = '<'
	DR Direction = '>'
	DA Direction = 'A'
)

type DirectionalMap map[Direction]image.Point

var directionalMap = DirectionalMap{
	DA: image.Pt(2, 1),
	DU: image.Pt(1, 1),
	DL: image.Pt(0, 0),
	DD: image.Pt(1, 0),
	DR: image.Pt(2, 0),
}

func parseLine(line string) []Code {
	var codes []Code
	for _, c := range line {
		codes = append(codes, Code(c))
	}
	return codes
}

func parseInput(input io.Reader) [][]Code {
	scanner := bufio.NewScanner(input)

	var codes [][]Code
	for scanner.Scan() {
		codes = append(codes, parseLine(scanner.Text()))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return codes
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/*
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//	   | 0 | A |
//	   +---+---+
*/
func getPressesForNumericPad(code []Code, start Code, numericalMap NumericalMap) []Direction {
	current := numericalMap[start]
	var output []Direction

	for _, c := range code {
		dest := numericalMap[c]
		diffX, diffY := dest.X-current.X, dest.Y-current.Y

		var horizontal, vertical []Direction

		for i := 0; i < abs(diffX); i++ {
			if diffX >= 0 {
				horizontal = append(horizontal, DR)
			} else {
				horizontal = append(horizontal, DL)
			}
		}

		for i := 0; i < abs(diffY); i++ {
			if diffY >= 0 {
				vertical = append(vertical, DU)
			} else {
				vertical = append(vertical, DD)
			}
		}

		if current.Y == 0 && dest.X == 0 {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		} else if current.X == 0 && dest.Y == 0 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else if diffX < 0 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		}

		current = dest
		output = append(output, DA)
	}
	return output
}

/*
//     +---+---+
//     | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+
*/
func getPressesForDirectionalPad(sequence []Direction, start Direction, directionalMap DirectionalMap) []Direction {
	current := directionalMap[start]
	var output []Direction

	for _, char := range sequence {
		dest := directionalMap[char]
		diffX, diffY := dest.X-current.X, dest.Y-current.Y

		var horizontal, vertical []Direction

		for i := 0; i < abs(diffX); i++ {
			if diffX >= 0 {
				horizontal = append(horizontal, DR)
			} else {
				horizontal = append(horizontal, DL)
			}
		}

		for i := 0; i < abs(diffY); i++ {
			if diffY >= 0 {
				vertical = append(vertical, DU)
			} else {
				vertical = append(vertical, DD)
			}
		}

		if current.X == 0 && dest.Y == 1 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else if current.Y == 1 && dest.X == 0 {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		} else if diffX < 0 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		}
		current = dest
		output = append(output, DA)
	}
	return output
}

func getCacheKey(input []Direction) string {
	return string(input)
}

func getCountAfterRobots(initSeq []Direction, maxRobots int, robot int, cache map[string][]int64, directionalMap DirectionalMap) int64 {
	cacheKey := getCacheKey(initSeq)
	if val, ok := cache[cacheKey]; ok {
		if val[robot-1] != 0 {
			return val[robot-1]
		}
	} else {
		cache[cacheKey] = make([]int64, maxRobots)
	}

	seq := getPressesForDirectionalPad(initSeq, DA, directionalMap)
	cache[cacheKey][0] = int64(len(seq))

	if robot == maxRobots {
		return int64(len(seq))
	}

	subSequences := getIndividualSteps(seq)

	var count int64
	for _, subSequence := range subSequences {
		c := getCountAfterRobots(subSequence, maxRobots, robot+1, cache, directionalMap)
		subCacheKey := getCacheKey(subSequence)
		if _, ok := cache[subCacheKey]; !ok {
			cache[subCacheKey] = make([]int64, maxRobots)
		}
		cache[subCacheKey][0] = c
		count += c
	}

	cache[cacheKey][robot-1] = count
	return count
}

func getIndividualSteps(sequence []Direction) [][]Direction {
	var output [][]Direction
	var current []Direction
	for _, dir := range sequence {
		current = append(current, dir)
		if dir == DA {
			output = append(output, current)
			current = []Direction{}
		}
	}
	return output
}

func getDigitsFromCode(code []Code) int64 {
	digits, errParse := strconv.ParseInt(string(code[:3]), 10, 64)
	if errParse != nil {
		log.Fatalf("Unable to parse digits: %v", errParse)
	}
	return digits
}

func getSequence(codes [][]Code, numericalMap NumericalMap, directionalMap DirectionalMap, robots int) int64 {
	var count int64
	cache := make(map[string][]int64)
	for _, code := range codes {
		firstSequence := getPressesForNumericPad(code, CA, numericalMap)
		num := getCountAfterRobots(firstSequence, robots, 1, cache, directionalMap)
		count += getDigitsFromCode(code) * num
	}
	return count
}

func getResult(input io.Reader) int64 {
	codes := parseInput(input)
	return getSequence(codes, numericalMap, directionalMap, 25)
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
