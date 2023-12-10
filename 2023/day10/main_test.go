package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	t.Run("part 1", func(t *testing.T) {
		t.Run("example 1", func(t *testing.T) {
			input := `.....
.S-7.
.|.|.
.L-J.
.....
`

			result := getResultPart1(strings.NewReader(input))
			if result != 4 {
				t.Errorf("Expected result to be 4, got %d", result)
			}
		})
		t.Run("example 2", func(t *testing.T) {
			input := `..F7.
.FJ|.
SJ.L7
|F--J
LJ...
`

			result := getResultPart1(strings.NewReader(input))
			if result != 8 {
				t.Errorf("Expected result to be 8, got %d", result)
			}
		})
	})

	t.Run("part 2", func(t *testing.T) {
		t.Run("example 1", func(t *testing.T) {
			input := `...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........
`

			result := getResultPart2(strings.NewReader(input))
			if result != 4 {
				t.Errorf("Expected result to be 4, got %d", result)
			}
		})
		t.Run("example 2", func(t *testing.T) {
			input := `.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...
`

			result := getResultPart2(strings.NewReader(input))
			if result != 8 {
				t.Errorf("Expected result to be 8, got %d", result)
			}
		})
		t.Run("example 3", func(t *testing.T) {
			input := `FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L
`

			result := getResultPart2(strings.NewReader(input))
			if result != 10 {
				t.Errorf("Expected result to be 10, got %d", result)
			}
		})
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("part 1", func(b *testing.B) {
		b.Run("small", func(b *testing.B) {
			input := `..F7.
.FJ|.
SJ.L7
|F--J
LJ...
`

			for n := 0; n < b.N; n++ {
				getResultPart1(strings.NewReader(input))
			}
		})

		b.Run("large", func(b *testing.B) {
			inputFile := loadFile()
			defer inputFile.Close()

			for n := 0; n < b.N; n++ {
				getResultPart1(inputFile)
			}
		})
	})

	b.Run("part 2", func(b *testing.B) {
		b.Run("small", func(b *testing.B) {
			input := `.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...
`

			for n := 0; n < b.N; n++ {
				getResultPart2(strings.NewReader(input))
			}
		})

		b.Run("large", func(b *testing.B) {
			inputFile := loadFile()
			defer inputFile.Close()

			for n := 0; n < b.N; n++ {
				getResultPart2(inputFile)
			}
		})
	})
}
