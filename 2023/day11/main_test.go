package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
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

		result := getResult(strings.NewReader(input), 2.0)
		if result != 374 {
			t.Errorf("Expected result to be 374, got %d", result)
		}
	})
	t.Run("part2", func(t *testing.T) {
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

		result := getResult(strings.NewReader(input), 100)
		if result != 8410 {
			t.Errorf("Expected result to be 8410, got %d", result)
		}
	})
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
			getResult(strings.NewReader(input), 2.0)
		}
	})

	b.Run("large", func(b *testing.B) {
		inputFile := loadFile()
		defer inputFile.Close()

		for n := 0; n < b.N; n++ {
			getResult(inputFile, 2.0)
		}
	})
}
