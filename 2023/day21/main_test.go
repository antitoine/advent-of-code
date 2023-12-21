package main

import (
	"strings"
	"testing"
)

const testingInput = `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........
`

func TestGetResults(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		const testingExpectedResult = 16
		result := getResultPart1(strings.NewReader(testingInput), 6)
		if result != testingExpectedResult {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
		}
	})
	t.Run("part2", func(t *testing.T) {
		const testingExpectedResult = 16733044
		result := getResultPart2(strings.NewReader(testingInput), 5000)
		if result != testingExpectedResult {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
		}
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResultPart1(strings.NewReader(testingInput), 6)
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResultPart1(inputFile, 64)
			inputFile.Close()
		}
	})
}
