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
	const testingExpectedResult = 16
	result := getResultPart1(strings.NewReader(testingInput), 6)
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
	}
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
