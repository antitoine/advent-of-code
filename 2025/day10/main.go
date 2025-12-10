package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Machine struct {
	numCounters int
	targets     []int
	buttons     [][]int
}

var buttonRegex = regexp.MustCompile(`\(([0-9,]+)\)`)
var joltageRegex = regexp.MustCompile(`\{([0-9,]+)\}`)

func parseMachine(line string) Machine {
	var m Machine

	joltageMatch := joltageRegex.FindStringSubmatch(line)
	if joltageMatch != nil {
		parts := strings.Split(joltageMatch[1], ",")
		m.numCounters = len(parts)
		m.targets = make([]int, len(parts))
		for i, p := range parts {
			m.targets[i], _ = strconv.Atoi(p)
		}
	}

	matches := buttonRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		var button []int
		parts := strings.Split(match[1], ",")
		for _, p := range parts {
			idx, _ := strconv.Atoi(p)
			button = append(button, idx)
		}
		m.buttons = append(m.buttons, button)
	}

	return m
}

type Rat struct{ n, d int64 }

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}

func newRat(n, d int64) Rat {
	if d == 0 {
		return Rat{0, 0}
	}
	if d < 0 {
		n, d = -n, -d
	}
	if n == 0 {
		return Rat{0, 1}
	}
	g := gcd(n, d)
	return Rat{n / g, d / g}
}

func (r Rat) sub(s Rat) Rat { return newRat(r.n*s.d-s.n*r.d, r.d*s.d) }
func (r Rat) mul(s Rat) Rat { return newRat(r.n*s.n, r.d*s.d) }
func (r Rat) div(s Rat) Rat {
	if s.n == 0 {
		return Rat{0, 0}
	}
	return newRat(r.n*s.d, r.d*s.n)
}

func minButtonPresses(m Machine) int {
	numButtons := len(m.buttons)
	numCounters := m.numCounters

	if numButtons == 0 {
		for _, t := range m.targets {
			if t != 0 {
				return -1
			}
		}
		return 0
	}

	// Build augmented matrix [A | t]
	matrix := make([][]Rat, numCounters)
	for i := 0; i < numCounters; i++ {
		matrix[i] = make([]Rat, numButtons+1)
		for j := 0; j <= numButtons; j++ {
			matrix[i][j] = Rat{0, 1}
		}
		matrix[i][numButtons] = Rat{int64(m.targets[i]), 1}
	}

	for j, button := range m.buttons {
		for _, idx := range button {
			matrix[idx][j] = Rat{1, 1}
		}
	}

	// Gaussian elimination to RREF
	pivotCol := make([]int, numCounters)
	for i := range pivotCol {
		pivotCol[i] = -1
	}

	row := 0
	for col := 0; col < numButtons && row < numCounters; col++ {
		pivotRow := -1
		for i := row; i < numCounters; i++ {
			if matrix[i][col].n != 0 {
				pivotRow = i
				break
			}
		}
		if pivotRow == -1 {
			continue
		}

		matrix[row], matrix[pivotRow] = matrix[pivotRow], matrix[row]
		pivotCol[row] = col

		pivot := matrix[row][col]
		for j := col; j <= numButtons; j++ {
			matrix[row][j] = matrix[row][j].div(pivot)
		}

		for i := 0; i < numCounters; i++ {
			if i == row || matrix[i][col].n == 0 {
				continue
			}
			factor := matrix[i][col]
			for j := col; j <= numButtons; j++ {
				matrix[i][j] = matrix[i][j].sub(factor.mul(matrix[row][j]))
			}
		}
		row++
	}

	rank := row

	// Check for inconsistency
	for i := rank; i < numCounters; i++ {
		if matrix[i][numButtons].n != 0 {
			return -1
		}
	}

	// Identify pivot and free variables
	isPivot := make([]bool, numButtons)
	pivotRowOf := make([]int, numButtons)
	for i := 0; i < rank; i++ {
		if pivotCol[i] >= 0 {
			isPivot[pivotCol[i]] = true
			pivotRowOf[pivotCol[i]] = i
		}
	}

	freeVars := []int{}
	for j := 0; j < numButtons; j++ {
		if !isPivot[j] {
			freeVars = append(freeVars, j)
		}
	}

	// Compute sum of all targets for upper bound
	sumTargets := 0
	for _, t := range m.targets {
		sumTargets += t
	}

	// If no free variables, unique solution
	if len(freeVars) == 0 {
		sum := 0
		for j := 0; j < numButtons; j++ {
			if isPivot[j] {
				r := matrix[pivotRowOf[j]][numButtons]
				if r.d != 1 || r.n < 0 {
					return -1
				}
				sum += int(r.n)
			}
		}
		return sum
	}

	// With free variables, enumerate
	bestSum := int(^uint(0) >> 1)
	freeVals := make([]int, len(freeVars))

	// Helper to evaluate solution for given free variable values
	evalSolution := func() (int, bool) {
		total := 0
		for i := range freeVars {
			total += freeVals[i]
		}
		for j := 0; j < numButtons; j++ {
			if isPivot[j] {
				r := pivotRowOf[j]
				val := matrix[r][numButtons]
				for i, fv := range freeVars {
					val = val.sub(matrix[r][fv].mul(Rat{int64(freeVals[i]), 1}))
				}
				if val.d != 1 || val.n < 0 {
					return 0, false
				}
				total += int(val.n)
			}
		}
		return total, true
	}

	maxFreeVal := sumTargets
	if maxFreeVal > 500 {
		maxFreeVal = 500
	}

	var enumerate func(idx int)
	enumerate = func(idx int) {
		if idx == len(freeVars) {
			if total, valid := evalSolution(); valid && total < bestSum {
				bestSum = total
			}
			return
		}

		for v := 0; v <= maxFreeVal; v++ {
			freeVals[idx] = v
			enumerate(idx + 1)
		}
	}

	enumerate(0)

	if bestSum == int(^uint(0)>>1) {
		return -1
	}
	return bestSum
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
		presses := minButtonPresses(machine)
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
