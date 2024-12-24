package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Operation string

const (
	And Operation = "AND"
	Or  Operation = "OR"
	Xor Operation = "XOR"
)

type Gate struct {
	wireA string
	op    Operation
	wireB string
}

type Wire struct {
	name  string
	value *bool
	gate  *Gate
}

// Line like "y01: 1"
func parseInitValue(line string) Wire {
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		log.Fatalf("Expected 2 parts, got %d for line %s", len(parts), line)
	}
	value := parts[1] == "1"
	return Wire{
		name:  parts[0],
		value: &value,
	}
}

// Line like "ntg XOR fgs -> mjb"
func parseConnection(line string) Wire {
	parts := strings.Split(line, " -> ")
	if len(parts) != 2 {
		log.Fatalf("Expected 2 parts, got %d for line %s", len(parts), line)
	}
	operationParts := strings.Split(parts[0], " ")
	if len(operationParts) != 3 {
		log.Fatalf("Expected 3 parts, got %d for line %s", len(operationParts), parts[0])
	}
	return Wire{
		name: parts[1],
		gate: &Gate{
			wireA: operationParts[0],
			op:    Operation(operationParts[1]),
			wireB: operationParts[2],
		},
	}
}

type System struct {
	wires      map[string]Wire
	finalWires []string
}

func (s System) getFinalValue() (int64, bool) {
	var binaryStr string
	for _, wireName := range s.finalWires {
		wire := s.wires[wireName]
		if wire.value == nil {
			return 0, false
		}
		wireValue := *wire.value
		if wireValue {
			binaryStr = "1" + binaryStr
		} else {
			binaryStr = "0" + binaryStr
		}
	}
	value, err := strconv.ParseInt(binaryStr, 2, 64)
	if err != nil {
		log.Fatalf("Unable to parse binary string %s: %v", binaryStr, err)
	}
	return value, true
}

func parseInput(input io.Reader) System {
	scanner := bufio.NewScanner(input)

	system := System{
		wires:      make(map[string]Wire),
		finalWires: []string{},
	}
	initValuesForWires := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			initValuesForWires = false
			continue
		}
		var wire Wire
		if initValuesForWires {
			wire = parseInitValue(line)
		} else {
			wire = parseConnection(line)
		}
		system.wires[wire.name] = wire
		if strings.HasPrefix(wire.name, "z") {
			system.finalWires = append(system.finalWires, wire.name)
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	sort.Strings(system.finalWires)

	return system
}

func getResult(input io.Reader) int64 {
	system := parseInput(input)

	finalValue, okFinalValue := system.getFinalValue()
	for !okFinalValue {
		for wireName, wire := range system.wires {
			if wire.value != nil {
				continue
			}
			gate := wire.gate
			var valueA, valueB bool
			if wireA, ok := system.wires[gate.wireA]; !ok || wireA.value == nil {
				continue
			} else {
				valueA = *wireA.value
			}
			if wireB, ok := system.wires[gate.wireB]; !ok || wireB.value == nil {
				continue
			} else {
				valueB = *wireB.value
			}
			var result bool
			switch gate.op {
			case And:
				result = valueA && valueB
			case Or:
				result = valueA || valueB
			case Xor:
				result = valueA != valueB
			default:
				log.Fatalf("Unknown operation %s", gate.op)
			}
			system.wires[wireName] = Wire{
				name:  wireName,
				value: &result,
			}
		}
		finalValue, okFinalValue = system.getFinalValue()
	}

	return finalValue
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
