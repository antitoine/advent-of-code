package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type isSpringDamaged = *bool

type Line struct {
	isSpringsDamaged       []isSpringDamaged
	damagedSpringsCounters []int64
}

func parseInput(input io.Reader) []Line {
	scanner := bufio.NewScanner(input)
	var lines []Line
	for scanner.Scan() {
		line := scanner.Text()
		lineSplit := strings.Split(line, " ")
		if len(lineSplit) != 2 {
			log.Fatalf("Invalid input line: %s", line)
		}

		states := strings.Split(lineSplit[0], "")
		initSprings := make([]isSpringDamaged, len(states))
		for i, state := range states {
			if state == "?" {
				initSprings[i] = nil
			} else {
				isDamaged := state == "#"
				initSprings[i] = &isDamaged
			}
		}
		springs := make([]isSpringDamaged, 4+(len(initSprings)*5))
		springs = initSprings
		for i := 0; i < 4; i++ {
			springs = append(append(springs, nil), initSprings...)
		}

		numbers := strings.Split(lineSplit[1], ",")
		initCounters := make([]int64, len(numbers))
		for i, number := range numbers {
			counter, errParsing := strconv.ParseInt(number, 10, 64)
			if errParsing != nil {
				log.Fatalf("Unable to parse input line provided '%s' for counter number '%s': %v", line, number, errParsing)
			}
			initCounters[i] = counter
		}
		counters := make([]int64, len(initSprings)*5)
		counters = initCounters
		for i := 0; i < 4; i++ {
			counters = append(counters, initCounters...)
		}

		lineStruct := Line{
			isSpringsDamaged:       springs,
			damagedSpringsCounters: counters,
		}

		lines = append(lines, lineStruct)
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return lines
}

func (l Line) String() string {
	var result strings.Builder
	for _, spring := range l.isSpringsDamaged {
		if spring == nil {
			result.WriteString("?")
		} else if *spring {
			result.WriteString("#")
		} else {
			result.WriteString(".")
		}
	}
	result.WriteString(" ")
	var numbersStr []string
	for _, counter := range l.damagedSpringsCounters {
		numbersStr = append(numbersStr, strconv.FormatInt(counter, 10))
	}
	result.WriteString(strings.Join(numbersStr, ","))
	return result.String()
}

var stateToArrangements = make(map[string]int64)

func getNumberOfArrangements(line Line) int64 {
	if arrangements, found := stateToArrangements[line.String()]; found {
		return arrangements
	}

	isSpringsDamaged, damagedSpringsCounters := line.isSpringsDamaged, line.damagedSpringsCounters

	//log.Printf("Compute arrangements for line: %s", line.String())

	if len(isSpringsDamaged) == 0 {
		stateToArrangements[line.String()] = 0
		return 0
	}
	var i int
	for ; i < len(isSpringsDamaged) && isSpringsDamaged[i] != nil && *isSpringsDamaged[i] == false; i++ {
	}
	var numberOfArrangements int64
	if i < len(isSpringsDamaged) && isSpringsDamaged[i] == nil {
		numberOfArrangements += getNumberOfArrangements(Line{isSpringsDamaged[i+1:], damagedSpringsCounters})
	}
	damagedSpringsCounter := damagedSpringsCounters[0]
	for ; i < len(isSpringsDamaged) && damagedSpringsCounter > 0; i++ {
		if isSpringsDamaged[i] == nil || *isSpringsDamaged[i] == true {
			damagedSpringsCounter--
		} else {
			break
		}
	}
	if damagedSpringsCounter == 0 {
		if i == len(isSpringsDamaged) {
			if len(damagedSpringsCounters) == 1 {
				numberOfArrangements++
			}
		} else if isSpringsDamaged[i] == nil || *isSpringsDamaged[i] == false {
			if len(damagedSpringsCounters) == 1 {
				for ; i < len(isSpringsDamaged) && (isSpringsDamaged[i] == nil || *isSpringsDamaged[i] == false); i++ {
				}
				if i == len(isSpringsDamaged) {
					numberOfArrangements++
				}
			} else if len(damagedSpringsCounters) > 1 {
				nextArrangements := getNumberOfArrangements(Line{isSpringsDamaged[i+1:], damagedSpringsCounters[1:]})
				if nextArrangements > 0 {
					numberOfArrangements += nextArrangements
				}
			}
		}
	}
	//log.Printf("Get %d arrangements for line: %s", numberOfArrangements, line.String())
	stateToArrangements[line.String()] = numberOfArrangements
	return numberOfArrangements
}

func getResult(input io.Reader) int64 {
	lines := parseInput(input)
	var result int64
	for _, line := range lines {
		result += getNumberOfArrangements(line)
		//log.Printf("------------------")
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
