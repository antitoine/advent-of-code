package main

import (
	"bufio"
	"io"
	"log"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Machine struct {
	numLights int
	target    uint64   // bitmask: bit i = 1 if light i should be ON
	buttons   []uint64 // each button is a bitmask of lights it toggles
}

var buttonRegex = regexp.MustCompile(`\(([0-9,]+)\)`)

func parseMachine(line string) Machine {
	var m Machine

	// Extract indicator light diagram [...]
	startBracket := strings.Index(line, "[")
	endBracket := strings.Index(line, "]")
	diagram := line[startBracket+1 : endBracket]

	m.numLights = len(diagram)
	m.target = 0
	for i, c := range diagram {
		if c == '#' {
			m.target |= 1 << i
		}
	}

	// Extract button wiring schematics (...)
	matches := buttonRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		var buttonMask uint64 = 0
		parts := strings.Split(match[1], ",")
		for _, p := range parts {
			lightIdx, _ := strconv.Atoi(p)
			buttonMask |= 1 << lightIdx
		}
		m.buttons = append(m.buttons, buttonMask)
	}

	return m
}

func minButtonPresses(m Machine) int {
	numButtons := len(m.buttons)
	minPresses := numButtons + 1 // worst case: more than all buttons

	// Try all 2^n combinations of button presses
	for combo := 0; combo < (1 << numButtons); combo++ {
		var state uint64 = 0 // all lights start OFF
		presses := 0

		for i := 0; i < numButtons; i++ {
			if (combo>>i)&1 == 1 {
				state ^= m.buttons[i]
				presses++
			}
		}

		if state == m.target && presses < minPresses {
			minPresses = presses
		}
	}

	if minPresses > numButtons {
		return -1 // no solution found (shouldn't happen for valid input)
	}
	return minPresses
}

func minButtonPressesBFS(m Machine) int {
	numButtons := len(m.buttons)
	if numButtons > 20 {
		return minButtonPresses(m) // fallback
	}

	// BFS to find minimum presses
	// State: current light configuration
	// We explore by number of button presses
	visited := make(map[uint64]bool)
	visited[0] = true

	if m.target == 0 {
		return 0
	}

	current := []uint64{0}
	presses := 0

	for len(current) > 0 {
		presses++
		next := make([]uint64, 0)
		seen := make(map[uint64]bool)

		for _, state := range current {
			for _, button := range m.buttons {
				newState := state ^ button
				if newState == m.target {
					return presses
				}
				if !visited[newState] && !seen[newState] {
					seen[newState] = true
					next = append(next, newState)
				}
			}
		}

		for s := range seen {
			visited[s] = true
		}
		current = next

		// Optimization: if we've explored more than numButtons presses,
		// the brute force approach is better
		if presses > numButtons {
			break
		}
	}

	// Fallback to brute force if BFS doesn't find quickly
	return minButtonPresses(m)
}

func minButtonPressesOptimized(m Machine) int {
	numButtons := len(m.buttons)

	// For small number of buttons, use brute force but iterate by hamming weight
	// This finds minimum presses faster
	for numPresses := 0; numPresses <= numButtons; numPresses++ {
		// Try all combinations with exactly numPresses buttons pressed
		for combo := 0; combo < (1 << numButtons); combo++ {
			if bits.OnesCount(uint(combo)) != numPresses {
				continue
			}

			var state uint64 = 0
			for i := 0; i < numButtons; i++ {
				if (combo>>i)&1 == 1 {
					state ^= m.buttons[i]
				}
			}

			if state == m.target {
				return numPresses
			}
		}
	}

	return -1 // no solution found
}

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	var totalPresses int64 = 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		machine := parseMachine(line)
		presses := minButtonPressesOptimized(machine)
		if presses >= 0 {
			totalPresses += int64(presses)
		}
	}

	return totalPresses
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
