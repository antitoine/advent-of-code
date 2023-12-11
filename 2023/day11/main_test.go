package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	input := `...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

	result := getResult(strings.NewReader(input))
	if result != 374 {
		t.Errorf("Expected result to be 374, got %d", result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		input := `...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(input))
		}
	})

	b.Run("large", func(b *testing.B) {
		inputFile := loadFile()
		defer inputFile.Close()

		for n := 0; n < b.N; n++ {
			getResult(inputFile)
		}
	})
}
