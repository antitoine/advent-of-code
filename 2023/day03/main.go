package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"
)

const matrixLength = 140

type Cell struct {
	linkedNumber *int64
}

type Row = []Cell

type Coords struct {
	rowIdx  int
	cellIdx int
}

type Matrix struct {
	rows          []Row
	detectedGears []Coords
}

type DetectedNumber struct {
	numberStr        string
	affectedCellsIdx []int
}

func (d *DetectedNumber) getNumber() *int64 {
	if d.numberStr == "" {
		return nil
	}
	number, errParsingNumber := strconv.ParseInt(d.numberStr, 10, 64)
	if errParsingNumber != nil {
		log.Fatalf("Unable to parse number '%s': %v", d.numberStr, errParsingNumber)
	}
	return &number
}

func (d *DetectedNumber) affectedNumber(row Row) {
	nb := d.getNumber()
	if nb == nil {
		return
	}
	for _, cellIdx := range d.affectedCellsIdx {
		row[cellIdx].linkedNumber = nb
	}
	d.numberStr = ""
	d.affectedCellsIdx = nil
}

func parseMatrix(file *os.File) Matrix {
	scanner := bufio.NewScanner(file)

	matrix := Matrix{
		rows: make([]Row, matrixLength),
	}
	for rowIdx := 0; scanner.Scan(); rowIdx++ {
		line := scanner.Text()
		row := make(Row, matrixLength)
		detectedNumber := &DetectedNumber{}
		for cellIdx, char := range line {
			var cell Cell
			if char != '.' {
				if char >= '0' && char <= '9' {
					detectedNumber.numberStr += string(char)
					detectedNumber.affectedCellsIdx = append(detectedNumber.affectedCellsIdx, cellIdx)
				} else {
					detectedNumber.affectedNumber(row)
				}
				if char == '*' {
					matrix.detectedGears = append(matrix.detectedGears, Coords{rowIdx, cellIdx})
				}
			} else {
				detectedNumber.affectedNumber(row)
			}
			row[cellIdx] = cell
		}
		detectedNumber.affectedNumber(row)
		matrix.rows[rowIdx] = row
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return matrix
}

type Ratio struct {
	linkedNumbers map[*int64]struct{}
	keptNumbers   []int64
}

func (ratio *Ratio) attachCell(cell Cell) {
	nb := cell.linkedNumber
	if cell.linkedNumber != nil {
		if _, ok := ratio.linkedNumbers[nb]; !ok {
			ratio.linkedNumbers[nb] = struct{}{}
			ratio.keptNumbers = append(ratio.keptNumbers, *nb)
		}
	}
}

func (ratio *Ratio) getNumberToSum() int64 {
	if len(ratio.keptNumbers) == 2 {
		return ratio.keptNumbers[0] * ratio.keptNumbers[1]
	}
	return 0
}

func getSumOfNumbersAttachedToSymbols(file *os.File) int64 {
	var finalSum int64
	matrix := parseMatrix(file)
	for _, coords := range matrix.detectedGears {
		ratio := &Ratio{
			linkedNumbers: make(map[*int64]struct{}),
		}
		if coords.rowIdx-1 >= 0 {
			if coords.cellIdx-1 >= 0 {
				ratio.attachCell(matrix.rows[coords.rowIdx-1][coords.cellIdx-1])
			}
			ratio.attachCell(matrix.rows[coords.rowIdx-1][coords.cellIdx])
			if coords.cellIdx+1 < matrixLength {
				ratio.attachCell(matrix.rows[coords.rowIdx-1][coords.cellIdx+1])
			}
		}
		if coords.cellIdx-1 >= 0 {
			ratio.attachCell(matrix.rows[coords.rowIdx][coords.cellIdx-1])
		}
		if coords.cellIdx+1 < matrixLength {
			ratio.attachCell(matrix.rows[coords.rowIdx][coords.cellIdx+1])
		}
		if coords.rowIdx+1 < matrixLength {
			if coords.cellIdx-1 >= 0 {
				ratio.attachCell(matrix.rows[coords.rowIdx+1][coords.cellIdx-1])
			}
			ratio.attachCell(matrix.rows[coords.rowIdx+1][coords.cellIdx])
			if coords.cellIdx+1 < matrixLength {
				ratio.attachCell(matrix.rows[coords.rowIdx+1][coords.cellIdx+1])
			}
		}
		finalSum += ratio.getNumberToSum()
	}

	return finalSum
}

func main() {
	startTime := time.Now()
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	finalSum := getSumOfNumbersAttachedToSymbols(inputFile)

	log.Printf("Final sum: %d", finalSum)
	log.Printf("Execution duration: %s", time.Since(startTime))
}
