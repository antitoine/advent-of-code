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

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return 0
	}

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	paddedLines := make([]string, len(lines))
	for i, line := range lines {
		paddedLines[i] = line + strings.Repeat(" ", maxWidth-len(line))
	}

	isSeparator := make([]bool, maxWidth)
	for col := 0; col < maxWidth; col++ {
		allSpace := true
		for row := 0; row < len(paddedLines); row++ {
			if paddedLines[row][col] != ' ' {
				allSpace = false
				break
			}
		}
		isSeparator[col] = allSpace
	}

	var problemRanges [][2]int
	start := -1
	for col := 0; col < maxWidth; col++ {
		if isSeparator[col] {
			if start != -1 {
				problemRanges = append(problemRanges, [2]int{start, col})
				start = -1
			}
		} else {
			if start == -1 {
				start = col
			}
		}
	}
	if start != -1 {
		problemRanges = append(problemRanges, [2]int{start, maxWidth})
	}

	var total int64 = 0
	operatorRow := len(paddedLines) - 1

	for _, pr := range problemRanges {
		startCol, endCol := pr[0], pr[1]

		var numbers []int64
		var op rune = 0

		for row := 0; row < len(paddedLines); row++ {
			substr := paddedLines[row][startCol:endCol]
			trimmed := strings.TrimSpace(substr)

			if row == operatorRow {
				if len(trimmed) > 0 {
					op = rune(trimmed[0])
				}
			} else {
				if trimmed != "" {
					num, err := strconv.ParseInt(trimmed, 10, 64)
					if err == nil {
						numbers = append(numbers, num)
					}
				}
			}
		}

		if len(numbers) > 0 && (op == '*' || op == '+') {
			var result int64
			if op == '*' {
				result = 1
				for _, n := range numbers {
					result *= n
				}
			} else {
				result = 0
				for _, n := range numbers {
					result += n
				}
			}
			total += result
		}
	}

	return total
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
