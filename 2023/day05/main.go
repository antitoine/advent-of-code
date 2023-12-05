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

type Range struct {
	from int64
	len  int64
}

type RangeMap struct {
	from int64
	to   int64
	len  int64
}

type Map struct {
	from    string
	to      string
	mapping []RangeMap
	nextMap *Map
}

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
	return numbers
}

var seedsLineRegex = regexp.MustCompile(`^seeds: (.*)$`)

func parseSeeds(line string) []Range {
	results := seedsLineRegex.FindStringSubmatch(line)
	if len(results) != 2 {
		log.Fatalf("Unable to parse seeds line '%s', getting results length of %d", line, len(results))
	}
	numbers := numbersStrToInts(results[1])
	if len(numbers)%2 != 0 {
		log.Fatalf("Unable to parse seeds line '%s', getting unexpected length of %d after parsing", line, len(numbers))
	}
	var seeds []Range
	for i := 0; i < len(numbers); i += 2 {
		seeds = append(seeds, Range{
			from: numbers[i],
			len:  numbers[i+1],
		})
	}
	return seeds
}

var mappingLineRegex = regexp.MustCompile(`^([^-]*)-to-([^ ]*) map:$`)

func parseMapping(line string) (string, string) {
	results := mappingLineRegex.FindStringSubmatch(line)
	if len(results) != 3 {
		log.Fatalf("Unable to parse mapping line '%s', getting results length of %d", line, len(results))
	}
	return results[1], results[2]
}

func parseInput(input io.Reader) ([]Range, *Map) {
	scanner := bufio.NewScanner(input)

	// Seeds
	if !scanner.Scan() {
		log.Fatalf("Unable to scan the first line of the input correctly: %v", scanner.Err())
	}
	seeds := parseSeeds(scanner.Text())

	scanner.Scan() // Empty line

	// Mapping
	var firstMap *Map
	var lastMap *Map
	for scanner.Scan() {
		from, to := parseMapping(scanner.Text())
		newMap := &Map{
			from: from,
			to:   to,
		}
		if firstMap == nil {
			firstMap = newMap
		}
		if lastMap != nil {
			lastMap.nextMap = newMap
		}
		for scanner.Scan() {
			mappingRangeStr := scanner.Text()
			if mappingRangeStr == "" {
				break
			}
			mappingRange := numbersStrToInts(mappingRangeStr)
			if len(mappingRange) != 3 {
				log.Fatalf("Unable to parse mapping range '%s', getting results length of %d", mappingRangeStr, len(mappingRange))
			}
			newRangeMap := RangeMap{
				from: mappingRange[1],
				to:   mappingRange[0],
				len:  mappingRange[2],
			}
			newMap.mapping = append(newMap.mapping, newRangeMap)
		}
		sort.Slice(newMap.mapping, func(i, j int) bool { return newMap.mapping[i].from < newMap.mapping[j].from })
		lastMap = newMap
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input correctly: %v", errScanningFile)
	}

	return seeds, firstMap
}

func getLowestLocation(input io.Reader) int64 {
	seeds, firstMap := parseInput(input)

	lowestLocation := int64(-1)
	for _, seedRange := range seeds {
		for seed := seedRange.from; seed < seedRange.from+seedRange.len; seed++ {
			currentMap := firstMap
			currentIdx := seed
			for {
				for _, mapping := range currentMap.mapping {
					if currentIdx >= mapping.from && currentIdx < mapping.from+mapping.len {
						currentIdx = (currentIdx - mapping.from) + mapping.to
						break
					}
				}
				if currentMap.nextMap != nil {
					currentMap = currentMap.nextMap
				} else {
					break
				}
			}
			if lowestLocation == -1 || currentIdx < lowestLocation {
				lowestLocation = currentIdx
			}
		}
	}

	return lowestLocation
}

func main() {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	lowestLocation := getLowestLocation(inputFile)

	log.Printf("Final lowest location: %d", lowestLocation)
}
