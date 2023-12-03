package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"
)

const matrixLength = 140

type LinkedNumber struct {
	number   int64
	included bool
}

type Cell struct {
	linkedNumber *LinkedNumber
}

type Row = []Cell

type Coords struct {
	rowIdx  int
	cellIdx int
}

type Matrix struct {
	rows            []Row
	detectedSymbols []Coords
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
	linkedNumber := &LinkedNumber{
		number: *nb,
	}
	for _, cellIdx := range d.affectedCellsIdx {
		row[cellIdx].linkedNumber = linkedNumber
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
					matrix.detectedSymbols = append(matrix.detectedSymbols, Coords{rowIdx, cellIdx})
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

func (cell Cell) getNumberToSum() int64 {
	if cell.linkedNumber != nil && !cell.linkedNumber.included {
		cell.linkedNumber.included = true
		return cell.linkedNumber.number
	} else {
		return 0
	}
}

func getSumOfNumbersAttachedToSymbols(file *os.File) int64 {
	var finalSum int64
	matrix := parseMatrix(file)
	for _, coords := range matrix.detectedSymbols {
		if coords.rowIdx-1 >= 0 {
			if coords.cellIdx-1 >= 0 {
				finalSum += matrix.rows[coords.rowIdx-1][coords.cellIdx-1].getNumberToSum()
			}
			finalSum += matrix.rows[coords.rowIdx-1][coords.cellIdx].getNumberToSum()
			if coords.cellIdx+1 < matrixLength {
				finalSum += matrix.rows[coords.rowIdx-1][coords.cellIdx+1].getNumberToSum()
			}
		}
		if coords.cellIdx-1 >= 0 {
			finalSum += matrix.rows[coords.rowIdx][coords.cellIdx-1].getNumberToSum()
		}
		if coords.cellIdx+1 < matrixLength {
			finalSum += matrix.rows[coords.rowIdx][coords.cellIdx+1].getNumberToSum()
		}
		if coords.rowIdx+1 < matrixLength {
			if coords.cellIdx-1 >= 0 {
				finalSum += matrix.rows[coords.rowIdx+1][coords.cellIdx-1].getNumberToSum()
			}
			finalSum += matrix.rows[coords.rowIdx+1][coords.cellIdx].getNumberToSum()
			if coords.cellIdx+1 < matrixLength {
				finalSum += matrix.rows[coords.rowIdx+1][coords.cellIdx+1].getNumberToSum()
			}
		}
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
